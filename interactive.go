package main

import (
	"bufio"
	"fmt"
	"gowordladder/solving"
	"gowordladder/words"
	"os"
	"strconv"
	"time"
)

const terminalColourRed = "\u001b[31m"
const terminalColourGreen = "\u001b[32m"
const terminalColourBlack = "\u001b[0m"

const appName = "WordLadder"
const prompt = appName + "> "
const minimumWordLength = 2
const maximumWordLength = 15
const minimumLadderLength = 1
const maximumLadderLength = 20

var steps = []struct {
	step   int
	prompt string
}{
	{0, prompt + "Enter start word: "},
	{1, prompt + "Enter final word: "},
	{2, prompt + "Maximum ladder length? [" + fmt.Sprintf("%d-%d", minimumLadderLength, maximumLadderLength) + ", or return]: "},
}

type Interactive struct {
	onStep              int
	dictionary          *words.Dictionary
	dictionaryLoadTime  int64
	startWord           *words.Word
	endWord             *words.Word
	maximumLadderLength int
}

func NewInteractive(args []string) (result *Interactive) {
	result = &Interactive{onStep: 0}
	// TODO fill steps from args
	return
}

func (i *Interactive) Run() {
	again := true
	for again {
		for i.onStep < len(steps) {
			i.processStepInput()
		}
		i.onStep = 0

		i.solve()

		reader := bufio.NewReader(os.Stdin)
		println("")
		print("Run again [y/n]: ")
		again = false
		if input, err := reader.ReadString('\n'); err == nil {
			again = input[:len(input)-1] == "y"
		}
	}
}

func (i *Interactive) processStepInput() {
	reader := bufio.NewReader(os.Stdin)
	print(steps[i.onStep].prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		println(red(fmt.Sprintf("Error: %s", err)))
		return
	}
	input = input[:len(input)-1]
	ok := false
	switch i.onStep {
	case 0:
		ok = i.setStartWord(input)
		break
	case 1:
		ok = i.setEndWord(input)
		break
	case 2:
		ok = i.setMaximumLadderLength(input)
		break
	}
	if ok {
		i.onStep++
	}
}

func (i *Interactive) loadDictionary(wordLength int) {
	startTime := time.Now().UnixNano()
	i.dictionary = words.LoadDictionary(wordLength)
	i.dictionaryLoadTime = time.Now().UnixNano() - startTime
}

func (i *Interactive) setStartWord(input string) bool {
	if len(input) < minimumWordLength || len(input) > maximumLadderLength {
		println(red(fmt.Sprintf("            Please enter a word with between %d and %d characters!", minimumWordLength, maximumWordLength)))
		return false
	}
	i.loadDictionary(len(input))
	if w, ok := i.dictionary.Word(input); !ok {
		println(red(fmt.Sprintf("            Word '%s' does not exist!", input)))
		return false
	} else if w.IsIsland() {
		println(red(fmt.Sprintf("            Word '%s' is an island word (cannot change single letter to form another word)", input)))
		return false
	} else {
		i.startWord = w
	}
	return true
}

func (i *Interactive) setEndWord(input string) bool {
	if len(input) != i.dictionary.WordLength() {
		println(red("            Final word length must match start word length!"))
		return false
	}
	if w, ok := i.dictionary.Word(input); !ok {
		println(red(fmt.Sprintf("            Word '%s' does not exist!", input)))
		return false
	} else if w.IsIsland() {
		println(red(fmt.Sprintf("            Word '%s' is an island word (cannot change single letter to form another word)", input)))
		return false
	} else {
		i.endWord = w
	}
	return true
}

func (i *Interactive) setMaximumLadderLength(input string) bool {
	if len(input) == 0 {
		println(green("            No answer - assuming auto calc of minimum ladder length"))
		i.maximumLadderLength = -1
		return true
	}
	if m, err := strconv.Atoi(input); err != nil {
		println(red("            Cannot convert '" + input + "' to int"))
		return false
	} else if m < minimumLadderLength || m > maximumLadderLength {
		println(red(fmt.Sprintf("            Please enter a number between %d and %d", minimumLadderLength, maximumLadderLength)))
		return false
	} else {
		i.maximumLadderLength = m
	}
	return true
}

func (i *Interactive) solve() {
	println("Took " + green(fmt.Sprintf("%dms", i.dictionaryLoadTime/1000000)) + " to load dictionary")
	puzzle := solving.NewPuzzle(i.startWord, i.endWord)
	if i.maximumLadderLength == -1 {
		startTime := time.Now().UnixNano()
		min, solvable := puzzle.CalculateMinimumLadderLength()
		took := time.Now().UnixNano() - startTime
		if !solvable {
			println(red("Cannot solve `" + i.startWord.ActualWord() + "' to '" + i.endWord.ActualWord() + "'!"))
		}
		i.maximumLadderLength = min
		println("Took " + green(fmt.Sprintf("%dms", took/1000000)) +
			" to determine minimum ladder length of " + green(fmt.Sprintf("%d", min)))
	}
	solver := solving.NewSolver(puzzle)
	startTime := time.Now().UnixNano()
	solutions := solver.Solve(i.maximumLadderLength, true)
	took := time.Now().UnixNano() - startTime
	if len(*solutions) == 0 {
		println(red(fmt.Sprintf("Took %dms to find no solutions (explored %d solutions)", took/1000000, solver.ExploredCount())))
	} else {
		println(
			"Took " + green(fmt.Sprintf("%dms", took/1000000)) +
				" to find " + green(fmt.Sprintf("%d", len(*solutions))) +
				" solutions (explored " + green(fmt.Sprintf("%d", solver.ExploredCount())) + " solutions)")
		i.displaySolutions(solutions)
	}
}

func (i *Interactive) displaySolutions(solutions *[]*solving.Solution) {
	solving.SortSolutions(*solutions)
	pageStart := 0
	length := len(*solutions)
	for pageStart < length {
		more := ""
		if pageStart > 0 {
			more = " more"
		}
		print("List" + more + " solutions? (Enter 'n' for no, 'y' or return for next 10, 'all' for all or how many): ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			input = " "
		}
		input = input[:len(input)-1]
		if "n" == input {
			return
		}
		limit := 10
		if "all" == input {
			limit = length
		} else if len(input) > 0 && "y" != input {
			if m, err := strconv.Atoi(input); err == nil {
				limit = m
			}
		}
		for s := 0; s < limit && (s+pageStart) < length; s++ {
			println(fmt.Sprintf("%d/%d", s+pageStart+1, length) + " " + green((*solutions)[s+pageStart].ToString()))
		}
		pageStart = pageStart + limit
	}
}

func green(msg string) string {
	return terminalColourGreen + msg + terminalColourBlack
}

func red(msg string) string {
	return terminalColourRed + msg + terminalColourBlack
}
