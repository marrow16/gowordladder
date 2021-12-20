package main

import (
	"fmt"
	"gowordladder/solving"
	"gowordladder/words"
	"os"
	"strconv"
	"time"
)

func main() {
	args := os.Args
	if len(args) >= 3 {
		solveNow(args[1:])
	} else {
		NewInteractive(args[1:]).Run()
	}
}

func solveNow(args []string) {
	first := args[0]
	var second = args[1]
	if len(first) != len(second) {
		panic("Words supplied as args must be same length")
	}
	startTime := time.Now().UnixNano()
	dictionary := words.LoadDictionary(len(first))
	took := time.Now().UnixNano() - startTime
	startWord, ok := dictionary.Word(first)
	if !ok {
		panic("Start word '" + first + "' does not exist in dictionary")
	}
	endWord, ok := dictionary.Word(second)
	if !ok {
		panic("End word '" + second + "' does not exist in dictionary")
	}
	println(fmt.Sprintf("Took %dms to load dictionary", took/1000000))

	puzzle := solving.NewPuzzle(startWord, endWord)
	var maxLadderLength int
	if len(args) > 2 {
		if i, err := strconv.Atoi(args[2]); err != nil {
			panic("Cannot convert '" + args[3] + "' to int")
		} else {
			maxLadderLength = i
		}
		println("Using specified maximum ladder length " + fmt.Sprintf("%d", maxLadderLength))
	} else {
		startTime = time.Now().UnixNano()
		min, solvable := puzzle.CalculateMinimumLadderLength()
		took = time.Now().UnixNano() - startTime
		if !solvable {
			panic("Cannot solve `" + first + "' to '" + second + "'!")
		}
		maxLadderLength = min
		println(fmt.Sprintf("Took %dms to determine minimum ladder length of %d", took/1000000, maxLadderLength))
	}
	solver := solving.NewSolver(puzzle)
	startTime = time.Now().UnixNano()
	solutions := solver.Solve(maxLadderLength, true)
	took = time.Now().UnixNano() - startTime
	println(fmt.Sprintf("Took %dms to find %d solutions (explored %d solutions)", took/1000000, len(*solutions), solver.ExploredCount()))
	solving.SortSolutions(*solutions)
	l := len(*solutions)
	for i, s := range *solutions {
		println(fmt.Sprintf("%d/%d %s", i+1, l, s.ToString()))
	}
}
