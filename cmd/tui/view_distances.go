package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"fmt"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/words"
	"os"
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
	scrollbar        layout.Scrollbar
}

func (v *viewDistances) render(sf layout.Surface, m *model) *tea.Cursor {
	const (
		prompt    = "Word:"
		promptLen = len(prompt)
		inputLen  = 15
	)
	v.scrollbar = nil
	v.wordsDisplayed = make(wordPoints)
	sf.TextRight(1, 0, promptLen+1, prompt)
	sf.TextFixed(1, promptLen+2, inputLen, v.input.value(), inputStyle)
	cp := v.input.cursorPos()
	sf.Block(2, 0, m.width, '─', helpStyle)
	if v.distancesResult != nil {
		if !v.distancesResult.inDictionary {
			if !v.distancesResult.hideError {
				sf.Text(1, promptLen+inputLen+3, "Not in my dictionary", errorStyle)
			}
		} else {
			header := layout.NewRuns(v.distancesResult.word, boldStyle)
			switch v.distancesResult.mode {
			case distancesNormal:
				switch v.distancesResult.maxLadderLength {
				case 1:
					sf.TextRun(3, 1, header.Add(" is an island word (no ladders)"))
				case 2:
					sf.TextRun(3, 1, header.Add(" Max ladder length: "+strconv.Itoa(v.distancesResult.maxLadderLength)))
					v.showDistances(m, sf)
				default:
					sf.TextRun(3, 1, header.Add(" Max ladder length: "+strconv.Itoa(v.distancesResult.maxLadderLength)+"  Max words: "+strconv.Itoa(v.distancesResult.maxWords)+" (@"+strconv.Itoa(v.distancesResult.maxWordsAt)+")"))
					v.showDistances(m, sf)
				}
			case distancesAnalysis:
				switch v.distancesResult.maxLadderLength {
				case 1:
					sf.TextRun(3, 1, header.Add(" is an island word (no ladders)"))
				case 2:
					sf.TextRun(3, 1, header.Add(" Max ladder length: "+strconv.Itoa(v.distancesResult.maxLadderLength)))
					v.showAnalysis(m, sf)
				default:
					sf.TextRun(3, 1, header.Add(" Max ladder length: "+strconv.Itoa(v.distancesResult.maxLadderLength)+"  Max words: "+strconv.Itoa(v.distancesResult.maxWords)+" (@"+strconv.Itoa(v.distancesResult.maxWordsAt)+")"))
					v.showAnalysis(m, sf)
				}
			case distancesLongests:
				sf.TextRun(3, 1, header.Add(": ").Add(commas(len(v.distancesResult.distances[1])), highlightStyle).Add(" words (word length "+strconv.Itoa(v.distancesResult.wordLength)+")"))
				v.showDistances(m, sf)
			case distancesOverall:
				if v.distancesResult.wordLength > 0 {
					sf.TextRun(3, 1, header.Add(": Word length "+strconv.Itoa(v.distancesResult.wordLength)))
				} else {
					sf.TextRun(3, 1, header)
				}
				v.showAnalysis(m, sf)
			default:
				if len(v.distancesResult.distances[1]) == 0 {
					sf.TextRun(3, 1, header.Add(": None for word length "+strconv.Itoa(v.distancesResult.wordLength)))
				} else if v.distancesResult.totalWords > 0 {
					perc := (float64(len(v.distancesResult.distances[1])) / float64(v.distancesResult.totalWords)) * 100.0
					sf.TextRun(3, 1, header.Add(": ").Add(commas(len(v.distancesResult.distances[1])), highlightStyle).Add(" words "+strconv.FormatFloat(perc, 'f', 1, 64)+"% (word length "+strconv.Itoa(v.distancesResult.wordLength)+")"))
					v.showDistances(m, sf)
				} else {
					sf.TextRun(3, 1, header.Add(": ").Add(commas(len(v.distancesResult.distances[1])), highlightStyle).Add(" words (word length "+strconv.Itoa(v.distancesResult.wordLength)+")"))
					v.showDistances(m, sf)
				}
			}
		}
	}
	return tea.NewCursor(cp+promptLen+2, 2)
}

