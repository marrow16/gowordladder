package solving

import (
	"github.com/marrow16/gowordladder/words"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateMinimumLadderLength(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("DOG")

	puzzle := NewPuzzle(startWord, endWord)
	m, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 4, m)
}

func TestCalculateMinimumLadderLengthOneLetterDifference(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("COT")

	puzzle := NewPuzzle(startWord, endWord)
	m, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 2, m)
}

func TestCalculateMinimumLadderLengthTwoLetterDifference(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("COP")

	puzzle := NewPuzzle(startWord, endWord)
	m, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 3, m)
}

func TestCalculateMinimumLadderLengthSameWord(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("CAT")

	puzzle := NewPuzzle(startWord, endWord)
	m, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 1, m)
}

func TestCalculateMinimumLadderLengthFlips(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("ANI")

	puzzle := NewPuzzle(startWord, endWord)
	m, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 5, m)
}
