package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"gowordladder/generator"
	"gowordladder/words"
	"strconv"
	"strings"
	"time"
)

type generateStep int

const (
	generateWordLength generateStep = iota
	generateLadderLength
	generateStartWord
	generateEndWord
	generateGenerated
)

type viewGenerate struct {
	step      generateStep
	wordLen   int
	ladderLen int
	startWord *words.Word
	endWord   *words.Word

	currentInput input
	currentError string

	puzzle             *generator.Puzzle
	puzzleGenerateTime time.Duration
}

func (v *viewGenerate) wordLength() int {
	if v.step > generateWordLength {
		return v.wordLen
	}
	return 0
}

func (v *viewGenerate) content(m *model) (string, *tea.Cursor) {
	const (
		promptWordLength   = "   Word length: "
		promptLadderLength = " Ladder length: "
		promptStartWord    = "    Start word: "
		promptEndWord      = "      End word: "
		promptLen          = len(promptEndWord)
	)
	var sb strings.Builder
	sb.WriteString("\n")
	lines := 1
	cpx := -1
	var s string
	switch v.step {
	case generateWordLength:
		sb.WriteString(promptWordLength)
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		}
		lines++
	case generateLadderLength:
		sb.WriteString(promptWordLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.wordLen)))
		sb.WriteString("\n" + promptLadderLength)
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		} else {
			sb.WriteString(helpStyle.Render(fmt.Sprintf("  (enter a number 3-%d)", m.dictionary.MaxSteps())))
		}
		lines += 2
	case generateStartWord:
		sb.WriteString(promptWordLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.wordLen)))
		sb.WriteString("\n" + promptLadderLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.ladderLen)))
		sb.WriteString("\n" + promptStartWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: v.wordLen}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		} else {
			sb.WriteString(helpStyle.Render("  (word, blank or '?' for random)"))
		}
		lines += 3
	case generateEndWord:
		sb.WriteString(promptWordLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.wordLen)))
		sb.WriteString("\n" + promptLadderLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.ladderLen)))
		sb.WriteString("\n" + promptStartWord)
		if v.startWord != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.startWord.String()))
		} else {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(""))
		}
		sb.WriteString("\n" + promptEndWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: v.wordLen}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		} else {
			sb.WriteString(helpStyle.Render("  (word, blank or '?' for random)"))
		}
		lines += 4
	case generateGenerated:
		sb.WriteString(promptWordLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.wordLen)))
		sb.WriteString("\n" + promptLadderLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.ladderLen)))
		sb.WriteString("\n" + promptStartWord)
		if v.startWord != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.startWord.String()))
		} else if v.puzzle != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.puzzle.StartWord.String()))
		} else {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(""))
		}
		sb.WriteString("\n" + promptEndWord)
		if v.endWord != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.endWord.String()))
		} else if v.puzzle != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.puzzle.EndWord.String()))
		} else {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(""))
		}
		lines += 4
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("\n\n  " + v.currentError))
			lines += 2
		} else {
			sb.WriteString("\n\n  Took " + highlightStyle.Render(v.puzzleGenerateTime.String()) + " to generate puzzle")
			sb.WriteString("\n  Max score: " + highlightStyle.Render(fmt.Sprintf("%.0f", v.puzzle.MaxScore)))
			sb.WriteString(", " + highlightStyle.Render(fmt.Sprintf("%d", len(v.puzzle.Solutions))) + " solutions")
			lines += 3
		}
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

func (v *viewGenerate) help() string {
	if v.step == generateGenerated {
		return "ctrl+n: New  •  ctrl+p: Play  •  enter: Solutions  •  ctrl+s: Solver"
	} else {
		return "ctrl+n: New  •  ctrl+s: Solver"
	}
}

func (v *viewGenerate) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.currentError = ""
	switch msg.String() {
	case "ctrl+n":
		v.currentInput = nil
		v.currentError = ""
		v.wordLen = 0
		v.ladderLen = 0
		v.startWord = nil
		v.endWord = nil
		v.puzzle = nil
		v.step = generateWordLength
	case "ctrl+p":
		if v.puzzle != nil {
			m.play(*v.puzzle)
		}
	case "enter":
		switch v.step {
		case generateWordLength:
			return v.enterWordLength(m)
		case generateLadderLength:
			return v.enterLadderLength(m)
		case generateStartWord:
			return v.enterStartWord(m)
		case generateEndWord:
			return v.enterEndWord(m)
		case generateGenerated:
			if v.puzzle != nil && len(v.puzzle.Solutions) > 0 {
				m.showSolutions(v.puzzle.Solutions)
				return nil
			}
		}
	case "?":
		switch v.step {
		case generateStartWord:
			return v.randomStartWord(m)
		case generateEndWord:
			return v.randomEndWord(m)
		}
	}
	if v.currentInput != nil {
		v.currentInput.key(msg)
	}
	return nil
}

type generateEnterResult struct {
	err      string
	nextStep generateStep
	update   func(v *viewGenerate)
}

