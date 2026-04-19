package words

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIslandWordHasLimitedMap(t *testing.T) {
	d := LoadDictionary(3)
	w, ok := d.Word("iwi")
	assert.True(t, ok)
	assert.True(t, w.IsIsland())

	wordDistMap := NewWordDistanceMap(w, nil)
	assert.Equal(t, 1, len(wordDistMap))
	dist, ok := wordDistMap.Distance(w)
	assert.True(t, ok)
	assert.Equal(t, 1, dist)
}

func TestCatMap(t *testing.T) {
	d := LoadDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)

	wordDistMap := NewWordDistanceMap(w, nil)
	assert.Equal(t, 1346, len(wordDistMap))
	dist, ok := wordDistMap.Distance(w)
	assert.True(t, ok)
	assert.Equal(t, 1, dist)

	endWord, ok := d.Word("dog")
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
	d := LoadDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)

	limit := 4
	wordDistMap := NewWordDistanceMap(w, &limit)
	assert.Equal(t, 1086, len(wordDistMap))
	endWord, _ := d.Word("dog")
	_, hasWord := wordDistMap.Distance(endWord)
	assert.True(t, hasWord)
	assert.True(t, wordDistMap.Reachable(endWord, 5))
	assert.True(t, wordDistMap.Reachable(endWord, 4))
	assert.False(t, wordDistMap.Reachable(endWord, 3))
	assert.False(t, wordDistMap.Reachable(endWord, 2))

	// limit further...
	limit = 3
	wordDistMap = NewWordDistanceMap(w, &limit)
	assert.Equal(t, 345, len(wordDistMap))
	_, hasWord = wordDistMap.Distance(endWord)
	assert.False(t, hasWord)
}
