package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"gowordladder/solving"
	"gowordladder/words"
	"strconv"
	"strings"
	"time"
)

type solveStep int

const (
	solveStartWord solveStep = iota
	solveEndWord
	solveMaxLadder
	solveSolved
)

type viewSolve struct {
	step         solveStep
	startWord    *words.Word
	endWord      *words.Word
	ladderLength int

	currentInput input
	currentError string

	dictionaryLoadTime time.Duration
	minLadderLength    int
	minLengthCalcTime  time.Duration
	solutions          []*solving.Solution
	solveTime          time.Duration
}

func (v *viewSolve) wordLength() int {
	if v.step >= solveStartWord && v.startWord != nil {
		return len(v.startWord.String())
	}
	return 0
}

func (v *viewSolve) content(m *model) (string, *tea.Cursor) {
	const (
		promptStartWord = "            Start word: "
		promptEndWord   = "              End word: "
		promptMaxLadder = " Maximum ladder length: "
		promptLen       = len(promptMaxLadder)
	)
	var sb strings.Builder
	sb.WriteString("\n")
	lines := 1
	cpx := -1
	var s string
	switch v.step {
	case solveStartWord:
		sb.WriteString(promptStartWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: 15}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		}
		lines++
	case solveEndWord:
		sb.WriteString(promptStartWord)
		sb.WriteString(inputStyle.Width(15).Render(v.startWord.String()))
		sb.WriteString("\n" + promptEndWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: 15}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		} else {
			sb.WriteString(helpStyle.Render("  (blank for random word)"))
		}
		lines += 2
	case solveMaxLadder:
		sb.WriteString(promptStartWord)
		sb.WriteString(inputStyle.Width(15).Render(v.startWord.String()))
		sb.WriteString("\n" + promptEndWord)
		sb.WriteString(inputStyle.Width(15).Render(v.endWord.String()))
		sb.WriteString("\n" + promptMaxLadder)
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		} else {
			sb.WriteString(helpStyle.Render("  (optional - blank for auto min)"))
		}
		lines += 3
	case solveSolved:
		sb.WriteString(promptStartWord)
		sb.WriteString(inputStyle.Width(15).Render(v.startWord.String()))
		sb.WriteString("\n" + promptEndWord)
		sb.WriteString(inputStyle.Width(15).Render(v.endWord.String()))
		sb.WriteString("\n" + promptMaxLadder)
		if v.ladderLength == -1 {
			sb.WriteString(inputStyle.Width(2).Render("??"))
		} else {
			sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.ladderLength)))
		}
		sb.WriteString("\n\n Took " + highlightStyle.Render(v.dictionaryLoadTime.String()) + " to load dictionary")
		lines += 5
		if v.ladderLength == -1 {
			sb.WriteString("\n Took " + highlightStyle.Render(v.minLengthCalcTime.String()) + " to determine min ladder length " + highlightStyle.Render(strconv.Itoa(v.minLadderLength)))
			lines++
		}
		sb.WriteString("\n Took " + highlightStyle.Render(v.solveTime.String()) + " to find " + highlightStyle.Render(strconv.Itoa(len(v.solutions))) + " solutions")
		lines++
	}

	if m.height-lines-2 > 0 {
		sb.WriteString(strings.Repeat("\n", m.height-lines-2))
	}
	var csr *tea.Cursor
	if cpx > -1 {
		csr = tea.NewCursor(promptLen+cpx, lines)
	}
	return sb.String(), csr
}

func (v *viewSolve) help() string {
	if v.step == solveSolved && len(v.solutions) > 0 {
		return "ctrl+n: New  •  enter: Solutions  •  ctrl+g: Generate"
	} else {
		return "ctrl+n: New  •  ctrl+g: Generate"
	}
}

func (v *viewSolve) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.currentError = ""
	switch msg.String() {
	case "ctrl+n":
		v.currentInput = nil
		v.currentError = ""
		v.solutions = nil
		v.startWord = nil
		v.endWord = nil
		v.step = solveStartWord
	case "enter":
		switch v.step {
		case solveStartWord:
			return v.enterStartWord(m)
		case solveEndWord:
			return v.enterEndWord(m)
		case solveMaxLadder:
			return v.enterMaxLadder(m)
		case solveSolved:
			if len(v.solutions) > 0 {
				m.showSolutions(v.solutions)
				return nil
			}
		}
	}
	if v.currentInput != nil {
		v.currentInput.key(msg)
	}
	return nil
}

