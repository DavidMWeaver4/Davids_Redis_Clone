package store

import (
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
