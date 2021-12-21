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

	wordDistMap := NewWordDistanceMap(word, nil)
	test.AssertEqualsInt(t, 1, len(wordDistMap.distances))
	dist, ok := wordDistMap.distances[word.ActualWord()]
	test.AssertTrue(t, ok)
	test.AssertEqualsInt(t, 1, dist)
}

func TestCatMap(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("cat")
	test.AssertTrue(t, ok)

	wordDistMap := NewWordDistanceMap(word, nil)
	test.AssertEqualsInt(t, 1346, len(wordDistMap.distances))
	dist, ok := wordDistMap.distances[word.ActualWord()]
	test.AssertTrue(t, ok)
	test.AssertEqualsInt(t, 1, dist)

	endWord, ok := dictionary.Word("dog")
	test.AssertTrue(t, ok)
	_, hasWord := wordDistMap.distances[endWord.ActualWord()]
	test.AssertTrue(t, hasWord)
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 5))
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 4))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 3))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 2))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 1))
}

func TestCatMapLimited(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("cat")
	test.AssertTrue(t, ok)

	limit := 4
	wordDistMap := NewWordDistanceMap(word, &limit)
	test.AssertEqualsInt(t, 1086, len(wordDistMap.distances))
	endWord, _ := dictionary.Word("dog")
	_, hasWord := wordDistMap.distances[endWord.ActualWord()]
	test.AssertTrue(t, hasWord)
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 5))
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 4))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 3))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 2))

	// limit further...
	limit = 3
	wordDistMap = NewWordDistanceMap(word, &limit)
	test.AssertEqualsInt(t, 345, len(wordDistMap.distances))
	_, hasWord = wordDistMap.distances[endWord.ActualWord()]
	test.AssertFalse(t, hasWord)
}
