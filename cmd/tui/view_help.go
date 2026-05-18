package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/words"
	"strings"
)

type helpView interface {
	view
	viewShowable
}

type viewHelp struct {
	backMode mode
	backView view
	offsetY  int
	maxRow   int
	clickF2Y int
}

func (v *viewHelp) render(sf layout.Surface, m *model) *tea.Cursor {
	rgn := sf.Region(1, 1, sf.Height(), sf.Width()-2)
	row := 0
	hLine := strings.Repeat(horizontal, rgn.Width())
	for s, section := range helpText {
		if s > 0 {
			rgn.Text(row-v.offsetY, 0, hLine, helpStyle)
			row++
		}
		rgn.TextCenter(row-v.offsetY, 0, rgn.Width(), section.header, boldStyle)
		row++
		row += rgn.TextRunWrapped(row-v.offsetY, 0, rgn.Width(), section.text)
	}
	rgn.Text(row+1-v.offsetY, 0, hLine, helpStyle)
	rgn.TextCenter(row+2-v.offsetY, 0, rgn.Width(), "Current Dictionary", boldStyle)
	rgn.TextCenter(row+3-v.offsetY, 0, rgn.Width(), words.CurrentDictionary(), highlightStyle)
	rgn.TextCenter(row+4-v.offsetY, 0, rgn.Width(), fSwitch+": Switch", helpStyle)
	v.clickF2Y = row + 6 - v.offsetY
	v.maxRow = row + 4
	return nil
}

func (v *viewHelp) helpLines() ([]string, *lipgloss.Style) {
	return []string{back + "/" + backspace + ": Back"}, nil
}

func (v *viewHelp) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back, backspace:
		m.restoreView(help, v.backMode, v.backView)
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case down:
		if v.offsetY < v.maxRow-4 {
			v.offsetY++
		}
	case home:
		v.offsetY = 0
	case end:
		v.offsetY = v.maxRow - 3
	case pageUp:
		v.offsetY -= m.height - 6
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case pageDown:
		v.offsetY += m.height - 6
		if v.offsetY >= v.maxRow-4 {
			v.offsetY = v.maxRow - 3
		}
	}
	return nil
}

func (v *viewHelp) menu() []menuItem {
	return nil
}

func (v *viewHelp) click(m *model, msg tea.Mouse) tea.Cmd {
	if v.clickF2Y != 0 && (msg.Y == v.clickF2Y || msg.Y == v.clickF2Y-1) {
		m.showDictionarySwitch()
	}
	return nil
}

func (v *viewHelp) update(m *model, msg tea.Msg) tea.Cmd {
	return nil
}

func (v *viewHelp) wordLength() int {
	return 0
}

func (v *viewHelp) currentWord() string {
	return ""
}

func (v *viewHelp) show(backMode mode, backView view) {
	v.backMode = backMode
	v.backView = backView
}

var (
	helpKeyStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaaa"))
	helpHighlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#008000"))
)

type helpSection struct {
	header string
	text   layout.Runs
}

