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

func (rs rows) render(used []bool) string {
	if len(rs) == 0 {
		return ""
	}
	spaces := bytes.Repeat([]byte{' '}, len(rs[0]))
	var buf bytes.Buffer
	for rn, r := range rs {
		if used[rn] {
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
		}
		buf.Write(cr)
	}
	return buf.String()
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
	st := col
	nd := col + seg.extent
	for c := 0; c < len(r) && c < nd; c++ {
		existing := r[c]
		if existing == nil {
			continue
		}
		exStart := c
		exEnd := c + existing.extent
		// no overlap
		if exEnd <= st || exStart >= nd {
			continue
		}
		// remove existing segment
		r[c] = nil
		// keep left remainder
		if exStart < st {
			leftLen := st - exStart
			r[exStart] = &surfaceSegment{
				text:   string([]rune(existing.text)[:leftLen]),
				extent: leftLen,
				style:  existing.style,
			}
		}
		// keep right remainder
		if exEnd > nd {
			rightOffset := nd - exStart
			r[nd] = &surfaceSegment{
				text:   string([]rune(existing.text)[rightOffset:]),
				extent: exEnd - nd,
				style:  existing.style,
			}
			// and we're done
			break
		}
	}
	r[col] = seg
	return true
}
