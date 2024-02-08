package vlc

import (
	"strings"
)

type DecodingTree struct {
	Value string
	Zero  *DecodingTree
	One   *DecodingTree
}

func (et encodingTable) DecodingTree() DecodingTree {
	res := DecodingTree{}
	res.One = nil
	res.Zero = nil
	res.Value = ""

	for ch, code := range et {
		res.Add(code, ch)
	}

	return res
}

func (dt *DecodingTree) Add(code string, value rune) {
	currentNode := dt

	for _, ch := range code {
		switch ch {
		case '0':
			if currentNode.Zero == nil {
				currentNode.Zero = &DecodingTree{}
				currentNode.Zero.One = nil
				currentNode.Zero.Zero = nil
				currentNode.Zero.Value = ""
			}

			currentNode = currentNode.Zero
		case '1':
			if currentNode.One == nil {
				currentNode.One = &DecodingTree{}
				currentNode.One.One = nil
				currentNode.One.Zero = nil
				currentNode.One.Value = ""
			}

			currentNode = currentNode.One
		}
	}

	currentNode.Value = string(value)
}

func (dt *DecodingTree) Decode(str string) string {
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
