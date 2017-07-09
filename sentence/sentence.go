package sentence

import (
	"regexp"
	"strings"
	"github.com/lucasmenendez/gobstract/language"
)

type Sentence struct {
	Text string
	raw_text string
	Lang *language.Language
	Tokens []string
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
	var tokens []string
	var score float64 = 0.0
	var sentence *Sentence = &Sentence{text, raw_text, lang, tokens, score, order}

	sentence.tokenize()
	return sentence
}

func (sentence *Sentence) HasToken(needle string) bool {
	for _, token := range sentence.Tokens {
		if needle == token {
			return true
		}
	}

	return false
}

func (sentence *Sentence) tokenize() {
	var rgx_clean = regexp.MustCompile(`\[|\]|\(|\)|\{|\}|“|”|«|»|,|´|’|-|_|—`)
	var cleaned string = rgx_clean.ReplaceAllString(sentence.Text, "")

	var rgx_word = regexp.MustCompile(` `)
	var tokens []string = rgx_word.Split(cleaned, -1)

	OUTTER:
	for _, raw_token := range tokens {
		var token string = strings.TrimSpace(raw_token)
		if len(token) > 0 {
			var lower_token string = strings.ToLower(token)
			for _, stopword := range sentence.Lang.Stopwords {
				if lower_token == stopword {
					continue OUTTER
				}
			}
			sentence.Tokens = append(sentence.Tokens, lower_token)
		}
	}
}
