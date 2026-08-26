package server

import (
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func hset(s *Server, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments in 'HSET'")
	}
	success, err := s.store.HSet(args[0], args[1], args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(success))
}

func hget(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments in 'HGET'")
	}
	value, found, err := s.store.HGet(args[0], args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !found {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}

func hdel(s *Server, args []string) protocol.Value {
	if len(args) < 2 {
		return protocol.NewError("need at least 2 arguments in 'HDEL'")
	}
	removed, err := s.store.HDel(args[0], args[1:]...)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(removed))
}

func hexists(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments in 'HEXISTS'")
	}
	found, err := s.store.HExists(args[0], args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !found {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}

func hlen(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument in 'HLEN'")
	}
	length, err := s.store.HLen(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}

func hgetall(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument in 'HGETALL'")
	}
	fields, err := s.store.HGetAll(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	results := make([]protocol.Value, 0, len(fields)*2)
	for field, value := range fields {
		results = append(
			results,
			protocol.NewBulkString(field),
			protocol.NewBulkString(value))
	}
	return protocol.NewArray(results)
}
