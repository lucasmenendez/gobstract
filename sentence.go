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

func (sentences Sentences) Delete(index int) {
	sentences = append(sentences[:index], sentences[index+1:]...)
}

func (sentences Sentences) SortScore() {
	for i := 0; i < len(sentences); i++ {
		for j:= i+1; j < len(sentences); j++ {
			if sentences[j].Score > sentences[i].Score {
				sentences[i], sentences[j] = sentences[j], sentences[i]
			}
		}
	}
}

func (sentences Sentences) SortOrder() {
	for i := 0; i < len(sentences); i++ {
		for j:= i+1; j < len(sentences); j++ {
			if sentences[j].Order < sentences[i].Order {
				sentences[i], sentences[j] = sentences[j], sentences[i]
			}
		}
	}
}

func NewSentence(text string, order int, lang *Language) *Sentence {
	var raw_text string = strings.ToLower(text)
	var tokens []*Token = GetTokens(raw_text, lang)
	var score float64

	return &Sentence{text, raw_text, lang, tokens, score, order}
}

func (sentence *Sentence) HasToken(t *Token) bool {
	return t.IsIn(sentence.Tokens)
}

func SplitSentences (input string) []string {
	var numbersPattern *regexp.Regexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)`)
	var numbersNeedle string = `$1*|*$2`
	var no_numbers string = numbersPattern.ReplaceAllString(input, numbersNeedle)

	var stopsPattern *regexp.Regexp = regexp.MustCompile(`[^..][!?.]\s`)
	var stopsNeedle string = `$0{stop}`
	var noStops string = stopsPattern.ReplaceAllString(no_numbers, stopsNeedle)

	var restorePattern *regexp.Regexp = regexp.MustCompile(`\*\|\*`)
	var restoreNeedle = `.`
	var text string = restorePattern.ReplaceAllString(noStops, restoreNeedle)

	var spliter *regexp.Regexp = regexp.MustCompile(`{stop}`)
	return spliter.Split(text, -1)
}