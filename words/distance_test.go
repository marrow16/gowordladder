package words

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordDistanceMap_IslandWord(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("iwi")
	assert.True(t, ok)
	assert.True(t, w.IsIsland())

	wordDistMap := NewWordDistanceMap(w, nil)
	assert.Equal(t, 1, len(wordDistMap))
	dist, ok := wordDistMap.Distance(w)
	assert.True(t, ok)
	assert.Equal(t, 1, dist)
}

func TestWordDistanceMap(t *testing.T) {
	d := NewDictionary(3)
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
	notReachable, _ := d.Word("iwi")
	assert.False(t, wordDistMap.Reachable(notReachable, 100))
}

func TestWordDistanceMap_Limited(t *testing.T) {
	d := NewDictionary(3)
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

func TestWordDistanceMap_Words(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)

	limit := 3
	wordDistMap := NewWordDistanceMap(w, &limit)
	assert.Len(t, wordDistMap.Words(), 344)
}

func TestWordDistanceMap_WordsAt(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)

	wordDistMap := NewWordDistanceMap(w, nil)
	assert.Len(t, wordDistMap.WordsAt(2), 33)
	assert.Len(t, wordDistMap.WordsAt(w.MaxSteps()), 1)
}

func TestWordDistanceMap_MaxDistance(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)

	wordDistMap := NewWordDistanceMap(w, nil)
	assert.Equal(t, w.MaxSteps(), wordDistMap.MaxDistance())
}
