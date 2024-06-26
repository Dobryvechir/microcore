/*
**********************************************************************
MicroCore
Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
***********************************************************************
*/
package dvtextutils

import (
	"fmt"
	"strconv"
	"strings"
)

var UINT8_COUNT = 256
var closingMap = createClosingMap()

func createClosingMap() []uint8 {
	res := make([]uint8, UINT8_COUNT)
	for i := 0; i < UINT8_COUNT; i++ {
		res[i] = 0
	}
	res['{'] = '}'
	res['['] = ']'
	res['('] = ')'
	res['"'] = '"'
	res['\''] = '\''
	return res
}

func AddNonRepeatingWords(s string, oldList, newList []string, imap map[string]int, plain, joiner string) string {
	if len(oldList) > 0 {
		l := len(newList)
		if l < 33 {
			n := uint(0)
			for _, v := range oldList {
				if k, ok := imap[v]; ok {
					n |= 1 << uint(k)
				}
			}
			for i, v := range newList {
				if (n & (1 << uint(i))) == 0 {
					s += joiner + v
				}
			}

		} else {
			n := make([]uint, (l+31)>>5)
			for _, v := range oldList {
				if k, ok := imap[v]; ok {
					n[k>>5] |= 1 << uint(k&31)
				}
			}
			for i, v := range newList {
				if (n[i>>5] & (1 << uint(i&31))) == 0 {
					s += joiner + v
				}
			}

		}
		return s
	} else {
		s = plain
	}
	return s
}

func AddStringListWithoutRepeats(src []string, newList []string) []string {
	m := len(newList)
addListMain:
	for i := 0; i < m; i++ {
		s := newList[i]
		n := len(src)
		for j := 0; j < n; j++ {
			if src[j] == s {
				continue addListMain
			}
		}
		src = append(src, s)
	}
	return src
}

func ReduceSpaceAndCountWords(str string) (string, int) {
	sarray := strings.Split(str, " ")
	count := 0
	for i := 0; i < len(sarray); i++ {
		if sarray[i] != "" {
			if i != count {
				sarray[count] = sarray[i]
			}
			count++
		}
	}
	sarray = sarray[:count]
	result := strings.Join(sarray, " ")
	return result, count
}

func PrepareAndMayQuoteParams(src []string) (string, int) {
	count := 0
	params := ""
	for _, dat := range src {
		if dat == "" {
			continue
		}
		count++
		if strings.IndexByte(dat, ' ') >= 0 {
			params += " \"" + dat + "\""
		} else {
			params += " " + dat
		}
	}
	return params, count
}

func ConvertToList(lst string) []string {
	return strings.Split(strings.TrimSpace(strings.Replace(strings.Replace(lst, ",", " ", -1), ";", " ", -1)), " ")
}

func ReduceListToNonEmptyList(res []string) []string {
	k := 0
	for i := 0; i < len(res); i++ {
		res[i] = strings.TrimSpace(res[i])
		if res[i] != "" {
			if k != i {
				res[k] = res[i]
			}
			k++
		}
	}
	return res[:k]
}

func ConvertToNonEmptyList(lst string) []string {
	return ReduceListToNonEmptyList(ConvertToList(lst))
}

func ConvertToNonEmptyListBySeparator(lst string, separator string) []string {
	r := strings.Split(lst, separator)
	return ReduceListToNonEmptyList(r)
}

func ConvertToNonEmptyListByEOL(lst string) []string {
	return ConvertToNonEmptyListBySeparator(lst, "\n")
}

func ConvertURLToList(url string) []string {
	return ReduceListToNonEmptyList(strings.Split(url, "/"))
}

func ConvertToNonEmptySemicolonList(lst string) []string {
	return ReduceListToNonEmptyList(strings.Split(lst, ";"))
}

