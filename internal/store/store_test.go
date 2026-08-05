package store

import (
	"fmt"
	"slices"
	"sync"
	"testing"
)

func TestStore_SetGet_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar")
	expected := "Bar"
	got, ok := c.Get("Foo")
	if !ok {
		t.Fatal("expected Get to return ok=true")
	}
	if got != expected {
		t.Fatalf("expected value %q, got %q", expected, got)
	}
}
func TestStore_SetOverwrite_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar")
	c.Set("Foo", "Baz")

	got, ok := c.Get("Foo")
	if !ok {
		t.Fatal("expected Get to return ok=true")
	}
	if got != "Baz" {
		t.Fatalf("expected overwritten value %q, got %q", "Baz", got)
	}
}
func TestStore_SetDoesNotCreateExtraKeys(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar")
	c.Set("Foo", "Baz")

	if _, ok := c.Get("Missing"); ok {
		t.Fatal("unexpected key found")
	}
}
func TestStore_MissingKey_Failure(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar")
	got, ok := c.Get("X")
	if ok {
		t.Fatal("expected Get to return ok=false")
	}
	if got != "" {
		t.Fatalf("expected empty string for missing key, got %q", got)
	}
}

func TestStore_DeleteExistingData_Success(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar")
	got := c.Delete("Foo")
	if got != 1 {
		t.Fatalf("expected delete to return 1, got %d", got)
	}
	value, ok := c.Get("Foo")
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

	c.Set("Foo", "Bar")
	c.Set("Hello", "World")

	deleted := c.Delete("Foo")
	if deleted != 1 {
		t.Fatalf("expected Delete to return 1, got %d", deleted)
	}
	if _, ok := c.Get("Foo"); ok {
		t.Fatal("expected Foo to be deleted")
	}
	value, ok := c.Get("Hello")
	if !ok {
		t.Fatal("expected Hello to still exist")
	}
	if value != "World" {
		t.Fatalf("expected World, got %q", value)
	}
}
func TestStore_Exists_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar")
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

func TestStore_ConcurrentSetAndGet(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	const workers = 100

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", i)
			value := fmt.Sprintf("value-%d", i)
			c.Set(key, value)
			got, ok := c.Get(key)
			if !ok {
				t.Errorf("worker %d: key missing", i)
				return
			}
			if got != value {
				t.Errorf("worker %d: expected %q, got %q", i, value, got)
			}
		}(i)
	}
	wg.Wait()
}

func TestStore_ConcurrentSetSameKey(t *testing.T) {
	c := New()
	const workers = 100
	values := make([]string, workers)
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		values[i] = fmt.Sprintf("value-%d", i)
		go func(i int) {
			defer wg.Done()
			c.Set("shared", values[i])
		}(i)
	}
	wg.Wait()
	got, ok := c.Get("shared")
	if !ok {
		t.Fatal("expected shared key to exist")
	}
	if !slices.Contains(values, got) {
		t.Fatalf("expected final value to be one of the written values, got %q", got)
	}

}
