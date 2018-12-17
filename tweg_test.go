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

	result = tw.Encode("i had a great day at the beach! #sunshine                ", "kidnapped by pirates")
	if result != "i haｄ a grｅａｔ daｙ at the ｂeaｃh! #sunshｉne                " {
		t.Fatal("String is incorrect")
	}
}

// Decode a message
func TestDecode(t *testing.T) {
	tw := NewTweg()

	result := tw.Decode("A kｏａla arrivｅs іn the great forest of Wumpalumpa")
	if result != "alpaga        " {
		t.Fatal("String is incorrect")
	}

	result = tw.Decode("i haｄ a grｅａｔ daｙ at the ｂeaｃh! #sunshｉne                ")
	if result != "kidnapped by pirates   " {
		t.Fatal("String is incorrect")
	}
}
