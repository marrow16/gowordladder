package words

import (
	"bufio"
	"embed"
	"strconv"
	"strings"
	"sync"
)

type Dictionary struct {
	wordLength int
	words      map[string]*Word
}

func NewDictionary(wordLength int) (result *Dictionary) {
	if existing, ok := cache.dictionaries[wordLength]; !ok {
		result = &Dictionary{}
		result.words = map[string]*Word{}
		result.wordLength = wordLength
		result.load()
		cache.dictionaries[wordLength] = result
	} else {
		result = existing
	}
	return
}

//go:embed *
var resources embed.FS

func (d *Dictionary) load() {
	file, err := resources.Open("resources/dictionary-" + strconv.Itoa(d.wordLength) + "-letter-words.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var builder = &wordLinkageBuilder{variations: map[string]*[]*Word{}}
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

func (d *Dictionary) WordLength() int {
	return d.wordLength
}

func (d *Dictionary) addWord(actualWord string, builder *wordLinkageBuilder) {
	if len(actualWord) == d.wordLength {
		var word = newWord(actualWord)
		d.words[word.ActualWord()] = word
		builder.link(word)
	}
}

type wordLinkageBuilder struct {
	variations map[string]*[]*Word
}

func (b *wordLinkageBuilder) link(word *Word) {
	for _, variant := range word.variations() {
		links := b.computeIfAbsent(variant)
		for _, link := range *links {
			link.addLink(word)
			word.addLink(link)
		}
		*links = append(*links, word)
	}
}

func (b *wordLinkageBuilder) computeIfAbsent(variant string) (result *[]*Word) {
	if existing, ok := b.variations[variant]; !ok {
		result = &[]*Word{}
		b.variations[variant] = result
	} else {
		result = existing
	}
	return
}

type dictionaryCache struct {
	dictionaries map[int]*Dictionary
}

var once sync.Once
var (
	cache *dictionaryCache
)

func init() {
	once.Do(func() {
		cache = &dictionaryCache{dictionaries: map[int]*Dictionary{}}
	})
}

func LoadDictionary(wordLength int) (result *Dictionary) {
	if existing, ok := cache.dictionaries[wordLength]; !ok {
		result = NewDictionary(wordLength)
		cache.dictionaries[wordLength] = result
	} else {
		result = existing
	}
	return
}
