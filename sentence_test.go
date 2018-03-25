package gobstract

import "testing"

func TestIsSimilar(t *testing.T) {
	var ss sentences = sentences{
		{tokens: []string{"st-token", "nd-token"}, lengthTokens: 2},
		{tokens: []string{"st-token", "nd-token", "rd-token"}, lengthTokens: 3},
		{tokens: []string{"xxx", "cccc", "zzzz"}, lengthTokens: 2},
	}

	if !ss[0].isSimilar(ss[1]) {
		t.Error("Expected true, got false")
	}

	if ss[0].isSimilar(ss[2]) {
		t.Error("Expected false, got true")
	}
}
