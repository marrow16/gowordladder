package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"strconv"
	"strings"
	"time"
)

type lingoStep int

const (
	lingoAdjustWordLength lingoStep = iota
	lingoAdjustMaxGuesses
	lingoAdjustTimer
	lingoAdjustDone
	lingoPlay
	lingoDone
)

type viewLingo struct {
	loaded         bool
	step           lingoStep
	onGuess        int
	displayWords   [][]layout.RunItem
	words          []string
	wordsDisplayed wordPoints
	currentInput   input
	currentError   string
	word           string
	bank           int
	wordLen        int
	maxGuesses     int
	timer          int
	timeRemaining  int
	success        bool
	failed         string
	value          int
	remaining      int
	deduct         int
}

const (
	lingoPrefBank       = "LingoBank"
	lingoPrefWordLength = "LingoWordLength"
	lingoPrefMaxGuesses = "LingoMaxGuesses"
	lingoPrefTimer      = "LingoTimer"
)

func (v *viewLingo) load(m *model) {
	if !v.loaded {
		v.loaded = true
		v.bank = m.prefs.getInt(lingoPrefBank)
		v.wordLen = m.prefs.getInt(lingoPrefWordLength)
		if v.wordLen < 2 || v.wordLen > 15 {
			v.wordLen = 4
		}
		v.maxGuesses = m.prefs.getInt(lingoPrefMaxGuesses)
		if v.maxGuesses < 4 || v.maxGuesses > 15 {
			v.maxGuesses = 5
		}
		v.timer = m.prefs.getInt(lingoPrefTimer)
		if v.timer < 0 || v.timer > 60 {
			v.timer = 30
		}
	}
}

