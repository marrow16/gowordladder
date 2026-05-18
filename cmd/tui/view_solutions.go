package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/solving"
	"os"
	"sort"
	"strconv"
	"strings"
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

const wordNumWidth = 3

func (v *viewSolutions) render(sf layout.Surface, m *model) *tea.Cursor {
	v.wordsDisplayed = make(wordPoints)
	if !v.showingAnalysis {
		v.solutionWidth = v.calculateSolutionWidth()
		across := (m.width - wordNumWidth) / v.solutionWidth
		solTot := strconv.Itoa(len(v.solutions))
		for s := 0; s < across && (s+v.offsetX) < len(v.solutions); s++ {
			sf.Text(0, (s*v.solutionWidth)+wordNumWidth+1, strconv.Itoa(s+v.offsetX+1)+"/"+solTot, helpStyle)
		}
		row := 1
		if v.offsetY == 0 {
			boxTop := topLeft + strings.Repeat(horizontal, v.wordLen) + topRight
			for s := 0; s < across && (s+v.offsetX) < len(v.solutions); s++ {
				sf.Text(row, (s*v.solutionWidth)+wordNumWidth, boxTop, helpStyle)
			}
			row++
		}
		maxLines := sf.Height() - row
		boxBottom := bottomLeft + strings.Repeat(horizontal, v.wordLen) + bottomRight
		for l := 0; l < maxLines && (l+v.offsetY) <= v.maxLadderLen; l++ {
			actualRow := l + v.offsetY
			if actualRow < v.maxLadderLen {
				sf.TextRight(l+row, 0, wordNumWidth-1, strconv.Itoa(actualRow+1), helpStyle)
			}
			for s := 0; s < across && (s+v.offsetX) < len(v.solutions); s++ {
				solution := v.solutions[s+v.offsetX]
				ladder := solution.Ladder()
				ladderLen := len(ladder)
				col := (s * v.solutionWidth) + wordNumWidth
				colWs := col + 1
				colEnd := colWs + v.wordLen
				switch {
				case actualRow == ladderLen:
					sf.Text(l+row, col, boxBottom, helpStyle)
				case actualRow == 0:
					sf.Text(l+row, col, vertical, helpStyle)
					sf.Text(l+row, colWs, ladder[actualRow].String())
					sf.Text(l+row, colEnd, vertical, helpStyle)
					v.wordsDisplayed.addWord(ladder[actualRow].String(), l+3, colWs)
				case actualRow < ladderLen:
					sf.Text(l+row, col, vertical, helpStyle)
					sf.Text(l+row, colWs, ladder[actualRow].String())
					sf.Text(l+row, colEnd, vertical, helpStyle)
					v.wordsDisplayed.addWord(ladder[actualRow].String(), l+3, colWs)
					prev := []rune(ladder[actualRow-1].String())
					word := []rune(ladder[actualRow].String())
					diffPos := 0
					for i, r := range word {
						if r != prev[i] {
							diffPos = i
							break
						}
					}
					sf.Text(l+row, colWs+diffPos, string(word[diffPos]), letterStyle)
				}
			}
		}
	} else {
		sf.Text(0, 1, fmt.Sprintf("Analysis of distinct words over %d solutions", len(v.solutions)))
		maxCount := 1
		for _, a := range v.analysis {
			if mx := len(a); mx > maxCount {
				maxCount = mx
			}
		}
		maxDigits := len(strconv.Itoa(maxCount)) + 1
		maxLines := sf.Height() - 1
		barWidth := m.width / 2
		for l := 0; l < maxLines && (l+v.offsetY) < v.maxLadderLen; l++ {
			row := l + v.offsetY
			sf.TextRight(l+1, 1, 3, strconv.Itoa(row+1)+":", helpStyle)
			count := len(v.analysis[row])
			sf.TextRight(l+1, 4, maxDigits, strconv.Itoa(count))
			sf.Block(l+1, maxDigits+5, (count*barWidth)/maxCount, '█', helpStyle)
		}
	}
	return nil
}

func (v *viewSolutions) helpLines() ([]string, *lipgloss.Style) {
	if v.showingAnalysis {
		return []string{"↑/↓: Scroll  •  " + back + ": Back"}, nil
	} else if len(v.solutions) > 1 {
		return []string{"←/→: Solutions  •  ↑/↓: Scroll  •  " + back + ": Back  •  " + ctrlAnalyse + ": Analyse"}, nil
	} else {
		return []string{"←/→: Solutions  •  ↑/↓: Scroll  •  " + back + ": Back"}, nil
	}
}

func (v *viewSolutions) menu() []menuItem {
	if !v.showingAnalysis {
		return []menuItem{
			{text: "Analyse", key: ctrlAnalyse},
		}
	}
	return nil
}

func (v *viewSolutions) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back, backspace:
		if v.showingAnalysis {
			v.showingAnalysis = false
		} else {
			m.restoreView(solutions, v.backMode, v.backView)
		}
	case ctrlAnalyse:
		if !v.showingAnalysis && len(v.solutions) > 1 {
			return v.analyse()
		}
	case ctrlExport:
		go v.export()
	case home:
		v.offsetX = 0
		v.offsetY = 0
	case end:
		if !v.showingAnalysis {
			v.offsetX = len(v.solutions) - 1
			if v.offsetX < 0 {
				v.offsetX = 0
			}
			v.offsetY = v.maxLadderLen + 5 - m.height
			if v.offsetY < 0 {
				v.offsetY = 0
			}
		}
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
				pgWd := ((m.width - wordNumWidth) / v.solutionWidth) - 1
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
				pgWd := ((m.width - wordNumWidth) / v.solutionWidth) - 1
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
		switch msg.Button {
		case tea.MouseLeft:
			return m.lookupWord(wd)
		case tea.MouseRight:
			return m.lookupDistance(wd)
		}
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
	width := v.wordLen + 2
	l := len(strconv.Itoa(len(v.solutions)))
	// check that "n/n" isn't bigger than the word length...
	if maxHdr := (l * 2) + 1; maxHdr > width {
		width = maxHdr
	}
	return width + 2
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
