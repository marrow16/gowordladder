package layout

import (
	"charm.land/lipgloss/v2"
	"strings"
	"unicode/utf8"
)

type region struct {
	parent               Surface
	offsetRow, offsetCol int
	height, width        int
}

func (r *region) place(row int, col int, text string, extent int, styles ...lipgloss.Style) {
	if row >= 0 && row < r.height && col < r.width && extent > 0 {
		style := inheritStyles(styles...)
		seg := &surfaceSegment{
			text:   text,
			extent: extent,
			style:  style,
		}
		if col < 0 {
			crop := -col
			if crop >= seg.extent {
				return
			}
			seg.text = string([]rune(seg.text)[crop:])
			seg.extent -= crop
			col = 0
		}
		if col+seg.extent > r.width {
			seg.extent = r.width - col
			seg.text = string([]rune(seg.text)[:seg.extent])
		}
		if seg.extent <= 0 {
			return
		}
		r.parent.place(r.offsetRow+row, r.offsetCol+col, seg.text, seg.extent, styles...)
	}
}

func (r *region) rows_() rows {
	parentRows := r.parent.rows_()[r.offsetRow : r.offsetRow+r.height]
	result := make(rows, len(parentRows))
	for y, pr := range parentRows {
		result[y] = pr[r.offsetCol : r.offsetCol+r.width]
	}
	return result
}

func (r *region) rowUsed(row int) bool {
	return r.parent.rowUsed(row + r.offsetRow)
}

func (r *region) Height() int {
	return r.height
}

func (r *region) Width() int {
	return r.width
}

func (r *region) Render() string {
	used := make([]bool, r.height)
	for y := range r.height {
		used[y] = r.parent.rowUsed(y + r.offsetRow)
	}
	return r.rows_().render(used)
}

func (r *region) Text(row int, col int, text string, styles ...lipgloss.Style) {
	if row >= 0 && row < r.height && col < r.width {
		r.place(row, col, text, utf8.RuneCountInString(text), styles...)
	}
}

func (r *region) TextRun(row int, col int, items Runs) {
	for _, item := range items {
		extent := utf8.RuneCountInString(item.Text)
		r.place(row, col, item.Text, extent, item.Styles...)
		col += extent
	}
}

func (r *region) TextRunWrapped(row, col, width int, items Runs) int {
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
					r.place(row, col, txt, n, item.Styles...)
					col += n
					lineUsed += n
					runes = runes[n:]
				}
				continue
			}
			r.place(row, col, part, extent, item.Styles...)
			col += extent
			lineUsed += extent
		}
	}
	if lineUsed > 0 || lines == 0 {
		lines++
	}
	return lines
}

func (r *region) TextFixed(row int, col int, width int, text string, styles ...lipgloss.Style) {
	if row >= 0 && row < r.height && col < r.width && width > 0 {
		l := utf8.RuneCountInString(text)
		switch {
		case l == width:
			r.place(row, col, text, width, styles...)
		case l < width:
			r.place(row, col, text+strings.Repeat(" ", width-l), width, styles...)
		default:
			r.place(row, col, string([]rune(text)[l-width+1:])+" ", width, styles...)
		}
	}
}

func (r *region) TextRight(row int, col int, width int, text string, styles ...lipgloss.Style) {
	if extent := utf8.RuneCountInString(text); extent > 0 {
		r.place(row, col+width-extent, text, extent, styles...)
	}
}

func (r *region) TextCenter(row int, col int, width int, text string, styles ...lipgloss.Style) {
	if extent := utf8.RuneCountInString(text); extent > 0 {
		r.place(row, col+(width-extent)/2, text, extent, styles...)
	}
}

