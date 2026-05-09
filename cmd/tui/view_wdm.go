package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"gowordladder/words"
	"slices"
	"strconv"
	"strings"
)

type viewWordDistances struct {
	backMode         mode
	backView         view
	offsetY, offsetX int
	input            input
	distancesResult  *distancesResult
}

func (v *viewWordDistances) content(m *model) (string, *tea.Cursor) {
	const (
		prompt      = " Word: "
		footerLines = 2
	)
	var sb strings.Builder
	sb.WriteString("\n" + prompt)
	s, cxp := v.input.render()
	sb.WriteString(s)
	csr := tea.NewCursor(cxp+len(prompt), 2)
	sb.WriteString("\n" + strings.Repeat("─", m.width) + "\n")
	lines := 4
	if v.distancesResult != nil {
		if !v.distancesResult.inDictionary {
			sb.WriteString(" " + errorStyle.Render("Word not in my dictionary") + "\n")
			lines++
		} else {
			sb.WriteString(fmt.Sprintf(" Max ladder length: %d  Max words: %d (@%d)\n", v.distancesResult.maxLadderLength, v.distancesResult.maxWords, v.distancesResult.maxWordsAt))
			lines++
			colWd := v.distancesResult.wordLength + 2
			colWd += 2
			hFmt := " %" + strconv.Itoa(colWd-4) + "d   "
			lmWd := len(strconv.Itoa(v.distancesResult.maxWords))
			lmFmt := " %" + strconv.Itoa(lmWd) + "d "
			sb.WriteString(strings.Repeat(" ", lmWd+2))
			wd := lmWd + 2
			for c := 0; c < v.distancesResult.maxLadderLength && (c+v.offsetX) < v.distancesResult.maxLadderLength && (wd+colWd) < m.width; c++ {
				col := c + v.offsetX + 1
				sb.WriteString(helpStyle.Render(fmt.Sprintf(hFmt, col)))
				wd += colWd
			}
			sb.WriteString("\n")
			lines++
			if v.offsetY < 1 {
				sb.WriteString(strings.Repeat(" ", lmWd+1))
				wd = lmWd + 1
				for c := 0; c < v.distancesResult.maxLadderLength && (c+v.offsetX) < v.distancesResult.maxLadderLength && (wd+colWd) < m.width; c++ {
					sb.WriteString(" ")
					sb.WriteString(helpStyle.Render(topLeft + strings.Repeat(topBottom, v.distancesResult.wordLength) + topRight))
					sb.WriteString(" ")
					wd += colWd
				}
				sb.WriteString("\n")
				lines++
			}
			maxLines := m.height - lines - footerLines
			for ll := 0; ll < maxLines; ll++ {
				wn := ll + v.offsetY
				if wn < v.distancesResult.maxWords {
					sb.WriteString(helpStyle.Render(fmt.Sprintf(lmFmt, wn+1)))
				} else {
					sb.WriteString(strings.Repeat(" ", lmWd+2))
				}
				wd = lmWd + 1
				for c := 0; c < v.distancesResult.maxLadderLength && (c+v.offsetX) < v.distancesResult.maxLadderLength && (wd+colWd) < m.width; c++ {
					cll := c + v.offsetX + 1
					wds := v.distancesResult.distances[cll]
					switch {
					case wn < len(wds):
						sb.WriteString(helpStyle.Render(vertical))
						sb.WriteString(wds[wn])
						sb.WriteString(helpStyle.Render(vertical))
					case wn == len(wds):
						sb.WriteString(helpStyle.Render(bottomLeft + strings.Repeat(topBottom, v.distancesResult.wordLength) + bottomRight))
					default:
						sb.WriteString(strings.Repeat(" ", v.distancesResult.wordLength+2))
					}
					sb.WriteString("  ")
					wd += colWd
				}
				sb.WriteString("\n")
				lines++
			}
		}
	}
	sb.WriteString(padLines(m.height - lines - footerLines))
	return sb.String(), csr
}

func (v *viewWordDistances) help() string {
	return "enter: Lookup  •  " + back + ": Back"
}

func (v *viewWordDistances) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back:
		m.restoreView(v.backMode, v.backView)
		return nil
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case pageUp:
		v.offsetY -= m.height - 9
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case down:
		if v.distancesResult != nil && v.offsetY+15 < v.distancesResult.maxWords {
			v.offsetY++
		}
	case pageDown:
		v.offsetY += m.height - 9
		if v.distancesResult != nil && v.offsetY+15 >= v.distancesResult.maxWords {
			v.offsetY = v.distancesResult.maxWords - 15
			if v.offsetY < 0 {
				v.offsetY = 0
			}
		}
	case left:
		if v.offsetX > 0 {
			v.offsetX--
		}
	case pageLeft:
		if v.offsetX > 0 {
			v.offsetX -= 10
			if v.offsetX < 0 {
				v.offsetX = 0
			}
		}
	case right:
		if v.distancesResult != nil && v.offsetX+1 < v.distancesResult.maxLadderLength {
			v.offsetX++
		}
	case pageRight:
		if v.distancesResult != nil {
			v.offsetX += 10
			if v.offsetX >= v.distancesResult.maxLadderLength {
				v.offsetX = v.distancesResult.maxLadderLength - 1
			}
		}
	case enter:
		return v.doLookup()
	}
	if v.input.key(msg) {
		v.distancesResult = nil
	}
	return nil
}

func (v *viewWordDistances) update(m *model, msg tea.Msg) tea.Cmd {
	if dr, ok := msg.(distancesResult); ok {
		v.distancesResult = &dr
	}
	return nil
}

func (v *viewWordDistances) wordLength() int {
	if v.input != nil {
		return len(v.input.value())
	}
	return 0
}

func (v *viewWordDistances) currentWord() string {
	if v.input != nil {
		return v.input.value()
	}
	return ""
}

func (v *viewWordDistances) lookupWord(word string, backMode mode, backView view) tea.Cmd {
	v.offsetY = 0
	v.offsetX = 0
	v.distancesResult = nil
	v.backMode = backMode
	v.backView = backView
	v.input = &wordInput{maxLength: 15, current: word}
	return v.doLookup()
}

type distancesResult struct {
	inDictionary    bool
	wordLength      int
	maxLadderLength int
	maxWords        int
	maxWordsAt      int
	distances       map[int][]string
}

func (v *viewWordDistances) doLookup() tea.Cmd {
	s := v.input.value()
	if l := len(s); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			if wd, ok := dict.Word(s); ok {
				wdm := words.NewWordDistanceMap(wd, nil)
				d1 := make(map[int][]string, len(wdm))
				for w, dist := range wdm {
					if dist > 0 {
						d1[dist] = append(d1[dist], w)
					}
				}
				result := make(map[int][]string, len(d1))
				maxWords := 0
				maxWordsAt := 0
				for d, sl := range d1 {
					slices.Sort(sl)
					result[d] = sl
					if nw := len(sl); nw > maxWords {
						maxWords = nw
						maxWordsAt = d
					}
				}
				return distancesResult{
					inDictionary:    true,
					wordLength:      l,
					maxLadderLength: wd.MaxSteps(),
					maxWords:        maxWords,
					maxWordsAt:      maxWordsAt,
					distances:       result,
				}
			} else {
				return distancesResult{inDictionary: false}
			}
		}
	}
	return nil
}
