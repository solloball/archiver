package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const chunkSize = 8

type encodingTable map[rune]string

type BynaryChunk string

type BynaryChunks []BynaryChunk

type HexChunk string

type HexChunks []HexChunk

func Encode(str string) string {
	str = prepareText(str)

	chunks := splitByChunks(encodeBinary(str), chunkSize)

	return chunks.ToHex().ToString()
}

// prepare text for ecnoding.
// changes upper character to ! +  lower character.
// i.g.: My number -> !my number.
func prepareText(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(ch))
		} else {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}

// encodes str inti binary codes without whitespace.
func encodeBinary(str string) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch))
	}

	return buf.String()
}

func bin(ch rune) string {
	table := getEncodingTable()

	res, ok := table[ch]

	if !ok {
		panic("unknown charachter: " + string(ch))
	}

	return res
}

func getEncodingTable() encodingTable {
	return encodingTable{
		' ': "11",
		't': "1001",
		'n': "10000",
		's': "0101",
		'r': "01000",
		'd': "00101",
		'!': "001000",
		'c': "000101",
		'm': "000011",
		'g': "0000100",
		'b': "0000010",
		'v': "00000001",
		'k': "0000000001",
		'q': "000000000001",
		'e': "101",
		'o': "10001",
		'a': "011",
		'i': "01001",
		'h': "0011",
		'l': "001001",
		'u': "00011",
		'f': "000100",
		'p': "0000101",
		'w': "0000011",
		'y': "000000",
		'j': "000000001",
		'x': "00000000001",
		'z': "000000000000",
	}
}

// splitByChunks splits string into chunks into size.
// i.g.: "0000000011111111" -> "00000000" "11111111", size = 8.
func splitByChunks(str string, chunckSize int) BynaryChunks {

	strLen := utf8.RuneCountInString(str)

	chuncksCount := strLen / chunckSize
	if strLen/chunckSize != 0 {
		chuncksCount++
	}

	res := make(BynaryChunks, 0, chuncksCount)

	var buf strings.Builder

	for i, ch := range str {
		buf.WriteString(string(ch))

		if (i+1)%chunckSize == 0 {
			res = append(res, BynaryChunk(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()

		lastChunk += strings.Repeat("0", chunckSize-len(lastChunk))
		res = append(res, BynaryChunk(lastChunk))
	}

	return res
}

func (bc BynaryChunk) ToHex() HexChunk {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil {
		panic("can't parse binary chunk: " + err.Error())
	}

	res := strings.ToUpper(fmt.Sprintf("%x", num))

	if len(res) == 1 {
		res = "0" + res
	}

	return HexChunk(res)
}

func (bcs BynaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		res = append(res, chunk.ToHex())
	}

	return res
}

func (hcs HexChunks) ToString() string {
	const sep = " "

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}

	var buf strings.Builder

	buf.WriteString(string(hcs[0]))
	for _, hc := range hcs[1:] {
		buf.WriteString(sep)
		buf.WriteString(string(hc))
	}

	return buf.String()
}
