package store

import "errors"

var (
	ErrNotIntegerOrOutOfRange = errors.New("value is not an integer or out of range")
	ErrNegativeDecrement      = errors.New("decrement must be non-negative")
	ErrWrongType              = errors.New("WRONGTYPE operation against a key holding the wrong kind of value")
	ErrIndexOutOfRange        = errors.New("index out of range")
	ErrNoKey                  = errors.New("no such key")
	ErrInvalidScore           = errors.New("invalid score")
)
