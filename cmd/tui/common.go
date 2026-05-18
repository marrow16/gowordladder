package main

import "charm.land/lipgloss/v2"

type pt [2]int //Y,X
// wordPoints is a map of words displayed on screen (for clicking)
type wordPoints map[pt]string

func (wp wordPoints) addWord(word string, y, x int) {
	for l := 0; l <= len(word); l++ {
		wp[pt{y, x + l}] = word
	}
}

// common styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#0000ff")).
			AlignHorizontal(lipgloss.Center)
	headerSelectedStyle = headerStyle.Background(lipgloss.Color("#8080ff"))
	helpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).
				AlignHorizontal(lipgloss.Center)
	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#008000"))
	errorStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("#ff0000"))
	inputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#0000ff"))
	boldStyle       = lipgloss.NewStyle().Bold(true)
	playCursorColor = lipgloss.Color("#ccccff")
	hintStyle       = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#008000"))
	warningStyle    = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#ff7f00"))
	wrongStyle      = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#ff0000"))
)
