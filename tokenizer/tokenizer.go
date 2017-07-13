package tokenizer

import (
	"regexp"
	"strings"
)

func SplitSentences (input string) []string {
	var titles_pattern *regexp.Regexp = regexp.MustCompile(`([A-Z][a-z]+)\.`)
	var titles_needle string = `$1*|*`
	var no_titles string = titles_pattern.ReplaceAllString(input, titles_needle)

	var numbers_pattern *regexp.Regexp = regexp.MustCompile(`([0-9]+)\.([0-9]+)`)
	var numbers_needle string = `$1*|*$2`
	var no_numbers string = numbers_pattern.ReplaceAllString(no_titles, numbers_needle)

	var stops_pattern *regexp.Regexp = regexp.MustCompile(`[^..][!?.]\s`)
	var stops_needle string = `$0{stop}`
	var no_stops string = stops_pattern.ReplaceAllString(no_numbers, stops_needle)

	var restore_pattern *regexp.Regexp = regexp.MustCompile(`\*\|\*`)
	var restore_needle = `.`
	var text string = restore_pattern.ReplaceAllString(no_stops, restore_needle)

	var spliter *regexp.Regexp = regexp.MustCompile(`{stop}`)
	return spliter.Split(text, -1)
}

func SplitWords (input string) []string {
	var rgx_clean *regexp.Regexp= regexp.MustCompile(`\[|]|\(|\)|\{|}|“|”|«|»|,|´|’|-|_|—|\.\.|:`)
	var cleaned string = rgx_clean.ReplaceAllString(input, "")

	var rgx_word = regexp.MustCompile(`\s`)
	var tokens []string = rgx_word.Split(cleaned, -1)

	for _, raw_token := range tokens {
		var token string = strings.TrimSpace(raw_token)
		if len(token) > 3 {
			tokens = append(tokens, strings.ToLower(token))
		}
	}

	return tokens
}