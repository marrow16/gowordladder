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
	var args = os.Args
	if len(args) >= 3 {
		solveNow(args)
	} else {
		NewInteractive(args).Run()
	}
}

func solveNow(args []string)  {
	var first = args[1]
	var second = args[2]
	if len(first) != len(second) {
		panic("Words supplied as args must be same length")
	}
	var startTime = time.Now().UnixNano()
	var dictionary = words.LoadDictionary(len(first))
	var took = time.Now().UnixNano() - startTime
	var startWord, ok = dictionary.Word(first)
	if !ok {
		panic("Start word '" + first + "' does not exist in dictionary")
	}
	var endWord, ok2 = dictionary.Word(second)
	if !ok2 {
		panic("End word '" + second + "' does not exist in dictionary")
	}
	println(fmt.Sprintf("Took %dms to load dictionary", took / 1000000))

	var puzzle = solving.NewPuzzle(startWord, endWord)
	var maxLadderLength int
	if len(args) > 3 {
		if i, err := strconv.Atoi(args[3]); err != nil {
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
		println(fmt.Sprintf("Took %dms to determine minimum ladder length of %d", took / 1000000, maxLadderLength))
	}
	var solver = solving.NewSolver(puzzle)
	startTime = time.Now().UnixNano()
	var solutions = solver.Solve(maxLadderLength, true)
	took = time.Now().UnixNano() - startTime
	println(fmt.Sprintf("Took %dms to find %d solutions (explored %d solutions)", took / 1000000, len(*solutions), solver.ExploredCount()))
	solving.SortSolutions(*solutions)
	var l = len(*solutions)
	for i, s := range *solutions {
		println(fmt.Sprintf("%d/%d %s", i + 1, l, s.ToString()))
	}
}
