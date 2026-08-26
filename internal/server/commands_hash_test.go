package server

import (
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_HSet_Success(t *testing.T) {
	s := &Server{store: store.New()}

	response := hset(s, []string{"Foo", "Bar", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_HSet_Update(t *testing.T) {
	s := &Server{store: store.New()}

	hset(s, []string{"Foo", "Bar", "Baz"})

	response := hset(s, []string{"Foo", "Bar", "NewBaz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0 when updating existing field, got %d", response.Int)
	}
}

func TestCommands_HSet_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "Bar"},
		{"Foo", "Bar", "Baz", "Extra"},
	}

	for _, args := range tests {
		response := hset(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_HGet_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hget(s, []string{"Foo", "Bar"})

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}

	if response.Str != "Baz" {
		t.Fatalf("expected Baz, got %q", response.Str)
	}
}

func TestCommands_HGet_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := hget(s, []string{"Foo", "Bar"})

	if response.Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Type)
	}
}

func TestCommands_HGet_MissingField(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hget(s, []string{"Foo", "Missing"})

	if response.Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Type)
	}
}

func TestCommands_HGet_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "Bar", "Baz"},
	}

	for _, args := range tests {
		response := hget(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_HDel_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.store.HSet("Foo", "Second", "Value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hdel(s, []string{"Foo", "Bar", "Second"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_HDel_MissingFields(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hdel(s, []string{"Foo", "Missing"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_HDel_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	response := hdel(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_HExists_Exists(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hexists(s, []string{"Foo", "Bar"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_HExists_Missing(t *testing.T) {
	s := &Server{store: store.New()}

	response := hexists(s, []string{"Foo", "Bar"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_HExists_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "Bar", "Baz"},
	}

	for _, args := range tests {
		response := hexists(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_HLen_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.store.HSet("Foo", "Second", "Value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hlen(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_HLen_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := hlen(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_HLen_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := hlen(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_HGetAll_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.store.HSet("Foo", "Second", "Value")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := hgetall(s, []string{"Foo"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if len(response.Array) != 4 {
		t.Fatalf("expected 4 array elements, got %d", len(response.Array))
	}

	pairs := make(map[string]string)

	for i := 0; i < len(response.Array); i += 2 {
		field := response.Array[i]
		value := response.Array[i+1]

		if field.Type != protocol.BulkString {
			t.Fatalf("expected field to be BulkString, got %v", field.Type)
		}

		if value.Type != protocol.BulkString {
			t.Fatalf("expected value to be BulkString, got %v", value.Type)
		}

		pairs[field.Str] = value.Str
	}

	expected := map[string]string{
		"Bar":    "Baz",
		"Second": "Value",
	}

	for field, expectedValue := range expected {
		actualValue, ok := pairs[field]
		if !ok {
			t.Fatalf("expected field %q to exist", field)
		}

		if actualValue != expectedValue {
			t.Fatalf(
				"expected field %q to have value %q, got %q",
				field,
				expectedValue,
				actualValue,
			)
		}
	}
}

func TestCommands_HGetAll_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := hgetall(s, []string{"Foo"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if len(response.Array) != 0 {
		t.Fatalf("expected empty array, got %d elements", len(response.Array))
	}
}

func TestCommands_HGetAll_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := hgetall(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}
