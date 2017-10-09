package gobstract

import (
	"bufio"
	"strings"
)

const (
	titleMin int = 5
	titleMax int = 100
)

type Paragraph struct {
	Title *Sentence
	Line string
	Sentences *Sentences
	Lang *Language
}

type Paragraphs []*Paragraph

func SplitText(text string, lang *Language) (paragraphs *Paragraphs) {
	var ps Paragraphs
	var r *strings.Reader = strings.NewReader(text)
	var s *bufio.Scanner = bufio.NewScanner(r)

	var o int = 1
	for s.Scan() {
		var title *Sentence
		var sentences *Sentences

		var text string = s.Text()
		var line string = strings.TrimSpace(text)
		if len(line) > 0 {
			if len(line) <= titleMin {
				continue
			} else if titleMin < len(line) && len(line) <= titleMax {
				title = NewSentence(line, o, lang)

				if s.Scan() {
					line = s.Text()
				} else {
					break
				}
			}

			var p *Paragraph = &Paragraph{title, line, sentences, lang}
			p.split(&o)

			ps = append(ps, p)
		}
	}

	return &ps
}

func (p *Paragraph) split(o *int) {
	var l []string = SplitSentences(p.Line)
	var s Sentences
	for _, r := range l {
		var c string = strings.TrimSpace(r)

		s = append(s, NewSentence(c, *o, p.Lang))
		*o++
	}

	p.Sentences = &s
}
