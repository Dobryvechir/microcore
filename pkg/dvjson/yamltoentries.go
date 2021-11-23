/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func readYamlValue(data []byte, pos int, endChar byte, indent int) (res []byte, nextPos int, err error, tp int) {
	n := len(data)
	b := data[pos]
	if b == '\'' {
		for nextPos = pos + 1; nextPos < n && data[nextPos] != '\''; nextPos++ {
		}
		if nextPos >= n {
			err = fmt.Errorf("Unclosed single quote string %s", getPositionErrorInfo(data, pos))
			return
		}
		res = data[pos+1 : nextPos]
		tp = dvevaluation.FIELD_STRING
		nextPos++
		return
	}
	if b == '"' {
		buf := make([]byte, 0, 200)
		for nextPos = pos + 1; nextPos < n && data[nextPos] != '"'; nextPos++ {
			b := data[nextPos]
			if b == '\\' {
				nextPos++
				if nextPos < n {
					buf = append(buf, data[nextPos])
				}
			} else {
				buf = append(buf, b)
			}
		}
		if nextPos >= n {
			err = fmt.Errorf("Unclosed double quote string %s", getPositionErrorInfo(data, pos))
			return
		}
		res = buf
		tp = dvevaluation.FIELD_STRING
		nextPos++
		return
	}
	res = make([]byte, 0, 200)
	for nextPos = pos; nextPos < n; nextPos++ {
		b = data[nextPos]
		if b == endChar && b != 0 && (nextPos+1 == n || data[nextPos+1] <= ' ') {
			break
		}
		if b == 10 || b == 13 {
			if nextPos > pos {
				newPart := data[pos:nextPos]
				partLen := len(newPart)
				if len(res) == 0 && newPart[partLen-1] == '|' {
					if partLen == 1 {
						partLen = 0
					} else if newPart[partLen-2] <= ' ' {
						newPart = bytes.TrimSpace(newPart[:partLen-2])
					}
				}
				if partLen > 0 {
					res = append(res, newPart...)
				}
			}
			newPos, newIndent := readYamlNonEmptyLine(data, nextPos)
			if newIndent > indent {
				nextPos = newPos
				pos = newPos
				if len(res) == 2 && res[0] == '>' && res[1] == '-' {
					res = res[:0]
				} else {
					res = append(res, ' ')
				}
			} else {
				pos = nextPos
				break
			}
		}
	}
	if nextPos > pos {
		res = append(res, data[pos:nextPos]...)
	}
	l := len(res)
	for ; l > 0 && res[l-1] <= 32; l-- {
	}
	res = res[:l]
	tp = getSimpleStringType(res)
	if tp < 0 {
		if dvtextutils.IsJsonNumber(res) {
			tp = dvevaluation.FIELD_NUMBER
		} else {
			tp = dvevaluation.FIELD_STRING
		}
	}
	return
}

func readYamlNonEmptyLine(data []byte, pos int) (nxtPos int, indent int) {
	indent = 0
	for prevPos := pos; prevPos > 0 && data[prevPos-1] != 10 && data[prevPos-1] != 13; prevPos-- {
		indent++
	}
	n := len(data)
	for nxtPos = pos; nxtPos < n; nxtPos++ {
		b := data[nxtPos]
		if b == 13 || b == 10 {
			indent = 0
			continue
		}
		if b <= 32 {
			indent++
			continue
		}
		if b == '%' || b == '#' {
			for ; nxtPos < n; nxtPos++ {
				if data[nxtPos] == 13 || data[nxtPos] == 10 {
					break
				}
			}
			indent = 0
			continue
		}
		break
	}
	return
}

func readYamlNextCurrentKey(data []byte, pos int) (currentKey []byte, indent int, nextPos int, quoted bool, err error) {
	n := len(data)
	nextPos, indent = readYamlNonEmptyLine(data, pos)
	if nextPos >= n {
		return
	}
	c := data[nextPos]
	if c == '"' || c == '\'' {
		quoted = true
		quoteChar := data[nextPos]
		nextPos++
		pos = nextPos
		for ; nextPos < n; nextPos++ {
			c = data[nextPos]
			if c == quoteChar {
				break
			} else if c == '\\' {
				nextPos++
			}
		}
		if nextPos >= n {
			err = fmt.Errorf("Quote is not closed %s", getPositionErrorInfo(data, pos))
			return
		}
		currentKey = data[pos:nextPos]
		nextPos++
		return
	}
	pos = nextPos
	for ; nextPos < n; nextPos++ {
		c := data[nextPos]
		if c == '[' || c == ':' || c == '}' || c == ']' || c == '|' || c == ',' || c == '{' {
			break
		}
	}
	for pos < nextPos && data[pos] <= ' ' {
		pos++
	}
	lastPos := nextPos
	for lastPos > pos && data[lastPos] <= ' ' {
		lastPos--
	}
	currentKey = data[pos:lastPos]
	return
}