func (r *region) LineColumns(row int, col int, width int, columns Runs, lineStyle lipgloss.Style) {
	r.Block(row, col, width, ' ', lineStyle)
	if l := len(columns); l == 1 {
		r.TextCenter(row, col, width, columns[0].Text, columns[0].baseStyle(lineStyle)...)
	} else if l > 1 {
		// left most column...
		r.Text(row, col, columns[0].Text, columns[0].baseStyle(lineStyle)...)
		// right most column...
		right := columns[l-1]
		rightWd := utf8.RuneCountInString(right.Text)
		r.place(row, col+width-rightWd, right.Text, rightWd, right.baseStyle(lineStyle)...)
		// other columns...
		for i := 1; i < l-1; i++ {
			ext := utf8.RuneCountInString(columns[i].Text)
			center := col + (width*i)/(l-1)
			r.place(row, center-(ext/2), columns[i].Text, ext, columns[i].baseStyle(lineStyle)...)
		}
	}
}

func (r *region) Block(row int, col int, width int, ch rune, styles ...lipgloss.Style) {
	if width > 0 {
		r.place(row, col, strings.Repeat(string(ch), width), width, styles...)
	}
}

func (r *region) Box(row int, col int, height int, width int, styles ...lipgloss.Style) {
	r.box(row, col, height, width, boxTL, boxTR, boxBL, boxBR, boxV, boxH, styles...)
}

func (r *region) BoxRounded(row int, col int, height int, width int, styles ...lipgloss.Style) {
	r.box(row, col, height, width, rndBoxTL, rndBoxTR, rndBoxBL, rndBoxBR, boxV, boxH, styles...)
}

func (r *region) box(row, col, height, width int, tl, tr, bl, br, v, h string, styles ...lipgloss.Style) {
	if width < 2 || height < 2 {
		return
	}
	r.place(row, col, tl+strings.Repeat(h, width-2)+tr, width, styles...)
	for y := 1; y < height-1 && row+y < r.height; y++ {
		r.place(row+y, col, v, 1, styles...)
		r.place(row+y, col+width-1, v, 1, styles...)
	}
	if row+height-1 < r.height {
		r.place(row+height-1, col, bl+strings.Repeat(h, width-2)+br, width, styles...)
	}
}

func (r *region) TextWrapped(row int, col int, width int, text string, styles ...lipgloss.Style) int {
	if width <= 0 || row >= r.height {
		return 0
	}
	lines := 0
	text = strings.TrimSpace(text)
	for row+lines < r.height {
		runes := []rune(text)
		if len(runes) == 0 {
			break
		}
		if len(runes) <= width {
			r.Text(row+lines, col, text, styles...)
			lines++
			break
		}
		cut := width
		line := string(runes[:width])
		if i := strings.LastIndex(line, " "); i > 0 {
			cut = utf8.RuneCountInString(line[:i])
		}
		r.Text(row+lines, col, strings.TrimSpace(string(runes[:cut])), styles...)
		text = strings.TrimLeft(string(runes[cut:]), " ")
		lines++
	}
	return lines
}

func (r *region) Fill(row int, col int, height int, width int, styles ...lipgloss.Style) {
	if height > 0 && width > 0 && row < r.height && col < r.width {
		for y := 0; y < height; y++ {
			r.place(row+y, col, strings.Repeat(" ", width), width, styles...)
		}
	}
}

func (r *region) Draw(row int, col int, other Surface) {
	for y, or := range other.rows_() {
		//if !other.usedRow(y) {
		//	continue
		//}
		for x, seg := range or {
			if seg == nil {
				continue
			}
			r.place(row+y, col+x, seg.text, seg.extent, styleValue(seg.style)...)
		}
	}
}

func (r *region) Region(row, col, height, width int) Surface {
	if row >= 0 && col >= 0 && height > 0 && width > 0 && row < r.height && col < r.width {
		if row+height > r.height {
			height = r.height - row
		}
		if col+width > r.width {
			width = r.width - col
		}
		return &region{
			parent:    r,
			offsetRow: row,
			offsetCol: col,
			height:    height,
			width:     width,
		}
	}
	return nil
}
