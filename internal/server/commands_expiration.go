package server

import (
	"strconv"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func ttl(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'TTL'")
	}
	seconds := s.store.TTL(args[0])
	return protocol.NewInteger(int64(seconds))
}

func expire(s *Server, args []string) protocol.Value {
	if len(args) != 2 {
		return protocol.NewError("need 2 arguments for 'EXPIRE'")
	}
	seconds, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return protocol.NewError("invalid expire time in 'EXPIRE'")
	}

	if s.store.Expire(args[0], time.Duration(seconds)*time.Second) {
		return protocol.NewInteger(1)
	}
	return protocol.NewInteger(0)
}

func persist(s *Server, args []string) protocol.Value {
	if len(args) != 1 {
		return protocol.NewError("need 1 argument for 'PERSIST'")
	}
	if !s.store.Persist(args[0]) {
		return protocol.NewInteger(0)
	}
	return protocol.NewInteger(1)
}
