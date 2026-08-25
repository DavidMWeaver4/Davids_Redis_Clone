package store

import "testing"

func TestList_LPush(t *testing.T) {
	l := List{
		values: []string{"C", "D"},
	}

	got := l.LPush("A", "B")

	if got != 4 {
		t.Fatalf("expected length 4, got %d", got)
	}

	expected := []string{"B", "A", "C", "D"}

	for i := range expected {
		if l.values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, l.values[i])
		}
	}
}

func TestList_RPush(t *testing.T) {
	l := List{
		values: []string{"A", "B"},
	}

	got := l.RPush("C", "D")

	if got != 4 {
		t.Fatalf("expected length 4, got %d", got)
	}

	expected := []string{"A", "B", "C", "D"}

	for i := range expected {
		if l.values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, l.values[i])
		}
	}
}

func TestList_LPop(t *testing.T) {
	l := List{
		values: []string{"A", "B", "C"},
	}

	value, ok := l.LPop()

	if !ok {
		t.Fatal("expected value")
	}

	if value != "A" {
		t.Fatalf("expected A, got %q", value)
	}

	expected := []string{"B", "C"}

	for i := range expected {
		if l.values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, l.values[i])
		}
	}
}

func TestList_RPop(t *testing.T) {
	l := List{
		values: []string{"A", "B", "C"},
	}

	value, ok := l.RPop()

	if !ok {
		t.Fatal("expected value")
	}

	if value != "C" {
		t.Fatalf("expected C, got %q", value)
	}

	expected := []string{"A", "B"}

	for i := range expected {
		if l.values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, l.values[i])
		}
	}
}

func TestList_LPop_Empty(t *testing.T) {
	l := List{}

	value, ok := l.LPop()

	if ok {
		t.Fatal("expected no value")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
}

func TestList_RPop_Empty(t *testing.T) {
	l := List{}

	value, ok := l.RPop()

	if ok {
		t.Fatal("expected no value")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
}

func TestList_LLen(t *testing.T) {
	l := List{
		values: []string{"A", "B", "C"},
	}

	if got := l.LLen(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestList_LIndex(t *testing.T) {
	l := List{
		values: []string{"A", "B", "C"},
	}

	value, ok := l.LIndex(1)

	if !ok {
		t.Fatal("expected value")
	}

	if value != "B" {
		t.Fatalf("expected B, got %q", value)
	}
}

func TestList_LIndex_OutOfRange(t *testing.T) {
	l := List{
		values: []string{"A", "B"},
	}

	_, ok := l.LIndex(10)

	if ok {
		t.Fatal("expected index to be out of range")
	}
}

func TestList_LSet(t *testing.T) {
	l := List{
		values: []string{"A", "B", "C"},
	}

	err := l.LSet(1, "X")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if l.values[1] != "X" {
		t.Fatalf("expected X, got %q", l.values[1])
	}
}

func TestList_LSet_OutOfRange(t *testing.T) {
	l := List{
		values: []string{"A", "B"},
	}

	err := l.LSet(10, "X")

	if err != ErrIndexOutOfRange {
		t.Fatalf("expected ErrIndexOutOfRange, got %v", err)
	}
}

func TestList_LTrim(t *testing.T) {
	l := List{
		values: []string{"A", "B", "C", "D"},
	}

	l.LTrim(1, 2)

	expected := []string{"B", "C"}

	if len(l.values) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(l.values))
	}

	for i := range expected {
		if l.values[i] != expected[i] {
			t.Fatalf("expected %q at index %d, got %q", expected[i], i, l.values[i])
		}
	}
}
