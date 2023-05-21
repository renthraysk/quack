package field

import (
	"math/bits"
	"time"

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
	nv             nameValues
	reqInsertCount uint64
	base           uint64
	capacity       uint64
}

func newEncoder(nv nameValues, reqInsertCount, base, capacity uint64) *Encoder {
	return &Encoder{
		nv:             nv,
		reqInsertCount: reqInsertCount,
		base:           base,
		capacity:       capacity,
	}
}

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

func (fe *Encoder) AppendResponse(p []byte, statusCode int, header map[string][]string) []byte {
	p = fe.appendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = appendStatus(p, statusCode)
	// Automagic the Date header if absent
	if _, ok := header["Date"]; !ok {
		p = appendDate(p, time.Now())
	}
	p = fe.appendFieldLines(p, header)
	return p
}

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
		// Operating with only static table.
		return append(p, 0, 0)
	}

	var encodedInsertCount uint64

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-required-insert-count
	if fe.reqInsertCount > 0 {
		maxEntries := fe.capacity / 32
		encodedInsertCount = (fe.reqInsertCount % (2 * maxEntries)) + 1
	}
	// https://www.rfc-editor.org/rfc/rfc9204.html#name-base
	deltaBase, sign := bits.Sub64(fe.base, fe.reqInsertCount, 0)
	if sign != 0 {
		deltaBase = fe.reqInsertCount - fe.base - 1
		sign = 0x80
	}
	p = varint.Append(p, encodedInsertCount, 0xFF, 0)
	return varint.Append(p, deltaBase, 0b0111_1111, byte(sign))
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
