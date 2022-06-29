package solving

import (
	"github.com/gammazero/deque"
	"gowordladder/words"
)

type WordDistanceMap interface {
	Contains(word words.Word) bool
	Distance(word words.Word) (dist int, ok bool)
	Reachable(word words.Word, maximumLadderLength int) bool
	Len() int
}

type wordDistanceMap struct {
	distances map[string]int
}

func NewWordDistanceMap(word words.Word, maximumLadderLength *int) WordDistanceMap {
	result := &wordDistanceMap{map[string]int{}}
	result.distances[word.ActualWord()] = 1

	var maxDistance = 255
	if maximumLadderLength != nil {
		maxDistance = *maximumLadderLength
	}
	var q deque.Deque
	q.PushBack(word)
	for q.Len() != 0 {
		nextWord := q.PopFront().(words.Word)
		distance := result.distanceGetOrDefault(nextWord) + 1
		if distance <= maxDistance {
			for _, linkedWord := range nextWord.LinkedWords() {
				if !result.Contains(linkedWord) {
					q.PushBack(linkedWord)
					result.computeIfAbsent(linkedWord, distance)
				}
			}
		}
	}
	return result
}

func (m *wordDistanceMap) Len() int {
	return len(m.distances)
}

func (m *wordDistanceMap) distanceGetOrDefault(word words.Word) int {
	result := 0
	if d, ok := m.distances[word.ActualWord()]; ok {
		result = d
	}
	return result
}

func (m *wordDistanceMap) computeIfAbsent(word words.Word, distance int) {
	if _, ok := m.distances[word.ActualWord()]; !ok {
		m.distances[word.ActualWord()] = distance
	}
}

func (m *wordDistanceMap) Contains(word words.Word) bool {
	_, ok := m.distances[word.ActualWord()]
	return ok
}

func (m *wordDistanceMap) Distance(word words.Word) (dist int, ok bool) {
	dist, ok = m.distances[word.ActualWord()]
	return
}

func (m *wordDistanceMap) Reachable(word words.Word, maximumLadderLength int) bool {
	if distance, ok := m.distances[word.ActualWord()]; ok {
		return distance <= maximumLadderLength
	}
	return false
}
