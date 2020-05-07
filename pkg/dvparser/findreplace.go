/***********************************************************************
MicroCore
Copyright 2020 -2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvparser

func FindSubStringSmartInByteArray(data []byte, needle []byte) (int, int) {
	needleLen := len(needle)
	n := len(data)
	ns := n - needleLen + 1
	c := needle[0]
bigLoop:
	for i := 0; i < ns; i++ {
		if data[i] == c {
			p := i
			for j := 0; j < needleLen; j++ {
				b := needle[j]
				if b == '\n' {
					for p+j < n && data[p+j] != b {
						if data[p+j] <= ' ' {
							p++
						} else {
							continue bigLoop
						}
					}
					for j < needleLen && (needle[j] == 10 || needle[j] == 13) {
						j++
						p--
					}
					for p+j < n && (data[p+j] == 10 || data[p+j] == 13) {
						p++
					}
					if j >= needleLen-1 {
						if j >= needleLen {
							p += j - needleLen + 1
						}
						continue
					}
					for j < needleLen && needle[j] <= ' ' {
						j++
						p--
					}
					for p+j < n && data[p+j] <= ' ' {
						p++
					}
					j--
				} else if p+j >= n || data[p+j] != b {
					continue bigLoop
				}
			}
			return i, p + needleLen
		}
	}
	return -1, -1
}

func ExtractFromBufByBeforeAfterKeys(data []byte, before []byte, after []byte) (res []byte, status int) {
	_, pos := FindSubStringSmartInByteArray(data, before)
	if pos < 0 {
		return data, -1
	}
	data = data[pos:]
	pos, _ = FindSubStringSmartInByteArray(data, after)
	if pos < 0 {
		return data, -2
	}
	return data[:pos], 0
}

func ReplaceTextInsideByteArray(data []byte, start int, end int, replacement []byte) []byte {
	r := len(replacement)
	n := len(data)
	if start < 0 {
		start = 0
	}
	if end > n {
		end = n
	}
	m := end - start
	if m < 0 {
		return data
	}
	if r > m {
		return append(append(data[:start:start], replacement...), data[end:]...)
	}
	for i := 0; i < r; i++ {
		data[start+i] = replacement[i]
	}
	dif := m - r
	if dif > 0 {
		n -= dif
		for i := start + r; i < n; i++ {
			data[i] = data[i+dif]
		}
		data = data[:n]
	}
	return data
}
