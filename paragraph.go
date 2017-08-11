package gobstract

import (
	"bufio"
	"strings"
)

const (
	title_min int = 5
	title_max int = 100
)

type Paragraph struct {
	Title *Sentence
	Line string
	Sentences *Sentences
	Lang *Language
}

type Paragraphs []*Paragraph

func SplitText(text string, lang *Language) *Paragraphs {
	var paragraphs Paragraphs

	var reader *strings.Reader = strings.NewReader(text)
	var scanner *bufio.Scanner = bufio.NewScanner(reader)

	var order int = 1
	for scanner.Scan() {
		var title *Sentence
		var sentences *Sentences

		var text string = scanner.Text()
		var line string = strings.TrimSpace(text)
		if len(line) > 0 {
			if len(line) <= title_min {
				continue
			} else if title_min < len(line) && len(line) <= title_max {
				title = NewSentence(line, order, lang)

				if scanner.Scan() {
					line = scanner.Text()
				} else {
					break
				}
			}

			var paragraph *Paragraph = &Paragraph{title, line, sentences, lang}
			paragraph.split(&order)

			paragraphs = append(paragraphs, paragraph)
		}
	}

	return &paragraphs
}

func (paragraph *Paragraph) split(order *int) {
	var rawLines []string = SplitSentences(paragraph.Line)
	var sentences Sentences
	for _, rawSentence := range rawLines {
		var content string = strings.TrimSpace(rawSentence)

		sentences = append(sentences, NewSentence(content, *order, paragraph.Lang))
		*order++
	}

	paragraph.Sentences = &sentences
}
