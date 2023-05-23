package quack

import (
	"sync/atomic"

	"github.com/renthraysk/quack/internal/field"
)

type Decoder struct {
	dt           field.DT
	fieldDecoder atomic.Pointer[field.Decoder]
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(in []byte, accept func(name, value string)) error {
	fd := d.fieldDecoder.Load()
	if err := fd.Decode(in, accept); err != nil {
		return ErrDecompressionFailed{err}
	}
	return nil
}

func (d *Decoder) ParseEncoderInstructions(in []byte) error {
	if err := d.dt.ParseEncoderInstructions(in); err != nil {
		return ErrEncoderStream{err}
	}
	return nil
}
