package ascii

const (
	sp     = 1 << ' '
	htab   = 1 << '\t'
	lower  = ((1 << 26) - 1) << 'a' // 'a'-'z'
	upper  = ((1 << 26) - 1) << 'A' // 'A'-'Z'
	alpha  = lower | upper
	digits = ((1 << 10) - 1) << '0'              // '0'-'9'
	vchar  = ((1 << (1 + '~' - '!')) - 1) << '!' // '!'-'~'

	/*
		tchar          = "!" / "#" / "$" / "%" / "&" / "'" / "*"
			/ "+" / "-" / "." / "^" / "_" / "`" / "|" / "~"
		    / DIGIT / ALPHA
			; any VCHAR, except delimiters
	*/
	special = 1<<'!' | 1<<'#' | 1<<'$' | 1<<'%' | 1<<'&' | 1<<'\'' | 1<<'*' |
		1<<'+' | 1<<'-' | 1<<'.' | 1<<'^' | 1<<'_' | 1<<'`' | 1<<'|' |
		1<<'~' // !#$%&'*+-.^_`|~
)

// isUpper returns true if c is a upper case ASCII chararacter.
func isUpper(c byte) bool { return c-'A' <= 'Z'-'A' }

// ToLower returns the ASCII lowercase version of c.
func ToLower(c byte) byte {
	var x int
	if isUpper(c) {
		x = 'a' - 'A'
	}
	return c + byte(x)
}

func isIn(c byte, lo, hi uint64) bool {
	m := lo
	if c >= 64 {
		m = hi
	}
	return (1<<(c%64))&m != 0
}

func isTokenChar(c byte) bool {
	// token is the set of characters allowed in pre HTTP/3 names.
	const token = alpha | digits | special

	return c < 0x80 && isIn(c, token%(1<<64), token>>64)
}

// isToken3Char returns true if byte c is a valid in a HTTP/3 name literal, false
// otherwise.
func isToken3Char(c byte) bool {
	// token3 is the set of characters allowed in HTTP/3 incoming names.
	const token3 = lower | digits | special

	return c < 0x80 && isIn(c, token3%(1<<64), token3>>64)
}

func IsNameValid[T string | []byte](n T) bool {
	for i := 0; i < len(n); i++ {
		if !isTokenChar(n[i]) {
			return false
		}
	}
	return true
}

func IsName3Valid[T string | []byte](n T) bool {
	for i := 0; i < len(n); i++ {
		if !isToken3Char(n[i]) {
			return false
		}
	}
	return true
}

// isFieldChar returns true if c is in the field-vchar set, false otherwise
// field-vchar    = VCHAR / obs-text
// obs-text       = %x80-FF
func isFieldVChar(c byte) bool {
	return c >= 0x80 || isIn(c, vchar%(1<<64), vchar>>64)
}

// isFieldContent returns true if c is in the field-content set, false otherwise
// field-content  = field-vchar
// [ 1*( SP / HTAB / field-vchar ) field-vchar ]
// obs-text       = %x80-FF
// https://www.rfc-editor.org/rfc/rfc9110#name-field-values
func isFieldContent(c byte) bool {
	const fieldContent = sp | htab | vchar

	return c >= 0x80 || isIn(c, fieldContent%(1<<64), fieldContent>>64)
}

// https://www.rfc-editor.org/rfc/rfc9110#section-5.5
func IsValueValid[T []byte | string](v T) bool {
	// An empty value is valid
	if len(v) <= 0 {
		return true
	}
	// Has to start with a field-vchar
	if !isFieldVChar(v[0]) {
		return false
	}
	if len(v) < 2 {
		return true
	}
	for i := 1; i < len(v); i++ {
		// double checking last char is ok
		// as field-vchar is subset of field-content
		if !isFieldContent(v[i]) {
			return false
		}
	}
	// Has to end with a field-vchar
	return isFieldVChar(v[len(v)-1])
}

// AppendLower appends the lower cased version of s to p.
func AppendLower(p []byte, s string) []byte {
	i := len(p)
	p = append(p, s...)
	for ; i < len(p); i++ {
		if c := p[i]; isUpper(c) {
			p[i] = c | 0x20
		}
	}
	return p
}

// Lower return the lower cased representation of string s.
func Lower(s string) string {
	// inlines
	var i int
	for i < len(s) && !isUpper(s[i]) {
		i++
	}
	if i >= len(s) {
		return s
	}
	b := append(make([]byte, 0, 32), s...)
	for ; i < len(b); i++ {
		if isUpper(b[i]) {
			b[i] |= 0x20
		}
	}
	return string(b)
}

func ToCanonical(b []byte) string {
	nextA := 'a'
	for i, c := range b {
		if c-byte(nextA) < 26 {
			b[i] = c ^ 0x20
		}
		nextA = 'A'
		if c == '-' {
			nextA = 'a'
		}
	}
	return string(b)
}
