package huffman

const (
	// Huffman code length limits
	minCodeLength = 5
	maxCodeLength = 30

	eosCode = 0x3fffffff
	eosLen  = 30
)

type errorString string

func (e errorString) Error() string { return (string)(e) }

const errExpectedEOS = errorString("huffman: expected EOS")
const errEOSEncoded = errorString("huffman: EOS encoded")