func (v *viewGenerate) update(m *model, msg tea.Msg) tea.Cmd {
	if result, ok := msg.(generateEnterResult); ok {
		if result.err != "" {
			v.currentError = result.err
		} else {
			if result.update != nil {
				result.update(v)
			}
			v.step = result.nextStep
		}
	}
	return nil
}

func (v *viewGenerate) enterWordLength(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		if n, err := strconv.Atoi(s); err == nil && n >= 2 && n <= 15 {
			return generateEnterResult{
				nextStep: generateLadderLength,
				update: func(v *viewGenerate) {
					m.loadDictionary(n)
					v.wordLen = n
					v.currentInput = nil
				},
			}
		} else {
			return generateEnterResult{err: "Please enter a number 2-15"}
		}
	}
}

func (v *viewGenerate) enterLadderLength(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		if n, err := strconv.Atoi(s); err == nil && n >= 3 && n <= m.dictionary.MaxSteps() {
			return generateEnterResult{
				nextStep: generateStartWord,
				update: func(v *viewGenerate) {
					v.ladderLen = n
					v.currentInput = nil
				},
			}
		} else {
			return generateEnterResult{err: fmt.Sprintf("Please enter a number 3-%d", m.dictionary.MaxSteps())}
		}
	}
}

func (v *viewGenerate) enterStartWord(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		if len(s) == 0 {
			start := time.Now()
			puzzle, err := generator.GeneratePuzzle(v.wordLen, v.ladderLen, nil, nil)
			dur := time.Since(start)
			if err != nil {
				return generateEnterResult{err: err.Error()}
			}
			return generateEnterResult{
				nextStep: generateGenerated,
				update: func(v *viewGenerate) {
					v.currentInput = nil
					v.puzzle = puzzle
					v.puzzleGenerateTime = dur
				},
			}
		} else if len(s) != v.wordLen {
			return generateEnterResult{err: fmt.Sprintf("Please enter a word with %d letters", v.wordLen)}
		}
		dict := m.loadDictionary(v.wordLen)
		if wd, ok := dict.Word(s); ok {
			return generateEnterResult{
				nextStep: generateEndWord,
				update: func(v *viewGenerate) {
					v.startWord = wd
					if wd.MaxSteps() < v.ladderLen {
						v.ladderLen = wd.MaxSteps()
					}
					v.currentInput = nil
				},
			}
		} else {
			return generateEnterResult{err: "Word not in dictionary"}
		}
	}
}

func (v *viewGenerate) enterEndWord(m *model) tea.Cmd {
	s := v.currentInput.value()
	v.currentError = ""
	return func() tea.Msg {
		if len(s) == 0 {
			start := time.Now()
			sw := v.startWord.String()
			puzzle, err := generator.GeneratePuzzle(v.wordLen, v.ladderLen, &sw, nil)
			dur := time.Since(start)
			if err != nil {
				return generateEnterResult{err: err.Error()}
			}
			return generateEnterResult{
				nextStep: generateGenerated,
				update: func(v *viewGenerate) {
					v.currentInput = nil
					v.puzzle = puzzle
					v.puzzleGenerateTime = dur
				},
			}
		} else if len(s) != v.wordLen {
			return generateEnterResult{err: fmt.Sprintf("Please enter a word with %d letters", v.wordLen)}
		}
		dict := m.loadDictionary(v.wordLen)
		if wd, ok := dict.Word(s); ok {
			wdm := words.NewWordDistanceMap(v.startWord, &v.ladderLen)
			dist, ok := wdm.Distance(wd)
			if !ok {
				return generateEnterResult{err: fmt.Sprintf("Cannot reach %q from %q", wd.String(), v.startWord.String())}
			}
			_ = dist
			start := time.Now()
			sw := v.startWord.String()
			puzzle, err := generator.GeneratePuzzle(v.wordLen, v.ladderLen, &sw, &s)
			dur := time.Since(start)
			if err != nil {
				return generateEnterResult{err: err.Error()}
			}
			return generateEnterResult{
				nextStep: generateGenerated,
				update: func(v *viewGenerate) {
					v.currentInput = nil
					v.puzzle = puzzle
					v.puzzleGenerateTime = dur
				},
			}
		} else {
			return generateEnterResult{err: "Word not in dictionary"}
		}
	}
}

func (v *viewGenerate) randomStartWord(m *model) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		dict := m.loadDictionary(v.wordLen)
		candidates := dict.WordsWithSteps(v.ladderLen)
		word := candidates[rng.Intn(len(candidates))]
		return generateEnterResult{
			nextStep: generateStartWord,
			update: func(v *viewGenerate) {
				v.currentInput.set(strings.ToUpper(word.String()))
			},
		}
	}
}

func (v *viewGenerate) randomEndWord(m *model) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		wdm := words.NewWordDistanceMap(v.startWord, &v.ladderLen)
		candidates := wdm.WordsAt(v.ladderLen)
		if len(candidates) == 0 {
			return generateEnterResult{err: "Unable to generate random word (try again)"}
		} else {
			word := candidates[rng.Intn(len(candidates))]
			return generateEnterResult{
				nextStep: generateEndWord,
				update: func(v *viewGenerate) {
					v.currentInput.set(strings.ToUpper(word))
				},
			}
		}
	}
}
