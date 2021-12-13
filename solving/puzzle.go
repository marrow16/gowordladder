package solving

import "gowordladder/words"

type Puzzle struct {
	startWord *words.Word
	endWord   *words.Word
}

func NewPuzzle(startWord *words.Word, endWord *words.Word) *Puzzle {
	return &Puzzle{
		startWord: startWord,
		endWord:   endWord,
	}
}

func (p *Puzzle) CalculateMinimumLadderLength() (min int, ok bool) {
	var start = p.startWord
	var end = p.endWord
	var diffs = start.Differences(end)
	switch diffs {
	case 0, 1:
		return diffs + 1, true
	case 2: {
		var startSet = make(map[string]*words.Word, len(*start.LinkedWords))
		for _, w := range *start.LinkedWords {
			startSet[w.ActualWord()] = w
		}
		for _, w := range *end.LinkedWords {
			if _, ok := startSet[w.ActualWord()]; ok {
				return 3, true
			}
		}
	}
	}
	if len(*start.LinkedWords) > len(*end.LinkedWords) {
		start = p.endWord
		end = p.startWord
	}
	return NewWordDistanceMap(start).Distance(end)
}
