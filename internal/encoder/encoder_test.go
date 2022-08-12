package encoder

import (
	"testing"
)

func TestStrToIntsToStr(t *testing.T) {
	var str = "Lorem ipsum dolor sit amet, consectetur adipiscing elit.\nDonec molestie justo eget leo convallis ullamcorper.\nNam eros enim, dapibus euismod sodales eget, condimentum id enim.\nPraesent ornare feugiat ultrices.\nDonec tortor velit, ornare a interdum at, viverra et urna."

	resultInts := StrToInts(str)
	resultStr := IntsToStr(resultInts)

	if str != resultStr {
		t.Errorf("Expected '%s', got '%s'", resultInts, resultStr)
		return
	}
}