func QuickLookJsonOption(s string, option string) string {
	n := len(s)
	p := strings.Index(s, "\""+option+"\"")
	if p < 0 {
		return ""
	}
	for p += 2 + len(option); p < n && s[p] <= ' '; p++ {
	}
	if p >= n {
		return ""
	}
	if s[p] != ':' {
		return QuickLookJsonOption(s[p:], option)
	}
	for p++; p < n && s[p] <= ' '; p++ {
	}
	if p >= n {
		return ""
	}
	e := p
	for ; e < n && s[e] != '}' && s[e] != ','; e++ {
	}
	return strings.TrimSpace(s[p:e])
}

func MakeUniqueStringList(lists ...[]string) []string {
	res := make([]string, 0, 20)
	exist := make(map[string]bool)
	n := len(lists)
	for i := 0; i < n; i++ {
		l := lists[i]
		m := len(l)
		for j := 0; j < m; j++ {
			s := l[j]
			if !exist[s] {
				exist[s] = true
				res = append(res, s)
			}
		}
	}
	return res
}

func GetNextWordBySpaceTable(s string, spaceTable map[byte]bool, allRest bool) string {
	n := len(s)
	i := 0
	for ; i < n && (s[i] <= ' ' || spaceTable[s[i]]); i++ {

	}
	if i == n {
		return ""
	}
	pos := i
	if allRest {
		return s[pos:]
	}
	for i++; i < n && s[i] > ' ' && !spaceTable[s[i]]; i++ {

	}
	return s[pos:i]
}

func GetNextWordInText(s string) string {
	n := len(s)
	i := 0
	for ; i < n && s[i] <= ' '; i++ {

	}
	if i == n {
		return ""
	}
	pos := i
	for i++; i < n && s[i] > ' '; i++ {

	}
	return s[pos:i]
}

var yamlControlMap = map[byte]bool{
	'|': true,
}

func GetNextWordExceptYamlControls(s string) string {
	return GetNextWordBySpaceTable(s, yamlControlMap, false)
}

func GetNextNonEmptyPartInYaml(s string) string {
	return GetNextWordBySpaceTable(s, yamlControlMap, true)
}

func GetStringArrayWithDefaults(src []string, defs []string) []string {
	n := len(src)
	m := len(defs)
	if n >= m {
		return src
	}
	res := make([]string, m)
	for i := 0; i <= n; i++ {
		res[i] = src[i]
	}
	for ; n <= m; n++ {
		res[n] = defs[n]
	}
	return res
}

func TryReadInteger(str string, def int) int {
	str = strings.TrimSpace(str)
	val := def
	if str != "" {
		n, err := strconv.Atoi(str)
		if err != nil {
			val = n
		}
	}
	return val
}

func InsertTextIntoBuffer(src []byte, posStart int, posEnd int, buf ...[]byte) (dst []byte, dif int) {
	n := len(src)
	k := len(buf)
	m := 0
	for i := 0; i < k; i++ {
		m += len(buf)
	}
	if posEnd < 0 {
		posEnd = 0
	}
	if posEnd > n {
		posEnd = n
	}
	if posStart > posEnd {
		posStart = posEnd
	}
	if posStart < 0 {
		posStart = 0
	}
	p := posEnd - posStart
	dif = m - p
	dst = make([]byte, n+dif)
	i := 0
	for ; i < posStart; i++ {
		dst[i] = src[i]
	}
	for j := 0; j < k; j++ {
		b := buf[j]
		l := len(b)
		for o := 0; o < l; o++ {
			dst[i] = b[o]
			i++
		}
	}
	for j := posEnd; j < n; j++ {
		dst[i] = src[j]
		i++
	}
	return
}

