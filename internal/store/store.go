package store

import (
	"context"
	"sync"
	"time"
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
)

type Entry struct {
	Type      ValueType
	String    string
	List      List
	Hash      Hash
	ExpiresAt time.Time
}

type List struct {
	values []string
}
type Hash struct {
	values map[string]string
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
