package tweg

import (
	"fmt"
	"strconv"
)

// Lookup the homoglyphs
func (t *Tweg) lookup() {
	for c, h := range Homoglyphs {
		homoglyphOptionsBitLength := len(strconv.FormatInt(int64(len(h)+1), 2)) - 1
		t.HomoglyphsLookup[c] = ""
		for i := 0; i < homoglyphOptionsBitLength; i++ {
			t.HomoglyphsLookup[c] += "0"
		}

		// Options
		i := 0
		for _, o := range h {
			characterCodeInDecimal, err := strconv.ParseInt(o, 16, 64)
			if err != nil {
				fmt.Println("Error while parsing the characterCodeInDecimal")
			}
			homoglyphOptionCharacter := string(characterCodeInDecimal)
			t.HomoglyphsLookup[homoglyphOptionCharacter] = zeropadding(strconv.FormatInt(int64(i+1), 2), homoglyphOptionsBitLength)
			i++
		}
	}
}

// ensureDivisible by the secretAlphabetBitLength
func ensureDivisible(value string, secretAlphabetBitLength int) string {
	result := value
	nb := secretAlphabetBitLength - (len(value) % secretAlphabetBitLength)
	if len(value)%secretAlphabetBitLength > 0 {
		for i := 0; i < nb; i++ {
			result += "0"
		}
	}

	return result
}

// indexOf returns the index of an element in a []string
func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}

	return -1
}

// zeropadding pads a string with zeros
func zeropadding(value string, length int) string {
	myString := ""
	valueLength := len(value)

	if valueLength >= length {
		return value
	}

	for i := 0; i < length-valueLength; i++ {
		myString += "0"
	}
	myString += value

	return myString
}
