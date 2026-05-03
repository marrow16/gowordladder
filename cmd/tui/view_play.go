package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"gowordladder/generator"
	"gowordladder/words"
	"math"
	"slices"
	"strconv"
	"strings"
)

type playView interface {
	view
	start(puzzle generator.Puzzle, dict *words.Dictionary)
}

type viewPlay struct {
	offsetY        int
	onStep, onChar int
	entries        []string
	dictionary     *words.Dictionary
	puzzle         generator.Puzzle
	currentScore   float64
	hint           string
	hintsGiven     map[string]bool
	warning        string
	wrong          string
	okWords        map[int]bool
	solved         bool
}

func (v *viewPlay) content(m *model) (string, *tea.Cursor) {
	const (
		topLeft     = "╭"
		topRight    = "╮"
		topBottom   = "─"
		bottomLeft  = "╰"
		bottomRight = "╯"
		vertical    = "│"
	)
	var sb strings.Builder
	sols := fmt.Sprintf(" Solutions: %d ", len(v.puzzle.Solutions))
	currScore := fmt.Sprintf(" Current score: %.0f ", v.currentScore)
	maxScore := fmt.Sprintf(" Maximum score: %.0f ", v.puzzle.MaxScore)
	midPadL := ""
	midPadR := ""
	if p := (m.width / 2) - (len(currScore) / 2) - len(sols); p > 0 {
		midPadL = strings.Repeat(" ", p)
	}
	if p := (m.width / 2) - (len(currScore) / 2) - len(maxScore); p > 0 {
		midPadR = strings.Repeat(" ", p)
	}
	if len(sols)+len(midPadL)+len(currScore)+len(midPadR)+len(maxScore) < m.width {
		midPadR += " "
	}
	sb.WriteString(headerStyle.Render(sols + midPadL + currScore + midPadR + maxScore))
	lines := 2

	var csr *tea.Cursor
	maxLines := m.height - lines - 3
	ladderWd := v.puzzle.WordLength + 2
	padL := strings.Repeat(" ", ((m.width-ladderWd)/2)-3)
	stop := false
	for l := 0; !stop && l < maxLines; l++ {
		sb.WriteString("\n")
		lines++
		rung := l + v.offsetY - 2
		sb.WriteString(padL)
		switch {
		case rung == -2:
			sb.WriteString(helpStyle.Render("   " + topLeft + strings.Repeat(topBottom, v.puzzle.WordLength) + topRight))
		case rung == -1:
			sb.WriteString(helpStyle.Render(" 1 " + vertical))
			if v.solved {
				sb.WriteString(highlightStyle.Render(v.puzzle.StartWord.String()))
			} else {
				sb.WriteString(v.puzzle.StartWord.String())
			}
			sb.WriteString(helpStyle.Render(vertical))
		case rung == v.puzzle.LadderLength-2:
			sb.WriteString(helpStyle.Render(fmt.Sprintf("%2d ", v.puzzle.LadderLength)))
			sb.WriteString(helpStyle.Render(vertical))
			if v.solved {
				sb.WriteString(highlightStyle.Render(v.puzzle.EndWord.String()))
			} else {
				sb.WriteString(v.puzzle.EndWord.String())
			}
			sb.WriteString(helpStyle.Render(vertical))
		case rung == v.puzzle.LadderLength-1:
			sb.WriteString(helpStyle.Render("   " + bottomLeft + strings.Repeat(topBottom, v.puzzle.WordLength) + bottomRight))
		case rung < v.puzzle.LadderLength:
			if rung == v.onStep {
				csr = tea.NewCursor(len(padL)+4+v.onChar, lines-1)
				csr.Color = lipgloss.Color("#ccccff")
			}
			sb.WriteString(helpStyle.Render(fmt.Sprintf("%2d ", rung+2)))
			sb.WriteString(helpStyle.Render(vertical))
			if v.solved || v.okWords[rung] {
				sb.WriteString(highlightStyle.Render(v.entries[rung]))
			} else {
				sb.WriteString(v.entries[rung])
			}
			sb.WriteString(helpStyle.Render(vertical))
		default:
			stop = true
		}
	}

	if m.height-lines-2 > 0 {
		sb.WriteString(strings.Repeat("\n", m.height-lines-2))
	}
	return sb.String(), csr
}

