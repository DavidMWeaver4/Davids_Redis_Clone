package store

import (
	"math"
	"testing"

	"github.com/DavidMWeaver4/Davids_Redis_Clone/internal/store/skiplist"
)

func TestStore_ZAdd_NewMember(t *testing.T) {
	s := New()

	added, err := s.ZAdd("Scores", 10, "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if added != 1 {
		t.Fatalf("ZAdd = %d, want 1", added)
	}

	score, found, err := s.ZScore("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected member to exist")
	}
	if score != 10 {
		t.Fatalf("score = %v, want 10", score)
	}
}

func TestStore_ZAdd_ExistingMemberSameScore(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")

	added, err := s.ZAdd("Scores", 10, "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if added != 0 {
		t.Fatalf("ZAdd = %d, want 0", added)
	}

	score, found, err := s.ZScore("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected member to exist")
	}
	if score != 10 {
		t.Fatalf("score = %v, want 10", score)
	}
}

func TestStore_ZAdd_ExistingMemberUpdatesScore(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")

	added, err := s.ZAdd("Scores", 20, "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if added != 0 {
		t.Fatalf("ZAdd = %d, want 0", added)
	}

	score, found, err := s.ZScore("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected member to exist")
	}
	if score != 20 {
		t.Fatalf("score = %v, want 20", score)
	}
}

func TestStore_ZAdd_MultipleMembers(t *testing.T) {
	s := New()

	tests := []struct {
		member string
		score  float64
	}{
		{"Alice", 10},
		{"Bob", 20},
		{"Charlie", 30},
	}

	for _, tt := range tests {
		added, err := s.ZAdd("Scores", tt.score, tt.member)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if added != 1 {
			t.Fatalf("ZAdd = %d, want 1", added)
		}
	}

	length, err := s.ZCard("Scores")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if length != 3 {
		t.Fatalf("ZCard = %d, want 3", length)
	}
}

func TestStore_ZAdd_InvalidScore(t *testing.T) {
	s := New()

	_, err := s.ZAdd("Scores", math.NaN(), "Alice")
	if err != ErrInvalidScore {
		t.Fatalf("expected ErrInvalidScore, got %v", err)
	}
}

func TestStore_ZAdd_WrongType(t *testing.T) {
	s := New()

	s.Set("Scores", "hello", 0)

	_, err := s.ZAdd("Scores", 10, "Alice")
	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_ZScore_ExistingMember(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 42.5, "Alice")

	score, found, err := s.ZScore("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected member to exist")
	}
	if score != 42.5 {
		t.Fatalf("score = %v, want 42.5", score)
	}
}

func TestStore_ZScore_MissingMember(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")

	score, found, err := s.ZScore("Scores", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected member to not exist")
	}
	if score != 0 {
		t.Fatalf("score = %v, want 0", score)
	}
}

func TestStore_ZScore_MissingKey(t *testing.T) {
	s := New()

	score, found, err := s.ZScore("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected member to not exist")
	}
	if score != 0 {
		t.Fatalf("score = %v, want 0", score)
	}
}

