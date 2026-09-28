package server

import (
	"errors"
	"strconv"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func incr(s *Server, client *Client, args []string) protocol.Value {
	newValue, err := s.executeIncr(args)
	if err != nil {
		if errors.Is(err, store.ErrNotIntegerOrOutOfRange) {
			return protocol.NewError(err.Error())
		}
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}
func (s *Server) executeIncr(args []string) (int64, error) {
	if len(args) != 1 {
		return 0, errors.New("need 1 argument for 'INCR'")
	}
	return s.store.Incr(args[0])
}
func (s *Server) replayIncr(args []string) error {
	_, err := s.executeIncr(args)
	return err
}

func decr(s *Server, client *Client, args []string) protocol.Value {
	newValue, err := s.executeDecr(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}
func (s *Server) executeDecr(args []string) (int64, error) {
	if len(args) != 1 {
		return 0, errors.New("need 1 argument for 'DECR'")
	}
	return s.store.Decr(args[0])
}
func (s *Server) replayDecr(args []string) error {
	_, err := s.executeDecr(args)
	return err
}

func incrby(s *Server, client *Client, args []string) protocol.Value {

	newValue, err := s.executeIncrby(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}
func (s *Server) executeIncrby(args []string) (int64, error) {
	if len(args) != 2 {
		return 0, errors.New("need 2 arguments for 'INCRBY'")
	}
	inc, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return 0, errors.New("cannot INCRBY non-integer")
	}
	return s.store.Incrby(args[0], inc)
}
func (s *Server) replayIncrby(args []string) error {
	_, err := s.executeIncrby(args)
	return err
}

func decrby(s *Server, client *Client, args []string) protocol.Value {
	newValue, err := s.executeDecrby(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(newValue)
}
func (s *Server) executeDecrby(args []string) (int64, error) {
	if len(args) != 2 {
		return 0, errors.New("need 2 arguments for 'DECRBY'")
	}
	dec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return 0, errors.New("cannot DECRBY non-integer")
	}
	if dec < 0 {
		return 0, errors.New("decrement must be non-negative")
	}
	return s.store.Decrby(args[0], dec)
}
func (s *Server) replayDecrby(args []string) error {
	_, err := s.executeDecrby(args)
	return err
}
