package quack

import (
	"errors"
)

type DT struct {
	headers []headerField
	base    uint64
	head    int
}

func (dt *DT) lookup(name, value string) (uint64, match) {
	return 0, matchNone
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoded-field-section-prefi
func (dt *DT) appendFieldSectionPrefix(p []byte, reqInsertCount uint64) []byte {
	var insertCount uint64

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-required-insert-count
	if reqInsertCount > 0 {
		maxEntries := uint64(cap(dt.headers)) / 32
		insertCount = (reqInsertCount % (2 * maxEntries)) + 1
	}
	p = appendVarint(p, insertCount, 0xFF, 0)

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-base
	var sign byte
	base := dt.base - reqInsertCount
	if dt.base < reqInsertCount {
		base = reqInsertCount - dt.base - 1
		sign = 0b1000_0000
	}
	return appendVarint(p, base, 0b0111_1111, sign)
}

func (dt *DT) insert(name, value string) int {
	i := dt.head
	if len(dt.headers) < cap(dt.headers) {
		// this never reallocs so cap() is const
		dt.headers = append(dt.headers, headerField{name, value})
	} else {
		dt.headers[i] = headerField{name, value}
	}
	dt.head = (i + 1) % cap(dt.headers)
	return i
}

func (dt *DT) nameIndex(index uint64) (string, error) {
	return "", errors.New("TODO unsupported")
}

func (dt *DT) baseNameIndex(index uint64) (string, error) {
	return "", errors.New("TODO unsupported")
}

func (dt *DT) lineIndex(index uint64) (string, string, error) {
	return "", "", errors.New("TODO unsupported")
}

func (dt *DT) baseLineIndex(index uint64) (string, string, error) {
	return "", "", errors.New("TODO unsupported")
}
