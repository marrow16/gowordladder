package words

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var expectedDictionarySizes = map[int]int{
	2:  127,
	3:  1347,
	4:  5638,
	5:  12972,
	6:  23033,
	7:  34342,
	8:  42150,
	9:  42933,
	10: 37235,
	11: 29027,
	12: 21025,
	13: 14345,
	14: 9397,
	15: 5925,
}

func TestCanLoadDictionariesFromFactory(t *testing.T) {
	for k, v := range expectedDictionarySizes {
		dictionary := LoadDictionary(k)
		assert.Equal(t, v, dictionary.Len())
		assert.Equal(t, k, dictionary.WordLength())
	}
}

func TestCanLoadDictionariesFromConstructor(t *testing.T) {
	for k, v := range expectedDictionarySizes {
		dictionary := NewDictionary(k)
		assert.Equal(t, v, dictionary.Len())
	}
}

func TestDictionaryFromFactorySameAsConstructed(t *testing.T) {
	newDict := NewDictionary(3)
	dictFromFactory := LoadDictionary(3)
	assert.Equal(t, newDict, dictFromFactory)
}

func TestFailsToLoadInvalidWordLengths(t *testing.T) {
	assert.Panics(t, func() {
		LoadDictionary(1)
	})
	assert.Panics(t, func() {
		LoadDictionary(16)
	})
}

func TestDictionaryWordHasVariants(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("cat")
	assert.True(t, ok)
	assert.Equal(t, 33, len(word.LinkedWords()))
	assert.False(t, word.IsIsland())
}

func TestDictionaryWordIsIslandWord(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("iwi")
	assert.True(t, ok)
	assert.True(t, word.IsIsland())
	assert.Equal(t, 0, len(word.LinkedWords()))
}

func TestDifferencesBetweenLinkedWords(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("cat")
	assert.True(t, ok)
	assert.True(t, len(word.LinkedWords()) > 0)
	for _, linkedWord := range word.LinkedWords() {
		assert.Equal(t, 1, word.Differences(linkedWord))
	}
}

func TestWordsAreInterlinked(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("cat")
	assert.True(t, ok)
	assert.True(t, len(word.LinkedWords()) > 0)
	for _, linkedWord := range word.LinkedWords() {
		assert.True(t, contains(linkedWord.LinkedWords(), word))
	}
}

func contains(s []Word, e Word) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
