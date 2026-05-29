package layout

import (
	"charm.land/lipgloss/v2"
)

type surfaceInternal interface {
	place(row int, col int, text string, extent int, styles ...lipgloss.Style) Placement
	rows_() rows
}

type Surface interface {
	surfaceInternal
	surfaceText
	Render() string
	Region(row, col, height, width int) Surface
	AbsoluteTop() int
	AbsoluteLeft() int
	Clear()
	ClearArea(row, col, height, width int)
}

type surface struct {
	textSurface
	rows rows
}

func NewSurface(height, width int) Surface {
	sf := &surface{
		rows: newRows(height, width),
	}
	sf.textSurface = textSurface{
		height: height,
		width:  width,
		placer: sf,
	}
	return sf
}

func (s *surface) Render() string {
	return s.rows.render()
}

func (s *surface) rows_() rows {
	return s.rows
}

func (s *surface) place(row, col int, text string, extent int, styles ...lipgloss.Style) (result Placement) {
	if row >= 0 && row < s.height && col < s.width && extent > 0 {
		style := inheritStyles(styles...)
		seg := &surfaceSegment{
			text:   text,
			extent: extent,
			style:  style,
		}
		s.rows[row].place(col, seg)
		result = Placement{
			Text:   text,
			Extent: extent,
			Row:    row,
			Col:    col,
		}
		result.Text = text
		result.Extent = extent
	}
	return result
}

func (s *surface) Region(row, col, height, width int) Surface {
	if row >= 0 && col >= 0 && height > 0 && width > 0 && row < s.height && col < s.width {
		if row+height > s.height {
			height = s.height - row
		}
		if col+width > s.width {
			width = s.width - col
		}
		return newRegion(s, row, col, height, width)
	}
	return nil
}

func (s *surface) AbsoluteTop() int {
	return 0
}

func (s *surface) AbsoluteLeft() int {
	return 0
}

func (s *surface) Clear() {
	s.rows = newRows(s.height, s.width)
}

func (s *surface) ClearArea(row, col, height, width int) {
	s.rows.clear(row, col, height, width)
}

func styleValue(style *lipgloss.Style) []lipgloss.Style {
	if style == nil {
		return nil
	}
	return []lipgloss.Style{*style}
}

func inheritStyles(styles ...lipgloss.Style) (result *lipgloss.Style) {
	if len(styles) > 0 {
		st := styles[0]
		for i := 1; i < len(styles); i++ {
			st = st.Inherit(styles[i])
		}
		result = &st
	}
	return result
}
