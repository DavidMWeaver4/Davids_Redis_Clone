package store

import (
	"errors"
	"math"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]Entry
}
type Entry struct {
	Value     string
	ExpiresAt time.Time
}
type KeyValue struct {
	Key   string
	Value string
}

var (
	ErrNotIntegerOrOutOfRange = errors.New("value is not an integer or out of range")
	ErrNegativeDecrement      = errors.New("decrement must be non-negative")
)

func New() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := Entry{Value: value}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}
	s.data[key] = entry
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return "", false
	}
	if isExpired(entry) {
		delete(s.data, key)
		return "", false
	}
	return entry.Value, true
}

func (s *Store) Delete(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]
	if !ok {
		return 0
	}
	if isExpired(entry) {
		delete(s.data, key)
		return 0
	}
	delete(s.data, key)
	return 1
}

func (s *Store) Exists(key string) bool {
	_, ok := s.Get(key)
	return ok
}

func (s *Store) TTL(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]
	if !ok {
		return -2
	}
	if entry.ExpiresAt.IsZero() {
		return -1
	}
	if isExpired(entry) {
		delete(s.data, key)
		return -2
	}
	return int(time.Until(entry.ExpiresAt) / time.Second)
}

func isExpired(entry Entry) bool {
	return !entry.ExpiresAt.IsZero() && !time.Now().Before(entry.ExpiresAt)
}

func (s *Store) Expire(key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return false
	}
	if isExpired(entry) {
		delete(s.data, key)
		return false
	}
	entry.ExpiresAt = time.Now().Add(ttl)
	s.data[key] = entry
	return true
}

func (s *Store) Persist(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return false
	}
	if isExpired(entry) {
		delete(s.data, key)
		return false
	}
	if entry.ExpiresAt.IsZero() {
		return false
	}
	entry.ExpiresAt = time.Time{}
	s.data[key] = entry
	return true
}

func (s *Store) Incr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{Value: "1"}
		return 1, nil
	}
	value, err := strconv.ParseInt(entry.Value, 10, 64)
	if err != nil {
		return 0, ErrNotIntegerOrOutOfRange
	}
	if value == math.MaxInt64 {
		return 0, ErrNotIntegerOrOutOfRange
	}
	value++
	entry.Value = strconv.FormatInt(value, 10)
	s.data[key] = entry
	return value, nil
}

func (s *Store) Decr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{Value: "-1"}
		return -1, nil
	}
	value, err := strconv.ParseInt(entry.Value, 10, 64)
	if err != nil {
		return 0, ErrNotIntegerOrOutOfRange
	}
	if value == math.MinInt64 {
		return 0, ErrNotIntegerOrOutOfRange
	}
	value--
	entry.Value = strconv.FormatInt(value, 10)
	s.data[key] = entry
	return value, nil
}

func (s *Store) Incrby(key string, inc int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{Value: strconv.FormatInt(inc, 10)}
		return inc, nil
	}
	value, err := strconv.ParseInt(entry.Value, 10, 64)
	if err != nil {
		return 0, ErrNotIntegerOrOutOfRange
	}
	if inc > 0 && value > math.MaxInt64-inc {
		return 0, ErrNotIntegerOrOutOfRange
	}
	if inc < 0 && value < math.MinInt64-inc {
		return 0, ErrNotIntegerOrOutOfRange
	}
	value += inc
	entry.Value = strconv.FormatInt(value, 10)
	s.data[key] = entry
	return value, nil
}

func (s *Store) Decrby(key string, dec int64) (int64, error) {
	if dec < 0 {
		return 0, ErrNegativeDecrement
	}
	return s.Incrby(key, -dec)
}

func (s *Store) Append(key, value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{Value: value}
		return len(value)
	}
	entry.Value += value
	s.data[key] = entry
	return len(entry.Value)

}

func (s *Store) Strlen(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return 0
	}
	if isExpired(entry) {
		delete(s.data, key)
		return 0
	}
	return len(entry.Value)
}

func (s *Store) Setnx(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if ok && !isExpired(entry) {
		return false
	}

	s.data[key] = Entry{Value: value}
	return true
}

func (s *Store) Mset(pairs []KeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, pair := range pairs {
		s.data[pair.Key] = Entry{Value: pair.Value}
	}

}
