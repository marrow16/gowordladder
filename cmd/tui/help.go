package main

import (
	"charm.land/lipgloss/v2"
	"fmt"
	"strings"
)

var (
	helpHeaderStyle    = lipgloss.NewStyle().Bold(true).AlignHorizontal(lipgloss.Center)
	helpKeyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa"))
	helpHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
)

type helpHeader string
type text string
type key string
type highlight string

func (h helpHeader) render(width int) string {
	return helpHeaderStyle.Width(width).Render(string(h))
}

func (h helpHeader) String() string {
	return string(h)
}

func (h text) render() string {
	return string(h)
}

func (h text) String() string {
	return string(h)
}

func (h key) render() string {
	return helpKeyStyle.Render(string(h))
}

func (h key) String() string {
	return string(h)
}

func (h highlight) render() string {
	return helpHighlightStyle.Render(string(h))
}

func (h highlight) String() string {
	return string(h)
}

type helpPart interface {
	render() string
	fmt.Stringer
}

type helpItem struct {
	header helpHeader
	lines  []helpLine
}
type helpData []helpItem

func (h helpData) render(width int) []string {
	result := make([]string, 0, len(h)*10)
	for i := 0; i < len(h); i++ {
		if i > 0 {
			result = append(result, " "+helpStyle.Render(strings.Repeat(topBottom, width-2)))
		}
		item := h[i]
		result = append(result, item.header.render(width))
		for _, line := range item.lines {
			result = append(result, line.render(width)...)
		}
	}
	return result
}

type helpLine []helpPart

func (l helpLine) render(width int) []string {
	switch {
	case len(l) == 0:
		return []string{""}
	case len(l) == 1 && len(l[0].String())+1 <= width:
		return []string{" " + l[0].render()}
	}
	result := make([]string, 0, len(l))
	var sb strings.Builder
	usedWidth := 1
	wrap := func() {
		if sb.Len() > 0 {
			result = append(result, " "+sb.String())
			sb.Reset()
		}
		usedWidth = 1
	}
	for _, part := range l {
		if wd := len(part.String()); wd+usedWidth <= width {
			sb.WriteString(part.render())
			usedWidth += wd
		} else {
			switch part.(type) {
			case key, highlight:
				wrap()
				sb.WriteString(part.render())
				usedWidth += len(part.String())
			default:
				for _, word := range strings.Fields(part.String()) {
					if len(word)+1+usedWidth > width {
						wrap()
					}
					if usedWidth != 1 {
						sb.WriteString(" ")
					}
					sb.WriteString(word)
					usedWidth += 1 + len(word)
				}
			}
		}
	}
	wrap()
	return result
}

var helpText = helpData{
	{
		"Word Ladder Game",
		[]helpLine{
			{text("Is a word puzzle game invented by Lewis Carroll.")},
			{text("The objective is to transform one word into another by changing just one letter at a time, with each step forming a valid word.")},
		},
	},
	{
		"General Keys",
		[]helpLine{
			{text("Use the following keys on any screen:")},
			{key(fHelp + "    "), text(" - to show this help screen.")},
			{key(ctrlHighScores), text(" - to show the high scores table.")},
			{key(ctrlWord), text(" - to lookup the meaning of a word.")},
			{key(ctrlGenerate), text(" - to navigate to Generate Screen.")},
			{key(ctrlSolver), text(" - to navigate to Solver Screen.")},
			{key(ctrlDistances), text(" - to display distances from a word.")},
			{key(exit + "   "), text(" - to exit the app.")},
		},
	},
	{
		"Play Screen",
		[]helpLine{
			{text("Use keys "), key("↑"), text(","), key("↓"), text(","), key("←"), text(","), key("→"), text(","), key(enter), text(","), key(tab), text(" & "), key(shiftTab), text(" to navigate around the word ladder.  Use"), key(" space"), text(" to clear the current rung.")},
			{},
			{text("Hint keys:")},
			{key(ctrlHelp), text(" - to show all possible solutions.")},
			{key(ctrlFill), text(" - to fill in the word for you.")},
			{key("?"), text(" - (on an empty rung)")},
			{key("    "), text("displays which letter positions can be changed.")},
			{key("?"), text(" - (on a _ character)")},
			{key("    "), text("displays hints of what the letter could be.")},
			{text("Note that using hints will aggressively deduct from your score!")},
		},
	},
	{
		"Generate Screen",
		[]helpLine{
			{text("Use this screen to generate a new word ladder puzzle.")},
			{},
			{text("Enter the desired word and ladder length. Also, optionally, enter the desired start and end word.")},
			{text("Once the puzzle has been successfully generated, press"), key(" " + ctrlPlay), text(" to play the puzzle. You can also press"), key(" enter"), text(" to show all possible solutions for the generated puzzle.")},
		},
	},
	{
		"Solver Screen",
		[]helpLine{
			{text("Use this screen to solve a word ladder puzzle.")},
			{},
			{text("Enter the start and end word of the puzzle. Then enter the maximum ladder length - if this is left blank, the app automatically calculate the minimum ladder length.")},
			{text("Once the puzzle has been successfully solved, press"), key(" " + ctrlPlay), text(" to play the puzzle. You can also press"), key(" enter"), text(" to show all possible solutions.")},
		},
	},
	{
		"Solutions Screen",
		[]helpLine{
			{text("Use this screen to see all the solutions for the current puzzle.")},
			{text("This screen can be reached from the Play, Solver and Generate screens.")},
			{},
			{text("Each solution ladder is shown as a column, each rung shows the letter that was changed.")},
		},
	},
	{
		"Word Lookup Screen",
		[]helpLine{
			{text("Use this screen to check a word exists within the built-in dictionary and display a word meaning.")},
			{text("This screen can be reached by pressing"), key(" " + ctrlWord), text(" from other screens.")},
			{},
			{text("This screen can also be used to lookup variations of a word, for example entering"), key(" `C_T`"), text(" will show variations of that word with the `_` transposed.")},
		},
	},
	{
		"Word Distances Screen",
		[]helpLine{
			{text("Use this screen to perform an edge analysis on a specified word.")},
			{text("This screen can be reached by pressing"), key(" " + ctrlDistances), text(" from other screens.")},
			{},
			{text("Having entered a word, the display shows the maximum ladder length from that word and the words that can be reached for each ladder length.")},
		},
	},
	{
		"About",
		[]helpLine{
			{text("The built-in dictionary is based on the Official Collins Scrabble Words (2024).")},
			{},
			{text("Word meanings are brought to you courtesy of "), highlight(dictionaryUrl)},
			{},
			{text("This app is written by Martin \"Marrow\" Rowlinson and published on GitHub at "), highlight("https://github.com/marrow16/gowordladder")},
		},
	},
}
