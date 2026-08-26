package store

import (
	"testing"
	"time"
)

func TestStore_SetWithTTL(t *testing.T) {
	c := New()
	defer c.Close()
	c.Set("Foo", "Bar", 50*time.Millisecond)
	_, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected stored value to still exist")
	}
	time.Sleep(75 * time.Millisecond)
	_, ok, err = c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected key to expire")
	}
}

func TestStore_GetExpiredKey(t *testing.T) {
	c := New()
	defer c.Close()
	c.Set("Foo", "Bar", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)
	_, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected store value to be deleted")
	}
}

func TestStore_ExpiredKeyIsDeleted(t *testing.T) {
	c := New()
	defer c.Close()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)
	_, ok, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
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
	defer c.Close()
	c.Set("Foo", "Bar", 0)
	if got := c.TTL("Foo"); got != -1 {
		t.Fatalf("TTL = %d, want -1", got)
	}
}
func TestStore_TTL_MissingKey(t *testing.T) {
	c := New()
	defer c.Close()
	if got := c.TTL("missing"); got != -2 {
		t.Fatalf("TTL = %d, want -2", got)
	}
}
func TestStore_TTL_WithExpiration(t *testing.T) {
	c := New()
	defer c.Close()
	c.Set("Foo", "Bar", 5*time.Second)
	got := c.TTL("Foo")
	if got < 4 || got > 5 {
		t.Fatalf("TTL = %d, want approximately 5", got)
	}
}
func TestStore_TTL_ExpiredKey(t *testing.T) {
	c := New()
	defer c.Close()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	got := c.TTL("Foo")
	if got != -2 {
		t.Fatalf("TTL = %d, want -2", got)
	}
}
func TestStore_Delete_ExpiredKey(t *testing.T) {
	c := New()
	defer c.Close()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	time.Sleep(20 * time.Millisecond)

	got := c.Delete("Foo")
	if got != 0 {
		t.Fatalf("Delete = %d, want 0", got)
	}
}

func TestStore_Expire_ExistingKey(t *testing.T) {
	c := New()
	defer c.Close()
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
	defer c.Close()
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
	defer c.Close()
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
	defer c.Close()
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
	defer c.Close()
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
	defer c.Close()
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
	defer c.Close()
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
	defer c.Close()
	c.Set("Foo", "Bar", 10*time.Millisecond)
	ok := c.Persist("Foo")
	if !ok {
		t.Fatal("expected persist to succeed")
	}
	time.Sleep(20 * time.Millisecond)
	value, exists, err := c.Get("Foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists || value != "Bar" {
		t.Fatal("expected persisted key to remain after original expiration time")
	}
}

func TestStore_ActiveExpiration_DeletesExpiredKey(t *testing.T) {
	c := New()
	defer c.Close()

	c.Set("Foo", "Bar", 50*time.Millisecond)

	time.Sleep(1100 * time.Millisecond)

	c.mu.RLock()
	_, exists := c.data["Foo"]
	c.mu.RUnlock()

	if exists {
		t.Fatal("expected active expiration to delete key")
	}
}

func TestStore_ActiveExpiration_PreservesUnexpiredKey(t *testing.T) {
	c := New()
	defer c.Close()

	c.Set("Foo", "Bar", 2*time.Second)

	time.Sleep(1100 * time.Millisecond)

	c.mu.RLock()
	_, exists := c.data["Foo"]
	c.mu.RUnlock()

	if !exists {
		t.Fatal("expected unexpired key to remain")
	}
}

func TestStore_ActiveExpiration_DeletesMultipleExpiredKeys(t *testing.T) {
	c := New()
	defer c.Close()

	c.Set("Foo", "Bar", 50*time.Millisecond)
	c.Set("Baz", "Qux", 50*time.Millisecond)

	time.Sleep(1100 * time.Millisecond)

	c.mu.RLock()
	_, fooExists := c.data["Foo"]
	_, bazExists := c.data["Baz"]
	c.mu.RUnlock()

	if fooExists || bazExists {
		t.Fatal("expected all expired keys to be deleted")
	}
}

func TestStore_Close_StopsExpirationLoop(t *testing.T) {
	c := New()

	c.Set("Foo", "Bar", 50*time.Millisecond)

	c.Close()

	// The important behavior here is that Close returns successfully.
	// The expiration goroutine should have stopped.
}
