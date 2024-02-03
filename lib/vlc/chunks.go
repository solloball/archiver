package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const chunkSize = 8

type encodingTable map[rune]string

type BinaryChunk string

type BinaryChunks []BinaryChunk

type HexChunk string

type HexChunks []HexChunk

const hexChunksSeparator = " "

// Join joins binary chunks into one string.
func (bcs BinaryChunks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}

	return buf.String()
}

func (hcs HexChunks) ToBin() BinaryChunks {
	res := make(BinaryChunks, 0, len(hcs))

	for _, hc := range hcs {
		res = append(res, hc.ToBin())
	}

	return res
}

func (hc HexChunk) ToBin() BinaryChunk {
	num, err := strconv.ParseUint(string(hc), 16, chunkSize)
	if err != nil {
		panic("can't parse hex chunk into number: " + err.Error())
	}

	return BinaryChunk(fmt.Sprintf("%08b", num))
}

func NewHexChunks(str string) HexChunks {
	parts := strings.Split(str, hexChunksSeparator)

	res := make(HexChunks, 0, len(parts))

	for _, part := range parts {
		res = append(res, HexChunk(part))
	}

	return res
}

// splitByChunks splits string into chunks into size.
// i.g.: "0000000011111111" -> "00000000" "11111111", size = 8.
func splitByChunks(str string, chunckSize int) BinaryChunks {

	strLen := utf8.RuneCountInString(str)

	chuncksCount := strLen / chunckSize
	if strLen/chunckSize != 0 {
		chuncksCount++
	}

	res := make(BinaryChunks, 0, chuncksCount)

	var buf strings.Builder

	for i, ch := range str {
		buf.WriteString(string(ch))

		if (i+1)%chunckSize == 0 {
			res = append(res, BinaryChunk(buf.String()))
			buf.Reset()
		}
	}

	if buf.Len() != 0 {
		lastChunk := buf.String()

		lastChunk += strings.Repeat("0", chunckSize-len(lastChunk))
		res = append(res, BinaryChunk(lastChunk))
	}

	return res
}

func (bc BinaryChunk) ToHex() HexChunk {
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

func (bcs BinaryChunks) ToHex() HexChunks {
	res := make(HexChunks, 0, len(bcs))

	for _, chunk := range bcs {
		res = append(res, chunk.ToHex())
	}

	return res
}

func (hcs HexChunks) ToString() string {

	switch len(hcs) {
	case 0:
		return ""
	case 1:
		return string(hcs[0])
	}

	var buf strings.Builder

	buf.WriteString(string(hcs[0]))
	for _, hc := range hcs[1:] {
		buf.WriteString(hexChunksSeparator)
		buf.WriteString(string(hc))
	}

	return buf.String()
}
