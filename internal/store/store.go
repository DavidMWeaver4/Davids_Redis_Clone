package store

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

//TODO
// Set(key, value)
// Get(key)(string, bool)
// Delete(key)
// Exists(key)(bool)
