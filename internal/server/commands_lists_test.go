package server

import (
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_LPush_Success(t *testing.T) {
	s := &Server{store: store.New()}

	response := lpush(s, []string{"Foo", "Bar", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_LPush_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
	}

	for _, args := range tests {
		response := lpush(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_RPush_Success(t *testing.T) {
	s := &Server{store: store.New()}

	response := rpush(s, []string{"Foo", "Bar", "Baz"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_RPush_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
	}

	for _, args := range tests {
		response := rpush(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LPop_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.LPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := lpop(s, []string{"Foo"})

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}

	if response.Str != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", response.Str)
	}
}

func TestCommands_LPop_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := lpop(s, []string{"Foo"})

	if response.Type != protocol.Null {
		t.Fatalf("expected NullBulkString, got %v", response.Type)
	}
}

func TestCommands_LPop_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := lpop(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_RPop_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := rpop(s, []string{"Foo"})

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}

	if response.Str != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", response.Str)
	}
}

func TestCommands_RPop_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := rpop(s, []string{"Foo"})

	if response.Type != protocol.Null {
		t.Fatalf("expected NullBulkString, got %v", response.Type)
	}
}

func TestCommands_RPop_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := rpop(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LLen_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := llen(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_LLen_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := llen(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_LLen_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := llen(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LRange_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar", "Baz", "Qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := lrange(s, []string{"Foo", "1", "2"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if len(response.Array) != 2 {
		t.Fatalf("expected 2 values, got %d", len(response.Array))
	}

	expected := []string{"Baz", "Qux"}

	for i, value := range expected {
		if response.Array[i].Type != protocol.BulkString {
			t.Fatalf(
				"expected element %d to be BulkString, got %v",
				i,
				response.Array[i].Type,
			)
		}

		if response.Array[i].Str != value {
			t.Fatalf(
				"expected element %d to be %q, got %q",
				i,
				value,
				response.Array[i].Str,
			)
		}
	}
}

func TestCommands_LRange_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := lrange(s, []string{"Foo", "0", "-1"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if len(response.Array) != 0 {
		t.Fatalf("expected empty array, got %d values", len(response.Array))
	}
}

func TestCommands_LRange_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "0"},
		{"Foo", "0", "1", "2"},
	}

	for _, args := range tests {
		response := lrange(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LRange_InvalidStart(t *testing.T) {
	s := &Server{store: store.New()}

	response := lrange(s, []string{"Foo", "abc", "2"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LRange_InvalidEnd(t *testing.T) {
	s := &Server{store: store.New()}

	response := lrange(s, []string{"Foo", "0", "abc"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LIndex_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar", "Baz", "Qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := lindex(s, []string{"Foo", "1"})

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}

	if response.Str != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", response.Str)
	}
}

func TestCommands_LIndex_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := lindex(s, []string{"Foo", "0"})

	if response.Type != protocol.Null {
		t.Fatalf("expected NullBulkString, got %v", response.Type)
	}
}

func TestCommands_LIndex_InvalidIndex(t *testing.T) {
	s := &Server{store: store.New()}

	response := lindex(s, []string{"Foo", "abc"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LIndex_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "1", "Bar"},
	}

	for _, args := range tests {
		response := lindex(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LSet_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := lset(s, []string{"Foo", "1", "Qux"})

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}

	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}

	value, ok, err := s.store.LIndex("Foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected value to exist")
	}

	if value != "Qux" {
		t.Fatalf("expected %q, got %q", "Qux", value)
	}
}

func TestCommands_LSet_InvalidIndex(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := lset(s, []string{"Foo", "10", "Baz"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LSet_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := lset(s, []string{"Foo", "0", "Bar"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LSet_InvalidArgumentCount(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "0"},
		{"Foo", "0", "Bar", "Baz"},
	}

	for _, args := range tests {
		response := lset(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LSet_InvalidInteger(t *testing.T) {
	s := &Server{store: store.New()}

	response := lset(s, []string{"Foo", "abc", "Bar"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LTrim_Success(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("Foo", "Bar", "Baz", "Qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	response := ltrim(s, []string{"Foo", "1", "2"})

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}

	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}

	values, err := s.store.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"Baz", "Qux"}

	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestCommands_LTrim_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := ltrim(s, []string{"Foo", "0", "-1"})

	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}

	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}
}

func TestCommands_LTrim_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "0"},
		{"Foo", "0", "1", "2"},
	}

	for _, args := range tests {
		response := ltrim(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_LTrim_InvalidStart(t *testing.T) {
	s := &Server{store: store.New()}

	response := ltrim(s, []string{"Foo", "abc", "2"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_LTrim_InvalidEnd(t *testing.T) {
	s := &Server{store: store.New()}

	response := ltrim(s, []string{"Foo", "0", "abc"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}
