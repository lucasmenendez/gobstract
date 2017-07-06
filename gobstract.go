package gobstract

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"github.com/lucasmenendez/gobstract/language"
	"github.com/lucasmenendez/gobstract/sentence"
)

var min_length int = 100
var languages []string = []string{"en", "es"}

type Gobstract struct {
	Text string
	Paragraphs []sentence.Sentences
	Sentences []string
	Lang *language.Language
}

func NewAbstract(text string, lang_label string) (*Gobstract, error) {
	var paragraphs []sentence.Sentences
	var sentences []string
	var lang *language.Language = language.GetLanguage(lang_label)

	var gobstract *Gobstract = &Gobstract{text, paragraphs, sentences, lang}
	gobstract.splitText()
	gobstract.calcScore()
	gobstract.selectHighlights()

	return gobstract, nil
}

func (gobstract *Gobstract) splitText() {
	var text string = strings.ToLower(gobstract.Text)
	
	var rgx_paragraph *regexp.Regexp = regexp.MustCompile(`\n`)
	var paragraphs []string	= rgx_paragraph.Split(text, -1)

	var rgx_sentence *regexp.Regexp = regexp.MustCompile(`\. `)
	for _, paragraph := range paragraphs {
		if len(paragraph) >= min_length {
			var raw_sentences []string = rgx_sentence.Split(paragraph, -1)

			var sentences []*sentence.Sentence
			for order, raw_content := range raw_sentences {
				var s *sentence.Sentence = sentence.NewSentence(raw_content, order, gobstract.Lang)
				sentences = append(sentences, s)
			}

			gobstract.Paragraphs = append(gobstract.Paragraphs, sentences)
		}
	}
}

func (gobstract *Gobstract) calcScore() {
	gobstract.neighboursScore()
	gobstract.keywordsScore()
}

func (gobstract *Gobstract) neighboursScore() {
	var dimension int = 0
	for _, paragraph := range gobstract.Paragraphs {
		dimension += len(paragraph)
	}
	
	for _, paragraph := range gobstract.Paragraphs {
		for i, s1 := range paragraph {
			var score float64 = 0.0
			for j, s2 := range paragraph {
				if i != j {
					var coincidences, average float64
					for _, t1 := range s1.Tokens {
						for _, t2 := range s2.Tokens {
							if t1 == t2 {
								coincidences += 1.0
							}
						}
					}
					average = float64(len(s1.Tokens) + len(s2.Tokens)) / 2.0
					score += (coincidences / average)
				}
			}

			s1.Score = score
		}
	}
}

func (gobstract *Gobstract) keywordsScore() {
	var tokens []string
	for _, paragraph := range gobstract.Paragraphs {
		for _, sentence := range paragraph {
			tokens = append(tokens, sentence.Tokens...)
		}
	}

	var sum_avg, count_avg int
	var weights map[string]int = make(map[string]int, len(tokens))
	for i, token := range tokens {
		var occurrences int
		for j, t := range tokens {
			if i != j && token == t {
				occurrences++
			}
		}

		if occurrences > 0 {
			sum_avg += occurrences
			count_avg++
		}
		weights[token] = occurrences
	}
	if count_avg > 0 {
		var keywords []string = make([]string, count_avg)
		var average int = (sum_avg/count_avg)
	
		for raw_token, weight := range weights {
			var token string = strings.TrimSpace(raw_token)
			if len(token) > 0 && weight >= average {
				keywords = append(keywords, token)
			}
		}
	
		for _, keyword := range keywords {
			fmt.Println(keyword)	
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
	var average float64 = (sum_avg/float64(count))

	var sum_var float64
	for _, sentence := range sentences {
		var sub float64 = sentence.Score - average
		sum_var += math.Pow(sub, 2)
	}
	sentences.SortOrder()

	var variance float64 = math.Sqrt(sum_var/float64(count))
	for _, sentence := range sentences {
		if sentence.Score > variance {
			gobstract.Sentences = append(gobstract.Sentences, sentence.Text)	
		}
	}
} 
