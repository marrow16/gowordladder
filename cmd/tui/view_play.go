package main

import (
	tea "charm.land/bubbletea/v2"
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
	wordsDisplayed wordPoints
	mousePositions map[int][2]int
}

const (
	playHeaderLines = 2
	playFooterLines = 3
)

func (v *viewPlay) content(m *model) (string, *tea.Cursor) {
	v.wordsDisplayed = make(wordPoints)
	v.mousePositions = make(map[int][2]int)
	var sb strings.Builder
	sb.WriteString(headerStyle.Width(m.width).Render(center3(
		m.width,
		fmt.Sprintf(" Solutions: %d ", len(v.puzzle.Solutions)),
		fmt.Sprintf("Current score: %.0f", v.currentScore),
		fmt.Sprintf(" Max score: %.0f ", v.puzzle.MaxScore))))
	lines := playHeaderLines

	var csr *tea.Cursor
	maxLines := m.height - lines - playFooterLines
	padL := strings.Repeat(" ", ((m.width-v.puzzle.WordLength+2)/2)-3)
	x := len(padL) + 4
	stop := false
	for l := 0; !stop && l < maxLines; l++ {
		sb.WriteString("\n")
		lines++
		rung := l + v.offsetY - 2
		sb.WriteString(padL)
		switch {
		case rung == -2:
			sb.WriteString(helpStyle.Render("   " + topLeft + strings.Repeat(horizontal, v.puzzle.WordLength) + topRight))
		case rung == -1:
			sb.WriteString(helpStyle.Render(" 1 " + vertical))
			if v.solved {
				sb.WriteString(highlightStyle.Render(v.puzzle.StartWord.String()))
			} else {
				sb.WriteString(v.puzzle.StartWord.String())
			}
			sb.WriteString(helpStyle.Render(vertical))
			v.wordsDisplayed.addWord(v.puzzle.StartWord.String(), lines-1, x)
		case rung == v.puzzle.LadderLength-2:
			sb.WriteString(helpStyle.Render(fmt.Sprintf("%2d ", v.puzzle.LadderLength)))
			sb.WriteString(helpStyle.Render(vertical))
			if v.solved {
				sb.WriteString(highlightStyle.Render(v.puzzle.EndWord.String()))
			} else {
				sb.WriteString(v.puzzle.EndWord.String())
			}
			sb.WriteString(helpStyle.Render(vertical))
			v.wordsDisplayed.addWord(v.puzzle.EndWord.String(), lines-1, x)
		case rung == v.puzzle.LadderLength-1:
			sb.WriteString(helpStyle.Render("   " + bottomLeft + strings.Repeat(horizontal, v.puzzle.WordLength) + bottomRight))
		case rung < v.puzzle.LadderLength:
			if rung == v.onStep {
				csr = tea.NewCursor(len(padL)+4+v.onChar, lines-1)
				csr.Color = playCursorColor
			}
			sb.WriteString(helpStyle.Render(fmt.Sprintf("%2d ", rung+2)))
			sb.WriteString(helpStyle.Render(vertical))
			if v.solved {
				sb.WriteString(highlightStyle.Render(v.entries[rung]))
				v.wordsDisplayed.addWord(v.entries[rung], lines-1, x)
			} else if v.okWords[rung] {
				sb.WriteString(highlightStyle.Render(v.entries[rung]))
				v.mousePositions[lines-1] = [2]int{rung, x}
			} else {
				sb.WriteString(v.entries[rung])
				v.mousePositions[lines-1] = [2]int{rung, x}
			}
			sb.WriteString(helpStyle.Render(vertical))
		default:
			stop = true
		}
	}

	sb.WriteString(padLines(m.height - lines - playHeaderLines))
	return sb.String(), csr
}

