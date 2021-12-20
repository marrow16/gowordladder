package words

import (
	"gowordladder/test"
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
		test.AssertEqualsInt(t, v, dictionary.Len())
	}
}

func TestCanLoadDictionariesFromConstructor(t *testing.T) {
	for k, v := range expectedDictionarySizes {
		dictionary := NewDictionary(k)
		test.AssertEqualsInt(t, v, dictionary.Len())
	}
}

func TestDictionaryFromFactorySameAsConstructed(t *testing.T) {
	newDict := NewDictionary(3)
	dictFromFactory := LoadDictionary(3)
	test.AssertTrue(t, newDict == dictFromFactory)
}

func TestFailsToLoadInvalidWordLengths(t *testing.T) {
	test.AssertPanic(t, func() {
		LoadDictionary(1)
	})
	test.AssertPanic(t, func() {
		LoadDictionary(16)
	})
}

func TestDictionaryWordHasVariants(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("cat")
	test.AssertTrue(t, ok)
	test.AssertEqualsInt(t, 33, len(*word.LinkedWords))
	test.AssertFalse(t, word.IsIsland())
}

func TestDictionaryWordIsIslandWord(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("iwi")
	test.AssertTrue(t, ok)
	test.AssertTrue(t, word.IsIsland())
	test.AssertEqualsInt(t, 0, len(*word.LinkedWords))
}

func TestDifferencesBetweenLinkedWords(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("cat")
	test.AssertTrue(t, ok)
	test.AssertTrue(t, len(*word.LinkedWords) > 0)
	for _, linkedWord := range *word.LinkedWords {
		test.AssertEqualsInt(t, 1, word.Differences(linkedWord))
	}
}

func TestWordsAreInterlinked(t *testing.T) {
	dictionary := NewDictionary(3)
	word, ok := dictionary.Word("cat")
	test.AssertTrue(t, ok)
	test.AssertTrue(t, len(*word.LinkedWords) > 0)
	for _, linkedWord := range *word.LinkedWords {
		test.AssertTrue(t, contains(*linkedWord.LinkedWords, word))
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