func ReadInsideBrackets(str string, pos int) (endPos int, err error) {
	opening := str[pos]
	var closing byte = closingMap[opening]
	if closing == 0 || closing == opening {
		return 0, fmt.Errorf("Unknown bracket in %s at %d", str, pos)
	}
	pos++
	count := 1
	stack := make([]byte, 32)
	stack[0] = closing
	n := len(str)
	for ; pos < n; pos++ {
		b := str[pos]
		c := closingMap[b]
		if b == stack[count-1] {
			count--
			if count == 0 {
				return pos, nil
			}
		} else if c != 0 {
			if b != c {
				if count < len(stack) {
					stack[count] = c
				} else {
					stack = append(stack, c)
				}
				count++
			} else if b != 0 {
				for pos++; pos < n; pos++ {
					d := str[pos]
					if d == '\\' {
						pos++
					} else if d == b {
						break
					}
				}
			}
		}
	}
	return 0, fmt.Errorf("Unclosed [ in %s at %d", str, pos)
}

func CopyStringMap(src map[string]string) map[string]string {
	dst := make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func ReplaceWordBySpaceOrSemicolonOrEnd(src string, repl string, pos int) string {
	endPos := pos
	n := len(src)
	for ; endPos < n; endPos++ {
		c := src[endPos]
		if c <= 32 || c == ';' {
			break
		}
	}
	return src[:pos] + repl + src[endPos:]
}

func SeparateChildExpression(name string) (res []string) {
	res = make([]string, 0, 7)
	n := len(name)
	if n == 0 {
		return
	}
	pos := 0
	for pos < n {
		for pos < n && (name[pos] <= ' ' || name[pos] == '.') {
			pos++
		}
		if pos == n {
			break
		}
		c := name[pos]
		if c == '[' || c == '{' {
			endPos, err := ReadInsideBrackets(name, pos)
			if err != nil {
				res = append(res, name[pos:])
				return
			}
			endPos++
			if endPos < n && name[endPos] == '?' {
				endPos++
			}
			res = append(res, name[pos:endPos])
			pos = endPos
		} else {
			endPos := pos + 1
			var err error
			for ; endPos < n; endPos++ {
				c := name[endPos]
				if c == '.' || c == '[' || c == '{' {
					break
				}
				if c == '(' {
					endPos, err = ReadInsideBrackets(name, pos)
					if err != nil {
						res = append(res, name[pos:])
						return
					}
					endPos++
					if endPos < n && name[endPos] == '?' {
						endPos++
					}
					break
				}
			}
			res = append(res, name[pos:endPos])
			pos = endPos
		}
	}
	return
}

func SmartReadStringList(s string, nonEmpty bool) []string {
	s = strings.TrimSpace(s)
	res := make([]string, 0, 16)
	if s == "" {
		return res
	}
	if s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
	}
	n := len(s)
	for i := 0; i < n; {
		c := s[i]
		if c <= ' ' {
			i++
			continue
		}
		if c == '"' || c == '\'' {
			pos, screened := SmartReadQuotedStringEndPos(s, i)
			if pos < 0 {
				t := s[i:]
				res = append(res, t)
				break
			} else {
				t := s[i+1 : pos]
				if screened {
					t = ReadScreenedString(t)
				}
				i = pos + 1
				if len(t) > 0 || !nonEmpty {
					res = append(res, t)
				}
			}
			for ; i < n; i++ {
				if s[i] == ',' {
					i++
					break
				}
			}
		} else {
			pos := SmartReadUnquotedStringEndPos(s, i)
			t := strings.TrimSpace(s[i:pos])
			i = pos + 1
			if len(t) > 0 || !nonEmpty {
				res = append(res, t)
			}
		}
	}
	return res
}

func SmartReadQuotedStringEndPos(s string, pos int) (int, bool) {
	n := len(s)
	c := s[pos]
	screened := false
	for pos++; pos < n; pos++ {
		t := s[pos]
		if t == c {
			return pos, screened
		} else if t == '\\' {
			pos++
			screened = true
		}
	}
	return -1, false
}

func SmartReadUnquotedStringEndPos(s string, pos int) int {
	n := len(s)
	for ; pos < n; pos++ {
		if s[pos] == ',' {
			return pos
		}
	}
	return n
}

