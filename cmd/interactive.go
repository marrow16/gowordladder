package main

import (
	"bufio"
	"fmt"
	"gowordladder/generator"
	"gowordladder/solving"
	"gowordladder/words"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	terminalColourRed   = "\u001b[31m"
	terminalColourGreen = "\u001b[32m"
	terminalColourBlack = "\u001b[0m"
)

const (
	appName             = "WordLadder"
	prompt              = appName + "> "
	promptStartWord     = prompt + "Enter start word: "
	promptFinalWord     = prompt + "Enter final word: "
	minimumWordLength   = 2
	maximumWordLength   = 15
	minimumLadderLength = 2
)

type interactive struct {
	onStep              int
	dictionary          *words.Dictionary
	dictionaryLoadTime  time.Duration
	startWord           *words.Word
	endWord             *words.Word
	maximumLadderLength int
}

func newInteractive(args []string) *interactive {
	result := &interactive{onStep: 0}
	if len(args) >= 1 {
		println(promptStartWord + green(args[0]))
		if result.setStartWord(args[0]) {
			result.onStep++
			if len(args) >= 2 {
				println(promptFinalWord + green(args[1]))
				if result.setEndWord(args[1]) {
					result.onStep++
				}
			}
		}
	}
	return result
}

func (i *interactive) Run() {
	again := true
	for again {
		for i.onStep < 3 {
			i.processStepInput()
		}
		i.onStep = 0

		i.solve()

		reader := bufio.NewReader(os.Stdin)
		println("")
		print(prompt + "Run again [y/n]: ")
		again = false
		if input, err := reader.ReadString('\n'); err == nil {
			again = input == "\n" || input[:len(input)-1] == "y"
			println("")
		}
	}
}

