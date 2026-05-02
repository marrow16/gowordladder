package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"strings"
)

type input interface {
	render() (string, int)
	key(msg tea.KeyPressMsg)
	value() string
	set(val string)
}

type wordInput struct {
	maxLength int
	current   string
}

func (i *wordInput) render() (string, int) {
	l := len(i.current)
	cp := l
	if cp > i.maxLength {
		cp -= 2
	}
	var pad string
	if l < i.maxLength {
		pad = strings.Repeat(" ", i.maxLength-l)
	}
	return inputStyle.Render(i.current + pad), cp
}

func (i *wordInput) key(msg tea.KeyPressMsg) {
	k := strings.ToLower(msg.String())
	switch {
	case k == "backspace" && len(i.current) > 0:
		i.current = i.current[:len(i.current)-1]
	case len(k) == 1 && k >= "a" && k <= "z":
		if len(i.current) < i.maxLength {
			i.current += strings.ToUpper(k)
		} else {
			i.current = i.current[:len(i.current)-1] + strings.ToUpper(k)
		}
	}
}

func (i *wordInput) value() string {
	return i.current
}

func (i *wordInput) set(val string) {
	i.current = strings.ToUpper(val)
}

type numberInput struct {
	maxLength int
	current   string
}

func (i *numberInput) render() (string, int) {
	l := len(i.current)
	cp := l
	if cp > i.maxLength {
		cp -= 2
	}
	var pad string
	if l < i.maxLength {
		pad = strings.Repeat(" ", i.maxLength-l)
	}
	return inputStyle.Render(i.current + pad), cp
}

func (i *numberInput) key(msg tea.KeyPressMsg) {
	k := msg.String()
	switch {
	case k == "backspace" && len(i.current) > 0:
		i.current = i.current[:len(i.current)-1]
	case len(k) == 1 && k >= "0" && k <= "9":
		if len(i.current) < i.maxLength {
			i.current += k
		} else {
			i.current = i.current[:len(i.current)-1] + k
		}
	}
}

func (i *numberInput) value() string {
	return i.current
}

func (i *numberInput) set(val string) {
	i.current = val
}

var (
	inputStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ffffff")).
		Background(lipgloss.Color("#0000ff"))
)
