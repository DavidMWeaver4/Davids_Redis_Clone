package server

import (
	"errors"
	"strconv"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func incr(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'INCR'")
	}
	newValue, err := s.store.Incr(args[0])
	if err != nil {
		if errors.Is(err, store.ErrNotIntegerOrOutOfRange) {
			return protocol.NewError(err.Error())
		}
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}

func decr(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'DECR'")
	}
	newValue, err := s.store.Decr(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}

func incrby(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'INCRBY'")
	}
	inc, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return protocol.NewError("cannot INCRBY non-integer")
	}
	newValue, err := s.store.Incrby(args[0], inc)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}

func decrby(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'DECRBY'")
	}
	dec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return protocol.NewError("cannot DECRBY non-integer")
	}
	if dec < 0 {
		return protocol.NewError("decrement must be non-negative")
	}
	newValue, err := s.store.Decrby(args[0], dec)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}
