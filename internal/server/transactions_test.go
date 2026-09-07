package server

import (
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_Multi(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	response := multi(s, c, []string{})

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}
	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}
	if !c.inTransaction {
		t.Fatal("expected client to be in transaction")
	}
}

func TestCommands_Multi_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	response := multi(s, c, []string{"extra"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
	if c.inTransaction {
		t.Fatal("expected client to remain outside transaction")
	}
}

func TestCommands_Multi_NestedTransaction(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})
	response := multi(s, c, []string{})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
	if !c.inTransaction {
		t.Fatal("expected client to remain in transaction")
	}
}

func TestCommands_QueueCommand(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	response := s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	}))

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}
	if response.Str != "QUEUED" {
		t.Fatalf("expected QUEUED, got %q", response.Str)
	}
	if len(c.queue) != 1 {
		t.Fatalf("expected 1 queued command, got %d", len(c.queue))
	}

	_, ok, err := s.store.Get("name")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected queued command to not modify store")
	}
}

func TestCommands_QueueMultipleCommands(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	commands := []protocol.Value{
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("name"),
			protocol.NewBulkString("David"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("language"),
			protocol.NewBulkString("Go"),
		}),
	}

	for _, command := range commands {
		response := s.execute(c, command)

		if response.Type != protocol.SimpleString {
			t.Fatalf("expected SimpleString, got %v", response.Type)
		}
		if response.Str != "QUEUED" {
			t.Fatalf("expected QUEUED, got %q", response.Str)
		}
	}

	if len(c.queue) != 2 {
		t.Fatalf("expected 2 queued commands, got %d", len(c.queue))
	}

	_, ok, err := s.store.Get("name")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected name to not exist before EXEC")
	}
	_, ok, err = s.store.Get("language")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected language to not exist before EXEC")
	}
}

func TestCommands_Exec(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	}))

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("language"),
		protocol.NewBulkString("Go"),
	}))

	response := exec(s, c, []string{})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}
	if len(response.Array) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(response.Array))
	}

	for i, value := range response.Array {
		if value.Type != protocol.SimpleString {
			t.Fatalf("response[%d]: expected SimpleString, got %v", i, value.Type)
		}
		if value.Str != "OK" {
			t.Fatalf("response[%d]: expected OK, got %q", i, value.Str)
		}
	}

	if c.inTransaction {
		t.Fatal("expected transaction to be complete")
	}
	if len(c.queue) != 0 {
		t.Fatalf("expected queue to be empty, got %d", len(c.queue))
	}

	value, ok, err := s.store.Get("name")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected name to exist after EXEC")
	}
	if value != "David" {
		t.Fatalf("expected David, got %q", value)
	}

	value, ok, err = s.store.Get("language")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected language to exist after EXEC")
	}
	if value != "Go" {
		t.Fatalf("expected Go, got %q", value)
	}
}

func TestCommands_Exec_EmptyTransaction(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	response := exec(s, c, []string{})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}
	if len(response.Array) != 0 {
		t.Fatalf("expected empty response array, got %d values", len(response.Array))
	}
	if c.inTransaction {
		t.Fatal("expected transaction to be complete")
	}
}

func TestCommands_Exec_OutsideTransaction(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	response := exec(s, c, []string{})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_Exec_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	response := exec(s, c, []string{"extra"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
	if !c.inTransaction {
		t.Fatal("expected transaction to remain active")
	}
}

func TestCommands_Exec_ExecutesCommandsInOrder(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("count"),
		protocol.NewBulkString("1"),
	}))

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("INCR"),
		protocol.NewBulkString("count"),
	}))

	response := exec(s, c, []string{})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}
	if len(response.Array) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(response.Array))
	}

	if response.Array[0].Type != protocol.SimpleString {
		t.Fatalf("response[0]: expected SimpleString, got %v", response.Array[0].Type)
	}
	if response.Array[0].Str != "OK" {
		t.Fatalf("response[0]: expected OK, got %q", response.Array[0].Str)
	}

	if response.Array[1].Type != protocol.Integer {
		t.Fatalf("response[1]: expected Integer, got %v", response.Array[1].Type)
	}
	if response.Array[1].Int != 2 {
		t.Fatalf("response[1]: expected 2, got %d", response.Array[1].Int)
	}

	value, ok, err := s.store.Get("count")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected count to exist")
	}
	if value != "2" {
		t.Fatalf("expected count to be 2, got %q", value)
	}
}

func TestCommands_Exec_CommandErrorDoesNotStopTransaction(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("count"),
		protocol.NewBulkString("not-an-integer"),
	}))

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("INCR"),
		protocol.NewBulkString("count"),
	}))

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("finished"),
		protocol.NewBulkString("yes"),
	}))

	response := exec(s, c, []string{})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}
	if len(response.Array) != 3 {
		t.Fatalf("expected 3 responses, got %d", len(response.Array))
	}

	if response.Array[0].Type != protocol.SimpleString {
		t.Fatalf("response[0]: expected SimpleString, got %v", response.Array[0].Type)
	}

	if response.Array[1].Type != protocol.Error {
		t.Fatalf("response[1]: expected Error, got %v", response.Array[1].Type)
	}

	if response.Array[2].Type != protocol.SimpleString {
		t.Fatalf("response[2]: expected SimpleString, got %v", response.Array[2].Type)
	}
	if response.Array[2].Str != "OK" {
		t.Fatalf("response[2]: expected OK, got %q", response.Array[2].Str)
	}

	value, ok, err := s.store.Get("finished")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected later commands to execute after an error")
	}
	if value != "yes" {
		t.Fatalf("expected yes, got %q", value)
	}
}

func TestCommands_Discard(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	s.execute(c, protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	}))

	response := discard(s, c, []string{})

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}
	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}
	if c.inTransaction {
		t.Fatal("expected transaction to be inactive")
	}
	if len(c.queue) != 0 {
		t.Fatalf("expected queue to be empty, got %d", len(c.queue))
	}

	_, ok, err := s.store.Get("name")
	if err != nil {
		t.Fatalf("got unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected discarded command to not modify store")
	}
}

func TestCommands_Discard_OutsideTransaction(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	response := discard(s, c, []string{})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_Discard_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}
	c := newTestClient()

	multi(s, c, []string{})

	response := discard(s, c, []string{"extra"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
	if !c.inTransaction {
		t.Fatal("expected transaction to remain active")
	}
}
