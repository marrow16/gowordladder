package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/marrow16/gowordladder/words"
	"slices"
	"strings"
)

type switchView interface {
	view
	viewShow
}

type viewSwitch struct {
	backMode     mode
	backView     view
	input        *stringInput
	currentError string
}

func (v *viewSwitch) content(m *model) (string, *tea.Cursor) {
	const (
		prompt      = " Source: "
		promptLen   = len(prompt)
		footerLines = 3
	)
	var sb strings.Builder
	sb.Grow(m.width * m.height)
	sb.WriteString("\n")
	lines := 1
	sb.WriteString(prompt)
	v.input.maxWidth = m.width - promptLen - 1
	s, cpx := v.input.render()
	csr := tea.NewCursor(promptLen+cpx, 2)
	sb.WriteString(s)
	sb.WriteString("\n")
	lines++
	if v.currentError != "" {
		sb.WriteString(errorStyle.Width(m.width - promptLen - 2).Render(strings.Repeat(" ", promptLen) + v.currentError))
		sb.WriteString("\n")
		lines++
	}
	sb.WriteString("\n")
	sb.WriteString(helpStyle.Width(m.width-2).Render(" Press "+enter+" to switch.") + "\n")
	sb.WriteString(helpStyle.Width(m.width-2).Render(" Use internal: \""+strings.Join(internalDicts, "\",\"")+"\".") + "\n")
	sb.WriteString(helpStyle.Width(m.width-2).Render(" Or a filepath to dictionary files.") + "\n")
	sb.WriteString("\n")
	sb.WriteString(helpStyle.Width(m.width-2).Render("Use ↑/↓ keys to see previous entries.") + "\n")
	lines += 5

	sb.WriteString(padLines(m.height - lines - footerLines))
	return sb.String(), csr
}

func (v *viewSwitch) help() string {
	return back + ": Back"
}

var internalDicts = []string{
	"default",
	"csw24",
	"cas19",
	"enwiktionary",
}

func (v *viewSwitch) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	v.currentError = ""
	switch msg.String() {
	case enter:
		if s := v.input.value(); len(s) > 0 {
			return m.switchDictionary(v.input.value())
		}
	case back:
		m.restoreView(switchDict, v.backMode, v.backView)
	case up:
		if s := v.input.value(); len(s) > 0 {
			pvs := distinctValues(append(internalDicts, m.prefs.UsedDictionaries...))
			idx := slices.Index(pvs, s)
			if idx < 0 {
				idx = 0
			} else {
				idx--
				if idx < 0 {
					idx = len(pvs) - 1
				}
			}
			v.input.set(pvs[idx])
		} else {
			v.input.set(internalDicts[0])
		}
	case down:
		if s := v.input.value(); len(s) > 0 {
			pvs := distinctValues(append(internalDicts, m.prefs.UsedDictionaries...))
			idx := slices.Index(pvs, s) + 1
			if idx < 0 || idx >= len(pvs) {
				idx = 0
			}
			v.input.set(pvs[idx])
		} else {
			v.input.set(internalDicts[0])
		}
	default:
		if v.input != nil {
			v.input.key(msg)
		}
	}
	return nil
}

func distinctValues(vs []string) []string {
	result := make([]string, 0, len(vs))
	m := make(map[string]struct{}, len(vs))
	for _, v := range vs {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			result = append(result, v)
		}
	}
	slices.Sort(result)
	return result
}

func (v *viewSwitch) update(m *model, msg tea.Msg) tea.Cmd {
	if r, ok := msg.(switchDictionaryResult); ok {
		v.currentError = r.err.Error()
	}
	return nil
}

func (v *viewSwitch) wordLength() int {
	return 0
}

func (v *viewSwitch) currentWord() string {
	return ""
}

func (v *viewSwitch) show(backMode mode, backView view) {
	v.backMode = backMode
	v.backView = backView
	v.currentError = ""
	if v.input == nil {
		v.input = &stringInput{maxWidth: 10}
	}
	v.input.current = words.CurrentDictionary()
}

func (v *viewSwitch) paste(m *model, msg tea.PasteMsg) {
	if v.input != nil {
		v.input.set(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(msg.Content, "\t", ""), "\r", ""), "\n", ""))
	}
}
