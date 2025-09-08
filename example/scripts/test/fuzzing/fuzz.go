package fuzz

import (
	"testing"
	"unicode/utf8"
)

// Reverse returns the reverse of the input string
func Reverse(s string) string {
	b := []byte(s)
	result := make([]byte, len(b))

	for i, width := 0, 0; i < len(b); i += width {
		_, width = utf8.DecodeRune(b[i:])
		for j := 0; j < width; j++ {
			result[len(b)-i-width+j] = b[i+j]
		}
	}
	return string(result)
}

// FuzzReverse is a fuzz test for the Reverse function
func FuzzReverse(f *testing.F) {
	// Provide seed corpus
	testcases := []string{"Hello, world", "", "!12345"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	// Fuzz test
	f.Fuzz(func(t *testing.T, orig string) {
		rev := Reverse(orig)
		doubleRev := Reverse(rev)

		if orig != doubleRev {
			t.Errorf("Reverse(Reverse(%q)) = %q, want %q", orig, doubleRev, orig)
		}

		// Additional property checks can go here
	})
}
