package store

import "time"

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
