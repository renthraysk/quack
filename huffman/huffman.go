package huffman

const (
	// Huffman code length limits
	minCodeLength = 5
	maxCodeLength = 30

	eosCode = 0x3fffffff
	eosLen  = 30

	// maxDecodedLength the absolute maximum that the decoder should write
	maxDecodedLength = 1 << 20 // 1MB

	// maxEncodedLength the maximum length of input we should attempt to decode
	// Best compression ratio is 5 (minCodeLength) bits of input results in
	// 8 bits of output.
	maxEncodedLength = minCodeLength * maxDecodedLength / 8
)

type errorString string

func (e errorString) Error() string { return (string)(e) }

const errExpectedEOS = errorString("huffman: expected EOS")
const errEOSEncoded = errorString("huffman: EOS encoded")
const errInputTooLong = errorString("huffman: input too long")
