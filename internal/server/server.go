package server

import (
	"bufio"
	"log"
	"net"
	"strings"

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
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		line = strings.TrimSpace(line)
		err = s.execute(line, conn)
		if err != nil {
			log.Printf("client error %v", err)
			return
		}
	}
}