func (v *viewDistances) showAnalysis(m *model, sf layout.Surface) {
	rgn := sf.Region(4, 0, sf.Height(), sf.Width())
	maxLines := rgn.Height()
	maxWords := v.distancesResult.maxWords
	maxNumWidth := len(strconv.Itoa(maxWords))
	barWidth := rgn.Width() - maxNumWidth - 10
	for l := 0; l < maxLines && l+v.offsetY < v.distancesResult.maxLadderLength; l++ {
		rgn.TextRight(l, 1, 4, strconv.Itoa(l+v.offsetY+1)+":", helpStyle)
		wordsCount := len(v.distancesResult.distances[l+v.offsetY+1])
		rgn.TextRight(l, 6, maxNumWidth, strconv.Itoa(wordsCount))
		rgn.Block(l, 7+maxNumWidth, (wordsCount*barWidth)/maxWords, '█', helpStyle)
	}
	if v.distancesResult.maxLadderLength > rgn.Height() {
		v.scrollbar = layout.NewVerticalScrollbar(v.scrollHandler(m))
		v.scrollbar.Draw(rgn, v.distancesResult.maxLadderLength, v.offsetY)
	}
}

func (v *viewDistances) showDistances(m *model, sf layout.Surface) {
	columnWidth := v.distancesResult.wordLength + 4
	numWidth := len(strconv.Itoa(v.distancesResult.maxWords))
	across := (m.width - numWidth - 3) / columnWidth
	boxTop := topLeft + strings.Repeat(horizontal, v.distancesResult.wordLength) + topRight
	wdY := 5
	showTop := v.offsetY == 0
	if v.distancesResult.mode == distancesNormal {
		for l := 0; l < across && l+v.offsetX < v.distancesResult.maxLadderLength; l++ {
			sf.Text(4, (l*columnWidth)+numWidth+3, strconv.Itoa(l+v.offsetX+1), helpStyle)
			if showTop {
				sf.Text(5, (l*columnWidth)+numWidth+2, boxTop, helpStyle)
			}
		}
	} else if showTop {
		wdY--
		sf.Text(4, numWidth+2, boxTop, helpStyle)
	}
	if showTop {
		wdY++
	}
	rgn := sf.Region(wdY, 0, sf.Height(), sf.Width())
	maxLines := rgn.Height()
	boxBottom := bottomLeft + strings.Repeat(horizontal, v.distancesResult.wordLength) + bottomRight
	for l := 0; l < maxLines; l++ {
		wn := l + v.offsetY
		if wn < v.distancesResult.maxWords {
			rgn.TextRight(l, 1, numWidth, strconv.Itoa(wn+1), helpStyle)
		}
		for s := 0; s < across && s+v.offsetX < v.distancesResult.maxLadderLength; s++ {
			wds := v.distancesResult.distances[s+v.offsetX+1]
			col := (s * columnWidth) + numWidth + 2
			colW := col + 1
			colE := colW + v.distancesResult.wordLength
			switch {
			case wn < len(wds):
				rgn.Text(l, col, vertical, helpStyle)
				v.wordsDisplayed.add(rgn.Text(l, colW, wds[wn]))
				rgn.Text(l, colE, vertical, helpStyle)
			case wn == len(wds):
				rgn.Text(l, (s*columnWidth)+numWidth+2, boxBottom, helpStyle)
			}
		}
	}
	srgn := sf.Region(5, 0, sf.Height(), sf.Width())
	if v.distancesResult.maxWords+2 > srgn.Height() {
		v.scrollbar = layout.NewVerticalScrollbar(v.scrollHandler(m))
		v.scrollbar.Draw(srgn, v.distancesResult.maxWords, v.offsetY)
	}
}

func (v *viewDistances) scroll(msg tea.Msg) (handled bool) {
	if v.scrollbar != nil {
		switch msg.(type) {
		case tea.MouseClickMsg:
			return v.scrollbar.Update(msg)
		}
	}
	return false
}

