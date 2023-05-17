package quack

// control controls finer details of how a specific headers should be encoded.
type control uint8

const (
	// neverIndex header field should never be put in the dynamic table.
	neverIndex control = 1 << iota
	// neverHuffman never compress the value field.
	neverHuffman
)

func (c control) shouldHuffman() bool { return c&neverHuffman == 0 }
func (c control) neverIndex() bool    { return c&neverIndex != 0 }

func headerControl(name string) control {
	switch name {
	case "Authorization", "Content-Md5":
		return neverIndex | neverHuffman
	case "Date",
		"Etag",
		"If-Modified-Since",
		"If-Unmodified-Since",
		"Last-Modified",
		"Location",
		"Match",
		"Range",
		"Retry-After",
		"Set-Cookie":
		return neverIndex
	}
	return 0
}
