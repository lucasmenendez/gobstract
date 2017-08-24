package gobstract

import (
	"sort"
	"strings"
)

const (
	scorableLength int = 100
	maxKeywords int = 5
)

type Score struct {
	sentence *Sentence
	value float64
}

type Scorer struct {
	paragraphs 	*Paragraphs
	sentences 	Sentences
}

func (scorer *Scorer) addScores(scores []*Score) {
	var maxValue, minValue float64
	for _, score := range scores {
		if score.value > maxValue {
			maxValue = score.value
		} else if score.value < minValue {
			minValue = score.value
		}
	}

	var rang float64 = maxValue - minValue
	for _, score := range scores {
		var value float64 = score.value / rang * 0.3
		score.sentence.Score += value
	}
}

func (scorer *Scorer) addNeighbours() {
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

func (scorer *Scorer) addKeywords() {
	var tokens Tokens
	for _, sentence := range scorer.sentences {
		tokens = append(tokens, sentence.Tokens...)
	}

	for _, paragraph := range *scorer.paragraphs {
		if paragraph.Title != nil {
			tokens = append(tokens, paragraph.Title.Tokens...)
		}
	}

	var max, min int = 0, len(tokens)
	for i, token := range tokens {
		var occurrences int = 0
		for j, t := range tokens {
			if i != j && token.Root == t.Root {
				occurrences++
			}
		}

		if occurrences > max {
			max = occurrences
		}
		if occurrences < min {
			min = occurrences
		}
		token.Score = occurrences
	}

	var keywords Tokens
	var limit int = (max - min) / 3
	for _, token := range tokens {
		if token.Score >= limit {
			keywords = append(keywords, token)
		}
	}


	var scores []*Score
	for _, sentence := range scorer.sentences {
		var value float64
		for _, keyword := range keywords {
			if sentence.HasToken(keyword) {
				value += float64(keyword.Score) / float64(len(keywords))
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
	var keywords Tokens
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
			var oldOrder int = sentence.Order - sum
			if oldOrder == 1 || oldOrder == paragraph_length {
				sentence.Score += 0.2
			}
			scores = append(scores, &Score{sentence: sentence, value: float64(oldOrder)})
		}
		sum += len(*paragraph.Sentences)
	}
	scorer.addScores(scores)
}

func NewScorer(paragraphs *Paragraphs) *Scorer {
	var sentences Sentences

	for _, paragraph := range *paragraphs {
		for _, s := range *paragraph.Sentences {
			if len(s.Text) > scorableLength {
				sentences = append(sentences, s)
			}
		}
	}
	return &Scorer{paragraphs: paragraphs, sentences: sentences}
}

func (scorer *Scorer) Calc() {
	scorer.addNeighbours()
	scorer.addKeywords()
	scorer.length()
	scorer.order()
}

func (scorer *Scorer) SelectBestSentence() string {
	scorer.sentences.SortScore()
	if len(scorer.sentences) > 0 {
		return scorer.sentences[0].Text
	}

	return ""
}

func (scorer *Scorer) SelectHighlights(max int) []string {
	var sum_avg float64
	var sentences_scored Sentences
	for _, sentence := range scorer.sentences {
		if sentence.Score > 0.0 {
			sentences_scored = append(sentences_scored, sentence)
			sum_avg += sentence.Score
		}
	}

	if max > -1 {
		sentences_scored.SortScore()
		sentences_scored = sentences_scored[:max]
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

func (s *Scorer) SelectKeywords() []string {
	var tokens Tokens
	for _, sentence := range s.sentences {
		tokens = append(tokens, sentence.Tokens...)
	}

	for _, paragraph := range *s.paragraphs {
		if paragraph.Title != nil {
			tokens = append(tokens, paragraph.Title.Tokens...)
		}
	}

	sort.Sort(tokens)
	var keywords Tokens
	for _, t := range tokens {
		var exists bool = false
		for _, keyword := range keywords {
			if t.Root == keyword.Root {
				exists = true
				break
			}
		}

		if !exists {
			keywords = append(keywords, t)
		}

		if len(keywords) == maxKeywords {
			break
		}
	}


	var raw_keywords []string
	for _, keyword := range keywords {
		if strings.TrimSpace(keyword.Raw) != "" {
			raw_keywords = append(raw_keywords, keyword.Raw)
		}
	}

	return raw_keywords
}
