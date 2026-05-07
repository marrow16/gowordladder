package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
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
	help
)

func (m mode) String() string {
	switch m {
	case solve:
		return "Solver"
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
	case help:
		return "Help"
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
type viewShow interface {
	show(backMode mode, backView view)
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
	viewHelp      helpView

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
	initialMode := play
	gv := &viewGenerate{}
	if puzzle, err := generator.GeneratePuzzle(p.WordLength, p.LadderLength, nil, nil); err == nil {
		pv.start(*puzzle, words.NewDictionary(puzzle.WordLength))
	} else {
		initialView = gv
		initialMode = generate
	}
	return &model{
		logger:              l,
		prefs:               p,
		mode:                initialMode,
		currentView:         initialView,
		viewSolve:           &viewSolve{},
		viewGenerate:        gv,
		viewPlay:            pv,
		viewSolutions:       &viewSolutions{},
		viewLookup:          &viewLookup{},
		viewScores:          &viewScores{},
		viewHelp:            &viewHelp{},
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
		case "f1":
			if m.mode != help {
				m.viewHelp.show(m.mode, m.currentView)
				m.mode = help
				m.currentView = m.viewHelp
			}
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
	hdr := "Go Word Ladder - " + m.mode.String()
	if m.mode != help {
		const helpHdr = " F1:help"
		return headerStyle.Width(m.width).Render(center3(m.width, helpHdr, hdr, ""))
	} else {
		return headerStyle.Width(m.width).Render(hdr)
	}
}

func (m *model) footerView() string {
	if h := m.currentView.help(); h != "" {
		return helpStyle.Width(m.width).Render(h + "  •  esc: Exit")
	}
	return helpStyle.Width(m.width).Render("esc: Exit")
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

func center3(wd int, left, mid, right string) string {
	ll, lm, lr := len(left), len(mid), len(right)
	lmw := lm / 2
	rmw := lm - lmw
	lw := wd / 2
	rw := wd - lw
	lpad, rpad := "", ""
	if w := lw - lmw - ll; w > 0 {
		lpad = strings.Repeat(" ", w)
	}
	if w := rw - rmw - lr; w > 0 {
		rpad = strings.Repeat(" ", w)
	}
	return left + lpad + mid + rpad + right
}

func padLines(lines int) string {
	if lines > 0 {
		return strings.Repeat("\n", lines)
	}
	return ""
}

func wrap(text string, maxWidth int) []string {
	if len(text) <= maxWidth {
		return []string{text}
	}
	wds := strings.Fields(text)
	lines := make([]string, 0, len(wds)/2)
	var current string
	for _, w := range wds {
		if current == "" {
			current = w
			continue
		}
		if len(current)+1+len(w) <= maxWidth {
			current += " " + w
		} else {
			lines = append(lines, current)
			current = w
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
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
