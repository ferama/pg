package history

import (
	"testing"
)

func TestHistory(t *testing.T) {
	h := GetInstance()
	h.Append("1")
	h.Append("2")
	h.Append("3")

	_, err := h.GoNext()
	if err == nil {
		t.Fail()
	}

	p, err := h.GoPrev()
	if err != nil {
		t.Fatal(err)
	}
	if p != "2" {
		t.Fail()
	}
}
