package vlc

import (
	"archiver/lib/comperssion/vlc/table"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"strings"
)

type EnocoderDecoder struct {
	tblGenerator table.Generator
}

func New(tblGenerator table.Generator) EnocoderDecoder {
	return EnocoderDecoder{tblGenerator: tblGenerator}
}

func (ed EnocoderDecoder) Encode(str string) []byte {

	tbl := ed.tblGenerator.NewTable(str)

	encoded := encodeBinary(str, tbl)

	return buildEncodedFile(tbl, encoded)
}

func (ed EnocoderDecoder) Decode(encodedData []byte) string {
	tbl, data := parseFile(encodedData)

	return tbl.Decode(data)
}

func parseFile(data []byte) (table.EncodingTable, string) {
	const (
		tableSizeBytesCount = 4
		dataSizeBytesCount  = 4
	)

	tableSizeBinary, data := data[:tableSizeBytesCount],
		data[tableSizeBytesCount:]

	dataSizeBinary, data := data[:dataSizeBytesCount],
		data[dataSizeBytesCount:]

	tableSize := binary.BigEndian.Uint32(tableSizeBinary)
	dataSize := binary.BigEndian.Uint32(dataSizeBinary)

	tblBinary, data := data[:tableSize], data[tableSize:]

	tbl := decodeTable(tblBinary)

	body := NewBinChunks(data).Join()

	return tbl, body[:dataSize]
}

func buildEncodedFile(tbl table.EncodingTable, data string) []byte {
	encodedTable := encodeTable(tbl)

	var buf bytes.Buffer

	buf.Write(encodeInt(len(encodedTable)))
	buf.Write(encodeInt(len(data)))
	buf.Write(encodedTable)
	buf.Write(splitByChunks(data, chunkSize).Bytes())

	return buf.Bytes()
}

func encodeInt(num int) []byte {
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, uint32(num))

	return res
}

func decodeTable(tblBinary []byte) table.EncodingTable {
	var tbl table.EncodingTable

	r := bytes.NewReader(tblBinary)
	if err := gob.NewDecoder(r).Decode(&tbl); err != nil {
		log.Fatal("can't decode table", err)
	}

	return tbl
}

func encodeTable(tbl table.EncodingTable) []byte {
	var tableBuf bytes.Buffer

	if err := gob.NewEncoder(&tableBuf).Encode(tbl); err != nil {
		log.Fatal("can't serialize table", err)
	}

	return tableBuf.Bytes()
}

// encodes str inti binary codes without whitespace.
func encodeBinary(str string, table table.EncodingTable) string {
	var buf strings.Builder

	for _, ch := range str {
		buf.WriteString(bin(ch, table))
	}

	return buf.String()
}

func bin(ch rune, table table.EncodingTable) string {
	res, ok := table[ch]

	if !ok {
		panic("unknown charachter: " + string(ch))
	}

	return res
}
