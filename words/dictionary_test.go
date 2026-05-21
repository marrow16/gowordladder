package words

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var expectedDictionarySizes = map[int]int{
	2:  127,
	3:  1285,
	4:  5191,
	5:  11798,
	6:  20952,
	7:  31498,
	8:  39095,
	9:  39967,
	10: 34327,
	11: 26061,
	12: 18679,
	13: 12541,
	14: 7908,
	15: 4773,
}
var expectedMaxSteps = map[int]int{
	2:  5,
	3:  9,
	4:  17,
	5:  30,
	6:  50,
	7:  53,
	8:  61,
	9:  32,
	10: 11,
	11: 7,
	12: 7,
	13: 5,
	14: 6,
	15: 4,
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

func TestSwitchCurrentDictionary(t *testing.T) {
	defer SwitchCurrentDictionary(defaultDictionary)

	d := NewDictionary(3)
	assert.Equal(t, 1285, d.Len())

	err := SwitchCurrentDictionary("enwiktionary")
	require.NoError(t, err)
	d = NewDictionary(3)
	assert.Equal(t, 2573, d.Len())

	err = SwitchCurrentDictionary("./resources/csw19-")
	require.NoError(t, err)
	d = NewDictionary(3)
	assert.Equal(t, 1347, d.Len())
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
	assert.Equal(t, 32, len(word.LinkedWords()))
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
