package persistence

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func TestAOF_RESPRoundTrip(t *testing.T) {
	path := "round-trip.aof"
	fullPath := filepath.Join(aofBaseDir, path)

	cleanupAOFFile(t, fullPath)

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncAlways,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	if err := aof.Append(command); err != nil {
		aof.Close()
		t.Fatalf("Append() error = %v", err)
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var got protocol.Value

	err = Replay(path, func(value protocol.Value) error {
		got = value
		return nil
	})
	if err != nil {
		t.Fatalf("Replay() error = %v", err)
	}

	if got.Type != protocol.Array {
		t.Fatalf("value type = %v, want Array", got.Type)
	}

	if len(got.Array) != 3 {
		t.Fatalf("array length = %d, want 3", len(got.Array))
	}

	want := []string{"SET", "foo", "bar"}

	for i, expected := range want {
		if got.Array[i].Str != expected {
			t.Errorf("value[%d] = %q, want %q", i, got.Array[i].Str, expected)
		}
	}
}

func TestReplay_EmptyAOF(t *testing.T) {
	path := "empty.aof"

	file, _ := createAOFTestFile(t, path)

	if err := file.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	called := false

	err := Replay(path, func(protocol.Value) error {
		called = true
		return nil
	})

	if err != nil {
		t.Fatalf("Replay() error = %v, want nil", err)
	}

	if called {
		t.Fatal("apply should not be called for an empty AOF")
	}
}

func TestReplay_StopsAfterApplyFails(t *testing.T) {
	path := "apply-error.aof"
	fullPath := filepath.Join(aofBaseDir, path)

	cleanupAOFFile(t, fullPath)

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncAlways,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}

	commands := []protocol.Value{
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("first"),
			protocol.NewBulkString("value"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("second"),
			protocol.NewBulkString("value"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("third"),
			protocol.NewBulkString("value"),
		}),
	}

	for _, command := range commands {
		if err := aof.Append(command); err != nil {
			aof.Close()
			t.Fatalf("Append() error = %v", err)
		}
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	wantErr := errors.New("apply failed")
	callCount := 0

	err = Replay(path, func(protocol.Value) error {
		callCount++

		if callCount == 2 {
			return wantErr
		}

		return nil
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Replay() error = %v, want %v", err, wantErr)
	}

	if callCount != 2 {
		t.Errorf("apply called %d times, want 2", callCount)
	}
}

func TestAOF_ReopenExistingAndAppend(t *testing.T) {
	path := "reopen.aof"
	fullPath := filepath.Join(aofBaseDir, path)

	cleanupAOFFile(t, fullPath)

	firstCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncAlways,
	})
	if err != nil {
		t.Fatalf("first NewAOF() error = %v", err)
	}

	if err := aof.Append(firstCommand); err != nil {
		aof.Close()
		t.Fatalf("first Append() error = %v", err)
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}

	secondCommand := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("hello"),
		protocol.NewBulkString("world"),
	})

	aof, err = NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncAlways,
	})
	if err != nil {
		t.Fatalf("second NewAOF() error = %v", err)
	}

	if err := aof.Append(secondCommand); err != nil {
		aof.Close()
		t.Fatalf("second Append() error = %v", err)
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}

	var got []protocol.Value

	err = Replay(path, func(command protocol.Value) error {
		got = append(got, command)
		return nil
	})
	if err != nil {
		t.Fatalf("Replay() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("replayed commands = %d, want 2", len(got))
	}

	if got[0].Array[0].Str != "SET" {
		t.Errorf("first command = %q, want SET", got[0].Array[0].Str)
	}

	if got[0].Array[1].Str != "foo" {
		t.Errorf("first key = %q, want foo", got[0].Array[1].Str)
	}

	if got[1].Array[0].Str != "SET" {
		t.Errorf("second command = %q, want SET", got[1].Array[0].Str)
	}

	if got[1].Array[1].Str != "hello" {
		t.Errorf("second key = %q, want hello", got[1].Array[1].Str)
	}
	if got[0].Array[2].Str != "bar" {
		t.Errorf("first value = %q, want bar", got[0].Array[2].Str)
	}
	if got[1].Array[2].Str != "world" {
		t.Errorf("second value = %q, want world", got[1].Array[2].Str)
	}
}
