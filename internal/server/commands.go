package server

import (
	"net"
	"strings"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

type commandHandler func(*Server, []string) protocol.Value

var commandHandlers = map[string]commandHandler{
	"PING":   ping,
	"SET":    set,
	"GET":    get,
	"DEL":    deleteCommand,
	"EXISTS": exists,
}

func (s *Server) execute(command string, conn net.Conn) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return protocol.Write(conn, protocol.NewError("ERR no command entered"))
	}
	handler, ok := commandHandlers[strings.ToUpper(parts[0])]
	if !ok {
		return protocol.Write(conn, protocol.NewError("ERR invalid command"))
	}
	response := handler(s, parts[1:])
	return protocol.Write(conn, response)
}
func ping(s *Server, args []string) protocol.Value {
	return protocol.NewSimpleString("PONG")
}

func set(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("ERR need 2 arguments for 'set'")

	}
	key := args[0]
	value := args[1]

	s.store.Set(key, value)
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
