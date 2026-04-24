package words

import (
	"bufio"
	"fmt"
	"gowordladder/words/resources"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Dictionary struct {
	wordLength   int
	words        map[string]*Word
	wordsBySteps map[int][]*Word
	maxSteps     int
}

func NewDictionary(wordLength int) (result *Dictionary) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	if existing, ok := cache.dictionaries[wordLength]; !ok {
		newDict := &Dictionary{
			wordLength:   wordLength,
			words:        make(map[string]*Word),
			wordsBySteps: make(map[int][]*Word),
		}
		newDict.load()
		cache.dictionaries[wordLength] = newDict
		result = newDict
	} else {
		result = existing
	}
	return result
}

func (d *Dictionary) load() {
	file, err := resources.Files.Open("dictionary-" + strconv.Itoa(d.wordLength) + "-letter-words.txt")
	if err != nil {
		panic(any(err.Error()))
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	builder := &wordLinkageBuilder{variations: map[string][]*Word{}}
	for scanner.Scan() {
		d.addWord(scanner.Text(), builder)
	}
}

func (d *Dictionary) Word(s string) (word *Word, ok bool) {
	word, ok = d.words[strings.ToUpper(s)]
	return
}

func (d *Dictionary) Len() int {
	return len(d.words)
}

func (d *Dictionary) MaxSteps() int {
	return d.maxSteps
}

func (d *Dictionary) WordLength() int {
	return d.wordLength
}

func (d *Dictionary) WordsWithSteps(steps int) []*Word {
	return slices.Clone(d.wordsBySteps[steps])
}

func (d *Dictionary) Words() (result []*Word) {
	result = make([]*Word, 0, len(d.words))
	for _, w := range d.words {
		result = append(result, w)
	}
	return result
}

func (d *Dictionary) addWord(line string, builder *wordLinkageBuilder) {
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
	variations map[string][]*Word
}

func (b *wordLinkageBuilder) link(word *Word) {
	for _, variant := range word.Variations() {
		for _, link := range b.variations[variant] {
			link.addLink(word)
			word.addLink(link)
		}
		b.variations[variant] = append(b.variations[variant], word)
	}
}

type dictionaryCache struct {
	dictionaries map[int]*Dictionary
	mutex        sync.Mutex
}

var (
	once  sync.Once
	cache *dictionaryCache
)

func init() {
	once.Do(func() {
		cache = &dictionaryCache{dictionaries: map[int]*Dictionary{}}
	})
}
