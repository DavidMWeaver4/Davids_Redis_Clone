package server

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestServer_HandleClient_Ping(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)
	err := protocol.Write(clientConn, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("PING"),
	}))
	if err != nil {
		t.Fatalf("failed to write command: %v", err)
	}

	got, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if got.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", got.Type)
	}
	if got.Str != "PONG" {
		t.Fatalf("expected PONG, got %q", got.Str)
	}

	clientConn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}
func TestServer_HandleClient_SetGet(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)

	setCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("Foo"),
		protocol.NewBulkString("Bar"),
	})

	err := protocol.Write(clientConn, setCommand)
	if err != nil {
		t.Fatalf("failed to write SET: %v", err)
	}
	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("failed to read SET response: %v", err)
	}
	if response.Type != protocol.SimpleString || response.Str != "OK" {
		t.Fatalf("expected OK, got %+v", response)
	}

	getCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("Foo"),
	})

	err = protocol.Write(clientConn, getCommand)
	if err != nil {
		t.Fatalf("failed to write GET: %v", err)
	}
	response, err = protocol.Read(reader)
	if err != nil {
		t.Fatalf("failed to read GET response: %v", err)
	}
	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}
	if response.Str != "Bar" {
		t.Fatalf("expected Bar, got %q", response.Str)
	}

	clientConn.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}
func TestServer_HandleClient_MultipleCommands(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)
	tests := []struct {
		command protocol.Value
		want    protocol.Value
	}{
		{
			command: protocol.NewArray([]protocol.Value{
				protocol.NewBulkString("PING"),
			}),
			want: protocol.NewSimpleString("PONG"),
		},
		{
			command: protocol.NewArray([]protocol.Value{
				protocol.NewBulkString("SET"),
				protocol.NewBulkString("Foo"),
				protocol.NewBulkString("Bar"),
			}),
			want: protocol.NewSimpleString("OK"),
		},
		{
			command: protocol.NewArray([]protocol.Value{
				protocol.NewBulkString("GET"),
				protocol.NewBulkString("Foo"),
			}),
			want: protocol.NewBulkString("Bar"),
		},
	}

	for _, tt := range tests {
		err := protocol.Write(clientConn, tt.command)
		if err != nil {
			t.Fatalf("failed to write command: %v", err)
		}
		got, err := protocol.Read(reader)
		if err != nil {
			t.Fatalf("failed to read response: %v", err)
		}
		if got.Type != tt.want.Type {
			t.Fatalf("expected type %v, got %v", tt.want.Type, got.Type)
		}
		if got.Str != tt.want.Str {
			t.Fatalf("expected %q, got %q", tt.want.Str, got.Str)
		}
	}

	clientConn.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}

func TestServer_HandleClient_InvalidRESP(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)

	_, err := clientConn.Write([]byte("@invalid\r\n"))
	if err != nil {
		t.Fatalf("failed to write invalid RESP: %v", err)
	}
	response, err := protocol.Read(reader)
	if err == nil {
		t.Fatalf("expected connection to close, got response: %+v", response)
	}
	clientConn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}

func TestServer_HandleClient_InvalidCommandType(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewInteger(123),
	})
	err := protocol.Write(clientConn, command)
	if err != nil {
		t.Fatalf("failed to write command: %v", err)
	}
	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
	if response.Str != "ERR command must be a bulk string" {
		t.Fatalf("expected command type error, got %q", response.Str)
	}
	clientConn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}

func TestServer_HandleClient_EmptyCommand(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)

	command := protocol.NewArray([]protocol.Value{})

	err := protocol.Write(clientConn, command)
	if err != nil {
		t.Fatalf("failed to write command: %v", err)
	}
	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
	if response.Str != "ERR no command entered" {
		t.Fatalf("expected no command error, got %q", response.Str)
	}
	clientConn.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}

func TestServer_HandleClient_ConnectionClosing(t *testing.T) {
	server := New("", store.New())
	clientConn, _, done := newTestConnection(t, server)

	clientConn.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close after client disconnected")
	}
}

func TestServer_HandleClient_SetWithTTL(t *testing.T) {
	server := New("", store.New())
	clientConn, reader, done := newTestConnection(t, server)

	_, err := clientConn.Write([]byte(
		"*5\r\n" +
			"$3\r\nSET\r\n" +
			"$3\r\nfoo\r\n" +
			"$3\r\nbar\r\n" +
			"$2\r\nEX\r\n" +
			"$1\r\n1\r\n",
	))
	if err != nil {
		t.Fatalf("failed to send RESP: %v", err)
	}

	got, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if got != "+OK\r\n" {
		t.Fatalf("expected +OK\\r\\n, got %q", got)
	}

	value, ok, err := server.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "bar" {
		t.Fatalf("expected value %q, got %q", "bar", value)
	}

	ttl := server.store.TTL("foo")
	if ttl != 0 && ttl != 1 {
		t.Fatalf("expected TTL of 0 or 1, got %d", ttl)
	}
	clientConn.Close()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server did not close client connection")
	}
}

func newTestConnection(t *testing.T, server *Server) (net.Conn, *bufio.Reader, <-chan struct{}) {
	t.Helper()
	serverConn, clientConn := net.Pipe()
	done := make(chan struct{})
	go func() {
		server.handleClient(serverConn)
		close(done)
	}()
	reader := bufio.NewReader(clientConn)
	t.Cleanup(func() {
		clientConn.Close()
	})
	return clientConn, reader, done
}
