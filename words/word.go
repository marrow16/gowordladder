package words

import (
	"strings"
)

const reservedPatternChar = '_'

type Word struct {
	actualWord  string
	LinkedWords *[]*Word
}

func newWord(actualWord string) *Word {
	if strings.Contains(actualWord, string(reservedPatternChar)) {
		panic("Word cannot contain reserved character ('" + string(reservedPatternChar) + "')")
	}
	return &Word{actualWord: strings.ToUpper(actualWord), LinkedWords: &[]*Word{}}
}

func (w *Word) variations() (result []string) {
	result = []string{}
	for i := range w.actualWord {
		patt := []rune(w.actualWord)
		patt[i] = '_'
		result = append(result, string(patt))
	}
	return
}

func (w *Word) addLink(otherWord *Word) {
	*w.LinkedWords = append(*w.LinkedWords, otherWord)
}

func (w *Word) IsIsland() bool {
	return len(*w.LinkedWords) == 0
}

func (w *Word) Differences(other *Word) int {
	diffs := 0
	for i := range w.actualWord {
		if w.actualWord[i] != other.actualWord[i] {
			diffs++
		}
	}
	return diffs
}

func (w *Word) ToString() string {
	return w.actualWord
}

func (w *Word) ActualWord() string {
	return w.actualWord
}
