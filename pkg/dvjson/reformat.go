/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strconv"
)

const ExtraReformatJsonProblemCleanExclamationInValue = 1

func reformatTransferCharacter(c byte, addSpace int, dst []byte, dstPnt int, crlf int) (newDst []byte, newDstPnt int) {
	n := len(dst)
	newDstPnt = dstPnt + 1 + addSpace + crlf
	if newDstPnt > n {
		dst = append(dst, make([]byte, 51970+newDstPnt-n)...)
	}
	for i := 0; i < addSpace; i++ {
		dst[dstPnt] = ' '
		dstPnt++
	}
	dst[dstPnt] = c
	dstPnt++
	if crlf > 0 {
		if crlf == 2 {
			dst[dstPnt] = 13
			dstPnt++
		}
		dst[dstPnt] = 10
	}
	newDst = dst
	return
}

func insertCrLfIfNotInserted(crlf int, dst []byte, dstPnt int) (newDst []byte, newDstPnt int) {
	if dstPnt == 0 || crlf <= 0 || dst[dstPnt-1] < ' ' {
		return dst, dstPnt
	}
	n := len(dst)
	newDstPnt = dstPnt + crlf
	if newDstPnt > n {
		dst = append(dst, make([]byte, 51970+newDstPnt-n)...)
	}
	if crlf == 2 {
		dst[dstPnt] = 13
		dstPnt++
	}
	dst[dstPnt] = 10
	newDst = dst
	return
}

func reformatTransferBuffer(buffer []byte, addSpace int, dst []byte, dstPnt int) (newDst []byte, newDstPnt int) {
	n := len(dst)
	b := len(buffer)
	newDstPnt = dstPnt + addSpace + b
	if newDstPnt > n {
		dst = append(dst, make([]byte, 51970+newDstPnt-n)...)
	}
	for i := 0; i < addSpace; i++ {
		dst[dstPnt] = ' '
		dstPnt++
	}
	for i := 0; i < b; i++ {
		dst[dstPnt] = buffer[i]
		dstPnt++
	}
	newDst = dst
	return
}

func getJsonStringEndPos(buffer []byte, pos int) (endPos int, ok bool) {
	n := len(buffer)
	for i := pos; i < n; i++ {
		c := buffer[i]
		if c == '"' {
			return i + 1, true
		}
		if c == '\\' {
			i++
		}
	}
	return 0, false
}

func ensureInJsonWord(buffer []byte, pos int, word []byte) (int, bool) {
	n := len(word)
	m := len(buffer)
	ok := pos+n <= m
	if ok {
		for i := 0; i < n; i++ {
			if word[i] != buffer[pos+i] {
				ok = false
				break
			}
		}
	}
	return pos + n, ok
}

func ensureInJsonNumber(buffer []byte, pos int) (int, bool) {
	n := len(buffer)
	for ; pos < n; pos++ {
		c := buffer[pos]
		if !(c >= '0' && c <= '9' || c == '+' || c == '-' || c == 'e' || c == 'E' || c == '.') {
			break
		}
	}
	return pos, true
}

func getJsonNonStringEndPos(buffer []byte, pos int) (endPos int, ok bool) {
	c := buffer[pos]
	switch c {
	case 'f':
		return ensureInJsonWord(buffer, pos, []byte("false"))
	case 't':
		return ensureInJsonWord(buffer, pos, []byte("true"))
	case 'n':
		return ensureInJsonWord(buffer, pos, []byte("null"))
	default:
		if c >= '0' && c <= '9' || c == '-' || c == '+' {
			return ensureInJsonNumber(buffer, pos)
		}
	}
	return 0, false
}

func getJsonTrashPart(pos int, json []byte, trash [][]byte) (int, int) {
	n := len(trash)
	c := json[pos]
	m := len(json)
looper:
	for i := 0; i < n; i++ {
		refuse := trash[i]
		if refuse[0] == c {
			r := len(refuse)
			if pos+r <= m {
				for j := 1; j < r; j++ {
					if json[pos+j] != refuse[j] {
						continue looper
					}
				}
				extra := 0
				if string(refuse) == "(MISSING)" {
					extra = ExtraReformatJsonProblemCleanExclamationInValue
				}
				return pos + r, extra
			}
		}
	}
	return -1, 0
}

