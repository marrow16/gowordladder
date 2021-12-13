package solving

import (
	"github.com/gammazero/deque"
	_ "github.com/gammazero/deque"
	"gowordladder/words"
)

type WordDistanceMap struct {
	distances map[string]int
}

func NewWordDistanceMap(word *words.Word) (result *WordDistanceMap) {
	result = &WordDistanceMap{map[string]int{}}
	result.distances[word.ActualWord()] = 1;

	var q deque.Deque
	q.PushBack(&*word)
	for q.Len() != 0 {
		var nextWord = q.PopFront().(*words.Word)
		var distance = result.distanceGetOrDefault(nextWord) + 1
		for _, linkedWord := range *nextWord.LinkedWords {
			if !result.Contains(linkedWord) {
				q.PushBack(&*linkedWord)
				result.computeIfAbsent(linkedWord, distance)
			}
		}
	}
	return
}

func (m *WordDistanceMap) distanceGetOrDefault(word *words.Word) (result int) {
	if d, ok := m.distances[word.ActualWord()]; !ok {
		result = 0
	} else {
		result = d
	}
	return
}

func (m *WordDistanceMap) computeIfAbsent(word *words.Word, distance int) {
	if _, ok := m.distances[word.ActualWord()]; !ok {
		m.distances[word.ActualWord()] = distance
	}
}

func (m *WordDistanceMap) Len() int {
	return len(m.distances)
}

func (m *WordDistanceMap) Contains(word *words.Word) bool {
	_, ok := m.distances[word.ActualWord()]
	return ok
}

func (m *WordDistanceMap) Distance(word *words.Word) (dist int, ok bool) {
	dist, ok = m.distances[word.ActualWord()]
	return
}

func (m *WordDistanceMap) Reachable(word *words.Word, maximumLadderLength int) bool {
	if distance, ok := m.distances[word.ActualWord()]; ok {
		return distance <= maximumLadderLength
	}
	return false
}