func (i *interactive) processStepInput() {
	switch i.onStep {
	case 0:
		print(promptStartWord + terminalColourGreen)
	case 1:
		print(promptFinalWord + terminalColourGreen)
	case 2:
		print(prompt + "Maximum ladder length? [" + fmt.Sprintf("%d-%d", minimumLadderLength, i.dictionary.MaxSteps()) + ", or return]: " + terminalColourGreen)
	}
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	print(terminalColourBlack)
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

func (i *interactive) loadDictionary(wordLength int) {
	start := time.Now()
	i.dictionary = words.NewDictionary(wordLength)
	i.dictionaryLoadTime = time.Now().Sub(start)
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func isLengths(input string) (is bool, wordLength int, ladderLength int) {
	ladderLength = -1
	if n, err := strconv.Atoi(input); err == nil && n >= minimumWordLength && n <= maximumWordLength {
		is = true
		wordLength = n
	} else if parts := strings.Split(input, ","); len(parts) == 2 {
		if n1, err := strconv.Atoi(parts[0]); err == nil && n1 >= minimumWordLength && n1 <= maximumWordLength {
			wordLength = n1
			if n2, err := strconv.Atoi(parts[1]); err == nil && n2 >= minimumLadderLength {
				ladderLength = n2
				is = true
			}
		}
	}
	return is, wordLength, ladderLength
}

func (i *interactive) setStartWord(input string) bool {
	if is, wordLength, ladderLength := isLengths(input); is {
		i.loadDictionary(wordLength)
		if ladderLength == -1 {
			candidates := i.dictionary.WordsWithSteps(3)
			i.startWord = candidates[rng.Intn(len(candidates))]
			println(green(fmt.Sprintf("                      Random: %s", i.startWord)))
			return true
		}
		if ladderLength > i.dictionary.MaxSteps() {
			println(red(fmt.Sprintf("            Dictionary %d-letters does not have ladders of length %d (max is %d)!", wordLength, ladderLength, i.dictionary.MaxSteps())))
			return false
		}
		puzzle, err := generator.GeneratePuzzle(wordLength, ladderLength, nil, nil)
		if err != nil {
			println(red(fmt.Sprintf("            %s", err.Error())))
			return false
		}
		i.startWord = puzzle.StartWord
		i.endWord = puzzle.EndWord
		i.maximumLadderLength = ladderLength
		println(green(fmt.Sprintf("                      Random: %s", i.startWord)))
		i.onStep++
		println(promptFinalWord + green(i.endWord.String()))
		println(prompt + "   Ladder length: " + green(strconv.Itoa(ladderLength)))
		i.onStep++
		return true
	}
	if len(input) < minimumWordLength || len(input) > maximumWordLength {
		println(red(fmt.Sprintf("            Enter a word with between %d and %d characters!", minimumWordLength, maximumWordLength)))
		println(red(fmt.Sprintf("            Or enter a number (%d to %d) for random words of that length", minimumWordLength, maximumWordLength)))
		return false
	}
	i.loadDictionary(len(input))
	if input == strings.Repeat("?", len(input)) {
		ladderLen := rng.Intn(i.dictionary.MaxSteps()-2) + 3
		candidates := i.dictionary.WordsWithSteps(ladderLen)
		i.startWord = candidates[rng.Intn(len(candidates))]
		println(green(fmt.Sprintf("                      Random: %s", i.startWord)))
		return true
	}
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

func (i *interactive) setEndWord(input string) bool {
	if len(input) == 0 {
		wm := words.NewWordDistanceMap(i.startWord, nil)
		candidates := wm.Words()
		if len(candidates) == 0 {
			println(red("            Cannot find random final word (try again)!"))
			return false
		}
		input = candidates[rng.Intn(len(candidates))]
		i.endWord, _ = i.dictionary.Word(input)
		println(green(fmt.Sprintf("                      Random: %s", i.endWord)))
		return true
	}
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

func (i *interactive) setMaximumLadderLength(input string) bool {
	if len(input) == 0 {
		println(green("            No answer - assuming auto calc of minimum ladder length"))
		i.maximumLadderLength = -1
		return true
	}
	if m, err := strconv.Atoi(input); err != nil {
		println(red("            Cannot convert '" + input + "' to int"))
		return false
	} else if m < minimumLadderLength || m > i.dictionary.MaxSteps() {
		println(red(fmt.Sprintf("            Please enter a number between %d and %d", minimumLadderLength, i.dictionary.MaxSteps())))
		return false
	} else {
		i.maximumLadderLength = m
	}
	return true
}

func (i *interactive) solve() {
	println("Took " + green(fmt.Sprintf("%s", i.dictionaryLoadTime)) + " to load dictionary")
	puzzle := solving.NewPuzzle(i.startWord, i.endWord)
	if i.maximumLadderLength == -1 {
		start := time.Now()
		minLen, solvable := puzzle.CalculateMinimumLadderLength()
		dur := time.Now().Sub(start)
		if !solvable {
			println(red("Cannot solve `" + i.startWord.String() + "' to '" + i.endWord.String() + "'!"))
		}
		i.maximumLadderLength = minLen
		println("Took " + green(fmt.Sprintf("%s", dur)) +
			" to determine minimum ladder length of " + green(fmt.Sprintf("%d", minLen)))
	}
	solver := solving.NewSolver(puzzle)
	start := time.Now()
	solutions := solver.Solve(i.maximumLadderLength)
	dur := time.Now().Sub(start)
	if len(solutions) == 0 {
		println(red(fmt.Sprintf("Took %s to find no solutions (explored %d candidates)", dur, solver.ExploredCount())))
	} else {
		println(
			"Took " + green(fmt.Sprintf("%s", dur)) +
				" to find " + green(fmt.Sprintf("%d", len(solutions))) +
				" solutions (explored " + green(fmt.Sprintf("%d", solver.ExploredCount())) + " candidates)")
		i.displaySolutions(solutions)
	}
}

func (i *interactive) displaySolutions(solutions []*solving.Solution) {
	SortSolutions(solutions)
	pageStart := 0
	length := len(solutions)
	for pageStart < length {
		more := ""
		if pageStart > 0 {
			more = " more"
		}
		print(prompt + "List" + more + " solutions? (Enter 'n' for no, 'y' or return for next 10, 'all' for all or how many): ")
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
			println(fmt.Sprintf("%d/%d", s+pageStart+1, length) + " " + green(solutions[s+pageStart].String()))
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
