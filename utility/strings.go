package utility

import "unicode/utf8"

func DivideString(s string, n int) []string {
	return DivideStringWith(s, n, '…')
}

func DivideStringWith(s string, n int, r rune) []string {
	runes := []rune(s)
	split := []string{}

	idx := 0
	for idx < len(runes) {
		count := len(runes) - idx
		if count > n {
			count = n
		}

		slice := runes[idx : idx+count]
		sliceCopy := make([]rune, count)
		copy(sliceCopy, slice)

		if count == n && idx+count < len(runes) {
			runeSize := utf8.RuneLen(r)
			sliceCopy = sliceCopy[0 : len(sliceCopy)-runeSize]
			sliceCopy = append(sliceCopy, r)
			count -= runeSize
		}

		split = append(split, string(sliceCopy))
		idx += count
	}

	return split
}
