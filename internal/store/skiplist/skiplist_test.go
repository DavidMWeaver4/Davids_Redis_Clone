package skiplist

import (
	"reflect"
	"testing"
)

func TestSkipList_Insert(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")

	node := sl.Search(10, "foo")
	if node == nil {
		t.Fatal("expected inserted node to be found")
	}

	if node.member != "foo" {
		t.Fatalf("member = %q, want %q", node.member, "foo")
	}

	if node.score != 10 {
		t.Fatalf("score = %v, want %v", node.score, 10)
	}

	if sl.size != 1 {
		t.Fatalf("size = %d, want 1", sl.size)
	}
}

func TestSkipList_Insert_Multiple(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(30, "baz")
	sl.Insert(10, "foo")
	sl.Insert(20, "bar")

	got := sl.RangeByRank(0, 2)

	want := []MemberScore{
		{Member: "foo", Score: 10},
		{Member: "bar", Score: 20},
		{Member: "baz", Score: 30},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}

	if sl.size != 3 {
		t.Fatalf("size = %d, want 3", sl.size)
	}
}

func TestSkipList_Insert_SameScoreOrdersByMember(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "charlie")
	sl.Insert(10, "alpha")
	sl.Insert(10, "bravo")

	got := sl.RangeByRank(0, 2)

	want := []MemberScore{
		{Member: "alpha", Score: 10},
		{Member: "bravo", Score: 10},
		{Member: "charlie", Score: 10},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestSkipList_Search_Found(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(10, "foo")

	node := sl.Search(10, "foo")
	if node == nil {
		t.Fatal("expected node to be found")
	}
}

func TestSkipList_Search_Missing(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(10, "foo")

	tests := []struct {
		name   string
		score  float64
		member string
	}{
		{
			name:   "missing score",
			score:  20,
			member: "foo",
		},
		{
			name:   "missing member",
			score:  10,
			member: "bar",
		},
		{
			name:   "missing score and member",
			score:  20,
			member: "bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := sl.Search(tt.score, tt.member)
			if node != nil {
				t.Fatal("expected node not to be found")
			}
		})
	}
}

func TestSkipList_Search_Empty(t *testing.T) {
	sl := NewSkipList()

	if node := sl.Search(10, "foo"); node != nil {
		t.Fatal("expected search on empty list to return nil")
	}
}

func TestSkipList_Delete(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")

	if !sl.Delete(10, "foo") {
		t.Fatal("expected Delete to succeed")
	}

	if sl.Search(10, "foo") != nil {
		t.Fatal("expected deleted node to be absent")
	}

	if sl.Search(20, "bar") == nil {
		t.Fatal("expected remaining node to still exist")
	}

	if sl.size != 1 {
		t.Fatalf("size = %d, want 1", sl.size)
	}
}

func TestSkipList_Delete_Missing(t *testing.T) {
	sl := NewSkipList()
	sl.Insert(10, "foo")

	if sl.Delete(20, "bar") {
		t.Fatal("expected Delete to fail")
	}

	if sl.size != 1 {
		t.Fatalf("size = %d, want 1", sl.size)
	}
}

func TestSkipList_Delete_Empty(t *testing.T) {
	sl := NewSkipList()

	if sl.Delete(10, "foo") {
		t.Fatal("expected Delete on empty list to fail")
	}

	if sl.size != 0 {
		t.Fatalf("size = %d, want 0", sl.size)
	}
}

func TestSkipList_Delete_SameScoreDifferentMember(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(10, "bar")

	if !sl.Delete(10, "foo") {
		t.Fatal("expected Delete to succeed")
	}

	if sl.Search(10, "foo") != nil {
		t.Fatal("expected foo to be deleted")
	}

	if sl.Search(10, "bar") == nil {
		t.Fatal("expected bar to remain")
	}

	if sl.size != 1 {
		t.Fatalf("size = %d, want 1", sl.size)
	}
}

