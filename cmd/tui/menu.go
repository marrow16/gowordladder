package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"strings"
	"unicode/utf8"
)

const (
	menuChar = "≡"
)

type menu struct {
	showing                     bool
	overallWidth, overallHeight int
	onItem                      int
	items                       []menuItem
	itemClickKeys               map[int]tea.KeyPressMsg
}

type menuItem struct {
	text string
	key  string
}

var baseMenuItems = []menuItem{
	{text: "Play", key: ctrlPlay},
	{text: "Generate", key: ctrlGenerate},
	{text: "Solver", key: ctrlSolver},
	{text: "Word lookup", key: ctrlWord},
	{text: "Word distances", key: ctrlDistances},
	{},
	{text: "High Scores", key: ctrlHighScores},
	{text: "Switch Dictionary", key: fSwitch},
	{text: "Help", key: fHelp},
	{text: "Exit", key: exit},
}

func (m *menu) draw(sf layout.Surface, current view) bool {
	if m.showing {
		m.itemClickKeys = make(map[int]tea.KeyPressMsg, 0)
		m.items = baseMenuItems
		if viewItems := current.menu(); len(viewItems) > 0 {
			m.items = append(viewItems, menuItem{})
			m.items = append(m.items, baseMenuItems...)
		}
		maxWd := 0
		for _, item := range m.items {
			if item.text != "" {
				if l := utf8.RuneCountInString(item.text + item.key); l > maxWd {
					maxWd = l
				}
			}
		}
		maxWd += 6
		m.overallWidth, m.overallHeight = maxWd+2, len(m.items)+2
		sf.Fill(0, sf.Width()-m.overallWidth, m.overallHeight, m.overallWidth, headerStyle)
		sf.BoxRounded(0, sf.Width()-m.overallWidth, m.overallHeight, m.overallWidth, headerStyle)
		sf.Text(0, sf.Width()-2, menuChar, headerStyle)
		lpos := sf.Width() - m.overallWidth + 2
		for i, item := range m.items {
			useStyle := headerStyle
			if i == m.onItem {
				useStyle = headerSelectedStyle
			}
			if item.text == "" {
				sf.Text(i+1, lpos-2, "├"+strings.Repeat(horizontal, m.overallWidth-2)+"┤", headerStyle)
			} else if item.key != "" {
				m.itemClickKeys[i+1] = tea.KeyPressMsg{Text: item.key}
				if len(item.key) == 1 {
					sf.LineColumns(i+1, lpos, maxWd-2, layout.Runs{{Text: item.text}, {Text: "(\"" + item.key + "\")"}}, useStyle)
				} else {
					sf.LineColumns(i+1, lpos, maxWd-2, layout.Runs{{Text: item.text}, {Text: "(" + item.key + ")"}}, useStyle)
				}
			} else {
				sf.Text(i+1, lpos, item.text, useStyle)
			}
		}
	}
	return m.showing
}

func (m *menu) update(mo *model, canBack bool, msg tea.Msg) (result tea.Msg, handled bool) {
	result = msg
	if !m.showing {
		return
	}
	switch mt := msg.(type) {
	case tea.WindowSizeMsg, switchDictionaryResult:
		return msg, false
	case tea.MouseClickMsg:
		mmsg := mt.Mouse()
		switch {
		case mmsg.Button == tea.MouseLeft && mmsg.Y == 0 && mmsg.X < 3 && canBack:
			m.toggle()
			result = tea.KeyPressMsg{Text: back}
		case mmsg.Button == tea.MouseLeft && mmsg.Y == 0 && mmsg.X >= mo.width-2:
			m.toggle()
			handled = true
		case mmsg.Y >= m.overallHeight || mmsg.X < mo.width-m.overallWidth:
			// outside
			m.toggle()
			handled = true
		case mmsg.Y > 0 && mmsg.Y < m.overallHeight && mmsg.X < mo.width-1 && mmsg.X > mo.width-m.overallWidth:
			if k, ok := m.itemClickKeys[mmsg.Y]; ok {
				m.toggle()
				result = k
			} else {
				handled = true
			}
		default:
			handled = true
		}
	case tea.MouseWheelMsg:
		handled = true
		mmsg := mt.Mouse()
		switch mmsg.Button {
		case tea.MouseWheelUp:
			if m.onItem == -1 {
				m.onItem = len(m.items) - 1
			} else {
				m.onItem--
				if m.onItem < 0 {
					m.onItem = len(m.items) - 1
				} else if m.items[m.onItem].text == "" {
					m.onItem--
				}
			}
		case tea.MouseWheelDown:
			if m.onItem == -1 {
				m.onItem = 0
			} else {
				m.onItem++
				if m.onItem == len(m.items) {
					m.onItem = 0
				} else if m.items[m.onItem].text == "" {
					m.onItem++
				}
			}
		}
	case tea.KeyPressMsg:
		handled = true
		switch mt.String() {
		case "ctrl+c", exit, fHelp, fSwitch, ctrlWord, ctrlDistances, ctrlHighScores:
			m.toggle()
			handled = false
		case ctrlOpen:
			m.toggle()
		case enter:
			if m.onItem != -1 {
				result = tea.KeyPressMsg{Text: m.items[m.onItem].key}
				handled = false
				m.toggle()
			}
		case up:
			if m.onItem == -1 {
				m.onItem = len(m.items) - 1
			} else {
				m.onItem--
				if m.onItem < 0 {
					m.onItem = len(m.items) - 1
				} else if m.items[m.onItem].text == "" {
					m.onItem--
				}
			}
		case down:
			if m.onItem == -1 {
				m.onItem = 0
			} else {
				m.onItem++
				if m.onItem == len(m.items) {
					m.onItem = 0
				} else if m.items[m.onItem].text == "" {
					m.onItem++
				}
			}
		}
	}
	return result, handled
}

func (m *menu) toggle() {
	m.showing = !m.showing
	m.onItem = -1
}
