package persistence

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func TestAOF_ReplayReadsCommands(t *testing.T) {
	fileName := createTempAOFFile(t)

	aof, err := NewAOF(AOFConfig{
		Path:        fileName,
		FsyncPolicy: FsyncAlways,
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

	var replayed []protocol.Value

	err = Replay(fileName, func(command protocol.Value) error {
		replayed = append(replayed, command)
		return nil
	})
	if err != nil {
		t.Fatalf("Replay() error: %v", err)
	}
	if len(replayed) != 2 {
		t.Fatalf("expected 2 replayed commands, got %d", len(replayed))
	}
	if !reflect.DeepEqual(replayed[0], setName) {
		t.Errorf("first replayed command = %#v, want %#v", replayed[0], setName)
	}
	if !reflect.DeepEqual(replayed[1], setLanguage) {
		t.Errorf("second replayed command = %#v, want %#v", replayed[1], setLanguage)
	}
}

func TestAOF_ReplayEmptyFile(t *testing.T) {
	fileName := createTempAOFFile(t)

	replayed := false

	err := Replay(fileName, func(command protocol.Value) error {
		replayed = true
		return nil
	})
	if err != nil {
		t.Fatalf("Replay() error: %v", err)
	}

	if replayed {
		t.Fatal("Replay() invoked callback for an empty AOF")
	}
}

func TestAOF_ReopenAndAppend(t *testing.T) {
	fileName := createTempAOFFile(t)

	aof, err := NewAOF(AOFConfig{
		Path:        fileName,
		FsyncPolicy: FsyncAlways,
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
		t.Fatalf("append SET name command: %v", err)
	}

	setLanguage := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("language"),
		protocol.NewBulkString("Go"),
	})
	err = aof.Append(setLanguage)
	if err != nil {
		t.Fatalf("append SET language command: %v", err)
	}

	err = aof.Close()
	if err != nil {
		t.Fatalf("close AOF: %v", err)
	}
	aof, err = NewAOF(AOFConfig{
		Path:        fileName,
		FsyncPolicy: FsyncAlways,
	})
	if err != nil {
		t.Fatalf("reopen AOF: %v", err)
	}

	setName2 := protocol.NewArray([]protocol.Value{
		protocol.NewBulkString("SET"),
		protocol.NewBulkString("name2"),
		protocol.NewBulkString("Linh"),
	})
	err = aof.Append(setName2)
	if err != nil {
		aof.Close()
		t.Fatalf("append SET name2 command: %v", err)
	}

	err = aof.Close()
	if err != nil {
		t.Fatalf("close reopened AOF: %v", err)
	}

	var replayed []protocol.Value
	err = Replay(fileName, func(command protocol.Value) error {
		replayed = append(replayed, command)
		return nil
	})
	if err != nil {
		t.Fatalf("Replay() error: %v", err)
	}

	if len(replayed) != 3 {
		t.Fatalf("expected 3 replayed commands, got %d", len(replayed))
	}
	if !reflect.DeepEqual(replayed[0], setName) {
		t.Errorf("first replayed command = %#v, want %#v", replayed[0], setName)
	}
	if !reflect.DeepEqual(replayed[1], setLanguage) {
		t.Errorf("second replayed command = %#v, want %#v", replayed[1], setLanguage)
	}
	if !reflect.DeepEqual(replayed[2], setName2) {
		t.Errorf("third replayed command = %#v, want %#v", replayed[2], setName2)
	}
}

func TestAOF_ReplayTruncatedFinalRecord(t *testing.T) {
	fileName := createTempAOFFile(t)

	aof, err := NewAOF(AOFConfig{
		Path:        fileName,
		FsyncPolicy: FsyncAlways,
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

	err = aof.Close()
	if err != nil {
		t.Fatalf("close AOF: %v", err)
	}
	file, err := os.OpenFile(
		filepath.Join("data", fileName),
		os.O_WRONLY|os.O_APPEND,
		0o644,
	)
	if err != nil {
		t.Fatalf("open AOF file: %v", err)
	}

	_, err = file.WriteString("*2\r\n$3\r\nSET\r\n$4\r\n")
	if err != nil {
		file.Close()
		t.Fatalf("append truncated record: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("close AOF file: %v", err)
	}

	var replayed []protocol.Value

	err = Replay(fileName, func(command protocol.Value) error {
		replayed = append(replayed, command)
		return nil
	})

	if err == nil {
		t.Fatal("expected Replay() to return an error for truncated final record")
	}
	if len(replayed) != 1 {
		t.Fatalf("expected 1 complete command to be replayed, got %d", len(replayed))
	}
	if !reflect.DeepEqual(replayed[0], setName) {
		t.Errorf("first replayed command = %#v, want %#v", replayed[0], setName)
	}
}

func createTempAOFFile(t *testing.T) string {
	t.Helper()

	if err := os.MkdirAll("data", 0o755); err != nil {
		t.Fatalf("create data directory: %v", err)
	}

	file, err := os.CreateTemp("data", "test-*.aof")
	if err != nil {
		t.Fatalf("create temp AOF file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close temp AOF file: %v", err)
	}

	fileName := filepath.Base(file.Name())

	t.Cleanup(func() {
		if err := os.Remove(filepath.Join("data", fileName)); err != nil && !os.IsNotExist(err) {
			t.Errorf("remove temp AOF file: %v", err)
		}
	})

	return fileName
}
