package gobstract

import "sort"

const minSummary int = 5
const summaryPercent float64 = 0.22

// scorer struct contains a set of senteces to score and a limit to generates a
// string extraction summary.
type scorer struct {
	phrases sentences
	limit   int
}

// newScorer function initializes a scorer struct calculating limit according to
// set of sentences provided.
func newScorer(ps sentences) *scorer {
	var s scorer = scorer{ps, minSummary}
	if limit := int(float64(len(ps)) * summaryPercent); limit > minSummary {
		s.limit = limit
	}
	return &s
}

// calcRelations function measures relations between sentences from scorer set
// checking relevant tokens shared. Store the result into weight sentence attr.
func (s *scorer) calcRelations() {
	for i := 0; i < len(s.phrases); i++ {
		var pi sentence = s.phrases[i]
		for j := 0; j < len(s.phrases); j++ {
			var pj sentence = s.phrases[j]
			if i != j && pi.isSimilar(pj) {
				s.phrases[i].weight++
			}
		}
	}
}

// calcLength function fits each sentence weight according to its length. Longer
// sentences usually contains relevant/complete information or explanations.
func (s *scorer) calcLength() {
	var total float64
	for _, p := range s.phrases {
		total += p.length
	}

	var limit float64 = total / float64(len(s.phrases))
	for i, p := range s.phrases {
		if p.length >= limit {
			s.phrases[i].weight *= p.length / limit
		}
	}
}

// calcPosition function, as calcLength function, fits each sentence weight
// according to sentence position into the text. Considers that sentences in
// first and latest positions of the text contains relevant information like
// topic introductions or conclusions.
func (s *scorer) calcPosition() {
	var (
		min   int     = len(s.phrases) / 10
		max   int     = len(s.phrases) - min
		limit float64 = float64(s.limit)
	)

	for i, p := range s.phrases {
		if p.order < min {
			s.phrases[i].weight += p.length * (float64(min-p.order) / limit)
		} else if p.order > max {
			var porder int = len(s.phrases) - p.order
			s.phrases[i].weight += p.length * (float64(min-porder) / limit)
		}
	}
}

// getSummary function invokes calc functions to measuring sentences weight into
// the text. Function chooses sentences until the determined limit sorting by
// calculated weight and store into a sentences subset. Then return that subset
// sorted by original order into the text.
func (s *scorer) getSummary() (r []string) {
	s.calcRelations()
	s.calcLength()
	s.calcPosition()

	var ps sentences = s.phrases
	sort.Sort(byWeight(ps))

	var rs sentences = ps[:s.limit]
	sort.Sort(byOrder(rs))

	for _, p := range rs {
		r = append(r, p.raw)
	}
	return
}
