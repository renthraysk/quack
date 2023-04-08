package huffman

var bounds = []struct {
	delta  uint32 //
	length uint8  // number of bits for the code
	offset uint8  // offset into longCodes for first sym that has this length
}{
	// {delta: 0x50000000, length: 5, offset: 0},
	// {delta: 0x68000000, length: 6, offset: 10},
	// {delta: 0x40000000, length: 7, offset: 36},
	// {delta: 0x06000000, length: 8, offset: 68},
	// {delta: 0x01400000, length: 10, offset: 74},
	// {delta: 0x00600000, length: 11, offset: 79},
	// {delta: 0x00200000, length: 12, offset: 82},
	// {delta: 0x00300000, length: 13, offset: 84},
	{delta: 0x00080000, length: 14, offset: 90},
	{delta: 0x00060000, length: 15, offset: 92},
	{delta: 0x00006000, length: 19, offset: 95},
	{delta: 0x00008000, length: 20, offset: 98},
	{delta: 0x00006800, length: 21, offset: 106},
	{delta: 0x00006800, length: 22, offset: 119},
	{delta: 0x00003a00, length: 23, offset: 145},
	{delta: 0x00000c00, length: 24, offset: 174},
	{delta: 0x00000200, length: 25, offset: 186},
	{delta: 0x000003c0, length: 26, offset: 190},
	{delta: 0x00000260, length: 27, offset: 205},
	{delta: 0x000001d0, length: 28, offset: 224},
}

