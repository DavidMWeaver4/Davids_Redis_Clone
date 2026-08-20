package server

import (
	"math"
	"strconv"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/protocol"
	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store"
)

func TestCommands_Incr_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 21 {
		t.Fatalf("expected 21, got %d", response.Int)
	}
}

func TestCommands_Incr_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 1 {
		t.Fatalf("expected 1, got %d", response.Int)
	}
}

func TestCommands_Incr_NonInteger(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Incr_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := incr(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Incr_Overflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MaxInt64, 10), 0)

	response := incr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}
func TestCommands_Decr_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 19 {
		t.Fatalf("expected 19, got %d", response.Int)
	}
}

func TestCommands_Decr_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != -1 {
		t.Fatalf("expected -1, got %d", response.Int)
	}
}

func TestCommands_Decr_NonInteger(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Decr_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo", "Bar"},
	}

	for _, args := range tests {
		response := decr(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Decr_Underflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	response := decr(s, []string{"Foo"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}
func TestCommands_Incrby_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := incrby(s, []string{"Foo", "5"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 25 {
		t.Fatalf("expected 25, got %d", response.Int)
	}
}

func TestCommands_Incrby_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := incrby(s, []string{"Foo", "5"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 5 {
		t.Fatalf("expected 5, got %d", response.Int)
	}
}

func TestCommands_Incrby_NonInteger(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := incrby(s, []string{"Foo", "5"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Incrby_InvalidIncrement(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := incrby(s, []string{"Foo", "abc"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != "cannot INCRBY non-integer" {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Incrby_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "5", "Bar"},
	}

	for _, args := range tests {
		response := incrby(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Incrby_PositiveOverflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MaxInt64, 10), 0)

	response := incrby(s, []string{"Foo", "1"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Incrby_NegativeOverflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	response := incrby(s, []string{"Foo", "-1"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Incrby_Zero(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := incrby(s, []string{"Foo", "0"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 20 {
		t.Fatalf("expected 20, got %d", response.Int)
	}
}

func TestCommands_Decrby_Success(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := decrby(s, []string{"Foo", "5"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 15 {
		t.Fatalf("expected 15, got %d", response.Int)
	}
}

func TestCommands_Decrby_MissingKey(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := decrby(s, []string{"Foo", "5"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != -5 {
		t.Fatalf("expected -5, got %d", response.Int)
	}
}

func TestCommands_Decrby_NonInteger(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "Bar", 0)

	response := decrby(s, []string{"Foo", "5"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Decrby_InvalidDecrement(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := decrby(s, []string{"Foo", "abc"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != "cannot DECRBY non-integer" {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Decrby_NegativeArgument(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	response := decrby(s, []string{"Foo", "-5"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNegativeDecrement.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Decrby_InvalidArgumentCount(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	tests := [][]string{
		{},
		{"Foo"},
		{"Foo", "5", "Bar"},
	}

	for _, args := range tests {
		response := decrby(s, args)

		if response.Type != protocol.Error {
			t.Fatalf("expected Error, got %v", response.Type)
		}
	}
}

func TestCommands_Decrby_Underflow(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	response := decrby(s, []string{"Foo", "1"})

	if response.Type != protocol.Error {
		t.Fatalf("expected Error, got %v", response.Type)
	}

	if response.Str != store.ErrNotIntegerOrOutOfRange.Error() {
		t.Fatalf("unexpected error: %q", response.Str)
	}
}

func TestCommands_Decrby_Zero(t *testing.T) {
	s := &Server{
		store: store.New(),
	}

	s.store.Set("Foo", "20", 0)

	response := decrby(s, []string{"Foo", "0"})

	if response.Type != protocol.Integer {
		t.Fatalf("expected Integer, got %v", response.Type)
	}

	if response.Int != 20 {
		t.Fatalf("expected 20, got %d", response.Int)
	}
}
