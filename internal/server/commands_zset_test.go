package server

import (
	"math"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_ZAdd(t *testing.T) {
	s := &Server{store: store.New()}

	response := zadd(s, []string{"scores", "10", "Alice"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_ZAdd_UpdateExistingMember(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	response := zadd(s, []string{"scores", "20", "Alice"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_ZAdd_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "too few arguments",
			args: []string{"scores", "10"},
		},
		{
			name: "too many arguments",
			args: []string{"scores", "10", "Alice", "extra"},
		},
		{
			name: "invalid score",
			args: []string{"scores", "abc", "Alice"},
		},
		{
			name: "nan score",
			args: []string{"scores", "NaN", "Alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := zadd(s, tt.args)

			if response.Type != protocol.Error {
				t.Fatalf("expected Error, got %v", response.Type)
			}
		})
	}
}

func TestCommands_ZScore(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10.5", "Alice"})

	response := zscore(s, []string{"scores", "Alice"})

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}
	if response.Str != "10.5" {
		t.Fatalf("expected 10.5, got %q", response.Str)
	}
}

func TestCommands_ZScore_MissingMember(t *testing.T) {
	s := &Server{store: store.New()}

	response := zscore(s, []string{"scores", "Alice"})

	if response.Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Type)
	}
}

func TestCommands_ZScore_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"scores"},
		{"scores", "Alice", "extra"},
	}

	for _, args := range tests {
		response := zscore(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_ZCard(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})

	response := zcard(s, []string{"scores"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_ZCard_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := zcard(s, []string{"scores"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_ZRem(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})

	response := zrem(s, []string{"scores", "Alice"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_ZRem_MultipleMembers(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})
	zadd(s, []string{"scores", "30", "Charlie"})

	response := zrem(s, []string{"scores", "Alice", "Charlie"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}
	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_ZRem_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	response := zrem(s, []string{"scores"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_ZRange(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "30", "Charlie"})
	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})

	response := zrange(s, []string{"scores", "0", "2"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	expected := []string{"Alice", "Bob", "Charlie"}

	if len(response.Array) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(response.Array))
	}

	for i, value := range response.Array {
		if value.Type != protocol.BulkString {
			t.Fatalf("response[%d]: expected BulkString, got %v", i, value.Type)
		}
		if value.Str != expected[i] {
			t.Fatalf("response[%d]: expected %q, got %q", i, expected[i], value.Str)
		}
	}
}

func TestCommands_ZRange_Empty(t *testing.T) {
	s := &Server{store: store.New()}

	response := zrange(s, []string{"scores", "0", "10"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if response.Array == nil {
		t.Fatal("expected non-nil empty array")
	}

	if len(response.Array) != 0 {
		t.Fatalf("expected empty array, got %d values", len(response.Array))
	}
}

func TestCommands_ZRange_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "too few arguments",
			args: []string{"scores", "0"},
		},
		{
			name: "too many arguments",
			args: []string{"scores", "0", "1", "extra"},
		},
		{
			name: "invalid start",
			args: []string{"scores", "abc", "1"},
		},
		{
			name: "invalid end",
			args: []string{"scores", "0", "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := zrange(s, tt.args)

			if response.Type != protocol.Error {
				t.Fatalf("expected Error, got %v", response.Type)
			}
		})
	}
}

func TestCommands_ZRange_NegativeIndexes(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})
	zadd(s, []string{"scores", "30", "Charlie"})

	response := zrange(s, []string{"scores", "-2", "-1"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	expected := []string{"Bob", "Charlie"}

	if len(response.Array) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(response.Array))
	}

	for i, value := range response.Array {
		if value.Str != expected[i] {
			t.Fatalf("response[%d]: expected %q, got %q", i, expected[i], value.Str)
		}
	}
}

func TestParseScore(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float64
		wantErr bool
	}{
		{
			name:  "integer",
			input: "10",
			want:  10,
		},
		{
			name:  "decimal",
			input: "10.5",
			want:  10.5,
		},
		{
			name:  "negative",
			input: "-5.25",
			want:  -5.25,
		},
		{
			name:    "invalid",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "nan",
			input:   "NaN",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseScore(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tt.want) > 0.000001 {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestScoreValue(t *testing.T) {
	response := scoreValue(10.5)

	if response.Type != protocol.BulkString {
		t.Fatalf("expected BulkString, got %v", response.Type)
	}

	if response.Str != "10.5" {
		t.Fatalf("expected 10.5, got %q", response.Str)
	}
}
func TestCommands_ZRange_WithScores(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "30", "Charlie"})
	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})

	response := zrange(s, []string{"scores", "0", "-1", "WITHSCORES"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	expected := []string{
		"Alice", "10",
		"Bob", "20",
		"Charlie", "30",
	}

	if len(response.Array) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(response.Array))
	}

	for i, value := range response.Array {
		if value.Type != protocol.BulkString {
			t.Fatalf("response[%d]: expected BulkString, got %v", i, value.Type)
		}
		if value.Str != expected[i] {
			t.Fatalf("response[%d]: expected %q, got %q", i, expected[i], value.Str)
		}
	}
}
func TestCommands_ZRange_InvalidOption(t *testing.T) {
	s := &Server{store: store.New()}

	response := zrange(s, []string{"scores", "0", "-1", "INVALID"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}
func TestCommands_ZRange_WithScores_CaseInsensitive(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})

	response := zrange(s, []string{"scores", "0", "-1", "withscores"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if len(response.Array) != 2 {
		t.Fatalf("expected 2 results, got %d", len(response.Array))
	}

	if response.Array[0].Str != "Alice" {
		t.Fatalf("expected Alice, got %q", response.Array[0].Str)
	}

	if response.Array[1].Str != "10" {
		t.Fatalf("expected 10, got %q", response.Array[1].Str)
	}
}
func TestCommands_ZRange_PastEnd(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})

	response := zrange(s, []string{"scores", "0", "10"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	expected := []string{"Alice", "Bob"}

	if len(response.Array) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(response.Array))
	}

	for i, value := range response.Array {
		if value.Type != protocol.BulkString {
			t.Fatalf("response[%d]: expected BulkString, got %v", i, value.Type)
		}

		if value.Str != expected[i] {
			t.Fatalf("response[%d]: expected %q, got %q", i, expected[i], value.Str)
		}
	}
}

