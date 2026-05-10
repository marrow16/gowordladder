package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"gowordladder/solving"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type solutionsView interface {
	view
	setSolutions(solutions []*solving.Solution, backMode mode, backView view)
}

type viewSolutions struct {
	backMode         mode
	backView         view
	offsetX, offsetY int
	wordLen          int
	maxLadderLen     int
	solutionWidth    int
	solutions        []*solving.Solution
	showingAnalysis  bool
	analysis         []map[string]struct{}
	wordsDisplayed   wordPoints
}

func (v *viewSolutions) content(m *model) (string, *tea.Cursor) {
	const (
		footerLines = 2
	)
	v.wordsDisplayed = make(wordPoints)
	var sb strings.Builder
	lines := 1
	if !v.showingAnalysis {
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
		maxLines := m.height - lines - footerLines
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
				x := 4 + (s * v.solutionWidth)
				if row == 0 {
					sb.WriteString(helpStyle.Render(vertical))
					sb.WriteString(ladder[row].String())
					sb.WriteString(helpStyle.Render(vertical))
					sb.WriteString(strings.Repeat(" ", v.solutionWidth-v.wordLen-2))
					v.wordsDisplayed.addWord(ladder[row].String(), lines-1, x)
				} else if row < len(ladder) {
					sb.WriteString(helpStyle.Render(vertical))
					prev := []rune(ladder[row-1].String())
					word := []rune(ladder[row].String())
					for i, r := range word {
						if r != prev[i] {
							sb.WriteString(letterStyle.Render(string(r)))
						} else {
							sb.WriteRune(r)
						}
					}
					sb.WriteString(helpStyle.Render(vertical))
					sb.WriteString(strings.Repeat(" ", v.solutionWidth-v.wordLen-2))
					v.wordsDisplayed.addWord(ladder[row].String(), lines-1, x)
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
	} else {
		sb.WriteString(fmt.Sprintf(" Analysis of distinct words over %d solutions", len(v.solutions)))
		lines++
		maxCount := 1
		for _, a := range v.analysis {
			if mx := len(a); mx > maxCount {
				maxCount = mx
			}
		}
		maxDigits := len(strconv.Itoa(maxCount)) + 1
		maxFmt := fmt.Sprintf("%%%dd ", maxDigits)
		maxLines := m.height - lines - footerLines
		for l := 0; l < maxLines && (l+v.offsetY) < v.maxLadderLen; l++ {
			row := l + v.offsetY
			sb.WriteString("\n")
			lines++
			sb.WriteString(helpStyle.Render(fmt.Sprintf("%2d: ", row+1)))
			count := len(v.analysis[row])
			sb.WriteString(fmt.Sprintf(maxFmt, count))
			const barWidth = 20
			sb.WriteString(helpStyle.Render(strings.Repeat("█", (count*barWidth)/maxCount)))
		}
	}
	sb.WriteString(padLines(m.height - lines - 1))
	return sb.String(), nil
}

func (v *viewSolutions) help() string {
	if v.showingAnalysis {
		return "↑/↓: Scroll  •  " + back + ": Back"
	} else if len(v.solutions) > 1 {
		return "←/→: Solutions  •  ↑/↓: Scroll  •  " + back + ": Back  •  " + ctrlAnalyse + ": Analyse"
	} else {
		return "←/→: Solutions  •  ↑/↓: Scroll  •  " + back + ": Back"
	}
}

func (v *viewSolutions) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back, backspace:
		if v.showingAnalysis {
			v.showingAnalysis = false
		} else {
			m.restoreView(v.backMode, v.backView)
		}
	case ctrlAnalyse:
		if !v.showingAnalysis && len(v.solutions) > 1 {
			return v.analyse()
		}
	case ctrlExport:
		go v.export()
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case pageUp:
		v.offsetY -= m.height
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case down:
		maxOffsetY := v.maxLadderLen + 5 - m.height
		if maxOffsetY < 0 {
			maxOffsetY = 0
		} else if v.offsetY < maxOffsetY {
			v.offsetY++
		}
	case pageDown:
		v.offsetY += m.height
		maxOffsetY := v.maxLadderLen + 5 - m.height
		if maxOffsetY < 0 {
			maxOffsetY = 0
		} else if v.offsetY > maxOffsetY {
			v.offsetY = maxOffsetY
		}
	case left:
		if !v.showingAnalysis {
			if v.offsetX > 0 {
				v.offsetX--
			}
		}
	case pageLeft:
		if !v.showingAnalysis {
			if v.offsetX > 0 && v.solutionWidth > 0 {
				pgWd := (m.width / v.solutionWidth) - 1
				if l := v.offsetX - pgWd; l >= 0 {
					v.offsetX = l
				} else {
					v.offsetX = 0
				}
			}
		}
	case right:
		if !v.showingAnalysis {
			if v.offsetX < len(v.solutions)-1 {
				v.offsetX++
			}
		}
	case pageRight:
		if !v.showingAnalysis {
			if v.solutionWidth > 0 {
				pgWd := (m.width / v.solutionWidth) - 1
				if l := v.offsetX + pgWd; l < len(v.solutions) {
					v.offsetX = l
				} else {
					v.offsetX = len(v.solutions) - 1
				}
			}
		}
	}
	return nil
}

func (v *viewSolutions) click(m *model, msg tea.Mouse) tea.Cmd {
	if wd, ok := v.wordsDisplayed[pt{msg.Y, msg.X}]; ok {
		return m.lookupWord(wd)
	}
	return nil
}

func (v *viewSolutions) analyse() tea.Cmd {
	return func() tea.Msg {
		analysis := make([]map[string]struct{}, v.maxLadderLen)
		for i := range v.maxLadderLen {
			analysis[i] = make(map[string]struct{})
		}
		analysis[0] = map[string]struct{}{"": {}}
		for _, solution := range v.solutions {
			ladder := solution.Ladder()
			for i := 1; i < len(ladder); i++ {
				analysis[i][ladder[i].String()] = struct{}{}
			}
		}
		return analysisResult{analysis: analysis}
	}
}

func (v *viewSolutions) export() {
	if len(v.solutions) > 0 {
		ladder := v.solutions[0].Ladder()
		fn := fmt.Sprintf("solutions-%s-%s.csv", ladder[0], ladder[len(ladder)-1])
		if f, err := os.Create(fn); err == nil {
			defer f.Close()
			for i, solution := range v.solutions {
				ladder = solution.Ladder()
				_, _ = fmt.Fprintf(f, "%d,%d", i+1, len(ladder))
				for _, w := range ladder {
					_, _ = fmt.Fprint(f, ",")
					_, _ = fmt.Fprint(f, w.String())
				}
				_, _ = f.WriteString("\n")
			}
		}
	}
}

type analysisResult struct {
	analysis []map[string]struct{}
}

func (v *viewSolutions) update(m *model, msg tea.Msg) tea.Cmd {
	if ar, ok := msg.(analysisResult); ok && !v.showingAnalysis {
		v.analysis = ar.analysis
		v.showingAnalysis = true
	}
	return nil
}

func (v *viewSolutions) wordLength() int {
	return v.wordLen
}

func (v *viewSolutions) currentWord() string {
	return v.solutions[v.offsetX].Ladder()[0].String()
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

func (v *viewSolutions) setSolutions(solutions []*solving.Solution, backMode mode, backView view) {
	v.backMode = backMode
	v.backView = backView
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
	if len(solutions) < 500_000 {
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
}

var letterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
