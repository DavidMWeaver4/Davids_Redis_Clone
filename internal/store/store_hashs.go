package store

import "maps"

func (s *Store) HSet(key, field, value string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		entry = Entry{
			Type: HashType,
			Hash: Hash{},
		}
	}
	if entry.Type != HashType {
		return 0, ErrWrongType
	}
	added := entry.Hash.HSet(field, value)
	s.data[key] = entry
	return added, nil
}
func (h *Hash) HSet(field, value string) int {
	if h.values == nil {
		h.values = make(map[string]string)
	}
	_, exists := h.values[field]
	h.values[field] = value
	if !exists {
		return 1
	}
	return 0
}
func (s *Store) HGet(key, field string) (string, bool, error) {
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
	if entry.Type != HashType {
		return "", false, ErrWrongType
	}
	value, found := entry.Hash.HGet(field)
	return value, found, nil
}
func (h *Hash) HGet(field string) (string, bool) {
	value, ok := h.values[field]
	return value, ok
}
func (s *Store) HDel(key string, fields ...string) (int, error) {
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
	if entry.Type != HashType {
		return 0, ErrWrongType
	}
	removed := entry.Hash.HDel(fields...)
	if len(entry.Hash.values) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = entry
	}
	return removed, nil
}
func (h *Hash) HDel(fields ...string) int {
	num := 0
	for _, field := range fields {
		_, ok := h.values[field]
		if !ok {
			continue
		}
		delete(h.values, field)
		num++
	}
	return num
}
func (s *Store) HExists(key, field string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return false, nil
	}
	if isExpired(entry) {
		delete(s.data, key)
		return false, nil
	}
	if entry.Type != HashType {
		return false, ErrWrongType
	}
	exists := entry.Hash.HExists(field)
	return exists, nil
}
func (h *Hash) HExists(field string) bool {
	_, exists := h.values[field]
	return exists
}
func (s *Store) HLen(key string) (int, error) {
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
	if entry.Type != HashType {
		return 0, ErrWrongType
	}
	length := entry.Hash.HLen()
	return length, nil
}
func (h *Hash) HLen() int {
	return len(h.values)
}
func (s *Store) HGetAll(key string) (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok {
		return map[string]string{}, nil
	}
	if isExpired(entry) {
		delete(s.data, key)
		return map[string]string{}, nil
	}
	if entry.Type != HashType {
		return map[string]string{}, ErrWrongType
	}
	values := entry.Hash.HGetAll()
	return values, nil
}
func (h *Hash) HGetAll() map[string]string {
	values := make(map[string]string, len(h.values))
	maps.Copy(values, h.values)
	return values
}
