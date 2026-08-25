package store

import (
	"testing"
	"time"
)

func TestStore_LPush_Success(t *testing.T) {
	s := New()

	got, err := s.LPush("Foo", "Bar", "Baz", "Qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	values, err := s.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"Qux", "Baz", "Bar"}

	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LPush_ExistingList(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.LPush("Foo", "Qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	values, err := s.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"Qux", "Bar", "Baz"}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LPush_ExpiredKey(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type:      ListType,
		List:      List{values: []string{"Old"}},
		ExpiresAt: time.Now().Add(-time.Second),
	}

	_, err := s.LPush("Foo", "New")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 1 || values[0] != "New" {
		t.Fatalf("expected [New], got %v", values)
	}
}

func TestStore_LPush_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type:   StringType,
		String: "Bar",
	}

	_, err := s.LPush("Foo", "Baz")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_RPush_Success(t *testing.T) {
	s := New()

	got, err := s.RPush("Foo", "Bar", "Baz", "Qux")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	values, err := s.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"Bar", "Baz", "Qux"}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_RPush_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type:   StringType,
		String: "Bar",
	}

	_, err := s.RPush("Foo", "Baz")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_LPop_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok, err := s.LPop("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected value to be found")
	}

	if value != "Bar" {
		t.Fatalf("expected %q, got %q", "Bar", value)
	}

	length, err := s.LLen("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if length != 1 {
		t.Fatalf("expected length 1, got %d", length)
	}
}

func TestStore_RPop_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok, err := s.RPop("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected value to be found")
	}

	if value != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", value)
	}
}

func TestStore_LPop_EmptyList(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type: ListType,
		List: List{values: []string{}},
	}

	value, ok, err := s.LPop("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected no value")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}

	if s.Exists("Foo") {
		t.Fatal("expected empty list key to be deleted")
	}
}

func TestStore_LPop_MissingKey(t *testing.T) {
	s := New()

	value, ok, err := s.LPop("Foo")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected no value")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
}

func TestStore_RPop_MissingKey(t *testing.T) {
	s := New()

	value, ok, err := s.RPop("Foo")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected no value")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
}

func TestStore_LPop_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	_, _, err := s.LPop("Foo")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_RPop_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	_, _, err := s.RPop("Foo")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_LLen_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "Bar", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.LLen("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 2 {
		t.Fatalf("expected length 2, got %d", got)
	}
}

func TestStore_LLen_MissingKey(t *testing.T) {
	s := New()

	got, err := s.LLen("Foo")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestStore_LLen_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	_, err := s.LLen("Foo")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_LRange_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"B", "C"}

	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LRange_NegativeIndexes(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", -3, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"B", "C", "D"}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LRange_OutOfBounds(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", -100, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"A", "B"}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LRange_InvalidRange(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", 2, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(values) != 0 {
		t.Fatalf("expected empty result, got %v", values)
	}
}

func TestStore_LRange_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	_, err := s.LRange("Foo", 0, -1)

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_LIndex_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok, err := s.LIndex("Foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected value to exist")
	}

	if value != "B" {
		t.Fatalf("expected %q, got %q", "B", value)
	}
}

func TestStore_LIndex_NegativeIndex(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok, err := s.LIndex("Foo", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected value to exist")
	}

	if value != "C" {
		t.Fatalf("expected %q, got %q", "C", value)
	}
}

func TestStore_LIndex_OutOfRange(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok, err := s.LIndex("Foo", 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected no value")
	}
}

func TestStore_LIndex_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	_, _, err := s.LIndex("Foo", 0)

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_LSet_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.LSet("Foo", 1, "X")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok, err := s.LIndex("Foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected value to exist")
	}

	if value != "X" {
		t.Fatalf("expected %q, got %q", "X", value)
	}
}

func TestStore_LSet_NegativeIndex(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.LSet("Foo", -1, "X")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, _, err := s.LIndex("Foo", -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if value != "X" {
		t.Fatalf("expected %q, got %q", "X", value)
	}
}

func TestStore_LSet_MissingKey(t *testing.T) {
	s := New()

	err := s.LSet("Foo", 0, "Bar")

	if err != ErrNoKey {
		t.Fatalf("expected ErrNoKey, got %v", err)
	}
}

func TestStore_LSet_OutOfRange(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.LSet("Foo", 10, "Bar")

	if err != ErrIndexOutOfRange {
		t.Fatalf("expected ErrIndexOutOfRange, got %v", err)
	}
}

func TestStore_LSet_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	err := s.LSet("Foo", 0, "Baz")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_LTrim_Success(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.LTrim("Foo", 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"B", "C"}

	if len(values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(values))
	}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LTrim_NegativeIndexes(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C", "D")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.LTrim("Foo", -3, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	values, err := s.LRange("Foo", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"B", "C", "D"}

	for i := range expected {
		if values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, values[i])
		}
	}
}

func TestStore_LTrim_InvalidRange_DeletesKey(t *testing.T) {
	s := New()

	_, err := s.RPush("Foo", "A", "B", "C")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = s.LTrim("Foo", 3, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Exists("Foo") {
		t.Fatal("expected key to be deleted")
	}
}

func TestStore_LTrim_MissingKey(t *testing.T) {
	s := New()

	err := s.LTrim("Foo", 0, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStore_LTrim_WrongType(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	err := s.LTrim("Foo", 0, 1)

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