func TestSkipList_RangeByRank(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")
	sl.Insert(30, "baz")
	sl.Insert(40, "qux")

	tests := []struct {
		name  string
		start int
		end   int
		want  []MemberScore
	}{
		{
			name:  "first two",
			start: 0,
			end:   1,
			want: []MemberScore{
				{Member: "foo", Score: 10},
				{Member: "bar", Score: 20},
			},
		},
		{
			name:  "middle",
			start: 1,
			end:   2,
			want: []MemberScore{
				{Member: "bar", Score: 20},
				{Member: "baz", Score: 30},
			},
		},
		{
			name:  "last",
			start: 3,
			end:   3,
			want: []MemberScore{
				{Member: "qux", Score: 40},
			},
		},
		{
			name:  "past end",
			start: 10,
			end:   20,
			want:  []MemberScore{},
		},
		{
			name:  "invalid range",
			start: 3,
			end:   1,
			want:  []MemberScore{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sl.RangeByRank(tt.start, tt.end)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSkipList_RangeByRank_Empty(t *testing.T) {
	sl := NewSkipList()

	got := sl.RangeByRank(0, 10)

	if !reflect.DeepEqual(got, []MemberScore{}) {
		t.Fatalf("got %+v, want empty slice", got)
	}
}

func TestSkipList_Range(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")
	sl.Insert(30, "baz")
	sl.Insert(40, "qux")

	tests := []struct {
		name     string
		minScore float64
		maxScore float64
		offset   int
		count    int
		want     []MemberScore
	}{
		{
			name:     "entire score range",
			minScore: 10,
			maxScore: 40,
			offset:   0,
			count:    -1,
			want: []MemberScore{
				{Member: "foo", Score: 10},
				{Member: "bar", Score: 20},
				{Member: "baz", Score: 30},
				{Member: "qux", Score: 40},
			},
		},
		{
			name:     "score range",
			minScore: 20,
			maxScore: 30,
			offset:   0,
			count:    -1,
			want: []MemberScore{
				{Member: "bar", Score: 20},
				{Member: "baz", Score: 30},
			},
		},
		{
			name:     "score range with offset and count",
			minScore: 10,
			maxScore: 40,
			offset:   1,
			count:    2,
			want: []MemberScore{
				{Member: "bar", Score: 20},
				{Member: "baz", Score: 30},
			},
		},
		{
			name:     "no matching scores",
			minScore: 50,
			maxScore: 60,
			offset:   0,
			count:    -1,
			want:     []MemberScore{},
		},
		{
			name:     "count zero",
			minScore: 10,
			maxScore: 40,
			offset:   0,
			count:    0,
			want:     []MemberScore{},
		},
		{
			name:     "offset past end",
			minScore: 10,
			maxScore: 40,
			offset:   10,
			count:    2,
			want:     []MemberScore{},
		},
		{
			name:     "offset with unlimited count",
			minScore: 10,
			maxScore: 40,
			offset:   2,
			count:    -1,
			want: []MemberScore{
				{Member: "baz", Score: 30},
				{Member: "qux", Score: 40},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sl.Range(
				tt.minScore,
				tt.maxScore,
				tt.offset,
				tt.count,
			)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}
func TestSkipList_Range_InvalidOffset(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")

	got := sl.Range(10, 20, -1, 2)

	want := []MemberScore{}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestSkipList_Range_InvalidScoreRange(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")

	got := sl.Range(20, 10, 0, -1)

	want := []MemberScore{}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}
func TestSkipList_Range_SameScore(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "charlie")
	sl.Insert(10, "alpha")
	sl.Insert(10, "bravo")

	got := sl.Range(10, 10, 0, -1)

	want := []MemberScore{
		{Member: "alpha", Score: 10},
		{Member: "bravo", Score: 10},
		{Member: "charlie", Score: 10},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestSkipList_Rank(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")
	sl.Insert(30, "baz")

	tests := []struct {
		name   string
		score  float64
		member string
		want   int
	}{
		{
			name:   "first",
			score:  10,
			member: "foo",
			want:   0,
		},
		{
			name:   "middle",
			score:  20,
			member: "bar",
			want:   1,
		},
		{
			name:   "last",
			score:  30,
			member: "baz",
			want:   2,
		},
		{
			name:   "missing",
			score:  40,
			member: "missing",
			want:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sl.Rank(tt.score, tt.member)

			if got != tt.want {
				t.Fatalf("Rank = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSkipList_Rank_SameScore(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "charlie")
	sl.Insert(10, "alpha")
	sl.Insert(10, "bravo")

	tests := []struct {
		member string
		want   int
	}{
		{member: "alpha", want: 0},
		{member: "bravo", want: 1},
		{member: "charlie", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.member, func(t *testing.T) {
			got := sl.Rank(10, tt.member)

			if got != tt.want {
				t.Fatalf("Rank = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSkipList_Delete_LastNode_ShrinksLevel(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")

	// We cannot guarantee the inserted node's random level,
	// so this test verifies the invariant after deletion.
	if !sl.Delete(10, "foo") {
		t.Fatal("expected Delete to succeed")
	}

	if sl.level != 1 {
		t.Fatalf("level = %d, want 1", sl.level)
	}

	if sl.size != 0 {
		t.Fatalf("size = %d, want 0", sl.size)
	}
}

func TestSkipList_Insert_IncreasesSize(t *testing.T) {
	sl := NewSkipList()

	for i := 0; i < 100; i++ {
		sl.Insert(float64(i), "member")
	}

	if sl.size != 100 {
		t.Fatalf("size = %d, want 100", sl.size)
	}
}

func TestSkipList_Delete_DecreasesSize(t *testing.T) {
	sl := NewSkipList()

	sl.Insert(10, "foo")
	sl.Insert(20, "bar")
	sl.Insert(30, "baz")

	sl.Delete(20, "bar")

	if sl.size != 2 {
		t.Fatalf("size = %d, want 2", sl.size)
	}
}
