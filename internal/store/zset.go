package store

import (
	"math"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store/skiplist"
)

func (s *Store) ZAdd(key string, score float64, member string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if math.IsNaN(score) {
		return 0, ErrInvalidScore
	}
	entry, ok := s.getEntryForWrite(key)
	if !ok {
		entry = Entry{
			Type: ZSetType,
			ZSet: ZSet{},
		}
	}
	if entry.Type != ZSetType {
		return 0, ErrWrongType
	}
	added := entry.ZSet.ZAdd(score, member)
	s.data[key] = entry
	return added, nil
}
func (z *ZSet) ZAdd(score float64, member string) int {
	if z.scores == nil {
		z.scores = make(map[string]float64)
		z.list = skiplist.NewSkipList()
	}
	oldScore, exists := z.scores[member]
	if exists && (oldScore == score) {
		return 0
	}
	if exists && (oldScore != score) {
		z.list.Delete(oldScore, member)
		z.list.Insert(score, member)
		z.scores[member] = score
		return 0
	}
	z.list.Insert(score, member)
	z.scores[member] = score
	return 1
}
func (s *Store) ZScore(key string, member string) (float64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.getEntry(key)
	if !ok {
		return 0, false, nil
	}
	if entry.Type != ZSetType {
		return 0, false, ErrWrongType
	}
	gotScore, found := entry.ZSet.ZScore(member)
	if !found {
		return 0, found, nil
	}
	return gotScore, found, nil
}
func (z *ZSet) ZScore(member string) (float64, bool) {
	score, found := z.scores[member]
	return score, found
}
func (s *Store) ZCard(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.getEntry(key)
	if !ok {
		return 0, nil
	}
	if entry.Type != ZSetType {
		return 0, ErrWrongType
	}
	return entry.ZSet.ZCard(), nil
}
func (z *ZSet) ZCard() int {
	return len(z.scores)
}
func (s *Store) ZRem(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.getEntry(key)
	if !ok {
		return 0, nil
	}
	if entry.Type != ZSetType {
		return 0, ErrWrongType
	}
	count := 0
	for _, mem := range members {
		count += entry.ZSet.ZRem(mem)
	}
	if entry.ZSet.ZCard() == 0 {
		delete(s.data, key)
	} else {
		s.data[key] = entry
	}
	return count, nil
}
func (z *ZSet) ZRem(member string) int {
	score, found := z.scores[member]
	if !found {
		return 0
	}
	z.list.Delete(score, member)
	delete(z.scores, member)
	return 1

}
func (s *Store) ZRange(key string, start, stop int) ([]skiplist.MemberScore, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.getEntry(key)
	if !ok {
		return []skiplist.MemberScore{}, nil
	}
	if entry.Type != ZSetType {
		return []skiplist.MemberScore{}, ErrWrongType
	}

	return entry.ZSet.ZRange(start, stop), nil
}

func (z *ZSet) ZRange(start, stop int) []skiplist.MemberScore {
	size := len(z.scores)
	if size == 0 {
		return []skiplist.MemberScore{}
	}
	if start < 0 {
		start = size + start
	}
	if stop < 0 {
		stop = size + stop
	}
	if start < 0 {
		start = 0
	}
	if start >= size || stop < 0 || start > stop {
		return []skiplist.MemberScore{}
	}
	if stop >= size {
		stop = size - 1
	}

	return z.list.RangeByRank(start, stop)
}

func (s *Store) ZRank(key string, member string) (int, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.getEntry(key)
	if !ok {
		return 0, false, nil
	}
	if entry.Type != ZSetType {
		return 0, false, ErrWrongType
	}
	rank, found := entry.ZSet.ZRank(member)
	return rank, found, nil
}
func (z *ZSet) ZRank(member string) (int, bool) {
	score, found := z.scores[member]
	if !found {
		return 0, false
	}
	rank := z.list.Rank(score, member)
	return rank, rank >= 0
}
func (s *Store) ZIncrby(key string, inc float64, member string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if math.IsNaN(inc) {
		return 0.0, ErrInvalidNumber
	}
	entry, ok := s.getEntryForWrite(key)
	if !ok {
		entry = Entry{
			Type: ZSetType,
			ZSet: ZSet{},
		}
	}
	if entry.Type != ZSetType {
		return 0.0, ErrWrongType
	}
	currentScore, found := entry.ZSet.scores[member]
	if !found {
		currentScore = 0
	}
	newScore := currentScore + inc
	if math.IsNaN(newScore) {
		return 0, ErrInvalidNumber
	}
	entry.ZSet.ZAdd(newScore, member)
	s.data[key] = entry
	return newScore, nil
}
