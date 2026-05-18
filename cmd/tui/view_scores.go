package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/generator"
	"strconv"
)

type scoresView interface {
	view
	viewShowable
}

type viewScores struct {
	backMode       mode
	backView       view
	offsetY        int
	wordsDisplayed wordPoints
}

func (v *viewScores) show(backMode mode, backView view) {
	v.offsetY = 0
	v.backMode = backMode
	v.backView = backView
}

func (v *viewScores) render(sf layout.Surface, m *model) *tea.Cursor {
	v.wordsDisplayed = make(wordPoints)
	rgn := sf.Region(1, 1, sf.Height(), sf.Width()-2)
	for i, s := range m.prefs.HighScores {
		row := (i - v.offsetY) * 2
		rgn.TextRight(row, 0, 3, strconv.Itoa(i+1)+".", boldStyle)
		rgn.Text(row, 4, strconv.FormatFloat(s.Score, 'f', 0, 64)+" ("+strconv.FormatFloat((s.Score/s.MaxScore)*100, 'f', 0, 64)+"%)", boldStyle)
		rgn.TextRun(row+1, 4, layout.NewRuns(s.Date+"  ", scoreDetailStyle).
			Add(s.StartWord, highlightStyle).Add(" to ").Add(s.EndWord, highlightStyle).
			Add("("+strconv.Itoa(s.LadderLength)+" rungs)", scoreDetailStyle))
		dtWidth := len(s.Date)
		v.wordsDisplayed.addWord(s.StartWord, row+3, dtWidth+7)
		v.wordsDisplayed.addWord(s.EndWord, row+3, dtWidth+11+len(s.StartWord))
	}
	return nil
}

var (
	scoreDetailStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))
)

func (v *viewScores) helpLines() ([]string, *lipgloss.Style) {
	return []string{ctrlNew + ": Clear  •  " + ctrlPlay + ": Play again  •  " + back + ": Back"}, nil
}

func (v *viewScores) menu() []menuItem {
	return nil
}

func (v *viewScores) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back, backspace:
		m.restoreView(highs, v.backMode, v.backView)
	case ctrlNew:
		v.offsetY = 0
		m.clearScores()
	case ctrlPlay:
		if h := v.offsetY / 2; h < len(m.prefs.HighScores) {
			hs := m.prefs.HighScores[h]
			if puzzle, err := generator.GeneratePuzzle(hs.WordLength, hs.LadderLength, &hs.StartWord, &hs.EndWord); err == nil {
				m.play(*puzzle)
			}
		}
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case down:
		if v.offsetY+1 < len(m.prefs.HighScores) {
			v.offsetY++
		}
	}
	return nil
}

func (v *viewScores) click(m *model, msg tea.Mouse) tea.Cmd {
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

func (v *viewScores) update(m *model, msg tea.Msg) tea.Cmd {
	return nil
}

func (v *viewScores) wordLength() int {
	return 0
}

func (v *viewScores) currentWord() string {
	return ""
}
