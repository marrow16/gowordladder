package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/marrow16/gowordladder/generator"
	"github.com/marrow16/gowordladder/words"
	"slices"
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
	step                     generateStep
	wordLen                  int
	ladderLen                int
	startWord                *words.Word
	endWord                  *words.Word
	currentInput             input
	currentError             string
	wasWordLen, wasLadderLen int
	puzzle                   *generator.Puzzle
	puzzleGenerateTime       time.Duration
	wordsDisplayed           wordPoints
}

func (v *viewGenerate) wordLength() int {
	if v.step > generateWordLength {
		return v.wordLen
	}
	return 0
}

func (v *viewGenerate) currentWord() string {
	if (v.step == generateStartWord || v.step == generateEndWord) && v.currentInput != nil {
		return v.currentInput.value()
	} else if v.puzzle != nil {
		return v.puzzle.StartWord.String()
	}
	return ""
}

func (v *viewGenerate) content(m *model) (string, *tea.Cursor) {
	const (
		promptWordLength   = "   Word length: "
		promptLadderLength = " Ladder length: "
		promptStartWord    = "    Start word: "
		promptEndWord      = "      End word: "
		promptLen          = len(promptEndWord)
		footerLines        = 3
	)
	v.wordsDisplayed = make(wordPoints)
	var sb strings.Builder
	sb.Grow(m.height * m.width)
	sb.WriteString("\n")
	lines := 1
	cpx := -1
	var s string
	switch v.step {
	case generateWordLength:
		sb.WriteString(promptWordLength)
		if v.currentInput == nil {
			initial := strconv.Itoa(m.prefs.WordLength)
			if v.wasWordLen > 0 {
				initial = strconv.Itoa(v.wasWordLen)
			}
			v.currentInput = &numberInput{maxLength: 2, current: initial}
		}
		s, cpx = v.currentInput.render()
		sb.WriteString(s)
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("  " + v.currentError))
		} else {
			sb.WriteString(helpStyle.Render("  (enter a number 2-15)"))
		}
		lines++
	case generateLadderLength:
		sb.WriteString(promptWordLength)
		sb.WriteString(inputStyle.Width(2).Render(fmt.Sprintf("%2d", v.wordLen)))
		sb.WriteString("\n" + promptLadderLength)
		if v.currentInput == nil {
			initial := strconv.Itoa(m.prefs.LadderLength)
			if v.wasLadderLen > 0 {
				initial = strconv.Itoa(v.wasLadderLen)
			}
			v.currentInput = &numberInput{maxLength: 2, current: initial}
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
			v.wordsDisplayed.addWord(v.startWord.String(), 4, promptLen)
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
			v.wordsDisplayed.addWord(v.startWord.String(), 4, promptLen)
		} else if v.puzzle != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.puzzle.StartWord.String()))
			v.wordsDisplayed.addWord(v.puzzle.StartWord.String(), 4, promptLen)
		} else {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(""))
		}
		sb.WriteString("\n" + promptEndWord)
		if v.endWord != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.endWord.String()))
			v.wordsDisplayed.addWord(v.endWord.String(), 5, promptLen)
		} else if v.puzzle != nil {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(v.puzzle.EndWord.String()))
			v.wordsDisplayed.addWord(v.puzzle.EndWord.String(), 5, promptLen)
		} else {
			sb.WriteString(inputStyle.Width(v.wordLen).Render(""))
		}
		lines += 4
		if v.currentError != "" {
			sb.WriteString(errorStyle.Render("\n\n  " + v.currentError))
			lines += 2
		} else {
			sb.WriteString("\n\n  Took " + highlightStyle.Render(truncateDuration(v.puzzleGenerateTime)) + " to generate puzzle")
			sb.WriteString("\n  Max score: " + highlightStyle.Render(fmt.Sprintf("%.0f", v.puzzle.MaxScore)))
			sb.WriteString(" (" + highlightStyle.Render(commas(len(v.puzzle.Solutions))) + " solutions)")
			lines += 3
		}
	}

	sb.WriteString(padLines(m.height - lines - footerLines))
	var csr *tea.Cursor
	if cpx > -1 {
		csr = tea.NewCursor(promptLen+cpx, lines)
	}
	return sb.String(), csr
}

