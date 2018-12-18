package tweg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrBinaryTooLong is raised when binary representation of character is too long
	ErrBinaryTooLong = errors.New("The binary representation of character is too long")

	// ErrParsingHexaToDecimal is raised when hexa parsing fails
	ErrParsingHexaToDecimal = errors.New("Error while parsing hexa to decimal")

	// ErrParsingBinaryToDecimal is raised when binary parsing fails
	ErrParsingBinaryToDecimal = errors.New("Error while parsing binary to decimal")

	// ErrInvalidCharacter is raised when an invalid character is used
	ErrInvalidCharacter = errors.New("Invalid character")
)

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// Tweg is a Twitter stegano
type Tweg struct {
	// SecretAlphabetString is the alphabet
	SecretAlphabetString string

	// SecretAlphabet is the exploded alphabet
	SecretAlphabet []string

	// SecretAlphabetBitLength is the length of the alphabet
	SecretAlphabetBitLength int

	// HomoglyphsLookup is the lookup table
	HomoglyphsLookup map[string]string
}

//------------------------------------------------------------------------------
// Factory
//------------------------------------------------------------------------------

// NewTweg returns a new Tweg
func NewTweg() *Tweg {
	t := &Tweg{
		SecretAlphabetString: " abcdefghijklmnopqrstuvwxyz123456789'0.:/\\%-_?&;",
		HomoglyphsLookup:     make(map[string]string),
	}
	t.SecretAlphabet = strings.Split(t.SecretAlphabetString, "")
	t.SecretAlphabetBitLength = len(strconv.FormatInt(int64(len(t.SecretAlphabet)), 2))

	// Lookup
	t.lookup()

	return t
}

//------------------------------------------------------------------------------
// Structure
//------------------------------------------------------------------------------

// Encode the tweet
func (t *Tweg) Encode(tweet, secret string) (string, error) {
	secret = strings.ToLower(secret) + " "
	secretBinary := ""
	tweetCovertextChars := 0
	result := ""

	// Process the secret
	for i := 0; i < len(secret); i++ {
		character := string(secret[i])
		secretAlphabetIndex := indexOf(character, t.SecretAlphabet)

		if secretAlphabetIndex >= 0 {
			secretCharacterBinary := zeropadding(strconv.FormatInt(int64(secretAlphabetIndex), 2), t.SecretAlphabetBitLength)
			if len(secretCharacterBinary) != t.SecretAlphabetBitLength {
				return "", ErrBinaryTooLong
			}
			secretBinary += secretCharacterBinary
		} else {
			return "", ErrInvalidCharacter
		}
	}

	// Ensure that the secret binary is divisible by alphabet bit length
	secretBinary = ensureDivisible(secretBinary, t.SecretAlphabetBitLength)

	// Process the tweet
	for i := 0; i < len(tweet); i++ {
		character := string(tweet[i])

		homoglyph, exists := Homoglyphs[character]
		if exists {
			homoglyphOptions := homoglyph
			homoglyphOptionsBitLength := len(strconv.FormatInt(int64(len(homoglyphOptions)+1), 2)) - 1
			tweetCovertextChars += homoglyphOptionsBitLength

			if len(secretBinary) > 0 {
				// Add the missing zeros if needed
				if len(secretBinary) < homoglyphOptionsBitLength {
					nb := homoglyphOptionsBitLength - len(secretBinary)
					for i := 0; i < nb; i++ {
						secretBinary += "0"
					}
				}
				secretBinaryToEncode := secretBinary[0:homoglyphOptionsBitLength]
				secretBinary = secretBinary[homoglyphOptionsBitLength:]
				secretBinaryToEncodeInDecimal, err := strconv.ParseInt(secretBinaryToEncode, 2, 64)
				if err != nil {
					return "", ErrParsingBinaryToDecimal
				}

				if secretBinaryToEncodeInDecimal > 0 {
					characterCodeInHexadecimal := homoglyphOptions[secretBinaryToEncodeInDecimal-1]
					characterCodeInDecimal, err := strconv.ParseInt(characterCodeInHexadecimal, 16, 64)
					if err != nil {
						return "", ErrParsingHexaToDecimal
					}
					character = string(characterCodeInDecimal)
				}
			}
		}

		result += character
	}

	return result, nil
}

// Decode the tweet
func (t *Tweg) Decode(tweet string) (string, error) {
	secretBinary := ""
	result := ""

	// Lookup all the characters of the tweet and forge the secret binary
	for _, character := range tweet {
		homoglyphLookup, exists := t.HomoglyphsLookup[string(character)]
		if exists {
			secretBinary += homoglyphLookup
		}
	}

	// Ensure that the secret binary is divisible by alphabet bit length
	secretBinary = ensureDivisible(secretBinary, t.SecretAlphabetBitLength)

	// Decode the secret binary
	for len(secretBinary) > 0 {
		secretCharacteInBinary := secretBinary[0:t.SecretAlphabetBitLength]
		if len(secretCharacteInBinary) > 0 {
			secretCharacteInBinary = zeropadding(secretCharacteInBinary, t.SecretAlphabetBitLength)
			if len(secretCharacteInBinary) != t.SecretAlphabetBitLength {
				return "", fmt.Errorf("ERROR: Unable to extract 5 characters (zeropadded) from string. ")
			}
			secretCharacterInDecimal, err := strconv.ParseInt(secretCharacteInBinary, 2, 64)
			if err != nil {
				return "", fmt.Errorf("Error while parsing the secretCharacteInBinary")
			}
			if secretCharacterInDecimal < int64(len(t.SecretAlphabet)) {
				result += t.SecretAlphabet[secretCharacterInDecimal]
			}
		}

		secretBinary = secretBinary[t.SecretAlphabetBitLength:]
	}

	return result, nil
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

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
