package server

import (
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func lpush(s *Server, args []string) protocol.Value {
	if len(args) < 2 {
		return protocol.NewError("need at least 2 arguments in 'LPUSH'")
	}
	length, err := s.store.LPush(args[0], args[1:]...)
	if err != nil {
		return protocol.NewError(err.Error())
	}

	return protocol.NewInteger(int64(length))
}

func rpush(s *Server, args []string) protocol.Value {
	if len(args) < 2 {
		return protocol.NewError("need at least 2 arguments in 'RPUSH'")
	}
	length, err := s.store.RPush(args[0], args[1:]...)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}

func lpop(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'LPOP'")
	}
	value, ok, err := s.store.LPop(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}
func rpop(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'RPOP'")
	}
	value, ok, err := s.store.RPop(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}
func llen(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'LLEN'")
	}
	length, err := s.store.LLen(args[0])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewInteger(int64(length))
}

func lrange(s *Server, args []string) protocol.Value {
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

func lindex(s *Server, args []string) protocol.Value {
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
func lset(s *Server, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments for 'LSET'")
	}
	index, err := parseIntHelper(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	err = s.store.LSet(args[0], index, args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewSimpleString("OK")
}

func ltrim(s *Server, args []string) protocol.Value {
	if len(args) != 3 {
		return protocol.NewError("need 3 arguments for 'LTRIM'")
	}
	start, err := parseIntHelper(args[1])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	end, err := parseIntHelper(args[2])
	if err != nil {
		return protocol.NewError(err.Error())
	}
	err = s.store.LTrim(args[0], start, end)
	if err != nil {
		return protocol.NewError(err.Error())
	}
	return protocol.NewSimpleString("OK")
}