func calculateLineColumn(buf []byte, pos int) string {
	n := len(buf)
	if pos+1 >= n {
		return "at very end"
	}
	line := 1
	column := 0
	for i := 0; i <= pos; i++ {
		c := buf[i]
		if c == 13 {
			line++
			column = 0
			if buf[i+1] == 10 {
				i++
			}
		} else if c == 10 {
			line++
			column = 0
		} else {
			column++
		}
	}
	return strconv.Itoa(line) + ":" + strconv.Itoa(column)
}

func reportProblemInJson(message string, pos int, json []byte, state int, level int, stack []byte) ([]byte, error) {
	start := pos - 20
	end := pos + 20
	prePoints := "..."
	postPoints := "..."
	if start <= 0 {
		start = 0
		prePoints = ""
	}
	if end >= len(json) {
		end = len(json)
		postPoints = ""
	}
	place := prePoints + string(json[start:pos]) + ">" + string(json[pos:pos+1]) + "<" + string(json[pos+1:end])
	place += postPoints + " [" + strconv.Itoa(state) + "," + strconv.Itoa(level) + " " + string(stack[:level+2]) + "]"
	return nil, errors.New(message + " at " + calculateLineColumn(json, pos) + " in " + place)
}

func reportNotExpectedCharProblemInJson(c byte, pos int, json []byte, state int, level int, stack []byte) ([]byte, error) {
	return reportProblemInJson(string([]byte{c})+" is not expected", pos, json, state, level, stack)
}

func reportExpectedAnotherCharProblemInJson(real byte, expected byte, pos int, json []byte, state int, level int, stack []byte) ([]byte, error) {
	return reportProblemInJson(string([]byte{real})+" is not expected (expected "+string([]byte{expected})+")", pos, json, state, level, stack)
}

func reformatJsonValueCleanExclamationInValue(dst []byte, dstPnt int) ([]byte, int) {
	pnt := dstPnt - 1
	for pnt > 0 && dst[pnt] <= ' ' {
		pnt--
	}
	if pnt < 3 || dst[pnt] != '"' || dst[pnt-1] != '!' {
		return dst, dstPnt
	}
	for ; pnt < dstPnt; pnt++ {
		dst[pnt-1] = dst[pnt]
	}
	return dst, dstPnt - 1
}

