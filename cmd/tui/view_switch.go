package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/words"
	"slices"
	"strings"
)

type switchView interface {
	view
	viewShowable
}

type viewSwitch struct {
	backMode     mode
	backView     view
	input        *stringInput
	currentError string
}

func (v *viewSwitch) render(sf layout.Surface, m *model) *tea.Cursor {
	const (
		prompt    = "Source:"
		promptLen = len(prompt)
	)
	sf.TextRight(1, 0, promptLen+1, prompt)
	iw := m.width - promptLen - 3
	sf.TextFixed(1, promptLen+2, iw, v.input.value(), inputStyle)
	v.input.maxWidth = iw
	row := 3
	if v.currentError != "" {
		row += sf.TextWrapped(2, promptLen+2, iw, v.currentError, errorStyle)
	}
	sf.TextCenter(row, 1, m.width-1, "Press "+enter+" to switch.", helpStyle)
	sf.TextCenter(row+1, 1, m.width-2, "Use internal: \""+strings.Join(internalDicts, "\",\"")+"\".", helpStyle)
	sf.TextCenter(row+2, 1, m.width-2, "Or a filepath to dictionary files.", helpStyle)
	return tea.NewCursor(v.input.cursorPos()+promptLen+2, 2)
}

func (v *viewSwitch) helpLines() ([]string, *lipgloss.Style) {
	return []string{back + ": Back"}, nil
}

func (v *viewSwitch) menu() []menuItem {
	return nil
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
