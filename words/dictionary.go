package words

import (
	"bufio"
	"fmt"
	"gowordladder/words/resources"
	"strconv"
	"strings"
	"sync"
)

type Dictionary interface {
	Word(s string) (word Word, ok bool)
	Len() int
	WordLength() int
	MaxSteps() int
	WordsWithSteps(steps int) []Word
	Words() []Word
}

type dictionary struct {
	wordLength   int
	words        map[string]Word
	wordsBySteps map[int][]Word
	maxSteps     int
}

func NewDictionary(wordLength int) Dictionary {
	var result Dictionary
	if existing, ok := cache.dictionaries[wordLength]; !ok {
		newDict := &dictionary{
			wordLength:   wordLength,
			words:        make(map[string]Word),
			wordsBySteps: make(map[int][]Word),
		}
		newDict.load()
		cache.dictionaries[wordLength] = newDict
		result = newDict
	} else {
		result = existing
	}
	return result
}

func (d *dictionary) load() {
	file, err := resources.Files.Open("dictionary-" + strconv.Itoa(d.wordLength) + "-letter-words.txt")
	if err != nil {
		panic(any(err.Error()))
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	builder := &wordLinkageBuilder{variations: map[string][]Word{}}
	for scanner.Scan() {
		d.addWord(scanner.Text(), builder)
	}
}

func (d *dictionary) Word(s string) (word Word, ok bool) {
	word, ok = d.words[strings.ToUpper(s)]
	return
}

func (d *dictionary) Len() int {
	return len(d.words)
}

func (d *dictionary) MaxSteps() int {
	return d.maxSteps
}

func (d *dictionary) WordLength() int {
	return d.wordLength
}

func (d *dictionary) WordsWithSteps(steps int) []Word {
	return d.wordsBySteps[steps]
}

func (d *dictionary) Words() (result []Word) {
	for _, w := range d.words {
		result = append(result, w)
	}
	return result
}

func (d *dictionary) addWord(line string, builder *wordLinkageBuilder) {
	if parts := strings.Split(line, "\t"); len(parts) == 2 {
		actualWord, n := parts[0], parts[1]
		maxSteps, _ := strconv.Atoi(n)
		if maxSteps > d.maxSteps {
			d.maxSteps = maxSteps
		}
		if len(actualWord) == d.wordLength {
			w := newWord(actualWord, maxSteps)
			d.words[w.ActualWord()] = w
			for i := 3; i <= maxSteps; i++ {
				d.wordsBySteps[i] = append(d.wordsBySteps[i], w)
			}
			builder.link(w)
		}
	} else {
		panic(fmt.Sprintf("invalid word input: %q", line))
	}
}

type wordLinkageBuilder struct {
	variations map[string][]Word
}

func (b *wordLinkageBuilder) link(word Word) {
	for _, variant := range word.Variations() {
		for _, link := range b.variations[variant] {
			link.addLink(word)
			word.addLink(link)
		}
		b.variations[variant] = append(b.variations[variant], word)
	}
}

type dictionaryCache struct {
	dictionaries map[int]Dictionary
}

var once sync.Once
var (
	cache *dictionaryCache
)

func init() {
	once.Do(func() {
		cache = &dictionaryCache{dictionaries: map[int]Dictionary{}}
	})
}

func LoadDictionary(wordLength int) (result Dictionary) {
	if existing, ok := cache.dictionaries[wordLength]; !ok {
		result = NewDictionary(wordLength)
		cache.dictionaries[wordLength] = result
	} else {
		result = existing
	}
	return
}
