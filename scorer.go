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

func (sc *Scorer) addScores(scores []*Score) {
	var max, min float64
	for _, score := range scores {
		if score.value > max {
			max = score.value
		} else if score.value < min {
			min = score.value
		}
	}

	var rang float64 = max - min
	for _, s := range scores {
		s.sentence.Score += s.value / rang * 0.3
	}
}

func (sc *Scorer) addNeighbours() {
	var scores []*Score
	for i, s1 := range sc.sentences {
		var score float64
		for j, s2 := range sc.sentences {
			if i != j {
				var o, avg float64
				for _, token := range s2.Tokens {
					if s1.HasToken(token) {
						o += 1.0
					}
				}
				avg = float64(len(s1.Tokens) + len(s2.Tokens))
				score += o / avg
			}
		}
		scores = append(scores, &Score{s1, score})
	}
	sc.addScores(scores)
}

func (sc *Scorer) addKeywords() {
	var tokens Tokens
	for _, s := range sc.sentences {
		tokens = append(tokens, s.Tokens...)
	}

	for _, p := range *sc.paragraphs {
		if p.Title != nil {
			tokens = append(tokens, p.Title.Tokens...)
		}
	}

	var max, min int = 0, len(tokens)
	for i, token := range tokens {
		var o int = 0
		for j, t := range tokens {
			if i != j && token.Root == t.Root {
				o++
			}
		}

		if o > max {
			max = o
		}
		if o < min {
			min = o
		}
		token.Score = o
	}

	var keywords Tokens
	var l int = (max - min) / 3
	for _, t := range tokens {
		if t.Score >= l {
			keywords = append(keywords, t)
		}
	}


	var scores []*Score
	for _, s := range sc.sentences {
		var v float64
		for _, k := range keywords {
			if s.HasToken(k) {
				v += float64(k.Score) / float64(len(keywords))
			}
		}

		scores = append(scores, &Score{s, v})
	}

	sc.addScores(scores)
}

func (sc *Scorer) length() {
	var max float64
	for _, s := range sc.sentences {
		if v := float64(len(s.Text)); v > max {
			max = v
		}
	}

	var scores []*Score
	for _, s := range sc.sentences {
		var value float64 = float64(len(s.Text)) / max
		scores = append(scores, &Score{s, value * value})
	}
	sc.addScores(scores)
}

func (sc *Scorer) titles() {
	var keywords Tokens
	for _, p := range *sc.paragraphs {
		if p.Title != nil {
			keywords = append(keywords, p.Title.Tokens...)
		}
	}

	var scores []*Score
	for _, s := range sc.sentences {
		var	o int
		for _, k := range keywords {
			if s.HasToken(k) {
				o++
			}
		}

		if o > 0 {
			var v float64 = float64(o / len(keywords))
			scores = append(scores, &Score{s, v})
		}
	}
	sc.addScores(scores)
}

func (sc *Scorer) order() {
	var scores []*Score

	var sum int
	for _, p := range *sc.paragraphs {
		var l int = len(*p.Sentences)
		for _, s := range *p.Sentences {
			var o int = s.Order - sum
			if o == 1 || o == l {
				s.Score += 0.2
			}
			scores = append(scores, &Score{sentence: s, value: float64(o)})
		}
		sum += len(*p.Sentences)
	}
	sc.addScores(scores)
}

func NewScorer(paragraphs *Paragraphs) *Scorer {
	var sentences Sentences

	for _, p := range *paragraphs {
		for _, s := range *p.Sentences {
			if len(s.Text) > scorableLength {
				sentences = append(sentences, s)
			}
		}
	}

	for i, s1 := range sentences {
		for j, s2 := range sentences {
			if i != j && s1.Text == s2.Text {
				sentences.Delete(i)
				break
			}
		}
	}

	return &Scorer{paragraphs: paragraphs, sentences: sentences}
}

func (sc *Scorer) Calc() {
	sc.addNeighbours()
	sc.addKeywords()
	sc.length()
	sc.order()
}

func (sc *Scorer) SelectBestSentence() string {
	sc.sentences.SortScore()
	if len(sc.sentences) > 0 {
		return sc.sentences[0].Text
	}

	return ""
}

func (sc *Scorer) SelectHighlights(max int) (highlights []string) {
	var sumAvg float64
	var scored Sentences
	for _, s := range sc.sentences {
		if s.Score > 0.0 {
			scored = append(scored, s)
			sumAvg += s.Score
		}
	}

	var avg float64 = sumAvg / float64(len(scored))
	if max > -1 && max < len(scored) {
		scored.SortScore()
		scored = scored[:max]
	}

	scored.SortOrder()
	for _, s := range scored {
		if s.Score > avg {
			highlights = append(highlights, s.Text)
		}
	}

	return highlights
}

func (sc *Scorer) SelectKeywords() (keywords []string) {
	var tokens Tokens
	for _, s := range sc.sentences {
		tokens = append(tokens, s.Tokens...)
	}

	for _, p := range *sc.paragraphs {
		if p.Title != nil {
			tokens = append(tokens, p.Title.Tokens...)
		}
	}

	sort.Sort(tokens)
	var raw Tokens
	for _, t := range tokens {
		var exists bool = false
		for _, keyword := range raw {
			if t.Root == keyword.Root {
				exists = true
				break
			}
		}

		if !exists {
			raw = append(raw, t)
		}

		if len(raw) == maxKeywords {
			break
		}
	}


	for _, keyword := range raw {
		if strings.TrimSpace(keyword.Raw) != "" {
			keywords = append(keywords, keyword.Raw)
		}
	}

	return keywords
}
