package main

import (
	tea "charm.land/bubbletea/v2"
	"strconv"
	"strings"
)

type input interface {
	render() (string, int)
	key(msg tea.KeyPressMsg) bool
	value() string
	set(val string)
	paste(msg tea.PasteMsg)
}

type wordInput struct {
	maxLength        int
	current          string
	allowUnderscores bool
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

func (i *wordInput) key(msg tea.KeyPressMsg) bool {
	k := strings.ToLower(msg.String())
	switch {
	case k == backspace && len(i.current) > 0:
		i.current = i.current[:len(i.current)-1]
		return true
	case len(k) == 1 && ((i.allowUnderscores && k == "_") || (k >= "a" && k <= "z")):
		if len(i.current) < i.maxLength {
			i.current += strings.ToUpper(k)
		} else {
			i.current = i.current[:len(i.current)-1] + strings.ToUpper(k)
		}
		return true
	}
	return false
}

func (i *wordInput) value() string {
	return i.current
}

func (i *wordInput) set(val string) {
	i.current = strings.ToUpper(val)
}

func (i *wordInput) paste(msg tea.PasteMsg) {
	s := strings.ToUpper(msg.Content)
	if isAllAZ(s) {
		if len(s) > i.maxLength {
			i.current = s[:i.maxLength]
		} else {
			i.current = s
		}
	}
}

func isAllAZ(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
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

func (i *numberInput) key(msg tea.KeyPressMsg) bool {
	k := msg.String()
	switch {
	case k == backspace && len(i.current) > 0:
		i.current = i.current[:len(i.current)-1]
		return true
	case k == up:
		if i.current == "" {
			i.current = "1"
			return true
		} else if n, err := strconv.Atoi(i.current); err == nil {
			if s := strconv.Itoa(n + 1); len(s) <= i.maxLength {
				i.current = s
				return true
			}
		}
	case k == down:
		if n, err := strconv.Atoi(i.current); err == nil && n > 0 {
			if s := strconv.Itoa(n - 1); len(s) <= i.maxLength {
				i.current = s
				return true
			}
		}
	case len(k) == 1 && k >= "0" && k <= "9":
		if len(i.current) < i.maxLength {
			i.current += k
		} else {
			i.current = i.current[:len(i.current)-1] + k
		}
		return true
	}
	return false
}

func (i *numberInput) value() string {
	return i.current
}

func (i *numberInput) set(val string) {
	i.current = val
}

func (i *numberInput) paste(msg tea.PasteMsg) {
	if isAllDigits(msg.Content) {
		if len(msg.Content) > i.maxLength {
			i.current = msg.Content[:i.maxLength]
		} else {
			i.current = msg.Content
		}
	}
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
