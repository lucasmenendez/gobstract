package gobstract

import (
	"fmt"
	"regexp"
	"github.com/lucasmenendez/gobstract/language"
	"github.com/lucasmenendez/gobstract/sentence"
)

var min_length int = 100
var languages []string = []string{"en", "es"}

type Gobstract struct {
	Text string
	Paragraphs [][]sentence.Sentence
	Highlights []sentence.Sentence
	Graph [][]float32
	Abstract string
	Lang *language.Language
}

func NewAbstract(text string, lang_label string) (*Gobstract, error) {
	var paragraphs [][]sentence.Sentence
	var highlights []sentence.Sentence
	var graph [][]float32
	var abstract string
	var lang *language.Language = language.GetLanguage(lang_label)

	var gobstract *Gobstract = &Gobstract{text, paragraphs, highlights, graph, abstract, lang}
	gobstract.splitText()
	gobstract.calcGraph()

	fmt.Println(gobstract.Graph)
	//gobstract.summarize()
	//gobstract.join()

	return gobstract, nil
}

func (gobstract *Gobstract) splitText() {
	var rgx_paragraph *regexp.Regexp = regexp.MustCompile(`\n`)
	var paragraphs []string	= rgx_paragraph.Split(gobstract.Text, -1)

	var rgx_sentence *regexp.Regexp = regexp.MustCompile(`\. `)
	for _, paragraph := range paragraphs {
		if len(paragraph) >= min_length {
			var raw_sentences []string = rgx_sentence.Split(paragraph, -1)

			var sentences []sentence.Sentence
			for _, raw_content := range raw_sentences {
				var s sentence.Sentence = sentence.NewSentence(raw_content, gobstract.Lang)
				sentences = append(sentences, s)
			}

			gobstract.Paragraphs = append(gobstract.Paragraphs, sentences)
		}
	}
}

func (gobstract *Gobstract) calcGraph() {
	var dimension int = len(gobstract.Paragraphs)
	var graph [][]float32 = make([][]float32, dimension)

	for _, paragraph := range gobstract.Paragraphs {
		for i, s1 := range paragraph {
			var row []float32 = make([]float32, dimension)

			fmt.Println(s1.Text)
			for j, s2 := range paragraph {
				fmt.Println(s2.Text)
				if i == j {
					row[j] = 0.0
				} else {
					var coincidences, average float32

					for _, t1 := range s1.Tokens {
						for _, t2 := range s2.Tokens {
							if t1 == t2 {
								coincidences += 1.0
							}
						}
					}

					average = float32(len(s1.Tokens) + len(s2.Tokens)) / 2.0
					row[j] = (coincidences / average)
				}
			}

			graph[i] = row
		}
	}

	gobstract.Graph = graph
}
