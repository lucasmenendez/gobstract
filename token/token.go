package token

import (
	"regexp"
	"strings"
	"github.com/lucasmenendez/gobstract/language"
)

type Token struct {
	Raw string
	root string
}

func GetTokens(text string, lang *language.Language) []Token {
	var tokens []Token
	
	var words []string = getWords(text)
	for _, word := range words {
		var token Token = Token{Raw: word}
		if token.isStopword(lang) {
			continue
		}

		token.getroot(lang)
		tokens = append(tokens, token)
	}

	return tokens
}

func (needle Token) IsIn(tokens []Token) bool {
	for _, token := range tokens {
		if token.root == needle.root {
			return true
		}
	}
	return false
}

func getWords(text string) []string {
	var words []string

	var rgx_clean *regexp.Regexp= regexp.MustCompile(`\[|]|\(|\)|\{|}|“|”|«|»|,|´|’|-|_|—|\.\.|:`)
	var cleaned string = rgx_clean.ReplaceAllString(text, "")

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

func (token Token) getroot(lang *language.Language) {
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

	token.root = root
}

func (token Token) isStopword(lang *language.Language) bool {
	for _, stopword := range lang.Stopwords {
		if token.Raw == stopword {
			return true
		}
	}
	return false
}
