package server

import (
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_Append_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := appendCommand(s, []string{"Foo", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 6 {
		t.Fatalf("expected 6, got %d", response.Int)
	}

	value, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "BarBaz" {
		t.Fatalf("expected %q, got %q", "BarBaz", value)
	}
}

func TestCommands_Append_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := appendCommand(s, []string{"Foo", "Bar"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 3 {
		t.Fatalf("expected 3, got %d", response.Int)
	}
}

func TestCommands_Append_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "Bar", "Baz"},
	}

	for _, args := range tests {
		response := appendCommand(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}
func TestCommands_Strlen_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := strlen(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 3 {
		t.Fatalf("expected 3, got %d", response.Int)
	}
}

func TestCommands_Strlen_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := strlen(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_Strlen_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := strlen(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}
func TestCommands_Setnx_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := setnx(s, []string{"Foo", "Bar"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_Setnx_ExistingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := setnx(s, []string{"Foo", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}

	value, _, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if value != "Bar" {
		t.Fatalf("expected original value %q, got %q", "Bar", value)
	}
}

func TestCommands_Setnx_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "Bar", "Baz"},
	}

	for _, args := range tests {
		response := setnx(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}
func TestCommands_Mget_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)
	s.store.Set("Baz", "Qux", 0)

	response := mget(s, []string{"Foo", "Missing", "Baz"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if len(response.Array) != 3 {
		t.Fatalf("expected 3 values, got %d", len(response.Array))
	}

	if response.Array[0].Type != protocol.BulkString ||
		response.Array[0].Str != "Bar" {
		t.Fatalf("unexpected first value: %+v", response.Array[0])
	}

	if response.Array[1].Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Array[1].Type)
	}

	if response.Array[2].Type != protocol.BulkString ||
		response.Array[2].Str != "Qux" {
		t.Fatalf("unexpected third value: %+v", response.Array[2])
	}
}

func TestCommands_Mget_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := mget(s, []string{})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}
func TestCommands_Mset_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := mset(s, []string{
		"Foo", "Bar",
		"Baz", "Qux",
	})

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}

	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}

	foo, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok || foo != "Bar" {
		t.Fatalf("expected Foo=Bar, got %q", foo)
	}

	baz, ok, err := s.store.Get("Baz")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok || baz != "Qux" {
		t.Fatalf("expected Baz=Qux, got %q", baz)
	}
}

func TestCommands_Mset_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "Bar", "Baz"},
	}

	for _, args := range tests {
		response := mset(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}
func TestCommands_Append_PreservesTTL(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 5*time.Second)

	response := appendCommand(s, []string{"Foo", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	ttl := s.store.TTL("Foo")
	if ttl < 4 || ttl > 5 {
		t.Fatalf("expected TTL around 5, got %d", ttl)
	}

	value, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "BarBaz" {
		t.Fatalf("expected %q, got %q", "BarBaz", value)
	}
}
func TestCommands_Setnx_ExpiredKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	response := setnx(s, []string{"Foo", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}

	value, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", value)
	}
}
