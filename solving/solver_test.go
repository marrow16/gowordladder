package solving

import (
	"gowordladder/test"
	"gowordladder/words"
	"testing"
)

func TestSolveCatToDog(t *testing.T) {
	var dict = words.LoadDictionary(3)
	var cat, _ = dict.Word("cat")
	var dog, _ = dict.Word("dog")
	var solver = NewSolver(NewPuzzle(cat, dog))
	var solutions = solver.Solve(4, true)
	test.AssertEqualsInt(t, 4, len(*solutions))
	test.AssertEqualsInt(t, 38, solver.ExploredCount())

	var midWords = &map[string]int{}
	for _, solution := range *solutions {
		test.AssertEqualsInt(t, 4, len(solution.ladder))
		test.AssertEqualsString(t, "CAT", solution.ladder[0].ActualWord())
		test.AssertEqualsString(t, "DOG", solution.ladder[3].ActualWord())
		if x, ok := (*midWords)[solution.ladder[1].ActualWord()]; ok {
			(*midWords)[solution.ladder[1].ActualWord()] = x + 1
		} else {
			(*midWords)[solution.ladder[1].ActualWord()] = 1
		}
		if x, ok := (*midWords)[solution.ladder[2].ActualWord()]; ok {
			(*midWords)[solution.ladder[2].ActualWord()] = x + 1
		} else {
			(*midWords)[solution.ladder[2].ActualWord()] = 1
		}
	}
	test.AssertEqualsInt(t, 5, len(*midWords))
	test.AssertEqualsInt(t, 2, (*midWords)["CAG"])
	test.AssertEqualsInt(t, 2, (*midWords)["COG"])
	test.AssertEqualsInt(t, 2, (*midWords)["COT"])
	test.AssertEqualsInt(t, 1, (*midWords)["DAG"])
	test.AssertEqualsInt(t, 1, (*midWords)["DOT"])
}

func TestSolveColdToWarmAndWarmToCold(t *testing.T) {
	var dict = words.LoadDictionary(4)
	var cold, _ = dict.Word("cold")
	var warm, _ = dict.Word("warm")
	var solver = NewSolver(NewPuzzle(cold, warm))
	var solutions = solver.Solve(5, false)
	test.AssertEqualsInt(t, 7, len(*solutions))
	var explored = solver.ExploredCount()
	test.AssertEqualsInt(t, 33, explored)

	// now do it the other way around..
	solver = NewSolver(NewPuzzle(warm, cold))
	solutions = solver.Solve(5, false)
	test.AssertEqualsInt(t, 7, len(*solutions))
	test.AssertEqualsInt(t, explored, solver.ExploredCount())
}

func TestSolveKataToJava(t *testing.T) {
	var dict = words.LoadDictionary(4)
	var kata, _ = dict.Word("kata")
	var java, _ = dict.Word("java")
	var solver = NewSolver(NewPuzzle(kata, java))
	var solutions = solver.Solve(3, false)
	test.AssertEqualsInt(t, 1, len(*solutions))

	var solution = (*solutions)[0]
	test.AssertEqualsInt(t, 3, len(solution.ladder))
	test.AssertEqualsString(t, "KATA", solution.ladder[0].ActualWord())
	test.AssertEqualsString(t, "KAVA", solution.ladder[1].ActualWord())
	test.AssertEqualsString(t, "JAVA", solution.ladder[2].ActualWord())
}

func TestSameWordSolvable(t *testing.T) {
	var dict = words.LoadDictionary(3)
	var cat, _ = dict.Word("cat")
	var solver = NewSolver(NewPuzzle(cat, cat))
	var solutions = solver.Solve(1, true)
	test.AssertEqualsInt(t, 1, len(*solutions))
	test.AssertEqualsInt(t, 0, solver.ExploredCount())
}

func TestOneLetterDifferenceIsSolvable(t *testing.T) {
	var dict = words.LoadDictionary(3)
	var cat, _ = dict.Word("cat")
	var cot, _ = dict.Word("cot")
	var solver = NewSolver(NewPuzzle(cat, cot))
	var solutions = solver.Solve(2, true)
	test.AssertEqualsInt(t, 1, len(*solutions))
	test.AssertEqualsInt(t, 0, solver.ExploredCount())
}

func TestTwoLettersDifferenceIsSolvable(t *testing.T) {
	var dict = words.LoadDictionary(3)
	var cat, _ = dict.Word("cat")
	var bar, _ = dict.Word("bar")
	var solver = NewSolver(NewPuzzle(cat, bar))
	var solutions = solver.Solve(3, true)
	test.AssertEqualsInt(t, 2, len(*solutions))
	test.AssertEqualsInt(t, 0, solver.ExploredCount())
}

func TestEverythingUnsolvableWithBadMaxLadderLength(t *testing.T) {
	var dict = words.LoadDictionary(3)
	var cat, _ = dict.Word("cat")
	var dog, _ = dict.Word("dog")
	var solver = NewSolver(NewPuzzle(cat, dog))
	var solutions = solver.Solve(-1, true)
	test.AssertEqualsInt(t, 0, len(*solutions))
	solutions = solver.Solve(0, true)
	test.AssertEqualsInt(t, 0, len(*solutions))
	solutions = solver.Solve(2, true)
	test.AssertEqualsInt(t, 0, len(*solutions))
	solutions = solver.Solve(3, true)
	test.AssertEqualsInt(t, 0, len(*solutions))
	solutions = solver.Solve(4, true)
	test.AssertTrue(t, len(*solutions) > 0)
}

func TestShortCircuitOnOneLetterDifference(t *testing.T) {
	var dict = words.LoadDictionary(3)
	var cat, _ = dict.Word("cat")
	var cot, _ = dict.Word("cot")
	var solver = NewSolver(NewPuzzle(cat, cot))
	var solutions = solver.Solve(3, true)
	test.AssertEqualsInt(t, 3, len(*solutions))
	test.AssertEqualsInt(t, 0, solver.ExploredCount())

	test.AssertEqualsInt(t, 2, len((*solutions)[0].ladder))
	test.AssertEqualsInt(t, 3, len((*solutions)[1].ladder))
	test.AssertEqualsInt(t, 3, len((*solutions)[2].ladder))
}
