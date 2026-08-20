package store

import (
	"fmt"
	"math"
	"slices"
	"strconv"
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
func TestStore_Delete_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	got := c.Delete("Foo")
	if got != 0 {
		t.Fatalf("Delete = %d, want 0", got)
	}
}

func TestStore_Expire_ExistingKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	ok := c.Expire("Foo", 5*time.Second)
	if !ok {
		t.Fatal("expected Expire to succeed")
	}
	_, exists := c.data["Foo"]
	if !exists {
		t.Fatal("expected key to still exist")
	}
	ttl := c.TTL("Foo")
	if ttl < 4 || ttl > 5 {
		t.Fatalf("expected TTL around 5, got %d", ttl)
	}
}

func TestStore_Expire_MissingKey(t *testing.T) {
	c := New()
	ok := c.Expire("Foo", 5*time.Second)
	if ok {
		t.Fatal("expected Expire to fail for missing key")
	}
	_, exists := c.data["Foo"]
	if exists {
		t.Fatal("expect missing key to not exists")
	}
}

func TestStore_Expire_AlreadyExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	ok := c.Expire("Foo", 5*time.Second)
	if ok {
		t.Fatalf("expected Expire to fail for expired key")
	}
	_, exists := c.data["Foo"]
	if exists {
		t.Fatal("expected expired key to have been deleted")
	}
}

func TestStore_Persist_ExistingKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	ok := c.Persist("Foo")
	if !ok {
		t.Fatal("expected persist to succeed")
	}
	_, exists := c.data["Foo"]
	if !exists {
		t.Fatal("expected key to exist")
	}
	ttl := c.TTL("Foo")
	if ttl != -1 {
		t.Fatalf("expected ttl to be -1, got %v", ttl)
	}
}
func TestStore_Persist_MissingKey(t *testing.T) {
	c := New()
	ok := c.Persist("Foo")
	if ok {
		t.Fatal("expected Persist to fail for missing key")
	}
	_, exists := c.data["Foo"]
	if exists {
		t.Fatal("expect missing key to not exists")
	}
}
func TestStore_Persist_AlreadyPersistant(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 0)
	prettl := c.TTL("Foo")
	if prettl != -1 {
		t.Fatal("setup failure")
	}
	ok := c.Persist("Foo")
	if ok {
		t.Fatal("expected persist to not succeed")
	}
	_, exists := c.data["Foo"]
	if !exists {
		t.Fatal("expected key to exist")
	}
	ttl := c.TTL("Foo")
	if ttl != -1 {
		t.Fatalf("expected ttl to be -1, got %v", ttl)
	}
}
func TestStore_Persist_ExpiredKey(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 5*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	ok := c.Persist("Foo")
	if ok {
		t.Fatal("expected persist to not succeed")
	}
	_, exists := c.data["Foo"]
	if exists {
		t.Fatal("expected expired key to be deleted")
	}
}
func TestStore_Persist_CheckAfterOriginalExpirationTime(t *testing.T) {
	c := New()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	ok := c.Persist("Foo")
	if !ok {
		t.Fatal("expected persist to succeed")
	}
	time.Sleep(20 * time.Millisecond)
	value, exists := c.Get("Foo")
	if !exists || value != "Bar" {
		t.Fatal("expected persisted key to remain after original expiration time")
	}
}
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

	value, ok := c.Get("Foo")
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

	value, ok := c.Get("Foo")
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
}
func TestStore_Incr_Overflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MaxInt64, 10), 0)

	_, err := c.Incr("Foo")
	if err == nil {
		t.Fatal("expected overflow error")
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

	value, ok := c.Get("Foo")
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
}
func TestStore_Decr_Overflow(t *testing.T) {
	c := New()
	c.Set("Foo", strconv.FormatInt(math.MinInt64, 10), 0)

	_, err := c.Decr("Foo")
	if err == nil {
		t.Fatal("expected overflow error")
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
func TestStore_ConcurrentIncr(t *testing.T) {
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

	value, ok := c.Get("Foo")
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "100" {
		t.Fatalf("expected 100, got %q", value)
	}
}
func TestStore_ConcurrentDecr(t *testing.T) {
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

	value, ok := c.Get("Foo")
	if !ok {
		t.Fatal("expected Foo to exist")
	}

	if value != "-100" {
		t.Fatalf("expected -100, got %q", value)
	}
}
