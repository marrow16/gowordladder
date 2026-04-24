package generator

import (
	"errors"
	"fmt"
	"gowordladder/solving"
	"gowordladder/words"
	"math"
	"math/rand"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

type Puzzle struct {
	WordLength            int
	LadderLength          int
	StartWord             *words.Word
	EndWord               *words.Word
	Solutions             []*solving.Solution
	MaxScore              float64
	RungScore             float64
	DeductionWholeWord    float64
	DeductionPatternHint  float64
	DeductionPositionHint float64
}

func GeneratePuzzle(wordLength int, ladderLength int, startWord, endWord *string) (result *Puzzle, err error) {
	if ladderLength < 3 {
		return nil, errors.New("ladderLength must be greater than 2")
	}
	dictionary := words.NewDictionary(wordLength)
	if ladderLength > dictionary.MaxSteps() {
		return nil, fmt.Errorf("ladderLength must not be greater than %d", dictionary.MaxSteps())
	}
	var start, end *words.Word
	flip := false
	useStartWord, useEndWord := startWord, endWord
	if startWord == nil && endWord != nil {
		flip = true
		useStartWord, useEndWord = endWord, startWord
	}
	if useStartWord != nil {
		w, ok := dictionary.Word(*useStartWord)
		if !ok {
			return nil, fmt.Errorf("word %q does not exist", *useStartWord)
		}
		if ladderLength > w.MaxSteps() {
			return nil, fmt.Errorf("ladderLength for word %q must not be greater than %d", *useStartWord, w.MaxSteps())
		}
		start = w
	} else {
		candidates := dictionary.WordsWithSteps(ladderLength)
		// pick a random start word...
		w := candidates[rng.Intn(len(candidates))]
		start = w
	}
	wdm := words.NewWordDistanceMap(start, &ladderLength)
	if useEndWord != nil {
		w, ok := dictionary.Word(*useEndWord)
		if !ok {
			return nil, fmt.Errorf("word %q does not exist", *useEndWord)
		}
		end = w
		if !wdm.Reachable(end, ladderLength) {
			return nil, fmt.Errorf("word %q cannot be reached from word %q", end.ActualWord(), start.ActualWord())
		}
	} else {
		// pick a random end word...
		candidates := wdm.WordsAt(ladderLength)
		sw := candidates[rng.Intn(len(candidates))]
		end, _ = dictionary.Word(sw)
	}
	if flip {
		result = &Puzzle{
			WordLength:   wordLength,
			LadderLength: ladderLength,
			StartWord:    end,
			EndWord:      start,
		}
	} else {
		result = &Puzzle{
			WordLength:   wordLength,
			LadderLength: ladderLength,
			StartWord:    start,
			EndWord:      end,
		}
	}
	s := solving.NewSolver(solving.NewPuzzle(result.StartWord, result.EndWord))
	candidates := s.Solve(ladderLength)
	result.Solutions = make([]*solving.Solution, 0, len(candidates))
	for _, candidate := range candidates {
		if len(candidate.Ladder()) == ladderLength {
			result.Solutions = append(result.Solutions, candidate)
		}
	}
	if len(result.Solutions) == 0 {
		return nil, fmt.Errorf("sorry, generated a puzzle with no solutions found")
	}
	result.MaxScore = math.Round(((float64(ladderLength) * 5.0) + (float64(wordLength) * 2.0) - (math.Log2(float64(len(result.Solutions))) * 2.5)) * 100)
	hiddenSteps := ladderLength - 2
	result.RungScore = math.Ceil(result.MaxScore / float64(hiddenSteps))
	result.DeductionWholeWord = result.RungScore
	result.DeductionPatternHint = math.Ceil(result.RungScore * 0.50)
	result.DeductionPositionHint = math.Ceil(result.RungScore / float64(wordLength))
	return result, err
}
