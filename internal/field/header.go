package field

type Header struct {
	Name  string
	Value string
}

func (h Header) size() uint64 {
	return headerSize(h.Name, h.Value)
}

// https://datatracker.ietf.org/doc/html/rfc9204#name-dynamic-table-size
func headerSize(name, value string) uint64 {
	return uint64(len(name)) + uint64(len(value)) + 32
}
