package words

import (
	"fmt"
	"slices"
	"strings"
)

const reservedPatternChar = '_'

type Word struct {
	actualWord  string
	linkedWords []*Word
	maxSteps    int
}

func newWord(actualWord string, maxSteps int) *Word {
	if strings.ContainsRune(actualWord, reservedPatternChar) {
		panic(any(fmt.Sprintf("Word cannot contain reserved character - '%v'", reservedPatternChar)))
	}
	return &Word{
		actualWord:  strings.ToUpper(actualWord),
		linkedWords: make([]*Word, 0),
		maxSteps:    maxSteps,
	}
}

func (w *Word) LinkedWords() []*Word {
	return slices.Clone(w.linkedWords)
}

func (w *Word) Variations() []string {
	result := make([]string, len(w.actualWord))
	for i := range w.actualWord {
		patt := []rune(w.actualWord)
		patt[i] = reservedPatternChar
		result[i] = string(patt)
	}
	return result
}

func (w *Word) addLink(otherWord *Word) {
	w.linkedWords = append(w.linkedWords, otherWord)
}

func (w *Word) IsIsland() bool {
	return len(w.linkedWords) == 0
}

func (w *Word) IsDoublet() bool {
	return w.maxSteps == 2
}

func (w *Word) Differences(other *Word) int {
	diffs := 0
	for i := range w.actualWord {
		if w.actualWord[i] != other.ActualWord()[i] {
			diffs++
		}
	}
	return diffs
}

func (w *Word) ActualWord() string {
	return w.actualWord
}

func (w *Word) MaxSteps() int {
	return w.maxSteps
}
