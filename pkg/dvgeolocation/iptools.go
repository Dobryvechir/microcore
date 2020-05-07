/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvgeolocation

import (
	"strconv"
	"strings"
)

func ReadStringToBuf(buf []byte, size int, ipStr string) bool {
	for i := 0; i < size; i++ {
		buf[i] = 0
	}
	n := len(ipStr)
	for j := 0; j < n; j++ {
		var carry int = int(ipStr[j]) - 48
		if carry < 0 || carry > 9 {
			return false
		}
		for i := 0; i < size; i++ {
			var val int = int(buf[i])*10 + carry
			buf[i] = byte(val & 255)
			carry = val >> 8
		}
		if carry > 0 {
			return false
		}
	}
	return true
}

func ReadString(ip string) (ipBuf []byte, ok bool) {
	ipBuf = make([]byte, 16)
	ok = ReadStringToBuf(ipBuf, 16, ip)
	return
}

func ReadIP4ToBuf(ip string, ipBuf []byte, pos int) bool {
	n := len(ip)
	p := 0
	for i := 0; i < 4; i++ {
		if p >= n {
			return false
		}
		v := 0
		for ; p < n && ip[p] != '.'; p++ {
			c := int(ip[p]) - 48
			v = v*10 + c
			if c < 0 || c > 9 || v > 255 {
				return false
			}
		}
		ipBuf[pos+3-i] = byte(v)
		if i < 3 {
			p++
		}
	}
	return p == n
}

func ReadIP4(ip string) (ipBuf []byte, ok bool) {
	ipBuf = make([]byte, 4)
	ok = ReadIP4ToBuf(ip, ipBuf, 0)
	return
}

func ReadIP(ip string) (ipBuf []byte, ok bool) {
	ipBuf, ok = ReadIP4(ip)
	if ok {
		return
	}
	ipBuf, ok = ReadIP6(ip)
	return
}

func ReadIPOrString(src string) (ipBuf []byte, ok bool) {
	ipBuf, ok = ReadIP(src)
	if !ok {
		ipBuf, ok = ReadString(src)
	}
	return
}
func WriteIP4(ipBuf []byte) string {
	if len(ipBuf) < 4 {
		return ""
	}
	return strconv.Itoa(int(ipBuf[3])) + "." + strconv.Itoa(int(ipBuf[2])) + "." + strconv.Itoa(int(ipBuf[1])) + "." + strconv.Itoa(int(ipBuf[0]))
}

func readIP6Internal(ip string, n int, countColon int) (ipBuf []byte, ok bool) {
	ipBuf = make([]byte, 16)
	pos := 0
	ok = false
	doubleColon := false
	for i := 7; i >= 0; i-- {
		if pos >= n {
			return
		}
		c := ip[pos]
		if c == ':' {
			leftColonsAhead := countColon - (8 - i)
			pos++
			if i == 7 {
				if pos >= n {
					return
				}
				c = ip[pos]
				pos++
				if c != ':' {
					return
				}
				leftColonsAhead--
			}
			if leftColonsAhead != 0 || pos < n {
				leftColonsAhead++
			}
			if doubleColon || i <= leftColonsAhead {
				return
			}
			doubleColon = true
			i = leftColonsAhead
		} else {
			w := 0
			for j := 0; pos < n && j < 4 && ip[pos] != ':'; j++ {
				d := int(ip[pos]) - 48
				pos++
				if d > 9 {
					d = (d - 7) & 15
				}
				w = w<<4 | d
			}
			ipBuf[i<<1|1] = byte(w >> 8)
			ipBuf[i<<1] = byte(w & 255)
			if i != 0 {
				if pos >= n {
					return
				}
				c = ip[pos]
				pos++
				if c != ':' {
					return
				}
			}
		}
	}
	ok = pos == n
	return
}

func ReadIP6(ip string) (ipBuf []byte, ok bool) {
	ok = false
	n := len(ip)
	countColon := 0
	countPeriod := 0
	for i := 0; i < n; i++ {
		c := ip[i]
		if c == ':' {
			if countPeriod != 0 {
				return
			}
			countColon++
		} else if c == '.' {
			if countColon < 3 {
				return
			}
			countPeriod++
		} else if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return
		}
	}
	if countPeriod == 3 {
		leftPos := strings.LastIndex(ip, ":") + 1
		ipRight := ip[leftPos:]
		ipLeft := ip[:leftPos] + "0:0"
		ipBuf, ok = readIP6Internal(ipLeft, leftPos+3, countColon+1)
		if ok {
			ok = IsIP6UsedAsIP4Container(ipBuf) && ReadIP4ToBuf(ipRight, ipBuf, 0)
		}
		return
	}
	if countColon < 2 || countPeriod != 0 {
		return
	}
	ipBuf, ok = readIP6Internal(ip, n, countColon)
	return
}

