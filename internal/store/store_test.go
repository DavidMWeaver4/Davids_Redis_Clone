package store

import (
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestStore_SetGet_Success(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
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

	c.Set("Foo", "Bar", 0)
	c.Set("Foo", "Baz", 0)

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

	c.Set("Foo", "Bar", 0)
	c.Set("Foo", "Baz", 0)

	if _, ok := c.Get("Missing"); ok {
		t.Fatal("unexpected key found")
	}
}
func TestStore_MissingKey_Failure(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 0)
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
	c.Set("Foo", "Bar", 0)
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

	c.Set("Foo", "Bar", 0)
	c.Set("Hello", "World", 0)

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
			c.Set(key, value, 0)
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
			c.Set("shared", values[i], 0)
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
func TestStore_SetWithTTL(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 50*time.Millisecond)
	_, ok := c.Get("Foo")
	if !ok {
		t.Fatal("expected stored value to still exist")
	}
	time.Sleep(75 * time.Millisecond)
	_, ok = c.Get("Foo")
	if ok {
		t.Fatal("expected key to expire")
	}
}

func TestStore_GetExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	_, ok := c.Get("Foo")
	if ok {
		t.Fatal("expected store value to be deleted")
	}
}

func TestStore_ExpiredKeyIsDeleted(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, ok := c.Get("Foo")
	if ok {
		t.Fatal("expected expired key to be unavailable")
	}
	_, exists := c.data["Foo"]
	if exists {
		t.Fatal("expected expired key to be deleted from store")
	}
}
func TestStore_TTL_NoExpiration(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 0)
	if got := c.TTL("Foo"); got != -1 {
		t.Fatalf("TTL = %d, want -1", got)
	}
}
func TestStore_TTL_MissingKey(t *testing.T) {
	c := New()
	if got := c.TTL("missing"); got != -2 {
		t.Fatalf("TTL = %d, want -2", got)
	}
}
func TestStore_TTL_WithExpiration(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 5*time.Second)
	got := c.TTL("Foo")
	if got < 4 || got > 5 {
		t.Fatalf("TTL = %d, want approximately 5", got)
	}
}
func TestStore_TTL_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	got := c.TTL("Foo")
	if got != -2 {
		t.Fatalf("TTL = %d, want -2", got)
	}
}
func TestStore_DeleteExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	got := c.Delete("Foo")
	if got != 0 {
		t.Fatalf("Delete = %d, want 0", got)
	}
}
