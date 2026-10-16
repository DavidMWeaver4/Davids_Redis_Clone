package server

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestTCP_Ping(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)

	reader := bufio.NewReader(conn)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("PING"),
	})

	err := protocol.Write(conn, command)
	if err != nil {
		t.Fatalf("write PING: %v", err)
	}

	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read PING response: %v", err)
	}

	if response.Type != protocol.SimpleString || response.Str != "PONG" {
		t.Fatalf("PING response = %#v, want simple string PONG", response)
	}
}

func TestTCP_SetAndGet(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)

	reader := bufio.NewReader(conn)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("tcp-key"),
		protocol.NewBulkString("hello"),
	})

	err := protocol.Write(conn, command)
	if err != nil {
		t.Fatalf("write SET: %v", err)
	}

	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read SET response: %v", err)
	}

	if response.Type != protocol.SimpleString || response.Str != "OK" {
		t.Fatalf("SET response = %#v, want simple string OK", response)
	}

	command = protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("tcp-key"),
	})
	err = protocol.Write(conn, command)
	if err != nil {
		t.Fatalf("write GET: %v", err)
	}
	response, err = protocol.Read(reader)
	if err != nil {
		t.Fatalf("read GET response: %v", err)
	}
	if response.Type != protocol.BulkString || response.Str != "hello" {
		t.Fatalf("GET response = %#v, want bulk string hello", response)
	}

}

func TestTCP_SetWrongArgumentCount(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)
	reader := bufio.NewReader(conn)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("only-a-key"),
	})

	err := protocol.Write(conn, command)
	if err != nil {
		t.Fatalf("write SET: %v", err)
	}

	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read SET response: %v", err)
	}

	if response.Type != protocol.Error {
		t.Fatalf("SET response = %#v, want error", response)
	}
}
func TestTCP_UnknownCommand(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)
	reader := bufio.NewReader(conn)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("NOPE"),
	})

	err := protocol.Write(conn, command)
	if err != nil {
		t.Fatalf("write unknown command: %v", err)
	}
	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read unknown-command response: %v", err)
	}
	if response.Type != protocol.Error {
		t.Fatalf("response = %#v, want RESP error", response)
	}
}
func TestTCP_MultipleCommandsOnOneConnection(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)
	reader := bufio.NewReader(conn)

	setCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("multi-key"),
		protocol.NewBulkString("hello"),
	})

	err := protocol.Write(conn, setCommand)
	if err != nil {
		t.Fatalf("write SET: %v", err)
	}

	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read SET response: %v", err)
	}

	if response.Type != protocol.SimpleString || response.Str != "OK" {
		t.Fatalf("SET response = %#v, want simple string OK", response)
	}

	getCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("multi-key"),
	})

	err = protocol.Write(conn, getCommand)
	if err != nil {
		t.Fatalf("write GET: %v", err)
	}

	response, err = protocol.Read(reader)
	if err != nil {
		t.Fatalf("read GET response: %v", err)
	}

	if response.Type != protocol.BulkString || response.Str != "hello" {
		t.Fatalf("GET response = %#v, want bulk string hello", response)
	}

	delCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("DEL"),
		protocol.NewBulkString("multi-key"),
	})

	err = protocol.Write(conn, delCommand)
	if err != nil {
		t.Fatalf("write DEL: %v", err)
	}

	response, err = protocol.Read(reader)
	if err != nil {
		t.Fatalf("read DEL response: %v", err)
	}

	if response.Type != protocol.Integer || response.Int != 1 {
		t.Fatalf("DEL response = %#v, want integer 1", response)
	}
	getCommand = protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("multi-key"),
	})

	err = protocol.Write(conn, getCommand)
	if err != nil {
		t.Fatalf("write GET after DEL: %v", err)
	}

	response, err = protocol.Read(reader)
	if err != nil {
		t.Fatalf("read GET after DEL response: %v", err)
	}

	if response.Type != protocol.Null {
		t.Fatalf("GET after DEL response = %#v, want null bulk string", response)
	}
}
func TestTCP_Expiration(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)
	reader := bufio.NewReader(conn)

	setCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("expire-key"),
		protocol.NewBulkString("temporary"),
		protocol.NewBulkString("EX"),
		protocol.NewBulkString("1"),
	})

	err := protocol.Write(conn, setCommand)
	if err != nil {
		t.Fatalf("write SET: %v", err)
	}

	response, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read SET response: %v", err)
	}
	if response.Type != protocol.SimpleString || response.Str != "OK" {
		t.Fatalf("SET response = %#v, want simple string OK", response)
	}

	deadline := time.Now().Add(5 * time.Second)

	for time.Now().Before(deadline) {
		getCommand := protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("GET"),
			protocol.NewBulkString("expire-key"),
		})

		err = protocol.Write(conn, getCommand)
		if err != nil {
			t.Fatalf("write GET: %v", err)
		}

		response, err = protocol.Read(reader)
		if err != nil {
			t.Fatalf("read GET response: %v", err)
		}

		if response.Type == protocol.Null {
			return
		}

		time.Sleep(25 * time.Millisecond)
	}

	t.Fatal("key did not expire before the deadline")
}

