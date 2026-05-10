package words

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const importFilename = "CSW24.txt"

// use this test to build dictionaries from inputFilename
func TestDictionaryBuild(t *testing.T) {
	t.Skip()
	input, err := os.Open(importFilename)
	require.NoError(t, err)
	name := strings.ToLower(strings.TrimSuffix(filepath.Base(input.Name()), filepath.Ext(input.Name())))
	outFmt := "./resources/" + name + "-%d-letters.txt"
	files := make(map[int]*os.File)
	for i := 2; i <= 15; i++ {
		f, err := os.Create(fmt.Sprintf(outFmt, i))
		require.NoError(t, err)
		files[i] = f
	}
	defer func() {
		_ = input.Close()
		for _, file := range files {
			_ = file.Close()
		}
	}()
	words := make(map[int][]*Word)
	vars := make(variations)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		if l := len(scanner.Text()); l >= 2 && l <= 15 {
			word := newWord(scanner.Text(), 0)
			words[l] = append(words[l], word)
			vars.link(word)
		} else {
			assert.True(t, l >= 3 && l <= 15)
		}
	}
	maxwdl := 100
	for wl := 2; wl <= 15; wl++ {
		t.Run(fmt.Sprintf("building %d", wl), func(t *testing.T) {
			for _, w := range words[wl] {
				wdm := NewWordDistanceMap(w, &maxwdl)
				_, err = files[wl].WriteString(w.String() + "\t" + strconv.Itoa(wdm.MaxDistance()) + "\n")
				require.NoError(t, err)
			}
		})
	}
}