func readYamlFlowCollectionValue(data []byte, pos int) (obj *dvevaluation.DvVariable, nextPos int, err error) {
	nextPos, _ = readYamlNonEmptyLine(data, pos)
	n := len(data)
	if nextPos >= n {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_UNDEFINED}, nextPos, nil
	}
	switch data[nextPos] {
	case '{':
		return readYamlFlowCollectionMap(data, nextPos, '}')
	case '[':
		return readYamlFlowCollectionArray(data, nextPos, ']')
	}
	currentKey, _, nextPos, quoted, err := readYamlNextCurrentKey(data, nextPos)
	if err != nil {
		return
	}
	kind := dvevaluation.FIELD_STRING
	l := len(currentKey)
	if !quoted && l > 0 {
		errMes, startPos, endPos, nxtPos, newKind := jsonReadValue(currentKey, 0, l)
		if startPos == 0 && errMes == "" && endPos == nxtPos && endPos == nextPos && newKind != dvevaluation.FIELD_UNDEFINED {
			kind = newKind
		}
	}
	value := make([]byte, l)
	for i := 0; i < l; i++ {
		value[i] = currentKey[i]
	}
	obj = &dvevaluation.DvVariable{Kind: kind, Value: value}
	return
}

func readYamlFlowCollectionMap(data []byte, pos int, endChar byte) (obj *dvevaluation.DvVariable, nextPos int, err error) {
	n := len(data)
	obj = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, 7)}
	var currentKey []byte
	for nextPos = pos + 1; nextPos < n; nextPos++ {
		pos = nextPos
		currentKey, _, nextPos, _, err = readYamlNextCurrentKey(data, pos)
		if err != nil {
			return
		}
		if nextPos >= n {
			err = fmt.Errorf("Uncomplete object %s", getPositionErrorInfo(data, pos))
			return
		}
		if len(currentKey) == 0 {
			if data[nextPos] == endChar {
				nextPos++
				return
			}
			err = fmt.Errorf("Expected %s at the end of the object %s", string([]byte{endChar}), getPositionErrorInfo(data, pos))
			return
		}
		pos = nextPos
		nextPos, _ = readYamlNonEmptyLine(data, pos)
		if nextPos >= n {
			err = fmt.Errorf("Uncomplete object %s", getPositionErrorInfo(data, pos))
			return
		}
		c := data[nextPos]
		if c == ',' || c == '}' {
			item := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Name: currentKey, Value: currentKey}
			obj.Fields = append(obj.Fields, item)
			nextPos++
			if c == '}' {
				return
			}
			continue
		}
		if c != ':' {
			err = fmt.Errorf("Expected : but found %s %s", string([]byte{c}), getPositionErrorInfo(data, nextPos))
			return
		}
		var item *dvevaluation.DvVariable
		item, nextPos, err = readYamlFlowCollectionValue(data, nextPos+1)
		if err != nil {
			return
		}
		item.Name = currentKey
		obj.Fields = append(obj.Fields, item)
	}
	err = fmt.Errorf("Expected } at the end of object %s", getPositionErrorInfo(data, nextPos))
	return
}

func readYamlFlowCollectionArray(data []byte, pos int, endChar byte) (obj *dvevaluation.DvVariable, nextPos int, err error) {
	n := len(data)
	obj = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, 7)}
	for nextPos = pos + 1; nextPos < n; nextPos++ {
		pos = nextPos
		nextPos, _ = readYamlNonEmptyLine(data, pos)
		if nextPos >= n {
			err = fmt.Errorf("Uncomplete array %s", getPositionErrorInfo(data, pos))
			return
		}
		c := data[nextPos]
		if c == endChar {
			nextPos++
			return
		}
		if c == ',' {
			nextPos++
		} else if len(obj.Fields) != 0 {
			err = fmt.Errorf("Expected a comma(,) but found %s %s", string([]byte{c}), getPositionErrorInfo(data, nextPos))
			return
		}
		var item *dvevaluation.DvVariable
		item, nextPos, err = readYamlFlowCollectionValue(data, nextPos+1)
		if err != nil {
			return
		}
		obj.Fields = append(obj.Fields, item)
	}
	err = fmt.Errorf("Expected } at the end of object %s", getPositionErrorInfo(data, nextPos))
	return
}

