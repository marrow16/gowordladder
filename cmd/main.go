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
		newInteractive(args).Run()
	}
}

func solveNow(args []string) {
	first := args[0]
	second := args[1]
	third := args[2]
	if len(first) != len(second) {
		panic(any("Words supplied as args must be same length"))
	}
	start := time.Now()
	dictionary := words.NewDictionary(len(first))
	dur := time.Now().Sub(start)
	startWord, ok := dictionary.Word(first)
	if !ok {
		panic(any(fmt.Sprintf("Start word '%s' does not exist in dictionary", first)))
	}
	endWord, ok := dictionary.Word(second)
	if !ok {
		panic(any(fmt.Sprintf("End word '%s' does not exist in dictionary", second)))
	}
	println(fmt.Sprintf("Took %s to load dictionary", dur))

	puzzle := solving.NewPuzzle(startWord, endWord)
	var maxLadderLength int
	if i, err := strconv.Atoi(third); err == nil && i > 0 {
		maxLadderLength = i
	} else {
		start := time.Now()
		maxLadderLength, solvable := puzzle.CalculateMinimumLadderLength()
		dur := time.Now().Sub(start)
		if !solvable {
			panic(any(fmt.Sprintf("Cannot solve `%s' to '%s'!", first, second)))
		}
		println(fmt.Sprintf("Took %s to determine minimum ladder length of %d", dur, maxLadderLength))
	}
	println(fmt.Sprintf("Solving %s to %s (maximum steps: %d)",
		green(startWord.String()), green(endWord.String()), maxLadderLength))
	solver := solving.NewSolver(puzzle)
	start = time.Now()
	solutions := solver.Solve(maxLadderLength)
	dur = time.Now().Sub(start)
	println(fmt.Sprintf("Took %s to find %d solutions (explored %d solutions)", dur, len(solutions), solver.ExploredCount()))
	SortSolutions(solutions)
	l := len(solutions)
	for i, s := range solutions {
		println(fmt.Sprintf("%d/%d %s", i+1, l, s.String()))
	}
}

func SortSolutions(solutions []*solving.Solution) {
	sort.Slice(solutions, func(i, j int) bool {
		if len(solutions[i].Ladder()) < len(solutions[j].Ladder()) {
			return true
		} else if len(solutions[i].Ladder()) == len(solutions[j].Ladder()) {
			for idx, w := range solutions[i].Ladder() {
				if w.String() < solutions[j].Ladder()[idx].String() {
					return true
				} else if w.String() > solutions[j].Ladder()[idx].String() {
					return false
				}
			}
		}
		return false
	})
}
