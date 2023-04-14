package quack

import "errors"

type DT struct {
	ents        []headerField
	capacity    uint
	size        uint
	insertCount uint
	head        uint
	tail        uint
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
