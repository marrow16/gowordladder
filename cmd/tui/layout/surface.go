package layout

import (
	"charm.land/lipgloss/v2"
	"strings"
	"unicode/utf8"
)

type surfaceInternal interface {
	place(row int, col int, text string, extent int, styles ...lipgloss.Style)
	rows_() rows
	rowUsed(row int) bool
}

type Surface interface {
	surfaceInternal
	Height() int
	Width() int
	Render() string
	Text(row int, col int, text string, styles ...lipgloss.Style)
	TextRun(row int, col int, items Runs)
	TextRunWrapped(row, col, width int, items Runs) int
	TextFixed(row int, col int, width int, text string, styles ...lipgloss.Style)
	TextRight(row int, col int, width int, text string, styles ...lipgloss.Style)
	TextCenter(row int, col int, width int, text string, styles ...lipgloss.Style)
	LineColumns(row int, col int, width int, columns Runs, lineStyle lipgloss.Style)
	Block(row int, col int, width int, ch rune, styles ...lipgloss.Style)
	Box(row int, col int, height int, width int, styles ...lipgloss.Style)
	BoxRounded(row int, col int, height int, width int, styles ...lipgloss.Style)
	TextWrapped(row int, col int, width int, text string, styles ...lipgloss.Style) int
	Fill(row int, col int, height int, width int, styles ...lipgloss.Style)
	Draw(row int, col int, other Surface)
	Region(row, col, height, width int) Surface
}

type surface struct {
	height, width int
	rows          rows
	usedRows      []bool
}

func NewSurface(height, width int) Surface {
	return &surface{
		height:   height,
		width:    width,
		rows:     newRows(height, width),
		usedRows: make([]bool, height),
	}
}

func (s *surface) Height() int {
	return s.height
}

func (s *surface) Width() int {
	return s.width
}

func (s *surface) Render() string {
	return s.rows.render(s.usedRows)
}

func (s *surface) rows_() rows {
	return s.rows
}

func (s *surface) rowUsed(row int) bool {
	return s.usedRows[row]
}

func (s *surface) place(row, col int, text string, extent int, styles ...lipgloss.Style) {
	if row >= 0 && row < s.height && col < s.width && extent > 0 {
		style := inheritStyles(styles...)
		seg := &surfaceSegment{
			text:   text,
			extent: extent,
			style:  style,
		}
		if s.rows[row].place(col, seg) {
			s.usedRows[row] = true
		}
	}
}

func (s *surface) Text(row, col int, text string, styles ...lipgloss.Style) {
	if row >= 0 && row < s.height && col < s.width {
		s.place(row, col, text, utf8.RuneCountInString(text), styles...)
	}
}

func (s *surface) TextRun(row, col int, items Runs) {
	for _, item := range items {
		extent := utf8.RuneCountInString(item.Text)
		s.place(row, col, item.Text, extent, item.Styles...)
		col += extent
	}
}

func (s *surface) TextRunWrapped(row, col, width int, items Runs) int {
	if width <= 0 {
		return 0
	}
	startCol := col
	lineUsed := 0
	lines := 0
	newLine := func() {
		row++
		col = startCol
		lineUsed = 0
		lines++
	}
	for _, item := range items {
		if item.Text == "\n" {
			newLine()
			continue
		}
		parts := splitWords(item.Text)
		for _, part := range parts {
			if part == "\n" {
				newLine()
				continue
			}
			extent := utf8.RuneCountInString(part)
			if extent == 0 {
				continue
			}
			isSpace := strings.TrimSpace(part) == ""
			if isSpace && lineUsed == 0 {
				continue
			}
			if lineUsed > 0 && lineUsed+extent > width {
				newLine()
				if isSpace {
					continue
				}
			}
			if extent > width {
				runes := []rune(part)
				for len(runes) > 0 {
					remaining := width - lineUsed
					if remaining <= 0 {
						newLine()
						remaining = width
					}
					n := len(runes)
					if n > remaining {
						n = remaining
					}
					txt := string(runes[:n])
					s.place(row, col, txt, n, item.Styles...)
					col += n
					lineUsed += n
					runes = runes[n:]
				}
				continue
			}
			s.place(row, col, part, extent, item.Styles...)
			col += extent
			lineUsed += extent
		}
	}
	if lineUsed > 0 || lines == 0 {
		lines++
	}
	return lines
}

func splitWords(text string) []string {
	var parts []string
	var sb strings.Builder
	sb.Grow(len(text))
	var lastSpace *bool
	for _, r := range text {
		if r == '\n' {
			if sb.Len() > 0 {
				parts = append(parts, sb.String())
				sb.Reset()
			}
			parts = append(parts, string(r))
			lastSpace = nil
			continue
		}
		isSpace := r == ' ' || r == '\t'
		if lastSpace != nil && isSpace != *lastSpace {
			parts = append(parts, sb.String())
			sb.Reset()
		}
		sb.WriteRune(r)
		lastSpace = &isSpace
	}
	if sb.Len() > 0 {
		parts = append(parts, sb.String())
	}
	return parts
}

