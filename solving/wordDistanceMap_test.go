package solving

import (
	"github.com/stretchr/testify/assert"
	"gowordladder/words"
	"testing"
)

func TestIslandWordHasLimitedMap(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("iwi")
	assert.True(t, ok)
	assert.True(t, word.IsIsland())

	wordDistMap := NewWordDistanceMap(word, nil)
	assert.Equal(t, 1, wordDistMap.Len())
	dist, ok := wordDistMap.Distance(word)
	assert.True(t, ok)
	assert.Equal(t, 1, dist)
}

func TestCatMap(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("cat")
	assert.True(t, ok)

	wordDistMap := NewWordDistanceMap(word, nil)
	assert.Equal(t, 1346, wordDistMap.Len())
	dist, ok := wordDistMap.Distance(word)
	assert.True(t, ok)
	assert.Equal(t, 1, dist)

	endWord, ok := dictionary.Word("dog")
	assert.True(t, ok)
	_, hasWord := wordDistMap.Distance(endWord)
	assert.True(t, hasWord)
	assert.True(t, wordDistMap.Reachable(endWord, 5))
	assert.True(t, wordDistMap.Reachable(endWord, 4))
	assert.False(t, wordDistMap.Reachable(endWord, 3))
	assert.False(t, wordDistMap.Reachable(endWord, 2))
	assert.False(t, wordDistMap.Reachable(endWord, 1))
}

func TestCatMapLimited(t *testing.T) {
	dictionary := words.LoadDictionary(3)
	word, ok := dictionary.Word("cat")
	assert.True(t, ok)

	limit := 4
	wordDistMap := NewWordDistanceMap(word, &limit)
	assert.Equal(t, 1086, wordDistMap.Len())
	endWord, _ := dictionary.Word("dog")
	_, hasWord := wordDistMap.Distance(endWord)
	assert.True(t, hasWord)
	assert.True(t, wordDistMap.Reachable(endWord, 5))
	assert.True(t, wordDistMap.Reachable(endWord, 4))
	assert.False(t, wordDistMap.Reachable(endWord, 3))
	assert.False(t, wordDistMap.Reachable(endWord, 2))

	// limit further...
	limit = 3
	wordDistMap = NewWordDistanceMap(word, &limit)
	assert.Equal(t, 345, wordDistMap.Len())
	_, hasWord = wordDistMap.Distance(endWord)
	assert.False(t, hasWord)
}
