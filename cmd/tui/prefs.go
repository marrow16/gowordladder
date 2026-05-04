package main

import (
	"encoding/json"
	"gowordladder/words"
	"os"
)

type prefs struct {
	WordLength   int `json:"wordLength"`
	LadderLength int `json:"ladderLength"`
}

const (
	prefsFilename       = "prefs.json"
	defaultWordLength   = 3
	defaultLadderLength = 5
)

func newPrefs() *prefs {
	result := &prefs{
		WordLength:   defaultWordLength,
		LadderLength: defaultLadderLength,
	}
	if f, err := os.Open(prefsFilename); err == nil {
		defer f.Close()
		if err = json.NewDecoder(f).Decode(result); err == nil {
			if result.WordLength < 2 || result.WordLength > 15 {
				result.WordLength = defaultWordLength
			}
			if result.LadderLength < 2 {
				result.LadderLength = defaultLadderLength
			}
			dict := words.NewDictionary(result.WordLength)
			if result.LadderLength > dict.MaxSteps() {
				result.LadderLength = dict.MaxSteps()
			}
		}
	}
	return result
}

func (p *prefs) save() {
	if f, err := os.Create(prefsFilename); err == nil {
		defer f.Close()
		_ = json.NewEncoder(f).Encode(p)
	}
}
