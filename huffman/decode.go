package huffman

// Decode decodes huffman encoded data in, and appends to dst and returns the
// result. Will return an error if the input is incorrectly padded, or if the
// EOS code has been encoded.
func Decode(dst, in []byte) ([]byte, error) {
	// The maximum code length is 30.
	var x uint64
	var n uint

	for len(in) >= 4 {
		x <<= 32
		x |= uint64(in[3]) | uint64(in[2])<<8 | uint64(in[1])<<16 | uint64(in[0])<<24
		in = in[4:]
		n += 32
		for n >= 32 {
			b, s := codeLookup(uint32(x >> (n % 32))) // n<=59 so %32 is fine
			if s == 0 {
				return nil, errEOSEncoded
			}
			n -= uint(s)
			dst = append(dst, b)
		}
	}
	// n < 32 and len(in) < 4, so can load up remaining bytes.
	for i := 0; i < len(in); i++ {
		x <<= 8
		x |= uint64(in[i])
		n += 8
	}
	for n >= 32 {
		b, s := codeLookup(uint32(x >> (n % 32)))
		if s == 0 {
			return nil, errEOSEncoded
		}
		n -= uint(s)
		dst = append(dst, b)
	}
	for y := uint32(x << (32 - n)); n >= minCodeLength; {
		b, s := codeLookup(y)
		if s == 0 {
			return nil, errEOSEncoded
		}
		if uint(s) > n {
			break
		}
		n -= uint(s)
		y <<= s
		dst = append(dst, b)
	}
	if m := uint64(1<<(n%64)) - 1; x&m != m {
		return dst, errExpectedEOS
	}
	return dst, nil
}

// codeLookup takes a 32 bit value with the code to be decoded in the most
// significant bits, will return the symbol and bit length.
// If it encounters the EOS symbol, it will return a length of 0.
func codeLookup(x uint32) (sym byte, length uint8) {
	// inlines, just.
	// Fast path for codes with lengths less than or equal to maxShortCode
	if i := x >> (32 - maxShortCode); i < uint32(len(shortCodes)) {
		b := shortCodes[i]
		return b, codeLengths[b]
	}
	// slow path for longer codes
	x -= delta
	for _, y := range bounds {
		if x < y.delta {
			// longCodes is const string 256 bytes long, using uint8 for BCE
			return longCodes[y.offset+uint8(x>>((32-y.length)%32))], y.length
		}
		x -= y.delta
	}
	// codes of bit length 30, which includes EOS
	x >>= 32 - 30
	x += offset30
	if x >= uint32(len(longCodes)) {
		return 0, 0 // x == 256, eosCode 0x3FFFFFFF found
	}
	return longCodes[x], 30
}
