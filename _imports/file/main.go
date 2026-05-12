package main

import (
	"bufio"
	"fmt"
	"github.com/marrow16/gowordladder/words"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "missing first arg - import filename")
		os.Exit(1)
	}
	impFilename := args[0]
	impFile, err := os.Open(impFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open import file %q", impFilename)
		os.Exit(1)
	}
	outFs := make(map[int]*os.File)
	defer func() {
		impFile.Close()
		for _, f := range outFs {
			f.Close()
		}
	}()
	name := strings.ToLower(strings.TrimSuffix(filepath.Base(impFile.Name()), filepath.Ext(impFile.Name())))
	outPath := ""
	if len(args) > 1 {
		outPath = args[1]
	}
	outFmt := filepath.Join(outPath, name+"-%d-letters.txt")
	importers := make(map[int]*words.DictionaryImporter)

	// scan import file into importers...
	scanner := bufio.NewScanner(impFile)
	for scanner.Scan() {
		line := scanner.Text()
		ln := len(line)
		imp, ok := importers[ln]
		if !ok {
			imp = words.NewDictionaryImporter(ln)
			importers[ln] = imp
		}
		if err = imp.AddWord(line); err != nil {
			fmt.Fprintf(os.Stderr, "unable to add word %q", line)
			os.Exit(1)
		}
	}

	// process each importer into dictionary files...
	for k, v := range importers {
		outf, ok := outFs[k]
		if !ok {
			outName := fmt.Sprintf(outFmt, k)
			outf, err = os.Create(outName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to create file %q", outName)
				os.Exit(1)
			}
			outFs[k] = outf
		}
		fmt.Printf("Processing %d-letter words\n", k)
		i := 0
		err = v.Process(func(word *words.Word) error {
			_, _ = outf.WriteString(word.String())
			_, _ = outf.WriteString("\t")
			_, _ = outf.WriteString(strconv.Itoa(word.MaxSteps()))
			_, _ = outf.WriteString("\n")
			if i > 0 && i%1000 == 0 {
				fmt.Printf("Written %d\n", i)
			}
			i++
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to procss - %s", err.Error())
			os.Exit(1)
		}
	}
}