func IsIP6UsedAsIP4Container(ipBuf []byte) bool {
	return len(ipBuf) >= 16 && ipBuf[15] == 0 && ipBuf[14] == 0 && ipBuf[13] == 0 &&
		ipBuf[12] == 0 && ipBuf[11] == 0 && ipBuf[10] == 0 && ipBuf[9] == 0 && ipBuf[8] == 0 &&
		ipBuf[7] == 0 && ipBuf[6] == 0 && ipBuf[5] == 255 && ipBuf[4] == 255
}

func WriteIP6(ipBuf []byte) string {
	if len(ipBuf) < 16 || IsIP6UsedAsIP4Container(ipBuf) {
		if len(ipBuf) >= 4 {
			return "::ffff:" + WriteIP4(ipBuf)
		}
		return ""
	}
	zeroOmit := 0
	res := make([]byte, 63)
	pos := 0
	for i := 7; i >= 0; i-- {
		w := int(ipBuf[i<<1|1])<<8 | int(ipBuf[i<<1])
		if w == 0 {
			if zeroOmit == 0 {
				zeroOmit = 1
				res[pos] = ':'
				pos++
				res[pos] = ':'
				pos++
			} else if zeroOmit != 1 {
				res[pos] = ':'
				pos++
				res[pos] = '0'
				pos++
			}
		} else {
			if zeroOmit == 1 {
				zeroOmit = 2
			} else {
				if i != 7 {
					res[pos] = ':'
					pos++
				}
			}
			start := true
			for j := 3; j >= 0; j-- {
				d := w>>uint32(j<<2)&15 | 48
				if d == '0' && start {
					continue
				}
				start = false
				if d > '9' {
					d += 39
				}
				res[pos] = byte(d)
				pos++
			}
		}
	}
	return string(res[:pos])
}

func WriteIP(ipBuf []byte, onlyOne bool) string {
	if len(ipBuf) < 16 {
		return WriteIP4(ipBuf)
	}
	isLow := true
	for i := 4; i < 16; i++ {
		if ipBuf[i] != 0 {
			isLow = false
			break
		}
	}
	if !isLow {
		return WriteIP6(ipBuf)
	}
	if onlyOne {
		return WriteIP4(ipBuf)
	}
	return WriteIP4(ipBuf) + " (" + WriteIP6(ipBuf) + ")"
}

func WriteBufAsHex(ipBuf []byte) string {
	last := len(ipBuf) - 1
	if last > 15 {
		last = 15
	}
	for ; last > 0 && ipBuf[last] == 0; last-- {
	}
	scope := (last + 1) * 2
	if ipBuf[last] < 16 {
		scope--
	}
	res := make([]byte, scope)
	pos := 0
	for ; last >= 0; last-- {
		v := int(ipBuf[last])
		if pos > 0 || v >= 16 {
			k := (v >> 4) + 48
			if k > 57 {
				k += 7
			}
			res[pos] = byte(k)
			pos++
		}
		v = (v & 15) + 48
		if v > 57 {
			v += 7
		}
		res[pos] = byte(v)
		pos++
	}
	return string(res)
}

func WriteBufAsNumber(ipBuf []byte) string {
	last := len(ipBuf) - 1
	if last > 15 {
		last = 15
	}
	v0 := uint64(0)
	vLast := last
	if vLast > 7 {
		vLast = 7
	}
	for ; vLast >= 0; vLast-- {
		v0 |= uint64(ipBuf[vLast]) << uint32(vLast<<3)
	}
	if last < 8 {
		return strconv.FormatUint(v0, 10)
	}
	v1 := uint64(0)
	scope := (last + 1) * 3
	for ; last >= 8; last-- {
		v1 |= uint64(ipBuf[last]) << uint32((last&7)<<3)
	}
	res := make([]byte, scope)
	for v1 != 0 || v0 != 0 {
		r := v1 % 10
		v1 = v1 / 10
		if r != 0 {
			d0 := v0 & uint64(0xFFFFFFFF)
			d1 := (v0 >> 32) | (r << 32)
			r = d1 % 10
			d1 = d1 / 10

			d0 |= r << 32
			r = d0 % 10
			d0 = d0 / 10
			v0 = (d1 << 32) | d0

		} else {
			r = v0 % 10
			v0 = v0 / 10
		}
		scope--
		res[scope] = byte(r + 48)
	}
	return string(res[scope:])
}

func WriteBufInAllPresentations(data []byte) string {
	return WriteIP(data, false) + " = " + WriteBufAsNumber(data) + " (0x" +
		WriteBufAsHex(data) + ")"
}
