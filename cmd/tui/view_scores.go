package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"gowordladder/generator"
	"strings"
)

type viewScores struct {
	backMode mode
	backView view
	offsetY  int
}

type scoresView interface {
	view
	show(backMode mode, backView view)
}

func (v *viewScores) show(backMode mode, backView view) {
	v.offsetY = 0
	v.backMode = backMode
	v.backView = backView
}

func (v *viewScores) content(m *model) (string, *tea.Cursor) {
	const (
		footerLines = 2
	)
	var sb strings.Builder
	sb.WriteString("\n")
	lines := 2
	if len(m.prefs.HighScores) == 0 {
		sb.WriteString(errorStyle.Render(" No high scores to show"))
	} else {
		maxLines := m.height - lines - footerLines
		showLines := make([]string, 0, len(m.prefs.HighScores)*2)
		for i, s := range m.prefs.HighScores {
			showLines = append(showLines,
				boldScoreStyle.Render(fmt.Sprintf(" %2d. %.0f (%.0f%%)   ", i+1, s.Score, (s.Score/s.MaxScore)*100)),
				scoreDetailStyle.Render("     "+s.Date+"  ")+
					highlightStyle.Render(s.StartWord)+scoreDetailStyle.Render(" to ")+highlightStyle.Render(s.EndWord)+
					scoreDetailStyle.Render(fmt.Sprintf(" (%d rungs)", s.LadderLength)),
			)
		}
		for l := 0; l < maxLines && (l+v.offsetY) < len(showLines); l++ {
			sb.WriteString(showLines[l+v.offsetY])
			sb.WriteString("\n")
			lines++
		}
	}
	sb.WriteString(padLines(m.height - lines - footerLines))
	return sb.String(), nil
}

var (
	boldScoreStyle   = lipgloss.NewStyle().Bold(true)
	scoreDetailStyle = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))
)

func (v *viewScores) help() string {
	return "ctrl+n: Clear  •  ctrl+p: Play again  •  ctrl+b: Back"
}

func (v *viewScores) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+b", "backspace":
		m.restoreView(v.backMode, v.backView)
	case "ctrl+n":
		v.offsetY = 0
		m.clearScores()
	case "ctrl+p":
		if h := v.offsetY / 2; h < len(m.prefs.HighScores) {
			hs := m.prefs.HighScores[h]
			if puzzle, err := generator.GeneratePuzzle(hs.WordLength, hs.LadderLength, &hs.StartWord, &hs.EndWord); err == nil {
				m.play(*puzzle)
			}
		}
	case "up":
		if v.offsetY > 0 {
			v.offsetY -= 2
		}
	case "down":
		v.offsetY += 2
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