func (v *viewGenerate) help() string {
	if v.step == generateGenerated {
		return ctrlPlay + ": Play  •  enter: Solutions\n" + ctrlNew + ": New  •  " + ctrlSolver + ": Solver"
	} else {
		return "\n" + ctrlNew + ": New  •  " + ctrlSolver + ": Solver"
	}
}

func (v *viewGenerate) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.currentError = ""
	switch msg.String() {
	case ctrlNew:
		v.reset(m)
	case ctrlPlay:
		if v.puzzle != nil {
			m.play(*v.puzzle)
		}
	case enter:
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
	case up:
		switch v.step {
		case generateStartWord:
			return v.scrollStartWord(m, -1)
		case generateEndWord:
			return v.scrollEndWord(m, -1)
		}
	case down:
		switch v.step {
		case generateStartWord:
			return v.scrollStartWord(m, 1)
		case generateEndWord:
			return v.scrollEndWord(m, 1)
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

func (v *viewGenerate) click(m *model, msg tea.Mouse) tea.Cmd {
	if wd, ok := v.wordsDisplayed[pt{msg.Y, msg.X}]; ok {
		return m.lookupWord(wd)
	}
	return nil
}

func (v *viewGenerate) paste(m *model, msg tea.PasteMsg) {
	if v.currentInput != nil {
		v.currentInput.paste(msg)
	}
}

func (v *viewGenerate) reset(m *model) {
	v.currentInput = nil
	v.currentError = ""
	v.wasWordLen, v.wasLadderLen = v.wordLen, v.ladderLen
	v.wordLen = 0
	v.ladderLen = 0
	v.startWord = nil
	v.endWord = nil
	v.puzzle = nil
	v.step = generateWordLength
	v.currentInput = nil
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

func findWordIndex(words []*words.Word, word string) int {
	for i, w := range words {
		if w.String() == word {
			return i
		}
	}
	return -1
}

func (v *viewGenerate) scrollStartWord(m *model, scroll int) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		dict := m.loadDictionary(v.wordLen)
		candidates := dict.WordsWithSteps(v.ladderLen)
		var word *words.Word
		if cw := v.currentInput.value(); len(cw) != v.wordLen {
			word = candidates[0]
		} else {
			idx := findWordIndex(candidates, cw) + scroll
			if idx < 0 {
				idx = len(candidates) - 1
			} else if idx >= len(candidates) {
				idx = 0
			}
			word = candidates[idx]
		}
		return generateEnterResult{
			nextStep: generateStartWord,
			update: func(v *viewGenerate) {
				v.currentInput.set(strings.ToUpper(word.String()))
			},
		}
	}
}

func (v *viewGenerate) scrollEndWord(m *model, scroll int) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		wdm := words.NewWordDistanceMap(v.startWord, &v.ladderLen)
		candidates := wdm.WordsAt(v.ladderLen)
		if len(candidates) == 0 {
			return generateEnterResult{err: "Unable to generate random word (try again)"}
		} else {
			word := ""
			if cw := v.currentInput.value(); len(cw) != v.wordLen {
				word = candidates[0]
			} else {
				idx := slices.Index(candidates, v.currentInput.value()) + scroll
				if idx < 0 {
					idx = len(candidates) - 1
				} else if idx >= len(candidates) {
					idx = 0
				}
				word = candidates[idx]
			}
			return generateEnterResult{
				nextStep: generateEndWord,
				update: func(v *viewGenerate) {
					v.currentInput.set(strings.ToUpper(word))
				},
			}
		}
	}
}
