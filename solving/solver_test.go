package solving

import (
	"github.com/stretchr/testify/assert"
	"gowordladder/words"
	"testing"
)

func TestSolveCatToDog(t *testing.T) {
	dict := words.LoadDictionary(3)
	cat, _ := dict.Word("cat")
	dog, _ := dict.Word("dog")
	solver := NewSolver(NewPuzzle(cat, dog))
	solutions := solver.Solve(4)
	assert.Equal(t, 4, len(*solutions))
	assert.Equal(t, 10, solver.ExploredCount())

	midWords := &map[string]int{}
	for _, solution := range *solutions {
		assert.Equal(t, 4, len(solution.Ladder()))
		assert.Equal(t, "CAT", solution.Ladder()[0].ActualWord())
		assert.Equal(t, "DOG", solution.Ladder()[3].ActualWord())
		if x, ok := (*midWords)[solution.Ladder()[1].ActualWord()]; ok {
			(*midWords)[solution.Ladder()[1].ActualWord()] = x + 1
		} else {
			(*midWords)[solution.Ladder()[1].ActualWord()] = 1
		}
		if x, ok := (*midWords)[solution.Ladder()[2].ActualWord()]; ok {
			(*midWords)[solution.Ladder()[2].ActualWord()] = x + 1
		} else {
			(*midWords)[solution.Ladder()[2].ActualWord()] = 1
		}
	}
	assert.Equal(t, 5, len(*midWords))
	assert.Equal(t, 2, (*midWords)["CAG"])
	assert.Equal(t, 2, (*midWords)["COG"])
	assert.Equal(t, 2, (*midWords)["COT"])
	assert.Equal(t, 1, (*midWords)["DAG"])
	assert.Equal(t, 1, (*midWords)["DOT"])
}

func TestSolveColdToWarmAndWarmToCold(t *testing.T) {
	dict := words.LoadDictionary(4)
	cold, _ := dict.Word("cold")
	warm, _ := dict.Word("warm")
	solver := NewSolver(NewPuzzle(cold, warm))
	solutions := solver.Solve(5)
	assert.Equal(t, 7, len(*solutions))
	explored := solver.ExploredCount()
	assert.Equal(t, 21, explored)

	// now do it the other way around..
	solver = NewSolver(NewPuzzle(warm, cold))
	solutions = solver.Solve(5)
	assert.Equal(t, 7, len(*solutions))
	assert.Equal(t, explored, solver.ExploredCount())
}

func TestSolveKataToJava(t *testing.T) {
	dict := words.LoadDictionary(4)
	kata, _ := dict.Word("kata")
	java, _ := dict.Word("java")
	solver := NewSolver(NewPuzzle(kata, java))
	solutions := solver.Solve(3)
	assert.Equal(t, 1, len(*solutions))

	solution := (*solutions)[0]
	assert.Equal(t, 3, len(solution.Ladder()))
	assert.Equal(t, "KATA", solution.Ladder()[0].ActualWord())
	assert.Equal(t, "KAVA", solution.Ladder()[1].ActualWord())
	assert.Equal(t, "JAVA", solution.Ladder()[2].ActualWord())
}

func TestSameWordSolvable(t *testing.T) {
	dict := words.LoadDictionary(3)
	cat, _ := dict.Word("cat")
	solver := NewSolver(NewPuzzle(cat, cat))
	solutions := solver.Solve(1)
	assert.Equal(t, 1, len(*solutions))
	assert.Equal(t, 0, solver.ExploredCount())
}

func TestOneLetterDifferenceIsSolvable(t *testing.T) {
	dict := words.LoadDictionary(3)
	cat, _ := dict.Word("cat")
	cot, _ := dict.Word("cot")
	solver := NewSolver(NewPuzzle(cat, cot))
	solutions := solver.Solve(2)
	assert.Equal(t, 1, len(*solutions))
	assert.Equal(t, 0, solver.ExploredCount())
}

func TestTwoLettersDifferenceIsSolvable(t *testing.T) {
	dict := words.LoadDictionary(3)
	cat, _ := dict.Word("cat")
	bar, _ := dict.Word("bar")
	solver := NewSolver(NewPuzzle(cat, bar))
	solutions := solver.Solve(3)
	assert.Equal(t, 2, len(*solutions))
	assert.Equal(t, 0, solver.ExploredCount())
}

func TestEverythingUnsolvableWithBadMaxLadderLength(t *testing.T) {
	dict := words.LoadDictionary(3)
	cat, _ := dict.Word("cat")
	dog, _ := dict.Word("dog")
	solver := NewSolver(NewPuzzle(cat, dog))
	solutions := solver.Solve(-1)
	assert.Equal(t, 0, len(*solutions))
	solutions = solver.Solve(0)
	assert.Equal(t, 0, len(*solutions))
	solutions = solver.Solve(2)
	assert.Equal(t, 0, len(*solutions))
	solutions = solver.Solve(3)
	assert.Equal(t, 0, len(*solutions))
	solutions = solver.Solve(4)
	assert.True(t, len(*solutions) > 0)
}

func TestShortCircuitOnOneLetterDifference(t *testing.T) {
	dict := words.LoadDictionary(3)
	cat, _ := dict.Word("cat")
	cot, _ := dict.Word("cot")
	solver := NewSolver(NewPuzzle(cat, cot))
	solutions := solver.Solve(3)
	assert.Equal(t, 3, len(*solutions))
	assert.Equal(t, 0, solver.ExploredCount())

	assert.Equal(t, 2, len((*solutions)[0].Ladder()))
	assert.Equal(t, 3, len((*solutions)[1].Ladder()))
	assert.Equal(t, 3, len((*solutions)[2].Ladder()))
}
