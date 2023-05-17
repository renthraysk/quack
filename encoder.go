package quack

import (
	"sync/atomic"

	"github.com/renthraysk/quack/huffman"
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

type Encoder struct {
	// Dynamic table
	dt DT

	// current is the encoder state required to encode headers.
	// If nil then will only encode using the static table. Will be updated
	// when receive the increment decoder instruction from peer.
	current atomic.Pointer[fieldEncoder]
}

func NewEncoder(capacity uint64) *Encoder {
	return &Encoder{
		dt: DT{capacity: capacity},
	}
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-decoder-instructions
func (e *Encoder) readDecoderInstructions(p []byte) error {
	var streamID, increment uint64
	var err error

	for len(p) > 0 {
		switch p[0] >> 6 {
		case 0b00:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-count-increment
			increment, p, err = readVarint(p, 0b0011_1111)
			if err != nil {
				return err
			}
			if err := e.increment(increment); err != nil {
				return err
			}

		case 0b01:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-stream-cancellation
			streamID, p, err = readVarint(p, 0b0011_1111)
			if err != nil {
				return err
			}
			e.streamInstruction(streamID, (*stream).streamCancellation, true)

		case 0b10, 0b11:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-section-acknowledgment
			streamID, p, err = readVarint(p, 0b0111_1111)
			if err != nil {
				return err
			}
			e.streamInstruction(streamID, (*stream).sectionAcknowledgement, false)
		}
	}
	return nil
}

func (e *Encoder) streamInstruction(streamID uint64, f func(s *stream), remove bool) {
}

func (e *Encoder) increment(increment uint64) error {
	return nil
}

// appendStringLiteral appends the QPACK encoded string literal s to p.
func appendStringLiteral(p []byte, s string, shouldHuffman bool) []byte {
	const (
		H = 0b1000_0000
		M = 0b0111_1111
	)

	n := uint64(len(s))
	if n > 2 && shouldHuffman {
		if h := huffman.EncodeLength(s); h < n {
			p = appendVarint(p, h, M, H)
			return huffman.AppendString(p, s)
		}
	}
	p = appendVarint(p, n, M, 0)
	return append(p, s...)
}

func t(b bool, t byte) byte {
	if b {
		return t
	}
	return 0
}