func (v *viewPlay) help() string {
	const (
		firstHelp   = "ctrl+h: Solutions  •  ?: Hint  •  ctrl+f: Fill  •  space: Clear"
		secondHelp  = "\nctrl+n: New  •  ctrl+g: Generate  •  ctrl+s: Solver"
		solvedHelp  = "\nctrl+n: New  •  ctrl+h: Solutions  •  ctrl+g: Generate"
		defaultHelp = firstHelp + secondHelp
	)
	switch {
	case v.solved:
		return hintStyle.Render("You solved it!") + solvedHelp
	case v.hint != "":
		return hintStyle.Render(v.hint) + secondHelp
	case v.warning != "":
		return warningStyle.Render(v.warning) + secondHelp
	case v.wrong != "":
		return wrongStyle.Render(v.wrong) + secondHelp
	default:
		return defaultHelp
	}
}

var (
	hintStyle    = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#008000"))
	warningStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#ff7f00"))
	wrongStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#ff0000"))
)

func (v *viewPlay) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.hint = ""
	v.warning = ""
	v.wrong = ""
	switch msg.String() {
	case "ctrl+n":
		return v.generateNew()
	case "ctrl+h":
		if !v.solved {
			v.hintDeduction(solutionsPeek)
		}
		m.showSolutions(v.puzzle.Solutions)
	case "up":
		if v.onStep > 0 {
			v.onStep--
			v.ensureCursorVisible(m)
		}
	case "down":
		if v.onStep < len(v.entries)-1 {
			v.onStep++
			v.ensureCursorVisible(m)
		}
	case "enter":
		if !v.solved {
			v.checkWord()
			if v.wrong == "" && v.warning == "" && v.onStep < len(v.entries)-1 {
				v.onChar = 0
				v.onStep++
				v.ensureCursorVisible(m)
			}
		} else if v.onStep < len(v.entries)-1 {
			v.onChar = 0
			v.onStep++
			v.ensureCursorVisible(m)
		}
	case "left":
		if v.onChar > 0 {
			v.onChar--
		} else if v.onStep > 0 {
			v.onStep--
			v.onChar = v.puzzle.WordLength - 1
			v.ensureCursorVisible(m)
		}
	case "right":
		if v.onChar < v.puzzle.WordLength-1 {
			v.onChar++
		} else if v.onStep < len(v.entries)-1 {
			v.onStep++
			v.onChar = 0
			v.ensureCursorVisible(m)
		}
	case "backspace":
		if !v.solved {
			delete(v.okWords, v.onStep)
			s := v.entries[v.onStep]
			s = s[:v.onChar] + " " + s[v.onChar+1:]
			v.entries[v.onStep] = s
			if v.onChar > 0 {
				v.onChar--
			}
		}
	case "space":
		if !v.solved {
			delete(v.okWords, v.onStep)
			v.entries[v.onStep] = strings.Repeat("_", v.puzzle.WordLength)
			v.onChar = 0
		}
	case "ctrl+f":
		if !v.solved {
			delete(v.okWords, v.onStep)
			v.fillWord()
		}
	case "?":
		if !v.solved {
			v.hintWord()
		}
	default:
		if !v.solved {
			delete(v.okWords, v.onStep)
			if k := strings.ToUpper(msg.String()); len(k) == 1 && (k == "." || k == "_" || k == "-" || (k >= "A" && k <= "Z")) {
				s := v.entries[v.onStep]
				s = s[:v.onChar] + k + s[v.onChar+1:]
				v.entries[v.onStep] = s
				if v.onChar < v.puzzle.WordLength-1 {
					v.onChar++
				}
				v.checkWord()
			}
		}
	}
	return nil
}

func (v *viewPlay) ensureCursorVisible(m *model) {
	const (
		headerLines = 2
		footerLines = 3
	)
	maxLines := m.height - headerLines - footerLines
	visibleL := v.onStep - v.offsetY + headerLines
	if visibleL < 0 {
		v.offsetY = v.onStep + headerLines
		// if we're close enough to the top, snap so that first word is visible...
		if v.offsetY < 3 {
			v.offsetY = 0
		}
	} else if visibleL >= maxLines {
		v.offsetY = v.onStep + headerLines - maxLines + 1
		lastRung := v.puzzle.LadderLength - 1
		maxOffsetY := lastRung + footerLines - maxLines
		// if we're close enough to the bottom, snap so the final word is visible...
		if maxOffsetY-v.offsetY < footerLines {
			v.offsetY = maxOffsetY
		}
	}
	if v.offsetY < 0 {
		v.offsetY = 0
	}
}

type newResult struct {
	puzzle *generator.Puzzle
	err    error
}

func (v *viewPlay) update(m *model, msg tea.Msg) tea.Cmd {
	switch mt := msg.(type) {
	case newResult:
		if mt.err == nil {
			v.start(*mt.puzzle, nil)
		}
	}
	return nil
}

