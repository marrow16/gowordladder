package words

import (
	"fmt"
	"strings"
)

const reservedPatternChar = '_'

type Word interface {
	addLink(otherWord Word)
	LinkedWords() []Word
	Variations() []string
	IsIsland() bool
	Differences(other Word) int
	ActualWord() string
}

type word struct {
	actualWord  string
	linkedWords []Word
}

func newWord(actualWord string) Word {
	if strings.Contains(actualWord, string(reservedPatternChar)) {
		panic(any(fmt.Sprintf("Word cannot contain reserved character - '%v'", reservedPatternChar)))
	}
	return &word{
		actualWord:  strings.ToUpper(actualWord),
		linkedWords: make([]Word, 0),
	}
}

func (w *word) LinkedWords() []Word {
	return w.linkedWords
}

func (w *word) Variations() []string {
	result := make([]string, len(w.actualWord))
	for i := range w.actualWord {
		patt := []rune(w.actualWord)
		patt[i] = reservedPatternChar
		result[i] = string(patt)
	}
	return result
}

func (w *word) addLink(otherWord Word) {
	w.linkedWords = append(w.linkedWords, otherWord)
}

func (w *word) IsIsland() bool {
	return len(w.linkedWords) == 0
}

func (w *word) Differences(other Word) int {
	diffs := 0
	for i := range w.actualWord {
		if w.actualWord[i] != other.ActualWord()[i] {
			diffs++
		}
	}
	return diffs
}

func (w *word) ActualWord() string {
	return w.actualWord
}
