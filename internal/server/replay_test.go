package server

import (
	"testing"
	"time"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestReplayCommand_Set(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "bar" {
		t.Errorf("value = %q, want %q", value, "bar")
	}
}

func TestReplayCommand_Del(t *testing.T) {
	s := &Server{store: store.New()}
	s.store.Set("foo", "bar", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("DEL"),
		protocol.NewBulkString("foo"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	_, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestReplayCommand_Incr(t *testing.T) {
	s := &Server{store: store.New()}
	s.store.Set("counter", "10", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("INCR"),
		protocol.NewBulkString("counter"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected counter to exist")
	}

	if value != "11" {
		t.Errorf("value = %q, want %q", value, "11")
	}
}

func TestReplayCommand_IncrBy(t *testing.T) {
	s := &Server{store: store.New()}
	s.store.Set("counter", "10", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("INCRBY"),
		protocol.NewBulkString("counter"),
		protocol.NewBulkString("5"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected counter to exist")
	}

	if value != "15" {
		t.Errorf("value = %q, want %q", value, "15")
	}
}

func TestReplayCommand_Append(t *testing.T) {
	s := &Server{store: store.New()}
	s.store.Set("foo", "hello", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("APPEND"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString(" world"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "hello world" {
		t.Errorf("value = %q, want %q", value, "hello world")
	}
}

func TestReplayCommand_MSet(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("MSET"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
		protocol.NewBulkString("hello"),
		protocol.NewBulkString("world"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || value != "bar" {
		t.Errorf("foo = %q, want %q", value, "bar")
	}

	value, ok, err = s.store.Get("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || value != "world" {
		t.Errorf("hello = %q, want %q", value, "world")
	}
}

func TestReplayCommand_List(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("LPUSH"),
		protocol.NewBulkString("mylist"),
		protocol.NewBulkString("hello"),
		protocol.NewBulkString("world"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	values, err := s.store.LRange("mylist", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	if len(values) != 2 {
		t.Fatalf("list length = %d, want 2", len(values))
	}

	if values[0] != "world" || values[1] != "hello" {
		t.Errorf("list = %v, want [world hello]", values)
	}
}

func TestReplayCommand_HSet(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("HSET"),
		protocol.NewBulkString("user"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected hash field to exist")
	}

	if value != "David" {
		t.Errorf("value = %q, want %q", value, "David")
	}
}

func TestReplayCommand_HDel(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("HSet() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("HDEL"),
		protocol.NewBulkString("user"),
		protocol.NewBulkString("name"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	_, ok, err := s.store.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected hash field to be deleted")
	}
}

func TestReplayCommand_ZAdd(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("ZADD"),
		protocol.NewBulkString("scores"),
		protocol.NewBulkString("10.5"),
		protocol.NewBulkString("alice"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	score, ok, err := s.store.ZScore("scores", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected sorted-set member to exist")
	}

	if score != 10.5 {
		t.Errorf("score = %v, want 10.5", score)
	}
}

func TestReplayCommand_ZRem(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.ZAdd("scores", 10, "alice")
	if err != nil {
		t.Fatalf("ZAdd() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("ZREM"),
		protocol.NewBulkString("scores"),
		protocol.NewBulkString("alice"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	_, ok, err := s.store.ZScore("scores", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected member to be removed")
	}
}

func TestReplayCommand_ZIncrBy(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.ZAdd("scores", 10, "alice")
	if err != nil {
		t.Fatalf("ZAdd() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("ZINCRBY"),
		protocol.NewBulkString("scores"),
		protocol.NewBulkString("2.5"),
		protocol.NewBulkString("alice"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	score, ok, err := s.store.ZScore("scores", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected member to exist")
	}

	if score != 12.5 {
		t.Errorf("score = %v, want 12.5", score)
	}
}
func TestReplayCommand_InvalidCommandType(t *testing.T) {
	s := &Server{store: store.New()}

	tests := []protocol.Value{
		protocol.NewSimpleString("SET"),
		protocol.NewInteger(42),
		protocol.NewBulkString("SET"),
	}

	for _, command := range tests {
		err := s.ReplayCommand(command)
		if err == nil {
			t.Fatal("ReplayCommand() error = nil, want error")
		}
	}
}

func TestReplayCommand_EmptyCommand(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{})

	err := s.ReplayCommand(command)
	if err == nil {
		t.Fatal("ReplayCommand() error = nil, want error")
	}
}

func TestReplayCommand_InvalidCommandTypeInArray(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewInteger(42),
		protocol.NewBulkString("foo"),
	})

	err := s.ReplayCommand(command)
	if err == nil {
		t.Fatal("ReplayCommand() error = nil, want error")
	}
}

func TestReplayCommand_InvalidArgumentType(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewInteger(42),
		protocol.NewBulkString("bar"),
	})

	err := s.ReplayCommand(command)
	if err == nil {
		t.Fatal("ReplayCommand() error = nil, want error")
	}
}

func TestReplayCommand_InvalidCommand(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("NOTACOMMAND"),
	})

	err := s.ReplayCommand(command)
	if err == nil {
		t.Fatal("ReplayCommand() error = nil, want error")
	}
}

func TestReplayCommand_ReadCommand(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("GET"),
		protocol.NewBulkString("foo"),
	})

	err := s.ReplayCommand(command)
	if err == nil {
		t.Fatal("ReplayCommand() error = nil, want error")
	}
}
func TestReplayCommand_Expire(t *testing.T) {
	s := &Server{store: store.New()}

	s.store.Set("foo", "bar", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("EXPIRE"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("60"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	ttl := s.store.TTL("foo")

	if ttl <= 0 {
		t.Errorf("TTL = %v, want positive TTL", ttl)
	}
}

func TestReplayCommand_Persist(t *testing.T) {
	s := &Server{store: store.New()}

	s.store.Set("foo", "bar", 0)
	s.store.Expire("foo", 60*time.Second)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("PERSIST"),
		protocol.NewBulkString("foo"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	ttl := s.store.TTL("foo")
	if ttl != -1 {
		t.Errorf("TTL = %v, want -1", ttl)
	}
}

func TestReplayCommand_Decr(t *testing.T) {
	s := &Server{store: store.New()}

	s.store.Set("counter", "10", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("DECR"),
		protocol.NewBulkString("counter"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected counter to exist")
	}

	if value != "9" {
		t.Errorf("value = %q, want %q", value, "9")
	}
}

func TestReplayCommand_DecrBy(t *testing.T) {
	s := &Server{store: store.New()}

	s.store.Set("counter", "10", 0)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("DECRBY"),
		protocol.NewBulkString("counter"),
		protocol.NewBulkString("3"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("counter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected counter to exist")
	}

	if value != "7" {
		t.Errorf("value = %q, want %q", value, "7")
	}
}

func TestReplayCommand_SetNX(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SETNX"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "bar" {
		t.Errorf("value = %q, want %q", value, "bar")
	}
}

func TestReplayCommand_RPush(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("RPUSH"),
		protocol.NewBulkString("mylist"),
		protocol.NewBulkString("hello"),
		protocol.NewBulkString("world"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	values, err := s.store.LRange("mylist", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	want := []string{"hello", "world"}

	if len(values) != len(want) {
		t.Fatalf("list length = %d, want %d", len(values), len(want))
	}

	for i := range want {
		if values[i] != want[i] {
			t.Errorf("values[%d] = %q, want %q", i, values[i], want[i])
		}
	}
}

func TestReplayCommand_LPop(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.LPush("mylist", "hello", "world")
	if err != nil {
		t.Fatalf("LPush() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("LPOP"),
		protocol.NewBulkString("mylist"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	values, err := s.store.LRange("mylist", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	if len(values) != 1 {
		t.Fatalf("list length = %d, want 1", len(values))
	}

	if values[0] != "hello" {
		t.Errorf("value = %q, want %q", values[0], "hello")
	}
}

func TestReplayCommand_RPop(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("mylist", "hello", "world")
	if err != nil {
		t.Fatalf("RPush() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("RPOP"),
		protocol.NewBulkString("mylist"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	values, err := s.store.LRange("mylist", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	if len(values) != 1 {
		t.Fatalf("list length = %d, want 1", len(values))
	}

	if values[0] != "hello" {
		t.Errorf("value = %q, want %q", values[0], "hello")
	}
}

func TestReplayCommand_LSet(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("mylist", "hello", "world")
	if err != nil {
		t.Fatalf("RPush() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("LSET"),
		protocol.NewBulkString("mylist"),
		protocol.NewBulkString("1"),
		protocol.NewBulkString("redis"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	values, err := s.store.LRange("mylist", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	if values[1] != "redis" {
		t.Errorf("values[1] = %q, want %q", values[1], "redis")
	}
}

func TestReplayCommand_LTrim(t *testing.T) {
	s := &Server{store: store.New()}

	_, err := s.store.RPush("mylist", "one", "two", "three")
	if err != nil {
		t.Fatalf("RPush() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("LTRIM"),
		protocol.NewBulkString("mylist"),
		protocol.NewBulkString("0"),
		protocol.NewBulkString("1"),
	})

	err = s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	values, err := s.store.LRange("mylist", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	want := []string{"one", "two"}

	if len(values) != len(want) {
		t.Fatalf("list length = %d, want %d", len(values), len(want))
	}

	for i := range want {
		if values[i] != want[i] {
			t.Errorf("values[%d] = %q, want %q", i, values[i], want[i])
		}
	}
}
func TestReplayCommand_CaseInsensitive(t *testing.T) {
	s := &Server{store: store.New()}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("set"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	err := s.ReplayCommand(command)
	if err != nil {
		t.Fatalf("ReplayCommand() error = %v", err)
	}

	value, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "bar" {
		t.Errorf("value = %q, want %q", value, "bar")
	}
}
