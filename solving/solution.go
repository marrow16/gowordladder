package solving

import (
	"gowordladder/words"
	"strings"
)

type Solution interface {
	Ladder() []words.Word
	String() string
}

type solution struct {
	ladder []words.Word
}

func (s *solution) Ladder() []words.Word {
	return s.ladder
}

func (s *solution) String() string {
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

func newSolution(words ...words.Word) Solution {
	return &solution{ladder: words}
}

type candidateSolution struct {
	solver    incrementable
	ladder    []words.Word
	seenWords map[string]bool
}

func newCandidateSolution(solver incrementable, startWord words.Word, nextWord words.Word) (result *candidateSolution) {
	result = &candidateSolution{
		solver:    solver,
		ladder:    make([]words.Word, 0),
		seenWords: map[string]bool{},
	}
	result.addWord(startWord)
	result.addWord(nextWord)
	solver.incrementExplored()
	return
}

func (s *candidateSolution) spawn(nextWord words.Word) (result *candidateSolution) {
	result = &candidateSolution{
		solver:    s.solver,
		ladder:    make([]words.Word, len(s.ladder)),
		seenWords: map[string]bool{},
	}
	for sw := range s.seenWords {
		result.seenWords[sw] = true
	}
	for i, w := range s.ladder {
		result.ladder[i] = w
	}
	result.addWord(nextWord)
	s.solver.incrementExplored()
	return
}

func (s *candidateSolution) lastWord() words.Word {
	return (s.ladder)[len(s.ladder)-1]
}

func (s *candidateSolution) seen(word words.Word) bool {
	return s.seenWords[word.ActualWord()]
}

func (s *candidateSolution) asFoundSolution(reversed bool) Solution {
	result := &solution{ladder: make([]words.Word, len(s.ladder))}
	if reversed {
		l := len(s.ladder) - 1
		for i, w := range s.ladder {
			result.ladder[l-i] = w
		}
	} else {
		copy(result.ladder, s.ladder)
	}
	return result
}

func (s *candidateSolution) addWord(word words.Word) {
	s.ladder = append(s.ladder, word)
	s.seenWords[word.ActualWord()] = true
}