func (v *viewDistances) scrollHandler(m *model) layout.ScrollHandler {
	return func(evt layout.ScrollEvent) {
		switch evt {
		case layout.ScrollUp:
			v.key(m, tea.KeyPressMsg{Text: up})
		case layout.ScrollDown:
			v.key(m, tea.KeyPressMsg{Text: down})
		case layout.ScrollPageUp:
			v.key(m, tea.KeyPressMsg{Text: pageUp})
		case layout.ScrollPageDown:
			v.key(m, tea.KeyPressMsg{Text: pageDown})
		case layout.ScrollHome:
			v.key(m, tea.KeyPressMsg{Text: home})
		case layout.ScrollEnd:
			v.key(m, tea.KeyPressMsg{Text: end})
		}
	}
}

func (v *viewDistances) helpLines() ([]string, *lipgloss.Style) {
	if v.distancesResult != nil && v.distancesResult.mode == distancesNormal {
		return []string{
			enter + ": Lookup  •  " + ctrlAnalyse + ": Analyse  •  1: Islands  •  2: Doublets",
			back + ": Back",
		}, nil
	}
	return []string{
		enter + ": Lookup  •  1: Islands  •  2: Doublets",
		back + ": Back"}, nil
}

func (v *viewDistances) menu() []menuItem {
	if l := len(v.input.value()); l < 2 || l > 15 {
		return nil
	}
	menuItems := []menuItem{
		{text: "Distances", key: enter},
		{text: "Islands", key: "1"},
		{text: "Doublets", key: "2"},
		{text: "Longest ladders", key: "0"},
	}
	if v.distancesResult != nil {
		switch v.distancesResult.mode {
		case distancesNormal:
			menuItems[0] = menuItem{text: "Analyse", key: ctrlAnalyse}
			menuItems = append(menuItems, menuItem{}, menuItem{text: "Export", key: ctrlExport})
		case distancesIslands, distancesDoublets, distancesLongests, distancesAnalysis:
			menuItems = append(menuItems, menuItem{}, menuItem{text: "Export", key: ctrlExport})
		}
	}
	return menuItems
}

func (v *viewDistances) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "0":
		return v.doLongests()
	case "1":
		return v.doIslands()
	case "2":
		return v.doDoublets()
	case "+":
		return v.doOverallAnalysisSpread()
	case "-":
		return v.doOverallAnalysisGraph()
	case "=":
		return v.doOverallAnalysisAdjacents()
	case "_":
		return v.doOverallAnalysisWordCounts()
	case ctrlAnalyse:
		if v.distancesResult != nil && v.distancesResult.mode == distancesNormal {
			v.offsetY = 0
			v.offsetX = 0
			v.distancesResult.mode = distancesAnalysis
		}
	case ctrlExport:
		if v.distancesResult != nil {
			switch v.distancesResult.mode {
			case distancesNormal, distancesAnalysis:
				go exportDistances(*v.distancesResult)
			case distancesIslands:
				go exportIslands(*v.distancesResult)
			case distancesDoublets:
				go exportDoublets(*v.distancesResult)
			case distancesLongests:
				go exportLongestLadders(*v.distancesResult)
			}
		}
	case back:
		if v.distancesResult != nil && v.distancesResult.mode == distancesAnalysis {
			v.offsetY = 0
			v.offsetX = 0
			v.distancesResult.mode = distancesNormal
		} else {
			m.restoreView(distances, v.backMode, v.backView)
		}
		return nil
	case home:
		v.offsetY = 0
		v.offsetX = 0
	case end:
		if v.distancesResult != nil && v.distancesResult.mode != distancesAnalysis {
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
		if v.distancesResult != nil {
			switch v.distancesResult.mode {
			case distancesAnalysis:
				if v.offsetY < v.distancesResult.maxLadderLength {
					v.offsetY++
				}
			default:
				if v.offsetY+m.height-10 < v.distancesResult.maxWords {
					v.offsetY++
				}
			}
		}
	case pageDown:
		v.offsetY += m.height - 9
		if v.distancesResult != nil && v.distancesResult.mode != distancesAnalysis && v.offsetY+15 >= v.distancesResult.maxWords {
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
			v.offsetX -= v.pageWidth(m)
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
			v.offsetX += v.pageWidth(m)
			if v.offsetX >= v.distancesResult.maxLadderLength {
				v.offsetX = v.distancesResult.maxLadderLength - 1
			}
		}
	case enter:
		return v.doLookup()
	default:
		if v.input.key(msg) && v.distancesResult != nil {
			v.distancesResult.hideError = true
		}
	}
	return nil
}

