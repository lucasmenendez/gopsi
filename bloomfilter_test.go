package psi

import "testing"

func TestFilter(t *testing.T) {
	items := []string{ "aaa", "bbb", "ccc" }

	filter := NewFilter(len(items), 0.001)
	for _, item := range items {
		filter.Add([]byte(item))
	}

	input := []byte(items[0])
	if !filter.Test(input) {
		t.Errorf("Expected that filter contains '%s'.", input)
	}

	input = []byte("aa")
	if filter.Test(input) {
		t.Errorf("Expected that filter not contains '%s'.", input)
	}

	input = []byte("aaaa")
	if filter.Test(input) {
		t.Errorf("Expected that filter not contains '%s'.", input)
	}

	input = []byte("")
	if filter.Test(input) {
		t.Errorf("Expected that filter not contains '%s'.", input)
	}
}