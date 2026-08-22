package store

import (
	"testing"
	"time"
)

func TestStore_Append_ExistingKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", 0)

	got := s.Append("Foo", "Baz")

	if got != 6 {
		t.Fatalf("expected length 6, got %d", got)
	}

	value, ok := s.Get("Foo")
	if !ok {
		t.Fatal("expected key to exist")
	}

	if value != "BarBaz" {
		t.Fatalf("expected %q, got %q", "BarBaz", value)
	}
}

func TestStore_Append_MissingKey(t *testing.T) {
	s := New()

	got := s.Append("Foo", "Bar")

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	value, ok := s.Get("Foo")
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

	got := s.Append("Foo", "")

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	value, _ := s.Get("Foo")
	if value != "Bar" {
		t.Fatalf("expected %q, got %q", "Bar", value)
	}
}
func TestStore_Append_ExpiredKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	got := s.Append("Foo", "Baz")

	if got != 3 {
		t.Fatalf("expected length 3, got %d", got)
	}

	value, ok := s.Get("Foo")
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

	got := s.Strlen("Foo")

	if got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestStore_Strlen_MissingKey(t *testing.T) {
	s := New()

	got := s.Strlen("Foo")

	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestStore_Strlen_EmptyValue(t *testing.T) {
	s := New()

	s.Set("Foo", "", 0)

	got := s.Strlen("Foo")

	if got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
func TestStore_Setnx_MissingKey(t *testing.T) {
	s := New()

	ok := s.Setnx("Foo", "Bar")

	if !ok {
		t.Fatal("expected SETNX to succeed")
	}

	value, found := s.Get("Foo")
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

	ok := s.Setnx("Foo", "Baz")

	if ok {
		t.Fatal("expected SETNX to fail")
	}

	value, _ := s.Get("Foo")

	if value != "Bar" {
		t.Fatalf("expected original value %q, got %q", "Bar", value)
	}
}

func TestStore_Setnx_ExpiredKey(t *testing.T) {
	s := New()

	s.Set("Foo", "Bar", time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	ok := s.Setnx("Foo", "Baz")

	if !ok {
		t.Fatal("expected SETNX to succeed for expired key")
	}

	value, found := s.Get("Foo")
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
		got, ok := s.Get(tt.key)

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

	got, ok := s.Get("Foo")
	if !ok {
		t.Fatal("expected key to exist")
	}

	if got != "New" {
		t.Fatalf("expected %q, got %q", "New", got)
	}
}
