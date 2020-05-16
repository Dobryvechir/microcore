// package dbcsv reads / writes data in csv format
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvcsv

import (
	"bytes"
	"io/ioutil"
)

const (
	CsvSeparatorComma = 1
)

func WriteCsvToFile(name string, m map[string][][]string, options int) error {
	data := WriteCsvToBytes(m, options)
	return ioutil.WriteFile(name, data, 0644)
}

func evaluateBufferSize(m map[string][][]string, columns int) int {
	r := 0
	for k, v := range m {
		n := len(v)
		e := columns*3 + 10 + len([]byte(k))
		for j := 0; j < n; j++ {
			p := v[j]
			l := len(p)
			r += e
			for i := 0; i < l; i++ {
				r += len([]byte(p[i]))
			}
		}
	}
	return r
}

func detectMaximumNumberOfColumns(m map[string][][]string) int {
	if m == nil {
		return 0
	}
	r := 0
	for _, v := range m {
		if len(v) == 0 {
			continue
		}
		n := len(v[0])
		if n > r {
			r = n
		}
	}
	return r
}

func WriteCsvToBytes(m map[string][][]string, options int) (b []byte) {
	columns := detectMaximumNumberOfColumns(m)
	b = make([]byte, 0, evaluateBufferSize(m, columns+1))
	sep := byte(';')
	if (options & CsvSeparatorComma) != 0 {
		sep = ','
	}
	for k, v := range m {
		key := []byte("\"" + k + "\"")
		n := len(v)
		if n == 0 {
			continue
		}
		for j := 0; j < n; j++ {
			s := v[j]
			p := len(s)
			if p != 0 {
				b = append(b, key...)
				for i := 0; i < columns; i++ {
					t := ""
					if i < p {
						t = s[i]
					}
					b = append(b, sep)
					if t != "" {
						b = placeStringToBufForCsv(b, []byte(t))
					}
				}
				b = append(b, 13, 10)
			}
		}
	}
	return
}

func placeStringToBufForCsv(b []byte, p []byte) []byte {
	b = append(b, '"')
	if bytes.IndexByte(p, '"') < 0 {
		b = append(b, p...)
	} else {
		m := len(p)
		for h := 0; h < m; h++ {
			c := p[h]
			if c == '"' {
				b = append(b, '"', '"')
			} else {
				b = append(b, c)
			}
		}
	}
	b = append(b, '"')
	return b
}
