package inst

import "github.com/renthraysk/quack/varint"

// Decoder instructions

// https://www.rfc-editor.org/rfc/rfc9204.html#name-section-acknowledgment
func AppendSectionAcknowledgement(q []byte, streamID uint64) []byte {
	const (
		P = 0b1000_0000
		M = 0b0111_1111
	)
	return varint.Append(q, streamID, M, P)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-stream-cancellation
func AppendStreamCancellation(q []byte, streamID uint64) []byte {
	const (
		P = 0b0100_0000
		M = 0b0011_1111
	)
	return varint.Append(q, streamID, M, P)
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-insert-count-increment
func AppendInsertCountIncrement(q []byte, increment uint64) []byte {
	const (
		P = 0b0000_0000
		M = 0b0011_1111
	)
	return varint.Append(q, increment, M, P)
}
