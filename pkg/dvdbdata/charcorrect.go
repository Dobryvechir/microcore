// package dvdbdata provides functions for sql query
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbdata

var extraC2DF = [][]byte{{0x80, 0xbf}}
var extraE0 = [][]byte{{0xa0, 0xbf}, {0x80, 0xbf}}
var extraE1EC = [][]byte{{0x80, 0xbf}, {0x80, 0xbf}}
var extraED = [][]byte{{0x80, 0x9f}, {0x80, 0xbf}}
var extraEEEF = extraE1EC
var extraF0 = [][]byte{{0x90, 0xbf}, {0x80, 0xbf}, {0x80, 0xbf}}
var extraF1F3 = [][]byte{{0x80, 0xbf}, {0x80, 0xbf}, {0x80, 0xbf}}
var extraF4 = [][]byte{{0x80, 0x8f}, {0x80, 0xbf}, {0x80, 0xbf}}
var extraMap = map[byte][][]byte{
	0xc2: extraC2DF,
	0xc3: extraC2DF,
	0xc4: extraC2DF,
	0xc5: extraC2DF,
	0xc6: extraC2DF,
	0xc7: extraC2DF,
	0xc8: extraC2DF,
	0xc9: extraC2DF,
	0xca: extraC2DF,
	0xcb: extraC2DF,
	0xcc: extraC2DF,
	0xcd: extraC2DF,
	0xce: extraC2DF,
	0xcf: extraC2DF,
	0xd0: extraC2DF,
	0xd1: extraC2DF,
	0xd2: extraC2DF,
	0xd3: extraC2DF,
	0xd4: extraC2DF,
	0xd5: extraC2DF,
	0xd6: extraC2DF,
	0xd7: extraC2DF,
	0xd8: extraC2DF,
	0xd9: extraC2DF,
	0xda: extraC2DF,
	0xdb: extraC2DF,
	0xdc: extraC2DF,
	0xdd: extraC2DF,
	0xde: extraC2DF,
	0xdf: extraC2DF,
	0xe0: extraE0,
	0xe1: extraE1EC,
	0xe2: extraE1EC,
	0xe3: extraE1EC,
	0xe4: extraE1EC,
	0xe5: extraE1EC,
	0xe6: extraE1EC,
	0xe7: extraE1EC,
	0xe8: extraE1EC,
	0xe9: extraE1EC,
	0xea: extraE1EC,
	0xeb: extraE1EC,
	0xec: extraE1EC,
	0xed: extraED,
	0xee: extraEEEF,
	0xef: extraEEEF,
	0xf0: extraF0,
	0xf1: extraF1F3,
	0xf2: extraF1F3,
	0xf3: extraF1F3,
	0xf4: extraF4,
}
var replacementMap = map[byte]byte{
	0x92: 0x60,
}

func GetHexValue(b byte) string {
	b1 := (b & 0xf) | 0x30
	b2 := ((b >> 4) & 0xf) | 0x30
	if b1 >= 0x3a {
		b1 += 7
	}
	if b2 >= 0x3a {
		b2 += 7
	}
	return string([]byte{b2, b1})
}

func CheckReplaceNonUtf8Characters(b []byte, onlyCheck bool) ([]byte, string) {
	res := ""
	n := len(b)
	for i := 0; i < n; i++ {
		c := b[i]
		if c >= 0x80 {
			d := extraMap[c]
			m := len(d)
			replace := m == 0 || i+m >= n
			a:=c
			if !replace {
				for j := 0; j < m; j++ {
					a = b[i+1+j]
					if a < d[j][0] || a > d[j][1] {
						replace = true
						break
					}
				}
			}
			if replace {
				res += " 0x" + GetHexValue(a)
				bn, ok := replacementMap[c]
				if !ok {
					bn = ' '
				}
				if !onlyCheck {
					b[i] = bn
				}
			} else {
				i += m
			}
		}
	}
	if res != "" {
		res = "(Invalid " + res + " in " + string(b) + ")"
	}
	return b, res
}

func CheckReplaceNonUtf8CharactersInStringArray(m []string, onlyCheck bool) string {
	n := len(m)
	res := ""
	for i := 0; i < n; i++ {
		b, r := CheckReplaceNonUtf8Characters([]byte(m[i]), onlyCheck)
		if r != "" {
			if !onlyCheck {
				m[i] = string(b)
			}
			res += r
		}
	}
	return res
}

func CheckReplaceNonUtf8CharactersInStringArrayOfArray(m [][]string, onlyCheck bool) string {
	n := len(m)
	res := ""
	for i := 0; i < n; i++ {
		r := CheckReplaceNonUtf8CharactersInStringArray(m[i], onlyCheck)
		if r != "" {
			res += r
		}
	}
	return res
}
