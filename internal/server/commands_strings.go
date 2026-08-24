package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func set(s *Server, args []string) protocol.Value {
	if len(args) != 2 && len(args) != 4 {
		return protocol.NewError("need 2 or 4 arguments for 'SET'")
	}
	var ttl time.Duration
	key := args[0]
	value := args[1]
	if len(args) == 4 {
		ok := strings.EqualFold(args[2], "EX")
		if !ok {
			return protocol.NewError("need EX as third argument in 'SET'")
		}
		seconds, err := strconv.Atoi(args[3])
		if err != nil {
			return protocol.NewError("invalid expire time in 'SET'")
		}
		if seconds <= 0 {
			return protocol.NewError("expire time must be positive in 'SET'")
		}
		ttl = time.Second * time.Duration(seconds)
	}
	s.store.Set(key, value, ttl)
	return protocol.NewSimpleString("OK")
}

func get(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'GET'")
	}
	value, ok, err := s.store.Get(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}
func appendCommand(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'APPEND'")
	}
	strlen, err := s.store.Append(args[0], args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(strlen))
}
func strlen(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'STRLEN'")
	}
	strlen, err := s.store.Strlen(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(strlen))
}
func setnx(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'SETNX'")
	}
	ok, err := s.store.Setnx(args[0], args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}
func mget(s *Server, args []string) protocol.Value {
	if len(args) < 1 {
		return protocol.NewError("need at least 1 argument for 'MGET'")
	}
	getResults := make([]protocol.Value, 0, len(args))
	for _, key := range args {
		value, found, err := s.store.Get(key)
		if err != nil {
			return protocol.NewError(err.Error())
		}
		if !found {
			getResults = append(getResults, protocol.NewNullBulkString())
			continue
		}
		getResults = append(getResults, protocol.NewBulkString(value))
	}
	return protocol.NewArray(getResults)
}

func mset(s *Server, args []string) protocol.Value {
	if len(args) < 2 || len(args)%2 != 0 {
		return protocol.NewError("need at least 2 or even amount of arguments for 'MSET'")
	}
	pairs := make([]store.KeyValue, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		pairs = append(pairs, store.KeyValue{
			Key:   args[i],
			Value: args[i+1],
		})
	}
	s.store.Mset(pairs)
	return protocol.NewSimpleString("OK")
}