var codes = [256]uint32{
	0x00001ff8, 0x007fffd8, 0x0fffffe2, 0x0fffffe3, 0x0fffffe4, 0x0fffffe5, 0x0fffffe6, 0x0fffffe7,
	0x0fffffe8, 0x00ffffea, 0x3ffffffc, 0x0fffffe9, 0x0fffffea, 0x3ffffffd, 0x0fffffeb, 0x0fffffec,
	0x0fffffed, 0x0fffffee, 0x0fffffef, 0x0ffffff0, 0x0ffffff1, 0x0ffffff2, 0x3ffffffe, 0x0ffffff3,
	0x0ffffff4, 0x0ffffff5, 0x0ffffff6, 0x0ffffff7, 0x0ffffff8, 0x0ffffff9, 0x0ffffffa, 0x0ffffffb,
	0x00000014, 0x000003f8, 0x000003f9, 0x00000ffa, 0x00001ff9, 0x00000015, 0x000000f8, 0x000007fa,
	0x000003fa, 0x000003fb, 0x000000f9, 0x000007fb, 0x000000fa, 0x00000016, 0x00000017, 0x00000018,
	0x00000000, 0x00000001, 0x00000002, 0x00000019, 0x0000001a, 0x0000001b, 0x0000001c, 0x0000001d,
	0x0000001e, 0x0000001f, 0x0000005c, 0x000000fb, 0x00007ffc, 0x00000020, 0x00000ffb, 0x000003fc,
	0x00001ffa, 0x00000021, 0x0000005d, 0x0000005e, 0x0000005f, 0x00000060, 0x00000061, 0x00000062,
	0x00000063, 0x00000064, 0x00000065, 0x00000066, 0x00000067, 0x00000068, 0x00000069, 0x0000006a,
	0x0000006b, 0x0000006c, 0x0000006d, 0x0000006e, 0x0000006f, 0x00000070, 0x00000071, 0x00000072,
	0x000000fc, 0x00000073, 0x000000fd, 0x00001ffb, 0x0007fff0, 0x00001ffc, 0x00003ffc, 0x00000022,
	0x00007ffd, 0x00000003, 0x00000023, 0x00000004, 0x00000024, 0x00000005, 0x00000025, 0x00000026,
	0x00000027, 0x00000006, 0x00000074, 0x00000075, 0x00000028, 0x00000029, 0x0000002a, 0x00000007,
	0x0000002b, 0x00000076, 0x0000002c, 0x00000008, 0x00000009, 0x0000002d, 0x00000077, 0x00000078,
	0x00000079, 0x0000007a, 0x0000007b, 0x00007ffe, 0x000007fc, 0x00003ffd, 0x00001ffd, 0x0ffffffc,
	0x000fffe6, 0x003fffd2, 0x000fffe7, 0x000fffe8, 0x003fffd3, 0x003fffd4, 0x003fffd5, 0x007fffd9,
	0x003fffd6, 0x007fffda, 0x007fffdb, 0x007fffdc, 0x007fffdd, 0x007fffde, 0x00ffffeb, 0x007fffdf,
	0x00ffffec, 0x00ffffed, 0x003fffd7, 0x007fffe0, 0x00ffffee, 0x007fffe1, 0x007fffe2, 0x007fffe3,
	0x007fffe4, 0x001fffdc, 0x003fffd8, 0x007fffe5, 0x003fffd9, 0x007fffe6, 0x007fffe7, 0x00ffffef,
	0x003fffda, 0x001fffdd, 0x000fffe9, 0x003fffdb, 0x003fffdc, 0x007fffe8, 0x007fffe9, 0x001fffde,
	0x007fffea, 0x003fffdd, 0x003fffde, 0x00fffff0, 0x001fffdf, 0x003fffdf, 0x007fffeb, 0x007fffec,
	0x001fffe0, 0x001fffe1, 0x003fffe0, 0x001fffe2, 0x007fffed, 0x003fffe1, 0x007fffee, 0x007fffef,
	0x000fffea, 0x003fffe2, 0x003fffe3, 0x003fffe4, 0x007ffff0, 0x003fffe5, 0x003fffe6, 0x007ffff1,
	0x03ffffe0, 0x03ffffe1, 0x000fffeb, 0x0007fff1, 0x003fffe7, 0x007ffff2, 0x003fffe8, 0x01ffffec,
	0x03ffffe2, 0x03ffffe3, 0x03ffffe4, 0x07ffffde, 0x07ffffdf, 0x03ffffe5, 0x00fffff1, 0x01ffffed,
	0x0007fff2, 0x001fffe3, 0x03ffffe6, 0x07ffffe0, 0x07ffffe1, 0x03ffffe7, 0x07ffffe2, 0x00fffff2,
	0x001fffe4, 0x001fffe5, 0x03ffffe8, 0x03ffffe9, 0x0ffffffd, 0x07ffffe3, 0x07ffffe4, 0x07ffffe5,
	0x000fffec, 0x00fffff3, 0x000fffed, 0x001fffe6, 0x003fffe9, 0x001fffe7, 0x001fffe8, 0x007ffff3,
	0x003fffea, 0x003fffeb, 0x01ffffee, 0x01ffffef, 0x00fffff4, 0x00fffff5, 0x03ffffea, 0x007ffff4,
	0x03ffffeb, 0x07ffffe6, 0x03ffffec, 0x03ffffed, 0x07ffffe7, 0x07ffffe8, 0x07ffffe9, 0x07ffffea,
	0x07ffffeb, 0x0ffffffe, 0x07ffffec, 0x07ffffed, 0x07ffffee, 0x07ffffef, 0x07fffff0, 0x03ffffee,
}

