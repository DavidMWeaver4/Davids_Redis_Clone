package server

import (
	"strconv"
	"strings"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

type commandHandler func(*Server, []string) protocol.Value

var commandHandlers = map[string]commandHandler{
	"PING":   ping,
	"SET":    set,
	"GET":    get,
	"DEL":    deleteCommand,
	"EXISTS": exists,
	"TTL":    ttl,
}

func (s *Server) execute(command protocol.Value) protocol.Value {
	if command.Type != protocol.Array {
		return protocol.NewError("ERR command must be an array")
	}
	if len(command.Array) == 0 {
		return protocol.NewError("ERR no command entered")
	}
	cmd := command.Array[0]
	if cmd.Type != protocol.BulkString {
		return protocol.NewError("ERR command must be a bulk string")
	}
	args := make([]string, len(command.Array)-1)

	for i, value := range command.Array[1:] {
		if value.Type != protocol.BulkString {
			return protocol.NewError("ERR arguments must be bulk strings")
		}
		args[i] = value.Str
	}
	handler, ok := commandHandlers[strings.ToUpper(cmd.Str)]
	if !ok {
		return protocol.NewError("ERR invalid command")
	}
	return handler(s, args)
}
func ping(s *Server, args []string) protocol.Value {
	return protocol.NewSimpleString("PONG")
}

func set(s *Server, args []string) protocol.Value {
	if len(args) != 2 && len(args) != 4 {
		return protocol.NewError("ERR need 2 or 4 arguments for 'set'")
	}
	var ttl time.Duration
	key := args[0]
	value := args[1]
	if len(args) == 4 {
		ok := strings.EqualFold(args[2], "EX")
		if !ok {
			return protocol.NewError("ERR need EX as third argument in 'set'")
		}
		seconds, err := strconv.Atoi(args[3])
		if err != nil {
			return protocol.NewError("ERR invalid expire time in 'set'")
		}
		if seconds <= 0 {
			return protocol.NewError("ERR expire time must be positive in 'set'")
		}
		ttl = time.Second * time.Duration(seconds)
	}
	s.store.Set(key, value, ttl)
	return protocol.NewSimpleString("OK")
}

func get(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("ERR need 1 argument for 'get'")
	}
	value, ok := s.store.Get(args[0])
	if !ok {
		return protocol.NewNullBulkString()
	}
	return protocol.NewBulkString(value)
}

func deleteCommand(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("ERR need 1 argument for 'del'")
	}
	deleted := s.store.Delete(args[0])
	return protocol.NewInteger(int64(deleted))
}

func exists(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("ERR need 1 argument for 'exists'")
	}
	if s.store.Exists(args[0]) {
		return protocol.NewInteger(1)
	}
	return protocol.NewInteger(0)
}

func ttl(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("ERR need 1 argument for 'ttl'")
	}
	seconds := s.store.TTL(args[0])
	return protocol.NewInteger(int64(seconds))
}
