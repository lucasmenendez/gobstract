package scorer

import (
	"github.com/lucasmenendez/gobstract/paragraph"
	"github.com/lucasmenendez/gobstract/sentence"
	"github.com/lucasmenendez/gobstract/token"
)

var scorable_length = 150

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

	var rang float64 = max_value - min_value
	for _, score := range scores {
		var value float64 = score.value / rang * 0.3
		score.sentence.Score += value
	}
}


func NewScorer(paragraphs *paragraph.Paragraphs) *Scorer {
	var sentences sentence.Sentences
	for _, paragraph := range *paragraphs {
		for _, s := range *paragraph.Sentences {
			if len(s.Text) > scorable_length {
				sentences = append(sentences, s)
			}
		}
	}
	return &Scorer{paragraphs, sentences}
}

func (scorer *Scorer) Calc() {
	scorer.neighbours()
	scorer.keywords()
	scorer.length()
	scorer.order()
}

func (scorer *Scorer) neighbours() {
	var scores []*Score
	for i, sentence1 := range scorer.sentences {
		var score float64
		for j, sentence2 := range scorer.sentences {
			if i != j {
				var coincidences, average float64
				for _, token := range sentence2.Tokens {
					if sentence1.HasToken(token) {
						coincidences += 1.0
					}
				}
				average = float64(len(sentence1.Tokens) + len(sentence2.Tokens))
				score += coincidences / average
			}
		}
		scores = append(scores, &Score{sentence1, score})
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) keywords() {
	var tokens []token.Token
	for _, sentence := range scorer.sentences {
		tokens = append(tokens, sentence.Tokens...)
	}

	for _, paragraph := range *scorer.paragraphs {
		if paragraph.Title != nil {
			tokens = append(tokens, paragraph.Title.Tokens...)
		}
	}

	var max, min int = 0, len(tokens)
	var weights map[token.Token]int = make(map[token.Token]int, len(tokens))
	for i, token := range tokens {
		var occurrences int = 0
		for j, t := range tokens {
			if i != j && token == t {
				occurrences++
			}
		}

		if occurrences > max {
			max = occurrences
		}
		if occurrences < min {
			min = occurrences
		}
		weights[token] = occurrences
	}

	var keywords []token.Token
	var limit int = (max - min) / 3
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
		var value float64 = float64(len(sentence.Text)) / max
		scores = append(scores, &Score{sentence, value * value})
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) titles() {
	var keywords []token.Token
	for _, paragraph := range *scorer.paragraphs {
		if paragraph.Title != nil {
			keywords = append(keywords, paragraph.Title.Tokens...)
		}
	}

	var scores []*Score
	for _, sentence := range scorer.sentences {
		var	occurrences int
		for _, keyword := range keywords {
			if sentence.HasToken(keyword) {
				occurrences++
			}
		}

		if occurrences > 0 {
			var value float64 = float64(occurrences / len(keywords))
			scores = append(scores, &Score{sentence, value})
		}
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) order() {
	var scores []*Score

	var sum int
	for _, paragraph := range *scorer.paragraphs {
		var paragraph_length int = len(*paragraph.Sentences)
		for _, sentence := range *paragraph.Sentences {
			var original_order int = sentence.Order - sum
			if original_order == 1 || original_order == paragraph_length {
				sentence.Score += 0.2
			}
			scores = append(scores, &Score{sentence, float64(original_order)})
		}
		sum += len(*paragraph.Sentences)
	}
	scorer.addScores(scores)
}

func (scorer *Scorer) SelectHighlights() []string {
	var sum_avg float64
	var sentences_scored sentence.Sentences
	for _, sentence := range scorer.sentences {
		if sentence.Score > 0.0 {
			sentences_scored = append(sentences_scored, sentence)
			sum_avg += sentence.Score
		}
	}

	sentences_scored.SortOrder()

	var highlights []string
	var avg float64 = sum_avg / float64(len(sentences_scored))
	for _, sentence := range sentences_scored {
		if sentence.Score > avg {
			highlights = append(highlights, sentence.Text)
		}
	}
	return highlights
}
