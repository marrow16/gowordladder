package words

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanCreateWord(t *testing.T) {
	w := newWord("cat")
	assert.Equal(t, "CAT", w.ActualWord())
}

func TestFailsWithInvalidPatternChar(t *testing.T) {
	assert.Panics(t, func() {
		newWord("c" + string(reservedPatternChar) + "t")
	})
}

func TestVariationPatternsCorrect(t *testing.T) {
	w := newWord("cat")
	patts := w.Variations()
	assert.Equal(t, 3, len(patts))
	assert.Equal(t, "_AT", patts[0])
	assert.Equal(t, "C_T", patts[1])
	assert.Equal(t, "CA_", patts[2])
}

func TestDifferencesAreCorrect(t *testing.T) {
	cat := newWord("cat")
	cot := newWord("cot")
	dog := newWord("dog")
	assert.Equal(t, 0, cat.Differences(cat))
	assert.Equal(t, 0, cot.Differences(cot))
	assert.Equal(t, 0, dog.Differences(dog))
	assert.Equal(t, 1, cat.Differences(cot))
	assert.Equal(t, 1, cot.Differences(cat))
	assert.Equal(t, 2, cot.Differences(dog))
	assert.Equal(t, 2, dog.Differences(cot))
	assert.Equal(t, 3, cat.Differences(dog))
	assert.Equal(t, 3, dog.Differences(cat))
}

func TestWord_LinkedWords(t *testing.T) {
	w := newWord("xxx")
	assert.Equal(t, 0, len(w.LinkedWords()))

	w.addLink(newWord("yyy"))
	assert.Equal(t, 1, len(w.LinkedWords()))
}

func TestWord_IsIsland(t *testing.T) {
	w := newWord("xxx")
	assert.True(t, w.IsIsland())

	w.addLink(newWord("yyy"))
	assert.False(t, w.IsIsland())
}
