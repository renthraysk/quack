package quack

import (
	"time"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/huffman"
)

type Encoder struct {
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-request-pseudo-header-field
func (e *Encoder) NewRequest(p []byte, method, scheme, authority, path string) []byte {
	p = append(p, 0, 0)
	p = e.appendMethod(p, method)
	p = e.appendScheme(p, scheme)
	if authority != "" {
		p = e.appendAuthority(p, authority)
	}
	// This (:path) pseudo-header field MUST NOT be empty for "http" or "https"
	// URIs; "http" or "https" URIs that do not contain a path component MUST
	// include a value of / (ASCII 0x2f).
	if path == "" && (scheme == "http" || scheme == "https") {
		path = "/"
		// An OPTIONS request that does not include a path component includes
		// the value *
		if method == "OPTIONS" {
			path = "*"
		}
	}
	p = e.appendPath(p, path)
	return p
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-the-connect-method
func (e *Encoder) NewConnect(p []byte, authority string) []byte {
	p = append(p, 0, 0)
	p = e.appendMethod(p, "CONNECT")
	p = e.appendAuthority(p, authority)
	return p
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-response-pseudo-header-fiel
func (e *Encoder) NewResponse(p []byte, statusCode int) []byte {
	p = append(p, 0, 0)
	p = e.appendStatus(p, statusCode)
	return p
}

func (e *Encoder) appendHeaderField(p []byte, name, value string, neverIndex bool) []byte {
	// @TODO Dynamic stuff goes here.

	// appendLiteralFieldWithoutNameReference
	p = appendLiteralName(p, name, neverIndex)
	return appendStringLiteral(p, value)
}

// appendLiteralName appends the QPACK encoding of the lower case of name to p,
// applying huffman encoding if it would result in savings.
// https://datatracker.ietf.org/doc/html/rfc9204#name-literal-field-line-with-lit
// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-fields
func appendLiteralName(p []byte, name string, neverIndex bool) []byte {
	const (
		// layout of the first byte of the length of a name literal
		Prefix         byte = 0b0010_0000
		NeverIndex     byte = 0b0001_0000
		HuffmanEncoded byte = 0b0000_1000
		NameLenBits    byte = 0b0000_0111
	)
	prefix := Prefix
	if neverIndex {
		prefix = Prefix | NeverIndex
	}
	n := uint64(len(name))
	if h := huffman.EncodeLengthLower(name); h < n {
		p = appendVarint(p, h, NameLenBits, prefix|HuffmanEncoded)
		return huffman.AppendStringLower(p, name)
	}
	p = appendVarint(p, n, NameLenBits, prefix)
	return ascii.AppendLower(p, name)
}

func appendStringLiteralLength(p []byte, n uint64, huffmanEncoded bool) []byte {
	const (
		// layout of the first byte of the length of a string literal
		HuffmanEncoded byte = 0b1000_0000
		StringLenBits  byte = 0b0111_1111
	)
	var prefix byte
	if huffmanEncoded {
		prefix = HuffmanEncoded
	}
	return appendVarint(p, n, StringLenBits, prefix)
}

// appendStringLiteral appens the QPACK encoded string literal s to p.
func appendStringLiteral(p []byte, s string) []byte {
	n := uint64(len(s))
	if h := huffman.EncodeLength(s); h < n {
		p = appendStringLiteralLength(p, h, true)
		return huffman.AppendString(p, s)
	}
	p = appendStringLiteralLength(p, n, false)
	return append(p, s...)
}

// appendInt appends the QPACK string literal representation of int64 i.
func appendInt(p []byte, i int64) []byte {
	const HuffmanEncoded byte = 0b1000_0000

	if -9 <= i && i <= 99 {
		// No savings from huffman encoding 2 characters.
		if i < 0 {
			return append(p, 2, '-', byte('0'-i))
		}
		if i <= 9 {
			return append(p, 1, byte(i)+'0')
		}
		j := i / 10
		return append(p, 2, byte(j)+'0', byte(i-10*j)+'0')
	}

	j := len(p)
	p = append(p, 0)
	p = huffman.AppendInt(p, i)
	p[j] = HuffmanEncoded | uint8(len(p)-j-1)
	return p
}

// appendTime appends the QPACK string literal encoded RFC1123 representation
// of t to p.
func appendTime(p []byte, t time.Time) []byte {
	// RFC1123 time length is less 0x7F so only need a single byte for length
	const HuffmanEncoded byte = 0b1000_0000

	i := len(p)
	p = append(p, 0)
	p = huffman.AppendRFC1123Time(p, t)
	p[i] = HuffmanEncoded | uint8(len(p)-i-1)
	return p
}
