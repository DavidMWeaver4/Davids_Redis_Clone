package store

import (
	"errors"
	"math"
	"strconv"
	"testing"
	"time"
)

func TestStore_Incr_Success(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 0)

	got, err := c.Incr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 21 {
		t.Fatalf("expected 21, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "21" {
		t.Fatalf("expected stored value %q, got %q", "21", value)
	}
}
func TestStore_Incr_PreservesTTL(t *testing.T) {
	c := New()

	c.Set("Foo", "20", 10*time.Second)

	got, err := c.Incr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 21 {
		t.Fatalf("expected 21, got %d", got)
	}

	ttl := c.TTL("Foo")
	if ttl < 8 || ttl > 10 {
		t.Fatalf("expected TTL around 10, got %d", ttl)
	}
}
func TestStore_Incr_MissingKey(t *testing.T) {
	c := New()

	got, err := c.Incr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "1" {
		t.Fatalf("expected value %q, got %q", "1", value)
	}
}
func TestStore_Incr_NegativeValue(t *testing.T) {
	c := New()
	c.Set("Foo", "-20", 0)

	got, err := c.Incr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -19 {
		t.Fatalf("expected -19, got %d", got)
	}
}
func TestStore_Incr_NonInteger(t *testing.T) {
	c := New()
	c.Set("Foo", "hello", 0)

	_, err := c.Incr("Foo")
	if err == nil {
		t.Fatal("expected error for non-integer value")
	}
	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}
func TestStore_Incr_Overflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MaxInt64, 10), 0)

	_, err := c.Incr("Foo")
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}
func TestStore_Incr_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 5*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	got, err := c.Incr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
func TestStore_Incr_Zero(t *testing.T) {
	c := New()
	c.Set("Foo", "0", 0)

	got, err := c.Incr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}
func TestStore_Decr_Success(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 0)

	got, err := c.Decr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 19 {
		t.Fatalf("expected 19, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "19" {
		t.Fatalf("expected stored value %q, got %q", "19", value)
	}
}
func TestStore_Decr_PreservesTTL(t *testing.T) {
	c := New()

	c.Set("Foo", "20", 10*time.Second)

	got, err := c.Decr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 19 {
		t.Fatalf("expected 19, got %d", got)
	}

	ttl := c.TTL("Foo")
	if ttl < 8 || ttl > 10 {
		t.Fatalf("expected TTL around 10, got %d", ttl)
	}
}
func TestStore_Decr_MissingKey(t *testing.T) {
	c := New()

	got, err := c.Decr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -1 {
		t.Fatalf("expected -1, got %d", got)
	}
}
func TestStore_Decr_NegativeValue(t *testing.T) {
	c := New()
	c.Set("Foo", "-20", 0)

	got, err := c.Decr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -21 {
		t.Fatalf("expected -21, got %d", got)
	}
}
func TestStore_Decr_NonInteger(t *testing.T) {
	c := New()
	c.Set("Foo", "hello", 0)

	_, err := c.Decr("Foo")
	if err == nil {
		t.Fatal("expected error for non-integer value")
	}
	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}
func TestStore_Decr_Overflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	_, err := c.Decr("Foo")
	if err == nil {
		t.Fatal("expected overflow error")
	}
	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}
func TestStore_Decr_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 5*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	got, err := c.Decr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -1 {
		t.Fatalf("expected -1, got %d", got)
	}
}
func TestStore_Decr_Zero(t *testing.T) {
	c := New()
	c.Set("Foo", "0", 0)

	got, err := c.Decr("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -1 {
		t.Fatalf("expected -1, got %d", got)
	}
}

/*
 *
 * INCRBY and DECRBY
 *
 *
 */
func TestStore_Incrby_MissingKey(t *testing.T) {
	c := New()

	got, err := c.Incrby("Foo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "5" {
		t.Fatalf("expected value %q, got %q", "5", value)
	}
}

