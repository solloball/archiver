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
			}

			currentNode = currentNode.Zero
		case '1':
			if currentNode.One == nil {
				currentNode.One = &DecodingTree{}
			}

			currentNode = currentNode.One
		}

		currentNode.Value = string(value)
	}
}

func (dt *DecodingTree) Decode(str string) string {
	var buf strings.Builder

	currentNode := dt

	for _, ch := range str {
		if dt.Value != "" {
			buf.WriteString(dt.Value)
			currentNode = dt
			continue
		}

		switch ch {
		case '0':
			currentNode = currentNode.Zero
		case '1':
			currentNode = currentNode.One
		}
	}

	if dt.Value != "" {
		buf.WriteString(dt.Value)
		currentNode = dt
	}

	return buf.String()
}
