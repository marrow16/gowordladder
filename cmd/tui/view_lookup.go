package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"encoding/json"
	"fmt"
	"github.com/marrow16/gowordladder/cmd/tui/layout"
	"github.com/marrow16/gowordladder/words"
	"net/http"
	"strconv"
	"strings"
)

type lookupView interface {
	view
	lookupWord(word string, backMode mode, backView view) tea.Cmd
}

type viewLookup struct {
	backMode         mode
	backView         view
	offsetY          int
	totalRows        int
	input            input
	lookupResult     *lookupResult
	variationsResult *variationsResult
}

func (v *viewLookup) render(sf layout.Surface, m *model) *tea.Cursor {
	const (
		prompt    = "Word:"
		promptLen = len(prompt)
		inputLen  = 15
	)
	sf.TextRight(1, 0, promptLen+1, prompt)
	sf.TextFixed(1, promptLen+2, inputLen, v.input.value(), inputStyle)
	cp := v.input.cursorPos()
	sf.Block(2, 0, m.width, '─', helpStyle)
	switch {
	case v.variationsResult != nil && len(v.variationsResult.variations) == 0:
		sf.Text(4, 1, "No variations found", errorStyle)
	case v.variationsResult != nil:
		const (
			title    = "Variations:"
			titleLen = len(title) + 2
		)
		sf.Text(4, 1, title)
		vs := make([]string, len(v.variationsResult.variations))
		for i, vwd := range v.variationsResult.variations {
			vs[i] = vwd.String()
		}
		sf.TextWrapped(4, titleLen, m.width-titleLen, strings.Join(vs, ", "))
	case v.lookupResult != nil:
		if !v.lookupResult.inDictionary {
			sf.Text(1, promptLen+inputLen+3, "Not in my dictionary", errorStyle)
		}
		switch {
		case v.lookupResult.apiError != nil:
			sf.Text(4, 1, "API error:", errorStyle)
			sf.TextWrapped(4, 12, m.width-13, v.lookupResult.apiError.Error())
		case len(v.lookupResult.apiResponse.Entries) == 0:
			sf.Text(4, 1, "No meanings found in API dictionary", errorStyle)
			if v.lookupResult.inDictionary {
				sf.Text(5, 1, "But word exists in my dictionary", highlightStyle)
			}
		default:
			rgn := sf.Region(3, 0, m.height, m.width)
			row := 0
			for _, entry := range v.lookupResult.apiResponse.Entries {
				rgn.Text(row-v.offsetY, 1, " • "+entry.PartOfSpeech, boldStyle)
				row++
				maxNumWidth := len(strconv.Itoa(len(entry.Senses))) + 2
				maxWidth := m.width - maxNumWidth - 5
				for i, sense := range entry.Senses {
					rgn.TextRight(row-v.offsetY, 3, maxNumWidth, strconv.Itoa(i+1)+".")
					row += rgn.TextWrapped(row-v.offsetY, maxNumWidth+4, maxWidth, sense.Definition)
				}
			}
			v.totalRows = row - 1
		}
	}
	return tea.NewCursor(cp+promptLen+2, 2)
}

func (v *viewLookup) helpLines() ([]string, *lipgloss.Style) {
	return []string{enter + ": Lookup  •  " + back + ": Back"}, nil
}

func (v *viewLookup) menu() []menuItem {
	return nil
}

func (v *viewLookup) key(m *model, msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case back:
		m.restoreView(lookup, v.backMode, v.backView)
		return nil
	case up:
		if v.offsetY > 0 {
			v.offsetY--
		}
	case down:
		v.offsetY++
		if v.offsetY > v.totalRows {
			v.offsetY = v.totalRows
		}
	case home:
		v.offsetY = 0
	case end:
		v.offsetY = v.totalRows - (m.height - 7) + 1
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case pageUp:
		v.offsetY -= m.height - 7
		if v.offsetY < 0 {
			v.offsetY = 0
		}
	case pageDown:
		v.offsetY += m.height - 7
		if v.offsetY > v.totalRows {
			v.offsetY = v.totalRows
		}
	case enter:
		return v.doLookup()
	}
	if v.input.key(msg) {
		v.lookupResult = nil
		v.variationsResult = nil
	}
	return nil
}

func (v *viewLookup) paste(m *model, msg tea.PasteMsg) {
	if v.input != nil {
		v.input.paste(msg)
	}
}

type lookupResult struct {
	inDictionary bool
	apiResponse  *dictionaryResponse
	apiError     error
}

type variationsResult struct {
	variations []*words.Word
}

func (v *viewLookup) update(m *model, msg tea.Msg) tea.Cmd {
	switch mt := msg.(type) {
	case lookupResult:
		v.lookupResult = &mt
	case variationsResult:
		v.variationsResult = &mt
	}
	return nil
}

func (v *viewLookup) wordLength() int {
	if v.input != nil {
		return len(v.input.value())
	}
	return 0
}

func (v *viewLookup) currentWord() string {
	if v.input != nil {
		return v.input.value()
	}
	return ""
}

func (v *viewLookup) lookupWord(word string, backMode mode, backView view) tea.Cmd {
	v.offsetY = 0
	v.lookupResult = nil
	v.variationsResult = nil
	v.backMode = backMode
	v.backView = backView
	v.input = &wordInput{maxLength: 15, current: word, allowUnderscores: true}
	return v.doLookup()
}

func (v *viewLookup) doLookup() tea.Cmd {
	s := v.input.value()
	if l := len(s); l >= 2 {
		v.lookupResult = nil
		v.variationsResult = nil
		v.offsetY = 0
		if strings.ContainsRune(s, '_') {
			return func() tea.Msg {
				dict := words.NewDictionary(l)
				return variationsResult{
					variations: dict.Variations(s),
				}
			}
		} else {
			return func() tea.Msg {
				dict := words.NewDictionary(l)
				_, found := dict.Word(s)
				r, err := v.apiLookup()
				return lookupResult{
					inDictionary: found,
					apiResponse:  r,
					apiError:     err,
				}
			}
		}
	}
	return nil
}

const (
	dictionaryUrl = "https://freedictionaryapi.com/api/v1"
)

func (v *viewLookup) apiLookup() (result *dictionaryResponse, err error) {
	var req *http.Request
	if req, err = http.NewRequest("GET", dictionaryUrl+"/entries/en/"+strings.ToLower(v.input.value()), nil); err == nil {
		var resp *http.Response
		if resp, err = http.DefaultClient.Do(req); err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				result = &dictionaryResponse{}
				if err = json.NewDecoder(resp.Body).Decode(result); err == nil {
					result.normalize()
				}
			} else {
				err = fmt.Errorf("unexpected response status: %d", resp.StatusCode)
			}
		}
	}
	return result, err
}

type dictionaryResponse struct {
	Entries []dictionaryEntry `json:"entries"`
}
type dictionaryEntry struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Senses       []struct {
		Definition string `json:"definition"`
	} `json:"senses"`
}

func (r *dictionaryResponse) normalize() {
	parts := map[string]int{}
	newEntries := make([]dictionaryEntry, 0, len(r.Entries))
	for _, entry := range r.Entries {
		if idx, found := parts[entry.PartOfSpeech]; found {
			newEntries[idx].Senses = append(newEntries[idx].Senses, entry.Senses...)
		} else {
			parts[entry.PartOfSpeech] = len(newEntries)
			newEntries = append(newEntries, entry)
		}
	}
	r.Entries = newEntries
}
