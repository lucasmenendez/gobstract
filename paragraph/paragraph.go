package paragraph

import (
	"fmt"
	"bufio"
	"strings"

	"github.com/lucasmenendez/gobstract/language"
	"github.com/lucasmenendez/gobstract/sentence"
)

var title_min int = 5
var title_max int = 100

type Paragraph struct {
	Title string
	Line string
	Sentences *sentence.Sentences
	Lang *language.Language
}

type Paragraphs []*Paragraph

func SplitText(text string, lang *language.Language) *Paragraphs {
	var paragraphs Paragraphs

	var reader *strings.Reader = strings.NewReader(text)
	var scanner *bufio.Scanner = bufio.NewScanner(reader)

	var order int = 0
	for scanner.Scan() {
		var title string
		var sentences *sentence.Sentences

		var line string = scanner.Text()
		if len(line) <= title_min {
			continue	
		} else if title_min < len(line) && len(line) <= title_max {
			title = line
			if scanner.Scan() {
				line = scanner.Text()
			} else {
				break
			}
		}

		var paragraph *Paragraph = &Paragraph{title, line, sentences, lang}
		paragraph.split(&order)

		paragraphs = append(paragraphs, paragraph)
		fmt.Println(order)
	}

	return &paragraphs
}

func (paragraph *Paragraph) split(order *int) {
	var sentences sentence.Sentences

	var order_value int = *order
	var raw_line []string = strings.Split(paragraph.Line, ". ")
	for _, raw_sentence := range raw_line {
		var content string = strings.TrimSpace(raw_sentence)
		sentences = append(sentences, sentence.NewSentence(content, order_value, paragraph.Lang))
		order_value++
	}

	paragraph.Sentences = &sentences
}
