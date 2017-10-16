package gobstract

import (
	"sort"
	"strings"
	"fmt"
)

type Score struct {
	sentence *Sentence
	value float64
}

type Scorer struct {
	paragraphs 	*Paragraphs
	sentences 	Sentences
	tags		Tokens
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

func (sc *Scorer) neighbours() {
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

func (sc *Scorer) keywords() {
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
		var included bool = false
		for _, k := range keywords {
			if t.diff(k) < maxLevenshtain {
				included = true
				break
			}
		}

		if t.Score >= l && !included {
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

	sc.tags = keywords
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
			if s.HasSimilarToken(k) {
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

func (sc *Scorer) printLevenshtain() {
	for _, sp := range sc.sentences {
		var s Sentence = *sp
		var matrix [][]float64 = make([][]float64, len(s.Tokens))
		for i := 0; i < len(matrix); i++ {
			matrix[i] = make([]float64, len(s.Tokens))
		}

		fmt.Println(s.Text)

		for i := 0; i < len(s.Tokens); i++ {
			var t1 *Token = s.Tokens[i]
			for j := 0; j < len(s.Tokens); j++ {
				var t2 *Token = s.Tokens[j]

				matrix[i][j] = t1.diff(t2)
			}
		}

		fmt.Print(" NULL \t")
		for _, t := range s.Tokens {
			fmt.Printf("%s\t", t.Raw)
		}
		fmt.Println("")

		for i := 0; i < len(s.Tokens); i++ {
			var t *Token = s.Tokens[i]
			fmt.Printf("%s\t", t.Raw)
			for j := 0; j < len(s.Tokens); j++ {
				fmt.Printf("%f\t", matrix[i][j])
			}
			fmt.Println("")
		}
	}
}

func NewScorer(paragraphs *Paragraphs) *Scorer {
	var a Sentences
	for _, p := range *paragraphs {
		for _, s := range *p.Sentences {
			if len(s.Text) > scorableLength {
				a = append(a, s)
			}
		}
	}

	var la int = len(a)
	for i := 0; i < la; i++ {
		var x *Sentence = a[i]
		for j := 0; j < la; j++ {
			var y *Sentence = a[j]
			if i != j && x.Text == y.Text {
				a.Delete(i)
				la--
				break
			}
		}
	}

	var s Sentences
	for _, x := range a {
		if x != nil {
			s = append(s, x)
		}
	}

	return &Scorer{paragraphs: paragraphs, sentences: s}
}

func (sc *Scorer) Calc() {
	sc.neighbours()
	sc.keywords()

	sc.titles()
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
	var tokens Tokens = sc.tags

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
