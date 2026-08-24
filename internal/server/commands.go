package server

import (
	"strings"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
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
