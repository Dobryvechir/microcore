/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

var JsonEscapeTable = map[byte]byte{
	'\\':     '\\',
	'"':      '"',
	byte(8):  'b',
	byte(9):  't',
	byte(10): 'r',
	byte(12): 'f',
	byte(13): 'n',
}

func ReadHexValue(v string) int64 {
	var r int64 = 0
	n := len(v)
	i := 0
	for i < n && v[i] <= ' ' {
		i++
	}
	for ; i < n; i++ {
		c := v[i]
		if c >= '0' && c <= '9' {
			c -= '0'
		} else if c >= 'a' && c <= 'f' {
			c -= 'a' - 10
		} else if c >= 'A' && c <= 'F' {
			c -= 'A' - 10
		} else {
			break
		}
		r = (r << 4) | int64(c)
	}
	return r
}

func Int64ToFullHex(v int64) string {
	b := make([]byte, 16)
	for i := 15; i >= 0; i-- {
		c := byte((v & 0xf) | 0x30)
		if c > '9' {
			c += 7
		}
		b[i] = c
		v >>= 4
	}
	return string(b)
}

func GetVersionIndex(version string) int64 {
	var r int64 = 0
	var m int64
	i := 0
	n := len(version)
	for k := 3; k >= 0; k-- {
		m = 0
		for ; i < n; i++ {
			c := version[i]
			if c >= '0' && c <= '9' {
				break
			}
		}
		if i >= n {
			break
		}
		for ; i < n && version[i] >= '0' && version[i] <= '9'; i++ {
			m = m*10 + int64(version[i]-48)
		}
		r |= m << uint32(k<<4)
	}
	return r
}

func PrintInt64ToByteBuffer(n int64, b []byte) int {
	k := 0
	if n < 0 {
		b[k] = '-'
		k++
	}
	if n < 10 {
		b[k] = byte(n + 48)
		k++
	} else {
		m := k
		for n > 0 {
			b[k] = byte(n % 10)
			k++
			n /= 10
		}
		i := k - 1
		for m < i {
			j := b[m]
			b[m] = b[i]
			b[i] = j
			m++
			i--
		}
	}
	return k
}

func GetCanonicalVersion(version int64) string {
	b := make([]byte, 6*4-1)
	k := 0
	for i := 3; i >= 0; i-- {
		r := (version >> uint32(i<<4)) & 0xffff
		k += PrintInt64ToByteBuffer(r, b[k:])
		if i > 0 {
			b[k] = '.'
			k++
		}
	}
	return string(b[:k])
}

func GetCanonicalVersionFromHexName(version string) string {
	return GetCanonicalVersion(ReadHexValue(version))
}

func IsAlphabeticalLowCase(s string) bool {
	n := len(s)
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		if !(s[i] >= 'a' && s[i] <= 'z') {
			return false
		}
	}
	return true
}

func IsDigitOnly(s string) bool {
	n := len(s)
	for i := 0; i < n; i++ {
		if !(s[i] >= '0' && s[i] <= '9') {
			return false
		}
	}
	return true
}

func IsSignAndDigitsOnly(s string) bool {
	n := len(s)
	i := 0
	if n > 0 && (s[0] == '-' || s[0] == '+') {
		i++
	}
	for ; i < n; i++ {
		if !(s[i] >= '0' && s[i] <= '9') {
			return false
		}
	}
	return true
}

func ConvertToUpperAlphaDigital(b []byte) string {
	n := len(b)
	for i := 0; i < n; i++ {
		c := b[i]
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		} else if !(c >= 'A' && c <= 'Z' || c >= '0' && c <= '9') {
			b[i] = '_'
		}
	}
	return string(b)
}

func IsUpperAlphaDigitalBytes(b []byte) bool {
	n := len(b)
	for i := 0; i < n; i++ {
		c := b[i]
		if !(c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' && i != 0 || c == '_') {
			return false
		}
	}
	return true
}

func IsUpperAlphaDigital(s string) bool {
	return IsUpperAlphaDigitalBytes([]byte(s))
}

func GetKeysFromStringIntMap(data map[string]int) []string {
	i := 0
	r := make([]string, len(data))
	for k, _ := range data {
		r[i] = k
		i++
	}
	return r
}

func IsJsonNumber(d []byte) bool {
	n := len(d)
	pnt := 0
	if n > 0 && (d[0] == '+' || d[0] == '-') {
		pnt++
	}
	if pnt >= n || n >= 24 {
		return false
	}
	b := d[pnt]
	if !(b >= '0' && b <= '9') {
		return false
	}
	for pnt < n && d[pnt] >= '0' && d[pnt] <= '9' {
		pnt++
	}
	if pnt == n {
		return true
	}
	if d[pnt] == '.' {
		for pnt++; pnt < n && d[pnt] >= '0' && d[pnt] <= '9'; pnt++ {
		}
		if pnt == n {
			return true
		}
	}
	b = d[pnt]
	pnt++
	if b != 'e' && b != 'E' || pnt == n {
		return false
	}
	b = d[pnt]
	pnt++
	if b == '+' || b == '-' {
		if pnt == n {
			return false
		}
		b = d[pnt]
		pnt++
	}
	if !(b >= '0' && b <= '9') || n-pnt > 2 {
		return false
	}
	for ; pnt < n; pnt++ {
		b = d[pnt]
		if !(b >= '0' && b <= '9') {
			return false
		}
	}
	return true
}

func IsJsonString(str []byte) bool {
	n := len(str)
	if n < 2 || str[0] != str[n-1] {
		return false
	}
	return true
}

func QuoteEscapedJsonBytes(b []byte) []byte {
	n := len(b)
	dstLen := n + 2
	res := make([]byte, dstLen)
	dst := 1
	res[0] = '"'
	for i := 0; i < n; i++ {
		d := b[i]
		if c, ok := JsonEscapeTable[d]; ok {
			if dst < dstLen {
				res[dst] = '\\'
				dst++
			} else {
				res = append(res, '\\')
			}
			d = c
		} else if d < ' ' {
			d = ' '
		}
		if dst < dstLen {
			res[dst] = d
			dst++
		} else {
			res = append(res, d)
		}
	}
	if dst < dstLen {
		res[dst] = '"'
		dst++
		if dst != dstLen {
			res = res[:dst]
		}
	} else {
		res = append(res, '"')
	}
	return res
}

func QuoteEscapedJsonBytesToString(b []byte) string {
	return string(QuoteEscapedJsonBytes(b))
}

func QuoteEscapedJsonString(s string) string {
	return QuoteEscapedJsonBytesToString([]byte(s))
}
