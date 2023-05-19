package quack

import (
	"sync/atomic"

	"github.com/renthraysk/quack/internal/field"
	"github.com/renthraysk/quack/varint"
)

type Encoder struct {
	// Dynamic table
	dt DT

	// fieldEncoder is the encoder state required to encode headers.
	// It is immutable once created by the dynamic table.
	// If nil then will only encode using the static table. Will be updated
	// when receive the increment decoder instruction from peer.
	fieldEncoder atomic.Pointer[field.Encoder]
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
			increment, p, err = varint.Read(p, 0b0011_1111)
			if err != nil {
				return err
			}
			if err := e.increment(increment); err != nil {
				return err
			}

		case 0b01:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-stream-cancellation
			streamID, p, err = varint.Read(p, 0b0011_1111)
			if err != nil {
				return err
			}
			e.streamInstruction(streamID, (*stream).streamCancellation, true)

		case 0b10, 0b11:
			// https://www.rfc-editor.org/rfc/rfc9204.html#name-section-acknowledgment
			streamID, p, err = varint.Read(p, 0b0111_1111)
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

func t(b bool, t byte) byte {
	if b {
		return t
	}
	return 0
}
