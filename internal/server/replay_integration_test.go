package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/persistence"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestAOFReplay_Integration(t *testing.T) {
	aofDir := filepath.Join(t.TempDir(), "data")

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}

	if err := os.Chdir(filepath.Dir(aofDir)); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}

	t.Cleanup(func() {
		_ = os.Chdir(oldDir)
	})

	if err := os.Mkdir("data", 0755); err != nil {
		t.Fatalf("os.Mkdir() error = %v", err)
	}

	aof, err := persistence.NewAOF(persistence.AOFConfig{
		Path:        "appendonly.aof",
		FsyncPolicy: persistence.FsyncAlways,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}

	commands := []protocol.Value{
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("foo"),
			protocol.NewBulkString("bar"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("counter"),
			protocol.NewBulkString("10"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("INCR"),
			protocol.NewBulkString("counter"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("HSET"),
			protocol.NewBulkString("user"),
			protocol.NewBulkString("name"),
			protocol.NewBulkString("David"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("ZADD"),
			protocol.NewBulkString("scores"),
			protocol.NewBulkString("100"),
			protocol.NewBulkString("alice"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("LPUSH"),
			protocol.NewBulkString("list"),
			protocol.NewBulkString("hello"),
			protocol.NewBulkString("world"),
		}),
	}

	for _, command := range commands {
		if err := aof.Append(command); err != nil {
			aof.Close()
			t.Fatalf("AOF.Append() error = %v", err)
		}
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("AOF.Close() error = %v", err)
	}

	s := &Server{
		store: store.New(),
	}

	err = persistence.Replay("appendonly.aof", s.ReplayCommand)
	if err != nil {
		t.Fatalf("Replay() error = %v", err)
	}

	value, ok, err := s.store.Get("foo")
	if err != nil {
		t.Fatalf("Get(foo) error = %v", err)
	}
	if !ok {
		t.Fatal("expected foo to exist")
	}
	if value != "bar" {
		t.Errorf("foo = %q, want %q", value, "bar")
	}

	value, ok, err = s.store.Get("counter")
	if err != nil {
		t.Fatalf("Get(counter) error = %v", err)
	}
	if !ok {
		t.Fatal("expected counter to exist")
	}
	if value != "11" {
		t.Errorf("counter = %q, want %q", value, "11")
	}

	value, ok, err = s.store.HGet("user", "name")
	if err != nil {
		t.Fatalf("HGet() error = %v", err)
	}
	if !ok {
		t.Fatal("expected user.name to exist")
	}
	if value != "David" {
		t.Errorf("user.name = %q, want %q", value, "David")
	}

	score, ok, err := s.store.ZScore("scores", "alice")
	if err != nil {
		t.Fatalf("ZScore() error = %v", err)
	}
	if !ok {
		t.Fatal("expected alice to exist")
	}
	if score != 100 {
		t.Errorf("alice score = %v, want 100", score)
	}

	values, err := s.store.LRange("list", 0, -1)
	if err != nil {
		t.Fatalf("LRange() error = %v", err)
	}

	want := []string{"world", "hello"}

	if len(values) != len(want) {
		t.Fatalf("list length = %d, want %d", len(values), len(want))
	}

	for i := range want {
		if values[i] != want[i] {
			t.Errorf("values[%d] = %q, want %q", i, values[i], want[i])
		}
	}
}

func TestAOF_ReplayStopsOnApplyError(t *testing.T) {
	aofDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(aofDir, "data"), 0o755); err != nil {
		t.Fatalf("create data directory: %v", err)
	}

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(aofDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	fileName := "test.aof"

	aof, err := persistence.NewAOF(persistence.AOFConfig{
		Path:        fileName,
		FsyncPolicy: persistence.FsyncAlways,
	})
	if err != nil {
		t.Fatalf("create AOF: %v", err)
	}

	setName := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	})
	err = aof.Append(setName)
	if err != nil {
		aof.Close()
		t.Fatalf("append SET name command: %v", err)
	}

	invalidCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("INVALID"),
	})
	err = aof.Append(invalidCommand)
	if err != nil {
		aof.Close()
		t.Fatalf("append invalid command: %v", err)
	}

	setLanguage := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("language"),
		protocol.NewBulkString("Go"),
	})
	err = aof.Append(setLanguage)
	if err != nil {
		aof.Close()
		t.Fatalf("append SET language command: %v", err)
	}

	err = aof.Close()
	if err != nil {
		t.Fatalf("close AOF: %v", err)
	}

	st := store.New()
	defer st.Close()

	srv := New("127.0.0.1:0", st, nil)

	err = persistence.Replay(fileName, srv.ReplayCommand)
	if err == nil {
		t.Fatal("expected Replay() to return an error for an invalid command")
	}

	value, ok, err := st.Get("name")
	if err != nil {
		t.Fatalf("get name: %v", err)
	}
	if !ok {
		t.Errorf("name key does not exist")
	}
	if value != "David" {
		t.Errorf("name = %q, want %q", value, "David")
	}

	_, ok, err = st.Get("language")
	if err != nil {
		t.Errorf("get language: %v", err)
	}
	if ok {
		t.Errorf("expected language key not to exist after replay stopped")
	}
}
