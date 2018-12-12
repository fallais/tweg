package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	tweet  = flag.String("tweet", "A koala arrives in the great forest of Wumpalumpa", "Tweet ?")
	secret = flag.String("secret", "alpaga", "Secret ?")
)

func init() {
	// Parse the flags
	flag.Parse()

	// Set localtime to UTC
	time.Local = time.UTC
}

func main() {
	tweet := *tweet
	secret := strings.ToLower(*secret)
	secretBinary := ""
	secretAlphabetString := " abcdefghijklmnopqrstuvwxyz123456789'0.:/\\%-_?&;"
	secretAlphabet := strings.Split(secretAlphabetString, "")
	secretAlphabetBitLength := len(strconv.FormatInt(int64(len(secretAlphabet)), 2))
	tweetCovertextChars := 0
	result := ""

	// Process the secret
	for i := 0; i < len(secret); i++ {
		character := string(secret[i])
		secretAlphabetIndex := indexOf(character, secretAlphabet)

		if secretAlphabetIndex >= 0 {
			secretCharacterBinary := zeropad(strconv.FormatInt(int64(secretAlphabetIndex), 2), secretAlphabetBitLength)
			if len(secretCharacterBinary) != secretAlphabetBitLength {
				fmt.Errorf("ERROR: binary representation of character too big")
			}
			secretBinary += secretCharacterBinary
		} else {
			fmt.Errorf("ERROR: secret contains invalid character '" + character + "' Ignored")
		}
	}

	fmt.Println("TWEET :", tweet)
	fmt.Println("SECRET :", secret)
	fmt.Println("SECRET BINARY :", secretBinary)
	fmt.Println("SECRET ALPHABET STRING :", secretAlphabetString)

	// Ensure
	if len(secretBinary)%secretAlphabetBitLength > 0 {
		secretBinary = zeropad("0", secretAlphabetBitLength-(len(secretBinary)%secretAlphabetBitLength))
	}

	// Process the tweet
	for i := 0; i < len(tweet); i++ {
		character := tweet[i]
		var character2 []byte

		homoglyph, exists := homoglyphs[string(character)]
		if exists {
			homoglyphOptions := homoglyph
			homoglyphOptionsBitLength := len(strconv.FormatInt(int64(len(homoglyphOptions)+1), 2)) - 1
			tweetCovertextChars += homoglyphOptionsBitLength

			if len(secretBinary) > 0 {
				fmt.Println(len(secretBinary), homoglyphOptionsBitLength)
				secretBinaryToEncode := secretBinary[0:homoglyphOptionsBitLength]
				secretBinary = secretBinary[:homoglyphOptionsBitLength]
				secretBinaryToEncodeInDecimal, _ := strconv.ParseInt(secretBinaryToEncode, 2, 64)

				if secretBinaryToEncodeInDecimal > 0 {
					characterCodeInHexadecimal := homoglyphOptions[secretBinaryToEncodeInDecimal-1]
					characterCodeInDecimal, _ := strconv.ParseInt(characterCodeInHexadecimal, 16, 64)
					//character = String.fromCharCode(character_code_in_decimal)
					character2 = []byte(strconv.Itoa(int(characterCodeInDecimal)))
				}
			}
		}
		result += string(character2)
		fmt.Println("PARTIAL RESULT :", result)
	}

	fmt.Println("RESULT :", result)
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func zeropad(value string, length int) string {
	return strPad(value, length, "0", "LEFT")
}

func strPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

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
