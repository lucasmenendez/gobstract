package gobstract

import (
	"bufio"
	"strings"
	"github.com/lucasmenendez/gobstract/language"
	"github.com/lucasmenendez/gobstract/sentence"
	"github.com/lucasmenendez/gobstract/scorer"
)

var min_length int = 100

type Gobstract struct {
	Text string
	Paragraphs []sentence.Sentences
	Sentences []string
	Lang *language.Language
}

func NewAbstract(text string, lang_label string) (*Gobstract) {
	var paragraphs []sentence.Sentences
	var sentences []string
	var lang *language.Language = language.GetLanguage(lang_label)

	var gobstract *Gobstract = &Gobstract{text, paragraphs, sentences, lang}
	gobstract.splitText()

	var scorer *scorer.Scorer = scorer.NewScorer(gobstract.Paragraphs)
	scorer.Calc()

	gobstract.selectHighlights()

	return gobstract
}

func (gobstract *Gobstract) splitText() {
	var reader *strings.Reader = strings.NewReader(gobstract.Text)
	var scanner *bufio.Scanner = bufio.NewScanner(reader)

	var order int
	for scanner.Scan() {
		var paragraph string = scanner.Text()

		if len(paragraph) >= min_length {
			var raw_sentences []string = strings.Split(paragraph, ". ")
			var sentences []*sentence.Sentence

			for _, raw_content := range raw_sentences {
				var content string = strings.TrimSpace(raw_content)
				var s *sentence.Sentence = sentence.NewSentence(content, order, gobstract.Lang)
				sentences = append(sentences, s)
				order++
			}

			gobstract.Paragraphs = append(gobstract.Paragraphs, sentences)
		}
	}
}

func (gobstract *Gobstract) selectHighlights() {
	var sentences sentence.Sentences

	for _, paragraph := range gobstract.Paragraphs {
		for _, sentence := range paragraph {
			if sentence.Score > 0.0 {
				sentences = append(sentences, sentence)
			}
		}

	}
	sentences.SortScore()

	var count int
	var sum_avg float64
	for _, sentence := range sentences {
		sum_avg += sentence.Score
		count++
	}

	sentences.SortOrder()

	var average float64 = (sum_avg/float64(count))
	for _, sentence := range sentences {
		if sentence.Score > average {
			gobstract.Sentences = append(gobstract.Sentences, sentence.Text)
		}
	}
}
