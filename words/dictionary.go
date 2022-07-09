package words

import (
	"bufio"
	"embed"
	"strconv"
	"strings"
	"sync"
)

type Dictionary interface {
	Word(s string) (word Word, ok bool)
	Len() int
	WordLength() int
}

type dictionary struct {
	wordLength int
	words      map[string]Word
}

func NewDictionary(wordLength int) Dictionary {
	var result Dictionary
	if existing, ok := cache.dictionaries[wordLength]; !ok {
		newDict := &dictionary{
			wordLength: wordLength,
			words:      map[string]Word{},
		}
		newDict.load()
		cache.dictionaries[wordLength] = newDict
		result = newDict
	} else {
		result = existing
	}
	return result
}

//go:embed *
var resources embed.FS

func (d *dictionary) load() {
	file, err := resources.Open("resources/dictionary-" + strconv.Itoa(d.wordLength) + "-letter-words.txt")
	if err != nil {
		panic(any(err.Error()))
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	var builder = &wordLinkageBuilder{variations: map[string][]Word{}}
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

func (d *dictionary) WordLength() int {
	return d.wordLength
}

func (d *dictionary) addWord(actualWord string, builder *wordLinkageBuilder) {
	if len(actualWord) == d.wordLength {
		var word = newWord(actualWord)
		d.words[word.ActualWord()] = word
		builder.link(word)
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
