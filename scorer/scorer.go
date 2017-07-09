package scorer

import (
	"fmt"
	"github.com/lucasmenendez/gobstract/paragraph"
	"github.com/lucasmenendez/gobstract/sentence"
)

type Score struct {
	sentence *sentence.Sentence
	value float64
}

type Scorer struct {
	paragraphs *paragraph.Paragraphs
	sentences sentence.Sentences
}

func (scorer *Scorer) addScores(scores []*Score) {
	var max_value, min_value float64
	for _, score := range scores {
		if score.value > max_value {
			max_value = score.value
		} else if score.value < min_value {
			min_value = score.value
		}
	}

	var rang float64 = (max_value - min_value)
	for _, score := range scores {
		var value float64 = (score.value/rang*0.5)
		score.sentence.Score += value	
	}
}


func NewScorer(paragraphs *paragraph.Paragraphs) *Scorer {
	var sentences sentence.Sentences
	for _, paragraph := range *paragraphs {
		for _, s := range *paragraph.Sentences {
			sentences = append(sentences, s)
		}
	}
	return &Scorer{paragraphs, sentences}
}

func (scorer *Scorer) Calc() {
	scorer.neighbours()
	scorer.keywords()
	scorer.length()
}

func (scorer *Scorer) neighbours() {
	var scores []*Score
	for i, s1 := range scorer.sentences {
		var score float64
		for j, s2 := range scorer.sentences {
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
		scores = append(scores, &Score{s1, score})
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) keywords() {
	var tokens []string
	for _, sentence := range scorer.sentences {
		for _, token := range sentence.Tokens {
			if len(token) >= 3 {
				tokens = append(tokens, token)
			}
		}
	}

	var max, min int = 0, len(tokens)
	var weights map[string]int = make(map[string]int, len(tokens))
	for i, token := range tokens {
		var occurrences int = 0
		for j, t := range tokens {
			if i != j && token == t {
				occurrences++
			}
		}

		if occurrences > max {
			max = occurrences
		} else if occurrences < min {
			min = occurrences
		}
		weights[token] = occurrences
	}

	var keywords []string
	var limit int = ((max - min)/2) 
	for token, weight := range weights {
		if weight >= limit {
			keywords = append(keywords, token)
		}
	}
	
	var scores []*Score
	for _, sentence := range scorer.sentences {
		var value float64
		for _, keyword := range keywords {
			var weight int = weights[keyword]
			if sentence.HasToken(keyword) {
				value += float64(weight) / float64(len(keywords))
			}
		}
		scores = append(scores, &Score{sentence, value})
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) length() {
	var max float64
	for _, sentence := range scorer.sentences {
		if value := float64(len(sentence.Text)); value > max {
			max = value
		}
	}

	var scores []*Score
	for _, sentence := range scorer.sentences {
		var value float64 = (float64(len(sentence.Text)) / max) 
		scores = append(scores, &Score{sentence, value})
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) title() {
	for _, paragraph := range *scorer.paragraphs {
		fmt.Println(paragraph.Title.Tokens)
	}
}

func (scorer *Scorer) SelectHighlights() []string {
	var sentences_scored sentence.Sentences
	for _, sentence := range scorer.sentences {
		if sentence.Score > 0.0 {
			sentences_scored = append(sentences_scored, sentence)
		}
	}
	sentences_scored.SortScore()

	var max, min float64 = 0.0, float64(len(sentences_scored))
	for _, sentence := range sentences_scored {
		if sentence.Score > max {
			max = sentence.Score
		} else if sentence.Score < min {
			min = sentence.Score
		}
	}

	sentences_scored.SortOrder()

	var highlights []string
	var limit float64 = ((max - min) / 2.0)
	for _, sentence := range sentences_scored {
		if sentence.Score > limit {
			highlights = append(highlights, sentence.Text)
		}
	}
	return highlights
}
