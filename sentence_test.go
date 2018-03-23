package gobstract

import "testing"

func TestIsSimilar(t *testing.T) {
	var ss sentences = sentences{
		{ tokens: []string{ "first-token", "second-token" } },
		{ tokens: []string{ "first-token", "second-token", "third-token" } },
		{ tokens: []string{ "xxx", "cccc", "zzzz" } },
	}

	if !ss[0].isSimilar(ss[1]) {
		t.Error("Expected true, got false")
	}

	if ss[0].isSimilar(ss[2]) {
		t.Error("Expected false, got true")
	}
}