func exportDistances(dr distancesResult) {
	if f, err := os.Create(fmt.Sprintf("word-distances-%s.csv", dr.word)); err == nil {
		defer f.Close()
		_, _ = f.WriteString("Word,Distance,To Word\n")
		for d := 1; d <= dr.maxLadderLength; d++ {
			if wds, ok := dr.distances[d]; ok {
				for _, wd := range wds {
					_, _ = f.WriteString(dr.word)
					_, _ = f.WriteString(",")
					_, _ = f.WriteString(strconv.Itoa(d))
					_, _ = f.WriteString(",")
					_, _ = f.WriteString(wd)
					_, _ = f.WriteString("\n")
				}
			}
		}
	}
}

func exportIslands(dr distancesResult) {
	if f, err := os.Create(fmt.Sprintf("islands-%d-letters.csv", dr.wordLength)); err == nil {
		defer f.Close()
		if wds, ok := dr.distances[1]; ok {
			for _, wd := range wds {
				_, _ = f.WriteString(wd)
				_, _ = f.WriteString("\n")
			}
		}
	}
}

func exportDoublets(dr distancesResult) {
	if f, err := os.Create(fmt.Sprintf("doublets-%d-letters.csv", dr.wordLength)); err == nil {
		defer f.Close()
		if wds, ok := dr.distances[1]; ok {
			for _, wd := range wds {
				_, _ = f.WriteString(wd)
				_, _ = f.WriteString("\n")
			}
		}
	}
}

func exportLongestLadders(dr distancesResult) {
	if f, err := os.Create(fmt.Sprintf("longest-ladders-%d-letters.csv", dr.wordLength)); err == nil {
		defer f.Close()
		if wds, ok := dr.distances[1]; ok {
			_, _ = f.WriteString("Word,Ladder length\n")
			for _, wd := range wds {
				_, _ = f.WriteString(wd)
				_, _ = f.WriteString(",")
				_, _ = f.WriteString(strconv.Itoa(dr.longestLadder))
				_, _ = f.WriteString("\n")
			}
		}
	}
}

func (v *viewDistances) pageWidth(m *model) int {
	if v.distancesResult != nil {
		numWidth := len(strconv.Itoa(v.distancesResult.maxWords))
		return ((m.width - numWidth - 3) / (v.distancesResult.wordLength + 4)) - 1
	}
	return 0
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
	} else if msg.Y > 4 && v.distancesResult != nil && v.distancesResult.mode == distancesAnalysis {
		ll := msg.Y - 4 + v.offsetY
		if ll > 0 && ll < v.distancesResult.maxLadderLength {
			v.distancesResult.mode = distancesNormal
			v.offsetY = 0
			v.offsetX = ll - 1
		}
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
	distancesLongests
	distancesAnalysis
	distancesOverall
)

type distancesResult struct {
	word            string
	inDictionary    bool
	hideError       bool
	wordLength      int
	maxLadderLength int
	maxWords        int
	maxWordsAt      int
	totalWords      int
	distances       map[int][]string
	longestLadder   int
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

func (v *viewDistances) doLongests() tea.Cmd {
	if l := len(v.input.value()); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			mxll := dict.MaxSteps()
			wds := dict.WordsWithSteps(mxll)
			longests := make([]string, len(wds))
			for i, wd := range wds {
				longests[i] = wd.String()
			}
			return distancesResult{
				word:            "Longest Ladder " + strconv.Itoa(mxll),
				inDictionary:    true,
				wordLength:      l,
				maxLadderLength: 1,
				maxWords:        len(longests),
				maxWordsAt:      1,
				distances: map[int][]string{
					1: longests,
				},
				longestLadder: mxll,
				mode:          distancesLongests,
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
				totalWords:      dict.Len(),
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
				totalWords:      dict.Len(),
				distances: map[int][]string{
					1: doublets,
				},
				mode: distancesDoublets,
			}
		}
	}
	return nil
}

