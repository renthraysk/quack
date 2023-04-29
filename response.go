package quack

import "time"

// https://www.rfc-editor.org/rfc/rfc9114.html#name-response-pseudo-header-fiel
func (e *Encoder) NewResponse(p []byte, statusCode int, header map[string][]string) ([]byte, error) {
	var reqInsertCount uint64

	p = e.dt.appendFieldSectionPrefix(p, reqInsertCount)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = appendStatus(p, statusCode)
	// Automagic the Date header if absent
	if _, ok := header["Date"]; !ok {
		p = appendDate(p, time.Now())
	}
	p = e.encodeHeader(p, header)
	return p, nil
}
