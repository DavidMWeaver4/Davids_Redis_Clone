package server

import (
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

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
	value, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}
	if value != "Bar" {
		t.Fatalf("expected value %q, got %q", "Bar", value)
	}
	ttl := s.store.TTL("Foo")
	if ttl < 0 || ttl > 1 {
		t.Fatalf("expected TTL around 1 second, got %d", ttl)
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
func TestCommands_TTL_InvalidArgumentCount(t *testing.T) {
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

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}

	_, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected key to be deleted")
	}
}
func TestCommands_Expire_ZeroExpiration(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := expire(s, []string{"Foo", "0"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}

	_, ok, err := s.store.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected key to be deleted")
	}
}
func TestCommands_Expire_ExpiredKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	response := expire(s, []string{"Foo", "10"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0 for expired key, got %d", response.Int)
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
func TestCommands_Persist_ExpiredKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", time.Millisecond)
	time.Sleep(2 * time.Millisecond)

	response := persist(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0 for expired key, got %d", response.Int)
	}
}
