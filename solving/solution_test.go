package solving

import (
	"github.com/stretchr/testify/assert"
	"gowordladder/words"
	"testing"
)

func TestNewSolution(t *testing.T) {
	s := newSolution()
	assert.Equal(t, 0, len(s.Ladder()))

	dictionary := words.NewDictionary(3)
	cat, _ := dictionary.Word("CAT")
	dog, _ := dictionary.Word("DOG")

	s = newSolution(cat, dog)
	assert.Equal(t, 2, len(s.Ladder()))
}

func TestSolution_ToString(t *testing.T) {
	dictionary := words.NewDictionary(3)
	cat, _ := dictionary.Word("CAT")
	dog, _ := dictionary.Word("DOG")

	s := newSolution(cat, dog)
	str := s.String()
	assert.Equal(t, "[CAT,DOG]", str)
}

func TestNewCandidateSolution(t *testing.T) {
	dictionary := words.NewDictionary(3)
	cat, _ := dictionary.Word("CAT")
	cot, _ := dictionary.Word("COT")

	m := &mockSolver{}
	assert.Equal(t, 0, m.exploredCount)

	s := newCandidateSolution(m, cat, cot)
	assert.Equal(t, 1, m.exploredCount)
	assert.Equal(t, 2, len(s.ladder))
	assert.Equal(t, 2, len(s.seenWords))
	assert.True(t, s.seen(cat))
	assert.True(t, s.seen(cot))
	assert.Equal(t, cot, s.lastWord())
	dog, _ := dictionary.Word("DOG")
	assert.False(t, s.seen(dog))
}

func TestCandidateSolutionSpawn(t *testing.T) {
	dictionary := words.NewDictionary(3)
	cat, _ := dictionary.Word("CAT")
	cot, _ := dictionary.Word("COT")

	m := &mockSolver{}
	assert.Equal(t, 0, m.exploredCount)

	s := newCandidateSolution(m, cat, cot)
	assert.Equal(t, 1, m.exploredCount)

	cog, _ := dictionary.Word("COG")
	s2 := s.spawn(cog)
	assert.Equal(t, 2, m.exploredCount)
	assert.Equal(t, cog, s2.lastWord())
	assert.Equal(t, 3, len(s2.ladder))
	assert.Equal(t, 3, len(s2.seenWords))
	assert.True(t, s2.seen(cat))
	assert.True(t, s2.seen(cot))
	assert.True(t, s2.seen(cog))
	assert.Equal(t, cog, s2.lastWord())

	// and check the original is unchanged...
	assert.Equal(t, 2, len(s.ladder))
	assert.Equal(t, 2, len(s.seenWords))
	assert.True(t, s.seen(cat))
	assert.True(t, s.seen(cot))
	assert.Equal(t, cot, s.lastWord())
	assert.False(t, s.seen(cog))
}

func TestCandidateSolutionToSolution(t *testing.T) {
	dictionary := words.NewDictionary(3)
	cat, _ := dictionary.Word("CAT")
	cot, _ := dictionary.Word("COT")

	m := &mockSolver{}
	s := newCandidateSolution(m, cat, cot)

	fs := s.asFoundSolution(false)
	assert.Equal(t, 2, len(fs.Ladder()))
	assert.Equal(t, cat, fs.Ladder()[0])
	assert.Equal(t, cot, fs.Ladder()[1])
}

func TestCandidateSolutionToSolutionReversed(t *testing.T) {
	dictionary := words.NewDictionary(3)
	cat, _ := dictionary.Word("CAT")
	cot, _ := dictionary.Word("COT")

	m := &mockSolver{}
	s := newCandidateSolution(m, cat, cot)

	fs := s.asFoundSolution(true)
	assert.Equal(t, 2, len(fs.Ladder()))
	assert.Equal(t, cat, fs.Ladder()[1])
	assert.Equal(t, cot, fs.Ladder()[0])
}

type mockSolver struct {
	exploredCount int
}

func (m *mockSolver) incrementExplored() {
	m.exploredCount++
}
