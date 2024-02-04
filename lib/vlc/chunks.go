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

// Join joins binary chunks into one string.
func (bcs BinaryChunks) Join() string {
	var buf strings.Builder

	for _, bc := range bcs {
		buf.WriteString(string(bc))
	}

	return buf.String()
}

func NewBinChunks(data []byte) BinaryChunks {
	res := make(BinaryChunks, 0, len(data))

	for _, part := range data {
		res = append(res, NewBinChunk(part))
	}

	return res
}

func NewBinChunk(code byte) BinaryChunk {
	return BinaryChunk(fmt.Sprintf("%08b", code))
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

func (bcs BinaryChunks) Bytes() []byte {
	res := make([]byte, 0, len(bcs))

	for _, bc := range bcs {
		res = append(res, bc.Byte())
	}

	return res
}

func (bc BinaryChunk) Byte() byte {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil {
		panic("can't parse Binary chunk into num" + err.Error())
	}

	return byte(num)
}
