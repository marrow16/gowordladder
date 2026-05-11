package words

import (
	"bufio"
	"compress/bzip2"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const (
	enwiktionaryLlatest = "https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles.xml.bz2"
	stage1Fmt           = "./resources/enwiktionary-%d-letters.txt.tmp"
	stage2Fmt           = "./resources/enwiktionary-%d-letters.txt"
)

var simpleWord = regexp.MustCompile(`^[a-z]+$`)

type mediaWikiPage struct {
	Title    string `xml:"title"`
	NS       int    `xml:"ns"`
	Revision struct {
		Text string `xml:"text"`
	} `xml:"revision"`
}

func TestEnWiktionaryProcess_Stage1(t *testing.T) {
	t.Skip() // because this takes time and should only be explicitly run
	files := make(map[int]*os.File)
	for i := 2; i <= 15; i++ {
		f, err := os.Create(fmt.Sprintf(stage1Fmt, i))
		require.NoError(t, err)
		files[i] = f
	}
	defer func() {
		for _, file := range files {
			_ = file.Close()
		}
	}()
	resp, err := http.Get(enwiktionaryLlatest)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	bz := bzip2.NewReader(resp.Body)
	decoder := xml.NewDecoder(bz)
	count := 0
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "page" {
			continue
		}
		var page mediaWikiPage
		require.NoError(t, decoder.DecodeElement(&page, &start))
		if page.NS != 0 {
			continue
		}
		if !simpleWord.MatchString(page.Title) {
			continue
		}
		if !containsEnglishSection(page.Revision.Text) {
			continue
		}
		wd := page.Title
		if ln := len(wd); ln >= 2 && ln <= 15 {
			_, _ = files[ln].WriteString(wd)
			_, _ = files[ln].WriteString("\n")
			count++
			if count%1000 == 0 {
				// log some progress...
				t.Logf("Written %d words", count)
			}
		}
	}
}

var rxEnglish = regexp.MustCompile(`(?m)^==English==\s*$`)

func containsEnglishSection(text string) bool {
	return rxEnglish.MatchString(text)
}

func TestEnWiktionaryProcess_Stage2(t *testing.T) {
	t.Skip() // because this takes time and should only be explicitly run
	for i := 2; i <= 15; i++ {
		t.Run(fmt.Sprintf("Building %d", i), func(t *testing.T) {
			inf, err := os.Open(fmt.Sprintf(stage1Fmt, i))
			require.NoError(t, err)
			outf, err := os.Create(fmt.Sprintf(stage2Fmt, i))
			require.NoError(t, err)
			defer func() {
				_ = inf.Close()
				_ = outf.Close()
			}()
			vars := make(variations)
			scanner := bufio.NewScanner(inf)
			words := make([]*Word, 0, 100_000)
			wordsSet := make(set)
			t.Log("Reading words")
			for scanner.Scan() {
				wd := strings.ToUpper(scanner.Text())
				if !wordsSet.hasAdd(wd) {
					word := newWord(wd, 0)
					words = append(words, word)
					vars.link(word)
				}
			}
			t.Logf("Sorting %d words", len(words))
			slices.SortFunc(words, func(a, b *Word) int {
				return strings.Compare(a.actualWord, b.actualWord)
			})
			t.Logf("Writing %d words", len(words))
			maxwdl := 100
			for n, word := range words {
				wdm := NewWordDistanceMap(word, &maxwdl)
				_, _ = outf.WriteString(word.actualWord)
				_, _ = outf.WriteString("\t")
				_, _ = outf.WriteString(strconv.Itoa(wdm.MaxDistance()))
				_, _ = outf.WriteString("\n")
				if n > 0 && n%1000 == 0 {
					// log some progress...
					t.Logf("Written %d words (%s)", n, word.actualWord)
				}
			}
		})
	}
}

type set map[string]struct{}

func (s set) hasAdd(word string) bool {
	if _, ok := s[word]; !ok {
		s[word] = struct{}{}
		return false
	}
	return true
}
