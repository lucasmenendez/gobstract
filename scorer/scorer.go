package scorer

import (
	"strings"
	"github.com/lucasmenendez/gobstract/sentence"
)

type Scorer struct {
	paragraphs []sentence.Sentences
}

func NewScorer(paragraphs []sentence.Sentences) *Scorer {
	return &Scorer{paragraphs}
}

func (scorer *Scorer) Calc() {
	scorer.neighbours()
	scorer.keywords()
}

func (scorer *Scorer) neighbours() {
	var dimension int = 0
	for _, paragraph := range scorer.paragraphs {
		dimension += len(paragraph)
	}

	for _, paragraph := range scorer.paragraphs {
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

func (scorer *Scorer) keywords() {
	var tokens []string
	for _, paragraph := range scorer.paragraphs {
		for _, sentence := range paragraph {
			tokens = append(tokens, sentence.Tokens...)
		}
	}

	var sum_avg, count_avg, max_avg int
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
				if occurrences > max_avg {
					max_avg = occurrences
				}
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

		for _, paragraph := range scorer.paragraphs {
			for _, sentence := range paragraph {
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
		}
	}
}
