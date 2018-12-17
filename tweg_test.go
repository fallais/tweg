package tweg

import (
	"testing"
)

// Encode a message
func TestEncode(t *testing.T) {
	tw := NewTweg()

	result, err := tw.Encode("A koala arrives in the great forest of Wumpalumpa", "alpaga")
	if err != nil {
		t.Fatal("Error while encoding :", err)
	}
	if result != "A kｏａla arrivｅs іn the great forest of Wumpalumpa" {
		t.Fatal("String is incorrect")
	}

	result, err = tw.Encode("i had a great day at the beach! #sunshine                ", "kidnapped by pirates")
	if err != nil {
		t.Fatal("Error while encoding :", err)
	}
	if result != "i haｄ a grｅａｔ daｙ at the ｂeaｃh! #sunshｉne                " {
		t.Fatal("String is incorrect")
	}
}

// Decode a message
func TestDecode(t *testing.T) {
	tw := NewTweg()

	result, err := tw.Decode("A kｏａla arrivｅs іn the great forest of Wumpalumpa")
	if err != nil {
		t.Fatal("Error while decoding :", err)
	}
	if result != "alpaga        " {
		t.Fatal("String is incorrect")
	}

	result, err = tw.Decode("i haｄ a grｅａｔ daｙ at the ｂeaｃh! #sunshｉne                ")
	if err != nil {
		t.Fatal("Error while decoding :", err)
	}
	if result != "kidnapped by pirates   " {
		t.Fatal("String is incorrect")
	}

	result, err = tw.Decode("Ｃhｏose  a  jοｂ  yоu  lονｅ,  and  you  ｗіｌl  ｎeｖｅｒ  have  tο  ｗｏrk  a  day  in  yοur  lіfｅ．                        ")
	if err != nil {
		t.Fatal("Error while decoding :", err)
	}
	if result != "rendezvous at grand central terminal on friday.    " {
		t.Fatal("String is incorrect")
	}
}
