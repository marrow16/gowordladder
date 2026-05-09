package main

import (
	tea "charm.land/bubbletea/v2"
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
}

func (v *viewHelp) content(m *model) (string, *tea.Cursor) {
	const (
		footerLines = 2
	)
	var sb strings.Builder
	lines := 1
	maxLines := m.height - lines - footerLines
	showLines := helpText.render(m.width)
	for l := 0; l < maxLines && (l+v.offsetY) < len(showLines); l++ {
		sb.WriteString("\n")
		sb.WriteString(showLines[l+v.offsetY])
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
		m.restoreView(v.backMode, v.backView)
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case down:
		v.offsetY++
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
