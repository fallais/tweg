package tweg

import (
	"testing"
)

// Encode a message
func TestEncode(t *testing.T) {
	tw := NewTweg()

	result := tw.Encode("A koala arrives in the great forest of Wumpalumpa", "alpaga")
	if result != "A kｏａla arrivｅs іn the great forest of Wumpalumpa" {
		t.Fatal("String is incorrect")
	}
}
