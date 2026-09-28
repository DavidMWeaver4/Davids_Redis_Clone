package server

import (
	"errors"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func lpush(s *Server, client *Client, args []string) protocol.Value {
	length, err := s.executeLPush(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}
func (s *Server) executeLPush(args []string) (int, error) {
	if len(args) < 2 {
		return 0, errors.New("need at least 2 arguments in 'LPUSH'")
	}
	return s.store.LPush(args[0], args[1:]...)
}
func (s *Server) replayLPush(args []string) error {
	_, err := s.executeLPush(args)
	return err
}

func rpush(s *Server, client *Client, args []string) protocol.Value {
	length, err := s.executeRPush(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}
func (s *Server) executeRPush(args []string) (int, error) {
	if len(args) < 2 {
		return 0, errors.New("need at least 2 arguments in 'RPUSH'")
	}
	return s.store.RPush(args[0], args[1:]...)
}
func (s *Server) replayRPush(args []string) error {
	_, err := s.executeRPush(args)
	return err
}

func lpop(s *Server, client *Client, args []string) protocol.Value {
	value, ok, err := s.executeLPop(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}
func (s *Server) executeLPop(args []string) (string, bool, error) {
	if len(args) != 1 {
		return "", false, errors.New("need 1 argument for 'LPOP'")
	}
	return s.store.LPop(args[0])
}
func (s *Server) replayLPop(args []string) error {
	_, _, err := s.executeLPop(args)
	return err
}
func rpop(s *Server, client *Client, args []string) protocol.Value {
	value, ok, err := s.executeRPop(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}
func (s *Server) executeRPop(args []string) (string, bool, error) {
	if len(args) != 1 {
		return "", false, errors.New("need 1 argument for 'RPOP'")
	}
	return s.store.RPop(args[0])
}
func (s *Server) replayRPop(args []string) error {
	_, _, err := s.executeRPop(args)
	return err
}
func llen(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'LLEN'")
	}
	length, err := s.store.LLen(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}

func lrange(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments for 'LRANGE'")
	}
	start, err := parseIntHelper(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	end, err := parseIntHelper(args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	values, err := s.store.LRange(args[0], start, end)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	result := make([]protocol.Value, 0, len(values))
	for _, value := range values {
		result = append(result, protocol.NewBulkString(value))
	}

	return protocol.NewArray(result)
}

func lindex(s *Server, client *Client, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'LINDEX'")
	}
	index, err := parseIntHelper(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	value, ok, err := s.store.LIndex(args[0], index)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}
func lset(s *Server, client *Client, args []string) protocol.Value {
	err := s.executeLSet(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewSimpleString("OK")
}
func (s *Server) executeLSet(args []string) error {
	if len(args) != 3 {
		return errors.New("need 3 arguments for 'LSET'")
	}
	index, err := parseIntHelper(args[1])
	if err != nil {
		return err
	}
	return s.store.LSet(args[0], index, args[2])
}
func (s *Server) replayLSet(args []string) error {
	return s.executeLSet(args)
}

func ltrim(s *Server, client *Client, args []string) protocol.Value {
	err := s.executeLTrim(args)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewSimpleString("OK")
}

func (s *Server) executeLTrim(args []string) error {
	if len(args) != 3 {
		return errors.New("need 3 arguments for 'LTRIM'")
	}
	start, err := parseIntHelper(args[1])
	if err != nil {
		return err
	}
	end, err := parseIntHelper(args[2])
	if err != nil {
		return err
	}
	return s.store.LTrim(args[0], start, end)
}
func (s *Server) replayLTrim(args []string) error {
	return s.executeLTrim(args)
}
