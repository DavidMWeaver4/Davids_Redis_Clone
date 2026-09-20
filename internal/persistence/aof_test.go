package persistence

import (
	"os"
	"sync"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func TestNewAOF(t *testing.T) {
	path := "test.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncAlways,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}
	defer aof.Close()

	if aof.file == nil {
		t.Fatal("expected AOF file to be open")
	}

	if aof.policy != FsyncAlways {
		t.Fatalf("policy = %v, want %v", aof.policy, FsyncAlways)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected AOF file to exist: %v", err)
	}
}

func TestNewAOF_DefaultFsyncPolicy(t *testing.T) {
	path := "test_default.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path: path,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}
	defer aof.Close()

	if aof.policy != FsyncEverySec {
		t.Fatalf("policy = %v, want %v", aof.policy, FsyncEverySec)
	}
}

func TestNewAOF_InvalidPath(t *testing.T) {
	_, err := NewAOF(AOFConfig{})

	if err != ErrInvalidPath {
		t.Fatalf("error = %v, want %v", err, ErrInvalidPath)
	}
}

func TestNewAOF_InvalidFsyncPolicy(t *testing.T) {
	path := "test_invalid_fsync.aof"

	_, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: "invalid",
	})

	if err != ErrInvalidFsync {
		t.Fatalf("error = %v, want %v", err, ErrInvalidFsync)
	}
}

func TestNewAOF_DirectoryTraversal(t *testing.T) {
	_, err := NewAOF(AOFConfig{
		Path: "../../test.aof",
	})

	if err != ErrDirTrav {
		t.Fatalf("error = %v, want %v", err, ErrDirTrav)
	}
}

func TestAOF_Append(t *testing.T) {
	path := "test_append.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncNo,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}
	defer aof.Close()

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	})

	if err := aof.Append(command); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	expected := "*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nDavid\r\n"

	if string(data) != expected {
		t.Fatalf("file contents = %q, want %q", string(data), expected)
	}
}

func TestAOF_Append_MultipleCommands(t *testing.T) {
	path := "test_multiple.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncNo,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}
	defer aof.Close()

	commands := []protocol.Value{
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("name"),
			protocol.NewBulkString("David"),
		}),
		protocol.NewArray([]protocol.Value{
			protocol.NewBulkString("DEL"),
			protocol.NewBulkString("name"),
		}),
	}

	for _, command := range commands {
		if err := aof.Append(command); err != nil {
			t.Fatalf("Append() error = %v", err)
		}
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	expected := "*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nDavid\r\n" +
		"*2\r\n$3\r\nDEL\r\n$4\r\nname\r\n"

	if string(data) != expected {
		t.Fatalf("file contents = %q, want %q", string(data), expected)
	}
}

func TestAOF_Append_AfterClose(t *testing.T) {
	path := "test_closed.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path: path,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("PING"),
	})

	err = aof.Append(command)
	if err != ErrAOFClosed {
		t.Fatalf("Append() error = %v, want %v", err, ErrAOFClosed)
	}
}

func TestAOF_DoubleClose(t *testing.T) {
	path := "test_double_close.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path: path,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("first Close() error = %v", err)
	}

	if err := aof.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
}

func TestAOF_ConcurrentAppend(t *testing.T) {
	path := "test_concurrent.aof"

	filePath, err := validatePath(path)
	if err != nil {
		t.Fatalf("validatePath() error = %v", err)
	}

	t.Cleanup(func() {
		os.Remove(filePath)
	})

	aof, err := NewAOF(AOFConfig{
		Path:        path,
		FsyncPolicy: FsyncNo,
	})
	if err != nil {
		t.Fatalf("NewAOF() error = %v", err)
	}
	defer aof.Close()

	command := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name"),
		protocol.NewBulkString("David"),
	})

	const goroutines = 10
	const appendsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < appendsPerGoroutine; j++ {
				if err := aof.Append(command); err != nil {
					t.Errorf("Append() error = %v", err)
				}
			}
		}()
	}

	wg.Wait()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	expectedCommand := "*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nDavid\r\n"

	expectedLength := len(expectedCommand) * goroutines * appendsPerGoroutine

	if len(data) != expectedLength {
		t.Fatalf("file length = %d, want %d", len(data), expectedLength)
	}
}
