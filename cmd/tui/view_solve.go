package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/generator"
	"github.com/marrow16/gowordladder/solving"
	"github.com/marrow16/gowordladder/words"
	"strconv"
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
	step               solveStep
	startWord          *words.Word
	endWord            *words.Word
	ladderLength       int
	currentInput       input
	currentError       string
	dictionaryLoadTime time.Duration
	minLadderLength    int
	minLengthCalcTime  time.Duration
	solutions          []*solving.Solution
	solveTime          time.Duration
	explored           int
	wordsDisplayed     wordPoints
}

func (v *viewSolve) wordLength() int {
	if v.step >= solveStartWord && v.startWord != nil {
		return len(v.startWord.String())
	}
	return 0
}

func (v *viewSolve) currentWord() string {
	if (v.step == solveStartWord || v.step == solveEndWord) && v.currentInput != nil {
		return v.currentInput.value()
	} else if v.step == solveSolved {
		return v.startWord.String()
	}
	return ""
}

func (v *viewSolve) render(sf layout.Surface, m *model) *tea.Cursor {
	const (
		promptStartWord = "Start word:"
		promptEndWord   = "End word:"
		promptMaxLadder = "Maximum ladder length:"
		promptLen       = len(promptMaxLadder)
	)
	v.wordsDisplayed = make(wordPoints)
	var csr *tea.Cursor
	switch v.step {
	case solveStartWord:
		sf.TextRight(1, 1, promptLen, promptStartWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: 15}
		}
		sf.TextFixed(1, promptLen+2, 15, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 2)
		if v.currentError != "" {
			sf.Text(1, promptLen+15+4, v.currentError, errorStyle)
		}
	case solveEndWord:
		wl := len(v.startWord.String())
		sf.TextRight(1, 1, promptLen, promptStartWord)
		v.wordsDisplayed.add(sf.TextFixed(1, promptLen+2, len(v.startWord.String()), v.startWord.String(), inputStyle))
		sf.TextRight(2, 1, promptLen, promptEndWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: wl}
		}
		sf.TextFixed(2, promptLen+2, wl, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 3)
		if v.currentError != "" {
			sf.Text(2, promptLen+2+wl+2, v.currentError, errorStyle)
		} else {
			sf.Text(2, promptLen+2+wl+2, "(blank or '?' for random)", helpStyle)
		}
	case solveMaxLadder:
		wl := len(v.startWord.String())
		sf.TextRight(1, 1, promptLen, promptStartWord)
		v.wordsDisplayed.add(sf.TextFixed(1, promptLen+2, wl, v.startWord.String(), inputStyle))
		sf.TextRight(2, 1, promptLen, promptEndWord)
		v.wordsDisplayed.add(sf.TextFixed(2, promptLen+2, wl, v.endWord.String(), inputStyle))
		sf.TextRight(3, 1, promptLen, promptMaxLadder)
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2}
		}
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 4)
		sf.TextFixed(3, promptLen+2, 2, v.currentInput.value(), inputStyle)
		if v.currentError != "" {
			sf.Text(3, promptLen+6, v.currentError, errorStyle)
		} else {
			sf.Text(3, promptLen+6, "(optional - blank for auto min)", helpStyle)
		}
	case solveSolved:
		wl := len(v.startWord.String())
		sf.TextRight(1, 1, promptLen, promptStartWord)
		v.wordsDisplayed.add(sf.TextFixed(1, promptLen+2, wl, v.startWord.String(), inputStyle))
		sf.TextRight(2, 1, promptLen, promptEndWord)
		v.wordsDisplayed.add(sf.TextFixed(2, promptLen+2, wl, v.endWord.String(), inputStyle))
		sf.TextRight(3, 1, promptLen, promptMaxLadder)
		if v.ladderLength == -1 {
			sf.TextFixed(3, promptLen+2, 2, "??", inputStyle)
		} else {
			sf.TextFixed(3, promptLen+2, 2, strconv.Itoa(v.ladderLength), inputStyle)
		}
		sf.TextRun(5, 1, layout.NewRuns("Took ").Add(truncateDuration(v.dictionaryLoadTime), highlightStyle).Add(" to load dictionary"))
		row := 6
		if v.ladderLength == -1 {
			sf.TextRun(row, 1, layout.NewRuns("Took ").Add(truncateDuration(v.minLengthCalcTime), highlightStyle).Add(" to determine minimum ladder length ").Add(strconv.Itoa(v.minLadderLength), highlightStyle))
			row++
		}
		sf.TextRun(row, 1, layout.NewRuns("Took ").
			Add(truncateDuration(v.solveTime), highlightStyle).
			Add(" to find ").
			Add(commas(len(v.solutions)), highlightStyle).Add(" solutions"))
	}
	return csr
}

func (v *viewSolve) helpLines() ([]string, *lipgloss.Style) {
	if v.step == solveSolved && len(v.solutions) > 0 {
		return []string{
			ctrlPlay + ": Play  •  " + enter + ": Solutions",
			ctrlNew + ": New"}, nil
	} else {
		return []string{ctrlNew + ": New"}, nil
	}
}

func (v *viewSolve) menu() []menuItem {
	if v.step == solveSolved {
		return []menuItem{
			{text: "Play", key: ctrlPlay},
			{text: "Solutions", key: enter},
			{text: "New", key: ctrlNew},
		}
	}
	return []menuItem{
		{text: "New", key: ctrlNew},
	}
}

func (v *viewSolve) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.currentError = ""
	switch msg.String() {
	case ctrlNew:
		v.currentInput = nil
		v.currentError = ""
		v.solutions = nil
		v.startWord = nil
		v.endWord = nil
		v.step = solveStartWord
	case ctrlPlay:
		if v.step == solveSolved && len(v.solutions) > 0 {
			sw, ew := v.startWord.String(), v.endWord.String()
			ll := len(v.solutions[0].Ladder())
			if puzzle, err := generator.GeneratePuzzle(len(sw), ll, &sw, &ew); err == nil {
				m.play(*puzzle)
			}
		}
	case enter:
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
	case "?":
		if v.step == solveEndWord && v.currentInput != nil {
			v.currentInput.set("")
			return v.enterEndWord(m)
		}
	}
	if v.currentInput != nil {
		v.currentInput.key(msg)
	}
	return nil
}

func (v *viewSolve) click(m *model, msg tea.Mouse) tea.Cmd {
	if wd, ok := v.wordsDisplayed[pt{msg.Y, msg.X}]; ok {
		switch msg.Button {
		case tea.MouseLeft:
			return m.lookupWord(wd)
		case tea.MouseRight:
			return m.lookupDistance(wd)
		}
	}
	return nil
}

func (v *viewSolve) paste(m *model, msg tea.PasteMsg) {
	if v.currentInput != nil {
		v.currentInput.paste(msg)
	}
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
					return solveEnterResult{err: fmt.Sprintf("Cannot reach %q from %q", s, v.startWord)}
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
						v.explored = solver.ExploredCount()
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
					v.explored = solver.ExploredCount()
				},
			}
		} else {
			return solveEnterResult{err: fmt.Sprintf("Please enter a number between 2 and %d", dict.MaxSteps())}
		}
	}
}
