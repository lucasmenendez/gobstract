package sentence

import (
	"regexp"
	"github.com/lucasmenendez/gobstract/language"
)

type Sentence struct {
	Text string
	Lang *language.Language
	Tokens []string
}

func NewSentence(text string, lang *language.Language) Sentence {
	var tokens []string
	var sentence Sentence = Sentence{text, lang, tokens}
	sentence.tokenize()
	return sentence
}

func (sentence *Sentence) tokenize() {
	var rgx_clean = regexp.MustCompile(`\[|\]|\(|\)|\{|\}`)
	var cleaned string = rgx_clean.ReplaceAllString(sentence.Text, "")

	var rgx_word = regexp.MustCompile(` `)
	var tokens []string = rgx_word.Split(cleaned, -1)

	OUTTER:
	for _, token := range tokens {
		for _, stopword := range sentence.Lang.Stopwords {
			if token == stopword {
				continue OUTTER
			}
		}
		sentence.Tokens = append(sentence.Tokens, token)
	}
}
