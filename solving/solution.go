package solving

import (
	"gowordladder/words"
	"sort"
	"strings"
)

type Solution struct {
	ladder []*words.Word
}

func (s *Solution) ToString() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, w := range s.ladder {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(w.ActualWord())
	}
	sb.WriteString("]")
	return sb.String()
}

func newSolution(words ...*words.Word) *Solution {
	return &Solution{ladder: words}
}

type candidateSolution struct {
	solver *Solver
	ladder *[]*words.Word
	seenWords map[string]bool
}

func newCandidateSolution(solver *Solver, startWord *words.Word, nextWord *words.Word) (result *candidateSolution) {
	result = &candidateSolution{
		solver:    solver,
		ladder:    &[]*words.Word{},
		seenWords: map[string]bool{},
	}
	result.addWord(startWord)
	result.addWord(nextWord)
	solver.incrementExplored()
	return
}

func (s *candidateSolution) spawn(nextWord *words.Word) (result *candidateSolution) {
	result = &candidateSolution{
		solver:    s.solver,
		ladder:    &[]*words.Word{},
		seenWords: map[string]bool{},
	}
	for sw := range s.seenWords {
		result.seenWords[sw] = true
	}
	for _, w := range *s.ladder {
		result.addWord(w)
	}
	result.addWord(nextWord)
	s.solver.incrementExplored()
	return
}

func (s *candidateSolution) lastWord() *words.Word {
	return (*s.ladder)[len(*s.ladder) - 1]
}

func (s *candidateSolution) seen(word *words.Word) bool {
	if _, ok := s.seenWords[word.ActualWord()]; ok {
		return true
	}
	return false
}

func (s *candidateSolution) asFoundSolution(reversed bool) (result *Solution) {
	result = &Solution{ladder: make([]*words.Word, len(*s.ladder))}
	if reversed {
		var l = len(*s.ladder) - 1
		for i, w := range *s.ladder {
			result.ladder[l - i] = w
		}
	} else {
		for i, w := range *s.ladder {
			result.ladder[i] = w
		}
	}
	return
}

func (s *candidateSolution) addWord(word *words.Word) {
	*s.ladder = append(*s.ladder, word)
	s.seenWords[word.ActualWord()] = true
}

func SortSolutions(solutions []*Solution) {
	sort.Slice(solutions, func(i, j int) bool {
		if len(solutions[i].ladder) < len(solutions[j].ladder) {
			return true
		} else if len(solutions[i].ladder) != len(solutions[j].ladder) {
			return false
		}
		for idx, w := range solutions[i].ladder {
			if w.ActualWord() < solutions[j].ladder[idx].ActualWord() {
				return true
			}
		}
		return false
	})
}