const codeLengths = "" +
	"\r\x17\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x18\x1e\x1c\x1c\x1e\x1c\x1c\x1c" +
	"\x1c\x1c\x1c\x1c\x1c\x1e\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x1c\x06" +
	"\n\n\f\r\x06\b\v\n\n\b\v\b\x06\x06\x06\x05\x05\x05\x06\x06\x06\x06" +
	"\x06\x06\x06\a\b\x0f\x06\f\n\r\x06\a\a\a\a\a\a\a\a\a\a\a\a\a\a\a" +
	"\a\a\a\a\a\a\a\b\a\b\r\x13\r\x0e\x06\x0f\x05\x06\x05\x06\x05\x06" +
	"\x06\x06\x05\a\a\x06\x06\x06\x05\x06\a\x06\x05\x05\x06\a\a\a\a\a" +
	"\x0f\v\x0e\r\x1c\x14\x16\x14\x14\x16\x16\x16\x17\x16\x17\x17\x17" +
	"\x17\x17\x18\x17\x18\x18\x16\x17\x18\x17\x17\x17\x17\x15\x16\x17" +
	"\x16\x17\x17\x18\x16\x15\x14\x16\x16\x17\x17\x15\x17\x16\x16\x18" +
	"\x15\x16\x17\x17\x15\x15\x16\x15\x17\x16\x17\x17\x14\x16\x16\x16" +
	"\x17\x16\x16\x17\x1a\x1a\x14\x13\x16\x17\x16\x19\x1a\x1a\x1a\x1b" +
	"\x1b\x1a\x18\x19\x13\x15\x1a\x1b\x1b\x1a\x1b\x18\x15\x15\x1a\x1a" +
	"\x1c\x1b\x1b\x1b\x14\x18\x14\x15\x16\x15\x15\x17\x16\x16\x19\x19" +
	"\x18\x18\x1a\x17\x1a\x1b\x1a\x1a\x1b\x1b\x1b\x1b\x1b\x1c\x1b\x1b" +
	"\x1b\x1b\x1b\x1a"

const shortCodes = "" +
	"0000000000000000000000000000000000000000000000000000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000000000" +
	"1111111111111111111111111111111111111111111111111111111111111111" +
	"1111111111111111111111111111111111111111111111111111111111111111" +
	"1111111111111111111111111111111111111111111111111111111111111111" +
	"1111111111111111111111111111111111111111111111111111111111111111" +
	"2222222222222222222222222222222222222222222222222222222222222222" +
	"2222222222222222222222222222222222222222222222222222222222222222" +
	"2222222222222222222222222222222222222222222222222222222222222222" +
	"2222222222222222222222222222222222222222222222222222222222222222" +
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
	"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
	"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc" +
	"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc" +
	"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc" +
	"cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc" +
	"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
	"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
	"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
	"eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee" +
	"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii" +
	"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii" +
	"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii" +
	"iiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiiii" +
	"oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo" +
	"oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo" +
	"oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo" +
	"oooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo" +
	"ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss" +
	"ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss" +
	"ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss" +
	"ssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssss" +
	"tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt" +
	"tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt" +
	"tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt" +
	"tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt" +
	"                                                                " +
	"                                                                " +
	"%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%" +
	"%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%" +
	"----------------------------------------------------------------" +
	"----------------------------------------------------------------" +
	"................................................................" +
	"................................................................" +
	"////////////////////////////////////////////////////////////////" +
	"////////////////////////////////////////////////////////////////" +
	"3333333333333333333333333333333333333333333333333333333333333333" +
	"3333333333333333333333333333333333333333333333333333333333333333" +
	"4444444444444444444444444444444444444444444444444444444444444444" +
	"4444444444444444444444444444444444444444444444444444444444444444" +
	"5555555555555555555555555555555555555555555555555555555555555555" +
	"5555555555555555555555555555555555555555555555555555555555555555" +
	"6666666666666666666666666666666666666666666666666666666666666666" +
	"6666666666666666666666666666666666666666666666666666666666666666" +
	"7777777777777777777777777777777777777777777777777777777777777777" +
	"7777777777777777777777777777777777777777777777777777777777777777" +
	"8888888888888888888888888888888888888888888888888888888888888888" +
	"8888888888888888888888888888888888888888888888888888888888888888" +
	"9999999999999999999999999999999999999999999999999999999999999999" +
	"9999999999999999999999999999999999999999999999999999999999999999" +
	"================================================================" +
	"================================================================" +
	"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
	"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
	"________________________________________________________________" +
	"________________________________________________________________" +
	"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" +
	"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb" +
	"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
	"dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd" +
	"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
	"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" +
	"gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg" +
	"gggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg" +
	"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh" +
	"hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh" +
	"llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll" +
	"llllllllllllllllllllllllllllllllllllllllllllllllllllllllllllllll" +
	"mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm" +
	"mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm" +
	"nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn" +
	"nnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnnn" +
	"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp" +
	"pppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppppp" +
	"rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr" +
	"rrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrr" +
	"uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu" +
	"uuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuuu" +
	"::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::" +
	"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB" +
	"CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC" +
	"DDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD" +
	"EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE" +
	"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" +
	"GGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGGG" +
	"HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH" +
	"IIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIIII" +
	"JJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJJ" +
	"KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK" +
	"LLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL" +
	"MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMMM" +
	"NNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN" +
	"OOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO" +
	"PPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPPP" +
	"QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ" +
	"RRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR" +
	"SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSS" +
	"TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT" +
	"UUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUU" +
	"VVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVV" +
	"WWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWWW" +
	"YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY" +
	"jjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjjj" +
	"kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk" +
	"qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq" +
	"vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv" +
	"wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww" +
	"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" +
	"yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy" +
	"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" +
	"&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&********************************" +
	",,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,,;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;" +
	"XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ" +
	"!!!!!!!!\"\"\"\"\"\"\"\"(((((((())))))))????????''''++++||||##>>" +
	"\x00$@[]~"

