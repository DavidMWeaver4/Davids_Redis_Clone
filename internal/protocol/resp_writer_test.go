package protocol

import (
	"bytes"
	"testing"
)

func TestWrite_SimpleString_Success(t *testing.T) {
	var buf bytes.Buffer

	err := Write(&buf, NewSimpleString("PONG"))
	if err != nil {
		t.Fatal(err)
	}

	expected := "+PONG\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}

func TestWrite_Error_Success(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, NewError("invalid command"))
	if err != nil {
		t.Fatal(err)
	}
	expected := "-ERR invalid command\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}

func TestWrite_Integer_Success(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, NewInteger(42))
	if err != nil {
		t.Fatal(err)
	}
	expected := ":42\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}
func TestWrite_BulkString_Success(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, NewBulkString("hello"))
	if err != nil {
		t.Fatal(err)
	}
	expected := "$5\r\nhello\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}
func TestWrite_NullBulkString_Success(t *testing.T) {
	var buf bytes.Buffer
	err := Write(&buf, NewNullBulkString())
	if err != nil {
		t.Fatal(err)
	}
	expected := "$-1\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}
func TestWrite_Array_Success(t *testing.T) {
	var buf bytes.Buffer
	data := NewArray([]Value{
		NewBulkString("Foo"),
		NewBulkString("Bar"),
	})
	err := Write(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	expected := "*2\r\n$3\r\nFoo\r\n$3\r\nBar\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}
func TestWrite_NestedArray_Success(t *testing.T) {
	var buf bytes.Buffer
	data := NewArray([]Value{
		NewBulkString("outer"),
		NewArray([]Value{
			NewBulkString("inner"),
			NewInteger(123),
		}),
	})
	err := Write(&buf, data)
	if err != nil {
		t.Fatal(err)
	}
	expected := "*2\r\n$5\r\nouter\r\n*2\r\n$5\r\ninner\r\n:123\r\n"
	if buf.String() != expected {
		t.Fatalf("expected %q, got %q", expected, buf.String())
	}
}
func TestWrite_UnknownType_ReturnsError(t *testing.T) {
	var buf bytes.Buffer

	err := Write(&buf, Value{Type: "mystery"})

	if err == nil {
		t.Fatal("expected an error for an unknown RESP type")
	}
}
