package crawler

import (
	"strings"
)

const placeholder = '-'

func Stencil(stencil string, toStrip string) string {
	var result []byte
	for idx, char := range stencil {
		if char != placeholder && idx < len(toStrip) {
			result = append(result, toStrip[idx])
		}
	}
	return string(result)
}

// StripDashes strips '-' from the given string
func stripDashes(toStrip string) string {
	return strings.Replace(toStrip, string(placeholder), "", -1)
}
