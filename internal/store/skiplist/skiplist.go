package skiplist

import (
	"math/rand/v2"
)

const maxLevel = 32
const probability = 0.25

type SkipList struct {
	head  *skipNode
	level int
	size  int
}
type skipNode struct {
	member  string
	score   float64
	forward [maxLevel]*skipNode
}
type MemberScore struct {
	Member string
	Score  float64
}

func NewSkipList() *SkipList {
	sl := &SkipList{}
	sl.head = &skipNode{}
	sl.level = 1
	return sl
}

func (sl *SkipList) Insert(score float64, member string) {
	update := sl.findPath(score, member)
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.head
		}
		sl.level = level
	}
	node := &skipNode{
		member: member,
		score:  score,
	}
	for i := range level {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}
	sl.size++
}
func (sl *SkipList) Delete(score float64, member string) bool {
	path := sl.findPath(score, member)
	target := path[0].forward[0]

	if target == nil || target.score != score || target.member != member {
		return false
	}
	for level := 0; level < sl.level; level++ {
		if path[level].forward[level] == target {
			path[level].forward[level] = target.forward[level]
		}
	}
	sl.size--
	sl.shrink()
	return true
}
func (sl *SkipList) Search(score float64, member string) *skipNode {
	path := sl.findPath(score, member)
	next := path[0].forward[0]

	if next != nil && next.score == score && next.member == member {
		return next
	}
	return nil
}
func (sl *SkipList) Range(minScore, maxScore float64, offset, count int) []MemberScore {
	if offset < 0 || count == 0 || minScore > maxScore {
		return []MemberScore{}
	}
	myNode := sl.findFirstAtOrAfter(minScore)
	index := 0
	results := make([]MemberScore, 0)

	for myNode != nil && myNode.score <= maxScore {
		if index >= offset {
			results = append(results, MemberScore{
				Member: myNode.member,
				Score:  myNode.score,
			})

			if count > 0 && len(results) == count {
				break
			}
		}
		index++
		myNode = myNode.forward[0]
	}
	return results
}
func (sl *SkipList) RangeByRank(start, end int) []MemberScore {
	if start > end {
		return []MemberScore{}
	}
	index := 0
	results := make([]MemberScore, 0)
	node := sl.head.forward[0]
	for node != nil && index <= end {
		if index >= start {
			results = append(results, MemberScore{
				Member: node.member,
				Score:  node.score,
			})
		}
		index++
		node = node.forward[0]
	}
	return results

}
func (sl *SkipList) Rank(score float64, member string) int {
	myNode := sl.head.forward[0]
	rank := 0

	for myNode != nil {
		if myNode.score == score && myNode.member == member {
			return rank
		}
		myNode = myNode.forward[0]
		rank++
	}
	return -1
}

/*
 *
 *
 * helpers
 *
 *
 *
 */
func before(score float64, member string, node *skipNode) bool {
	return node.score < score ||
		(node.score == score && node.member < member)
}
func randomLevel() int {
	level := 1
	for level < maxLevel && rand.Float64() < probability {
		level++
	}
	return level
}
func (sl *SkipList) findPath(score float64, member string) []*skipNode {
	update := make([]*skipNode, maxLevel)
	prev := sl.head

	for level := sl.level - 1; level >= 0; level-- {
		for prev.forward[level] != nil &&
			before(score, member, prev.forward[level]) {
			prev = prev.forward[level]
		}
		update[level] = prev
	}
	return update
}

func (sl *SkipList) shrink() {
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
}

func (sl *SkipList) findFirstAtOrAfter(minScore float64) *skipNode {
	prev := sl.head
	for level := sl.level - 1; level >= 0; level-- {
		for prev.forward[level] != nil && prev.forward[level].score < minScore {
			prev = prev.forward[level]
		}
	}
	return prev.forward[0]
}
