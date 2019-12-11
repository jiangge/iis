package storage

import (
	"bytes"
	"strconv"
)

type Bitmap []byte

func NewBitmap(n int) Bitmap {
	x := n / 8
	if x*8 != n {
		x++
	}
	return make(Bitmap, x)
}

func (m Bitmap) String() string {
	p := bytes.Buffer{}
	for _, v := range m {
		for i := 0; i < 8; i++ {
			s := strconv.Itoa(int((v >> i) & 1))
			p.WriteString(s)
		}
	}
	return p.String()
}

var firstZero = [256]int{
	0, // 00000000
	1, // 00000001
	0, // 00000010
	2, // 00000011
	0, // 00000100
	1, // 00000101
	0, // 00000110
	3, // 00000111
	0, // 00001000
	1, // 00001001
	0, // 00001010
	2, // 00001011
	0, // 00001100
	1, // 00001101
	0, // 00001110
	4, // 00001111
	0, // 00010000
	1, // 00010001
	0, // 00010010
	2, // 00010011
	0, // 00010100
	1, // 00010101
	0, // 00010110
	3, // 00010111
	0, // 00011000
	1, // 00011001
	0, // 00011010
	2, // 00011011
	0, // 00011100
	1, // 00011101
	0, // 00011110
	5, // 00011111
	0, // 00100000
	1, // 00100001
	0, // 00100010
	2, // 00100011
	0, // 00100100
	1, // 00100101
	0, // 00100110
	3, // 00100111
	0, // 00101000
	1, // 00101001
	0, // 00101010
	2, // 00101011
	0, // 00101100
	1, // 00101101
	0, // 00101110
	4, // 00101111
	0, // 00110000
	1, // 00110001
	0, // 00110010
	2, // 00110011
	0, // 00110100
	1, // 00110101
	0, // 00110110
	3, // 00110111
	0, // 00111000
	1, // 00111001
	0, // 00111010
	2, // 00111011
	0, // 00111100
	1, // 00111101
	0, // 00111110
	6, // 00111111
	0, // 01000000
	1, // 01000001
	0, // 01000010
	2, // 01000011
	0, // 01000100
	1, // 01000101
	0, // 01000110
	3, // 01000111
	0, // 01001000
	1, // 01001001
	0, // 01001010
	2, // 01001011
	0, // 01001100
	1, // 01001101
	0, // 01001110
	4, // 01001111
	0, // 01010000
	1, // 01010001
	0, // 01010010
	2, // 01010011
	0, // 01010100
	1, // 01010101
	0, // 01010110
	3, // 01010111
	0, // 01011000
	1, // 01011001
	0, // 01011010
	2, // 01011011
	0, // 01011100
	1, // 01011101
	0, // 01011110
	5, // 01011111
	0, // 01100000
	1, // 01100001
	0, // 01100010
	2, // 01100011
	0, // 01100100
	1, // 01100101
	0, // 01100110
	3, // 01100111
	0, // 01101000
	1, // 01101001
	0, // 01101010
	2, // 01101011
	0, // 01101100
	1, // 01101101
	0, // 01101110
	4, // 01101111
	0, // 01110000
	1, // 01110001
	0, // 01110010
	2, // 01110011
	0, // 01110100
	1, // 01110101
	0, // 01110110
	3, // 01110111
	0, // 01111000
	1, // 01111001
	0, // 01111010
	2, // 01111011
	0, // 01111100
	1, // 01111101
	0, // 01111110
	7, // 01111111
	0, // 10000000
	1, // 10000001
	0, // 10000010
	2, // 10000011
	0, // 10000100
	1, // 10000101
	0, // 10000110
	3, // 10000111
	0, // 10001000
	1, // 10001001
	0, // 10001010
	2, // 10001011
	0, // 10001100
	1, // 10001101
	0, // 10001110
	4, // 10001111
	0, // 10010000
	1, // 10010001
	0, // 10010010
	2, // 10010011
	0, // 10010100
	1, // 10010101
	0, // 10010110
	3, // 10010111
	0, // 10011000
	1, // 10011001
	0, // 10011010
	2, // 10011011
	0, // 10011100
	1, // 10011101
	0, // 10011110
	5, // 10011111
	0, // 10100000
	1, // 10100001
	0, // 10100010
	2, // 10100011
	0, // 10100100
	1, // 10100101
	0, // 10100110
	3, // 10100111
	0, // 10101000
	1, // 10101001
	0, // 10101010
	2, // 10101011
	0, // 10101100
	1, // 10101101
	0, // 10101110
	4, // 10101111
	0, // 10110000
	1, // 10110001
	0, // 10110010
	2, // 10110011
	0, // 10110100
	1, // 10110101
	0, // 10110110
	3, // 10110111
	0, // 10111000
	1, // 10111001
	0, // 10111010
	2, // 10111011
	0, // 10111100
	1, // 10111101
	0, // 10111110
	6, // 10111111
	0, // 11000000
	1, // 11000001
	0, // 11000010
	2, // 11000011
	0, // 11000100
	1, // 11000101
	0, // 11000110
	3, // 11000111
	0, // 11001000
	1, // 11001001
	0, // 11001010
	2, // 11001011
	0, // 11001100
	1, // 11001101
	0, // 11001110
	4, // 11001111
	0, // 11010000
	1, // 11010001
	0, // 11010010
	2, // 11010011
	0, // 11010100
	1, // 11010101
	0, // 11010110
	3, // 11010111
	0, // 11011000
	1, // 11011001
	0, // 11011010
	2, // 11011011
	0, // 11011100
	1, // 11011101
	0, // 11011110
	5, // 11011111
	0, // 11100000
	1, // 11100001
	0, // 11100010
	2, // 11100011
	0, // 11100100
	1, // 11100101
	0, // 11100110
	3, // 11100111
	0, // 11101000
	1, // 11101001
	0, // 11101010
	2, // 11101011
	0, // 11101100
	1, // 11101101
	0, // 11101110
	4, // 11101111
	0, // 11110000
	1, // 11110001
	0, // 11110010
	2, // 11110011
	0, // 11110100
	1, // 11110101
	0, // 11110110
	3, // 11110111
	0, // 11111000
	1, // 11111001
	0, // 11111010
	2, // 11111011
	0, // 11111100
	1, // 11111101
	0, // 11111110
	8, // 11111111
}

