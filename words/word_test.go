package words

import (
	"gowordladder/test"
	"testing"
)

func TestCanCreateWord(t *testing.T) {
	var word = newWord("cat")
	test.AssertEqualsString(t, "CAT", word.ActualWord())
}

func TestFailsWithInvalidPatternChar(t *testing.T) {
	test.AssertPanic(t, func() {
		newWord("c_t")
	})
}

func TestVariationPatternsCorrect(t *testing.T) {
	var word = newWord("cat")
	var patts = word.variations()
	test.AssertEqualsInt(t, 3, len(patts))
	test.AssertEqualsString(t, "_AT", patts[0])
	test.AssertEqualsString(t, "C_T", patts[1])
	test.AssertEqualsString(t, "CA_", patts[2])
}

func TestDifferencesAreCorrect(t *testing.T) {
	var cat = newWord("cat")
	var cot = newWord("cot")
	var dog = newWord("dog")
	test.AssertEqualsInt(t, 0, cat.Differences(cat))
	test.AssertEqualsInt(t, 0, cot.Differences(cot))
	test.AssertEqualsInt(t, 0, dog.Differences(dog))
	test.AssertEqualsInt(t, 1, cat.Differences(cot))
	test.AssertEqualsInt(t, 1, cot.Differences(cat))
	test.AssertEqualsInt(t, 2, cot.Differences(dog))
	test.AssertEqualsInt(t, 2, dog.Differences(cot))
	test.AssertEqualsInt(t, 3, cat.Differences(dog))
	test.AssertEqualsInt(t, 3, dog.Differences(cat))
}
