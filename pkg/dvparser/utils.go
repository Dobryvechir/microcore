/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvparser

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
