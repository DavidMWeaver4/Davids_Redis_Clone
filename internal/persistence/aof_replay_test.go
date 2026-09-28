package persistence

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func createAOFTestFile(t *testing.T, path string) (*os.File, string) {
	t.Helper()

	if err := os.MkdirAll(aofBaseDir, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	fullPath := filepath.Join(aofBaseDir, path)

	t.Cleanup(func() {
		_ = os.Remove(fullPath)
		_ = os.Remove(aofBaseDir)
	})

	file, err := os.Create(fullPath)
	if err != nil {
		t.Fatalf("os.Create() error = %v", err)
	}

	return file, fullPath
}

func TestReplay_FileDoesNotExist(t *testing.T) {
	err := Replay("does-not-exist.aof", func(protocol.Value) error {
		t.Fatal("apply should not be called")
		return nil
	})

	if err != nil {
		t.Fatalf("Replay() error = %v, want nil", err)
	}
}

func TestReplay_SingleCommand(t *testing.T) {
	path := "replay-single.aof"

	file, _ := createAOFTestFile(t, path)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	if err := protocol.Write(file, command); err != nil {
		file.Close()
		t.Fatalf("protocol.Write() error = %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("file.Close() error = %v", err)
	}

	var got protocol.Value

	err := Replay(path, func(value protocol.Value) error {
		got = value
		return nil
	})

	if err != nil {
		t.Fatalf("Replay() error = %v", err)
	}

	if got.Type != protocol.Array {
		t.Fatalf("command type = %v, want Array", got.Type)
	}

	if len(got.Array) != 3 {
		t.Fatalf("command length = %d, want 3", len(got.Array))
	}

	if got.Array[0].Str != "SET" {
		t.Errorf("command name = %q, want %q", got.Array[0].Str, "SET")
	}

	if got.Array[1].Str != "foo" {
		t.Errorf("key = %q, want %q", got.Array[1].Str, "foo")
	}

	if got.Array[2].Str != "bar" {
		t.Errorf("value = %q, want %q", got.Array[2].Str, "bar")
	}
}

func TestReplay_MultipleCommands(t *testing.T) {
	path := "replay-multiple.aof"
	file, _ := createAOFTestFile(t, path)

	commands := []protocol.Value{
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("foo"),
			protocol.NewBulkString("bar"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("INCR"),
			protocol.NewBulkString("counter"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("DEL"),
			protocol.NewBulkString("foo"),
		}),
	}

	for _, command := range commands {
		if err := protocol.Write(file, command); err != nil {
			file.Close()
			t.Fatalf("protocol.Write() error = %v", err)
		}
	}

	if err := file.Close(); err != nil {
		t.Fatalf("file.Close() error = %v", err)
	}

	var got []protocol.Value

	err := Replay(path, func(value protocol.Value) error {
		got = append(got, value)
		return nil
	})

	if err != nil {
		t.Fatalf("Replay() error = %v", err)
	}

	if len(got) != len(commands) {
		t.Fatalf("number of commands = %d, want %d", len(got), len(commands))
	}

	for i, command := range got {
		if command.Type != protocol.Array {
			t.Errorf("command %d type = %v, want Array", i, command.Type)
		}
	}

	if got[0].Array[0].Str != "SET" {
		t.Errorf("first command = %q, want SET", got[0].Array[0].Str)
	}

	if got[1].Array[0].Str != "INCR" {
		t.Errorf("second command = %q, want INCR", got[1].Array[0].Str)
	}

	if got[2].Array[0].Str != "DEL" {
		t.Errorf("third command = %q, want DEL", got[2].Array[0].Str)
	}
}

func TestReplay_MalformedRESP(t *testing.T) {
	path := "replay-malformed.aof"
	file, _ := createAOFTestFile(t, path)

	_, err := file.WriteString("*2\r\n$4\r\nSET\r\n$3\r\nfoo")
	if err != nil {
		file.Close()
		t.Fatalf("file.WriteString() error = %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("file.Close() error = %v", err)
	}

	err = Replay(path, func(protocol.Value) error {
		t.Fatal("apply should not be called for malformed RESP")
		return nil
	})

	if err == nil {
		t.Fatal("Replay() error = nil, want error")
	}
}

func TestReplay_ApplyError(t *testing.T) {
	path := "replay-apply-error.aof"
	file, _ := createAOFTestFile(t, path)

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("foo"),
		protocol.NewBulkString("bar"),
	})

	if err := protocol.Write(file, command); err != nil {
		file.Close()
		t.Fatalf("protocol.Write() error = %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("file.Close() error = %v", err)
	}

	wantErr := errors.New("apply failed")

	err := Replay(path, func(protocol.Value) error {
		return wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Replay() error = %v, want %v", err, wantErr)
	}
}

func TestReplay_InvalidPath(t *testing.T) {
	tests := []struct {
		name string
		path string
	}{
		{
			name: "empty path",
			path: "",
		},
		{
			name: "whitespace path",
			path: "   ",
		},
		{
			name: "directory traversal",
			path: "../replay.aof",
		},
		{
			name: "absolute path",
			path: "/tmp/replay.aof",
		},
		{
			name: "null byte",
			path: "replay\x00.aof",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Replay(tt.path, func(protocol.Value) error {
				t.Fatal("apply should not be called")
				return nil
			})

			if err == nil {
				t.Fatal("Replay() error = nil, want error")
			}
		})
	}
}
