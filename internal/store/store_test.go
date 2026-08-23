package store

import (
	"testing"
)

func TestStore_SetGet_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
	expected := "Bar"
	got, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Get to return ok=true")
	}
	if got != expected {
		t.Fatalf("expected value %q, got %q", expected, got)
	}
}
func TestStore_SetOverwrite_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
	c.Set("Foo", "Baz", 0)

	got, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Get to return ok=true")
	}
	if got != "Baz" {
		t.Fatalf("expected overwritten value %q, got %q", "Baz", got)
	}
}
func TestStore_SetDoesNotCreateExtraKeys(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
	c.Set("Foo", "Baz", 0)

	_, ok, err := c.Get("Missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {

		t.Fatal("unexpected key found")
	}
}
func TestStore_MissingKey_Failure(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
	got, ok, err := c.Get("X")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Get to return ok=false")
	}
	if got != "" {
		t.Fatalf("expected empty string for missing key, got %q", got)
	}
}

func TestStore_DeleteExistingData_Success(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 0)
	got := c.Delete("Foo")
	if got != 1 {
		t.Fatalf("expected delete to return 1, got %d", got)
	}
	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected key to be deleted")
	}
	if value != "" {
		t.Fatalf("expected deleted value to be empty, got %q", value)
	}
}

func TestStore_DeleteMissing_Failure(t *testing.T) {
	c := New()
	got := c.Delete("Foo")
	if got != 0 {
		t.Fatalf("expected Delete to return 0, got %d", got)
	}
}
func TestStore_DeleteOnlySpecifiedKey(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
	c.Set("Hello", "World", 0)

	deleted := c.Delete("Foo")
	if deleted != 1 {
		t.Fatalf("expected Delete to return 1, got %d", deleted)
	}
	_, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected Foo to be deleted")
	}
	value, ok, err := c.Get("Hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Hello to still exist")
	}
	if value != "World" {
		t.Fatalf("expected World, got %q", value)
	}
}
func TestStore_Exists_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
	exists := c.Exists("Foo")
	if !exists {
		t.Fatal("expected Exists to return true")
	}
}

func TestStore_Exists_Failure(t *testing.T) {
	c := New()

	exists := c.Exists("Foo")
	if exists {
		t.Fatal("expected Exists to return false")
	}
}