func ReadScreenedString(s string) string {
	t := []byte(s)
	n := len(t)
	pos := 0
	for i := 0; i < n; i++ {
		b := t[i]
		if b == '\\' && i < n-1 {
			i++
			t[pos] = t[i]
			pos++
		} else {
			t[pos] = b
			pos++
		}
	}
	return string(t[:pos])
}

func SeparateBytesToUTF8Chars(b []byte) [][]byte {
	n := len(b)
	res := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		c := b[i]
		if c < 0xc0 || i == n-1 {
			res = append(res, []byte{c})
		} else if c < 0xe0 || i == n-2 {
			i++
			res = append(res, []byte{c, b[i]})
		} else if c < 0xf0 || i == n-3 {
			i++
			res = append(res, []byte{c, b[i], b[i+1]})
			i++
		} else if c < 0xf8 || i == n-4 {
			i++
			res = append(res, []byte{c, b[i], b[i+1], b[i+2]})
			i += 2
		} else if c < 0xfc || i == n-5 {
			i++
			res = append(res, []byte{c, b[i], b[i+1], b[i+2], b[i+3]})
			i += 3
		} else if c < 0xfe || i == n-6 {
			i++
			res = append(res, []byte{c, b[i], b[i+1], b[i+2], b[i+3], b[i+4]})
			i += 4
		} else {
			res = append(res, []byte{c})
		}
	}
	return res
}

func GetCodePoint(b []byte) int {
	n := len(b)
	if n == 0 {
		return 0
	}
	c := b[0]
	if c < 0xc0 {
		return int(c)
	} else if c < 0xe0 {
		return createCodePointFromBytes(c&0x1f, b[1:], 1)
	} else if c < 0xf0 {
		return createCodePointFromBytes(c&0xf, b[1:], 2)
	} else if c < 0xf8 {
		return createCodePointFromBytes(c&0x7, b[1:], 3)
	} else if c < 0xfc {
		return createCodePointFromBytes(c&0x3, b[1:], 4)
	} else if c < 0xfe {
		return createCodePointFromBytes(c&0x1, b[1:], 5)
	} else {
		return int(c)
	}
}

func createCodePointFromBytes(val byte, b []byte, m int) int {
	res := int(val)
	n := len(b)
	for i := 0; i < m; i++ {
		d := 0
		if i < n {
			d = int(b[i]) & 0x3f
		}
		res = (res << 6) | d
	}
	return res
}

func GetBytesFromPointCode(v int) []byte {
	v = v & 0x7fffffff
	if v <= 0x7f {
		return []byte{byte(v)}
	}
	if v <= 0x7ff {
		return createBytesFromPointCode(v, 0xc0, 5, 1)
	}
	if v <= 0xffff {
		return createBytesFromPointCode(v, 0xe0, 4, 2)
	}
	if v <= 0x1fffff {
		return createBytesFromPointCode(v, 0xf0, 3, 3)
	}
	if v <= 0x3ffffff {
		return createBytesFromPointCode(v, 0xf8, 2, 4)
	}
	return createBytesFromPointCode(v, 0xfc, 1, 5)
}

func createBytesFromPointCode(v int, first byte, bits int, rest int) []byte {
	buf := make([]byte, rest+1)
	for i := 0; i < rest; i++ {
		buf[rest-i] = byte(v&0x3f) | 0x80
		v >>= 6
	}
	buf[0] = first | byte(v&((1<<bits)-1))
	return buf
}

func FindNonEmptyLastString(data []string) string {
	n := len(data) - 1
	for ; n >= 0; n-- {
		if len(data[n]) > 0 {
			return data[n]
		}
	}
	return ""
}

func GetLowCaseExtension(name string) string {
	p := strings.LastIndex(name, ".")
	if p < 0 {
		return ""
	}
	return strings.ToLower(name[p:])
}
