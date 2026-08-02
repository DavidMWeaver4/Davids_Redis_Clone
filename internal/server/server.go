package server

import "https://github.com/DavidMWeaver4/Davids_Redis_Clone/redis_clone/internal/store/store.go"

type Server struct {
	addr  string
	store *store.Store
}

func New(addr string, s *store.Store) *Server {
	return &Server{
		addr:  addr,
		store: s,
	}
}
