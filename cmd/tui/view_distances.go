package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/marrow16/gowordladder/words"
	"slices"
	"strconv"
	"strings"
)

type viewDistances struct {
	backMode         mode
	backView         view
	offsetY, offsetX int
	input            input
	distancesResult  *distancesResult
	wordsDisplayed   wordPoints
}

func (v *viewDistances) content(m *model) (string, *tea.Cursor) {
	const (
		prompt      = " Word: "
		footerLines = 3
	)
	v.wordsDisplayed = make(wordPoints)
	var sb strings.Builder
	sb.Grow(m.height * m.width)
	sb.WriteString("\n" + prompt)
	s, cxp := v.input.render()
	sb.WriteString(s)
	csr := tea.NewCursor(cxp+len(prompt), 2)
	sb.WriteString("\n" + helpStyle.Render(strings.Repeat(horizontal, m.width)) + "\n")
	lines := 4
	if v.distancesResult != nil {
		if !v.distancesResult.inDictionary {
			sb.WriteString(" " + errorStyle.Render("Word not in my dictionary") + "\n")
			lines++
		} else {
			skipDisplay := false
			if v.distancesResult.mode == distancesNormal {
				switch v.distancesResult.maxLadderLength {
				case 1:
					sb.WriteString(" " + boldStyle.Render(v.distancesResult.word) + helpStyle.Render(" is an island word (no ladders)\n"))
					skipDisplay = true
				case 2:
					sb.WriteString(" " + boldStyle.Render(v.distancesResult.word) + fmt.Sprintf(" Max ladder length: %d\n", v.distancesResult.maxLadderLength))
				default:
					sb.WriteString(" " + boldStyle.Render(v.distancesResult.word) + fmt.Sprintf(" Max ladder length: %d  Max words: %d (@%d)\n", v.distancesResult.maxLadderLength, v.distancesResult.maxWords, v.distancesResult.maxWordsAt))
				}
			} else if len(v.distancesResult.distances[1]) == 0 {
				skipDisplay = true
				sb.WriteString(" " + boldStyle.Render(v.distancesResult.word+": ") + helpStyle.Render(fmt.Sprintf("None for word length %d\n", v.distancesResult.wordLength)))
			} else {
				sb.WriteString(" " + boldStyle.Render(v.distancesResult.word+": ") + highlightStyle.Render(commas(len(v.distancesResult.distances[1]))) + fmt.Sprintf(" words (for word length %d)\n", v.distancesResult.wordLength))
			}
			lines++
			if !skipDisplay {
				colWd := v.distancesResult.wordLength + 2
				colWd += 2
				hFmt := " %" + strconv.Itoa(colWd-4) + "d   "
				lmWd := len(strconv.Itoa(v.distancesResult.maxWords))
				lmFmt := " %" + strconv.Itoa(lmWd) + "d "
				sb.WriteString(strings.Repeat(" ", lmWd+2))
				wd := lmWd + 2
				if v.distancesResult.mode == distancesNormal {
					for c := 0; c < v.distancesResult.maxLadderLength && (c+v.offsetX) < v.distancesResult.maxLadderLength && (wd+colWd) < m.width; c++ {
						col := c + v.offsetX + 1
						sb.WriteString(helpStyle.Render(fmt.Sprintf(hFmt, col)))
						wd += colWd
					}
				} else {
					sb.WriteString(v.distancesResult.word)
				}
				sb.WriteString("\n")
				lines++
				if v.offsetY < 1 {
					sb.WriteString(strings.Repeat(" ", lmWd+1))
					wd = lmWd + 1
					for c := 0; c < v.distancesResult.maxLadderLength && (c+v.offsetX) < v.distancesResult.maxLadderLength && (wd+colWd) < m.width; c++ {
						sb.WriteString(" ")
						sb.WriteString(helpStyle.Render(topLeft + strings.Repeat(horizontal, v.distancesResult.wordLength) + topRight))
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
							v.wordsDisplayed.addWord(wds[wn], lines, wd+2)
						case wn == len(wds):
							sb.WriteString(helpStyle.Render(bottomLeft + strings.Repeat(horizontal, v.distancesResult.wordLength) + bottomRight))
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
	}
	sb.WriteString(padLines(m.height - lines - footerLines))
	return sb.String(), csr
}

func (v *viewDistances) help() string {
	//	return "enter: Lookup  •  " + back + ": Back"
	return "enter: Lookup  •  1: Islands  •  2: Doublets\n" + back + ": Back"
}

func (v *viewDistances) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "1":
		return v.doIslands()
	case "2":
		return v.doDoublets()
	case back:
		m.restoreView(distances, v.backMode, v.backView)
		return nil
	case home:
		v.offsetY = 0
		v.offsetX = 0
	case end:
		if v.distancesResult != nil {
			v.offsetY = v.distancesResult.maxWords + 10 - m.height
			if v.offsetY < 0 {
				v.offsetY = 0
			}
			v.offsetX = v.distancesResult.maxWordsAt - 1
			if v.offsetX < 0 {
				v.offsetX = 0
			}
		}
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
	default:
		v.input.key(msg)
	}
	return nil
}

func (v *viewDistances) paste(m *model, msg tea.PasteMsg) {
	if v.input != nil {
		v.input.paste(msg)
	}
}

func (v *viewDistances) click(m *model, msg tea.Mouse) tea.Cmd {
	if msg.Y == 2 && msg.X >= 7 && msg.X <= 22 {
		if s := v.input.value(); len(s) > 0 {
			return m.lookupWord(s)
		}
		return nil
	} else if msg.Y == 4 && msg.X > 0 && v.distancesResult != nil && v.distancesResult.mode == distancesNormal {
		if msg.X < v.distancesResult.wordLength+1 {
			return m.lookupWord(v.distancesResult.word)
		} else {
			v.offsetX = v.distancesResult.maxWordsAt - 1
			if v.offsetX < 0 {
				v.offsetX = 0
			}
		}
		return nil
	}
	if wd, ok := v.wordsDisplayed[pt{msg.Y, msg.X}]; ok {
		switch msg.Button {
		case tea.MouseLeft:
			return v.doLookupWord(wd)
		case tea.MouseRight:
			return m.lookupWord(wd)
		}
	}
	return nil
}

func (v *viewDistances) update(m *model, msg tea.Msg) tea.Cmd {
	if dr, ok := msg.(distancesResult); ok {
		if dr.mode == distancesNormal {
			v.input.set(dr.word)
		}
		v.distancesResult = &dr
	}
	return nil
}

func (v *viewDistances) wordLength() int {
	if v.input != nil {
		return len(v.input.value())
	}
	return 0
}

func (v *viewDistances) currentWord() string {
	if v.input != nil {
		return v.input.value()
	}
	return ""
}

func (v *viewDistances) lookupWord(word string, backMode mode, backView view) tea.Cmd {
	if v.distancesResult != nil && v.distancesResult.word == word {
		return nil
	}
	v.offsetY = 0
	v.offsetX = 0
	v.distancesResult = nil
	v.backMode = backMode
	v.backView = backView
	v.input = &wordInput{maxLength: 15, current: word}
	return v.doLookup()
}

type distancesMode int

const (
	distancesNormal distancesMode = iota
	distancesIslands
	distancesDoublets
)

type distancesResult struct {
	word            string
	inDictionary    bool
	wordLength      int
	maxLadderLength int
	maxWords        int
	maxWordsAt      int
	distances       map[int][]string
	mode            distancesMode
}

func (v *viewDistances) doLookup() tea.Cmd {
	return v.doLookupWord(v.input.value())
}

func (v *viewDistances) doLookupWord(s string) tea.Cmd {
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
					word:            s,
					inDictionary:    true,
					wordLength:      l,
					maxLadderLength: wd.MaxSteps(),
					maxWords:        maxWords,
					maxWordsAt:      maxWordsAt,
					distances:       result,
				}
			} else {
				return distancesResult{word: s, inDictionary: false}
			}
		}
	}
	return nil
}

func (v *viewDistances) doIslands() tea.Cmd {
	if l := len(v.input.value()); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			islands := make([]string, len(dict.Islands))
			for i, island := range dict.Islands {
				islands[i] = island.String()
			}
			return distancesResult{
				word:            "Islands",
				inDictionary:    true,
				wordLength:      l,
				maxLadderLength: 1,
				maxWords:        len(islands),
				maxWordsAt:      1,
				distances: map[int][]string{
					1: islands,
				},
				mode: distancesIslands,
			}
		}
	}
	return nil
}

func (v *viewDistances) doDoublets() tea.Cmd {
	if l := len(v.input.value()); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			doublets := make([]string, len(dict.Doublets))
			for i, doublet := range dict.Doublets {
				doublets[i] = doublet.String()
			}
			return distancesResult{
				word:            "Doublets",
				inDictionary:    true,
				wordLength:      l,
				maxLadderLength: 1,
				maxWords:        len(doublets),
				maxWordsAt:      1,
				distances: map[int][]string{
					1: doublets,
				},
				mode: distancesDoublets,
			}
		}
	}
	return nil
}
