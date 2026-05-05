package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"gowordladder/generator"
	"gowordladder/solving"
	"gowordladder/words"
	"log/slog"
	"math/rand"
	"os"
	"strings"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type mode int

const (
	solve mode = iota
	generate
	solutions
	play
	lookup
	highs
)

func (m mode) String() string {
	switch m {
	case solve:
		return "Solve"
	case generate:
		return "Generate"
	case solutions:
		return "Solutions"
	case play:
		return "Play"
	case lookup:
		return "Lookup Word"
	case highs:
		return "High Scores"
	}
	return ""
}

type view interface {
	content(m *model) (string, *tea.Cursor)
	help() string
	key(m *model, msg tea.KeyPressMsg) tea.Cmd
	update(m *model, msg tea.Msg) tea.Cmd
	wordLength() int
	currentWord() string
}

type model struct {
	logger      *slog.Logger
	prefs       *prefs
	mode        mode
	width       int
	height      int
	currentView view
	// mode views...
	viewSolve     view
	viewGenerate  view
	viewPlay      playView
	viewSolutions solutionsView
	viewLookup    lookupView
	viewScores    scoresView

	dictionary          *words.Dictionary
	dictionaryLoadTimes map[int]time.Duration
}

func newModel(withLogging bool) *model {
	var l *slog.Logger
	if withLogging {
		if lf, err := os.OpenFile("tui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
			l = slog.New(slog.NewJSONHandler(lf, &slog.HandlerOptions{Level: slog.LevelInfo}))
		}
	}
	p := newPrefs()
	pv := &viewPlay{}
	var initialView view = pv
	gv := &viewGenerate{}
	if puzzle, err := generator.GeneratePuzzle(p.WordLength, p.LadderLength, nil, nil); err == nil {
		pv.start(*puzzle, words.NewDictionary(puzzle.WordLength))
	} else {
		initialView = gv
	}
	return &model{
		logger:              l,
		prefs:               p,
		mode:                solve,
		currentView:         initialView,
		viewSolve:           &viewSolve{},
		viewGenerate:        gv,
		viewPlay:            pv,
		viewSolutions:       &viewSolutions{},
		viewLookup:          &viewLookup{},
		viewScores:          &viewScores{},
		dictionaryLoadTimes: map[int]time.Duration{},
	}
}

func (m *model) log(msg string, args ...any) {
	if m.logger != nil {
		m.logger.Info(msg, args...)
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch mt := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = mt.Width
		m.height = mt.Height
	//	case tea.ResumeMsg:
	//		m.suspending = false
	//		return m, nil
	case tea.KeyPressMsg:
		switch mt.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+w":
			if m.mode != lookup {
				cw := m.currentView.currentWord()
				cmd := m.viewLookup.lookupWord(cw, m.mode, m.currentView)
				m.mode = lookup
				m.currentView = m.viewLookup
				return m, cmd
			}
		case "ctrl+t":
			if m.mode != highs {
				m.viewScores.show(m.mode, m.currentView)
				m.mode = highs
				m.currentView = m.viewScores
			}
		default:
			if !m.viewSwitch(mt.String()) {
				return m, m.currentView.key(m, mt)
			}
		}
	default:
		return m, m.currentView.update(m, msg)
	}
	return m, nil
}

func (m *model) viewSwitch(key string) bool {
	switch {
	case key == "ctrl+s" && m.mode != solve:
		m.mode = solve
		m.currentView = m.viewSolve
		if wl := m.currentView.wordLength(); wl > 0 {
			m.loadDictionary(wl)
		}
		return true
	case key == "ctrl+g" && m.mode != generate:
		m.mode = generate
		m.currentView = m.viewGenerate
		if wl := m.currentView.wordLength(); wl > 0 {
			m.loadDictionary(wl)
		}
		return true
	}
	return false
}

func (m *model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.View{AltScreen: true}
	}
	vc, csr := m.currentView.content(m)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.headerView(),
		vc,
		m.footerView(),
	)
	v := tea.NewView(content)
	v.Cursor = csr
	v.AltScreen = true
	return v
}

func (m *model) headerView() string {
	return headerStyle.Width(m.width).Render(fmt.Sprintf("Go Word Ladder - %s", m.mode))
}

func (m *model) footerView() string {
	return helpStyle.Width(m.width).Render(m.currentView.help() + "  •  esc: Exit")
}

func (m *model) loadDictionary(wordLength int) *words.Dictionary {
	if m.dictionary == nil || m.dictionary.WordLength() != wordLength {
		start := time.Now()
		m.dictionary = words.NewDictionary(wordLength)
		dur := time.Since(start)
		if _, ok := m.dictionaryLoadTimes[wordLength]; !ok {
			m.dictionaryLoadTimes[wordLength] = dur
		}
	}
	return m.dictionary
}

func (m *model) play(puzzle generator.Puzzle) {
	m.prefs.WordLength = puzzle.WordLength
	m.prefs.LadderLength = puzzle.LadderLength
	m.prefs.save()
	m.viewPlay.start(puzzle, m.loadDictionary(puzzle.WordLength))
	m.mode = play
	m.currentView = m.viewPlay
}

func (m *model) playScore(score float64, p generator.Puzzle) {
	if score > 0 {
		m.prefs.addScore(score, p.MaxScore, p.WordLength, p.LadderLength, p.StartWord.String(), p.EndWord.String())
	}
}

func (m *model) clearScores() {
	m.prefs.clearScores()
}

func (m *model) showSolutions(s []*solving.Solution) {
	m.viewSolutions.setSolutions(s, m.mode, m.currentView)
	m.mode = solutions
	m.currentView = m.viewSolutions
}

func (m *model) restoreView(rm mode, rv view) {
	m.currentView = rv
	m.mode = rm
}

func padLines(lines int) string {
	if lines > 0 {
		return strings.Repeat("\n", lines)
	}
	return ""
}

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#0000ff")).
			AlignHorizontal(lipgloss.Center)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).
			AlignHorizontal(lipgloss.Center)
	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008000"))
	errorStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#ff0000"))
)