func (v *viewLingo) render(sf layout.Surface, m *model) *tea.Cursor {
	const (
		promptWordLength = "Word Length: "
		promptMaxGuesses = "Max Guesses: "
		promptTimer      = "Timer: "
		promptLen        = len(promptWordLength) + 1
	)
	v.load(m)
	v.wordsDisplayed = wordPoints{}
	hdr := layout.Runs{
		{Text: " "},
		{Text: " "},
		{Text: "Bank: £" + commas(v.bank) + " "},
	}
	var csr *tea.Cursor
	switch v.step {
	case lingoAdjustWordLength:
		hdr[1] = layout.RunItem{Text: " Adjust Game "}
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2, current: strconv.Itoa(v.wordLen)}
		}
		sf.TextRight(2, 0, promptLen, promptWordLength)
		sf.TextFixed(2, promptLen, 2, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(sf.AbsoluteLeft()+promptLen+v.currentInput.cursorPos(), sf.AbsoluteTop()+2)
		if v.currentError != "" {
			sf.Text(2, promptLen+4, v.currentError, errorStyle)
		} else {
			sf.Text(2, promptLen+4, "(2-15)", helpStyle)
		}
	case lingoAdjustMaxGuesses:
		hdr[1] = layout.RunItem{Text: " Adjust Game "}
		sf.TextRight(2, 0, promptLen, promptWordLength)
		sf.TextFixed(2, promptLen, 2, strconv.Itoa(v.wordLen), inputStyle)
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2, current: strconv.Itoa(v.maxGuesses)}
		}
		sf.TextRight(3, 0, promptLen, promptMaxGuesses)
		sf.TextFixed(3, promptLen, 2, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(sf.AbsoluteLeft()+promptLen+v.currentInput.cursorPos(), sf.AbsoluteTop()+3)
		if v.currentError != "" {
			sf.Text(3, promptLen+4, v.currentError, errorStyle)
		} else {
			sf.Text(3, promptLen+4, "(4-15)", helpStyle)
		}
	case lingoAdjustTimer:
		hdr[1] = layout.RunItem{Text: " Adjust Game "}
		sf.TextRight(2, 0, promptLen, promptWordLength)
		sf.TextFixed(2, promptLen, 2, strconv.Itoa(v.wordLen), inputStyle)
		sf.TextRight(3, 0, promptLen, promptMaxGuesses)
		sf.TextFixed(3, promptLen, 2, strconv.Itoa(v.maxGuesses), inputStyle)
		if v.currentInput == nil {
			v.currentInput = &numberInput{maxLength: 2, current: strconv.Itoa(v.timer)}
		}
		sf.TextRight(4, 0, promptLen, promptTimer)
		sf.TextFixed(4, promptLen, 2, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(sf.AbsoluteLeft()+promptLen+v.currentInput.cursorPos(), sf.AbsoluteTop()+4)
		if v.currentError != "" {
			sf.Text(4, promptLen+4, v.currentError, errorStyle)
		} else {
			sf.Text(4, promptLen+4, "(0-60 seconds - 0 for no timer)", helpStyle)
		}
	case lingoAdjustDone:
		hdr[1] = layout.RunItem{Text: " Adjust Game "}
		sf.TextRight(2, 0, promptLen, promptWordLength)
		sf.TextFixed(2, promptLen, 2, strconv.Itoa(v.wordLen), inputStyle)
		sf.TextRight(3, 0, promptLen, promptMaxGuesses)
		sf.TextFixed(3, promptLen, 2, strconv.Itoa(v.maxGuesses), inputStyle)
		sf.TextRight(4, 0, promptLen, promptTimer)
		sf.TextFixed(4, promptLen, 2, strconv.Itoa(v.timer), inputStyle)
	case lingoPlay:
		hdr[0] = layout.RunItem{Text: " Max: £" + strconv.Itoa(v.value) + " "}
		hdr[1] = layout.RunItem{Text: " Remaining: £" + strconv.Itoa(v.remaining) + " "}
		leftPos := (m.width - (v.wordLen + 1)) / 2
		sf.BoxRounded(2, leftPos, v.maxGuesses+3, v.wordLen+2, helpStyle)
		sf.Text(3, leftPos+1, v.word[:1])
		sf.Text(3, leftPos+2, strings.Repeat("_", v.wordLen-1), helpStyle)
		for g := 0; g < v.onGuess; g++ {
			sf.TextRun(4+g, leftPos+1, v.displayWords[g])
		}
		if v.currentInput == nil {
			v.currentInput = &wordInput{maxLength: v.wordLen, current: v.word[:1]}
		}
		sf.TextFixed(4+v.onGuess, leftPos+1, v.wordLen, v.currentInput.value(), inputStyle)
		csr = tea.NewCursor(sf.AbsoluteLeft()+leftPos+1+v.currentInput.cursorPos(), sf.AbsoluteTop()+4+v.onGuess)
		if v.timer != 0 {
			leftPos = (m.width - v.timer) / 2
			sf.FillWith(5+v.maxGuesses, leftPos, 1, v.timer, '\u00A0', lingoCorrectStyle)
			sf.FillWith(5+v.maxGuesses, leftPos+v.timeRemaining, 1, v.timer-v.timeRemaining, '\u00A0', lingoErrorStyle)
		}
	case lingoDone:
		hdr[0] = layout.RunItem{Text: " Max: £" + strconv.Itoa(v.value) + " "}
		leftPos := (m.width - (v.wordLen + 1)) / 2
		sf.BoxRounded(2, leftPos, v.maxGuesses+3, v.wordLen+2, helpStyle)
		for g := 0; g < v.onGuess; g++ {
			sf.TextRun(4+g, leftPos+1, v.displayWords[g])
			if v.words[g] != "" {
				v.wordsDisplayed.addWord(v.words[g], sf.AbsoluteTop()+4+g, sf.AbsoluteLeft()+leftPos+1)
			}
		}
		v.wordsDisplayed.addWord(v.word, sf.AbsoluteTop()+3, sf.AbsoluteLeft()+leftPos+1)
		if v.success {
			hdr[1] = layout.RunItem{Text: " Banked: £" + strconv.Itoa(v.remaining) + " "}
			sf.Text(3, leftPos+1, v.word, lingoCorrectStyle)
			sf.TextCenter(5+v.maxGuesses, 1, sf.Width()-2, "Completed! Well Done!", highlightStyle)
		} else {
			sf.Text(3, leftPos+1, v.word)
			sf.TextCenter(5+v.maxGuesses, 1, sf.Width()-2, "Sorry, "+v.failed, errorStyle)
		}
	}
	sf.LineColumns(0, 0, m.width, hdr, headerStyle)
	return csr
}