func (m Bitmap) FirstUnset(start int) int {
	for i := start % 8; i < 8; i++ {
		if (m[start/8]>>i)&1 == 0 {
			return start + i
		}
	}
	for i := start/8 + 1; i < len(m); i++ {
		if m[i] == 255 {
			continue
		}
		return i*8 + firstZero[m[i]]
	}
	return -1
}

func (m *Bitmap) Set(start, end int) {
	if end <= start {
		return
	}
	if end == start+1 {
		(*m)[start/8] |= 1 << (start % 8)
		return
	}
	if end/8 == start/8 {
		for i := start % 8; i < end%8; i++ {
			(*m)[start/8] |= 1 << i
		}
		return
	}
	for i := start % 8; i < 8; i++ {
		(*m)[start/8] |= 1 << i
	}
	for i := start/8 + 1; i < end/8; i++ {
		(*m)[i] = 0xFF
	}
	for i := 0; i < end%8; i++ {
		(*m)[end/8] |= 1 << i
	}
}

func (m *Bitmap) Unset(start, end int) {
	if end <= start {
		return
	}
	if end == start+1 {
		(*m)[start/8] &= ^(1 << (start % 8))
		return
	}
	if end/8 == start/8 {
		for i := start % 8; i < end%8; i++ {
			(*m)[start/8] &= ^(1 << i)
		}
		return
	}
	for i := start % 8; i < 8; i++ {
		(*m)[start/8] &= ^(1 << i)
	}
	for i := start/8 + 1; i < end/8; i++ {
		(*m)[i] = 0
	}
	for i := 0; i < end%8; i++ {
		(*m)[end/8] &= ^(1 << i)
	}
}