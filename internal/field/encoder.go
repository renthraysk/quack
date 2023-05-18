package field

import (
	"math/bits"

	"github.com/renthraysk/quack/internal/inst"
	"github.com/renthraysk/quack/varint"
)

// match returned status of a table search
type Match uint

const (
	// MatchNone No match
	MatchNone Match = iota
	// MatchName Matched name only
	MatchName
	// MatchNameValue Matched name & value
	MatchNameValue
)

type Value struct {
	Value string
	Index uint64
}

type NameValues map[string][]Value

// fieldEncoder field line encoder, immutable once created
type Encoder struct {
	nv             NameValues
	reqInsertCount uint64
	base           uint64
	capacity       uint64
}

func New(nv NameValues, reqInsertCount, base, capacity uint64) *Encoder {
	return &Encoder{
		nv:             nv,
		reqInsertCount: reqInsertCount,
		base:           base,
		capacity:       capacity,
	}
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoded-field-section-prefi
func (fe *Encoder) AppendFieldSectionPrefix(p []byte) []byte {
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

func (fe *Encoder) Lookup(name, value string) (uint64, bool, Match) {
	i, m := staticLookup(name, value)
	if fe == nil || m == MatchNameValue {
		// Operating with only static table or have the best match already.
		return i, true, m
	}
	values, ok := fe.nv[name]
	if ok {
		for _, v := range values {
			if v.Value == value {
				return v.Index, false, MatchNameValue
			}
		}
	}
	switch {
	case m == MatchName:
		return i, true, MatchName
	case ok && len(values) > 0:
		return values[0].Index, false, MatchName
	}
	return 0, false, MatchNone
}

func (fe *Encoder) AppendFieldLines(p []byte, header map[string][]string) []byte {
	for name, values := range header {
		for _, value := range values {
			p = fe.appendFieldLine(p, name, value)
		}
	}
	return p
}

func (fe *Encoder) appendFieldLine(p []byte, name, value string) []byte {
	ctrl := HeaderControl(name)
	switch i, isStatic, m := fe.Lookup(name, value); m {
	case MatchNameValue:
		if isStatic {
			return inst.AppendIndexedLine(p, i, true)
		}
		return inst.AppendIndexedLinePostBase(p, i)

	case MatchName:
		p = inst.AppendNamedReference(p, i, ctrl.NeverIndex(), isStatic)
	case MatchNone:
		p = inst.AppendLiteralName(p, name, ctrl.NeverIndex())
	}
	return inst.AppendStringLiteral(p, value, ctrl.ShouldHuffman())
}

/* */
