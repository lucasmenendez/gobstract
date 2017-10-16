package gobstract

import (
	"regexp"
	"strings"
)

type Token struct {
	Raw string
	Root string
	Score int
}

type Tokens []*Token

func (t Tokens) Len() int {
	return len(t)
}

func (t Tokens) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t Tokens) Less(i, j int) bool {
	return t[i].Score > t[j].Score
}

func getWords(text string) (words []string) {
	var rgxClean *regexp.Regexp= regexp.MustCompile(`\[|]|\(|\)|\{|}|“|”|«|»|,|´|’|-|_|—|\.\.|:`)
	var cleaned string = rgxClean.ReplaceAllString(text, "")

	var rgxWord = regexp.MustCompile(`\s`)
	var p []string = rgxWord.Split(cleaned, -1)

	for _, r := range p {
		var w string = strings.TrimSpace(r)
		if len(w) > 3 {
			words = append(words, strings.ToLower(w))
		}
	}
	return words
}

func (t *Token) getRoot(lang *Language) {
	var root string = t.Raw
	for _, prefix := range lang.Prefixes {
		if strings.HasPrefix(root, prefix) {
			root = strings.Replace(root, prefix, "", -1)
			break
		} 
	}
	for _, suffix := range lang.Suffixes {
		if strings.HasSuffix(root, suffix) {
			root = strings.Replace(root, suffix, "", -1)
		}
	}

	t.Root = root
}

func (t *Token) isStopword(lang *Language) bool {
	for _, s := range lang.Stopwords {
		if t.Raw == s {
			return true
		}
	}
	return false
}

func (o *Token) diff(t *Token) (diff float64) {
	var distRaw, distRoot int = levenshtain(o.Raw, t.Raw), levenshtain(o.Root, t.Root)
	var lenRaw, lenRoot float64 = float64(len(o.Raw) + len(t.Raw)) / 2.0, float64(len(o.Root) + len(t.Root)) / 2.0

	var diffRaw, diffRoot float64 = float64(distRaw) / lenRaw, float64(distRoot) / lenRoot
	diff = (diffRaw + diffRoot) / 2.0

	//diff = float64(levenshtain(o.Raw, t.Raw)) / lenRaw

	return diff
}

func GetTokens(text string, lang *Language) (tokens []*Token) {
	var words []string = getWords(text)
	for _, w := range words {
		w = strings.TrimSpace(w)
		var t *Token = &Token{Raw: w}

		if t.isStopword(lang) || len(t.Raw) == 0 {
			continue
		}

		t.getRoot(lang)
		tokens = append(tokens, t)
	}

	return tokens
}

func (o *Token) IsIn(tokens []*Token) bool {
	for _, t := range tokens {
		if t.Root == o.Root {
			return true
		}
	}
	return false
}