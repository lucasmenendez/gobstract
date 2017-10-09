package gobstract

import (
	"regexp"
	"strings"
)

type Sentence struct {
	Text string
	raw_text string
	Lang *Language
	Tokens []*Token
	Score float64
	Order int
}

type Sentences []*Sentence

func (s Sentences) Delete(index int) {
	s = append(s[:index], s[index+1:]...)
}

func (s Sentences) SortScore() {
	for i := 0; i < len(s); i++ {
		for j := i+1; j < len(s); j++ {
			if s[j].Score > s[i].Score {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

func (s Sentences) SortOrder() {
	for i := 0; i < len(s); i++ {
		for j := i+1; j < len(s); j++ {
			if s[j].Order < s[i].Order {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

func NewSentence(txt string, o int, l *Language) (s *Sentence) {
	var r string = strings.ToLower(txt)
	var t []*Token = GetTokens(r, l)
	var sc float64

	return &Sentence{txt, r, l, t, sc, o}
}

func (s *Sentence) HasToken(t *Token) bool {
	return t.IsIn(s.Tokens)
}

func SplitSentences (i string) (sentences []string) {
	var numbersPattern *regexp.Regexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)`)
	var numbersNeedle string = `$1*|*$2`
	var no_numbers string = numbersPattern.ReplaceAllString(i, numbersNeedle)

	var stopsPattern *regexp.Regexp = regexp.MustCompile(`[^..][!?.]\s`)
	var stopsNeedle string = `$0{stop}`
	var noStops string = stopsPattern.ReplaceAllString(no_numbers, stopsNeedle)

	var restorePattern *regexp.Regexp = regexp.MustCompile(`\*\|\*`)
	var restoreNeedle = `.`
	var text string = restorePattern.ReplaceAllString(noStops, restoreNeedle)

	var spliter *regexp.Regexp = regexp.MustCompile(`{stop}`)
	return spliter.Split(text, -1)
}