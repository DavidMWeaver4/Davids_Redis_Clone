package server

import (
	"math"
	"strconv"
	"testing"
	"time"

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

func TestCommands_Set_Success(t *testing.T) {
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
	s := &Server{
		store: store.New(),
	}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no arguments",
			args: []string{},
		},
		{
			name: "one argument",
			args: []string{"Foo"},
		},
		{
			name: "three arguments",
			args: []string{"Foo", "Bar", "Baz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := set(s, tt.args)

			if got.Type != protocol.Error {
				t.Fatalf("expected Error, got %v", got.Type)
			}
		})
	}
}
func TestCommands_Get_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := get(s, []string{"Foo"})

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}

	if response.Str != "Bar" {
		t.Fatalf("expected %q, got %q", "Bar", response.Str)
	}
}

func TestCommands_Get_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := get(s, []string{"Foo"})

	if response.Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Type)
	}
}

func TestCommands_Get_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := get(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Del_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := deleteCommand(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}

	_, ok := s.store.Get("Foo")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCommands_Del_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := deleteCommand(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_Del_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := deleteCommand(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Exists_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := exists(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_Exists_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := exists(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_Exists_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := exists(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_SetWithTTL_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := set(s, []string{"Foo", "Bar", "EX", "1"})
	if response.Type != protocol.SimpleString {
		t.Fatalf("expected SimpleString, got %v", response.Type)
	}
	if response.Str != "OK" {
		t.Fatalf("expected OK, got %q", response.Str)
	}
	if _, ok := s.store.Get("Foo"); !ok {
		t.Fatal("expected key to exist")
	}
}
func TestCommands_SetWithTTL_InvalidOption(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := set(s, []string{"Foo", "Bar", "PX", "1000"})
	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}
func TestCommands_SetWithTTL_InvalidExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := set(s, []string{"Foo", "Bar", "EX", "abc"})
	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}
func TestCommands_SetWithTTL_NegativeExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := []string{"0", "-1"}
	for _, seconds := range tests {
		t.Run(seconds, func(t *testing.T) {
			response := set(s, []string{"Foo", "Bar", "EX", seconds})
			if response.Type != protocol.Error {
				t.Fatalf("expected Error, got %v", response.Type)
			}
		})
	}
}
func TestCommands_TTL_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 5*time.Second)

	response := ttl(s, []string{"Foo"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int < 4 || response.Int > 5 {
		t.Fatalf("expected TTL around 5, got %d", response.Int)
	}
}
func TestCommands_TTL_NoExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)

	response := ttl(s, []string{"Foo"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != -1 {
		t.Fatalf("expected -1, got %d", response.Int)
	}
}
func TestCommands_TTL_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := ttl(s, []string{"Foo"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != -2 {
		t.Fatalf("expected -2, got %d", response.Int)
	}
}
func TestCommands_TTL_InvaludArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := ttl(s, args)
		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Expire_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)
	response := expire(s, []string{"Foo", "10"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
	ttl := s.store.TTL("Foo")
	if ttl < 9 || ttl > 10 {
		t.Fatalf("expected TTL around 10, got %d", ttl)
	}
}

func TestCommands_Expire_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	response := expire(s, []string{"Foo", "10"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 0 {
		t.Fatalf("expected failure due to missing key, got %d", response.Int)
	}
}
func TestCommands_Expire_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)
	response := expire(s, []string{"Foo", "Bar", "10", "EXP"})
	if response.Type != protocol.Error {
		t.Fatalf("expected error, got %v", response.Type)
	}
}
func TestCommands_Expire_InvalidExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)
	response := expire(s, []string{"Foo", "x10x1x"})
	if response.Type != protocol.Error {
		t.Fatalf("expected error, got %v", response.Type)
	}
}
func TestCommands_Expire_NegativeExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)
	response := expire(s, []string{"Foo", "-100"})
	if response.Type != protocol.Error {
		t.Fatalf("expected error, got %v", response.Type)
	}
}
func TestCommands_Expire_ZeroExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)
	response := expire(s, []string{"Foo", "0"})
	if response.Type != protocol.Error {
		t.Fatalf("expected error, got %v", response.Type)
	}
}

func TestCommands_Persist_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 10*time.Second)
	response := persist(s, []string{"Foo"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
	ttl := s.store.TTL("Foo")
	if ttl != -1 {
		t.Fatalf("expected ttl to be -1, got %d", ttl)
	}
}

func TestCommands_Persist_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	response := persist(s, []string{"Foo"})
	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 0 {
		t.Fatalf("expected failure due to missing key, got %d", response.Int)
	}
}

func TestCommands_Persist_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 10*time.Second)
	response := persist(s, []string{"Foo", "Bar", "10"})
	if response.Type != protocol.Error {
		t.Fatalf("expected error, got %v", response.Type)
	}
}

func TestCommands_Persist_AlreadyPersistent(t *testing.T) {
	s := &Server{
		store: store.New(),
	}
	s.store.Set("Foo", "Bar", 0)

	response := persist(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 0 {
		t.Fatalf("expected 0 for already persistent key, got %d", response.Int)
	}
}
func TestCommands_Incr_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 21 {
		t.Fatalf("expected 21, got %d", response.Int)
	}
}

func TestCommands_Incr_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_Incr_NonInteger(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != "ERR value is not an integer" {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Incr_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := incr(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Incr_Overflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MaxInt64, 10), 0)

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != "ERR increment would cause integer overflow" {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}
func TestCommands_Decr_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 19 {
		t.Fatalf("expected 19, got %d", response.Int)
	}
}

func TestCommands_Decr_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != -1 {
		t.Fatalf("expected -1, got %d", response.Int)
	}
}

func TestCommands_Decr_NonInteger(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != "ERR value is not an integer" {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Decr_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := decr(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Decr_Underflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != "ERR decrement would cause integer overflow" {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}
