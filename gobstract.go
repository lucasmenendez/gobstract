// Package gobstract provide simple abstract algorithm and keywords extractor
// using simple linguistic weight scorers
package gobstract

const (
	titleMin int = 5
	titleMax int = 100
	scorableLength int = 100
	maxKeywords int = 6
	maxLevenshtain = 0.5
)

type Gobstract struct {
	Text		string
	Paragraphs	*Paragraphs
	Sentences	[]string
	Scorer		*Scorer
	Lang		*Language
}

func NewAbstract(text string, lang string) (*Gobstract, error) {
	var l *Language
	var e error
	if l, e = GetLanguage(lang); e != nil {
		return nil, e
	}

	var s []string
	var p *Paragraphs = SplitText(text, l)
	var g *Gobstract = &Gobstract{Text: text, Paragraphs: p, Sentences: s, Lang: l}

	g.Scorer = NewScorer(g.Paragraphs)
	g.Scorer.Calc()

	return g, nil
}

func (g *Gobstract) GetBestSentence() string {
	return g.Scorer.SelectBestSentence()
}

func (g *Gobstract) GetHightlights(max ...int) []string {
	var l int = -1
	if len(max) == 1 {
		l = max[0]
	}

	return g.Scorer.SelectHighlights(l)
}

func (g *Gobstract) GetKeywords() []string {
	return g.Scorer.SelectKeywords()
}