type solveEnterResult struct {
	err      string
	nextStep solveStep
	update   func(v *viewSolve)
}

func (v *viewSolve) update(m *model, msg tea.Msg) tea.Cmd {
	if result, ok := msg.(solveEnterResult); ok {
		if result.err != "" {
			v.currentError = result.err
		} else {
			result.update(v)
			v.step = result.nextStep
		}
	}
	return nil
}

func (v *viewSolve) enterStartWord(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		if len(s) < 2 || len(s) > 15 {
			return solveEnterResult{err: "Please enter a word 2-15 letters"}
		}
		dict := m.loadDictionary(len(s))
		wd, ok := dict.Word(s)
		if !ok {
			return solveEnterResult{err: "Word not in dictionary"}
		}
		if wd.IsIsland() {
			return solveEnterResult{err: fmt.Sprintf("%q is an island word", wd.String())}
		}
		return solveEnterResult{
			nextStep: solveEndWord,
			update: func(v *viewSolve) {
				v.currentInput = nil
				v.startWord = wd
			},
		}
	}
}

func (v *viewSolve) enterEndWord(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		if len(s) == 0 {
			wm := words.NewWordDistanceMap(v.startWord, nil)
			candidates := wm.Words()
			if len(candidates) == 0 {
				return solveEnterResult{err: "Cannot find random final word (try again)!"}
			}
			word := candidates[rng.Intn(len(candidates))]
			return solveEnterResult{
				nextStep: solveEndWord,
				update: func(v *viewSolve) {
					v.currentInput.set(word)
				},
			}
		} else if len(s) == len(v.startWord.String()) {
			dict := m.loadDictionary(len(s))
			if wd, ok := dict.Word(s); ok {
				wdm := words.NewWordDistanceMap(v.startWord, nil)
				if _, ok = wdm.Distance(wd); ok {
					return solveEnterResult{
						nextStep: solveMaxLadder,
						update: func(v *viewSolve) {
							v.currentInput = nil
							v.endWord = wd
						},
					}
				} else {
					return solveEnterResult{err: fmt.Sprintf("Cannot reach %q from %q", v.endWord, v.startWord)}
				}
			} else {
				return solveEnterResult{err: "Word not in dictionary"}
			}
		} else {
			return solveEnterResult{err: fmt.Sprintf("Please enter a word with %d letters", len(v.startWord.String()))}
		}
	}
}

func (v *viewSolve) enterMaxLadder(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		dict := m.loadDictionary(len(v.startWord.String()))
		if len(s) == 0 {
			start := time.Now()
			wdm := words.NewWordDistanceMap(v.startWord, nil)
			minLengthCalcTime := time.Since(start)
			if dist, ok := wdm.Distance(v.endWord); !ok {
				return solveEnterResult{err: fmt.Sprintf("Cannot reach %q from %q", s, v.startWord)}
			} else {
				solver := solving.NewSolver(solving.NewPuzzle(v.startWord, v.endWord))
				start = time.Now()
				solutions := solver.Solve(dist)
				solveTime := time.Since(start)
				return solveEnterResult{
					nextStep: solveSolved,
					update: func(v *viewSolve) {
						v.dictionaryLoadTime = m.dictionaryLoadTimes[len(v.startWord.String())]
						v.minLengthCalcTime = minLengthCalcTime
						v.minLadderLength = dist
						v.ladderLength = -1
						v.solutions = solutions
						v.solveTime = solveTime
					},
				}
			}
		} else if n, err := strconv.Atoi(s); err == nil && n > 2 && n <= dict.MaxSteps() {
			solver := solving.NewSolver(solving.NewPuzzle(v.startWord, v.endWord))
			start := time.Now()
			solutions := solver.Solve(n)
			solveTime := time.Since(start)
			return solveEnterResult{
				nextStep: solveSolved,
				update: func(v *viewSolve) {
					v.dictionaryLoadTime = m.dictionaryLoadTimes[len(v.startWord.String())]
					v.ladderLength = n
					v.solutions = solutions
					v.solveTime = solveTime
				},
			}
		} else {
			return solveEnterResult{err: fmt.Sprintf("Please enter a number between 2 and %d", dict.MaxSteps())}
		}
	}
}
