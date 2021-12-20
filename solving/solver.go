package solving

import (
	"gowordladder/words"
	"sync"
)

type Solver struct {
	puzzle              *Puzzle
	exploredCount       int
	solutions           *[]*Solution
	startWord           *words.Word
	endWord             *words.Word
	reversed            bool
	maximumLadderLength int
	endDistances        *WordDistanceMap
	sync                *sync.Mutex
}

func NewSolver(puzzle *Puzzle) (result *Solver) {
	result = &Solver{puzzle: puzzle}
	return
}

func (s *Solver) Solve(maximumLadderLength int, async bool) (result *[]*Solution) {
	s.exploredCount = 0
	s.solutions = &[]*Solution{}
	result = &*s.solutions
	s.maximumLadderLength = maximumLadderLength
	if maximumLadderLength < 1 {
		return
	}
	s.startWord = s.puzzle.startWord
	s.endWord = s.puzzle.endWord
	s.reversed = false

	diffs := s.startWord.Differences(s.endWord)
	switch diffs {
	case 0:
		s.addSolution(s.startWord)
		return
	case 1:
		s.addSolution(s.startWord, s.endWord)
		switch maximumLadderLength {
		case 2:
			// maximum ladder is 2 so we already have the only answer...
			return
		case 3:
			s.shortCircuitLadderLength3()
			return
		}
	case 2:
		if maximumLadderLength == 3 {
			s.shortCircuitLadderLength3()
			return
		}
	}

	// begin with the word that has the least number of linked words...
	// (this reduces the number of pointless solution candidates explored!)
	s.reversed = len(*s.startWord.LinkedWords) > len(*s.endWord.LinkedWords)
	if s.reversed {
		s.startWord = s.puzzle.endWord
		s.endWord = s.puzzle.startWord
	}
	s.endDistances = NewWordDistanceMap(s.endWord)

	if async {
		s.sync = &sync.Mutex{}
		var waitGroup sync.WaitGroup
		defer waitGroup.Wait()
		waitGroup.Add(1)
		for _, linkedWord := range *s.startWord.LinkedWords {
			if s.endDistances.Reachable(linkedWord, maximumLadderLength) {
				waitGroup.Add(1)
				go s.solveAsync(newCandidateSolution(s, s.startWord, linkedWord), &waitGroup)
			}
		}
		waitGroup.Done()
	} else {
		for _, linkedWord := range *s.startWord.LinkedWords {
			if s.endDistances.Reachable(linkedWord, maximumLadderLength) {
				s.solve(newCandidateSolution(s, s.startWord, linkedWord))
			}
		}
	}
	return
}

func (s *Solver) solve(candidate *candidateSolution) {
	lastWord := candidate.lastWord()
	if *lastWord == *s.endWord {
		*s.solutions = append(*s.solutions, candidate.asFoundSolution(s.reversed))
		return
	}
	ladderLen := len(*candidate.ladder)
	if ladderLen < s.maximumLadderLength {
		newMax := s.maximumLadderLength - ladderLen
		for _, linkedWord := range *lastWord.LinkedWords {
			if !candidate.seen(linkedWord) && s.endDistances.Reachable(linkedWord, newMax) {
				s.solve(candidate.spawn(linkedWord))
			}
		}
	}
}

func (s *Solver) solveAsync(candidate *candidateSolution, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	lastWord := candidate.lastWord()
	if *lastWord == *s.endWord {
		s.sync.Lock()
		defer s.sync.Unlock()
		*s.solutions = append(*s.solutions, candidate.asFoundSolution(s.reversed))
		return
	}
	ladderLen := len(*candidate.ladder)
	if ladderLen < s.maximumLadderLength {
		newMax := s.maximumLadderLength - ladderLen
		for _, linkedWord := range *lastWord.LinkedWords {
			if !candidate.seen(linkedWord) && s.endDistances.Reachable(linkedWord, newMax) {
				waitGroup.Add(1)
				go s.solveAsync(candidate.spawn(linkedWord), waitGroup)
			}
		}
	}
}

func (s *Solver) shortCircuitLadderLength3() {
	// we can determine solutions by convergence of the two linked word sets...
	startSet := make(map[string]*words.Word, len(*s.startWord.LinkedWords))
	for _, w := range *s.startWord.LinkedWords {
		startSet[w.ActualWord()] = w
	}
	for _, w := range *s.endWord.LinkedWords {
		if _, ok := startSet[w.ActualWord()]; ok {
			s.addSolution(s.startWord, w, s.endWord)
		}
	}
}

func (s *Solver) addSolution(words ...*words.Word) {
	*s.solutions = append(*s.solutions, newSolution(words...))
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
