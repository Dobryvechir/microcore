/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import (
	"bytes"
	"io/ioutil"
)

func placeReadyItemToMap(f []string, m map[string][][]string, p int) {
	key := f[0]
	h := make([]string, p-1)
	for j := 1; j < p; j++ {
		h[j-1] = f[j]
	}
	m[key] = append(m[key], h)
}

func countCharEntries(data []byte, n int, b byte) int {
	c := 0
	for i := 0; i < n; i++ {
		if data[i] == b {
			c++
		}
	}
	return c
}

func DetectColumnSeparator(data []byte) byte {
	n := len(data) - 1
	p := bytes.IndexByte(data, '"')
	if p >= 0 {
		if p > 0 {
			c := data[p-1]
			if c == ',' || c == ';' {
				return c
			}
		}
		p1 := bytes.IndexByte(data[p+1:], '"')
		if p1 < 0 {
			p = p1
		} else {
			p = p1 + p + 1
		}
		for p >= 0 && p < n {
			c := data[p+1]
			if c == ',' || c == ';' {
				return c
			}
			if c != '"' {
				break
			}
			p1 = bytes.IndexByte(data[p+2:], '"')
			if p1 < 0 {
				p = p1
			} else {
				p = p1 + p + 2
			}
		}
	}
	n1 := countCharEntries(data, n, ';')
	n2 := countCharEntries(data, n, ',')
	if n1 > n2 {
		return ';'
	}
	return ','
}

func ReadCsvFromBytes(b []byte) (m map[string][][]string, err error) {
	m = make(map[string][][]string)
	sep := DetectColumnSeparator(b)
	n := len(b)
	maxp := 100
	f := make([]string, maxp)
	p := 0
	for i := 0; i < n; {
		r := b[i]
		if r == 10 || r == 13 {
			if p > 0 {
				placeReadyItemToMap(f, m, p)
				p = 0
			}
			i++
			continue
		}
		if r == sep {
			if p >= maxp {
				f = append(f, "")
				maxp++
			} else {
				f[p] = ""
			}
			p++
			i++
		} else if r == '"' {
			i++
			ps := i
			pe := i
			for ; i < n; i++ {
				if b[i] == '"' {
					if i+1 < n && b[i+1] == '"' {
						b[pe] = '"'
						pe++
						i++
					} else {
						s := string(b[ps:pe])
						if p >= maxp {
							f = append(f, s)
							maxp++
						} else {
							f[p] = s
						}
						p++
						i++
						break
					}
				} else {
					b[pe] = b[i]
					pe++
				}
			}
			for ; i < n; i++ {
				b1 := b[i]
				if b1 == 10 || b1 == 13 {
					break
				}
				if b1 == sep {
					i++
					break
				}
			}
		} else {
			ps := i
			for ; i < n; i++ {
				if b[i] == sep || b[i] == 13 || b[i] == 10 {
					s := string(b[ps:i])
					if p >= maxp {
						f = append(f, s)
						maxp++
					} else {
						f[p] = s
					}
					p++
					if b[i] == sep {
						i++
					}
					break
				}
			}
		}
	}
	if p > 0 {
		placeReadyItemToMap(f, m, p)
	}
	return
}

func ReadCsvFromFile(name string) (m map[string][][]string, err error) {
	var data []byte
	data, err = ioutil.ReadFile(name)
	if err != nil {
		return
	}
	return ReadCsvFromBytes(data)
}
