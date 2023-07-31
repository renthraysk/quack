//go:build ignore

package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Value struct {
	index  uint64
	values map[string]uint64
}

func intern(ss []string) (string, map[string]int) {
	var intern strings.Builder
	pos := make(map[string]int, len(ss))

	sort.Slice(ss, func(i, j int) bool {
		if len(ss[j]) < len(ss[i]) {
			return true
		}
		return len(ss[i]) == len(ss[j]) && ss[i] < ss[j]
	})

	for _, s := range ss {
		if i := strings.Index(intern.String(), s); i >= 0 {
			pos[s] = i
			continue
		}
		pos[s] = intern.Len()
		intern.WriteString(s)
	}
	return intern.String(), pos
}

func main() {
	r := strings.NewReader(http3)
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	ss := make([]string, 0, 100)
	for row, err := cr.Read(); err == nil; row, err = cr.Read() {
		if len(row) < 3 {
			continue
		}
		ss = append(ss, http.CanonicalHeaderKey(row[1]))
		if len(row[2]) > 0 {
			ss = append(ss, row[2])
		}
	}

	intern, pos := intern(ss)

	a := &strings.Builder{}
	r.Seek(0, io.SeekStart)
	cr = csv.NewReader(r)
	cr.Comma = '\t'

	nameValues := make(map[string]Value, 100)
	for row, err := cr.Read(); err == nil; row, err = cr.Read() {
		if len(row) < 3 {
			continue
		}

		index, _ := strconv.ParseUint(row[0], 10, 32)
		name := http.CanonicalHeaderKey(row[1])
		namePos := pos[name]
		if len(row[2]) > 0 {
			value := row[2]
			valuePos := pos[value]
			fmt.Fprintf(a, "\t{name: intern[%d:%d], value: intern[%d:%d]}, // %d %s: %s\n",
				namePos, namePos+len(name),
				valuePos, valuePos+len(value),
				index,
				name, value)

			if _, ok := nameValues[name]; !ok {
				nameValues[name] = Value{index, make(map[string]uint64)}
			}
			nameValues[name].values[value] = index
		} else {
			fmt.Fprintf(a, "\t{name: intern[%d:%d]}, // %d %s \n",
				namePos, namePos+len(name),
				index, name)
			nameValues[name] = Value{index, make(map[string]uint64)}
			nameValues[name].values[""] = index
		}
	}

	w := os.Stdout

	fmt.Fprint(w, "package field\n\n")

	if true {

		type nameValue struct {
			name  string
			value Value
		}

		ordered := make([]nameValue, 0, len(nameValues))

		for key, value := range nameValues {
			ordered = append(ordered, nameValue{name: key, value: value})
		}
		sort.Slice(ordered, func(i, j int) bool { return ordered[i].name < ordered[j].name })

		fmt.Fprintln(w, `func staticLookup(name, value string) (index uint64, m match) {`)
		fmt.Fprintln(w, "\tswitch name {")
		for _, v := range ordered {
			key := v.name
			v := v.value

			if key[0] == ':' {
				continue
			}

			fmt.Fprintf(w, "\tcase %q:\n", key)

			type vs struct {
				index uint64
				value string
			}

			values := make([]vs, 0, 8)

			for value, index := range v.values {
				values = append(values, vs{index, value})
			}

			switch len(values) {
			case 1:
				fmt.Fprintf(w, "\t\tif value == %q {\n", values[0].value)
				fmt.Fprintf(w, "\t\t\treturn %d, matchNameValue\n", values[0].index)
			default:
				fmt.Fprintf(w, "\t\tswitch value {\n")
				for _, vs := range values {
					fmt.Fprintf(w, "\t\t\tcase %q:\n", vs.value)
					fmt.Fprintf(w, "\t\t\t\treturn %d, matchNameValue\n", vs.index)
				}
			}
			fmt.Fprintf(w, "\t\t}\n")
			fmt.Fprintf(w, "\t\treturn %d, matchName\n", values[0].index)
		}
		fmt.Fprintln(w, "\t}")
		fmt.Fprintln(w, "\treturn 0, matchNone")
		fmt.Fprintln(w, `}`)
	}
	fmt.Fprintln(w)
	fmt.Fprint(w, "const intern string = \"\"+\n")
	printGoString(w, intern, 70)
	fmt.Fprintln(w)
	fmt.Fprintln(w, `var staticTable = [...]header{`)
	fmt.Fprint(w, a.String())
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}

