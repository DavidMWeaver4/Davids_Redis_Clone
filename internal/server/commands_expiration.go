package server

import (
	"errors"
	"strconv"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func ttl(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'TTL'")
	}
	seconds := s.store.TTL(args[0])
	return protocol.NewInteger(int64(seconds))
}

func expire(s *Server, client *Client, args []string) protocol.Value {
	ok, err := s.executeExpire(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}
func (s *Server) executeExpire(args []string) (bool, error) {
	if len(args) != 2 {
		return false, errors.New("need 2 arguments for 'EXPIRE'")
	}
	seconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return false, errors.New("invalid expire time in 'EXPIRE'")
	}
	ok := s.store.Expire(args[0], time.Duration(seconds)*time.Second)
	return ok, nil
}
func (s *Server) replayExpire(args []string) error {
	_, err := s.executeExpire(args)
	return err
}
func persist(s *Server, client *Client, args []string) protocol.Value {
	ok, err := s.executePersist(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}
func (s *Server) executePersist(args []string) (bool, error) {
	if len(args) != 1 {
		return false, errors.New("need 1 argument for 'PERSIST'")
	}
	ok := s.store.Persist(args[0])
	return ok, nil
}
func (s *Server) replayPersist(args []string) error {
	_, err := s.executePersist(args)
	return err
}