func TestTCP_Pipelining(t *testing.T) {
	_, conn, _ := tcpTestServerHelper(t)
	reader := bufio.NewReader(conn)

	setCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("pipeline-key"),
		protocol.NewBulkString("pipeline-value"),
	})

	getCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("pipeline-key"),
	})

	err := protocol.Write(conn, setCommand)
	if err != nil {
		t.Fatalf("write SET: %v", err)
	}
	err = protocol.Write(conn, getCommand)
	if err != nil {
		t.Fatalf("write GET: %v", err)
	}

	setResponse, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read SET response: %v", err)
	}
	if setResponse.Type != protocol.SimpleString || setResponse.Str != "OK" {
		t.Fatalf("SET response = %#v, want simple string OK", setResponse)
	}

	getResponse, err := protocol.Read(reader)
	if err != nil {
		t.Fatalf("read GET response: %v", err)
	}
	if getResponse.Type != protocol.BulkString || getResponse.Str != "pipeline-value" {
		t.Fatalf("GET response = %#v, want bulk string pipeline-value", getResponse)
	}
}
func TestTCP_TwoClientsShareState(t *testing.T) {
	_, connA, addr := tcpTestServerHelper(t)
	connB := tcpTestDialHelper(t, addr)

	readerA := bufio.NewReader(connA)
	readerB := bufio.NewReader(connB)

	setCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("shared-key"),
		protocol.NewBulkString("shared-value"),
	})

	err := protocol.Write(connA, setCommand)
	if err != nil {
		t.Fatalf("client A write SET: %v", err)
	}

	response, err := protocol.Read(readerA)
	if err != nil {
		t.Fatalf("client A read SET response: %v", err)
	}
	if response.Type != protocol.SimpleString || response.Str != "OK" {
		t.Fatalf("client A SET response = %#v, want simple string OK", response)
	}

	getCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("shared-key"),
	})

	err = protocol.Write(connB, getCommand)
	if err != nil {
		t.Fatalf("client B write GET: %v", err)
	}

	response, err = protocol.Read(readerB)
	if err != nil {
		t.Fatalf("client B read GET response: %v", err)
	}
	if response.Type != protocol.BulkString || response.Str != "shared-value" {
		t.Fatalf("client B GET response = %#v, want bulk string shared-value", response)
	}
}

/*
 *
 *
 * Helpers
 *
 */
func tcpTestDialHelper(t *testing.T, addr string) net.Conn {
	t.Helper()

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("DialTimeout() error = %v", err)
	}

	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		closeerr := conn.Close()
		if closeerr != nil {
			t.Errorf("conn close err: %v", closeerr)
		}
		t.Fatalf("SetDeadline() error = %v", err)
	}

	t.Cleanup(func() {
		err := conn.Close()
		if err != nil {
			t.Errorf("conn.Close() error = %v", err)
		}
	})

	return conn
}

func tcpTestServerHelper(t *testing.T) (*Server, net.Conn, string) {
	t.Helper()
	st := store.New()
	srv := New("127.0.0.1:0", st, nil)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		st.Close()
		t.Fatal(err)
	}

	serveDone := make(chan error, 1)
	go func() {
		serveDone <- srv.serve(listener)
	}()
	addr := listener.Addr().String()
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownerr := srv.Shutdown(ctx)
		if shutdownerr != nil {
			t.Errorf("Shutdown() after dial failure: %v", shutdownerr)
		}

		st.Close()
		t.Fatalf("DialTimeout() error: %v", err)
	}

	err = conn.SetDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		connErr := conn.Close()
		if connErr != nil {
			t.Log(connErr)
		}
		t.Fatal(err)
	}

	t.Cleanup(func() {
		err = conn.Close()
		if err != nil {
			t.Errorf("conn.Close() error = %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			t.Errorf("Shutdown() error = %v", err)
		}

		select {
		case err := <-serveDone:
			if err != nil {
				t.Errorf("serve() error = %v", err)
			}
		case <-ctx.Done():
			t.Error("server did not stop before timeout")
		}

		st.Close()
	})

	return srv, conn, addr
}