const longCodes = "" +
	"012aceiost %-./3456789=A_bdfghlmnpru:BCDEFGHIJKLMNOPQRSTUVWYjkqv" +
	"wxyz&*,;XZ!\"()?'+|#>\x00$@[]~^}<`{\\\xc3\xd0\x80\x82\x83\xa2\xb8" +
	"\xc2\xe0\xe2\x99\xa1\xa7\xac\xb0\xb1\xb3\xd1\xd8\xd9\xe3\xe5\xe6" +
	"\x81\x84\x85\x86\x88\x92\x9a\x9c\xa0\xa3\xa4\xa9\xaa\xad\xb2\xb5" +
	"\xb9\xba\xbb\xbd\xbe\xc4\xc6\xe4\xe8\xe9\x01\x87\x89\x8a\x8b\x8c" +
	"\x8d\x8f\x93\x95\x96\x97\x98\x9b\x9d\x9e\xa5\xa6\xa8\xae\xaf\xb4" +
	"\xb6\xb7\xbc\xbf\xc5\xe7\xef\t\x8e\x90\x91\x94\x9f\xab\xce\xd7\xe1" +
	"\xec\xed\xc7\xcf\xea\xeb\xc0\xc1\xc8\xc9\xca\xcd\xd2\xd5\xda\xdb" +
	"\xee\xf0\xf2\xf3\xff\xcb\xcc\xd3\xd4\xd6\xdd\xde\xdf\xf1\xf4\xf5" +
	"\xf6\xf7\xf8\xfa\xfb\xfc\xfd\xfe\x02\x03\x04\x05\x06\a\b\v\f\x0e" +
	"\x0f\x10\x11\x12\x13\x14\x15\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f" +
	"\x7f\xdc\xf9\n\r\x16"

