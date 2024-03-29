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
	var x byte
	if isUpper(c) {
		x = 1
	}
	return c + x*('a'-'A')
}

// isIn returns if byte c is in a 128 bit set represented in a lo, the lower 64
// bit mask, and hi the upper 64 bits of the set.
func isIn(c byte, lo, hi uint64) bool {
	var mask uint64
	if c < 128 {
		mask = hi
	}
	if c < 64 {
		mask = lo
	}
	return (1<<(c%64))&mask != 0
}

func isTokenChar(c byte) bool {
	// token is the set of characters allowed in pre HTTP/3 names.
	const token = alpha | digits | special

	return isIn(c, token%(1<<64), token>>64)
}

// isToken3Char returns true if byte c is a valid in a HTTP/3 name literal, false
// otherwise.
func isToken3Char(c byte) bool {
	// token3 is the set of characters allowed in HTTP/3 incoming names.
	const token3 = lower | digits | special

	return isIn(c, token3%(1<<64), token3>>64)
}

func IsNameValid[T string | []byte](n T) bool {
	for i := range len(n) {
		if !isTokenChar(n[i]) {
			return false
		}
	}
	return true
}

func IsName3Valid[T string | []byte](n T) bool {
	for i := range len(n) {
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
	// return c >= '!' && c != '\x7F'
	// split to avoid branching
	x := c >= '!'
	y := c != '\x7F'
	return x && y
}

// isFieldContent returns true if c is in the field-content set, false otherwise
// field-content  = field-vchar
// [ 1*( SP / HTAB / field-vchar ) field-vchar ]
// field-vchar    = VCHAR / obs-text
// obs-text       = %x80-FF
// https://www.rfc-editor.org/rfc/rfc9110#name-field-values
func isFieldContent(c byte) bool {
	// return (c >= ' ' && c != '\x7F') || c == '\t'
	// split to avoid branching
	x := c >= ' '
	if c == '\t' {
		x = true
	}
	y := c != '\x7F'
	return x && y
}

// https://www.rfc-editor.org/rfc/rfc9110#section-5.5
func IsValueValid[T ~[]byte | ~string](v T) bool {
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
			b[i] = c ^ 0x20 // toggle case
		}
		nextA = 'A'
		if c == '-' {
			nextA = 'a'
		}
	}
	return string(b)
}
