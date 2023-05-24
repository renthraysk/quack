package quack

import (
	"fmt"
)

const (
	QpackDecompressionFailed = 0x0200
	QpackEncoderStreamError  = 0x0201
	QpackDecoderStreamError  = 0x0202
)

// https://www.rfc-editor.org/rfc/rfc9204.html#name-error-handling

type Error interface {
	error
	ErrorCode() uint16
}

type ErrDecompressionFailed struct {
	err error
}

func (e ErrDecompressionFailed) ErrorCode() uint16 { return QpackDecompressionFailed }

func (e ErrDecompressionFailed) Unwrap() error {
	return e.err
}

func (e ErrDecompressionFailed) Error() string {
	return fmt.Sprintf("qpack: decompression failed: %s", e.err.Error())
}

type ErrEncoderStream struct {
	err error
}

func (e ErrEncoderStream) ErrorCode() uint16 { return QpackEncoderStreamError }

func (e ErrEncoderStream) Unwrap() error {
	return e.err
}

func (e ErrEncoderStream) Error() string {
	return fmt.Sprintf("qpack: encoder stream error: %s", e.err.Error())
}

type ErrDecoderStream struct {
	err error
}

func (e ErrDecoderStream) ErrorCode() uint16 { return QpackDecoderStreamError }

func (e ErrDecoderStream) Unwrap() error {
	return e.err
}

func (e ErrDecoderStream) Error() string {
	return fmt.Sprintf("qpack: decoder stream error: %s", e.err.Error())
}
