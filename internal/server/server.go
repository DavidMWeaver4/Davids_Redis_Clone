package server

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

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

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		value, err := protocol.Read(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("client read error: %v", err)
			return
		}

		response := s.execute(value)
		err = protocol.Write(conn, response)
		if err != nil {
			log.Printf("client error: %v", err)
			return
		}
	}
}
