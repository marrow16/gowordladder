package main

import (
	"compress/bzip2"
	"encoding/xml"
	"fmt"
	"github.com/marrow16/gowordladder/words"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	enWiktionaryLlatest = "https://dumps.wikimedia.org/enwiktionary/latest/enwiktionary-latest-pages-articles.xml.bz2"
	stage1Fmt           = "enwiktionary-%d-letters.txt.tmp"
	stage2Fmt           = "enwiktionary-%d-letters.txt"
)

func main() {
	err := stage1()
	if err != nil {
		log.Fatal(err)
	}

	err = stage2()
	if err != nil {
		log.Fatal(err)
	}
}

func stage2() error {
	fmt.Println("Writing output dictionaries...")
	inFs := map[int]*os.File{}
	outFs := map[int]*os.File{}
	defer func() {
		delNames := make([]string, 0, len(inFs))
		for _, inF := range inFs {
			delNames = append(delNames, inF.Name())
			inF.Close()
		}
		for _, name := range delNames {
			os.Remove(name)
		}
		for _, outF := range outFs {
			outF.Close()
		}
	}()
	total := 0
	for wl := 2; wl <= 15; wl++ {
		if inF, err := os.Open(fmt.Sprintf(stage1Fmt, wl)); err == nil {
			inFs[wl] = inF
			imp := words.NewDictionaryImporter(wl)
			if err = imp.Import(inF); err != nil {
				return err
			}
			outF, err := os.Create(fmt.Sprintf(stage2Fmt, wl))
			if err != nil {
				return err
			}
			outFs[wl] = outF
			fmt.Printf("Stage 2 - writing %d-letter dictionary (%d words)\n", wl, len(imp.Words))
			count := 0
			if err = imp.Process(func(w *words.Word) error {
				_, _ = outF.WriteString(w.String())
				_, _ = outF.WriteString("\t")
				_, _ = outF.WriteString(strconv.Itoa(w.MaxSteps()))
				_, _ = outF.WriteString("\n")
				if count > 0 && count%1000 == 0 {
					// log some progress...
					fmt.Printf("Written %d words\n", count)
				}
				count++
				return nil
			}); err != nil {
				return err
			}
			total += count
		}
	}
	fmt.Printf("Finished downloading %d words\n", total)
	return nil
}

func stage1() error {
	fmt.Printf("Reading wiktionary from %q\n", enWiktionaryLlatest)
	files := make(map[int]*os.File)
	var resp *http.Response
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
		for _, file := range files {
			file.Close()
		}
	}()
	resp, err := http.Get(enWiktionaryLlatest)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	bz := bzip2.NewReader(resp.Body)
	decoder := xml.NewDecoder(bz)
	count := 0
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		start, ok := tok.(xml.StartElement)
		if !ok || start.Name.Local != "page" {
			continue
		}
		var page mediaWikiPage
		err = decoder.DecodeElement(&page, &start)
		if err != nil {
			return err
		}
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
			f, ok := files[ln]
			if !ok {
				f, err = os.Create(fmt.Sprintf(stage1Fmt, ln))
				if err != nil {
					return err
				}
				files[ln] = f
			}
			_, _ = f.WriteString(wd)
			_, _ = f.WriteString("\n")
			if count > 0 && count%1000 == 0 {
				// log some progress...
				fmt.Printf("Downloaded %d words\n", count)
			}
			count++
		}
	}
	return nil
}
