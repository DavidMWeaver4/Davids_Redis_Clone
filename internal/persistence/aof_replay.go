package persistence

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
)

func Replay(path string, apply func(protocol.Value) error) error {
	path, err := validatePath(path)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		command, err := protocol.Read(reader)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		err = apply(command)
		if err != nil {
			return err
		}

	}
	return nil
}
