package server

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/persistence"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

type Server struct {
	addr     string
	store    *store.Store
	listener net.Listener

	wg sync.WaitGroup

	mu           sync.Mutex
	clients      map[net.Conn]*Client
	shuttingDown bool

	aof               *persistence.AOF
	persistenceFailed bool
}

func New(addr string, s *store.Store, aof *persistence.AOF) *Server {
	return &Server{
		addr:    addr,
		store:   s,
		clients: make(map[net.Conn]*Client),
		aof:     aof,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.listener = listener
	s.mu.Unlock()

	if s.aof != nil {
		go s.monitorPersistence()
	}
	log.Printf("Redis clone listening on %s", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			log.Printf("error: %v", err)
			return err
		}
		client := &Client{conn: conn}
		s.mu.Lock()
		if s.shuttingDown {
			s.mu.Unlock()
			conn.Close()
			continue
		}
		s.clients[conn] = client
		s.wg.Add(1)
		s.mu.Unlock()
		go func() {
			defer s.wg.Done()
			s.handleClient(client)
		}()
	}
}

func (s *Server) handleClient(client *Client) {
	defer client.conn.Close()
	defer func() {
		s.mu.Lock()
		delete(s.clients, client.conn)
		s.mu.Unlock()
	}()

	reader := bufio.NewReader(client.conn)

	for {
		value, err := protocol.Read(reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			log.Printf("client read error: %v", err)
			return
		}

		response := s.execute(client, value)
		err = protocol.Write(client.conn, response)
		if err != nil {
			log.Printf("client error: %v", err)
			return
		}
	}
}
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	s.shuttingDown = true
	listener := s.listener
	s.mu.Unlock()
	if listener != nil {
		_ = listener.Close()
	}
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.mu.Lock()
		clients := make([]*Client, 0, len(s.clients))
		for _, client := range s.clients {
			clients = append(clients, client)
		}
		s.mu.Unlock()

		for _, client := range clients {
			_ = client.conn.Close()
		}
		<-done
		return ctx.Err()
	}
}