func ReformatJson(json []byte, indent int, trash [][]byte, crlf int) ([]byte, error) {
	n := len(json)
	dstPnt := 0
	dst := make([]byte, n<<3)
	level := 0
	stack := make([]byte, 20, 20)
	currentIndent := 0
	state := STATE_INITIAL
	for i := 0; i < n; i++ {
		for i < n && json[i] <= ' ' {
			i++
		}
		if i == n {
			break
		}
		c := json[i]
		switch c {
		case '[', '{':
			space := currentIndent
			if state == STATE_VALUE_STARTED {
				space = 1
			} else if state != STATE_ARRAY_STARTED && state != STATE_INITIAL {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			back := byte(']')
			if c == '{' {
				back = '}'
				state = STATE_OBJECT_STARTED
			} else {
				state = STATE_ARRAY_STARTED
			}
			dst, dstPnt = reformatTransferCharacter(c, space, dst, dstPnt, crlf)
			if level == len(stack) {
				stack = append(stack, back)
			} else {
				stack[level] = back
			}
			level++
			currentIndent += indent
		case '}', ']':
			if state == STATE_ARRAY_COMMA_OR_END_EXPECTED || state == STATE_OBJECT_COMMA_OR_END_EXPECTED {
				dst, dstPnt = insertCrLfIfNotInserted(crlf, dst, dstPnt)
			} else if state != STATE_ARRAY_STARTED && state != STATE_OBJECT_STARTED {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			level--
			if stack[level] != c {
				return reportExpectedAnotherCharProblemInJson(c, stack[level], i, json, state, level, stack)
			}
			currentIndent -= indent
			if level == 0 {
				state = STATE_WAITING_END
			} else {
				if stack[level-1] == '}' {
					state = STATE_OBJECT_COMMA_OR_END_EXPECTED
				} else {
					state = STATE_ARRAY_COMMA_OR_END_EXPECTED
				}
			}
			dst, dstPnt = reformatTransferCharacter(c, currentIndent, dst, dstPnt, crlf)
		case ',':
			if state == STATE_ARRAY_COMMA_OR_END_EXPECTED {
				state = STATE_ARRAY_STARTED
			} else if state == STATE_OBJECT_COMMA_OR_END_EXPECTED {
				state = STATE_OBJECT_STARTED
			} else {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			for dstPnt > 0 && dst[dstPnt-1] <= ' ' {
				dstPnt--
			}
			dst, dstPnt = reformatTransferCharacter(c, 0, dst, dstPnt, crlf)
		case '"':
			space := currentIndent
			if state == STATE_OBJECT_STARTED {
				state = STATE_OBJECT_COLON
			} else if state == STATE_ARRAY_STARTED {
				state = STATE_ARRAY_COMMA_OR_END_EXPECTED
			} else if state == STATE_VALUE_STARTED {
				state = STATE_OBJECT_COMMA_OR_END_EXPECTED
				space = 1
			} else {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			endPos, ok := getJsonStringEndPos(json, i+1)
			if !ok {
				return reportProblemInJson("Quote is unclosed", i, json, state, level, stack)
			}
			dst, dstPnt = reformatTransferBuffer(json[i:endPos], space, dst, dstPnt)
			i = endPos - 1
		case ':':
			if state != STATE_OBJECT_COLON {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			dst, dstPnt = reformatTransferCharacter(c, 0, dst, dstPnt, 0)
			state = STATE_VALUE_STARTED
		default:
			endPos, extra := getJsonTrashPart(i, json, trash)
			if endPos >= 0 {
				i = endPos - 1
				switch extra {
				case ExtraReformatJsonProblemCleanExclamationInValue:
					dst, dstPnt = reformatJsonValueCleanExclamationInValue(dst, dstPnt)
				}
				continue
			}
			space := currentIndent
			if state == STATE_ARRAY_STARTED {
				state = STATE_ARRAY_COMMA_OR_END_EXPECTED
			} else if state == STATE_VALUE_STARTED {
				state = STATE_OBJECT_COMMA_OR_END_EXPECTED
				space = 1
			} else {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			endPos, ok := getJsonNonStringEndPos(json, i)
			if !ok {
				return reportNotExpectedCharProblemInJson(c, i, json, state, level, stack)
			}
			dst, dstPnt = reformatTransferBuffer(json[i:endPos], space, dst, dstPnt)
			i = endPos - 1
		}
	}
	if state != STATE_INITIAL && state != STATE_WAITING_END {
		return reportProblemInJson("Json is incomplete", n-1, json, state, level, stack)
	}
	return dst[:dstPnt], nil
}

func ConvertStringArrayToByteByteArray(data []string) [][]byte {
	n := len(data)
	res := make([][]byte, n)
	for i := 0; i < n; i++ {
		res[i] = []byte(data[i])
	}
	return res
}

func ConvertFieldItemArrayIntoMap(info *dvevaluation.DvVariable, keyField string, valueField string, defValue string) (res map[string]string, err error) {
	res = make(map[string]string)
	if info == nil || info.Kind != dvevaluation.FIELD_ARRAY {
		return
	}
	n := len(info.Fields)
	for i := 0; i < n; i++ {
		if info.Fields[i] != nil {
			keyNode := info.Fields[i].ReadSimpleChild(keyField)
			if keyNode == nil {
				continue
			}
			key := string(keyNode.Value)
			valueNode := info.Fields[i].ReadSimpleChild(valueField)
			value := defValue
			if valueNode != nil {
				value = dvtextutils.GetUnquotedString(valueNode.GetStringValue())
			}
			res[key] = value
		}
	}
	return
}
