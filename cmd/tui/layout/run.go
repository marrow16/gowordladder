package layout

import "charm.land/lipgloss/v2"

type RunItem struct {
	Text   string
	Styles []lipgloss.Style
}

func (i RunItem) baseStyle(s lipgloss.Style) []lipgloss.Style {
	return append(i.Styles, s)
}

type Runs []RunItem

func NewRuns(text string, styles ...lipgloss.Style) Runs {
	return Runs{
		{
			Text:   text,
			Styles: styles,
		},
	}
}
func (r Runs) Add(text string, styles ...lipgloss.Style) Runs {
	return append(r, RunItem{
		Text:   text,
		Styles: styles,
	})
}