func TestStore_Incrby_ExistingInteger(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 0)

	got, err := c.Incrby("Foo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 25 {
		t.Fatalf("expected 25, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "25" {
		t.Fatalf("expected value %q, got %q", "25", value)
	}
}

func TestStore_Incrby_NonInteger(t *testing.T) {
	c := New()
	c.Set("Foo", "hello", 0)

	_, err := c.Incrby("Foo", 5)
	if err == nil {
		t.Fatal("expected error for non-integer value")
	}

	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}

func TestStore_Incrby_PositiveOverflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MaxInt64, 10), 0)

	_, err := c.Incrby("Foo", 1)
	if err == nil {
		t.Fatal("expected overflow error")
	}

	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}

func TestStore_Incrby_NegativeOverflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	_, err := c.Incrby("Foo", -1)
	if err == nil {
		t.Fatal("expected overflow error")
	}

	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}

func TestStore_Incrby_Zero(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 0)

	got, err := c.Incrby("Foo", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 20 {
		t.Fatalf("expected 20, got %d", got)
	}
}

func TestStore_Incrby_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 5*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	got, err := c.Incrby("Foo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "5" {
		t.Fatalf("expected value %q, got %q", "5", value)
	}
}
func TestStore_Incrby_MaxInt64(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MaxInt64-1, 10), 0)

	got, err := c.Incrby("Foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != math.MaxInt64 {
		t.Fatalf("expected %d, got %d", int64(math.MaxInt64), got)
	}
}
func TestStore_Decrby_MissingKey(t *testing.T) {
	c := New()

	got, err := c.Decrby("Foo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -5 {
		t.Fatalf("expected -5, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "-5" {
		t.Fatalf("expected value %q, got %q", "-5", value)
	}
}

func TestStore_Decrby_ExistingInteger(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 0)

	got, err := c.Decrby("Foo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 15 {
		t.Fatalf("expected 15, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "15" {
		t.Fatalf("expected value %q, got %q", "15", value)
	}
}

func TestStore_Decrby_NonInteger(t *testing.T) {
	c := New()
	c.Set("Foo", "hello", 0)

	_, err := c.Decrby("Foo", 5)
	if err == nil {
		t.Fatal("expected error for non-integer value")
	}

	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}

func TestStore_Decrby_NegativeArgument(t *testing.T) {
	c := New()

	_, err := c.Decrby("Foo", -5)
	if err == nil {
		t.Fatal("expected error for negative decrement")
	}

	if !errors.Is(err, ErrNegativeDecrement) {
		t.Fatalf("expected ErrNegativeDecrement, got %v", err)
	}
}

func TestStore_Decrby_MinInt64Underflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	_, err := c.Decrby("Foo", 1)
	if err == nil {
		t.Fatal("expected underflow error")
	}

	if !errors.Is(err, ErrNotIntegerOrOutOfRange) {
		t.Fatalf("expected ErrNotIntegerOrOutOfRange, got %v", err)
	}
}

func TestStore_Decrby_Zero(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 0)

	got, err := c.Decrby("Foo", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 20 {
		t.Fatalf("expected 20, got %d", got)
	}
}

func TestStore_Decrby_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "20", 5*time.Millisecond)

	time.Sleep(10 * time.Millisecond)

	got, err := c.Decrby("Foo", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != -5 {
		t.Fatalf("expected -5, got %d", got)
	}

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "-5" {
		t.Fatalf("expected value %q, got %q", "-5", value)
	}
}
func TestStore_Incrby_MinInt64(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MinInt64+1, 10), 0)

	got, err := c.Incrby("Foo", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != math.MinInt64 {
		t.Fatalf("expected %d, got %d", int64(math.MinInt64), got)
	}
}
func TestStore_Decrby_MinInt64(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MinInt64+1, 10), 0)

	got, err := c.Decrby("Foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != math.MinInt64 {
		t.Fatalf("expected %d, got %d", int64(math.MinInt64), got)
	}
}
func TestStore_Incr_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type: ListType,
		List: List{
			values: []string{"Bar"},
		},
	}

	_, err := s.Incr("Foo")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
