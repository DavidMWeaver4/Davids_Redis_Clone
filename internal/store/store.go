package store

import (
	"context"
	"sync"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store/skiplist"
)

type Store struct {
	mu     sync.RWMutex
	data   map[string]Entry
	cancel context.CancelFunc
	wg     sync.WaitGroup
}
type ValueType int

const (
	InvalidType ValueType = iota
	StringType
	ListType
	HashType
	ZSetType
)

type Entry struct {
	Type      ValueType
	String    string
	List      List
	Hash      Hash
	ZSet      ZSet
	ExpiresAt time.Time
}

type List struct {
	values []string
}
type Hash struct {
	values map[string]string
}
type ZSet struct {
	scores map[string]float64
	list   *skiplist.SkipList
}

type KeyValue struct {
	Key   string
	Value string
}

func New() *Store {
	return NewStore()
}
func NewStore() *Store {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Store{
		data:   make(map[string]Entry),
		cancel: cancel,
	}

	s.wg.Add(1)
	go s.expirationLoop(ctx)

	return s
}
func (s *Store) Close() {
	s.cancel()
	s.wg.Wait()
}

func (s *Store) getEntry(key string) (Entry, bool) {
	entry, ok := s.data[key]
	if !ok {
		return Entry{}, false
	}
	if isExpired(entry) {
		delete(s.data, key)
		return Entry{}, false
	}
	return entry, true
}
func (s *Store) getEntryForWrite(key string) (Entry, bool) {
	entry, ok := s.data[key]
	if !ok || isExpired(entry) {
		return Entry{}, false
	}
	return entry, true
}
