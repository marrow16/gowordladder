package solving

import (
	"gowordladder/test"
	"gowordladder/words"
	"testing"
)

func TestIslandWordHasLimitedMap(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("iwi")
	test.AssertTrue(t, ok)
	test.AssertTrue(t, word.IsIsland())

	wordDistMap := NewWordDistanceMap(word)
	test.AssertEqualsInt(t, 1, len(wordDistMap.distances))
	dist, ok := wordDistMap.distances[word.ActualWord()]
	test.AssertTrue(t, ok)
	test.AssertEqualsInt(t, 1, dist)
}

func TestCatMap(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("cat")
	test.AssertTrue(t, ok)

	wordDistMap := NewWordDistanceMap(word)
	test.AssertEqualsInt(t, 1346, len(wordDistMap.distances))
	dist, ok := wordDistMap.distances[word.ActualWord()]
	test.AssertTrue(t, ok)
	test.AssertEqualsInt(t, 1, dist)

	endWord, ok := dictionary.Word("dog")
	test.AssertTrue(t, ok)
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 5))
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 4))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 3))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 2))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 1))
}
