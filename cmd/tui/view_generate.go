package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
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

func (v *viewGenerate) render(sf layout.Surface, m *model) *tea.Cursor {
	const (
		promptWordLength   = "Word length:"
		promptLadderLength = "Ladder length:"
		promptStartWord    = "Start word:"
		promptEndWord      = "End word:"
		promptLen          = len(promptLadderLength)
	)
	v.wordsDisplayed = make(wordPoints)
	var csr *tea.Cursor
	switch v.step {
	case generateWordLength:
		sf.TextRight(1, 1, promptLen, promptWordLength)
		if v.currentInput == nil {
			initial := strconv.Itoa(m.prefs.WordLength)
			if v.wasWordLen > 0 {
				initial = strconv.Itoa(v.wasWordLen)
			}
			v.currentInput = &numberInput{maxLength: 2, current: initial}
		}
		sf.TextFixed(1, promptLen+2, 2, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 2)
		if v.currentError != "" {
			sf.Text(1, promptLen+6, v.currentError, errorStyle)
		} else {
			sf.Text(1, promptLen+6, "(enter a number 2-15)", helpStyle)
		}
	case generateLadderLength:
		sf.TextRight(1, 1, promptLen, promptWordLength)
		sf.TextFixed(1, promptLen+2, 2, strconv.Itoa(v.wordLen), inputStyle)
		sf.TextRight(2, 1, promptLen, promptLadderLength)
		if v.currentInput == nil {
			initial := strconv.Itoa(m.prefs.LadderLength)
			if v.wasLadderLen > 0 {
				initial = strconv.Itoa(v.wasLadderLen)
			}
			v.currentInput = &numberInput{maxLength: 2, current: initial}
		}
		sf.TextFixed(2, promptLen+2, 2, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 3)
		if v.currentError != "" {
			sf.Text(2, promptLen+6, v.currentError, errorStyle)
		} else {
			sf.Text(2, promptLen+6, "(enter a number 2-"+strconv.Itoa(m.dictionary.MaxSteps())+")", helpStyle)
		}
	case generateStartWord:
		sf.TextRight(1, 1, promptLen, promptWordLength)
		sf.TextFixed(1, promptLen+2, 2, strconv.Itoa(v.wordLen), inputStyle)
		sf.TextRight(2, 1, promptLen, promptLadderLength)
		sf.TextFixed(2, promptLen+2, 2, strconv.Itoa(v.ladderLen), inputStyle)
		sf.TextRight(3, 1, promptLen, promptStartWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: v.wordLen}
		}
		sf.TextFixed(3, promptLen+2, v.wordLen, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 4)
		if v.currentError != "" {
			sf.Text(3, promptLen+4+v.wordLen, v.currentError, errorStyle)
		} else {
			sf.Text(3, promptLen+4+v.wordLen, "(word, blank or '?' for random)", helpStyle)
		}
	case generateEndWord:
		sf.TextRight(1, 1, promptLen, promptWordLength)
		sf.TextFixed(1, promptLen+2, 2, strconv.Itoa(v.wordLen), inputStyle)
		sf.TextRight(2, 1, promptLen, promptLadderLength)
		sf.TextFixed(2, promptLen+2, 2, strconv.Itoa(v.ladderLen), inputStyle)
		sf.TextRight(3, 1, promptLen, promptStartWord)
		sw := ""
		if v.startWord != nil {
			sw = v.startWord.String()
		}
		v.wordsDisplayed.addWord(sw, 4, promptLen+2)
		sf.TextFixed(3, promptLen+2, v.wordLen, sw, inputStyle)
		sf.TextRight(4, 1, promptLen, promptEndWord)
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: v.wordLen}
		}
		sf.TextFixed(4, promptLen+2, v.wordLen, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(promptLen+v.currentInput.cursorPos()+2, 5)
		if v.currentError != "" {
			sf.Text(4, promptLen+4+v.wordLen, v.currentError, errorStyle)
		} else {
			sf.Text(4, promptLen+4+v.wordLen, "(word, blank or '?' for random)", helpStyle)
		}
	case generateGenerated:
		sf.TextRight(1, 1, promptLen, promptWordLength)
		sf.TextFixed(1, promptLen+2, 2, strconv.Itoa(v.wordLen), inputStyle)
		sf.TextRight(2, 1, promptLen, promptLadderLength)
		sf.TextFixed(2, promptLen+2, 2, strconv.Itoa(v.ladderLen), inputStyle)
		sf.TextRight(3, 1, promptLen, promptStartWord)
		sw := ""
		if v.startWord != nil {
			sw = v.startWord.String()
		} else if v.puzzle != nil {
			sw = v.puzzle.StartWord.String()
		}
		v.wordsDisplayed.addWord(sw, 4, promptLen+2)
		sf.TextFixed(3, promptLen+2, v.wordLen, sw, inputStyle)
		sf.TextRight(4, 1, promptLen, promptEndWord)
		ew := ""
		if v.endWord != nil {
			ew = v.endWord.String()
		} else if v.puzzle != nil {
			ew = v.puzzle.EndWord.String()
		}
		v.wordsDisplayed.addWord(ew, 5, promptLen+2)
		sf.TextFixed(4, promptLen+2, v.wordLen, ew, inputStyle)
		if v.currentError != "" {
			sf.Text(6, 1, v.currentError, errorStyle)
		} else {
			sf.TextRun(6, 2, layout.NewRuns("Took ").Add(truncateDuration(v.puzzleGenerateTime), highlightStyle).Add(" to generate puzzle"))
			sf.TextRun(7, 2, layout.NewRuns("Max score: ").
				Add(strconv.FormatFloat(v.puzzle.MaxScore, 'f', 0, 64), highlightStyle).
				Add(" (").Add(commas(len(v.puzzle.Solutions)), highlightStyle).Add(" solutions)"))
		}
	}
	return csr
}

func (v *viewGenerate) helpLines() ([]string, *lipgloss.Style) {
	if v.step == generateGenerated {
		return []string{
			ctrlPlay + ": Play  •  " + enter + ": Solutions",
			ctrlNew + ": New  •  " + ctrlSolver + ": Solver",
		}, nil
	} else {
		return []string{ctrlNew + ": New  •  " + ctrlSolver + ": Solver"}, nil
	}
}

func (v *viewGenerate) menu() []menuItem {
	if v.step == generateGenerated {
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
