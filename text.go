// Gobstract package make extraction summaries from text provided. The algorithm
// measures sentence relations, position and similarity to pick the most
// important text sentences.
package gobstract

import (
	"errors"
	"github.com/lucasmenendez/gopostagger"
	"github.com/lucasmenendez/gotokenizer"
	"strings"
)

const minSummarized = 280
const nTag string = "NOUN"
const aTag string = "ADJ"

// Text struct evolves text sentences and associated language to use them across
// the algorithm.
type Text struct {
	sentences sentences
	lang      language
}

// NewText function initializes Text struct splitting sentences, checking text
// length and loading the according language.
func NewText(i, c string) (*Text, error) {
	var t Text
	if len(i) < minSummarized {
		return &t, errors.New("input text too short")
	}

	var e error
	if t.lang, e = loadLanguage(c); e != nil {
		return &t, e
	}

	t.buildSentences(i)
	return &t, nil
}

// buildSentences function splits text sentences and initializes sentence
// structs measuring its length and order into the full text.
func (t *Text) buildSentences(i string) {
	for o, rs := range gotokenizer.Sentences(i) {
		var s sentence = sentence{
			tokens: t.getTokens(rs),
			raw:    rs,
			length: float64(len(rs)),
			order:  o + 1,
		}

		if len(s.raw) > 0 {
			t.sentences = append(t.sentences, s)
		}
	}
}

// getTokens function split sentence into single tokens. If exists a model for
// associated text language, function tags each one with part-of-speech tags to
// extract only NOUN's and ADJ's. Else, extract all non-stopwords tokens.
func (t *Text) getTokens(s string) (r []string) {
	var tr []string
	var ts []string = gotokenizer.Words(s)
	if m, e := gopostagger.LoadModel(t.lang.model); e != nil {
		tr = append(r, ts...)
	} else {
		tagger := gopostagger.NewTagger(m)

		var pts [][]string = tagger.Tag(ts)
		for _, i := range pts {
			var k, v string = i[1], i[0]
			if strings.Contains(k, nTag) || strings.Contains(k, aTag) {
				tr = append(r, v)
			}
		}
	}

	for _, i := range tr {
		if !t.lang.isStopword(i) {
			r = append(r, i)
		}
	}

	return
}

// Summarize function initializes a Scorer and return scoring process result.
func (t *Text) Summarize() []string {
	var scorer *scorer = newScorer(t.sentences)
	return scorer.getSummary()
}
