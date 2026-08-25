package store

import (
	"testing"
	"time"
)

func TestStore_Append_ExistingKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	got, err := s.Append("Foo", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 6 {
		t.Fatalf("expected length 6, got %d", got)
	}

	value, ok, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "BarBaz" {
		t.Fatalf("expected %q, got %q", "BarBaz", value)
	}
}

func TestStore_Append_MissingKey(t *testing.T) {
	s := New()

	got, err := s.Append("Foo", "Bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	value, ok, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "Bar" {
		t.Fatalf("expected %q, got %q", "Bar", value)
	}
}

func TestStore_Append_EmptyValue(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	got, err := s.Append("Foo", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	value, _, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "Bar" {
		t.Fatalf("expected %q, got %q", "Bar", value)
	}
}
func TestStore_Append_ExpiredKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	got, err := s.Append("Foo", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	value, ok, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", value)
	}
}
func TestStore_Strlen_ExistingKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	got, err := s.Strlen("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestStore_Strlen_MissingKey(t *testing.T) {
	s := New()

	got, err := s.Strlen("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestStore_Strlen_EmptyValue(t *testing.T) {
	s := New()

	s.Set("Foo", "", 0)

	got, err := s.Strlen("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
func TestStore_Setnx_MissingKey(t *testing.T) {
	s := New()

	ok, err := s.Setnx("Foo", "Bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected SETNX to succeed")
	}

	value, found, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected key to exist")
	}

	if value != "Bar" {
		t.Fatalf("expected %q, got %q", "Bar", value)
	}
}

func TestStore_Setnx_ExistingKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	ok, err := s.Setnx("Foo", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("expected SETNX to fail")
	}

	value, _, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if value != "Bar" {
		t.Fatalf("expected original value %q, got %q", "Bar", value)
	}
}

func TestStore_Setnx_ExpiredKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	ok, err := s.Setnx("Foo", "Baz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("expected SETNX to succeed for expired key")
	}

	value, found, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected key to exist")
	}

	if value != "Baz" {
		t.Fatalf("expected %q, got %q", "Baz", value)
	}
}
func TestStore_Mset(t *testing.T) {
	s := New()

	s.Mset([]KeyValue{
		{Key: "Foo", Value: "Bar"},
		{Key: "Baz", Value: "Qux"},
	})

	tests := []struct {
		key      string
		expected string
	}{
		{"Foo", "Bar"},
		{"Baz", "Qux"},
	}

	for _, tt := range tests {
		got, ok, err := s.Get(tt.key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !ok {
			t.Fatalf("expected key %q to exist", tt.key)
		}

		if got != tt.expected {
			t.Fatalf("expected %q, got %q", tt.expected, got)
		}
	}
}

func TestStore_Mset_OverwritesExistingKeys(t *testing.T) {
	s := New()

	s.Set("Foo", "Old", 0)

	s.Mset([]KeyValue{
		{Key: "Foo", Value: "New"},
	})

	got, ok, err := s.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected key to exist")
	}

	if got != "New" {
		t.Fatalf("expected %q, got %q", "New", got)
	}
}
func TestStore_Get_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type: ListType,
		List: List{
			values: []string{"Bar"},
		},
	}

	_, _, err := s.Get("Foo")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
func TestStore_Append_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type: ListType,
		List: List{
			values: []string{"Bar"},
		},
	}

	_, err := s.Append("Foo", "Baz")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
func TestStore_Strlen_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type: ListType,
		List: List{
			values: []string{"Bar"},
		},
	}

	_, err := s.Strlen("Foo")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
func TestStore_Setnx_WrongType(t *testing.T) {
	s := New()

	s.data["Foo"] = Entry{
		Type: ListType,
		List: List{
			values: []string{"Bar"},
		},
	}

	_, err := s.Setnx("Foo", "Baz")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
