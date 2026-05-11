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
	"strconv"
	"strings"
	"testing"
)

const enwiktionaryLlatest = "https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles.xml.bz2"

var simpleWord = regexp.MustCompile(`^[a-z]+$`)

type mediaWikiPage struct {
	Title    string `xml:"title"`
	NS       int    `xml:"ns"`
	Revision struct {
		Text string `xml:"text"`
	} `xml:"revision"`
}

func TestEnWiktionaryProcess_Stage1(t *testing.T) {
	t.Skip()
	files := make(map[int]*os.File)
	outFmt := "./resources/enwiktionary-%d-letters.txt.tmp"
	for i := 2; i <= 15; i++ {
		f, err := os.Create(fmt.Sprintf(outFmt, i))
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
		}
	}
}

var rxEnglish = regexp.MustCompile(`(?m)^==English==\s*$`)

func containsEnglishSection(text string) bool {
	return rxEnglish.MatchString(text)
}

func TestEnWiktionaryProcess_Stage2(t *testing.T) {
	t.Skip()
	inFmt := "./resources/enwiktionary-%d-letters.txt.tmp"
	outFmt := "./resources/enwiktionary-%d-letters.txt"
	for i := 2; i <= 15; i++ {
		t.Run(fmt.Sprintf("Building %d", i), func(t *testing.T) {
			inf, err := os.Open(fmt.Sprintf(inFmt, i))
			require.NoError(t, err)
			outf, err := os.Create(fmt.Sprintf(outFmt, i))
			require.NoError(t, err)
			defer func() {
				_ = inf.Close()
				_ = outf.Close()
			}()
			vars := make(variations)
			scanner := bufio.NewScanner(inf)
			words := map[string]*Word{}
			for scanner.Scan() {
				wd := strings.ToUpper(scanner.Text())
				word := words[wd]
				if word == nil {
					word = newWord(wd, 0)
					words[wd] = word
				}
				vars.link(word)
			}
			maxwdl := 100
			for k, v := range words {
				wdm := NewWordDistanceMap(v, &maxwdl)
				_, _ = outf.WriteString(k)
				_, _ = outf.WriteString("\t")
				_, _ = outf.WriteString(strconv.Itoa(wdm.MaxDistance()))
				_, _ = outf.WriteString("\n")
			}
		})
	}
}
