package server

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

type commandHandler func(*Server, []string) protocol.Value

var commandHandlers = map[string]commandHandler{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"DEL":     deleteCommand,
	"EXISTS":  exists,
	"TTL":     ttl,
	"EXPIRE":  expire,
	"PERSIST": persist,
	"INCR":    incr,
	"DECR":    decr,
	"INCRBY":  incrby,
	"APPEND":  appendCommand,
	"STRLEN":  strlen,
	"SETNX":   setnx,
	"MGET":    mget,
	"MSET":    mset,
}

func (s *Server) execute(command protocol.Value) protocol.Value {
	if command.Type != protocol.Array {
		return protocol.NewError("command must be an array")
	}
	if len(command.Array) == 0 {
		return protocol.NewError("no command entered")
	}
	cmd := command.Array[0]
	if cmd.Type != protocol.BulkString {
		return protocol.NewError("command must be a bulk string")
	}
	args := make([]string, len(command.Array)-1)

	for i, value := range command.Array[1:] {
		if value.Type != protocol.BulkString {
			return protocol.NewError("arguments must be bulk strings")
		}
		args[i] = value.Str
	}
	handler, ok := commandHandlers[strings.ToUpper(cmd.Str)]
	if !ok {
		return protocol.NewError("invalid command")
	}
	return handler(s, args)
}
func ping(s *Server, args []string) protocol.Value {
	return protocol.NewSimpleString("PONG")
}

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

func deleteCommand(s *Server, args []string) protocol.Value {
	if len(args) < 1 {
		return protocol.NewError("need at least 1 argument for 'DEL'")
	}
	deleted := 0
	for _, key := range args {
		deleted += s.store.Delete(key)
	}
	return protocol.NewInteger(int64(deleted))
}

func exists(s *Server, args []string) protocol.Value {
	if len(args) < 1 {
		return protocol.NewError("need at least 1 argument for 'EXISTS'")
	}
	count := 0
	for _, key := range args {
		if s.store.Exists(key) {
			count++
		}
	}
	return protocol.NewInteger(int64(count))
}

func ttl(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'TTL'")
	}
	seconds := s.store.TTL(args[0])
	return protocol.NewInteger(int64(seconds))
}

func expire(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'EXPIRE'")
	}
	seconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return protocol.NewError("invalid expire time in 'EXPIRE'")
	}

	if s.store.Expire(args[0], time.Duration(seconds)*time.Second) {
		return protocol.NewInteger(1)
	}
	return protocol.NewInteger(0)
}

func persist(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'PERSIST'")
	}
	if !s.store.Persist(args[0]) {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}

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
