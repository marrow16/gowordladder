package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gowordladder/solving"
	"testing"
)

func TestGeneratePuzzle_5x10(t *testing.T) {
	p, err := GeneratePuzzle(5, 10, nil, nil)
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestGeneratePuzzle_cat2dog(t *testing.T) {
	start := "cat"
	end := "dog"
	p, err := GeneratePuzzle(3, 5, &start, &end)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, 3, p.WordLength)
	assert.Equal(t, 5, p.LadderLength)
	assert.Equal(t, 221, len(p.Solutions))
	assert.Equal(t, 1153.0, p.MaxScore)
	assert.Equal(t, 385.0, p.RungScore)

	p, err = GeneratePuzzle(3, 4, &start, &end)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, 3, p.WordLength)
	assert.Equal(t, 4, p.LadderLength)
	assert.Equal(t, 4, len(p.Solutions))
	assert.Equal(t, 2100.0, p.MaxScore)
	assert.Equal(t, 1050.0, p.RungScore)
}

func TestGeneratePuzzle_filpEnd(t *testing.T) {
	end := "dog"
	p, err := GeneratePuzzle(3, 5, nil, &end)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, 3, p.WordLength)
	assert.Equal(t, 5, p.LadderLength)
	assert.True(t, len(p.Solutions) > 0)
}

func TestGeneratePuzzle_code2java(t *testing.T) {
	start := "code"
	end := "java"
	p, err := GeneratePuzzle(4, 5, &start, &end)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, 4, p.WordLength)
	assert.Equal(t, 5, p.LadderLength)
	assert.Equal(t, 2, len(p.Solutions))
	assert.Equal(t, 3050.0, p.MaxScore)
	assert.Equal(t, 1017.0, p.RungScore)
}

func TestGeneratePuzzle_randomStartWord(t *testing.T) {
	for wl := 2; wl <= 15; wl++ {
		p, err := GeneratePuzzle(wl, 5, nil, nil)
		require.NoError(t, err)
		require.NotNil(t, p.StartWord)
		require.NotNil(t, p.EndWord)

		puzzle := solving.NewPuzzle(p.StartWord, p.EndWord)
		solver := solving.NewSolver(puzzle)
		solutions := solver.Solve(5)
		require.True(t, len(solutions) > 0)
	}
}

func TestGeneratePuzzle_errors(t *testing.T) {
	t.Run("Ladder length too small", func(t *testing.T) {
		_, err := GeneratePuzzle(2, 2, nil, nil)
		require.Error(t, err)
	})
	t.Run("Ladder length too big", func(t *testing.T) {
		for wl := 2; wl < 16; wl++ {
			_, err := GeneratePuzzle(wl, 500, nil, nil)
			require.Error(t, err)
		}
	})
	t.Run("start word does not exist", func(t *testing.T) {
		w := "xxx"
		_, err := GeneratePuzzle(3, 5, &w, nil)
		require.Error(t, err)
	})
	t.Run("end word does not exist", func(t *testing.T) {
		sw := "cat"
		ew := "xxx"
		_, err := GeneratePuzzle(3, 5, &sw, &ew)
		require.Error(t, err)
	})
	t.Run("end word not reachable", func(t *testing.T) {
		sw := "cat"
		ew := "iwi"
		_, err := GeneratePuzzle(3, 5, &sw, &ew)
		require.Error(t, err)
	})
	t.Run("ladder length too big for start word", func(t *testing.T) {
		w := "iwi"
		_, err := GeneratePuzzle(3, 5, &w, nil)
		require.Error(t, err)
	})
}
