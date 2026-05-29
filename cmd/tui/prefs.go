package main

import (
	"encoding/json"
	"github.com/marrow16/gowordladder/words"
	"os"
	"slices"
	"time"
)

type prefs struct {
	WordLength       int            `json:"wordLength"`
	LadderLength     int            `json:"ladderLength"`
	HighScores       []highScore    `json:"highScores,omitempty"`
	MaxScores        int            `json:"maxScores"`
	UsedDictionaries []string       `json:"usedDictionaries"`
	Miscellaneous    map[string]any `json:"miscellaneous"`
}
type highScore struct {
	Score        float64 `json:"score"`
	MaxScore     float64 `json:"maxScore"`
	WordLength   int     `json:"wordLength"`
	LadderLength int     `json:"ladderLength"`
	StartWord    string  `json:"startWord"`
	EndWord      string  `json:"endWord"`
	Date         string  `json:"date"`
}

const (
	prefsFilename       = "prefs.json"
	defaultWordLength   = 3
	defaultLadderLength = 5
	defaultMaxScores    = 10
)

func newPrefs() *prefs {
	result := &prefs{
		WordLength:   defaultWordLength,
		LadderLength: defaultLadderLength,
		MaxScores:    defaultMaxScores,
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
			if result.MaxScores < defaultMaxScores {
				result.MaxScores = defaultMaxScores
			} else if result.MaxScores > 100 {
				result.MaxScores = 100
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
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		_ = enc.Encode(p)
	}
}

func (p *prefs) addScore(score, maxScore float64, wordLength, ladderLength int, startWord, endWord string) {
	dupe := false
	for _, hs := range p.HighScores {
		if dupe = hs.StartWord == startWord && hs.EndWord == endWord && hs.LadderLength == ladderLength && int(hs.Score) == int(score); dupe {
			break
		}
	}
	if !dupe {
		scores := append(p.HighScores, highScore{Score: score, MaxScore: maxScore,
			WordLength: wordLength, LadderLength: ladderLength,
			StartWord: startWord, EndWord: endWord,
			Date: time.Now().Format("Mon 02 Jan 2006 15:04")})
		slices.SortFunc(scores, func(a, b highScore) int {
			if a.Score < b.Score {
				return 1
			} else if a.Score > b.Score {
				return -1
			}
			return 0
		})
		if len(scores) > p.MaxScores {
			scores = scores[:p.MaxScores]
		}
		p.HighScores = scores
		p.save()
	}
}

func (p *prefs) clearScores() {
	p.HighScores = []highScore{}
	p.save()
}

func (p *prefs) getInt(name string) int {
	if v, ok := p.Miscellaneous[name]; ok {
		switch vt := v.(type) {
		case int:
			return vt
		case float64:
			return int(vt)
		case json.Number:
			if i, err := vt.Int64(); err == nil {
				return int(i)
			}
		}
	}
	return 0
}

func (p *prefs) setInt(name string, value int) {
	if p.Miscellaneous == nil {
		p.Miscellaneous = make(map[string]any)
	}
	p.Miscellaneous[name] = value
	p.save()
}
