package words

import (
	"github.com/gammazero/deque"
)

type WordDistanceMap map[string]int

func NewWordDistanceMap(word Word, maximumLadderLength *int) WordDistanceMap {
	result := make(WordDistanceMap)
	result[word.ActualWord()] = 1

	var maxDistance = 100
	if maximumLadderLength != nil {
		maxDistance = *maximumLadderLength
	}
	var q deque.Deque
	q.PushBack(word)
	for q.Len() != 0 {
		nextWord := q.PopFront().(Word)
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

func (m WordDistanceMap) distanceGetOrDefault(word Word) int {
	result := 0
	if d, ok := m[word.ActualWord()]; ok {
		result = d
	}
	return result
}

func (m WordDistanceMap) computeIfAbsent(word Word, distance int) {
	if _, ok := m[word.ActualWord()]; !ok {
		m[word.ActualWord()] = distance
	}
}

func (m WordDistanceMap) Contains(word Word) bool {
	_, ok := m[word.ActualWord()]
	return ok
}

func (m WordDistanceMap) Distance(word Word) (dist int, ok bool) {
	dist, ok = m[word.ActualWord()]
	return
}

func (m WordDistanceMap) Reachable(word Word, maximumLadderLength int) bool {
	if distance, ok := m[word.ActualWord()]; ok {
		return distance <= maximumLadderLength
	}
	return false
}

func (m WordDistanceMap) WordsAt(ladderLength int) (result []string) {
	for k, v := range m {
		if v == ladderLength {
			result = append(result, k)
		}
	}
	return result
}

func (m WordDistanceMap) Words() (result []string) {
	for k, v := range m {
		if v > 1 {
			result = append(result, k)
		}
	}
	return result
}

func (m WordDistanceMap) MaxDistance() (result int) {
	for _, v := range m {
		if v > result {
			result = v
		}
	}
	return result
}
