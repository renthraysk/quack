package ascii

const (
	sp     = 1 << ' '
	htab   = 1 << '\t'
	lower  = ((1 << 26) - 1) << 'a'              // 'a'-'z'
	upper  = ((1 << 26) - 1) << 'A'              // 'A'-'Z'
	digits = ((1 << 10) - 1) << '0'              // '0'-'9'
	vchar  = ((1 << (1 + '~' - '!')) - 1) << '!' // '!'-'~'

	special = 1<<'!' | 1<<'#' | 1<<'$' | 1<<'%' | 1<<'&' | 1<<'\'' | 1<<'*' |
		1<<'+' | 1<<'-' | 1<<'.' | 1<<'^' | 1<<'_' | 1<<'`' | 1<<'|' |
		1<<'~' // !#$&'*+-.^_`|~

	// token is the set of characters allowed in pre HTTP/3 names.
	token = lower | upper | digits | special

	// token3 is the set of characters allowed in HTTP/3 incoming names.
	token3 = lower | digits | special

	// field-content  = field-vchar
	// [ 1*( SP / HTAB / field-vchar ) field-vchar ]
	// field-vchar    = VCHAR / obs-text
	// obs-text       = %x80-FF
	// https://www.rfc-editor.org/rfc/rfc9110#name-field-values
	fieldContent = sp | htab | vchar
)

// isUpper returns true if c is a upper case ASCII chararacter.
func isUpper(c byte) bool { return c-'A' <= 'Z'-'A' }

// ToLower returns the ASCII lowercase version of c.
func ToLower(c byte) byte {
	if isUpper(c) {
		c += 'a' - 'A'
	}
	return c
}

func isIn(c byte, lo, hi uint64) bool {
	m := lo
	if c >= 64 {
		m = hi
	}
	return (1<<(c%64))&m != 0
}

// isToken3Char returns true if byte c is a valid in a HTTP/3 name literal, false
// otherwise.
func isToken3Char(c byte) bool {
	return c < 0x80 && isIn(c, token3%(1<<64), token3>>64)
}

func IsNameValid(n []byte) bool {
	for _, c := range n {
		if !isToken3Char(c) {
			return false
		}
	}
	return true
}

// https://www.rfc-editor.org/rfc/rfc9110#section-5.5
func IsValueValid(v []byte) bool {
	if len(v) <= 0 {
		return false
	}
	// Has to start with a field-vchar
	// field-vchar    = VCHAR / obs-text
	// obs-text       = %x80-FF
	if v[0] < 0x80 && !isIn(v[0], vchar%(1<<64), vchar>>64) {
		return false
	}
	// Subsequent characters can include horizontal spaces.
	// field-content  = field-vchar
	// [ 1*( SP / HTAB / field-vchar ) field-vchar ]
	for _, c := range v[1:] {
		if c < 0x80 && !isIn(c, fieldContent%(1<<64), fieldContent>>64) {
			return false
		}
	}
	return true
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

// ToCanonical
func ToCanonical(b []byte) string {
	nextA := 'a'
	for i, c := range b {
		if c-byte(nextA) < 26 {
			// wrong cased letter
			b[i] = c ^ 0x20 // toggle case
		}
		nextA = 'A'
		if c == '-' {
			nextA = 'a'
		}
	}
	return string(b)
}

func AppendCanonical(p []byte, s string) []byte {
	i := len(p)
	p = append(p, s...)
	nextA := 'a'
	for ; i < len(p); i++ {
		c := p[i]
		if c-byte(nextA) < 26 {
			// wrong cased letter
			p[i] = c ^ 0x20 // toggle case
		}
		nextA = 'A'
		if c == '-' {
			nextA = 'a'
		}
	}
	return p
}
