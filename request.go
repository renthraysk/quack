package quack

import (
	"errors"

	"github.com/renthraysk/quack/internal/field"
)

// https://www.rfc-editor.org/rfc/rfc9114.html#name-request-pseudo-header-field
func (e *Encoder) NewRequest(p []byte, method, scheme, authority, path string, header map[string][]string) ([]byte, error) {
	if scheme == "http" || scheme == "https" {

		// Clients that generate HTTP/3 requests directly SHOULD use the
		// :authority pseudo-header field instead of the Host header field.
		if authority == "" {
			return p, errors.New("empty :authority")
		}
		delete(header, "Host")

		// This pseudo-header field MUST NOT be empty for "http" or "https" URIs;
		// "http" or "https" URIs that do not contain a path component MUST
		// include a value of / (ASCII 0x2f).
		if path == "" {
			path = "/"
			// An OPTIONS request that does not include a path component includes
			// the value * (ASCII 0x2a) for the :path pseudo-header field
			if method == "OPTIONS" {
				path = "*"
			}
		}
	}

	fe := e.current.Load()
	p = fe.AppendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = field.AppendMethod(p, method)
	p = field.AppendScheme(p, scheme)
	p = field.AppendAuthority(p, authority)
	p = field.AppendPath(p, path)
	p = fe.AppendFieldLines(p, header)
	return p, nil
}

// https://www.rfc-editor.org/rfc/rfc9114.html#name-the-connect-method
func (e *Encoder) NewConnect(p []byte, authority string, header map[string][]string) ([]byte, error) {
	fe := e.current.Load()
	p = fe.AppendFieldSectionPrefix(p)
	// All pseudo-header fields MUST appear in the header section before regular header fields.
	// https://www.rfc-editor.org/rfc/rfc9114.html#name-http-control-data
	p = field.AppendMethod(p, "CONNECT")
	p = field.AppendAuthority(p, authority)
	p = fe.AppendFieldLines(p, header)
	return p, nil
}
