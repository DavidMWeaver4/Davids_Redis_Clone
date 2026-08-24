package store

import (
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

func New() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}
