package server

import "github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"

func ping(s *Server, args []string) protocol.Value {
	return protocol.NewSimpleString("PONG")
}

func deleteCommand(s *Server, args []string) protocol.Value {
	if len(args) < 1 {
		return protocol.NewError("need at least 1 argument for 'DEL'")
	}
	deleted := 0
	for _, key := range args {
		deleted += s.store.Delete(key)
	}
	return protocol.NewInteger(int64(deleted))
}

func exists(s *Server, args []string) protocol.Value {
	if len(args) < 1 {
		return protocol.NewError("need at least 1 argument for 'EXISTS'")
	}
	count := 0
	for _, key := range args {
		if s.store.Exists(key) {
			count++
		}
	}
	return protocol.NewInteger(int64(count))
}
