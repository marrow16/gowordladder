package layout

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Scrollbar interface {
	ItemSize(size int) Scrollbar
	Draw(sf Surface, maximum int, position int)
	Update(msg tea.Msg) (handled bool)
	PageSize() int
	MaxPosition() int
	NewPosition(current int, evt ScrollEvent) int
}

type ScrollEvent int

const (
	ScrollUp ScrollEvent = iota
	ScrollDown
	ScrollPageUp
	ScrollPageDown
	ScrollHome
	ScrollEnd
)

type ScrollHandler func(evt ScrollEvent)

func NewVerticalScrollbar(handler ScrollHandler) Scrollbar {
	return &verticalScrollbar{
		handler:    handler,
		style:      defaultScrollbarStyle,
		thumbStyle: lipgloss.NewStyle().Foreground(defaultScrollbarStyle.GetBackground()).Background(defaultScrollbarStyle.GetForeground()),
	}
}

type verticalScrollbar struct {
	handler    ScrollHandler
	style      lipgloss.Style
	thumbStyle lipgloss.Style
	maximum    int
	itemHeight int
	// used to record thumb for click testing...
	rowOffset     int
	colOffset     int
	xPosition     int
	overallHeight int
	thumbTop      int
	thumbHeight   int
}

var defaultScrollbarStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#999999")).
	Background(lipgloss.Color("#eeeeee"))

const (
	scrollbarTopArrow    = "▲"
	scrollbarBottomArrow = "▼"
	scrollbarBlock       = '\u00A0'
)

func (s *verticalScrollbar) ItemSize(size int) Scrollbar {
	s.itemHeight = size
	return s
}

func (s *verticalScrollbar) PageSize() int {
	return s.overallHeight - 1
}

func (s *verticalScrollbar) MaxPosition() int {
	if s.maximum < s.overallHeight {
		return 0
	}
	return s.maximum - s.overallHeight + 1
}

func (s *verticalScrollbar) NewPosition(current int, evt ScrollEvent) int {
	n := current
	switch evt {
	case ScrollUp:
		n--
		if n < 0 {
			n = 0
		}
	case ScrollDown:
		n++
		if n > s.MaxPosition() {
			n = s.MaxPosition()
		}
	case ScrollPageUp:
		n -= s.PageSize()
		if n < 0 {
			n = 0
		}
	case ScrollPageDown:
		n += s.PageSize()
		if n > s.MaxPosition() {
			n = s.MaxPosition()
		}
	case ScrollHome:
		n = 0
	case ScrollEnd:
		n = s.MaxPosition()
	}
	return n
}

func (s *verticalScrollbar) Draw(sf Surface, maximum int, position int) {
	if s.itemHeight < 1 {
		s.itemHeight = 1
	}
	s.rowOffset = sf.AbsoluteTop()
	s.colOffset = sf.AbsoluteLeft()
	s.maximum = maximum
	s.overallHeight = sf.Height()
	column := sf.Width() - 1
	sf.Text(0, column, scrollbarTopArrow, s.style)
	sf.FillWith(1, column, sf.Height(), 1, scrollbarBlock, s.style)
	trackHeight := s.overallHeight - 2
	contentHeight := maximum * s.itemHeight
	scrollPosition := position * s.itemHeight
	if trackHeight <= 0 || contentHeight <= s.overallHeight {
		sf.Text(s.overallHeight-1, column, scrollbarBottomArrow, s.style)
		s.xPosition = s.colOffset + column
		return
	}
	s.thumbHeight = (trackHeight * trackHeight) / contentHeight
	if s.thumbHeight < 1 {
		s.thumbHeight = 1
	}
	if s.thumbHeight > trackHeight {
		s.thumbHeight = trackHeight
	}
	maxPos := contentHeight - s.overallHeight + 1
	if scrollPosition < 0 {
		scrollPosition = 0
	}
	if scrollPosition > maxPos {
		scrollPosition = maxPos
	}
	s.thumbTop = 1 + (scrollPosition*(trackHeight-s.thumbHeight))/maxPos
	if s.thumbTop == 1 && scrollPosition > 0 {
		s.thumbTop++
	}
	if s.thumbTop+s.thumbHeight-1 == trackHeight && scrollPosition < maxPos {
		s.thumbTop--
	}
	sf.FillWith(s.thumbTop, column, s.thumbHeight, 1, scrollbarBlock, s.thumbStyle)
	sf.Text(s.overallHeight-1, column, scrollbarBottomArrow, s.style)
	s.xPosition = s.colOffset + column
}

func (s *verticalScrollbar) Update(msg tea.Msg) (handled bool) {
	if s.handler == nil {
		return false
	}
	switch mt := msg.(type) {
	case tea.MouseClickMsg:
		x, y := mt.Mouse().X, mt.Mouse().Y
		if x == s.xPosition && y >= s.rowOffset && y < s.rowOffset+s.overallHeight {
			handled = true
			switch {
			case y == s.rowOffset:
				if mt.Mouse().Button == tea.MouseLeft {
					s.handler(ScrollUp)
				} else {
					s.handler(ScrollHome)
				}
			case y == s.rowOffset+s.overallHeight-1:
				if mt.Mouse().Button == tea.MouseLeft {
					s.handler(ScrollDown)
				} else {
					s.handler(ScrollEnd)
				}
			case y < s.rowOffset+s.thumbTop:
				s.handler(ScrollPageUp)
			case y > s.rowOffset+s.thumbTop+s.thumbHeight-1:
				s.handler(ScrollPageDown)
			default:
				handled = false
			}
		}
	case tea.MouseWheelMsg:
		handled = true
		switch mt.Mouse().Button {
		case tea.MouseWheelDown:
			s.handler(ScrollDown)
		case tea.MouseWheelUp:
			s.handler(ScrollUp)
		default:
			handled = false
		}
	case tea.KeyPressMsg:
		handled = true
		switch mt.String() {
		case "up":
			s.handler(ScrollUp)
		case "down":
			s.handler(ScrollDown)
		case "pgup":
			s.handler(ScrollPageUp)
		case "pgdown":
			s.handler(ScrollPageDown)
		case "home":
			s.handler(ScrollHome)
		case "end":
			s.handler(ScrollEnd)
		default:
			handled = false
		}
	}
	return handled
}
