package store

import (
	"math"
	"strconv"
)

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
