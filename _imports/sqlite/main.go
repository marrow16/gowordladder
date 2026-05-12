package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/marrow16/gowordladder/words"
	_ "modernc.org/sqlite"
)

// Imports words from SQLite database (i.e. as created by NASPA Zyzzyva application)
// The first arg is the filepath of the database.
//
// Writes them as text files that can be used by dictionary
// (i.e. copy the text files to words/resources)
func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: expected sqlite db filepath")
		os.Exit(1)
	}
	db, err := sql.Open("sqlite", args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer db.Close()
	name := strings.ToLower(strings.TrimSuffix(filepath.Base(args[0]), filepath.Ext(args[0])))
	for i := 2; i <= 15; i++ {
		if err := importWords(name, db, i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func importWords(name string, db *sql.DB, wordLen int) error {
	const query = `SELECT word
		FROM words
		WHERE length = ?
		ORDER BY word`
	fmt.Printf("Importing: %s (%d letters)\n", name, wordLen)
	if rows, err := db.Query(query, wordLen); err == nil {
		defer rows.Close()
		imp := words.NewDictionaryImporter(wordLen)
		for rows.Next() {
			var wd string
			if err = rows.Scan(&wd); err != nil {
				return err
			}
			if err = imp.AddWord(wd); err != nil {
				return err
			}
		}
		if f, err := os.Create(fmt.Sprintf("%s-%d-letters.txt", name, wordLen)); err != nil {
			return err
		} else {
			defer f.Close()
			count := 0
			err = imp.Process(func(word *words.Word) error {
				_, _ = f.WriteString(word.String())
				_, _ = f.WriteString("\t")
				_, _ = f.WriteString(strconv.Itoa(word.MaxSteps()))
				_, _ = f.WriteString("\n")
				if count > 0 && count%1000 == 0 {
					// log some progress...
					fmt.Printf("Written %d words\n", count)
				}
				count++
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	/*
		vars := make(variations)
		words := make([]*Word, 0)
		fmt.Printf("Importing: %s (%d letters)\n", name, wordLen)
		if rows, err = db.Query(query, wordLen); err == nil {
			defer rows.Close()
			for rows.Next() {
				var wd string
				if err = rows.Scan(&wd); err != nil {
					break
				}
				word := &Word{actualWord: strings.ToUpper(wd)}
				vars.link(word)
				words = append(words, word)
			}
			if err == nil {
				var f *os.File
				if f, err = os.Create(fmt.Sprintf("%s-%d-letters.txt", name, wordLen)); err == nil {
					maxllen := 100
					for i, word := range words {
						wdm := NewWordDistanceMap(word, &maxllen)
						longest := wdm.MaxDistance()
						_, _ = f.WriteString(word.actualWord)
						_, _ = f.WriteString("\t")
						_, _ = f.WriteString(strconv.Itoa(longest))
						_, _ = f.WriteString("\n")
						if i > 0 && i%1000 == 0 {
							fmt.Printf("Written %d words (%s)\n", i, word.actualWord)
						}
					}
				}
			}
		}
	*/
	return nil
}
