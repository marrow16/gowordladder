package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"gowordladder/solving"
	"sort"
	"strings"
	"unicode/utf8"
)

type solutionsView interface {
	view
	setSolutions(solutions []*solving.Solution)
}

type viewSolutions struct {
	offsetX, offsetY int
	wordLen          int
	maxLadderLen     int
	solutionWidth    int
	solutions        []*solving.Solution
}

func (v *viewSolutions) content(m *model) (string, *tea.Cursor) {
	const (
		topLeft     = "╭"
		topRight    = "╮"
		topBottom   = "─"
		bottomLeft  = "╰"
		bottomRight = "╯"
	)
	var sb strings.Builder
	lines := 1
	v.solutionWidth = v.calculateSolutionWidth()

	numSolutions := m.width / v.solutionWidth
	lines++
	sb.WriteString(" ")
	for s := 0; s < numSolutions && (s+v.offsetX) < len(v.solutions); s++ {
		hdr := fmt.Sprintf("  %d/%d", s+v.offsetX+1, len(v.solutions))
		sb.WriteString(helpStyle.Render(hdr))
		sb.WriteString(strings.Repeat(" ", v.solutionWidth-len(hdr)))
	}
	if v.offsetY == 0 {
		sb.WriteString("\n ")
		lines++
		hdr := "  " + topLeft + strings.Repeat(topBottom, v.wordLen) + topRight
		hdr += strings.Repeat(" ", v.solutionWidth-utf8.RuneCountInString(hdr))
		for s := 0; s < numSolutions && (s+v.offsetX) < len(v.solutions); s++ {
			sb.WriteString(helpStyle.Render(hdr))
		}
	}
	maxLines := m.height - lines - 2
	for l := 0; l < maxLines && (l+v.offsetY) <= v.maxLadderLen; l++ {
		sb.WriteString("\n")
		row := l + v.offsetY
		if row == v.maxLadderLen {
			sb.WriteString("   ")
		} else {
			sb.WriteString(helpStyle.Render(fmt.Sprintf("%2d ", v.offsetY+l+1)))
		}
		lines++
		for s := 0; s < numSolutions && (s+v.offsetX) < len(v.solutions); s++ {
			solution := v.solutions[s+v.offsetX]
			ladder := solution.Ladder()
			if row == 0 {
				sb.WriteString(helpStyle.Render("|"))
				sb.WriteString(ladder[row].String())
				sb.WriteString(helpStyle.Render("|"))
				sb.WriteString(strings.Repeat(" ", v.solutionWidth-v.wordLen-2))
			} else if row < len(ladder) {
				sb.WriteString(helpStyle.Render("|"))
				prev := []rune(ladder[row-1].String())
				word := []rune(ladder[row].String())
				for i, r := range word {
					if r != prev[i] {
						sb.WriteString(letterStyle.Render(string(r)))
					} else {
						sb.WriteRune(r)
					}
				}
				sb.WriteString(helpStyle.Render("|"))
				sb.WriteString(strings.Repeat(" ", v.solutionWidth-v.wordLen-2))
			} else if row == len(ladder) {
				sb.WriteString(helpStyle.Render(bottomLeft))
				sb.WriteString(helpStyle.Render(strings.Repeat(topBottom, v.wordLen)))
				sb.WriteString(helpStyle.Render(bottomRight))
				sb.WriteString(strings.Repeat(" ", v.solutionWidth-v.wordLen-2))
			} else {
				sb.WriteString(strings.Repeat(" ", v.solutionWidth))
			}
		}
	}

	sb.WriteString(strings.Repeat("\n", m.height-lines-1))
	return sb.String(), nil
}

func (v *viewSolutions) help() string {
	return "←/→: Solutions  •  ↑/↓: Scroll  •  backspace: Back"
}

func (v *viewSolutions) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "backspace":
		m.back()
	case "up":
		v.scrollUp(m)
	case "shift+up":
		v.offsetY = 0
	case "down":
		v.scrollDown(m)
	case "left":
		v.panLeft(m)
	case "shift+left":
		v.pageLeft(m)
	case "right":
		v.panRight(m)
	case "shift+right":
		v.pageRight(m)
	}
	return nil
}

func (v *viewSolutions) scrollUp(m *model) {
	if v.offsetY > 0 {
		v.offsetY--
	}
}

func (v *viewSolutions) scrollDown(m *model) {
	if v.offsetY < v.maxLadderLen && (v.offsetY+m.height) < (v.maxLadderLen+5) {
		v.offsetY++
	}
}

func (v *viewSolutions) panLeft(m *model) {
	if v.offsetX > 0 {
		v.offsetX--
	}
}

func (v *viewSolutions) pageLeft(m *model) {
	if v.offsetX > 0 && v.solutionWidth > 0 {
		pgWd := (m.width / v.solutionWidth) - 1
		if l := v.offsetX - pgWd; l >= 0 {
			v.offsetX = l
		} else {
			v.offsetX = 0
		}
	}
}

func (v *viewSolutions) panRight(m *model) {
	if v.offsetX < len(v.solutions)-1 {
		v.offsetX++
	}
}

func (v *viewSolutions) pageRight(m *model) {
	if v.solutionWidth > 0 {
		pgWd := (m.width / v.solutionWidth) - 1
		if l := v.offsetX + pgWd; l < len(v.solutions) {
			v.offsetX = l
		} else {
			v.offsetX = len(v.solutions) - 1
		}
	}
}

func (v *viewSolutions) update(m *model, msg tea.Msg) tea.Cmd {
	return nil
}

func (v *viewSolutions) wordLength() int {
	return v.wordLen
}

func (v *viewSolutions) calculateSolutionWidth() int {
	width := v.wordLen + 6
	l := len(v.solutions)
	maxHdr := fmt.Sprintf("  %d/%d  ", l, l)
	if len(maxHdr) > width {
		width = len(maxHdr)
	}
	return width
}

func (v *viewSolutions) setSolutions(solutions []*solving.Solution) {
	v.solutions = solutions
	v.solutionWidth = v.calculateSolutionWidth()
	sortSolutions(v.solutions)
	v.maxLadderLen = len(v.solutions[len(v.solutions)-1].Ladder())
	if len(v.solutions) > 0 {
		v.wordLen = len(v.solutions[0].Ladder()[0].String())
	} else {
		v.wordLen = 0
	}
	v.offsetX, v.offsetY = 0, 0
}

func sortSolutions(solutions []*solving.Solution) {
	sort.Slice(solutions, func(i, j int) bool {
		if len(solutions[i].Ladder()) < len(solutions[j].Ladder()) {
			return true
		} else if len(solutions[i].Ladder()) == len(solutions[j].Ladder()) {
			for idx, w := range solutions[i].Ladder() {
				if w.String() < solutions[j].Ladder()[idx].String() {
					return true
				} else if w.String() > solutions[j].Ladder()[idx].String() {
					return false
				}
			}
		}
		return false
	})
}

var letterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
