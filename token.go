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

func getWords(text string) []string {
	var words []string

	var rgxClean *regexp.Regexp= regexp.MustCompile(`\[|]|\(|\)|\{|}|“|”|«|»|,|´|’|-|_|—|\.\.|:`)
	var cleaned string = rgxClean.ReplaceAllString(text, "")

	var rgx_word = regexp.MustCompile(`\s`)
	var parts []string = rgx_word.Split(cleaned, -1)

	for _, raw_word := range parts {
		var word string = strings.TrimSpace(raw_word)
		if len(word) > 3 {
			words = append(words, strings.ToLower(word))
		}
	}
	return words
}

func (token *Token) getRoot(lang *Language) {
	var root string = token.Raw
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

	token.Root = root
	//token.root = token.Raw
}

func (token *Token) isStopword(lang *Language) bool {
	for _, stopword := range lang.Stopwords {
		if token.Raw == stopword {
			return true
		}
	}
	return false
}

func GetTokens(text string, lang *Language) []*Token {
	var tokens []*Token

	var words []string = getWords(text)
	for _, word := range words {
		word = strings.TrimSpace(word)
		var token *Token = &Token{Raw: word}
		if token.isStopword(lang) || len(token.Raw) == 0 {
			continue
		}

		token.getRoot(lang)
		tokens = append(tokens, token)
	}

	return tokens
}

func (needle *Token) IsIn(tokens []*Token) bool {
	for _, token := range tokens {
		if token.Root == needle.Root {
			return true
		}
	}
	return false
}