func TestStore_ZScore_WrongType(t *testing.T) {
	s := New()

	s.Set("Scores", "hello", 0)

	_, _, err := s.ZScore("Scores", "Alice")
	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_ZCard_ExistingKey(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")
	s.ZAdd("Scores", 30, "Charlie")

	length, err := s.ZCard("Scores")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if length != 3 {
		t.Fatalf("ZCard = %d, want 3", length)
	}
}

func TestStore_ZCard_MissingKey(t *testing.T) {
	s := New()

	length, err := s.ZCard("Scores")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if length != 0 {
		t.Fatalf("ZCard = %d, want 0", length)
	}
}

func TestStore_ZCard_WrongType(t *testing.T) {
	s := New()

	s.Set("Scores", "hello", 0)

	_, err := s.ZCard("Scores")
	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_ZRem_ExistingMember(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")

	removed, err := s.ZRem("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 1 {
		t.Fatalf("ZRem = %d, want 1", removed)
	}

	_, found, err := s.ZScore("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected Alice to be removed")
	}

	length, err := s.ZCard("Scores")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if length != 1 {
		t.Fatalf("ZCard = %d, want 1", length)
	}
}

func TestStore_ZRem_MissingMember(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")

	removed, err := s.ZRem("Scores", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 0 {
		t.Fatalf("ZRem = %d, want 0", removed)
	}
}

func TestStore_ZRem_MultipleMembers(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")
	s.ZAdd("Scores", 30, "Charlie")

	removed, err := s.ZRem("Scores", "Alice", "Charlie")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 2 {
		t.Fatalf("ZRem = %d, want 2", removed)
	}

	length, err := s.ZCard("Scores")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if length != 1 {
		t.Fatalf("ZCard = %d, want 1", length)
	}
}

func TestStore_ZRem_MissingKey(t *testing.T) {
	s := New()

	removed, err := s.ZRem("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 0 {
		t.Fatalf("ZRem = %d, want 0", removed)
	}
}

func TestStore_ZRem_WrongType(t *testing.T) {
	s := New()

	s.Set("Scores", "hello", 0)

	_, err := s.ZRem("Scores", "Alice")
	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}

func TestStore_ZRem_DeletesEmptyKey(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")

	removed, err := s.ZRem("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if removed != 1 {
		t.Fatalf("ZRem = %d, want 1", removed)
	}

	_, exists := s.data["Scores"]
	if exists {
		t.Fatal("expected empty sorted set key to be deleted")
	}
}

func TestStore_ZRange(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 30, "Charlie")
	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")

	got, err := s.ZRange("Scores", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []skiplist.MemberScore{
		{Member: "Alice", Score: 10},
		{Member: "Bob", Score: 20},
		{Member: "Charlie", Score: 30},
	}

	if len(got) != len(want) {
		t.Fatalf("ZRange length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ZRange[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestStore_ZRange_SubRange(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")
	s.ZAdd("Scores", 30, "Charlie")
	s.ZAdd("Scores", 40, "David")

	got, err := s.ZRange("Scores", 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []skiplist.MemberScore{
		{Member: "Bob", Score: 20},
		{Member: "Charlie", Score: 30},
	}

	if len(got) != len(want) {
		t.Fatalf("ZRange length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ZRange[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestStore_ZRange_NegativeIndexes(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")
	s.ZAdd("Scores", 30, "Charlie")

	got, err := s.ZRange("Scores", -2, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []skiplist.MemberScore{
		{Member: "Bob", Score: 20},
		{Member: "Charlie", Score: 30},
	}

	if len(got) != len(want) {
		t.Fatalf("ZRange length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ZRange[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestStore_ZRange_PastEnd(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")

	got, err := s.ZRange("Scores", 0, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []skiplist.MemberScore{
		{Member: "Alice", Score: 10},
		{Member: "Bob", Score: 20},
	}

	if len(got) != len(want) {
		t.Fatalf("ZRange length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ZRange[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestStore_ZRange_InvalidRange(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")

	got, err := s.ZRange("Scores", 2, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("ZRange = %v, want empty result", got)
	}
}

func TestStore_ZRange_MissingKey(t *testing.T) {
	s := New()

	got, err := s.ZRange("Scores", 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("ZRange = %v, want empty result", got)
	}
}

func TestStore_ZRange_WrongType(t *testing.T) {
	s := New()

	s.Set("Scores", "hello", 0)

	_, err := s.ZRange("Scores", 0, -1)
	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
func TestStore_ZRank_ExistingMember(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 30, "Charlie")
	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")

	rank, found, err := s.ZRank("Scores", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected member to exist")
	}
	if rank != 1 {
		t.Fatalf("ZRank = %d, want 1", rank)
	}
}

func TestStore_ZRank_MissingMember(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")

	rank, found, err := s.ZRank("Scores", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected member to not exist")
	}
	if rank != 0 {
		t.Fatalf("ZRank = %d, want 0", rank)
	}
}

func TestStore_ZRank_MissingKey(t *testing.T) {
	s := New()

	rank, found, err := s.ZRank("Scores", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected member to not exist")
	}
	if rank != 0 {
		t.Fatalf("ZRank = %d, want 0", rank)
	}
}

func TestStore_ZRank_WrongType(t *testing.T) {
	s := New()

	s.Set("Scores", "hello", 0)

	_, _, err := s.ZRank("Scores", "Alice")
	if err != ErrWrongType {
		t.Fatalf("expected ErrWrongType, got %v", err)
	}
}
func TestStore_ZRank_TiedScores(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Charlie")
	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 10, "Bob")

	rank, found, err := s.ZRank("Scores", "Bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Fatal("expected member to exist")
	}
	if rank != 1 {
		t.Fatalf("ZRank = %d, want 1", rank)
	}
}
func TestStore_ZRangeByScore(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 30, "Charlie")
	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")

	got, err := s.ZRangeByScore("Scores", 10, 30, 0, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []skiplist.MemberScore{
		{Member: "Alice", Score: 10},
		{Member: "Bob", Score: 20},
		{Member: "Charlie", Score: 30},
	}

	if len(got) != len(want) {
		t.Fatalf("ZRangeByScore length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ZRangeByScore[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
func TestStore_ZRangeByScore_SubRange(t *testing.T) {
	s := New()

	s.ZAdd("Scores", 10, "Alice")
	s.ZAdd("Scores", 20, "Bob")
	s.ZAdd("Scores", 30, "Charlie")
	s.ZAdd("Scores", 40, "David")

	got, err := s.ZRangeByScore("Scores", 10, 40, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []skiplist.MemberScore{
		{Member: "Bob", Score: 20},
		{Member: "Charlie", Score: 30},
	}

	if len(got) != len(want) {
		t.Fatalf("ZRangeByScore length = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ZRangeByScore[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}
