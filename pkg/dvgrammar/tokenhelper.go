/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

func readNumber(data []byte, i int, amount int) int {
	state := 0
	for i < amount {
		c := data[i]
		if c == '.' {
			if state != 0 {
				return i - 1
			}
			state = 1
			i++
			continue
		}
		if c == 'e' || c == 'E' && state < 2 {
			state = 2
			i++
			if i < amount {
				c = data[i]
				if (c == '+' || c == '-') && i+1 < amount {
					i++
					c = data[i]
					if c >= '0' && c <= '9' {
						i++
						continue
					}
				}
				if c >= '0' && c <= '9' {
					i++
					continue
				}
			}
			return i - 2
		}
		if c >= '0' && c <= '9' {
			i++
			continue
		}
		break
	}
	return i - 1
}

func GetEscapedString(data []byte) string {
	l := len(data)
	res := make([]byte, l)
	k := 0
	for i := 0; i < l; i++ {
		c := data[i]
		if c != '\\' {
			res[k] = c
			k++
		} else {
			i++
			c = data[i]
			switch c {
			case 'n':
				c = 13
			case 'r':
				c = 10
			case 't':
				c = 9
			}
			res[k] = c
			k++
		}
	}
	return string(res[:k])
}

func (b *TokenBlock) findOperator(grammar *GrammarTable, data []byte, pos int) (oper string, nextPos int) {
	if b.MaxLength == 0 {
		return
	}
	l := len(data)
	n := l - pos
	if n > b.MaxLength {
		n = b.MaxLength
	}
	for n > 0 {
		k := n
		n--
		if b.Controls[n] != nil {
			var s string
			if b.MaxWords == 1 {
				s = string(data[pos : pos+k])
			} else {
				j := k
				i := 0
				sb := make([]byte, j)
				for k = 0; i < j && pos+k < l; {
					c := data[pos+k]
					k++
					sb[i] = c
					i++
					if grammar.TokenMap[c].Kind == BYTE_SPACE {
						sb[i-1] = ' '
						for pos+k < l && grammar.TokenMap[data[k]].Kind == BYTE_SPACE {
							k++
						}
					}
				}
				if i < j {
					continue
				}
				s = string(sb)
			}
			res, ok := b.Controls[n][s]
			if ok {
				endPos := pos + k
				if endPos < l {
					c1 := grammar.TokenMap[data[endPos-1]].Kind
					c2 := grammar.TokenMap[data[endPos]].Kind
					if (c1 == BYTE_DIGIT || c1 == BYTE_LETTER) && (c2 == BYTE_DIGIT || c2 == BYTE_LETTER) {
						continue
					}
				}
				return res, endPos
			}
		}
	}
	return "", -1
}
