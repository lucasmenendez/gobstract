package gobstract

import (
	"github.com/lucasmenendez/gobstract/language"
	"github.com/lucasmenendez/gobstract/paragraph"
	"github.com/lucasmenendez/gobstract/scorer"
)

var min_length int = 100

type Gobstract struct {
	Text string
	Paragraphs *paragraph.Paragraphs
	Sentences []string
	Lang *language.Language
}

func NewAbstract(text string, lang_label string) (*Gobstract, error) {
	var sentences []string

	var err error
	var lang *language.Language
	if lang, err = language.GetLanguage(lang_label); err != nil {
		return nil, err
	}

	var paragraphs *paragraph.Paragraphs = paragraph.SplitText(text, lang)
	var gobstract *Gobstract = &Gobstract{text, paragraphs, sentences, lang}

	var scorer *scorer.Scorer = scorer.NewScorer(gobstract.Paragraphs)
	scorer.Calc()

	gobstract.Sentences = scorer.SelectHighlights()
	return gobstract, nil
}

