# TWEG

A **Twitter Steganography** built with **Golang** and based on the great work of [Steg Of The Dump](https://github.com/holloway/steg-of-the-dump)

## How it works

Given this :

- Tweet : `A koala arrives in the great forest of Wumpalumpa`
- Hidden message : `alpaga`
- Alphabet : ` abcdefghijklmnopqrstuvwxyz123456789'0.:/\\%-_?&`

### Step 1 : generate binary representation of the hidden message

First letter is `a`, its position in the alphabet is `1`. So its binary reprsentation is `000001`.  
Second letter is `l`, its position in the alphabet is `12`. So its binary reprsentation is `001100`.  
Third letter is `p`, its position in the alphabet is `16`. So its binary reprsentation is `010000`.  
Forth letter is `a`, its position in the alphabet is `1`. So its binary reprsentation is `000001`.  
Fifth letter is `g`, its position in the alphabet is `7`. So its binary reprsentation is `000111`.  
Sixth letter is `a`, its position in the alphabet is `1`. So its binary reprsentation is `000001`.  
A space is added (I do not know why..), its position in the alphabet is `0`. So its binary reprsentation is `000000`.  
Binary of hidden message is `000001001100010000000001000111000001000000`.

Then, it needs to ensure that the binary is divisible by the alphabet bitlength, which is `6`. To do so, it adds as many `0` as needed at the end.

Binary of hidden message is finally `0000010011000100000000010001110000010000000000`.

### Step 2 : hide the message in the tweet

#### First character

First character of the tweet is `A`, it has `3` homoglyphs, that means `4` possibilities. The binary of `4` is `100`, which has a length of `3`. It is substracting 1 so we get `2`. That means we can use `2 bit` to represent the `4 possibilities` of the character `A`.

So, we can get the first `2 bit` of the binary representation of the hidden message. Which are `00`. This is `0` in decimal. That means that this character is left intact.

Final tweet starts with `A`. The first 2 bit of the binary representation of the hidden message are removed.

#### Second character

Second character of the tweet is ` ` (space), it has `15` homoglyphs, that means `16` possibilities. The binary of `16` is `10000`, which has a length of `5`. It is substracting 1 soto st we get `4`. That means we can use `4 bit` to represent the `16 possibilities` of the character ` `.

So, we can get the next `4 bit` of the binary representation of the hidden message. Which are `0001`. This is `1` in decimal. It is subtracting `1` and searching the hexadecimal code of the homoglyph at this position, the hexadecimal code of the homoglyph of ` ` at position `2` is `2000`. So the character is `U+2000`, which is a homoplygh of a space.

Final tweet is now `A `. The first 4 bit of the binary representation of the hidden message are removed.

#### Again and again..

Repeat the step.