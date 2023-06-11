package quack

import (
	"errors"
	"sync/atomic"

	"github.com/renthraysk/quack/ascii"
	"github.com/renthraysk/quack/internal/field"
	"github.com/renthraysk/quack/varint"
)

type Encoder struct {
	// Dynamic table
	dt field.DT

	// fieldEncoder is the encoder state required to encode headers.
	// It is immutable once created by the dynamic table.
	// If nil then will only encode using the static table. Will be updated
	// when receive the increment decoder instruction from peer.
	fieldEncoder atomic.Pointer[field.Encoder]
}

func NewEncoder() *Encoder {
	return &Encoder{}
}

func allEqual[T comparable](s []T, one T) bool {
	for _, v := range s {
		if v != one {
			return false
		}
	}
	return true
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-request-pseudo-header-field
func (e *Encoder) AppendRequest(p []byte, method, scheme, authority, path string, header map[string][]string) ([]byte, error) {

	if ascii.EqualI(scheme, "https") || ascii.EqualI(scheme, "http") {

		// Clients that generate HTTP/3 requests directly SHOULD use the
		// :authority pseudo-header field instead of the Host header field.
		if authority == "" {
			return p, errors.New("empty :authority")
		} else if hosts, ok := header["Host"]; ok && !allEqual(hosts, authority) {
			return p, errors.New(":authority and Host header are inconsistent")
		}

		// This pseudo-header field MUST NOT be empty for "http" or "https" URIs;
		// "http" or "https" URIs that do not contain a path component MUST
		// include a value of / (ASCII 0x2f).
		if path == "" {
			path = "/"
			// An OPTIONS request that does not include a path component includes
			// the value * (ASCII 0x2a) for the :path pseudo-header field
			if method == "OPTIONS" {
				path = "*"
			}
		}
	}

	fe := e.fieldEncoder.Load()
	p = fe.AppendRequest(p, method, scheme, authority, path, header)
	return p, nil
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-the-connect-method
func (e *Encoder) AppendConnect(p []byte, authority string, header map[string][]string) ([]byte, error) {
	fe := e.fieldEncoder.Load()
	p = fe.AppendConnect(p, authority, header)
	return p, nil
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-response-pseudo-header-fiel
func (e *Encoder) AppendResponse(p []byte, statusCode int, header map[string][]string) ([]byte, error) {
	fe := e.fieldEncoder.Load()
	p = fe.AppendResponse(p, statusCode, header)
	return p, nil
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
