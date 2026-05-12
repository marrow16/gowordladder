package words

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type DictionaryImporter struct {
	WordLen    int
	Words      []*Word
	seen       map[string]struct{}
	variations variations
}

func NewDictionaryImporter(wordLen int) *DictionaryImporter {
	return &DictionaryImporter{
		WordLen:    wordLen,
		Words:      make([]*Word, 0, 100_000),
		seen:       make(map[string]struct{}),
		variations: variations{},
	}
}

func (di *DictionaryImporter) Import(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if w := strings.ToUpper(scanner.Text()); len(w) == di.WordLen && isAllAZ(w) {
			if _, seen := di.seen[w]; !seen {
				di.seen[w] = struct{}{}
				word := newWord(w, 0)
				di.Words = append(di.Words, word)
				di.variations.link(word)
			}
		} else {
			return fmt.Errorf("invalid word: %s", w)
		}
	}
	return scanner.Err()
}

func (di *DictionaryImporter) AddWord(word string) error {
	if w := strings.ToUpper(word); len(w) == di.WordLen && isAllAZ(w) {
		if _, seen := di.seen[w]; !seen {
			di.seen[w] = struct{}{}
			wd := newWord(w, 0)
			di.Words = append(di.Words, wd)
			di.variations.link(wd)
		}
	} else {
		return fmt.Errorf("invalid word: %s", word)
	}
	return nil
}

func (di *DictionaryImporter) Process(out func(word *Word) error) error {
	maxwdl := 100
	for _, word := range di.Words {
		wdm := NewWordDistanceMap(word, &maxwdl)
		word.maxSteps = wdm.MaxDistance()
		if err := out(word); err != nil {
			return err
		}
	}
	return nil
}

func isAllAZ(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}
