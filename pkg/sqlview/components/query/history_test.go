package query

import (
	"testing"
)

func TestHistory(t *testing.T) {
	h := newHistory()
	h.append("1")
	h.append("2")
	h.append("3")

	_, err := h.getNext()
	if err == nil {
		t.Fail()
	}

	p, err := h.getPrev()
	if err != nil {
		t.Fatal(err)
	}
	if p != "2" {
		t.Fail()
	}
}
