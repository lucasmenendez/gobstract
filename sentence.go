package gobstract

const tokenSimilarityThreshold float64 = 0.55
const sentenceSimilarityThreshold float64 = 0.7

// sentence struct evolves sentence data such as relevant tokens, raw content,
// sentence order and weight into the text. Also stores tokens length and raw
// content length in float64 to use them across the algorithm.
type sentence struct {
	raw          string
	lengthRaw    float64
	tokens       []string
	lengthTokens float64
	weight       float64
	order        int
}

// isSimilar function determine if provided sentence is similar to referenced
// sentence. Check similarity between both sentences relevant tokens. If
// accumulated similarity is greater than sentenceSimilarityThreshold const,
// both sentences are definitely similar.
func (s sentence) isSimilar(s2 sentence) bool {
	var c float64
	for _, t1 := range s.tokens {
		var d float64
		for _, t2 := range s2.tokens {
			d += strDistance(t1, t2)
		}

		if d /= s2.lengthTokens; d > tokenSimilarityThreshold {
			c++
		}
	}

	return c/s.lengthTokens > sentenceSimilarityThreshold
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
