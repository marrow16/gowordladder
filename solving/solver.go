package solving

import (
	"gowordladder/words"
	"sync"
)

type Solver struct {
	puzzle              *Puzzle
	exploredCount       int
	solutions           []Solution
	startWord           words.Word
	endWord             words.Word
	reversed            bool
	maximumLadderLength int
	endDistances        WordDistanceMap
	sync                *sync.Mutex
	waitGroup           *sync.WaitGroup
}

type incrementable interface {
	incrementExplored()
}

func NewSolver(puzzle *Puzzle) (result *Solver) {
	return &Solver{
		puzzle: puzzle,
		sync:   &sync.Mutex{},
	}
}

func (s *Solver) Solve(maximumLadderLength int) []Solution {
	s.exploredCount = 0
	s.solutions = make([]Solution, 0)
	s.maximumLadderLength = maximumLadderLength
	if maximumLadderLength < 1 {
		return s.solutions
	}
	s.startWord = s.puzzle.startWord
	s.endWord = s.puzzle.endWord
	s.reversed = false

	diffs := s.startWord.Differences(s.endWord)
	switch diffs {
	case 0:
		s.addSolution(newSolution(s.startWord))
		return s.solutions
	case 1:
		s.addSolution(newSolution(s.startWord, s.endWord))
		switch maximumLadderLength {
		case 2:
			// maximum ladder is 2 - so we already have the only answer...
			return s.solutions
		case 3:
			s.shortCircuitLadderLength3()
			return s.solutions
		}
	case 2:
		if maximumLadderLength == 3 {
			s.shortCircuitLadderLength3()
			return s.solutions
		}
	}

	// begin with the word that has the least number of linked words...
	// (this reduces the number of pointless solution candidates explored!)
	s.reversed = len(s.startWord.LinkedWords()) > len(s.endWord.LinkedWords())
	if s.reversed {
		s.startWord = s.puzzle.endWord
		s.endWord = s.puzzle.startWord
	}
	limit := maximumLadderLength - 1
	s.endDistances = NewWordDistanceMap(s.endWord, &limit)

	s.waitGroup = &sync.WaitGroup{}
	for _, linkedWord := range s.startWord.LinkedWords() {
		if s.endDistances.Reachable(linkedWord, maximumLadderLength) {
			s.waitGroup.Add(1)
			go s.explore(newCandidateSolution(s, s.startWord, linkedWord))
		}
	}
	s.waitGroup.Wait()
	return s.solutions
}

func (s *Solver) explore(candidate *candidateSolution) {
	defer s.waitGroup.Done()
	lastWord := candidate.lastWord()
	if lastWord == s.endWord {
		s.sync.Lock()
		defer s.sync.Unlock()
		s.addSolution(candidate.asFoundSolution(s.reversed))
		return
	}
	ladderLen := len(candidate.ladder)
	if ladderLen < s.maximumLadderLength {
		newMax := s.maximumLadderLength - ladderLen
		for _, linkedWord := range lastWord.LinkedWords() {
			if !candidate.seen(linkedWord) && s.endDistances.Reachable(linkedWord, newMax) {
				s.waitGroup.Add(1)
				go s.explore(candidate.spawn(linkedWord))
			}
		}
	}
}

func (s *Solver) shortCircuitLadderLength3() {
	// we can determine solutions by convergence of the two linked word sets...
	startSet := make(map[string]words.Word, len(s.startWord.LinkedWords()))
	for _, w := range s.startWord.LinkedWords() {
		startSet[w.ActualWord()] = w
	}
	for _, w := range s.endWord.LinkedWords() {
		if _, ok := startSet[w.ActualWord()]; ok {
			s.addSolution(newSolution(s.startWord, w, s.endWord))
		}
	}
}

func (s *Solver) addSolution(sol Solution) {
	s.solutions = append(s.solutions, sol)
}

func (s *Solver) incrementExplored() {
	if s.sync != nil {
		defer s.sync.Unlock()
		s.sync.Lock()
	}
	s.exploredCount++
}

func (s *Solver) ExploredCount() int {
	return s.exploredCount
}
