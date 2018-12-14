# TWEG

A **Twitter Steganography** based on the great work of [Steg Of The Dump](https://github.com/holloway/steg-of-the-dump)

## How it works

Given this :

- Tweet : `A koala arrives in the great forest of Wumpalumpa`
- Hidden message : `alpaga`
- Alphabet : ` abcdefghijklmnopqrstuvwxyz123456789'0.:/\\%-_?&`

### Step 1 : binary of the hidden message

First letter is `a`, its position in the alphabet is `1`. So its binary reprsentation is `000001`.

Second letter is `l`, its position in the alphabet is `12`. So its binary reprsentation is `001100`.

Third letter is `p`, its position in the alphabet is `16`. So its binary reprsentation is `010000`.

Forth letter is `a`, its position in the alphabet is `1`. So its binary reprsentation is `000001`.

Fifth letter is `g`, its position in the alphabet is `7`. So its binary reprsentation is `000111`.

Sixth letter is `a`, its position in the alphabet is `1`. So its binary reprsentation is `000001`.

A space is added, its position in the alphabet is `0`. So its binary reprsentation is `000000`.

Binary of hidden message is `000001001100010000000001000111000001000000`.

Then, it needs to ensure that the binary is divisible by the alphabet bitlength, which is `6`. To do so, it adds as many `0` as needed at the end.

Binary of hidden message is finally `0000010011000100000000010001110000010000000000`.

### Step 2