func (v *viewLingo) helpLines() ([]string, *lipgloss.Style) {
	switch v.step {
	case lingoAdjustWordLength, lingoAdjustMaxGuesses, lingoAdjustTimer:
		return []string{ctrlAgain + ": Adjust"}, nil
	case lingoAdjustDone:
		return []string{enter + ": Play  •  " + ctrlAgain + ": Adjust"}, nil
	default:
		return []string{ctrlNew + ": New  •  " + ctrlAgain + ": Adjust"}, nil
	}
}

func (v *viewLingo) menu() []menuItem {
	switch v.step {
	case lingoAdjustWordLength, lingoAdjustMaxGuesses, lingoAdjustTimer:
		return []menuItem{
			{text: "Adjust", key: ctrlAgain},
		}
	case lingoAdjustDone:
		return []menuItem{
			{text: "Play", key: enter},
			{text: "Adjust", key: ctrlAgain},
		}
	default:
		return []menuItem{
			{text: "New", key: ctrlNew},
			{text: "Adjust", key: ctrlAgain},
		}
	}
}

func (v *viewLingo) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.currentError = ""
	switch msg.String() {
	case ctrlAgain:
		v.step = lingoAdjustWordLength
		v.currentInput = nil
	case ctrlNew:
		return v.enterStart(m)
	case "?":
		if v.step == lingoPlay {
			return v.doHintWord(m)
		}
	case enter:
		switch v.step {
		case lingoAdjustWordLength:
			if n, err := strconv.Atoi(v.currentInput.value()); err == nil && n >= 2 && n <= 15 {
				v.wordLen = n
				v.step = lingoAdjustMaxGuesses
				v.currentInput = nil
				m.prefs.setInt(lingoPrefWordLength, v.wordLen)
			} else {
				v.currentError = "Please enter a number (2-15)"
			}
		case lingoAdjustMaxGuesses:
			if n, err := strconv.Atoi(v.currentInput.value()); err == nil && n >= 4 && n <= 15 {
				v.maxGuesses = n
				v.step = lingoAdjustTimer
				v.currentInput = nil
				m.prefs.setInt(lingoPrefMaxGuesses, v.maxGuesses)
			} else {
				v.currentError = "Please enter a number (4-15)"
			}
		case lingoAdjustTimer:
			if n, err := strconv.Atoi(v.currentInput.value()); err == nil && n >= 0 && n <= 60 {
				v.timer = n
				v.step = lingoAdjustDone
				v.currentInput = nil
				m.prefs.setInt(lingoPrefTimer, v.timer)
			} else {
				v.currentError = "Please enter a number (0-60)"
			}
		case lingoAdjustDone:
			return v.enterStart(m)
		}
	}
	if v.currentInput != nil {
		v.currentInput.key(msg)
		if v.step == lingoPlay && len(v.currentInput.value()) == v.wordLen {
			return v.doGuess(m, v.currentInput.value())
		}
	}
	return nil
}

func (v *viewLingo) click(m *model, msg tea.Mouse) tea.Cmd {
	if v.step == lingoDone {
		if wd, ok := v.wordsDisplayed[pt{msg.Y, msg.X}]; ok {
			switch msg.Button {
			case tea.MouseLeft:
				return m.lookupWord(wd)
			case tea.MouseRight:
				return m.lookupDistance(wd)
			}
		}
	}
	return nil
}

var (
	lingoCorrectStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#00ff00"))
	lingoErrorStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#ff0000"))
	lingoMisplacedStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ffa500"))
)