func (v *viewPlay) help() string {
	const (
		firstHelp   = ctrlHelp + ": Solutions  •  ?: Hint  •  " + ctrlFill + ": Fill  •  space: Clear"
		secondHelp  = "\n" + ctrlNew + ": New  •  " + ctrlGenerate + ": Generate  •  " + ctrlSolver + ": Solver"
		solvedHelp  = "\n" + ctrlNew + ": New  •  " + ctrlHelp + ": Solutions  •  " + ctrlGenerate + ": Generate"
		defaultHelp = firstHelp + secondHelp
	)
	switch {
	case v.solved:
		return hintStyle.Render(fmt.Sprintf("You solved it!  Score: %.0f (%0.f%%)", v.currentScore, (v.currentScore/v.puzzle.MaxScore)*100)) + solvedHelp
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

func (v *viewPlay) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.hint = ""
	v.warning = ""
	v.wrong = ""
	switch msg.String() {
	case ctrlNew:
		return v.generateNew()
	case ctrlHelp:
		if !v.solved {
			v.hintDeduction(solutionsPeek)
		}
		m.showSolutions(v.puzzle.Solutions)
	case up:
		if v.onStep > 0 {
			v.onStep--
			v.ensureCursorVisible(m)
		}
	case down:
		if v.onStep < len(v.entries)-1 {
			v.onStep++
			v.ensureCursorVisible(m)
		}
	case enter, tab:
		if !v.solved {
			v.checkWord(m)
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
	case shiftTab:
		if !v.solved {
			v.checkWord(m)
			if v.wrong == "" && v.warning == "" && v.onStep > 0 {
				v.onStep--
				v.onChar = 0
				v.ensureCursorVisible(m)
			}
		} else if v.onStep < len(v.entries)-1 {
			v.onStep++
			v.onChar = 0
			v.ensureCursorVisible(m)
		}
	case left:
		if v.onChar > 0 {
			v.onChar--
		} else if v.onStep > 0 {
			v.onStep--
			v.onChar = v.puzzle.WordLength - 1
			v.ensureCursorVisible(m)
		}
	case right:
		if v.onChar < v.puzzle.WordLength-1 {
			v.onChar++
		} else if v.onStep < len(v.entries)-1 {
			v.onStep++
			v.onChar = 0
			v.ensureCursorVisible(m)
		}
	case backspace:
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
	case ctrlFill:
		if !v.solved {
			delete(v.okWords, v.onStep)
			v.fillWord(m)
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
				if v.onChar < v.puzzle.WordLength-1 && k != "_" {
					v.onChar++
				}
				v.checkWord(m)
			}
		}
	}
	return nil
}

func (v *viewPlay) click(m *model, msg tea.Mouse) tea.Cmd {
	if wd, ok := v.wordsDisplayed[pt{msg.Y, msg.X}]; ok {
		return m.lookupWord(wd)
	}
	if mp, ok := v.mousePositions[msg.Y]; ok {
		return func() tea.Msg {
			v.onStep = mp[0]
			v.onChar = msg.X - mp[1]
			if v.onChar < 0 {
				v.onChar = 0
			} else if v.onChar >= v.puzzle.WordLength {
				v.onChar = v.puzzle.WordLength - 1
			}
			return nil
		}
	}
	return nil
}

func (v *viewPlay) ensureCursorVisible(m *model) {
	maxLines := m.height - playHeaderLines - playFooterLines
	visibleL := v.onStep - v.offsetY + playHeaderLines
	if visibleL < 0 {
		v.offsetY = v.onStep + playHeaderLines
		// if we're close enough to the top, snap so that first word is visible...
		if v.offsetY < 3 {
			v.offsetY = 0
		}
	} else if visibleL >= maxLines {
		v.offsetY = v.onStep + playHeaderLines - maxLines + 1
		lastRung := v.puzzle.LadderLength - 1
		maxOffsetY := lastRung + playFooterLines - maxLines
		// if we're close enough to the bottom, snap so the final word is visible...
		if maxOffsetY-v.offsetY < playFooterLines {
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

func (v *viewPlay) currentWord() string {
	s := v.entries[v.onStep]
	if len(s) >= 2 && (strings.Count(s, "_") == 1 || isAllAZ(s)) {
		return s
	}
	return ""
}

func (v *viewPlay) checkWord(m *model) {
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
					m.playScore(v.currentScore, v.puzzle)
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

func (v *viewPlay) fillWord(m *model) {
	prevWord := v.previousWord(true)
	nextWord := v.nextWord(true)
	switch {
	case prevWord == nil && nextWord == nil:
		wd := v.puzzle.Solutions[0].Ladder()[v.onStep+1]
		v.entries[v.onStep] = wd.String()
		v.onChar = 0
		v.checkWord(m)
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
			v.checkWord(m)
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
			v.checkWord(m)
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
			v.checkWord(m)
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
			if ladder[v.onStep].String() == prevWord.String() {
				templateDifferences(template, ladder[v.onStep+1])
			}
		}
		v.entries[v.onStep] = string(template)
	default:
		// only next word is known
		template := []rune(nextWord.String())
		for _, solution := range v.puzzle.Solutions {
			ladder := solution.Ladder()
			if ladder[v.onStep+2].String() == nextWord.String() {
				templateDifferences(template, ladder[v.onStep+1])
			}
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
