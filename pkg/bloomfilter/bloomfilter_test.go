package bloomfilter

import "testing"

func TestFilter(t *testing.T) {
	items := [][]byte{
		[]byte("aaa"),
		[]byte("bbb"),
		[]byte("ccc"),
	}

	filter := NewFilter(len(items), 0.001)
	filter.Add(items...)

	input := items[0]
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
