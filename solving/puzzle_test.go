package solving

import (
	"github.com/stretchr/testify/assert"
	"gowordladder/words"
	"testing"
)

func TestCalculateMinimumLadderLength(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("DOG")

	puzzle := NewPuzzle(startWord, endWord)
	min, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 4, min)
}

func TestCalculateMinimumLadderLengthOneLetterDifference(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("COT")

	puzzle := NewPuzzle(startWord, endWord)
	min, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 2, min)
}

func TestCalculateMinimumLadderLengthTwoLetterDifference(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("COP")

	puzzle := NewPuzzle(startWord, endWord)
	min, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 3, min)
}

func TestCalculateMinimumLadderLengthSameWord(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("CAT")

	puzzle := NewPuzzle(startWord, endWord)
	min, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 1, min)
}

func TestCalculateMinimumLadderLengthFlips(t *testing.T) {
	dictionary := words.NewDictionary(3)
	startWord, _ := dictionary.Word("CAT")
	endWord, _ := dictionary.Word("ANI")

	puzzle := NewPuzzle(startWord, endWord)
	min, ok := puzzle.CalculateMinimumLadderLength()
	assert.True(t, ok)
	assert.Equal(t, 5, min)
}
