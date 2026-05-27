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
	error          string
	wordsDisplayed wordPoints
	scrollbar      layout.Scrollbar
	maxItems       int
}

func (v *viewScores) show(backMode mode, backView view) {
	v.offsetY = 0
	v.backMode = backMode
	v.backView = backView
}

func (v *viewScores) render(sf layout.Surface, m *model) *tea.Cursor {
	v.wordsDisplayed = make(wordPoints)
	v.maxItems = len(m.prefs.HighScores)
	for i, s := range m.prefs.HighScores {
		row := (i - v.offsetY) * 2
		sf.TextRight(row, 1, 3, strconv.Itoa(i+1)+".", boldStyle)
		sf.Text(row, 5, strconv.FormatFloat(s.Score, 'f', 0, 64)+" ("+strconv.FormatFloat((s.Score/s.MaxScore)*100, 'f', 0, 64)+"%)", boldStyle)
		sf.TextRun(row+1, 5, layout.NewRuns(s.Date+"  ", scoreDetailStyle).
			Add(s.StartWord, highlightStyle).Add(" to ").Add(s.EndWord, highlightStyle).
			Add(" ("+strconv.Itoa(s.LadderLength)+" rungs)", scoreDetailStyle))
		dtWidth := len(s.Date)
		v.wordsDisplayed.addWord(s.StartWord, sf.AbsoluteTop()+row+1, sf.AbsoluteLeft()+dtWidth+7)
		v.wordsDisplayed.addWord(s.EndWord, sf.AbsoluteTop()+row+1, sf.AbsoluteLeft()+dtWidth+11+len(s.StartWord))
	}
	if v.maxItems > 1 {
		v.scrollbar = layout.NewVerticalScrollbar(v.scrollHandler).ItemSize(2)
		v.scrollbar.Draw(sf, v.maxItems+(sf.Height()/2)-2, v.offsetY)
	} else {
		v.scrollbar = nil
	}
	return nil
}

func (v *viewScores) scroll(msg tea.Msg) (handled bool) {
	if v.scrollbar != nil {
		return v.scrollbar.Update(msg)
	}
	return false
}

func (v *viewScores) scrollHandler(evt layout.ScrollEvent) {
	v.error = ""
	switch evt {
	case layout.ScrollUp, layout.ScrollPageUp:
		v.offsetY--
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case layout.ScrollDown, layout.ScrollPageDown:
		v.offsetY++
		if v.offsetY > v.maxItems-1 {
			v.offsetY = v.maxItems - 1
		}
	case layout.ScrollHome:
		v.offsetY = 0
	case layout.ScrollEnd:
		v.offsetY = v.maxItems - 1
	}
}

var (
	scoreDetailStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))
)

func (v *viewScores) helpLines() ([]string, *lipgloss.Style) {
	if v.error != "" {
		return []string{
			v.error,
			ctrlAgain + ": Play again  •  " + ctrlNew + ": Clear  •  " + back + ": Back",
		}, &errorStyle
	}
	return []string{ctrlAgain + ": Play again  •  " + ctrlNew + ": Clear  •  " + back + ": Back"}, nil
}

func (v *viewScores) menu() []menuItem {
	return []menuItem{
		{text: "Play again", key: ctrlAgain},
		{text: "Clear", key: ctrlNew},
	}
}

func (v *viewScores) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.error = ""
	switch msg.String() {
	case back, backspace:
		m.restoreView(highs, v.backMode, v.backView)
	case ctrlAgain:
		if v.offsetY >= 0 && v.offsetY < len(m.prefs.HighScores) {
			hs := m.prefs.HighScores[v.offsetY]
			if puzzle, err := generator.GeneratePuzzle(hs.WordLength, hs.LadderLength, &hs.StartWord, &hs.EndWord); err == nil {
				m.play(*puzzle)
			} else {
				v.error = err.Error()
			}
		}
	case ctrlNew:
		v.offsetY = 0
		m.clearScores()
	}
	return nil
}

func (v *viewScores) click(m *model, msg tea.Mouse) tea.Cmd {
	v.error = ""
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
