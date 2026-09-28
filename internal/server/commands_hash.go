package server

import (
	"errors"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func hset(s *Server, client *Client, args []string) protocol.Value {
	success, err := s.executeHSet(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(success))
}
func (s *Server) executeHSet(args []string) (int, error) {
	if len(args) != 3 {
		return 0, errors.New("need 3 arguments in 'HSET'")
	}
	return s.store.HSet(args[0], args[1], args[2])
}
func (s *Server) replayHSet(args []string) error {
	_, err := s.executeHSet(args)
	return err
}

func hget(s *Server, client *Client, args []string) protocol.Value {
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

func hdel(s *Server, client *Client, args []string) protocol.Value {
	removed, err := s.executeHDel(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(removed))
}
func (s *Server) executeHDel(args []string) (int, error) {
	if len(args) < 2 {
		return 0, errors.New("need at least 2 arguments in 'HDEL'")
	}
	return s.store.HDel(args[0], args[1:]...)
}
func (s *Server) replayHDel(args []string) error {
	_, err := s.executeHDel(args)
	return err
}
func hexists(s *Server, client *Client, args []string) protocol.Value {
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

func hlen(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument in 'HLEN'")
	}
	length, err := s.store.HLen(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}

func hgetall(s *Server, client *Client, args []string) protocol.Value {
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
