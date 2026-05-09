package main

import "charm.land/lipgloss/v2"

// box characters
const (
	topLeft     = "╭"
	topRight    = "╮"
	topBottom   = "─"
	bottomLeft  = "╰"
	bottomRight = "╯"
	vertical    = "│"
)

// common keys
const (
	enter          = "enter"
	tab            = "tab"
	shiftTab       = "shift+tab"
	back           = "ctrl+b"
	backspace      = "backspace"
	up             = "up"
	down           = "down"
	left           = "left"
	right          = "right"
	pageLeft       = "shift+left"
	pageRight      = "shift+right"
	pageUp         = "pgup"
	pageDown       = "pgdown"
	ctrlWord       = "ctrl+w"
	ctrlGenerate   = "ctrl+g"
	ctrlSolver     = "ctrl+s"
	ctrlDistances  = "ctrl+d"
	ctrlHighScores = "ctrl+t"
	ctrlNew        = "ctrl+n"
	ctrlPlay       = "ctrl+p"
	ctrlHelp       = "ctrl+h"
	ctrlFill       = "ctrl+f"
	ctrlAnalyse    = "ctrl+a"
	exit           = "esc"
	fHelp          = "f1"
)

// common styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#0000ff")).
			AlignHorizontal(lipgloss.Center)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa")).
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
