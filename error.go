package quack

import (
	"fmt"
)

type Error interface {
	error
	ErrorCode() uint16
}

type ErrDecompressionFailed struct {
	err error
}

func (e ErrDecompressionFailed) ErrorCode() uint16 { return 0x0200 }

func (e ErrDecompressionFailed) Unwrap() error {
	return e.err
}

func (e ErrDecompressionFailed) Error() string {
	return fmt.Sprintf("qpack: decompression failed: %s", e.err.Error())
}

type ErrEncoderStream struct {
	err error
}

func (e ErrEncoderStream) ErrorCode() uint16 { return 0x0201 }

func (e ErrEncoderStream) Unwrap() error {
	return e.err
}

func (e ErrEncoderStream) Error() string {
	return fmt.Sprintf("qpack: encoder stream error: %s", e.err.Error())
}

type ErrDecoderStream struct {
	err error
}

func (e ErrDecoderStream) ErrorCode() uint16 { return 0x0202 }

func (e ErrDecoderStream) Unwrap() error {
	return e.err
}

func (e ErrDecoderStream) Error() string {
	return fmt.Sprintf("qpack: decoder stream error: %s", e.err.Error())
}
