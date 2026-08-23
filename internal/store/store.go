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
type ValueType int

const (
	InvalidType ValueType = iota
	StringType
	ListType
)

type Entry struct {
	Type      ValueType
	String    string
	List      []string
	ExpiresAt time.Time
}

type KeyValue struct {
	Key   string
	Value string
}

var (
	ErrNotIntegerOrOutOfRange = errors.New("value is not an integer or out of range")
	ErrNegativeDecrement      = errors.New("decrement must be non-negative")
	ErrWrongType              = errors.New("WRONGTYPE operation against a key holding the wrong kind of value")
)

func New() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := Entry{
		Type:   StringType,
		String: value,
	}
	if ttl > 0 {
		entry.ExpiresAt = time.Now().Add(ttl)
	}
	s.data[key] = entry
}

func (s *Store) Get(key string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return "", false, nil
	}
	if isExpired(entry) {
		delete(s.data, key)
		return "", false, nil
	}
	if entry.Type != StringType {
		return "", false, ErrWrongType
	}
	return entry.String, true, nil
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
	return true
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
	if ttl <= 0 {
		delete(s.data, key)
		return true
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
		s.data[key] = Entry{
			Type:   StringType,
			String: "1",
		}
		return 1, nil
	}
	if entry.Type != StringType {
		return 0, ErrWrongType
	}
	value, err := strconv.ParseInt(entry.String, 10, 64)
	if err != nil {
		return 0, ErrNotIntegerOrOutOfRange
	}
	if value == math.MaxInt64 {
		return 0, ErrNotIntegerOrOutOfRange
	}
	value++
	entry.String = strconv.FormatInt(value, 10)
	s.data[key] = entry
	return value, nil
}

func (s *Store) Decr(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{
			Type:   StringType,
			String: "-1",
		}
		return -1, nil
	}
	if entry.Type != StringType {
		return 0, ErrWrongType
	}
	value, err := strconv.ParseInt(entry.String, 10, 64)
	if err != nil {
		return 0, ErrNotIntegerOrOutOfRange
	}
	if value == math.MinInt64 {
		return 0, ErrNotIntegerOrOutOfRange
	}
	value--
	entry.String = strconv.FormatInt(value, 10)
	s.data[key] = entry
	return value, nil
}

func (s *Store) Incrby(key string, inc int64) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{
			Type:   StringType,
			String: strconv.FormatInt(inc, 10),
		}
		return inc, nil
	}
	if entry.Type != StringType {
		return 0, ErrWrongType
	}
	value, err := strconv.ParseInt(entry.String, 10, 64)
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
	entry.String = strconv.FormatInt(value, 10)
	s.data[key] = entry
	return value, nil
}

func (s *Store) Decrby(key string, dec int64) (int64, error) {
	if dec < 0 {
		return 0, ErrNegativeDecrement
	}
	return s.Incrby(key, -dec)
}

func (s *Store) Append(key, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if ok && isExpired(entry) {
		delete(s.data, key)
		ok = false
	}
	if !ok {
		s.data[key] = Entry{
			Type:   StringType,
			String: value,
		}
		return len(value), nil
	}
	if entry.Type != StringType {
		return 0, ErrWrongType
	}
	entry.String += value
	s.data[key] = entry
	return len(entry.String), nil
}

func (s *Store) Strlen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return 0, nil
	}
	if isExpired(entry) {
		delete(s.data, key)
		return 0, nil
	}
	if entry.Type != StringType {
		return 0, ErrWrongType
	}
	return len(entry.String), nil
}

func (s *Store) Setnx(key, value string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		s.data[key] = Entry{
			Type:   StringType,
			String: value,
		}
		return true, nil
	}
	if entry.Type != StringType {
		return false, ErrWrongType
	}

	return false, nil
}

func (s *Store) Mset(pairs []KeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, pair := range pairs {
		s.data[pair.Key] = Entry{
			Type:   StringType,
			String: pair.Value,
		}
	}

}
