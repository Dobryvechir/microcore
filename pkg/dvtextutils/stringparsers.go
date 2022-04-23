/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
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
	res := make([]string, 0, 16)
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
		} else {
			pos := SmartReadUnquotedStringEndPos(s, i)
			t := strings.TrimSpace(s[i:pos])
			i = pos
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
