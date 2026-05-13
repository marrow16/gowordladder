package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marrow16/gowordladder/words"
	"strings"
)

type helpView interface {
	view
	viewShow
}

type viewHelp struct {
	backMode mode
	backView view
	offsetY  int

	cachedLines []string
	cachedWidth int
}

func (v *viewHelp) content(m *model) (string, *tea.Cursor) {
	const (
		footerLines = 2
	)
	var sb strings.Builder
	sb.Grow(m.height * m.width)
	lines := 1
	maxLines := m.height - lines - footerLines
	if m.width != v.cachedWidth {
		v.cachedLines = helpText.render(m.width)
		v.cachedWidth = m.width
	}
	useLines := append(v.cachedLines,
		helpStyle.Render(" "+strings.Repeat(horizontal, m.width-2)+" "),
		helpHeaderStyle.Width(m.width).Render("Current Dictionary"),
		highlightStyle.Width(m.width).AlignHorizontal(lipgloss.Center).Render(" "+words.CurrentDictionary()),
		helpStyle.Width(m.width).Render(" "+fSwitch+": Switch"),
	)
	for l := 0; l < maxLines && (l+v.offsetY) < len(useLines); l++ {
		sb.WriteString("\n")
		sb.WriteString(useLines[l+v.offsetY])
		lines++
	}
	sb.WriteString(padLines(m.height - lines - footerLines))
	return sb.String(), nil
}

func (v *viewHelp) help() string {
	return back + "/" + backspace + ": Back"
}

func (v *viewHelp) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back, backspace:
		m.restoreView(help, v.backMode, v.backView)
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case down:
		if v.offsetY < len(v.cachedLines)-m.height+10 {
			v.offsetY++
		}
	case home:
		v.offsetY = 0
	case pageUp:
		v.offsetY -= m.height - 4
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case pageDown:
		v.offsetY += m.height - 4
	}
	return nil
}

func (v *viewHelp) update(m *model, msg tea.Msg) tea.Cmd {
	return nil
}

func (v *viewHelp) wordLength() int {
	return 0
}

func (v *viewHelp) currentWord() string {
	return ""
}

func (v *viewHelp) show(backMode mode, backView view) {
	v.backMode = backMode
	v.backView = backView
}
