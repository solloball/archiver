package table

import "strings"

type Generator interface {
	NewTable(text string) EncodingTable
}

type EncodingTable map[rune]string

func (et EncodingTable) Decode(str string) string {
	dt := et.decodingTree()

	return dt.Decode(str)
}

type decodingTree struct {
	Value string
	Zero  *decodingTree
	One   *decodingTree
}

func (et EncodingTable) decodingTree() decodingTree {
	res := decodingTree{}
	res.One = nil
	res.Zero = nil
	res.Value = ""

	for ch, code := range et {
		res.add(code, ch)
	}

	return res
}

func (dt *decodingTree) add(code string, value rune) {
	currentNode := dt

	for _, ch := range code {
		switch ch {
		case '0':
			if currentNode.Zero == nil {
				currentNode.Zero = &decodingTree{}
				currentNode.Zero.One = nil
				currentNode.Zero.Zero = nil
				currentNode.Zero.Value = ""
			}

			currentNode = currentNode.Zero
		case '1':
			if currentNode.One == nil {
				currentNode.One = &decodingTree{}
				currentNode.One.One = nil
				currentNode.One.Zero = nil
				currentNode.One.Value = ""
			}

			currentNode = currentNode.One
		}
	}

	currentNode.Value = string(value)
}

func (dt *decodingTree) Decode(str string) string {
	var buf strings.Builder

	currentNode := dt

	for _, ch := range str {
		if currentNode.Value != "" {
			buf.WriteString(currentNode.Value)
			currentNode = dt
		}

		switch ch {
		case '0':
			currentNode = currentNode.Zero
		case '1':
			currentNode = currentNode.One
		}
	}

	if currentNode.Value != "" {
		buf.WriteString(currentNode.Value)
		currentNode = dt
	}

	return buf.String()
}
