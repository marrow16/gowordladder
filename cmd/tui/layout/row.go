package layout

import (
	"bytes"
	"charm.land/lipgloss/v2"
)

type rows []row

func newRows(height, width int) rows {
	result := make([]row, height)
	for r := 0; r < height; r++ {
		result[r] = make(row, width)
	}
	return result
}

var cr = []byte{'\n'}

func (rs rows) render() string {
	if len(rs) == 0 {
		return ""
	}
	spaces := bytes.Repeat([]byte{' '}, len(rs[0]))
	var buf bytes.Buffer
	for _, r := range rs {
		lastCol := 0
		for c, seg := range r {
			if seg == nil {
				continue
			}
			if pad := c - lastCol; pad > 0 {
				buf.Write(spaces[:pad])
			}
			if seg.style != nil {
				buf.WriteString(seg.style.Render(seg.text))
			} else {
				buf.WriteString(seg.text)
			}
			lastCol = c + seg.extent
		}
		buf.Write(cr)
	}
	return buf.String()
}

func (rs rows) clear(row, col, height, width int) {
	if row < 0 {
		height += row
		row = 0
	}
	if col < 0 {
		width += col
		col = 0
	}
	for y := 0; y < height && y+row < len(rs); y++ {
		rs[y+row].clear(col, width)
	}
}

type surfaceSegment struct {
	text   string
	extent int
	style  *lipgloss.Style
}

type row []*surfaceSegment

func (r row) place(col int, seg *surfaceSegment) bool {
	if seg == nil || seg.extent <= 0 || col >= len(r) {
		return false
	}
	if col < 0 {
		crop := -col
		if crop >= seg.extent {
			return false
		}
		seg.text = string([]rune(seg.text)[crop:])
		seg.extent -= crop
		col = 0
	}
	if col+seg.extent > len(r) {
		seg.extent = len(r) - col
		seg.text = string([]rune(seg.text)[:seg.extent])
	}
	if seg.extent <= 0 {
		return false
	}
	start := col
	end := col + seg.extent
	l := len(r)
	for c := 0; c < l && c < end; c++ {
		existing := r[c]
		if existing == nil {
			continue
		}
		exStart := c
		exEnd := c + existing.extent
		// no overlap
		if exEnd <= start || exStart >= end {
			continue
		}
		// remove existing segment
		r[c] = nil
		existingRunes := []rune(existing.text)
		// keep left remainder
		if exStart < start {
			leftLen := start - exStart
			r[exStart] = &surfaceSegment{
				text:   string(existingRunes[:leftLen]),
				extent: leftLen,
				style:  existing.style,
			}
		}
		// keep right remainder
		if exEnd > end {
			rightOffset := end - exStart
			r[end] = &surfaceSegment{
				text:   string(existingRunes[rightOffset:]),
				extent: exEnd - end,
				style:  existing.style,
			}
			// and we're done
			break
		}
	}
	r[col] = seg
	return true
}

func (r row) clear(col int, width int) {
	if width <= 0 || col >= len(r) {
		return
	}
	if col < 0 {
		width += col
		col = 0
		if width <= 0 {
			return
		}
	}
	if col+width > len(r) {
		width = len(r) - col
	}
	if width <= 0 {
		return
	}
	start := col
	end := col + width
	l := len(r)
	for c := 0; c < l && c < end; c++ {
		existing := r[c]
		if existing == nil {
			continue
		}
		exStart := c
		exEnd := c + existing.extent
		// no overlap
		if exEnd <= start || exStart >= end {
			continue
		}
		// remove existing segment
		r[c] = nil
		existingRunes := []rune(existing.text)
		// keep left remainder
		if exStart < start {
			leftLen := start - exStart
			r[exStart] = &surfaceSegment{
				text:   string(existingRunes[:leftLen]),
				extent: leftLen,
				style:  existing.style,
			}
		}
		// keep right remainder
		if exEnd > end {
			rightOffset := end - exStart
			r[end] = &surfaceSegment{
				text:   string(existingRunes[rightOffset:]),
				extent: exEnd - end,
				style:  existing.style,
			}
			break
		}
	}
}