var codes00To99 = [100]struct {
	length uint16
	code   uint16
}{
	{length: 10, code: 0x0000}, // 00
	{length: 10, code: 0x0001}, // 01
	{length: 10, code: 0x0002}, // 02
	{length: 11, code: 0x0019}, // 03
	{length: 11, code: 0x001a}, // 04
	{length: 11, code: 0x001b}, // 05
	{length: 11, code: 0x001c}, // 06
	{length: 11, code: 0x001d}, // 07
	{length: 11, code: 0x001e}, // 08
	{length: 11, code: 0x001f}, // 09
	{length: 10, code: 0x0020}, // 10
	{length: 10, code: 0x0021}, // 11
	{length: 10, code: 0x0022}, // 12
	{length: 11, code: 0x0059}, // 13
	{length: 11, code: 0x005a}, // 14
	{length: 11, code: 0x005b}, // 15
	{length: 11, code: 0x005c}, // 16
	{length: 11, code: 0x005d}, // 17
	{length: 11, code: 0x005e}, // 18
	{length: 11, code: 0x005f}, // 19
	{length: 10, code: 0x0040}, // 20
	{length: 10, code: 0x0041}, // 21
	{length: 10, code: 0x0042}, // 22
	{length: 11, code: 0x0099}, // 23
	{length: 11, code: 0x009a}, // 24
	{length: 11, code: 0x009b}, // 25
	{length: 11, code: 0x009c}, // 26
	{length: 11, code: 0x009d}, // 27
	{length: 11, code: 0x009e}, // 28
	{length: 11, code: 0x009f}, // 29
	{length: 11, code: 0x0320}, // 30
	{length: 11, code: 0x0321}, // 31
	{length: 11, code: 0x0322}, // 32
	{length: 12, code: 0x0659}, // 33
	{length: 12, code: 0x065a}, // 34
	{length: 12, code: 0x065b}, // 35
	{length: 12, code: 0x065c}, // 36
	{length: 12, code: 0x065d}, // 37
	{length: 12, code: 0x065e}, // 38
	{length: 12, code: 0x065f}, // 39
	{length: 11, code: 0x0340}, // 40
	{length: 11, code: 0x0341}, // 41
	{length: 11, code: 0x0342}, // 42
	{length: 12, code: 0x0699}, // 43
	{length: 12, code: 0x069a}, // 44
	{length: 12, code: 0x069b}, // 45
	{length: 12, code: 0x069c}, // 46
	{length: 12, code: 0x069d}, // 47
	{length: 12, code: 0x069e}, // 48
	{length: 12, code: 0x069f}, // 49
	{length: 11, code: 0x0360}, // 50
	{length: 11, code: 0x0361}, // 51
	{length: 11, code: 0x0362}, // 52
	{length: 12, code: 0x06d9}, // 53
	{length: 12, code: 0x06da}, // 54
	{length: 12, code: 0x06db}, // 55
	{length: 12, code: 0x06dc}, // 56
	{length: 12, code: 0x06dd}, // 57
	{length: 12, code: 0x06de}, // 58
	{length: 12, code: 0x06df}, // 59
	{length: 11, code: 0x0380}, // 60
	{length: 11, code: 0x0381}, // 61
	{length: 11, code: 0x0382}, // 62
	{length: 12, code: 0x0719}, // 63
	{length: 12, code: 0x071a}, // 64
	{length: 12, code: 0x071b}, // 65
	{length: 12, code: 0x071c}, // 66
	{length: 12, code: 0x071d}, // 67
	{length: 12, code: 0x071e}, // 68
	{length: 12, code: 0x071f}, // 69
	{length: 11, code: 0x03a0}, // 70
	{length: 11, code: 0x03a1}, // 71
	{length: 11, code: 0x03a2}, // 72
	{length: 12, code: 0x0759}, // 73
	{length: 12, code: 0x075a}, // 74
	{length: 12, code: 0x075b}, // 75
	{length: 12, code: 0x075c}, // 76
	{length: 12, code: 0x075d}, // 77
	{length: 12, code: 0x075e}, // 78
	{length: 12, code: 0x075f}, // 79
	{length: 11, code: 0x03c0}, // 80
	{length: 11, code: 0x03c1}, // 81
	{length: 11, code: 0x03c2}, // 82
	{length: 12, code: 0x0799}, // 83
	{length: 12, code: 0x079a}, // 84
	{length: 12, code: 0x079b}, // 85
	{length: 12, code: 0x079c}, // 86
	{length: 12, code: 0x079d}, // 87
	{length: 12, code: 0x079e}, // 88
	{length: 12, code: 0x079f}, // 89
	{length: 11, code: 0x03e0}, // 90
	{length: 11, code: 0x03e1}, // 91
	{length: 11, code: 0x03e2}, // 92
	{length: 12, code: 0x07d9}, // 93
	{length: 12, code: 0x07da}, // 94
	{length: 12, code: 0x07db}, // 95
	{length: 12, code: 0x07dc}, // 96
	{length: 12, code: 0x07dd}, // 97
	{length: 12, code: 0x07de}, // 98
	{length: 12, code: 0x07df}, // 99
}
