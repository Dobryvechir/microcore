/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strconv"
)

func readJsonStringPart(data []byte, pos int) ([]byte, int, error) {
	escaped := false
	pos++
	start := pos
	var buf []byte
	n := len(data)
	for ; pos < n; pos++ {
		b := data[pos]
		if b == '\\' {
			if !escaped {
				escaped = true
				buf = make([]byte, pos-start, pos-start+1024)
				for i := start; i < pos; i++ {
					buf[i-start] = data[i]
				}
			}
			pos++
			if pos < n {
				buf = append(buf, data[pos])
			}
		} else if b == '"' {
			break
		} else if escaped {
			buf = append(buf, b)
		}
	}
	if pos >= n {
		return nil, n, fmt.Errorf("Quote is not closed %s", getPositionErrorInfo(data, start))
	}
	var res []byte
	if escaped {
		res = buf
	} else {
		res = data[start:pos]
	}
	return res, pos + 1, nil
}

var nullWord = []byte("null")
var trueWord = []byte("true")
var falseWord = []byte("false")

func getSimpleStringType(data []byte) int {
	if bytes.Equal(data, nullWord) {
		return dvevaluation.FIELD_NULL
	}
	if bytes.Equal(data, trueWord) || bytes.Equal(data, falseWord) {
		return dvevaluation.FIELD_BOOLEAN
	}
	return -1
}

func readJsonSimplePart(data []byte, pos int) ([]byte, int, error, int) {
	b := data[pos]
	n := len(data)
	start := pos
	isNumber := b >= '0' && b <= '9' || b == '-' || b == '+' || b == '.'
	isWord := b == 'f' || b == 't' || b == 'n'
	if isNumber {
		for ; pos < n; pos++ {
			b = data[pos]
			if !(b >= '0' && b <= '9' || b == '+' || b == '-' || b == '.' || b == 'e' || b == 'E') {
				break
			}
		}
		return data[start:pos], pos, nil, dvevaluation.FIELD_NUMBER
	}
	if isWord {
		for ; pos < n; pos++ {
			b := data[pos]
			if !(b >= 'a' && b <= 'z') {
				break
			}
		}
		str := data[start:pos]
		kind := getSimpleStringType(str)
		if kind >= 0 {
			return str, pos, nil, kind
		}
		return nil, n, fmt.Errorf("Incorrect word %s %s", str, getPositionErrorInfo(data, start)), 0
	}
	return nil, n, fmt.Errorf("Unexpected character %s %s", string([]byte{b}), getPositionErrorInfo(data, start)), 0
}

func readJsonNextNonSpace(data []byte, pos int, n int) int {
	for ; pos < n; pos++ {
		if data[pos] > 32 {
			break
		}
	}
	return pos
}

func readJsonPart(data []byte, i int) (*dvevaluation.DvVariable, int, error) {
	n := len(data)
	for ; i < n && data[i] <= 32; i++ {
	}
	if i >= n {
		return nil, n, fmt.Errorf("Empty json ")
	}
	switch data[i] {
	case '{':
		fields := make([]*dvevaluation.DvVariable, 0, 20)
		i = readJsonNextNonSpace(data, i+1, n)
		for i < n && data[i] != '}' {
			if data[i] == '"' {
				key, nextPos, err := readJsonStringPart(data, i)
				if err != nil {
					return nil, n, err
				}
				i = readJsonNextNonSpace(data, nextPos, n)
				if i >= n || data[i] != ':' {
					return nil, n, fmt.Errorf("Expected colon %s", getPositionErrorInfo(data, n))
				}
				dvEntry, nxtPos, err1 := readJsonPart(data, i+1)
				if err1 != nil {
					return nil, n, err1
				}
				i = nxtPos
				dvEntry.Name = key
				fields = append(fields, dvEntry)
			} else {
				return nil, n, fmt.Errorf("Unexpected character %s %s", string(data[i:i+1]), getPositionErrorInfo(data, i))
			}
			i = readJsonNextNonSpace(data, i, n)
			if i < n && data[i] == ',' {
				i = readJsonNextNonSpace(data, i+1, n)
				continue
			}
			if data[i] != '}' {
				return nil, n, fmt.Errorf("Unexpected character %s %s (expected } or comma)", string(data[i:i+1]), getPositionErrorInfo(data, i))
			}
		}
		if i >= n {
			return nil, n, fmt.Errorf("Expected } at the end %s", string(data[i:i+1]), getPositionErrorInfo(data, n))
		}
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: fields}, i + 1, nil
	case '[':
		fields := make([]*dvevaluation.DvVariable, 0, 20)
		i = readJsonNextNonSpace(data, i+1, n)
		for i < n && data[i] != ']' {
			dvEntry, nextPos, err := readJsonPart(data, i)
			if err != nil {
				return nil, n, err
			}
			fields = append(fields, dvEntry)
			i = readJsonNextNonSpace(data, nextPos, n)
			if i < n && data[i] == ',' {
				i = readJsonNextNonSpace(data, i+1, n)
				continue
			}
			if data[i] != ']' {
				return nil, n, fmt.Errorf("Unexpected character %s %s (expected ] or comma)", string(data[i:i+1]), getPositionErrorInfo(data, i))
			}
		}
		if i >= n {
			return nil, n, fmt.Errorf("Expected ] at the end %s", string(data[i:i+1]), getPositionErrorInfo(data, n))
		}
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: fields}, i + 1, nil
	case '"':
		str, nextPos, err := readJsonStringPart(data, i)
		if err != nil {
			return nil, n, err
		}
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Value: str}, nextPos, nil
	default:
		str, nextPos, err, kind := readJsonSimplePart(data, i)
		if err != nil {
			return nil, n, err
		}
		return &dvevaluation.DvVariable{Kind: kind, Value: str}, nextPos, nil
	}

}

func ReadJsonAsDvFieldInfo(data []byte) (*dvevaluation.DvVariable, error) {
	dvEntry, pos, err := readJsonPart(data, 0)
	if err != nil {
		return nil, err
	}
	l := len(data)
	for ; pos < l; pos++ {
		if data[pos] > 32 {
			return nil, fmt.Errorf("\nUnexpected extra characters %s", getPositionErrorInfo(data, pos))
		}
	}
	return dvEntry, nil
}

func getPositionErrorInfo(data []byte, pos int) string {
	line := 1
	column := 1
	for i := 0; i < pos; i++ {
		b := data[i]
		if b == 13 || b == 10 {
			if b == 13 && i+1 < pos && data[i+1] == 10 {
				i++
			}
			line++
			column = 1
		} else {
			column++
		}
	}
	endPos := pos + 20
	if endPos > len(data) {
		endPos = len(data)
	}
	bufLen := endPos - pos
	addInfo := ""
	if bufLen > 0 {
		buf := make([]byte, bufLen)
		for i := 0; i < bufLen; i++ {
			buf[i] = data[pos+i]
			if buf[i] < 32 {
				buf[i] = '.'
			}
		}
		addInfo = string(buf)
	}
	return " at " + strconv.Itoa(line) + ":" + strconv.Itoa(column) + " [" + addInfo + "] "
}