func (v *viewPlay) wordLength() int {
	return v.puzzle.WordLength
}

func (v *viewPlay) checkWord() {
	s := v.entries[v.onStep]
	if isAllAZ(s) {
		wd, ok := v.dictionary.Word(s)
		if !ok {
			v.wrong = "That word is not in my dictionary!"
			return
		}
		if prevWord := v.previousWord(true); prevWord != nil {
			if wd.Differences(prevWord) != 1 {
				v.wrong = "Not one letter different to previous word"
				return
			}
		}
		if nextWord := v.nextWord(true); nextWord != nil {
			if wd.Differences(nextWord) != 1 {
				v.wrong = "Not one letter different to next word"
				return
			}
		}
		// check if word is any of the solutions...
		ok = false
		for _, solution := range v.puzzle.Solutions {
			if solution.Ladder()[v.onStep+1].String() == s {
				ok = true
				break
			}
		}
		if !ok {
			delete(v.okWords, v.onStep)
			v.warning = "That word is not in any solution!"
			return
		} else {
			v.okWords[v.onStep] = true
		}
		// see if all words have been filled...
		allWords := true
		for _, e := range v.entries {
			if !isAllAZ(e) {
				allWords = false
				break
			}
		}
		if allWords {
			// all words are filled - check if they match any solution...
			expectCount := len(v.entries)
			for _, solution := range v.puzzle.Solutions {
				count := 0
				ladder := solution.Ladder()
				for i, e := range v.entries {
					if ladder[i+1].String() == e {
						count++
					} else {
						break
					}
				}
				if count == expectCount {
					v.solved = true
					return
				}
			}
		}
	}
}

func (v *viewPlay) previousWord(incStart bool) *words.Word {
	switch {
	case v.onStep == 0 && incStart:
		return v.puzzle.StartWord
	case v.onStep > 0:
		if wd, ok := v.dictionary.Word(v.entries[v.onStep-1]); ok {
			return wd
		}
	}
	return nil
}

func (v *viewPlay) nextWord(incEnd bool) *words.Word {
	switch {
	case v.onStep < len(v.entries)-1:
		if wd, ok := v.dictionary.Word(v.entries[v.onStep+1]); ok {
			return wd
		}
	case incEnd:
		return v.puzzle.EndWord
	}
	return nil
}

func isAllAZ(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func (v *viewPlay) fillWord() {
	prevWord := v.previousWord(true)
	nextWord := v.nextWord(true)
	switch {
	case prevWord == nil && nextWord == nil:
		wd := v.puzzle.Solutions[0].Ladder()[v.onStep+1]
		v.entries[v.onStep] = wd.String()
		v.onChar = 0
		v.checkWord()
	case prevWord != nil && nextWord != nil:
		var wd *words.Word
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			if ladder[v.onStep].String() == prevWord.String() && ladder[v.onStep+2].String() == nextWord.String() {
				wd = ladder[v.onStep+1]
			}
		}
		if wd != nil {
			v.entries[v.onStep] = wd.String()
			v.onChar = 0
			v.checkWord()
		} else {
			v.warning = "Sorry, no words fit here (mistake above/below?)"
			return
		}
	case prevWord != nil:
		var wd *words.Word
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			if ladder[v.onStep].String() == prevWord.String() {
				wd = ladder[v.onStep+1]
			}
		}
		if wd != nil {
			v.entries[v.onStep] = wd.String()
			v.onChar = 0
			v.checkWord()
		} else {
			v.warning = "Sorry, no words fit here (mistake above?)"
			return
		}
	default:
		var wd *words.Word
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			if ladder[v.onStep+2].String() == nextWord.String() {
				wd = ladder[v.onStep+1]
			}
		}
		if wd != nil {
			v.entries[v.onStep] = wd.String()
			v.onChar = 0
			v.checkWord()
		} else {
			v.warning = "Sorry, no words fit here (mistake below?)"
			return
		}
	}
	v.hintDeduction(wordSuggest)
}

func (v *viewPlay) hintWord() {
	s := v.entries[v.onStep]
	if len(strings.Trim(s, " -_.")) == 0 {
		delete(v.okWords, v.onStep)
		v.hintWordTemplate()
	} else if ch := s[v.onChar : v.onChar+1]; ch == "." || ch == "_" || ch == "-" {
		letters := make([]string, 0, 26)
		lMap := make(map[string]bool)
		for _, solution := range v.puzzle.Solutions {
			l := solution.Ladder()[v.onStep+1].String()[v.onChar : v.onChar+1]
			if !lMap[l] {
				lMap[l] = true
				letters = append(letters, l)
			}
		}
		slices.Sort(letters)
		v.hint = "Try " + strings.Join(letters, ",")
		v.hintDeduction(letter)
	}
}

