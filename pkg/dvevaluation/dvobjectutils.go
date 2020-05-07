/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

func (obj *DvObject) CalculateExpression(b []byte) ([]byte, error) {
	expr := ParseForDvObject(b, obj, 0, 0, string(b))
	if expr.Err != nil {
		return nil, expr.Err
	}
	return []byte(AnyToString(expr.FinalResult)), nil
}

func (obj *DvObject) CalculateStringWithBrackets(b []byte, level int) (string, error) {
	n := len(b)
	k := n + 1024
	r := make([]byte, k)
	m := 0
	v := 0
	for i := 0; i < n; i++ {
		c := b[i]
		if c == '{' {
			v++
			if v == level {
				m -= level
				k = len(r)
				p := level
				for i++; i < n && b[i] == '{'; i++ {
					p++
				}
				j := i
				v = 0
				for ; i < n; i++ {
					if b[i] == '}' {
						v++
						if v == p {
							break
						}
					} else {
						v = 0
					}
				}
				if v == p {
					d, err := obj.CalculateExpression(b[j : i+1-v])
					if err != nil {
						return "", err
					}
					s := len(d)
					if m+s < k {
						for j := 0; j < s; j++ {
							r[m] = d[j]
							m++
						}
					} else {
						r = append(r[:m], d...)
						k = len(r)
					}
				} else {
					r = append(r[:m], b[j-p:]...)
					m = len(r)
					i = n
				}
			} else {
				if m < k {
					r[m] = c
				} else {
					r = append(r, c)
				}
				m++
			}
		} else {
			if m < k {
				r[m] = c
			} else {
				r = append(r, c)
			}
			m++
			v = 0
		}
	}
	return string(r[:m]), nil
}

func (obj *DvObject) CalculateString(str string) (string, error) {
	return obj.CalculateStringWithBrackets([]byte(str), 2)
}