func (s *surface) TextFixed(row, col, width int, text string, styles ...lipgloss.Style) {
	if row >= 0 && row < s.height && col < s.width && width > 0 {
		l := utf8.RuneCountInString(text)
		switch {
		case l == width:
			s.place(row, col, text, width, styles...)
		case l < width:
			s.place(row, col, text+strings.Repeat(" ", width-l), width, styles...)
		default:
			s.place(row, col, string([]rune(text)[l-width+1:])+" ", width, styles...)
		}
	}
}

func (s *surface) TextRight(row, col, width int, text string, styles ...lipgloss.Style) {
	if extent := utf8.RuneCountInString(text); extent > 0 {
		s.place(row, col+width-extent, text, extent, styles...)
	}
}

func (s *surface) TextCenter(row, col, width int, text string, styles ...lipgloss.Style) {
	if extent := utf8.RuneCountInString(text); extent > 0 {
		s.place(row, col+((width-extent)/2), text, extent, styles...)
	}
}

func (s *surface) LineColumns(row, col, width int, columns Runs, lineStyle lipgloss.Style) {
	s.Block(row, col, width, ' ', lineStyle)
	if l := len(columns); l == 1 {
		s.TextCenter(row, col, width, columns[0].Text, columns[0].baseStyle(lineStyle)...)
	} else if l > 1 {
		// left most column...
		s.Text(row, col, columns[0].Text, columns[0].baseStyle(lineStyle)...)
		// right most column...
		right := columns[l-1]
		rightWd := utf8.RuneCountInString(right.Text)
		s.place(row, col+width-rightWd, right.Text, rightWd, right.baseStyle(lineStyle)...)
		// other columns...
		for i := 1; i < l-1; i++ {
			ext := utf8.RuneCountInString(columns[i].Text)
			center := col + (width*i)/(l-1)
			s.place(row, center-(ext/2), columns[i].Text, ext, columns[i].baseStyle(lineStyle)...)
		}
	}
}

func (s *surface) Block(row, col, width int, ch rune, styles ...lipgloss.Style) {
	if width > 0 {
		s.place(row, col, strings.Repeat(string(ch), width), width, styles...)
	}
}

func (s *surface) Box(row, col, height, width int, styles ...lipgloss.Style) {
	s.box(row, col, height, width, boxTL, boxTR, boxBL, boxBR, boxV, boxH, styles...)
}

func (s *surface) BoxRounded(row, col, height, width int, styles ...lipgloss.Style) {
	s.box(row, col, height, width, rndBoxTL, rndBoxTR, rndBoxBL, rndBoxBR, boxV, boxH, styles...)
}

func (s *surface) box(row, col, height, width int, tl, tr, bl, br, v, h string, styles ...lipgloss.Style) {
	if width < 2 || height < 2 {
		return
	}
	s.place(row, col, tl+strings.Repeat(h, width-2)+tr, width, styles...)
	for r := 1; r < height-1 && row+r < s.height; r++ {
		s.place(row+r, col, v, 1, styles...)
		s.place(row+r, col+width-1, v, 1, styles...)
	}
	if row+height-1 < s.height {
		s.place(row+height-1, col, bl+strings.Repeat(h, width-2)+br, width, styles...)
	}
}

func (s *surface) TextWrapped(row, col, width int, text string, styles ...lipgloss.Style) int {
	if width <= 0 || row >= s.height {
		return 0
	}
	lines := 0
	text = strings.TrimSpace(text)
	for row+lines < s.height {
		runes := []rune(text)
		if len(runes) == 0 {
			break
		}
		if len(runes) <= width {
			s.Text(row+lines, col, text, styles...)
			lines++
			break
		}
		cut := width
		line := string(runes[:width])
		if i := strings.LastIndex(line, " "); i > 0 {
			cut = utf8.RuneCountInString(line[:i])
		}
		s.Text(row+lines, col, strings.TrimSpace(string(runes[:cut])), styles...)
		text = strings.TrimLeft(string(runes[cut:]), " ")
		lines++
	}
	return lines
}

func (s *surface) Fill(row, col, height, width int, styles ...lipgloss.Style) {
	if height > 0 && width > 0 && row < s.height && col < s.width {
		for r := 0; r < height; r++ {
			s.place(row+r, col, strings.Repeat(" ", width), width, styles...)
		}
	}
}

func (s *surface) Draw(row, col int, other Surface) {
	for y, r := range other.rows_() {
		//		if !other.usedRow(y) {
		//			continue
		//		}
		for x, seg := range r {
			if seg == nil {
				continue
			}
			s.place(row+y, col+x, seg.text, seg.extent, styleValue(seg.style)...)
		}
	}
}

func (s *surface) Region(row, col, height, width int) Surface {
	if row >= 0 && col >= 0 && height > 0 && width > 0 && row < s.height && col < s.width {
		if row+height > s.height {
			height = s.height - row
		}
		if col+width > s.width {
			width = s.width - col
		}
		return &region{
			parent:    s,
			offsetRow: row,
			offsetCol: col,
			height:    height,
			width:     width,
		}
	}
	return nil
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
