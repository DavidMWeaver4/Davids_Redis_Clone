package server

import (
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_Ping_Success(t *testing.T) {
	s := &Server{}

	got := ping(s, nil)

	if got.Type != protocol.SimpleString {
		t.Fatalf("expected type %v, got %v", protocol.SimpleString, got.Type)
	}

	if got.Str != "PONG" {
		t.Fatalf("expected %q, got %q", "PONG", got.Str)
	}
}

func TestCommands_Set_Sucess(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	response := set(s, []string{"Foo", "Bar"})
	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}
	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}
	value, ok := s.store.Get("Foo")
	if !ok {
		t.Fatal("expected key to exist after SET")
	}
	if value != "Bar" {
		t.Fatalf("expected value %q, got %q", "Bar", value)
	}
}
func TestCommands_Set_InvalidArgumentCount(t *testing.T) {
	t.Skip()
}
func TestCommands_Get_Success(t *testing.T) {
	t.Skip()
}
func TestCommands_Get_MissingKey(t *testing.T) {
	t.Skip()
}
func TestCommands_Del_Success_(t *testing.T) {
	t.Skip()
}
func TestCommands_Del_MissingKey(t *testing.T) {
	t.Skip()
}
func TestCommands_Exists_Success(t *testing.T) {
	t.Skip()
}
func TestCommands_Exists_MissingKey(t *testing.T) {
	t.Skip()
}