func TestCommands_ZRange_InvalidRange(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})
	zadd(s, []string{"scores", "30", "Charlie"})

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "start greater than end",
			args: []string{"scores", "2", "1"},
		},
		{
			name: "start past end of set",
			args: []string{"scores", "10", "20"},
		},
		{
			name: "end before beginning",
			args: []string{"scores", "0", "-10"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := zrange(s, tt.args)

			if response.Type != protocol.Array {
				t.Fatalf("expected Array, got %v", response.Type)
			}

			if response.Array == nil {
				t.Fatal("expected non-nil empty array")
			}

			if len(response.Array) != 0 {
				t.Fatalf("expected empty array, got %d values", len(response.Array))
			}
		})
	}
}

func TestCommands_ZRange_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := zrange(s, []string{"scores", "0", "-1"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	if response.Array == nil {
		t.Fatal("expected non-nil empty array")
	}

	if len(response.Array) != 0 {
		t.Fatalf("expected empty array, got %d values", len(response.Array))
	}
}

func TestCommands_ZRange_WrongType(t *testing.T) {
	s := &Server{store: store.New()}

	s.store.Set("scores", "value", 0)

	response := zrange(s, []string{"scores", "0", "-1"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}

func TestCommands_ZRange_WithScores_SubRange(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})
	zadd(s, []string{"scores", "30", "Charlie"})

	response := zrange(s, []string{"scores", "1", "2", "WITHSCORES"})

	if response.Type != protocol.Array {
		t.Fatalf("expected Array, got %v", response.Type)
	}

	expected := []string{
		"Bob", "20",
		"Charlie", "30",
	}

	if len(response.Array) != len(expected) {
		t.Fatalf("expected %d results, got %d", len(expected), len(response.Array))
	}

	for i, value := range response.Array {
		if value.Type != protocol.BulkString {
			t.Fatalf("response[%d]: expected BulkString, got %v", i, value.Type)
		}

		if value.Str != expected[i] {
			t.Fatalf("response[%d]: expected %q, got %q", i, expected[i], value.Str)
		}
	}
}

func TestCommands_ZRank(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "30", "Charlie"})
	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})

	response := zrank(s, []string{"scores", "Bob"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_ZRank_FirstMember(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})
	zadd(s, []string{"scores", "30", "Charlie"})

	response := zrank(s, []string{"scores", "Alice"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 0 {
		t.Fatalf("expected 0, got %d", response.Int)
	}
}

func TestCommands_ZRank_LastMember(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "20", "Bob"})
	zadd(s, []string{"scores", "30", "Charlie"})

	response := zrank(s, []string{"scores", "Charlie"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 2 {
		t.Fatalf("expected 2, got %d", response.Int)
	}
}

func TestCommands_ZRank_MissingMember(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Alice"})

	response := zrank(s, []string{"scores", "Bob"})

	if response.Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Type)
	}
}

func TestCommands_ZRank_MissingKey(t *testing.T) {
	s := &Server{store: store.New()}

	response := zrank(s, []string{"scores", "Alice"})

	if response.Type != protocol.Null {
		t.Fatalf("expected Null, got %v", response.Type)
	}
}

func TestCommands_ZRank_InvalidArguments(t *testing.T) {
	s := &Server{store: store.New()}

	tests := [][]string{
		{},
		{"scores"},
		{"scores", "Alice", "extra"},
	}

	for _, args := range tests {
		response := zrank(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_ZRank_WrongType(t *testing.T) {
	s := &Server{store: store.New()}

	s.store.Set("scores", "value", 0)

	response := zrank(s, []string{"scores", "Alice"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}
}
func TestCommands_ZRank_TiedScores(t *testing.T) {
	s := &Server{store: store.New()}

	zadd(s, []string{"scores", "10", "Charlie"})
	zadd(s, []string{"scores", "10", "Alice"})
	zadd(s, []string{"scores", "10", "Bob"})

	response := zrank(s, []string{"scores", "Bob"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}
