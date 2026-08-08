package protocol

import (
	"bufio"
	"strings"
	"testing"
)

func newReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

func TestRead_SimpleString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "OK",
			input: "+OK\r\n",
			want:  "OK",
		},
		{
			name:  "Long",
			input: "+supercalifragilisticexpialidocious\r\n",
			want:  "supercalifragilisticexpialidocious",
		},
		{
			name:  "contains spaces",
			input: "+test value\r\n",
			want:  "test value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(newReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Type != SimpleString {
				t.Fatalf("expected SimpleString, got %v", got.Type)
			}

			if got.Str != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got.Str)
			}
		})
	}
}

func TestRead_Error(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "invalid",
			input: "-ERR invalid command\r\n",
			want:  "ERR invalid command",
		},
		{
			name:  "something went wrong",
			input: "-ERR something went wrong\r\n",
			want:  "ERR something went wrong",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(newReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Type != Error {
				t.Fatalf("expected Error, got %v", got.Type)
			}

			if got.Str != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got.Str)
			}
		})
	}
}

func TestRead_Integer_Success(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{
			name:  "positive num",
			input: ":23\r\n",
			want:  23,
		},
		{
			name:  "negative",
			input: ":-23\r\n",
			want:  -23,
		},
		{
			name:  "zero",
			input: ":0\r\n",
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(newReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Type != Integer {
				t.Fatalf("expected Integer, got %v", got.Type)
			}

			if got.Int != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got.Int)
			}
		})
	}
}

func TestRead_Integer_Invalid(t *testing.T) {
	tests := []string{
		":dmw\r\n",
		":123abc\r\n",
		":\r\n",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Read(newReader(input))
			if err == nil {
				t.Fatalf("expected error for input %q", input)
			}
		})
	}
}

func TestRead_BulkString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "standard",
			input: "$5\r\nhello\r\n",
			want:  "hello",
		},
		{
			name:  "empty",
			input: "$0\r\n\r\n",
			want:  "",
		},
		{
			name:  "space",
			input: "$11\r\nhello world\r\n",
			want:  "hello world",
		},
		{
			name:  "contains CRLF",
			input: "$12\r\nhello\r\nworld\r\n",
			want:  "hello\r\nworld",
		},
		{
			name:  "UTF-8",
			input: "$15\r\nこんにちは\r\n",
			want:  "こんにちは",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Read(newReader(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.Type != BulkString {
				t.Fatalf("expected BulkString, got %v", got.Type)
			}

			if got.Str != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got.Str)
			}
		})
	}
}

func TestRead_NullBulkString(t *testing.T) {
	got, err := Read(newReader("$-1\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Type != Null {
		t.Fatalf("expected Null, got %v", got.Type)
	}
}

func TestRead_BulkString_InvalidLength(t *testing.T) {
	tests := []string{
		"$dmw\r\n",
		"$-2\r\n",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Read(newReader(input))
			if err == nil {
				t.Fatalf("expected error for input %q", input)
			}
		})
	}
}

func TestRead_BulkString_IncompletePayload(t *testing.T) {
	tests := []string{
		"$5\r\nhel\r\n",
		"$5\r\nhello",
		"$5\r\nhelloXX",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Read(newReader(input))
			if err == nil {
				t.Fatalf("expected error for input %q", input)
			}
		})
	}
}

func TestRead_Array_Success(t *testing.T) {
	input := "*2\r\n$5\r\nhello\r\n:23\r\n"

	got, err := Read(newReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Type != Array {
		t.Fatalf("expected Array, got %v", got.Type)
	}

	if len(got.Array) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(got.Array))
	}

	if got.Array[0].Type != BulkString {
		t.Fatalf("expected first element to be BulkString")
	}

	if got.Array[0].Str != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got.Array[0].Str)
	}

	if got.Array[1].Type != Integer {
		t.Fatalf("expected second element to be Integer")
	}

	if got.Array[1].Int != 23 {
		t.Fatalf("expected 23, got %d", got.Array[1].Int)
	}
}

func TestRead_EmptyArray(t *testing.T) {
	got, err := Read(newReader("*0\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Type != Array {
		t.Fatalf("expected Array, got %v", got.Type)
	}

	if len(got.Array) != 0 {
		t.Fatalf("expected empty array, got %d elements", len(got.Array))
	}
}

func TestRead_NestedArray_Success(t *testing.T) {
	input := "*2\r\n*1\r\n:23\r\n+OK\r\n"

	got, err := Read(newReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Type != Array {
		t.Fatalf("expected Array, got %v", got.Type)
	}

	if len(got.Array) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(got.Array))
	}

	if got.Array[0].Type != Array {
		t.Fatalf("expected first element to be Array")
	}

	if len(got.Array[0].Array) != 1 {
		t.Fatalf(
			"expected nested array to contain 1 element, got %d",
			len(got.Array[0].Array),
		)
	}

	if got.Array[0].Array[0].Type != Integer {
		t.Fatalf("expected nested element to be Integer")
	}

	if got.Array[0].Array[0].Int != 23 {
		t.Fatalf(
			"expected nested integer 23, got %d",
			got.Array[0].Array[0].Int,
		)
	}

	if got.Array[1].Type != SimpleString {
		t.Fatalf("expected second element to be SimpleString")
	}

	if got.Array[1].Str != "OK" {
		t.Fatalf("expected OK, got %q", got.Array[1].Str)
	}
}

func TestRead_Array_InvalidLength(t *testing.T) {
	tests := []string{
		"*dmw\r\n",
		"*-2\r\n",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Read(newReader(input))
			if err == nil {
				t.Fatalf("expected error for input %q", input)
			}
		})
	}
}

func TestRead_NullArray(t *testing.T) {
	got, err := Read(newReader("*-1\r\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Type != NullArray {
		t.Fatalf("expected NullArray, got %v", got.Type)
	}
}

func TestRead_InvalidPrefix(t *testing.T) {
	tests := []string{
		"@\r\n",
		"hello\r\n",
		"\r\n",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Read(newReader(input))
			if err == nil {
				t.Fatalf("expected error for input %q", input)
			}
		})
	}
}

func TestRead_LineRequiresCRLF(t *testing.T) {
	tests := []string{
		"+OK\n",
		"+OK\r",
		"+OK",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := Read(newReader(input))
			if err == nil {
				t.Fatalf("expected error for invalid line %q", input)
			}
		})
	}
}
