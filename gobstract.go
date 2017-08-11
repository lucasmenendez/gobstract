package gobstract

type Gobstract struct {
	Text		string
	Paragraphs	*Paragraphs
	Sentences	[]string
	Scorer		*Scorer
	Lang		*Language
}

func NewAbstract(text string, lang_label string) (*Gobstract, error) {
	var sentences []string

	var err error
	var lang *Language
	if lang, err = GetLanguage(lang_label); err != nil {
		return nil, err
	}

	var paragraphs *Paragraphs = SplitText(text, lang)
	var g *Gobstract = &Gobstract{Text: text, Paragraphs: paragraphs, Sentences: sentences, Lang: lang}

	g.Scorer = NewScorer(g.Paragraphs)
	g.Scorer.Calc()

	return g, nil
}

func (g *Gobstract) GetBestSentence() string {
	return g.Scorer.SelectBestSentence()
}

func (g *Gobstract) GetHightlights() []string {
	return g.Scorer.SelectHighlights()
}

func (g *Gobstract) GetKeywords() []string {
	return g.Scorer.SelectKeywords()
}