func (v *viewDistances) doOverallAnalysisSpread() tea.Cmd {
	if l := len(v.input.value()); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			mxll := dict.MaxSteps()
			tots := make(map[int]int, mxll)
			for _, wd := range dict.Words() {
				mwl := wd.MaxSteps()
				tots[mwl] = tots[mwl] + 1
			}
			overall := make(map[int][]string, mxll)
			maxWords := 0
			for ll := 1; ll <= mxll; ll++ {
				if tots[ll] > maxWords {
					maxWords = tots[ll]
				}
				overall[ll] = make([]string, tots[ll])
			}
			return distancesResult{
				word:            "Overall Distances (spread)",
				inDictionary:    true,
				wordLength:      l,
				maxLadderLength: mxll,
				maxWords:        maxWords,
				maxWordsAt:      1,
				distances:       overall,
				mode:            distancesOverall,
			}
		}
	}
	return nil
}

func (v *viewDistances) doOverallAnalysisGraph() tea.Cmd {
	if l := len(v.input.value()); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			mxll := dict.MaxSteps()
			tots := make(map[int]int, mxll)
			for _, wd := range dict.Words() {
				mwl := wd.MaxSteps()
				if mwl == 1 || mwl == 2 {
					tots[mwl] = tots[mwl] + 1
				} else {
					for i := 3; i <= mwl; i++ {
						tots[i] = tots[i] + 1
					}
				}
			}
			overall := make(map[int][]string, mxll)
			maxWords := 0
			for ll := 1; ll <= mxll; ll++ {
				if tots[ll] > maxWords {
					maxWords = tots[ll]
				}
				overall[ll] = make([]string, tots[ll])
			}
			return distancesResult{
				word:            "Overall Distances (graph)",
				inDictionary:    true,
				wordLength:      l,
				maxLadderLength: mxll,
				maxWords:        maxWords,
				maxWordsAt:      1,
				distances:       overall,
				mode:            distancesOverall,
			}
		}
	}
	return nil
}

func (v *viewDistances) doOverallAnalysisAdjacents() tea.Cmd {
	if l := len(v.input.value()); l >= 2 {
		v.distancesResult = nil
		v.offsetX = 0
		v.offsetY = 0
		return func() tea.Msg {
			dict := words.NewDictionary(l)
			tots := make(map[int]int)
			maxAdjs := 0
			for _, wd := range dict.Words() {
				if adjs := len(wd.LinkedWords()); adjs > 0 {
					if adjs > maxAdjs {
						maxAdjs = adjs
					}
					tots[adjs] = tots[adjs] + 1
				}
			}
			maxWords := 0
			overall := make(map[int][]string, maxAdjs)
			for k, val := range tots {
				if val > maxWords {
					maxWords = val
				}
				overall[k] = make([]string, val)
			}
			return distancesResult{
				word:            "Overall Adjacent Words",
				inDictionary:    true,
				wordLength:      l,
				maxLadderLength: maxAdjs,
				maxWords:        maxWords,
				maxWordsAt:      1,
				distances:       overall,
				mode:            distancesOverall,
			}
		}
	}
	return nil
}

func (v *viewDistances) doOverallAnalysisWordCounts() tea.Cmd {
	v.distancesResult = nil
	v.offsetX = 0
	v.offsetY = 0
	return func() tea.Msg {
		maxWords := 0
		overall := make(map[int][]string)
		for wl := 2; wl <= 15; wl++ {
			dict := words.NewDictionary(wl)
			wc := dict.Len()
			overall[wl] = make([]string, wc)
			if wc > maxWords {
				maxWords = wc
			}
		}
		return distancesResult{
			word:            "Overall Word Counts",
			inDictionary:    true,
			wordLength:      0,
			maxLadderLength: 15,
			maxWords:        maxWords,
			maxWordsAt:      1,
			distances:       overall,
			mode:            distancesOverall,
		}
	}
}
