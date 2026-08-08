package protocol

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Read(reader *bufio.Reader) (Value, error) {
	prefix, err := reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch prefix {
	case '+':
		return readSimpleString(reader)
	case '-':
		return readError(reader)
	case ':':
		return readInteger(reader)
	case '$':
		return readBulkString(reader)
	case '*':
		return readArray(reader)
	default:
		return Value{}, fmt.Errorf("invalid RESP prefix %q", prefix)
	}

}

func readLine(reader *bufio.Reader) (string, error) {
	//strip the \r\n and send the next line
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read RESP line: %w", err)
	}
	if len(line) < 2 {
		return "", fmt.Errorf("wrong byte length in readLine got %v, %q", len(line), line)
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", fmt.Errorf("invalid RESP input, got %v", line)
	}
	return strings.TrimSuffix(line, "\r\n"), nil
}

func readSimpleString(reader *bufio.Reader) (Value, error) {
	str, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}
	return NewSimpleString(str), nil
}

func readError(reader *bufio.Reader) (Value, error) {
	str, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}
	return NewError(str), nil
}

func readInteger(reader *bufio.Reader) (Value, error) {
	str, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}
	convertedInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return Value{}, fmt.Errorf("parsing RESP integer: %w", err)
	}
	return NewInteger(convertedInt), nil
}

func readBulkString(reader *bufio.Reader) (Value, error) {
	length, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}
	convertLength, err := strconv.Atoi(length)
	if err != nil {
		return Value{}, fmt.Errorf("parsing RESP bulkstring integer: %w", err)
	}
	if convertLength == -1 {
		return NewNullBulkString(), nil
	}
	if convertLength < -1 {
		return Value{}, fmt.Errorf("invalid bulkstring length, got %d", convertLength)
	}
	payload := make([]byte, convertLength)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return Value{}, fmt.Errorf("read bulk-string payload: %w", err)
	}

	terminator := make([]byte, 2)
	_, err = io.ReadFull(reader, terminator)
	if err != nil {
		return Value{}, fmt.Errorf("read bulk-string terminator: %w", err)
	}
	if !bytes.Equal(terminator, []byte("\r\n")) {
		return Value{}, fmt.Errorf("invalid bulk-string terminator: %q", terminator)
	}

	return NewBulkString(string(payload)), nil
}

func readArray(reader *bufio.Reader) (Value, error) {
	length, err := readLine(reader)
	if err != nil {
		return Value{}, err
	}
	convertLength, err := strconv.Atoi(length)
	if err != nil {
		return Value{}, fmt.Errorf("parsing RESP array integer: %w", err)
	}
	if convertLength == -1 {
		return NewNullArray(), nil
	}
	if convertLength < -1 {
		return Value{}, fmt.Errorf("invalid array length, got %d", convertLength)
	}
	values := make([]Value, convertLength)
	for i := range convertLength {
		next, err := Read(reader)
		if err != nil {
			return Value{}, fmt.Errorf("reading error on index %d, got %w", i, err)
		}
		values[i] = next
	}
	return NewArray(values), nil
}
