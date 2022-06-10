package radix

import (
	"fmt"
	"unicode/utf8"
)

// longestCommonPrefix finds the longest common prefix.
// This also implies that the common prefix contains no ':' or '*'
// since the existing key can't contain those chars.
func longestCommonPrefix(a, b string) int {
	leftBracket := rune(123)

	i := 0
	size := 0
	max := min(utf8.RuneCountInString(a), utf8.RuneCountInString(b))
	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)
		rb, sizeB := utf8.DecodeRuneInString(b)

		if ra == leftBracket || rb == leftBracket {
			return size
		}

		a = a[sizeA:]
		b = b[sizeB:]

		if ra != rb {
			return i
		}

		i++
		size += sizeA
	}

	return size
}

func findParamEnd(a string) int {
	leftBracket := rune(123)
	rightBracket := rune(125)
	slash := rune(47)

	i := 0
	size := 0
	max := utf8.RuneCountInString(a)
	var prevA rune
	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)

		if ra == slash {
			return -1
		}
		if prevA == leftBracket && ra == rightBracket {
			return -1
		}

		if ra == rightBracket {
			return size + sizeA
		}

		a = a[sizeA:]

		i++
		size += sizeA
		prevA = ra
	}

	return -1
}

func findParamStart(a string) int {
	leftBracket := rune(123)

	i := 0
	size := 0
	max := utf8.RuneCountInString(a)
	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)

		if ra == leftBracket {
			return size
		}

		a = a[sizeA:]

		i++
		size += sizeA
	}

	return -1
}

func findSlashOrEnd(a string) int {
	slash := rune(47)

	i := 0
	size := 0
	max := utf8.RuneCountInString(a)
	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)

		if ra == slash {
			return size
		}

		a = a[sizeA:]

		i++
		size += sizeA
	}

	return size
}

func updateParamNode(pn Node, path string, key uint64) Node {
	end := findParamEnd(path)
	if end == -1 {
		panic(fmt.Sprintf("no right bracket: %v", path))
	}

	if pn.kind != param {
		panic("node must be kind param")
	}

	if pn.path != "" && pn.path != path[:end] {
		panic(ErrParamNameConflict)
	}

	if pn.path == "" {
		pn.path = path[:end]
	}

	path = path[end:]
	if path != "" {
		pn = pn.Insert(path, key)
	} else {
		if pn.key > 0 && pn.key != key {
			panic(ErrPathAlreadyTaken)
		}

		pn.key = key
	}

	return pn
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
