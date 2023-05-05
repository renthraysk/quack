package quack

import (
	"errors"
	"math/bits"
)

type DT struct {
	headers  []headerField
	base     uint64
	capacity uint64
	size     uint64
	evicted  uint64
}

func New(maxCapacity uint64) DT {
	return DT{
		capacity: maxCapacity,
	}
}

func (dt *DT) maxEntries() uint64 { return dt.capacity / 32 }

func (dt *DT) lookup(name, value string) (uint64, match) {
	i := len(dt.headers) - 1
	for i >= 0 && dt.headers[i].Name != name {
		i--
	}
	if i < 0 || i >= len(dt.headers) {
		return 0, matchNone
	}
	if dt.headers[i].Value == value {
		return uint64(i), matchNameValue
	}
	j := i - 1
	for j >= 0 && dt.headers[j].Name != name || dt.headers[j].Value != value {
		j--
	}
	if j < 0 {
		return uint64(i), matchName
	}
	return uint64(j), matchNameValue
}

// evict evicts headerFields until size is less than targetSize, returns the
// number of headerFields evicted.
func (dt *DT) evict(targetSize uint64) (uint64, bool) {
	var i int

	size := dt.size
	for i < len(dt.headers) && size > targetSize {
		size -= dt.headers[i].size()
		i++
	}
	if i == 0 {
		return 0, true
	}
	if i >= len(dt.headers) {
		return 0, false
	}
	n := uint64(i)
	if ec, c := bits.Add64(dt.evicted, n, 0); c != 0 {
		return 0, false
	} else {
		dt.evicted = ec
	}
	dt.headers = append(dt.headers[:0], dt.headers[n:]...)
	dt.size = size
	return n, true
}

func (dt *DT) insert(name, value string) (int, bool) {
	s := size(name, value)
	if s > dt.capacity-dt.size {
		if s > dt.capacity {
			return 0, false
		}
		if _, ok := dt.evict(dt.capacity - s); !ok {
			return 0, false
		}
	}
	i := len(dt.headers)
	dt.headers = append(dt.headers, headerField{name, value})
	dt.size += s
	return i, true
}

// https://www.rfc-editor.org/rfc/rfc9204.html#name-encoded-field-section-prefi
func (dt *DT) appendFieldSectionPrefix(p []byte, reqInsertCount uint64) []byte {
	var insertCount uint64

	// https://www.rfc-editor.org/rfc/rfc9204.html#name-required-insert-count
	if reqInsertCount > 0 {
		insertCount = (reqInsertCount % (2 * dt.maxEntries())) + 1
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
