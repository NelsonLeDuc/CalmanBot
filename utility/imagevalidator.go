package utility

import (
	"encoding/hex"
	"io"
)

func ValidateImage(r io.Reader) bool {
	buf := make([]byte, 8)
	num, err := r.Read(buf)
	if err != nil || num < 8 {
		return false
	}

	gif := convertHexSlice([]string{"47", "49", "46"})
	jpg := convertHexSlice([]string{"FF", "D8", "FF", "E0"})
	png := convertHexSlice([]string{"89", "50", "4E", "47", "D", "A", "1A", "A"})

	return byteSliceSubset(buf, gif) || byteSliceSubset(buf, jpg) || byteSliceSubset(buf, png)
}

//b is a subset of a
func byteSliceSubset(a, b []byte) bool {
	for i, el := range b {
		if el != a[i] {
			return false
		}
	}
	
	return true
}

func convertHexSlice(s []string) []byte {
	b := []byte{}
	for _, hexStr := range s {
		result, err := hex.DecodeString(hexStr)
		if err == nil {
			b = append(b, result[0])
		}
	}

	return b
}