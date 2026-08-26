package store

import (
	"reflect"
	"testing"
	"time"
)

func TestStore_HSet_NewField(t *testing.T) {
	s := New()

	added, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if added != 1 {
		t.Fatalf("expected 1, got %d", added)
	}

	entry, ok := s.data["user"]
	if !ok {
		t.Fatal("expected key to exist")
	}

	if entry.Type != HashType {
		t.Fatalf("expected HashType, got %v", entry.Type)
	}

	if entry.Hash.values["name"] != "David" {
		t.Fatalf("expected David, got %q", entry.Hash.values["name"])
	}
}

func TestStore_HSet_UpdateField(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	added, err := s.HSet("user", "name", "Dave")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if added != 0 {
		t.Fatalf("expected 0, got %d", added)
	}

	value, found, err := s.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !found {
		t.Fatal("expected field to exist")
	}

	if value != "Dave" {
		t.Fatalf("expected Dave, got %q", value)
	}
}

func TestStore_HSet_MultipleFields(t *testing.T) {
	s := New()

	tests := []struct {
		field string
		value string
	}{
		{"name", "David"},
		{"age", "30"},
		{"city", "Tokyo"},
	}

	for _, tt := range tests {
		added, err := s.HSet("user", tt.field, tt.value)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if added != 1 {
			t.Fatalf("expected 1 for new field, got %d", added)
		}
	}

	if len(s.data["user"].Hash.values) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(s.data["user"].Hash.values))
	}
}

func TestStore_HSet_WrongType(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:   StringType,
		String: "hello",
	}

	_, err := s.HSet("user", "name", "David")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_HGet(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, found, err := s.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !found {
		t.Fatal("expected field to exist")
	}

	if value != "David" {
		t.Fatalf("expected David, got %q", value)
	}
}

func TestStore_HGet_MissingKey(t *testing.T) {
	s := New()

	value, found, err := s.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found {
		t.Fatal("expected field not to exist")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
}

func TestStore_HGet_MissingField(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, found, err := s.HGet("user", "age")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found {
		t.Fatal("expected field not to exist")
	}

	if value != "" {
		t.Fatalf("expected empty value, got %q", value)
	}
}

func TestStore_HGet_WrongType(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:   StringType,
		String: "hello",
	}

	_, _, err := s.HGet("user", "name")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_HDel(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.HSet("user", "age", "30")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	removed, err := s.HDel("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if removed != 1 {
		t.Fatalf("expected 1, got %d", removed)
	}

	_, found, err := s.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found {
		t.Fatal("expected field to be deleted")
	}
}

func TestStore_HDel_MultipleFields(t *testing.T) {
	s := New()

	fields := map[string]string{
		"name": "David",
		"age":  "30",
		"city": "Tokyo",
	}

	for field, value := range fields {
		_, err := s.HSet("user", field, value)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	removed, err := s.HDel("user", "name", "age")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if removed != 2 {
		t.Fatalf("expected 2, got %d", removed)
	}

	if len(s.data["user"].Hash.values) != 1 {
		t.Fatalf("expected 1 remaining field, got %d", len(s.data["user"].Hash.values))
	}
}

func TestStore_HDel_MissingField(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	removed, err := s.HDel("user", "age")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if removed != 0 {
		t.Fatalf("expected 0, got %d", removed)
	}
}

func TestStore_HDel_RemovesKeyWhenEmpty(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	removed, err := s.HDel("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if removed != 1 {
		t.Fatalf("expected 1, got %d", removed)
	}

	if _, ok := s.data["user"]; ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestStore_HDel_MissingKey(t *testing.T) {
	s := New()

	removed, err := s.HDel("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if removed != 0 {
		t.Fatalf("expected 0, got %d", removed)
	}
}

func TestStore_HDel_WrongType(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:   StringType,
		String: "hello",
	}

	_, err := s.HDel("user", "name")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_HExists(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exists, err := s.HExists("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !exists {
		t.Fatal("expected field to exist")
	}
}

func TestStore_HExists_MissingField(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	exists, err := s.HExists("user", "age")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Fatal("expected field not to exist")
	}
}

func TestStore_HExists_MissingKey(t *testing.T) {
	s := New()

	exists, err := s.HExists("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if exists {
		t.Fatal("expected field not to exist")
	}
}

func TestStore_HExists_WrongType(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:   StringType,
		String: "hello",
	}

	_, err := s.HExists("user", "name")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_HLen(t *testing.T) {
	s := New()

	fields := map[string]string{
		"name": "David",
		"age":  "30",
		"city": "Tokyo",
	}

	for field, value := range fields {
		_, err := s.HSet("user", field, value)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	length, err := s.HLen("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if length != 3 {
		t.Fatalf("expected 3, got %d", length)
	}
}

func TestStore_HLen_MissingKey(t *testing.T) {
	s := New()

	length, err := s.HLen("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if length != 0 {
		t.Fatalf("expected 0, got %d", length)
	}
}

func TestStore_HLen_WrongType(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:   StringType,
		String: "hello",
	}

	_, err := s.HLen("user")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_HGetAll(t *testing.T) {
	s := New()

	expected := map[string]string{
		"name": "David",
		"age":  "30",
		"city": "Tokyo",
	}

	for field, value := range expected {
		_, err := s.HSet("user", field, value)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	got, err := s.HGetAll("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestStore_HGetAll_MissingKey(t *testing.T) {
	s := New()

	got, err := s.HGetAll("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected empty map, got %v", got)
	}
}

func TestStore_HGetAll_WrongType(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:   StringType,
		String: "hello",
	}

	_, err := s.HGetAll("user")

	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_HGetAll_ReturnsCopy(t *testing.T) {
	s := New()

	_, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.HGetAll("user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got["name"] = "Modified"

	value, _, err := s.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if value != "David" {
		t.Fatalf("expected original value David, got %q", value)
	}
}

func TestStore_HSet_ExpiredKey(t *testing.T) {
	s := New()

	s.data["user"] = Entry{
		Type:      HashType,
		Hash:      Hash{values: map[string]string{"name": "Old"}},
		ExpiresAt: time.Now().Add(-time.Second),
	}

	added, err := s.HSet("user", "name", "David")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if added != 1 {
		t.Fatalf("expected 1, got %d", added)
	}

	value, found, err := s.HGet("user", "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !found {
		t.Fatal("expected field to exist")
	}

	if value != "David" {
		t.Fatalf("expected David, got %q", value)
	}
}
