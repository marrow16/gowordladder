package main

import (
	tea "charm.land/bubbletea/v2"
	"encoding/json"
	"fmt"
	"gowordladder/words"
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
	input            input
	lookupResult     *lookupResult
	variationsResult *variationsResult
}

func (v *viewLookup) content(m *model) (string, *tea.Cursor) {
	const (
		prompt      = " Word: "
		footerLines = 2
	)
	var sb strings.Builder
	sb.WriteString("\n" + prompt)
	s, cxp := v.input.render()
	sb.WriteString(s)
	csr := tea.NewCursor(cxp+len(prompt), 2)
	sb.WriteString("\n" + strings.Repeat("─", m.width) + "\n")
	lines := 4
	if v.variationsResult != nil {
		if len(v.variationsResult.variations) == 0 {
			sb.WriteString(" " + errorStyle.Render("No variations found") + "\n")
			lines++
		} else {
			vs := make([]string, len(v.variationsResult.variations))
			for i, vwd := range v.variationsResult.variations {
				vs[i] = vwd.String()
			}
			const title = " Variations: "
			vlns := wrap(strings.Join(vs, ", "), m.width-len(title))
			maxLines := m.height - lines - footerLines
			for l := 0; l < maxLines && (l+v.offsetY) < len(vlns); l++ {
				sb.WriteString("\n")
				if l == 0 {
					sb.WriteString(title)
				} else {
					sb.WriteString(strings.Repeat(" ", len(title)))
				}
				sb.WriteString(vlns[l+v.offsetY])
				lines++
			}
		}
	} else if v.lookupResult != nil {
		if !v.lookupResult.inDictionary {
			sb.WriteString(" " + errorStyle.Render("Not in my dictionary") + "\n")
			lines++
		}
		if v.lookupResult.apiError != nil {
			sb.WriteString(errorStyle.Render(" API error: "+v.lookupResult.apiError.Error()) + "\n")
			lines++
		} else if len(v.lookupResult.apiResponse.Entries) == 0 {
			sb.WriteString(" " + errorStyle.Render("No meanings found in API dictionary") + "\n")
			lines++
			if v.lookupResult.inDictionary {
				sb.WriteString(" " + highlightStyle.Render("But word exists in my dictionary") + "\n")
				lines++
			}
		} else {
			maxLines := m.height - lines - footerLines
			showLines := v.lookupResult.apiResponse.buildLines(m.width)
			for l := 0; l < maxLines && (l+v.offsetY) < len(showLines); l++ {
				sb.WriteString("\n")
				sb.WriteString(showLines[l+v.offsetY])
				lines++
			}
		}
	}
	sb.WriteString(padLines(m.height - lines - footerLines))
	return sb.String(), csr
}

func (v *viewLookup) help() string {
	return "enter: Lookup  •  " + back + ": Back"
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
				err = fmt.Errorf("unexpected response status: %d", resp.Status)
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

func (r *dictionaryResponse) buildLines(width int) []string {
	result := make([]string, 0)
	for _, entry := range r.Entries {
		result = append(result, boldStyle.Render(" • "+entry.PartOfSpeech))
		maxNWd := len(strconv.Itoa(len(entry.Senses) + 1))
		numFmt := "   %" + strconv.Itoa(maxNWd) + "d. "
		pad := strings.Repeat(" ", 3+maxNWd+2)
		for i, sense := range entry.Senses {
			num := fmt.Sprintf(numFmt, i+1)
			wrapped := wrap(sense.Definition, width-len(num))
			for w, s := range wrapped {
				if w == 0 {
					result = append(result, num+s)
				} else {
					result = append(result, pad+s)
				}
			}
		}
	}
	return result
}
