package sentence

import (
	"regexp"
	"strings"
	"github.com/lucasmenendez/gobstract/language"
	"github.com/lucasmenendez/gobstract/token"
)

type Sentence struct {
	Text string
	raw_text string
	Lang *language.Language
	Tokens []token.Token
	Score float64
	Order int
}

type Sentences []*Sentence

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

func NewSentence(text string, order int, lang *language.Language) *Sentence {
	var raw_text string = strings.ToLower(text)
	var tokens []token.Token = token.GetTokens(raw_text, lang)
	var score float64

	return &Sentence{text, raw_text, lang, tokens, score, order}
}

func (sentence *Sentence) HasToken(t token.Token) bool {
	return t.IsIn(sentence.Tokens)
}

func SplitSentences (input string) []string {
	var titles_pattern *regexp.Regexp = regexp.MustCompile(`([A-Z][a-z]+)\.`)
	var titles_needle string = `$1*|*`
	var no_titles string = titles_pattern.ReplaceAllString(input, titles_needle)

	var numbers_pattern *regexp.Regexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)`)
	var numbers_needle string = `$1*|*$2`
	var no_numbers string = numbers_pattern.ReplaceAllString(no_titles, numbers_needle)

	var stops_pattern *regexp.Regexp = regexp.MustCompile(`[^..][!?.]\s`)
	var stops_needle string = `$0{stop}`
	var no_stops string = stops_pattern.ReplaceAllString(no_numbers, stops_needle)

	var restore_pattern *regexp.Regexp = regexp.MustCompile(`\*\|\*`)
	var restore_needle = `.`
	var text string = restore_pattern.ReplaceAllString(no_stops, restore_needle)

	var spliter *regexp.Regexp = regexp.MustCompile(`{stop}`)
	return spliter.Split(text, -1)
}
