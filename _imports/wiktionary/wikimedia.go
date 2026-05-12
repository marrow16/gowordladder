package main

import "regexp"

var (
	simpleWord = regexp.MustCompile(`^[a-z]+$`)
	rxEnglish  = regexp.MustCompile(`(?m)^==English==\s*$`)
)

type mediaWikiPage struct {
	Title    string `xml:"title"`
	NS       int    `xml:"ns"`
	Revision struct {
		Text string `xml:"text"`
	} `xml:"revision"`
}

func containsEnglishSection(text string) bool {
	return rxEnglish.MatchString(text)
}
