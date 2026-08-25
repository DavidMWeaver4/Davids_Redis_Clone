package store

type popFunc func(*List) (string, bool)

func (s *Store) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		entry = Entry{
			Type: ListType,
		}
	}
	if entry.Type != ListType {
		return 0, ErrWrongType
	}
	length := entry.List.LPush(values...)
	s.data[key] = entry
	return length, nil
}

func (l *List) LPush(values ...string) int {
	size := len(values) + len(l.values)
	var newList = make([]string, size)
	for i := range values {
		newList[i] = values[len(values)-1-i]
	}
	copy(newList[len(values):], l.values)
	l.values = newList
	return len(newList)
}

func (s *Store) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		entry = Entry{
			Type: ListType,
		}
	}
	if entry.Type != ListType {
		return 0, ErrWrongType
	}
	length := entry.List.RPush(values...)
	s.data[key] = entry
	return length, nil
}
func (l *List) RPush(values ...string) int {
	l.values = append(l.values, values...)
	return len(l.values)
}
func (s *Store) LPop(key string) (string, bool, error) {
	return s.pop(key, (*List).LPop)
}

func (s *Store) RPop(key string) (string, bool, error) {
	return s.pop(key, (*List).RPop)
}
func (s *Store) pop(key string, popFromList popFunc) (string, bool, error) {
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
	if entry.Type != ListType {
		return "", false, ErrWrongType
	}

	value, found := popFromList(&entry.List)
	if !found {
		delete(s.data, key)
		return "", false, nil
	}

	if len(entry.List.values) == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = entry
	}

	return value, true, nil
}
func (l *List) LPop() (string, bool) {
	if len(l.values) == 0 {
		return "", false
	}
	value := l.values[0]
	l.values[0] = ""
	l.values = l.values[1:]
	return value, true
}
func (l *List) RPop() (string, bool) {
	if len(l.values) == 0 {
		return "", false
	}
	last := len(l.values) - 1
	value := l.values[last]
	l.values[last] = ""
	l.values = l.values[:last]
	return value, true
}

func (s *Store) LLen(key string) (int, error) {
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
	if entry.Type != ListType {
		return 0, ErrWrongType
	}

	return entry.List.LLen(), nil
}
func (l *List) LLen() int {
	return len(l.values)
}
func (s *Store) LRange(key string, start, end int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]
	if !ok {
		return []string{}, nil
	}
	if isExpired(entry) {
		delete(s.data, key)
		return []string{}, nil
	}
	if entry.Type != ListType {
		return []string{}, ErrWrongType
	}
	n := len(entry.List.values)
	if n == 0 {
		return []string{}, nil
	}
	if start < 0 {
		start = n + start
	}
	if end < 0 {
		end = n + end
	}
	if start < 0 {
		start = 0
	}
	if end >= n {
		end = n - 1
	}
	if start >= n || end < 0 || start > end {
		return []string{}, nil
	}
	return entry.List.LRange(start, end), nil
}
func (l *List) LRange(start, end int) []string {
	var values = make([]string, end-start+1)
	copy(values, l.values[start:end+1])
	return values
}
func (s *Store) LIndex(key string, index int) (string, bool, error) {
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
	if entry.Type != ListType {
		return "", false, ErrWrongType
	}
	if index < 0 {
		index = len(entry.List.values) + index
	}
	value, ok := entry.List.LIndex(index)
	if !ok {
		return "", false, nil
	}
	return value, true, nil
}
func (l *List) LIndex(index int) (string, bool) {
	if index < 0 || index >= len(l.values) {
		return "", false
	}
	return l.values[index], true
}
func (s *Store) LSet(key string, index int, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]
	if !ok {
		return ErrNoKey
	}
	if isExpired(entry) {
		delete(s.data, key)
		return ErrNoKey
	}
	if entry.Type != ListType {
		return ErrWrongType
	}
	if index < 0 {
		index = len(entry.List.values) + index
	}
	err := entry.List.LSet(index, value)
	if err != nil {
		return err
	}
	s.data[key] = entry
	return nil
}
func (l *List) LSet(index int, value string) error {
	if index < 0 || index >= len(l.values) {
		return ErrIndexOutOfRange
	}
	l.values[index] = value
	return nil
}
func (s *Store) LTrim(key string, start, end int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.data[key]
	if !ok {
		return nil
	}
	if isExpired(entry) {
		delete(s.data, key)
		return nil
	}
	if entry.Type != ListType {
		return ErrWrongType
	}
	n := len(entry.List.values)
	if n == 0 {
		return nil
	}
	if start < 0 {
		start = n + start
	}
	if end < 0 {
		end = n + end
	}
	if start < 0 {
		start = 0
	}
	if end >= n {
		end = n - 1
	}
	if start >= n || end < 0 || start > end {
		delete(s.data, key)
		return nil
	}
	entry.List.LTrim(start, end)
	s.data[key] = entry
	return nil
}
func (l *List) LTrim(start, end int) {
	trimmed := make([]string, end-start+1)
	copy(trimmed, l.values[start:end+1])
	l.values = trimmed
}
