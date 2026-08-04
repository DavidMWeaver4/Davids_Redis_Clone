package protocol

import (
	"fmt"
	"net"
)

func Write(conn net.Conn, v Value) error {
	switch v.Type {
	case SimpleString:
		_, err := fmt.Fprintf(conn, "+%s\r\n", v.Str)
		return err
	case Error:
		_, err := fmt.Fprintf(conn, "-%s\r\n", v.Str)
		return err
	case Integer:
		_, err := fmt.Fprintf(conn, ":%d\r\n", v.Int)
		return err
	case BulkString:
		_, err := fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(v.Str), v.Str)
		return err
	case Null:
		_, err := fmt.Fprint(conn, "$-1\r\n")
		return err
	case Array:
		return writeArray(conn, v.Array)
	default:
		return fmt.Errorf("ERR unknown RESP type %s", v.Type)
	}
}
func writeArray(conn net.Conn, values []Value) error {
	_, err := fmt.Fprintf(conn, "*%d\r\n", len(values))
	if err != nil {
		return err
	}

	for _, value := range values {
		err := Write(conn, value)
		if err != nil {
			return err
		}
	}

	return nil
}
