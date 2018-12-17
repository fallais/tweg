package tweg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
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

	// Print some useful values
	logrus.WithFields(logrus.Fields{
		"SECRET ALPHABET BIT LENGTH": t.SecretAlphabetBitLength,
		"SECRET LENGTH":              len(secret),
		"TWEET":                      tweet,
		"SECRET":                     secret,
		"SECRET BINARY":              secretBinary,
		"SECRET BINARY LENGTH":       len(secretBinary),
		"SECRET ALPHABET STRING":     t.SecretAlphabetString,
	}).Debugln("Debug values")

	// Ensure
	atoi64, err := strconv.ParseInt(secretBinary, 2, 64)
	if err != nil {
		return "", fmt.Errorf("Error while converting : %s", err)
	}
	fmt.Println(atoi64)
	if int(atoi64)%t.SecretAlphabetBitLength > 0 {
		for i := 0; i < t.SecretAlphabetBitLength-(int(atoi64)%t.SecretAlphabetBitLength); i++ {
			secretBinary += "0"
		}
	}

	logrus.Debugln("SECRET BINARY AFTER ENSURE :", secretBinary)

	// Process the tweet
	for i := 0; i < len(tweet); i++ {
		character := string(tweet[i])

		homoglyph, exists := homoglyphs[character]
		if exists {
			homoglyphOptions := homoglyph
			homoglyphOptionsBitLength := len(strconv.FormatInt(int64(len(homoglyphOptions)+1), 2)) - 1
			tweetCovertextChars += homoglyphOptionsBitLength

			if len(secretBinary) > 0 {
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
		logrus.Debugln("PARTIAL RESULT :", result, character)
	}

	return result, nil
}

// Decode the tweet
func (t *Tweg) Decode(tweet string) string {
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
	nb := t.SecretAlphabetBitLength - (len(secretBinary) % t.SecretAlphabetBitLength)
	if len(secretBinary)%t.SecretAlphabetBitLength > 0 {
		for i := 0; i < nb; i++ {
			secretBinary += "0"
		}
	}

	// Decode the secret binary
	for len(secretBinary) > 0 {
		secretCharacteInBinary := secretBinary[0:t.SecretAlphabetBitLength]
		if len(secretCharacteInBinary) > 0 {
			secretCharacteInBinary = zeropadding(secretCharacteInBinary, t.SecretAlphabetBitLength)
			if len(secretCharacteInBinary) != t.SecretAlphabetBitLength {
				logrus.Errorln("ERROR: Unable to extract 5 characters (zeropadded) from string. ")
			}
			secretCharacterInDecimal, err := strconv.ParseInt(secretCharacteInBinary, 2, 64)
			if err != nil {
				fmt.Println("Error while parsing the secretCharacteInBinary")
			}
			if secretCharacterInDecimal < int64(len(t.SecretAlphabet)) {
				result += t.SecretAlphabet[secretCharacterInDecimal]
			}
		}

		secretBinary = secretBinary[t.SecretAlphabetBitLength:]
	}

	return result
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

// Lookup the homoglyphs
func (t *Tweg) lookup() {
	for c, h := range homoglyphs {
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

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}

	return -1
}

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

//------------------------------------------------------------------------------
// Const
//------------------------------------------------------------------------------

var homoglyphs = map[string][]string{
	"!":  []string{"FF01"},
	"\"": []string{"FF02"},
	"$":  []string{"FF04"},
	"%":  []string{"FF05"},
	"&":  []string{"FF06"},
	"'":  []string{"FF07"},
	"(":  []string{"FF08"},
	")":  []string{"FF09"},
	"*":  []string{"FF0A"},
	"+":  []string{"FF0B"},
	",":  []string{"FF0C"},
	"-":  []string{"FF0D"},
	".":  []string{"FF0E"},
	"/":  []string{"FF0F"},
	"0":  []string{"FF10"},
	"1":  []string{"FF11"},
	"2":  []string{"FF12"},
	"3":  []string{"FF13"},
	"4":  []string{"FF14"},
	"5":  []string{"FF15"},
	"6":  []string{"FF16"},
	"7":  []string{"FF17"},
	"8":  []string{"FF18"},
	"9":  []string{"FF19"},
	":":  []string{"FF1A"},
	";":  []string{"FF1B"},
	"<":  []string{"FF1C"},
	"=":  []string{"FF1D"},
	">":  []string{"FF1E"},
	"?":  []string{"FF1F"},
	"@":  []string{"FF20"},
	"A":  []string{"FF21", "0391", "0410"},
	"B":  []string{"FF22", "0392", "0412"},
	"C":  []string{"FF23", "03F9", "216D"},
	"D":  []string{"FF24"},
	"E":  []string{"FF25", "0395", "0415"},
	"F":  []string{"FF26"},
	"G":  []string{"FF27"},
	"H":  []string{"FF28", "0397", "041D"},
	"I":  []string{"FF29", "0399", "0406"},
	"J":  []string{"FF2A"},
	"K":  []string{"FF2B", "039A", "212A"},
	"L":  []string{"FF2C"},
	"M":  []string{"FF2D", "039C", "041C"},
	"N":  []string{"FF2E"},
	"O":  []string{"FF2F", "039F", "041E"},
	"P":  []string{"FF30", "03A1", "0420"},
	"Q":  []string{"FF31"},
	"R":  []string{"FF32"},
	"S":  []string{"FF33"},
	"T":  []string{"FF34", "03A4", "0422"},
	"U":  []string{"FF35"},
	"V":  []string{"FF36", "0474", "2164"},
	"W":  []string{"FF37"},
	"X":  []string{"FF38", "03A7", "2169"},
	"Y":  []string{"FF39", "03A5", "04AE"},
	"Z":  []string{"FF3A"},
	"[":  []string{"FF3B"},
	"\\": []string{"FF3C"},
	"]":  []string{"FF3D"},
	"^":  []string{"FF3E"},
	"_":  []string{"FF3F"},
	"`":  []string{"FF40"},
	"a":  []string{"FF41"},
	"b":  []string{"FF42"},
	"c":  []string{"FF43", "03F2", "0441"},
	"d":  []string{"FF44"},
	"e":  []string{"FF45"},
	"f":  []string{"FF46"},
	"g":  []string{"FF47"},
	"h":  []string{"FF48"},
	"i":  []string{"FF49", "0456", "2170"},
	"j":  []string{"FF4A"},
	"k":  []string{"FF4B"},
	"l":  []string{"FF4C"},
	"m":  []string{"FF4D"},
	"n":  []string{"FF4E"},
	"o":  []string{"FF4F", "03BF", "043E"},
	"p":  []string{"FF50"},
	"q":  []string{"FF51"},
	"r":  []string{"FF52"},
	"s":  []string{"FF53"},
	"t":  []string{"FF54"},
	"u":  []string{"FF55"},
	"v":  []string{"FF56", "03BD", "2174"},
	"w":  []string{"FF57"},
	"x":  []string{"FF58", "0445", "2179"},
	"y":  []string{"FF59"},
	"z":  []string{"FF5A"},
	"{":  []string{"FF5B"},
	"|":  []string{"FF5C"},
	"}":  []string{"FF5D"},
	"~":  []string{"FF5E"},
	" ":  []string{"2000", "2001", "2002", "2003", "2004", "2005", "2006", "2007", "2008", "2009", "200A", "2028", "2029", "202F", "205F"},
}
