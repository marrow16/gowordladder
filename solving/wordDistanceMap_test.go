package solving

import (
	"gowordladder/test"
	"gowordladder/words"
	"testing"
)

func TestIslandWordHasLimitedMap(t *testing.T) {
	var dictionary = words.LoadDictionary(3)
	var word, ok = dictionary.Word("iwi")
	test.AssertTrue(t, ok)
	test.AssertTrue(t, word.IsIsland())

	var wordDistMap = NewWordDistanceMap(word)
	test.AssertEqualsInt(t, 1, len(wordDistMap.distances))
	var dist, ok2 = wordDistMap.distances[word.ActualWord()]
	test.AssertTrue(t, ok2)
	test.AssertEqualsInt(t, 1, dist)
}

func TestCatMap(t *testing.T) {
	var dictionary = words.LoadDictionary(3)
	var word, ok = dictionary.Word("cat")
	test.AssertTrue(t, ok)

	var wordDistMap = NewWordDistanceMap(word)
	test.AssertEqualsInt(t, 1346, len(wordDistMap.distances))
	var dist, ok2 = wordDistMap.distances[word.ActualWord()]
	test.AssertTrue(t, ok2)
	test.AssertEqualsInt(t, 1, dist)

	var endWord, ok3 = dictionary.Word("dog")
	test.AssertTrue(t, ok3)
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 5))
	test.AssertTrue(t, wordDistMap.Reachable(endWord, 4))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 3))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 2))
	test.AssertFalse(t, wordDistMap.Reachable(endWord, 1))
}
