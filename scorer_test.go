package gobstract

import (
	"testing"
)

func TestNewScorer(t *testing.T) {
	var s *scorer = newScorer(make(sentences, 10))
	if s.limit != minSummary {
		t.Errorf("Expected %d (min), got %d", minSummary, s.limit)
	}

	s = newScorer(make(sentences, 50))
	var expected int = int(50 * summaryPercent)
	if s.limit != expected {
		t.Errorf("Expected %d, got %d", expected, s.limit)
	}
}

func TestCalcRelations(t *testing.T) {
	var ps sentences = sentences{
		{tokens: []string{"aaaa"}, lengthTokens: 1},
		{tokens: []string{"aaaa"}, lengthTokens: 1},
		{tokens: []string{"bbbb"}, lengthTokens: 1},
		{tokens: []string{"cccc"}, lengthTokens: 1},
		{tokens: []string{"dddd"}, lengthTokens: 1},
	}

	var s *scorer = newScorer(ps)
	s.calcRelations()
	if s.phrases[0].weight != 1 {
		t.Errorf("Expected 1, got %f", s.phrases[0].weight)
	}

	if s.phrases[1].weight != 1 {
		t.Errorf("Expected 1, got %f", s.phrases[1].weight)
	}

	for _, p := range s.phrases[2:] {
		if p.weight != 0 {
			t.Errorf("Expected 0, got %f", p.weight)
		}
	}

	s = newScorer(ps[2:])
	s.calcRelations()

	for _, p := range s.phrases {
		if p.weight != 0 {
			t.Errorf("Expected 0, got %f", p.weight)
		}
	}
}

func TestCalcLength(t *testing.T) {
	var ps sentences = sentences{
		{lengthRaw: 170, weight: 3},
		{lengthRaw: 140, weight: 3},
		{lengthRaw: 90, weight: 1},
		{lengthRaw: 90, weight: 1},
		{lengthRaw: 90, weight: 1},
	}

	var s *scorer = newScorer(ps)
	s.calcLength()

	if 3 >= s.phrases[0].weight {
		t.Errorf("Expected > 3, got %f", s.phrases[0].weight)
	}

	if 3 >= s.phrases[1].weight {
		t.Errorf("Expected > 3, got %f", s.phrases[1].weight)
	}

	for _, p := range s.phrases[2:] {
		if p.weight != 1 {
			t.Errorf("Expected 1, got %f", p.weight)
		}
	}

	ps = sentences{
		{lengthRaw: 170, weight: 0},
		{lengthRaw: 140, weight: 0},
		{lengthRaw: 90, weight: 0},
		{lengthRaw: 90, weight: 0},
		{lengthRaw: 90, weight: 0},
	}
	s = newScorer(ps)
	s.calcLength()

	for _, p := range s.phrases {
		if p.weight != 0 {
			t.Errorf("Expected 0, got %f", p.weight)
		}
	}
}

func TestCalcPosition(t *testing.T) {
	var ps sentences = sentences{
		{weight: 3, order: 1},
		{weight: 3, order: 2},
		{weight: 3, order: 3},
		{weight: 1, order: 4},
		{weight: 1, order: 5},
		{weight: 1, order: 6},
		{weight: 1, order: 7},
		{weight: 3, order: 8},
		{weight: 3, order: 9},
		{weight: 3, order: 10},
	}

	var s *scorer = newScorer(ps)
	s.calcPosition()

	for _, p := range s.phrases[3:7] {
		if p.weight != 1 {
			t.Errorf("Expected 1, got %f", p.weight)
		}
	}

	for _, p := range append(s.phrases[:3], s.phrases[7:]...) {
		if p.weight <= 3 {
			t.Errorf("Expected > 3, got %f", p.weight)
		}
	}
}
