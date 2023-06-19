package field

type header struct {
	name  string
	value string
}

func (h header) size() uint64 {
	return headerSize(h.name, h.value)
}

// https://datatracker.ietf.org/doc/html/rfc9204#name-dynamic-table-size
func headerSize(name, value string) uint64 {
	return uint64(len(name)) + uint64(len(value)) + 32
}