func readYamlPartMap(data []byte, pos int, indent int, endChar byte, currentKey []byte) (*dvevaluation.DvVariable, int, error) {
	fields := make([]*dvevaluation.DvVariable, 0, 20)
	n := len(data)
	for pos < n {
		nextPos, newIndent := readYamlNonEmptyLine(data, pos)
		b := byte(0)
		if nextPos >= n {
			nextPos = n
		} else {
			b = data[nextPos]
		}
		if dvparser.FindEol(data[pos:nextPos]) < 0 && b != '[' && b != '{' {
			str, nextPos, err, tp := readYamlValue(data, nextPos, 0, indent)
			if err != nil {
				return nil, n, err
			}
			pos = nextPos
			fields = append(fields, &dvevaluation.DvVariable{Kind: tp, Value: str, Name: currentKey})
		} else {
			if newIndent != indent || b != '-' {
				newIndent = indent + 1
			}
			dvEntry, nextPos, err := readYamlAllFromIndent(data, pos, newIndent, endChar)
			if err != nil {
				return nil, n, err
			}
			fields = append(fields, dvEntry)
			dvEntry.Name = currentKey
			pos = nextPos
		}
		nextPos, newIndent = readYamlNonEmptyLine(data, pos)
		if newIndent != indent || nextPos >= n || data[nextPos] == '-' || data[nextPos] == endChar {
			break
		}
		var err error
		currentKey, nextPos, err, _ = readYamlValue(data, nextPos, ':', indent)
		if err != nil {
			return nil, n, err
		}
		pos, newIndent = readYamlNonEmptyLine(data, nextPos)
		if pos >= n || data[pos] != ':' || newIndent <= indent {
			return nil, n, fmt.Errorf("Expected colon  %s", getPositionErrorInfo(data, pos))
		}
		pos++
	}
	return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: fields}, pos, nil
}

func readYamlPartArray(data []byte, pos int, indent int, endChar byte) (*dvevaluation.DvVariable, int, error) {
	fields := make([]*dvevaluation.DvVariable, 0, 20)
	n := len(data)
	for pos < n && data[pos] == '-' {
		dvEntry, nxtPos, err := readYamlAllFromIndent(data, pos+1, indent+1, endChar)
		if err != nil {
			return nil, n, err
		}
		fields = append(fields, dvEntry)
		pos = nxtPos
		nextPos, newIndent := readYamlNonEmptyLine(data, nxtPos)
		if newIndent != indent {
			break
		}
		pos = nextPos
	}
	return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: fields}, pos, nil
}

func readYamlAllFromIndent(data []byte, pos int, indent int, endChar byte) (*dvevaluation.DvVariable, int, error) {
	n := len(data)
	nxtPos, newIndent := readYamlNonEmptyLine(data, pos)
	if nxtPos >= n || newIndent < indent {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Value: make([]byte, 0)}, nxtPos, nil
	}
	oldIndent := indent
	indent = newIndent
	b := data[nxtPos]
	if b == endChar && b != 0 {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Value: make([]byte, 0)}, nxtPos + 1, nil
	}
	switch b {
	case '-':
		if nxtPos+1 == n || data[nxtPos+1] <= ' ' {
			return readYamlPartArray(data, nxtPos, indent, endChar)
		}
	case '{', '[':
		return readYamlFlowCollectionValue(data, nxtPos)
	case '}', ']':
		return nil, n, fmt.Errorf("Unexpected characters: %s ", string(data[nxtPos:nxtPos+1]))
	case '|', '?', ':':
		return nil, n, fmt.Errorf("%s is unimplemented for yaml yet", string(data[nxtPos:nxtPos+1]))
	}
	str, nextPos, err, tp := readYamlValue(data, nxtPos, ':', oldIndent)
	if err != nil {
		return nil, n, err
	}
	nexterPos, newIndent := readYamlNonEmptyLine(data, nextPos)
	if newIndent <= indent || nexterPos >= n || data[nexterPos] != ':' {
		return &dvevaluation.DvVariable{Kind: tp, Value: str}, nextPos, nil
	}
	return readYamlPartMap(data, nexterPos+1, indent, endChar, str)
}

func ReadYamlAsDvFieldInfo(data []byte) (*dvevaluation.DvVariable, error) {
	dvEntry, pos, err := readYamlAllFromIndent(data, 0, 0, byte(0))
	if err != nil {
		return nil, err
	}
	l := len(data)
	for ; pos < l; pos++ {
		if data[pos] > 32 {
			return nil, fmt.Errorf("Unexpected extra characters %s", getPositionErrorInfo(data, pos))
		}
	}
	if dvEntry == nil {
		dvEntry = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Value: make([]byte, 0)}
	}
	return dvEntry, nil
}
