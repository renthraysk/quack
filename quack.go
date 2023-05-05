package quack

type headerField struct {
	Name  string
	Value string
}

func (hf headerField) size() uint64 {
	return size(hf.Name, hf.Value)
}

func size(name, value string) uint64 {
	return uint64(len(name)) + uint64(len(value)) + 32
}