var helpText = []helpSection{
	{
		"Word Ladder Game",
		layout.NewRuns("Is a word puzzle game invented by Lewis Carroll.\n").
			Add("The objective is to transform one word into another by changing just one letter at a time, with each step forming a valid word."),
	},
	{
		"General Keys",
		layout.NewRuns("Use the following keys on any screen:\n").
			Add(fHelp+"    ", helpKeyStyle).Add(" - to show this help screen.\n").
			Add(ctrlHighScores, helpKeyStyle).Add(" - to show the high scores table.\n").
			Add(ctrlWord, helpKeyStyle).Add(" - to lookup the meaning of a word.\n").
			Add(ctrlGenerate, helpKeyStyle).Add(" - to navigate to Generate Screen.\n").
			Add(ctrlSolver, helpKeyStyle).Add(" - to navigate to Solver Screen.\n").
			Add(ctrlDistances, helpKeyStyle).Add(" - to display distances from a word.\n").
			Add(fSwitch+"    ", helpKeyStyle).Add(" - to switch dictionary.\n").
			Add(ctrlOpen, helpKeyStyle).Add(" - to open the app menu.\n").
			Add(exit+"   ", helpKeyStyle).Add(" - to exit the app.\n").
			Add("\n").
			Add("Lookup word meaning can also be displayed by mouse clicking on a displayed word."),
	},
	{
		"Play Screen",
		layout.NewRuns("Use keys ").
			Add("↑", helpKeyStyle).Add(",").
			Add("↓", helpKeyStyle).Add(",").
			Add("←", helpKeyStyle).Add(",").
			Add("→", helpKeyStyle).Add(",").
			Add(enter, helpKeyStyle).Add(",").
			Add(tab, helpKeyStyle).Add(" & ").Add(shiftTab, helpKeyStyle).
			Add(" to navigate around the word ladder. Use ").Add("space", helpKeyStyle).Add(" to clear the current rung.\n").
			Add("\n").
			Add("Hint keys:\n").
			Add(ctrlHelp, helpKeyStyle).Add(" - to show all possible solutions.\n").
			Add(ctrlFill, helpKeyStyle).Add(" - to fill in the word for you.\n").
			Add("?", helpKeyStyle).Add(" - (on an empty rung) displays which letter positions can be changed.\n").
			Add("?", helpKeyStyle).Add(" - (on a _ character) displays hints of what the letter could be.\n").
			Add("\n").
			Add("Note that using hints will aggressively deduct from your score!"),
	},
	{
		"Generate Screen",
		layout.NewRuns("Use this screen to generate a new word ladder puzzle.\n").
			Add("\n").
			Add("Enter the desired word and ladder length. Also, optionally, enter the desired start and end word.\n").
			Add("Once the puzzle has been successfully generated, press ").
			Add(ctrlPlay, helpKeyStyle).Add(" to play the puzzle. You can also press ").
			Add(enter, helpKeyStyle).Add(" to show all possible solutions for the generated puzzle."),
	},
	{
		"Solver Screen",
		layout.NewRuns("Use this screen to solve a word ladder puzzle.\n").
			Add("\n").
			Add("Enter the start and end word of the puzzle. Then enter the maximum ladder length - if this is left blank, the app automatically calculate the minimum ladder length.\n").
			Add("Once the puzzle has been successfully solved, press ").
			Add(ctrlPlay, helpKeyStyle).Add(" to play the puzzle. You can also press ").
			Add(enter, helpKeyStyle).Add(" to show all possible solutions."),
	},
	{
		"Solutions Screen",
		layout.NewRuns("Use this screen to see all the solutions for the current puzzle.\n").
			Add("This screen can be reached from the Play, Solver and Generate screens.\n").
			Add("\n").
			Add("Each solution ladder is shown as a column, each rung shows the letter that was changed.\n").
			Add("\n").
			Add("Press ").Add(ctrlExport, helpKeyStyle).Add(" to export the current solutions to a CSV file."),
	},
	{
		"Word Lookup Screen",
		layout.NewRuns("Use this screen to check a word exists within the built-in dictionary and display a word meaning.\n").
			Add("This screen can be reached by pressing ").Add(ctrlWord, helpKeyStyle).Add(" from other screens.\n").
			Add("\n").
			Add("This screen can also be used to lookup variations of a word, for example entering `C_T` will show variations of that word with the `_` transposed."),
	},
	{
		"Word Distances Screen",
		layout.NewRuns("Use this screen to perform an edge analysis on a specified word.\n").
			Add("This screen can be reached by pressing ").Add(ctrlDistances, helpKeyStyle).Add(" from other screens.\n").
			Add("\n").
			Add("Having entered a word, the display shows the maximum ladder length from that word and the words that can be reached for each ladder length.\n").
			Add("\n").
			Add("You can press ").
			Add("1", helpKeyStyle).Add(" to show island words, ").
			Add("2", helpKeyStyle).Add(" to show doublet words, or ").
			Add("0", helpKeyStyle).Add(" to show words for longest possible ladders."),
	},
	{
		"About",
		layout.NewRuns("The default built-in dictionary is based on the Official Collins Scrabble Words (2024).\n").
			Add("\n").
			Add("Word meanings are brought to you courtesy of ").Add(dictionaryUrl, helpHighlightStyle).Add("\n").
			Add("\n").
			Add("This app is written by Martin \"Marrow\" Rowlinson and published on GitHub at ").Add("https://github.com/marrow16/gowordladder", helpHighlightStyle),
	},
}
