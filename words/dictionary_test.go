package words

import (
	"fmt"
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
var expectedMaxSteps = map[int]int{
	2:  5,
	3:  9,
	4:  16,
	5:  27,
	6:  43,
	7:  65,
	8:  80,
	9:  34,
	10: 11,
	11: 27,
	12: 7,
	13: 5,
	14: 7,
	15: 5,
}

func TestCanLoadDictionariesFromFactory(t *testing.T) {
	for k, v := range expectedDictionarySizes {
		t.Run(fmt.Sprintf("%d-letters", k), func(t *testing.T) {
			d := NewDictionary(k)
			assert.Equal(t, v, d.Len())
			assert.Equal(t, k, d.WordLength())
			assert.Len(t, d.Words(), v)
			assert.Equal(t, expectedMaxSteps[k], d.MaxSteps())
			for i := 3; i <= expectedMaxSteps[k]; i++ {
				assert.True(t, len(d.WordsWithSteps(i)) > 0)
			}
			assert.Len(t, d.WordsWithSteps(expectedMaxSteps[k]+1), 0)
		})
	}
}

func TestCanLoadDictionariesFromConstructor(t *testing.T) {
	for k, v := range expectedDictionarySizes {
		d := NewDictionary(k)
		assert.Equal(t, v, d.Len())
	}
}

func TestDictionaryFromFactorySameAsConstructed(t *testing.T) {
	newDict := NewDictionary(3)
	dictFromFactory := NewDictionary(3)
	assert.Equal(t, newDict, dictFromFactory)
}

func TestFailsToLoadInvalidWordLengths(t *testing.T) {
	assert.Panics(t, func() {
		NewDictionary(1)
	})
	assert.Panics(t, func() {
		NewDictionary(16)
	})
}

func TestDictionaryWordHasVariants(t *testing.T) {
	d := NewDictionary(3)
	word, ok := d.Word("cat")
	assert.True(t, ok)
	assert.Equal(t, 33, len(word.LinkedWords()))
	assert.False(t, word.IsIsland())
}

func TestDictionaryWordIsIslandWord(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("iwi")
	assert.True(t, ok)
	assert.True(t, w.IsIsland())
	assert.Equal(t, 0, len(w.LinkedWords()))
}

func TestDifferencesBetweenLinkedWords(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)
	assert.True(t, len(w.LinkedWords()) > 0)
	for _, linkedWord := range w.LinkedWords() {
		assert.Equal(t, 1, w.Differences(linkedWord))
	}
}

func TestWordsAreInterlinked(t *testing.T) {
	d := NewDictionary(3)
	w, ok := d.Word("cat")
	assert.True(t, ok)
	assert.True(t, len(w.LinkedWords()) > 0)
	for _, linkedWord := range w.LinkedWords() {
		assert.True(t, contains(linkedWord.LinkedWords(), w))
	}
}

func contains(s []*Word, e *Word) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
