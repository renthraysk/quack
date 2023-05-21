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
	return fd.Decode(in, accept)
}

func (d *Decoder) ParseEncoderInstructions(in []byte) error {
	return d.dt.ParseEncoderInstructions(in)
}
