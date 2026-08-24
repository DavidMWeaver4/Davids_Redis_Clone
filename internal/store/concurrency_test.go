package store

import (
	"fmt"
	"slices"
	"strconv"
	"sync"
	"testing"
)

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
			got, ok, err := c.Get(key)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
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
	got, ok, err := c.Get("shared")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected shared key to exist")
	}
	if !slices.Contains(values, got) {
		t.Fatalf("expected final value to be one of the written values, got %q", got)
	}

}

func TestStore_Incr_Concurrent(t *testing.T) {
	c := New()
	c.Set("Foo", "0", 0)

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			_, err := c.Incr("Foo")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "100" {
		t.Fatalf("expected 100, got %q", value)
	}
}
func TestStore_Decr_Concurrent(t *testing.T) {
	c := New()
	c.Set("Foo", "0", 0)

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			_, err := c.Decr("Foo")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "-100" {
		t.Fatalf("expected -100, got %q", value)
	}
}
func TestStore_Incrby_Concurrent(t *testing.T) {
	c := New()
	c.Set("Foo", "0", 0)

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			_, err := c.Incrby("Foo", 5)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "500" {
		t.Fatalf("expected 500, got %q", value)
	}
}

func TestStore_Decrby_Concurrent(t *testing.T) {
	c := New()
	c.Set("Foo", "0", 0)

	const workers = 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()

			_, err := c.Decrby("Foo", 5)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	value, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "-500" {
		t.Fatalf("expected -500, got %q", value)
	}
}
func TestStore_ConcurrentAccess(t *testing.T) {
	s := New()

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			s.Set("Foo", strconv.Itoa(i), 0)
			s.Get("Foo")
			s.Exists("Foo")
			s.Strlen("Foo")
		}(i)
	}

	wg.Wait()
}