func isIn(c byte, lo, hi uint64) bool {
	if c > 64 {
		lo = hi
	}
	return (1 << (c % 64) & lo) != 0
}

func writeByte(w io.ByteWriter, c byte) error {
	const hex = "0123456789abcdef"
	const escChars = (1 << 0x20) - 1 | 1<<'"' | 1<<'\\' | 1<<'\x7F'

	if c > '~' || isIn(c, escChars%(1<<64), escChars>>64) {
		w.WriteByte('\\')
		switch c {
		case '\a':
			c = 'a'
		case '\b':
			c = 'b'
		case '\t':
			c = 't'
		case '\n':
			c = 'n'
		case '\v':
			c = 'v'
		case '\f':
			c = 'f'
		case '\r':
			c = 'r'
		case '"', '\\':
		default:
			w.WriteByte('x')
			w.WriteByte(hex[c>>4])
			c = hex[c&0xF]
		}
	}
	return w.WriteByte(c)
}

func printGoString[T ~string | ~[]byte](w io.Writer, s T, width int) {

	const tabWidth = 4

	if len(s) == 0 {
		return
	}
	var b bytes.Buffer

	ww := width - tabWidth
	for i := 0; i < len(s); {
		b.Reset()
		b.WriteString("\t\"")
		for ; i < len(s) && b.Len() < ww; i++ {
			writeByte(&b, s[i])
		}
		eol := "\" +\n"
		if i >= len(s) {
			eol = "\"\n"
		}
		b.WriteString(eol)
		b.WriteTo(w)
	}
}

const http3 = `0	:authority	
1	:path	/
2	age	0
3	content-disposition	
4	content-length	0
5	cookie	
6	date	
7	etag	
8	if-modified-since	
9	if-none-match	
10	last-modified	
11	link	
12	location	
13	referer	
14	set-cookie	
15	:method	CONNECT
16	:method	DELETE
17	:method	GET
18	:method	HEAD
19	:method	OPTIONS
20	:method	POST
21	:method	PUT
22	:scheme	http
23	:scheme	https
24	:status	103
25	:status	200
26	:status	304
27	:status	404
28	:status	503
29	accept	*/*
30	accept	application/dns-message
31	accept-encoding	gzip, deflate, br
32	accept-ranges	bytes
33	access-control-allow-headers	cache-control
34	access-control-allow-headers	content-type
35	access-control-allow-origin	*
36	cache-control	max-age=0
37	cache-control	max-age=2592000
38	cache-control	max-age=604800
39	cache-control	no-cache
40	cache-control	no-store
41	cache-control	public, max-age=31536000
42	content-encoding	br
43	content-encoding	gzip
44	content-type	application/dns-message
45	content-type	application/javascript
46	content-type	application/json
47	content-type	application/x-www-form-urlencoded
48	content-type	image/gif
49	content-type	image/jpeg
50	content-type	image/png
51	content-type	text/css
52	content-type	text/html; charset=utf-8
53	content-type	text/plain
54	content-type	text/plain;charset=utf-8
55	range	bytes=0-
56	strict-transport-security	max-age=31536000
57	strict-transport-security	max-age=31536000; includesubdomains
58	strict-transport-security	max-age=31536000; includesubdomains; preload
59	vary	accept-encoding
60	vary	origin
61	x-content-type-options	nosniff
62	x-xss-protection	1; mode=block
63	:status	100
64	:status	204
65	:status	206
66	:status	302
67	:status	400
68	:status	403
69	:status	421
70	:status	425
71	:status	500
72	accept-language	
73	access-control-allow-credentials	FALSE
74	access-control-allow-credentials	TRUE
75	access-control-allow-headers	*
76	access-control-allow-methods	get
77	access-control-allow-methods	get, post, options
78	access-control-allow-methods	options
79	access-control-expose-headers	content-length
80	access-control-request-headers	content-type
81	access-control-request-method	get
82	access-control-request-method	post
83	alt-svc	clear
84	authorization	
85	content-security-policy	script-src 'none'; object-src 'none'; base-uri 'none'
86	early-data	1
87	expect-ct	
88	forwarded	
89	if-range	
90	origin	
91	purpose	prefetch
92	server	
93	timing-allow-origin	*
94	upgrade-insecure-requests	1
95	user-agent	
96	x-forwarded-for	
97	x-frame-options	deny
98	x-frame-options	sameorigin
`
