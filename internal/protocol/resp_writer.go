package protocol

import (
	"fmt"
	"io"
)

func Write(w io.Writer, v Value) error {
	switch v.Type {
	case SimpleString:
		_, err := fmt.Fprintf(w, "+%s\r\n", v.Str)
		return err
	case Error:
		_, err := fmt.Fprintf(w, "-%s\r\n", v.Str)
		return err
	case Integer:
		_, err := fmt.Fprintf(w, ":%d\r\n", v.Int)
		return err
	case BulkString:
		_, err := fmt.Fprintf(w, "$%d\r\n%s\r\n", len([]byte(v.Str)), v.Str)
		return err
	case Null:
		_, err := fmt.Fprint(w, "$-1\r\n")
		return err
	case Array:
		return writeArray(w, v.Array)
	default:
		return fmt.Errorf("ERR unknown RESP type %s", v.Type)
	}
}
func writeArray(w io.Writer, values []Value) error {
	_, err := fmt.Fprintf(w, "*%d\r\n", len(values))
	if err != nil {
		return err
	}

	for _, value := range values {
		err := Write(w, value)
		if err != nil {
			return err
		}
	}

	return nil
}
