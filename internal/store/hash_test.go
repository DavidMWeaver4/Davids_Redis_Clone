package store

import (
	"reflect"
	"testing"
)

func TestHash_HSet(t *testing.T) {
	h := Hash{}

	if got := h.HSet("name", "David"); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}

	if got := h.HSet("name", "Dave"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestHash_HGet(t *testing.T) {
	h := Hash{
		values: map[string]string{
			"name": "David",
		},
	}

	value, ok := h.HGet("name")
	if !ok {
		t.Fatal("expected field to exist")
	}

	if value != "David" {
		t.Fatalf("expected David, got %q", value)
	}
}

func TestHash_HDel(t *testing.T) {
	h := Hash{
		values: map[string]string{
			"name": "David",
			"age":  "30",
		},
	}

	removed := h.HDel("name", "missing")

	if removed != 1 {
		t.Fatalf("expected 1, got %d", removed)
	}

	if _, ok := h.values["name"]; ok {
		t.Fatal("expected name to be deleted")
	}
}

func TestHash_HExists(t *testing.T) {
	h := Hash{
		values: map[string]string{
			"name": "David",
		},
	}

	if !h.HExists("name") {
		t.Fatal("expected name to exist")
	}

	if h.HExists("age") {
		t.Fatal("expected age not to exist")
	}
}

func TestHash_HLen(t *testing.T) {
	h := Hash{
		values: map[string]string{
			"name": "David",
			"age":  "30",
		},
	}

	if got := h.HLen(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestHash_HGetAll(t *testing.T) {
	expected := map[string]string{
		"name": "David",
		"age":  "30",
	}

	h := Hash{
		values: expected,
	}

	got := h.HGetAll()

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected %v, got %v", expected, got)
	}

	got["name"] = "Modified"

	if h.values["name"] != "David" {
		t.Fatal("expected HGetAll to return a copy")
	}
}
