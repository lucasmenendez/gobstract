package scorer

import (
	"strings"
	"github.com/lucasmenendez/gobstract/paragraph"
	"github.com/lucasmenendez/gobstract/sentence"
)

type Scorer struct {
	paragraphs *paragraph.Paragraphs
}

func NewScorer(paragraphs *paragraph.Paragraphs) *Scorer {
	return &Scorer{paragraphs}
}

func (scorer *Scorer) Calc() {
	scorer.neighbours()
	scorer.keywords()
}

func (scorer *Scorer) neighbours() {
	var paragraphs paragraph.Paragraphs = *scorer.paragraphs
	
	var dimension int = 0
	for _, paragraph := range paragraphs {
		dimension += len(*paragraph.Sentences)
	}

	for _, paragraph := range paragraphs {
		var sentences = *paragraph.Sentences
		for i, s1 := range sentences {
			var score float64 = 0.0
			for j, s2 := range sentences {
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
		paragraph.Sentences = &sentences
	}

	scorer.paragraphs = &paragraphs
}

func (scorer *Scorer) keywords() {
	var paragraphs paragraph.Paragraphs = *scorer.paragraphs

	var tokens []string
	for _, paragraph := range paragraphs {
		var sentences sentence.Sentences = *paragraph.Sentences
		for _, sentence := range sentences {
			tokens = append(tokens, sentence.Tokens...)
		}
	}

	var sum_avg, count_avg int
	var weights map[string]int = make(map[string]int, len(tokens))
	for i, raw_token := range tokens {
		var token string = strings.TrimSpace(raw_token)
		if len(token) > 0 {
			var occurrences int = 0
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
	}

	if count_avg > 0 {
		var keywords []string
		var average int = (sum_avg/count_avg)

		for raw_token, weight := range weights {
			var token string = strings.TrimSpace(raw_token)
			if len(token) > 0 && weight >= average {
				keywords = append(keywords, token)
			}
		}

		for _, paragraph := range paragraphs {
			var sentences = *paragraph.Sentences
			for _, sentence := range sentences {
				var score float64
				for _, keyword := range keywords {
					var weight int = weights[keyword]

					for _, token := range sentence.Tokens {
						if token == keyword {
							score += float64(weight) / float64(len(keywords))
						}
					}
				}
				sentence.Score += score
			}
			paragraph.Sentences = &sentences
		}
	}
	scorer.paragraphs = &paragraphs
}

func (scorer *Scorer) SelectHighlights() []string {
	var paragraphs paragraph.Paragraphs = *scorer.paragraphs
	var sentences_scored sentence.Sentences

	for _, paragraph := range paragraphs {
		var sentences = *paragraph.Sentences
		for _, sentence := range sentences {
			if sentence.Score > 0.0 {
				sentences_scored = append(sentences_scored, sentence)
			}
		}

	}
	sentences_scored.SortScore()

	var count int
	var sum_avg float64
	for _, sentence := range sentences_scored {
		sum_avg += sentence.Score
		count++
	}

	sentences_scored.SortOrder()

	var highlights []string
	var average float64 = (sum_avg/float64(count))
	for _, sentence := range sentences_scored {
		if sentence.Score > average {
			highlights = append(highlights, sentence.Text)
		}
	}

	return highlights
}
