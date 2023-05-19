package quack

import (
	"time"

	"github.com/renthraysk/quack/internal/field"
)

// https://www.rfc-editor.org/rfc/rfc9114.html#name-response-pseudo-header-fiel
func (e *Encoder) NewResponse(p []byte, statusCode int, header map[string][]string) ([]byte, error) {

	fe := e.fieldEncoder.Load()
	p = fe.AppendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = field.AppendStatus(p, statusCode)
	// Automagic the Date header if absent
	if _, ok := header["Date"]; !ok {
		p = field.AppendDate(p, time.Now())
	}
	p = fe.AppendFieldLines(p, header)
	return p, nil
}
