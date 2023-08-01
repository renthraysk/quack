package field

import (
	"math/bits"
	"time"

	"github.com/renthraysk/quack/huffman"
	"github.com/renthraysk/quack/internal/inst"
	"github.com/renthraysk/quack/varint"
)

// match returned status of a table search
type match uint

const (
	// matchNone No match
	matchNone match = iota
	// matchName Matched name only
	matchName
	// matchNameValue Matched name & value
	matchNameValue
)

type value struct {
	value string
	index uint64
}

type nameValues map[string][]value

// Encoder field line encoder, immutable once created
type Encoder struct {
	nv          nameValues
	base        uint64
	insertCount uint64
	capacity    uint64
}

func newEncoder(nv nameValues, base, insertCount, capacity uint64) *Encoder {
	return &Encoder{
		nv:          nv,
		base:        base,
		insertCount: insertCount,
		capacity:    capacity,
	}
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-request-pseudo-header-field
func (fe *Encoder) AppendRequest(p []byte, method, scheme, authority, path string, header map[string][]string) []byte {
	p = fe.appendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = appendMethod(p, method)
	p = appendScheme(p, scheme)
	p = appendAuthority(p, authority)
	p = appendPath(p, path)
	p = fe.appendFieldLines(p, header)
	return p
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-response-pseudo-header-fiel
func (fe *Encoder) AppendResponse(p []byte, statusCode int, header map[string][]string) []byte {
	p = fe.appendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = appendStatus(p, statusCode)

	if statusCode <= 100 || statusCode >= 200 {
		// Automagic the Date header if absent
		if _, ok := header["Date"]; !ok {
			p = appendDate(p, time.Now())
		}
	}
	p = fe.appendFieldLines(p, header)
	return p
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-the-connect-method
func (fe *Encoder) AppendConnect(p []byte, authority string, header map[string][]string) []byte {
	p = fe.appendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = appendMethod(p, "CONNECT")
	p = appendAuthority(p, authority)
	p = fe.appendFieldLines(p, header)
	return p
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoded-field-section-prefi
func (fe *Encoder) appendFieldSectionPrefix(p []byte) []byte {
	if fe == nil {
		// Operating with only static table
		return append(p, 0, 0)
	}

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-required-insert-count
	maxEntries := fe.capacity / 32
	p = varint.Append(p, 0, 0xFF, (fe.insertCount%(2*maxEntries))+1)

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-base
	deltaBase, sign := bits.Sub64(fe.base, fe.insertCount, 0)
	if sign != 0 {
		deltaBase = fe.insertCount - fe.base - 1
		sign = 0x80
	}
	return varint.Append(p, byte(sign), 0x7F, deltaBase)
}

func (fe *Encoder) lookup(name, value string) (index uint64, isStatic bool, m match) {
	index, m = staticLookup(name, value)
	if fe == nil || m == matchNameValue {
		// Operating with only static table or have the best match already.
		return index, true, m
	}
	values, ok := fe.nv[name]
	if ok {
		for _, v := range values {
			if v.value == value {
				return v.index, false, matchNameValue
			}
		}
	}
	switch {
	case m == matchName:
		return index, true, matchName
	case ok && len(values) > 0:
		return values[0].index, false, matchName
	}
	return 0, false, matchNone
}

func (fe *Encoder) appendFieldLines(p []byte, header map[string][]string) []byte {
	for name, values := range header {
		for _, value := range values {
			p = fe.appendFieldLine(p, name, value)
		}
	}
	return p
}

func (fe *Encoder) appendFieldLine(p []byte, name, value string) []byte {
	ctrl := headerControl(name)
	switch i, isStatic, m := fe.lookup(name, value); m {
	case matchNameValue:
		if isStatic {
			return inst.AppendIndexedLine(p, i, true)
		}
		return inst.AppendIndexedLinePostBase(p, i)

	case matchName:
		p = inst.AppendNamedReference(p, i, ctrl.neverIndex(), isStatic)
	case matchNone:
		p = inst.AppendLiteralName(p, name, ctrl.neverIndex())
	}
	return inst.AppendStringLiteral(p, value, ctrl.shouldHuffman())
}

/* */

// The pseudo headers

// appendAuthority appends an :authority pseudo header field to p
func appendAuthority(p []byte, authority string) []byte {
	if authority == "" {
		return inst.AppendIndexedLine(p, 0, true)
	}
	p = inst.AppendNamedReference(p, 0, false, true)
	return inst.AppendStringLiteral(p, authority, true)
}

// appendPath appends a :path pseudo header field to p
func appendPath(p []byte, path string) []byte {
	if path == "/" {
		return inst.AppendIndexedLine(p, 1, true)
	}
	p = inst.AppendNamedReference(p, 1, false, true)
	return inst.AppendStringLiteral(p, path, true)
}

// appendStatus appends a :status pseudo header field to p
func appendStatus(p []byte, statusCode int) []byte {
	switch statusCode {
	case 100:
		return inst.AppendIndexedLine(p, 63, true)
	case 103:
		return inst.AppendIndexedLine(p, 24, true)
	case 200:
		return inst.AppendIndexedLine(p, 25, true)
	case 204:
		return inst.AppendIndexedLine(p, 64, true)
	case 206:
		return inst.AppendIndexedLine(p, 65, true)
	case 302:
		return inst.AppendIndexedLine(p, 66, true)
	case 304:
		return inst.AppendIndexedLine(p, 26, true)
	case 400:
		return inst.AppendIndexedLine(p, 67, true)
	case 403:
		return inst.AppendIndexedLine(p, 68, true)
	case 404:
		return inst.AppendIndexedLine(p, 27, true)
	case 421:
		return inst.AppendIndexedLine(p, 69, true)
	case 425:
		return inst.AppendIndexedLine(p, 70, true)
	case 500:
		return inst.AppendIndexedLine(p, 71, true)
	case 503:
		return inst.AppendIndexedLine(p, 28, true)
	}
	p = inst.AppendNamedReference(p, 24, false, true)
	return appendInt(p, int64(statusCode))
}

// appendInt appends the QPACK string literal representation of int64 i.
func appendInt(p []byte, i int64) []byte {
	// H HuffmanEncoded
	const H = 0b1000_0000

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
	p[j] = H | uint8(len(p)-j-1)
	return p
}

// appendMethod appends a :method pseudo header field to p
func appendMethod(p []byte, method string) []byte {
	switch method {
	case "CONNECT":
		return inst.AppendIndexedLine(p, 15, true)
	case "DELETE":
		return inst.AppendIndexedLine(p, 16, true)
	case "GET":
		return inst.AppendIndexedLine(p, 17, true)
	case "HEAD":
		return inst.AppendIndexedLine(p, 18, true)
	case "OPTIONS":
		return inst.AppendIndexedLine(p, 19, true)
	case "POST":
		return inst.AppendIndexedLine(p, 20, true)
	case "PUT":
		return inst.AppendIndexedLine(p, 21, true)
	}
	p = inst.AppendNamedReference(p, 15, false, true)
	return inst.AppendStringLiteral(p, method, true)
}

// appendScheme appends a :scheme pseudo header field to p
func appendScheme(p []byte, scheme string) []byte {
	switch scheme {
	case "http":
		return inst.AppendIndexedLine(p, 22, true)
	case "https":
		return inst.AppendIndexedLine(p, 23, true)
	}
	p = inst.AppendNamedReference(p, 22, false, true)
	return inst.AppendStringLiteral(p, scheme, true)
}

// Regular headers

// appendDate appends a Date header field with time t.
func appendDate(p []byte, t time.Time) []byte {
	const StaticTableIndex = 6

	const H = 0b1000_0000

	p = inst.AppendNamedReference(p, StaticTableIndex, false, true)
	// RFC1123 time length is less 0x7F so only need a single byte for length
	i := len(p)
	p = append(p, 0)
	p = huffman.AppendHttpTime(p, t)
	p[i] = H | uint8(len(p)-i-1)
	return p
}
