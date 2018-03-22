package gobstract

const tokenSimilarityThreshold float64 = 0.55
const sentenceSimilarityThreshold float64 = 0.2

type sentence struct {
	tokens []string
	raw    string
	weight float64
	length float64
	order  int
}

func (s sentence) isSimilar(s2 sentence) bool {
	var c float64
	for _, t1 := range s.tokens {
		var d float64
		var rate float64 = float64(len(t1)) / float64(len(s.tokens))
		for _, t2 := range s2.tokens {
			d += strSimilarity(t1, t2) * rate
		}

		if d > tokenSimilarityThreshold {
			c++
		}
	}
	return c/float64(len(s.tokens)) > sentenceSimilarityThreshold
}

type sentences []sentence

type byWeight []sentence

func (s byWeight) Len() int           { return len(s) }
func (s byWeight) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byWeight) Less(i, j int) bool { return s[i].weight > s[j].weight }

type byOrder []sentence

func (s byOrder) Len() int           { return len(s) }
func (s byOrder) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byOrder) Less(i, j int) bool { return s[i].order < s[j].order }
