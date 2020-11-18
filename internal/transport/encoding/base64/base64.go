// Package base64 provides support for Demon's Souls base64 encoding. In many
// cases the games sends data with broken base64 encoded data.
package base64

import (
	"encoding/base64"
	"strings"
)

// stdAlpha represents the set of valid base64 standard encoding characters.
const stdAlpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// StdEncoding is the standard base64 encoding, as defined in RFC 4648.
var StdEncoding = NewEncoding(stdAlpha)

// Encoding is a Demon's Souls base64 encoding handler.
type Encoding struct {
	characters map[rune]struct{}
}

// NewEncoding returns a new base64 decoder.
func NewEncoding(alpha string) *Encoding {
	c := make(map[rune]struct{}, len(alpha))
	for _, a := range alpha {
		c[a] = struct{}{}
	}

	return &Encoding{characters: c}
}

// DecodeString returns the bytes represented by the base64 string s.
//
// It wraps the standard library DecodeString method while fixing the broken
// characters in the source string.
func (e *Encoding) DecodeString(s string) ([]byte, error) {
	var sb strings.Builder

	// Filter out invalid characters and fix whitespace.
	for _, c := range s {
		if _, ok := e.characters[c]; ok {
			sb.WriteRune(c)
		} else if string(c) == " " {
			sb.WriteString("+")
		} else {
			break
		}
	}

	// Fix the ending.
	switch sb.Len() % 4 {
	case 3:
		sb.WriteString("=")
	case 2:
		sb.WriteString("==")
	case 1:
		sb.WriteString("A==")
	}

	return base64.StdEncoding.DecodeString(sb.String())
}

// EncodeToString returns the base64 encoding of src.
//
// It wraps the standard library EncodeToString method while altering the
// encoded values to the broken versions that Demon's Souls expects.
func (e *Encoding) EncodeToString(src []byte) string {
	s := base64.StdEncoding.EncodeToString(src)
	return strings.Replace(s, "+", " ", -1)
}
