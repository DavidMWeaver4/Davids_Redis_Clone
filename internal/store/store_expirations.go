package store

import "time"

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
