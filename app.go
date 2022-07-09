package main

import (
	"fmt"
	"gowordladder/solving"
	"gowordladder/words"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) >= 3 {
		solveNow(args)
	} else {
		NewInteractive(args).Run()
	}
}

func solveNow(args []string) {
	first := args[0]
	second := args[1]
	third := args[2]
	if len(first) != len(second) {
		panic(any("Words supplied as args must be same length"))
	}
	startTime := time.Now().UnixNano()
	dictionary := words.LoadDictionary(len(first))
	took := time.Now().UnixNano() - startTime
	startWord, ok := dictionary.Word(first)
	if !ok {
		panic(any(fmt.Sprintf("Start word '%s' does not exist in dictionary", first)))
	}
	endWord, ok := dictionary.Word(second)
	if !ok {
		panic(any(fmt.Sprintf("End word '%s' does not exist in dictionary", second)))
	}
	println(fmt.Sprintf("Took %dms to load dictionary", took/1000000))

	puzzle := solving.NewPuzzle(startWord, endWord)
	var maxLadderLength int
	if i, err := strconv.Atoi(third); err == nil && i > 0 {
		maxLadderLength = i
	} else {
		startTime = time.Now().UnixNano()
		min, solvable := puzzle.CalculateMinimumLadderLength()
		took = time.Now().UnixNano() - startTime
		if !solvable {
			panic(any(fmt.Sprintf("Cannot solve `%s' to '%s'!", first, second)))
		}
		maxLadderLength = min
		println(fmt.Sprintf("Took %dms to determine minimum ladder length of %d", took/1000000, maxLadderLength))
	}
	println(fmt.Sprintf("Solving %s to %s (maximum steps: %d)",
		green(startWord.ActualWord()), green(endWord.ActualWord()), maxLadderLength))
	solver := solving.NewSolver(puzzle)
	startTime = time.Now().UnixNano()
	solutions := solver.Solve(maxLadderLength)
	took = time.Now().UnixNano() - startTime
	println(fmt.Sprintf("Took %dms to find %d solutions (explored %d solutions)", took/1000000, len(solutions), solver.ExploredCount()))
	SortSolutions(solutions)
	l := len(solutions)
	for i, s := range solutions {
		println(fmt.Sprintf("%d/%d %s", i+1, l, s.String()))
	}
}

func SortSolutions(solutions []solving.Solution) {
	sort.Slice(solutions, func(i, j int) bool {
		if len(solutions[i].Ladder()) < len(solutions[j].Ladder()) {
			return true
		} else if len(solutions[i].Ladder()) == len(solutions[j].Ladder()) {
			for idx, w := range solutions[i].Ladder() {
				if w.ActualWord() < solutions[j].Ladder()[idx].ActualWord() {
					return true
				} else if w.ActualWord() > solutions[j].Ladder()[idx].ActualWord() {
					return false
				}
			}
		}
		return false
	})
}
