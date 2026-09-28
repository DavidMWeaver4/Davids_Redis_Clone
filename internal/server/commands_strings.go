package server

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func set(s *Server, client *Client, args []string) protocol.Value {
	err := s.executeSet(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}

	return protocol.NewSimpleString("OK")
}

func (s *Server) executeSet(args []string) error {
	if len(args) != 2 && len(args) != 4 {
		return errors.New("need 2 or 4 arguments for 'SET'")
	}

	key := args[0]
	value := args[1]

	var ttl time.Duration

	if len(args) == 4 {
		if !strings.EqualFold(args[2], "EX") {
			return errors.New("need EX as third argument in 'SET'")
		}

		seconds, err := strconv.Atoi(args[3])
		if err != nil {
			return errors.New("invalid expire time in 'SET'")
		}

		if seconds <= 0 {
			return errors.New("expire time must be positive in 'SET'")
		}

		ttl = time.Duration(seconds) * time.Second
	}

	s.store.Set(key, value, ttl)
	return nil
}
func (s *Server) replaySet(args []string) error {
	return s.executeSet(args)
}

func get(s *Server, client *Client, args []string) protocol.Value {
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
func appendCommand(s *Server, client *Client, args []string) protocol.Value {
	strlen, err := s.executeAppendCommand(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(strlen))
}
func (s *Server) executeAppendCommand(args []string) (int, error) {
	if len(args) != 2 {
		return 0, errors.New("need 2 arguments for 'APPEND'")
	}
	return s.store.Append(args[0], args[1])
}
func (s *Server) replayAppendCommand(args []string) error {
	_, err := s.executeAppendCommand(args)
	return err
}
func strlen(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'STRLEN'")
	}
	strlen, err := s.store.Strlen(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(strlen))
}
func setnx(s *Server, client *Client, args []string) protocol.Value {
	ok, err := s.executeSetnx(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}
func (s *Server) executeSetnx(args []string) (bool, error) {
	if len(args) != 2 {
		return false, errors.New("need 2 arguments for 'SETNX'")
	}
	return s.store.Setnx(args[0], args[1])
}
func (s *Server) replaySetnx(args []string) error {
	_, err := s.executeSetnx(args)
	return err
}

func mget(s *Server, client *Client, args []string) protocol.Value {
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

func mset(s *Server, client *Client, args []string) protocol.Value {
	err := s.executeMSet(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}

	return protocol.NewSimpleString("OK")
}

func (s *Server) executeMSet(args []string) error {
	if len(args) < 2 || len(args)%2 != 0 {
		return errors.New("need at least 2 or even amount of arguments for 'MSET'")
	}

	pairs := make([]store.KeyValue, 0, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		pairs = append(pairs, store.KeyValue{
			Key:   args[i],
			Value: args[i+1],
		})
	}

	s.store.Mset(pairs)

	return nil
}

func (s *Server) replayMSet(args []string) error {
	return s.executeMSet(args)
}