func (v *viewLingo) doHintWord(m *model) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		dict := m.loadDictionary(v.wordLen)
		candidates := dict.Words()
		word := candidates[rng.Intn(len(candidates))]
		if word.String() != v.word {
			v.remaining -= v.deduct / 2
			return lingoGuessResult{
				display: wordGuessDisplay(v.word, word.String()),
				word:    word.String(),
			}
		}
		return nil
	}
}

func (v *viewLingo) doGuess(m *model, word string) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		if strings.ToUpper(word) == strings.ToUpper(v.word) {
			return lingoGuessResult{
				correct: true,
				display: []layout.RunItem{{Text: word, Styles: []lipgloss.Style{lingoCorrectStyle}}},
				word:    word,
			}
		}
		dict := m.loadDictionary(v.wordLen)
		if _, ok := dict.Word(word); !ok {
			return lingoGuessResult{
				error:   "invalid word!",
				display: []layout.RunItem{{Text: word, Styles: []lipgloss.Style{lingoErrorStyle}}},
			}
		}
		return lingoGuessResult{
			display: wordGuessDisplay(v.word, word),
			word:    word,
		}
	}
}

func wordGuessDisplay(actualWord, guessedWord string) []layout.RunItem {
	result := make([]layout.RunItem, len(actualWord))
	for i := 0; i < len(actualWord); i++ {
		ach := actualWord[i]
		gch := guessedWord[i]
		if ach == gch {
			result[i] = layout.RunItem{Text: string(ach), Styles: []lipgloss.Style{lingoCorrectStyle}}
		} else if strings.Contains(actualWord, string(gch)) {
			result[i] = layout.RunItem{Text: string(gch), Styles: []lipgloss.Style{lingoMisplacedStyle}}
		} else {
			result[i] = layout.RunItem{Text: string(gch)}
		}
	}
	return result
}

func (v *viewLingo) enterStart(m *model) tea.Cmd {
	v.currentError = ""
	return func() tea.Msg {
		dict := m.loadDictionary(v.wordLen)
		candidates := dict.Words()
		word := candidates[rng.Intn(len(candidates))]
		return lingoStartResult{
			word: word.String(),
		}
	}
}

type lingoGuessResult struct {
	correct bool
	error   string
	display []layout.RunItem
	word    string
}

type lingoStartResult struct {
	word string
}

func (v *viewLingo) update(m *model, msg tea.Msg) tea.Cmd {
	switch mt := msg.(type) {
	case lingoStartResult:
		v.timeRemaining = v.timer
		v.currentInput = nil
		v.word = mt.word
		v.onGuess = 0
		v.displayWords = nil
		v.words = nil
		v.step = lingoPlay
		v.value = v.wordLen * 100
		v.remaining = v.value
		v.deduct = v.value / (v.maxGuesses + 1)
		if v.timer != 0 {
			return tick()
		}
	case lingoGuessResult:
		v.timeRemaining = v.timer
		v.currentInput = nil
		v.onGuess++
		v.displayWords = append(v.displayWords, mt.display)
		v.words = append(v.words, mt.word)
		if mt.correct {
			v.success = true
			v.step = lingoDone
			v.bank += v.remaining
			m.prefs.setInt(lingoPrefBank, v.bank)
		} else if mt.error != "" {
			v.success = false
			v.failed = mt.error
			v.step = lingoDone
		} else if v.onGuess >= v.maxGuesses {
			v.success = false
			v.step = lingoDone
			v.failed = "ran out of guesses!"
		} else {
			v.remaining -= v.deduct
		}
	case lingoTickMsg:
		if v.step == lingoPlay && v.timer != 0 {
			v.timeRemaining--
			if v.timeRemaining < 1 {
				v.success = false
				v.step = lingoDone
				v.failed = "ran out of time!"
			} else {
				return tick()
			}
		}
	}
	return nil
}

func (v *viewLingo) wordLength() int {
	return v.wordLen
}

func (v *viewLingo) currentWord() string {
	if v.step == lingoDone {
		return v.word
	}
	return ""
}

type lingoTickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return lingoTickMsg(t)
	})
}