func (v *viewPlay) hintWordTemplate() {
	prevWord := v.previousWord(true)
	nextWord := v.nextWord(true)
	switch {
	case prevWord == nil && nextWord == nil:
		template := []rune(v.puzzle.Solutions[0].Ladder()[v.onStep+1].String())
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			templateDifferences(template, ladder[v.onStep+1], ladder[v.onStep], ladder[v.onStep+2])
		}
		v.entries[v.onStep] = string(template)
	case prevWord != nil && nextWord != nil:
		template := []rune(prevWord.String())
		templateDifferences(template, nextWord)
		v.entries[v.onStep] = string(template)
	case v.onStep == 0:
		template := []rune(v.puzzle.StartWord.String())
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			templateDifferences(template, ladder[1])
		}
		v.entries[v.onStep] = string(template)
	case v.onStep == len(v.entries)-1:
		template := []rune(v.puzzle.EndWord.String())
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			templateDifferences(template, ladder[v.onStep+1])
		}
		v.entries[v.onStep] = string(template)
	case prevWord != nil:
		template := []rune(prevWord.String())
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			templateDifferences(template, ladder[v.onStep+1])
		}
		v.entries[v.onStep] = string(template)
	default:
		// only next word is known
		template := []rune(nextWord.String())
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			templateDifferences(template, ladder[v.onStep+1])
		}
		v.entries[v.onStep] = string(template)
	}
	if v.entries[v.onStep] == strings.Repeat("_", v.puzzle.WordLength) {
		v.warning = "Sorry, couldn't give a useful hint"
		v.onChar = 0
		return
	}
	if uAt := strings.IndexRune(v.entries[v.onStep], '_'); uAt != -1 {
		v.onChar = uAt
	} else {
		v.onChar = 0
	}
	v.hintDeduction(wordPattern)
}

func templateDifferences(template []rune, otherWords ...*words.Word) {
	for _, w := range otherWords {
		for i, r := range []rune(w.String()) {
			if r != template[i] {
				template[i] = '_'
			}
		}
	}
}

type hint int

const (
	solutionsPeek hint = iota
	wordSuggest
	wordPattern
	letter
)

func (v *viewPlay) hintDeduction(h hint) {
	switch h {
	case solutionsPeek:
		v.currentScore = 0
	case wordSuggest:
		if !v.hintsGiven[strconv.Itoa(v.onStep)] {
			v.hintsGiven[strconv.Itoa(v.onStep)] = true
			v.currentScore = math.Floor(v.currentScore - v.puzzle.RungScore)
		}
	case wordPattern:
		if !v.hintsGiven[strconv.Itoa(v.onStep)] {
			v.hintsGiven[strconv.Itoa(v.onStep)] = true
			v.currentScore = math.Floor(v.currentScore - v.puzzle.DeductionPatternHint)
		}
	case letter:
		if !v.hintsGiven[strconv.Itoa(v.onStep)+":"+strconv.Itoa(v.onChar)] {
			v.hintsGiven[strconv.Itoa(v.onStep)+":"+strconv.Itoa(v.onChar)] = true
			v.currentScore = math.Floor(v.currentScore - v.puzzle.DeductionPositionHint)
		}
	}
	if v.currentScore < 0 {
		v.currentScore = 0
	}
}

func (v *viewPlay) generateNew() tea.Cmd {
	return func() tea.Msg {
		puzzle, err := generator.GeneratePuzzle(v.puzzle.WordLength, v.puzzle.LadderLength, nil, nil)
		return newResult{puzzle: puzzle, err: err}
	}
}

func (v *viewPlay) start(puzzle generator.Puzzle, dict *words.Dictionary) {
	v.puzzle = puzzle
	if dict != nil {
		v.dictionary = dict
	}
	v.offsetY, v.onStep, v.onChar = 0, 0, 0
	v.currentScore = puzzle.MaxScore
	v.hint = ""
	v.warning = ""
	v.wrong = ""
	v.hintsGiven = make(map[string]bool)
	v.solved = false
	v.okWords = make(map[int]bool)
	v.entries = make([]string, puzzle.LadderLength-2)
	for i := range puzzle.LadderLength - 2 {
		v.entries[i] = strings.Repeat("_", puzzle.WordLength)
	}
}
