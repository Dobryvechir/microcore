/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvjson

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strconv"
	"strings"
)

type DvCrudDetails struct {
	storeFile      string
	lookedUrl      string
	initCrud       DvCrud
	userIdPossible bool
	userIdInteger  bool
	userIdName     []byte
	urlParts       []string
	dataUrlIndex   int
	autoIdPrefix   string
	table          *DvCrudParsingInfo
	idMap          map[string]int
	nextId         int
}

type DvFieldInfo struct {
	Name          []byte
	Value         []byte
	Kind          int
	Fields        []*DvFieldInfo
	valueStartPos int
	valueEndPos   int
	FieldStatus   int
	posStart      int
}

type DvCrudItem struct {
	DvFieldInfo
	itemBody []byte
	Id       []byte
	posEnd   int
}

type DvCrudParsingInfo struct {
	Items []*DvCrudItem
	Err   string
	Kind  int
	Value []byte
}

type ExpressionResolver func(string) (string, error)

const (
	OPTIONS_QUICK           = iota
	OPTIONS_FIELDS_ENTRIES  = iota
	OPTIONS_FIELDS_DETAILED = iota
	OPTIONS_ANY_OBJECT      = 1 << 7
	OPTIONS_LOW_MASK        = 1<<7 - 1
)

const (
	FIELD_EMPTY = iota
	FIELD_NULL
	FIELD_OBJECT
	FIELD_NUMBER
	FIELD_BOOLEAN
	FIELD_STRING
	FIELD_ARRAY
)

const (
	STATE_INITIAL = iota
	STATE_INITIAL_ARRAY_STARTED
	STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END
	STATE_ARRAY_STARTED
	STATE_OBJECT_STARTED
	STATE_OBJECT_COLON
	STATE_VALUE_STARTED
	STATE_WAITING_END
	STATE_OBJECT_COMMA_OR_END_EXPECTED
	STATE_ARRAY_COMMA_OR_END_EXPECTED
)

const (
	POST_STATE_NONE = iota
	POST_STATE_ARRAY_CLOSE
	POST_STATE_OBJECT_CLOSE
)

const (
	FIELD_STATUS_IS_ID = 1 << iota
)

var _stateDebugInfo = [...]string{"STATE_INITIAL", "STATE_INITIAL_ARRAY_STARTED", "STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END", "STATE_ARRAY_STARTED", "STATE_OBJECT_STARTED", "STATE_OBJECT_COLON", "STATE_VALUE_STARTED",
	"STATE_WAITING_END", "STATE_OBJECT_COMMA_OR_END_EXPECTED", "STATE_ARRAY_COMMA_OR_END_EXPECTED", "N"}

func jsonConvertToItems(body []byte, crudDetails *DvCrudDetails, options int) *DvCrudParsingInfo {
	return JsonQuickParser(body, crudDetails, options&OPTIONS_ANY_OBJECT != 0, options&OPTIONS_LOW_MASK)
}

func updateGroupOffset(details *DvFieldInfo, startOffset int, dif int) {
	if details.posStart > startOffset {
		details.posStart += dif
	}
	if details.valueStartPos > startOffset {
		details.valueStartPos += dif
	}
	if details.valueEndPos > startOffset {
		details.valueEndPos += dif
	}
	for _, v := range details.Fields {
		updateGroupOffset(v, startOffset, dif)
	}
}

func updateItemOffsets(item *DvCrudItem, startOffset int, dif int) {
	if dif == 0 {
		return
	}
	if item.posStart > startOffset {
		item.posStart += dif
	}
	if item.posEnd > startOffset {
		item.posEnd += dif
	}
	for _, v := range item.Fields {
		updateGroupOffset(v, startOffset, dif)
	}
}

func changeItemId(item *DvCrudItem, fieldValue []byte, crudInfo *DvCrudDetails) {
	knd := FIELD_STRING
	if crudInfo.userIdInteger {
		knd = FIELD_NUMBER
	}
	changeItemField(item, fieldValue, crudInfo.userIdName, knd, true)
}

func changeItemField(item *DvCrudItem, fieldValue []byte, fieldName []byte, fieldKind int, isIdField bool) int {
	fieldValueLen := len(fieldValue)
	dif := 0
	if isIdField {
		if bytes.Equal(item.Id, fieldValue) || fieldValueLen == 0 {
			return 0
		}
		item.Id = fieldValue
	}
	l := len(item.Fields)
	if l > 0 {
		for _, v := range item.Fields {
			if bytes.Equal(v.Name, fieldName) {
				v.Value = fieldValue
				dif = v.valueEndPos - v.valueStartPos - fieldValueLen
				res := make([]byte, 0, len(item.itemBody)+dif)
				res = append(res, item.itemBody[:v.valueStartPos]...)
				res = append(res, fieldValue...)
				res = append(res, item.itemBody[v.valueEndPos:]...)
				item.itemBody = res
				updateItemOffsets(item, v.valueStartPos, dif)
				return dif
			}
		}
	}
	bracket := "\""
	if fieldKind != FIELD_STRING {
		bracket = ""
	}
	insertion := "\"" + string(fieldName) + "\":" + bracket + string(fieldValue) + bracket
	valueEndOffset := len(insertion)
	pos := bytes.LastIndexByte(item.itemBody, '}')
	subPos := 0
	if pos < 0 {
		insertion = "{" + insertion + "}"
		item.Fields = nil
		dif = len(item.itemBody) - len(insertion)
		item.itemBody = []byte(insertion)
		pos = 0
	} else {
		if isObjectNotEmptyFromBack(item.itemBody, pos-1, '{') {
			insertion = "," + insertion
			subPos = 1
		}
		dif = len(insertion)
		res := make([]byte, 0, len(item.itemBody)+dif)
		res = append(res, item.itemBody[:pos]...)
		res = append(res, insertion...)
		res = append(res, item.itemBody[pos:]...)
		item.itemBody = res
	}
	valueStartPos := item.posStart + pos + subPos
	updateItemOffsets(item, valueStartPos, dif)
	newField := &DvFieldInfo{posStart: strings.LastIndex(insertion, "\":\"") + valueStartPos - 1, Name: fieldName, Value: fieldValue,
		valueStartPos: valueStartPos, valueEndPos: valueStartPos + valueEndOffset, Kind: fieldKind}
	item.Fields = append(item.Fields, newField)
	return dif
}

func isObjectNotEmptyFromBack(data []byte, pos int, startByte byte) bool {
	for pos >= 0 {
		if data[pos] == startByte {
			return false
		}
		if data[pos] > 32 {
			return true
		}
		pos--
	}
	return false
}

func jsonWholeItems(info *DvCrudParsingInfo) []byte {
	l := len(info.Items)
	extra := info.Kind == FIELD_ARRAY || l != 1
	size := 0
	if extra {
		size = l + 1
		if l == 0 {
			size = 2
		}
	}
	for _, v := range info.Items {
		size += len(v.itemBody)
	}
	b := make([]byte, size)
	pos := 0
	if extra {
		b[0] = '['
		b[size-1] = ']'
		pos = 1
	}
	for k, v := range info.Items {
		if k != 0 {
			b[pos] = ','
			pos++
		}
		for _, m := range v.itemBody {
			b[pos] = m
			pos++
		}
	}
	return b
}

func makeIdStatistics(crudInfo *DvCrudDetails) {
	crudInfo.nextId = 1
	crudInfo.idMap = make(map[string]int)
	items := crudInfo.table.Items
	l := len(crudInfo.autoIdPrefix)
	for i, item := range items {
		id := string(item.Id)
		crudInfo.idMap[id] = i
		if l == 0 || id[:l] == crudInfo.autoIdPrefix {
			n, err := strconv.Atoi(id[l:])
			if err == nil && n >= crudInfo.nextId {
				crudInfo.nextId = n + 1
			}
		}
	}
}

func jsonCheckWord(str []byte, pos int, maxOffset int, sample string, kind int) (string, int, int, int, int) {
	l := len(sample)
	if pos+l > maxOffset {
		return getExpectedButFound(sample, 255, pos, str), maxOffset, maxOffset, maxOffset, FIELD_EMPTY
	}
	for i := 0; i < l; i++ {
		c1 := str[i+pos]
		c2 := sample[i]
		if c1 != c2 && c1 != c2+32 {
			return getExpectedButFound(sample, c1, pos, str), maxOffset, maxOffset, maxOffset, FIELD_EMPTY
		}
	}
	return "", pos, pos + l, pos + l, kind
}

func jsonReadValue(str []byte, pos int, maxOffset int) (string, int, int, int, int) {
	for pos < maxOffset && str[pos] <= 32 {
		pos++
	}
	if pos >= maxOffset {
		return "", pos, pos, pos, FIELD_EMPTY
	}
	c := str[pos]
	startPos := pos
	if c == '"' {
		startPos++
		pos++
		for pos < maxOffset {
			c = str[pos]
			if c == '"' {
				return "", startPos, pos, pos + 1, FIELD_STRING
			}
			if c == '\\' {
				pos++
			}
			pos++
		}
		return getExpectedButFound("end quote", 255, pos, str), maxOffset, maxOffset, maxOffset, FIELD_EMPTY
	} else if c == 'n' || c == 'N' {
		return jsonCheckWord(str, pos, maxOffset, "null", FIELD_NULL)
	} else if c == 'f' || c == 'F' {
		return jsonCheckWord(str, pos, maxOffset, "false", FIELD_BOOLEAN)
	} else if c == 't' || c == 'T' {
		return jsonCheckWord(str, pos, maxOffset, "true", FIELD_BOOLEAN)
	} else if c == '+' || c == '-' || c == '.' || c >= '0' && c <= '9' {
		for pos < maxOffset {
			c = str[pos]
			if c == '+' || c == '-' || c == '.' || c >= '0' && c <= '9' || c == 'e' || c == 'E' {
				pos++
			} else {
				break
			}
		}
		return "", startPos, pos, pos, FIELD_NUMBER
	}
	return getExpectedButFound("value", c, pos, str), maxOffset, maxOffset, maxOffset, FIELD_EMPTY
}

func getErrorSample(main string, pos int, body []byte) string {
	pre := pos - 30
	post := pos + 30
	l := len(body)
	if pre < 0 {
		pre = 0
	}
	if post > l {
		post = l
	}
	if pos > l {
		pos = l
	}
	str := main + " at " + strconv.Itoa(pos) + " (" + string(body[pre:pos]) + "???" + string(body[pos:post]) + ")"
	str = strings.Replace(str, "\"", "&qt;", -1)
	str = strings.Replace(str, "\\", "&#92;", -1)
	return str
}

func getExpectedButFound(expected string, c byte, pos int, body []byte) string {
	mes := "Expected " + expected + ", but found "
	if c < 255 {
		if (c > 34 || c == 33) && c != '\\' {
			mes += string([]byte{c})
		}
		mes += "(" + strconv.Itoa(int(c)) + ")"
	} else {
		mes += "the end "
	}
	return getErrorSample(mes, pos, body)
}

func presentCurrentState(place string, stack []*DvFieldInfo, stackSign []int, level int, levelDif int, state int, postState int, i int, c byte, currentItem *DvCrudItem) {
	log.Printf("%s level: %d dif: %d state: %s postState: %d i:%d c:%c(%d) StackSign: %v Stack: %v %s", place, level, levelDif, _stateDebugInfo[state], postState, i, int(c), int(c), stackSign[:level+levelDif+1], stack[:level+1], logInfoForItem(currentItem))
}

func JsonQuickParser(body []byte, crudDetails *DvCrudDetails, highLevelObject bool, mode int) *DvCrudParsingInfo {
	stack := make([]*DvFieldInfo, 20, 20)
	stackSign := make([]int, 20, 20)
	l := len(body)
	level := -1
	levelDif := 0
	state := STATE_INITIAL
	postState := POST_STATE_NONE
	r := &DvCrudParsingInfo{}
	currentItem := &DvCrudItem{}
	isCurrentFieldId := false
	for i := 0; i < l; i++ {
		c := body[i]
		if c <= 32 {
			continue
		}
		switch state {
		case STATE_OBJECT_STARTED:
			if c == '}' {
				postState = POST_STATE_OBJECT_CLOSE
			} else if c == '"' {
				err, startPos, endPos, nxtPos, _ := jsonReadValue(body, i, l)
				if err != "" {
					r.Err = err
					return r
				}
				key := body[startPos:endPos]
				i = nxtPos - 1
				state = STATE_OBJECT_COLON
				if level == 0 {
					isCurrentFieldId = bytes.Equal(key, crudDetails.userIdName)
				}
				if (mode == OPTIONS_FIELDS_ENTRIES || isCurrentFieldId) && level == 0 || mode >= OPTIONS_FIELDS_DETAILED {
					currentField := &DvFieldInfo{Name: key, posStart: nxtPos}
					if isCurrentFieldId {
						currentField.FieldStatus |= FIELD_STATUS_IS_ID
					}
					if level == 0 {
						currentItem.Fields = append(currentItem.Fields, currentField)
						currentField = currentItem.Fields[len(currentItem.Fields)-1]
					} else {
						stack[level-1].Fields = append(stack[level-1].Fields, currentField)
						currentField = stack[level-1].Fields[len(stack[level-1].Fields)-1]
					}
					stack[level] = currentField
				}
			} else {
				r.Err = getExpectedButFound("} or opening quote", c, i, body)
				return r
			}
		case STATE_INITIAL_ARRAY_STARTED:
			if c == ']' {
				state = STATE_WAITING_END
			} else if c == '{' {
				state = STATE_OBJECT_STARTED
				currentItem = &DvCrudItem{DvFieldInfo: DvFieldInfo{Kind: FIELD_OBJECT, posStart: i}}
				level = 0
				stackSign[levelDif] = FIELD_OBJECT
			} else {
				if highLevelObject {
					r.Err = getExpectedButFound("] or { ", c, i, body)
					return r
				} else if c == '[' {
					state = STATE_VALUE_STARTED
					currentItem = &DvCrudItem{DvFieldInfo: DvFieldInfo{Kind: FIELD_ARRAY, posStart: i}}
					level = 0
					stackSign[levelDif] = FIELD_ARRAY
				} else {
					err, startPos, endPos, nxtPos, kind := jsonReadValue(body, i, l)
					if err != "" {
						r.Err = err
					} else {
						val := body[startPos:endPos]
						currentItem = &DvCrudItem{DvFieldInfo: DvFieldInfo{Kind: kind, posStart: startPos, Value: val}, itemBody: val, Id: val, posEnd: endPos}
						state = STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END
						r.Items = append(r.Items, currentItem)
						i = nxtPos - 1
					}
				}

			}
		case STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END:
			if c == ',' {
				state = STATE_INITIAL_ARRAY_STARTED
			} else if c == ']' {
				state = STATE_WAITING_END
			} else {
				r.Err = getExpectedButFound("] or , ", c, i, body)
				return r
			}
		case STATE_VALUE_STARTED:
			previousKind := stackSign[level+levelDif]
			if previousKind == FIELD_ARRAY && c != ']' && (mode >= OPTIONS_FIELDS_DETAILED || mode == OPTIONS_FIELDS_ENTRIES && level == 0) {
				indexCurrent := 0
				currentField := &DvFieldInfo{posStart: i}
				if level == 0 {
					indexCurrent = len(currentItem.Fields)
					currentItem.Fields = append(currentItem.Fields, currentField)
					currentField = currentItem.Fields[indexCurrent]
				} else {
					indexCurrent = len(stack[level-1].Fields)
					stack[level-1].Fields = append(stack[level-1].Fields, currentField)
					currentField = stack[level-1].Fields[indexCurrent]
				}
				stack[level] = currentField
				currentField.Name = []byte(strconv.Itoa(indexCurrent))
			}
			if c == '{' || c == '[' {
				var knd int
				if c == '{' {
					state = STATE_OBJECT_STARTED
					knd = FIELD_OBJECT
				} else {
					knd = FIELD_ARRAY
				}
				if level+levelDif >= cap(stackSign) {
					stackSign = append(stackSign, 0)
				}
				if mode >= OPTIONS_FIELDS_DETAILED || mode == OPTIONS_FIELDS_ENTRIES && level == 0 {
					if level >= cap(stack) {
						stack = append(stack, nil)
					}
				}
				if stack[level] != nil {
					stack[level].valueStartPos = i
					stack[level].valueEndPos = i
					stack[level].Kind = knd
				}
				level++
				stack[level] = nil
				stackSign[level+levelDif] = knd
			} else if c == ']' && stackSign[level+levelDif] == FIELD_ARRAY {
				postState = POST_STATE_ARRAY_CLOSE
			} else {
				err, startPos, endPos, nxtPos, kind := jsonReadValue(body, i, l)
				if err != "" {
					r.Err = err
					return r
				}
				if stack[level] != nil {
					stack[level].valueStartPos = startPos
					stack[level].valueEndPos = endPos
					stack[level].Kind = kind
					stack[level].Value = body[startPos:endPos]
					if isCurrentFieldId && level == 0 {
						currentItem.Id = stack[level].Value
					}
				}
				i = nxtPos - 1
				if previousKind == FIELD_OBJECT {
					state = STATE_OBJECT_COMMA_OR_END_EXPECTED
				} else {
					state = STATE_ARRAY_COMMA_OR_END_EXPECTED
				}
			}
		case STATE_OBJECT_COLON:
			if c == ':' {
				state = STATE_VALUE_STARTED
			} else {
				r.Err = getExpectedButFound(":", c, i, body)
				return r
			}
		case STATE_OBJECT_COMMA_OR_END_EXPECTED:
			if c == ',' {
				state = STATE_OBJECT_STARTED
			} else if c == '}' {
				postState = POST_STATE_OBJECT_CLOSE
			} else {
				r.Err = getExpectedButFound(", or }", c, i, body)
				return r
			}
		case STATE_ARRAY_COMMA_OR_END_EXPECTED:
			if c == ',' {
				state = STATE_VALUE_STARTED
			} else if c == ']' {
				postState = POST_STATE_ARRAY_CLOSE
			} else {
				r.Err = getExpectedButFound(", or ]", c, i, body)
				return r
			}
		case STATE_INITIAL:
			if c == '[' {
				state = STATE_INITIAL_ARRAY_STARTED
				r.Kind = FIELD_ARRAY
				levelDif = 1
				stackSign[0] = FIELD_ARRAY
			} else if c == '{' {
				state = STATE_OBJECT_STARTED
				r.Kind = FIELD_OBJECT
				currentItem = &DvCrudItem{DvFieldInfo: DvFieldInfo{Kind: FIELD_OBJECT, posStart: i}}
				level = 0
				levelDif = 0
				stackSign[0] = FIELD_OBJECT
			} else {
				if highLevelObject {
					r.Err = getExpectedButFound("[ or {", c, i, body)
					return r
				} else {
					err, startPos, endPos, nxtPos, kind := jsonReadValue(body, i, l)
					r.Err = err
					r.Kind = kind
					r.Value = body[startPos:endPos]
					state = STATE_WAITING_END
					i = nxtPos - 1
				}
			}
		case STATE_WAITING_END:
			r.Err = getExpectedButFound("End", c, i, body)
			return r
		}
		if LogJson && dvlog.CurrentLogLevel >= dvlog.LogTrace {
			presentCurrentState("SP", stack, stackSign, level, levelDif, state, postState, i, c, currentItem)
		}
		if postState != POST_STATE_NONE {
			switch postState {
			case POST_STATE_ARRAY_CLOSE:
				stack[level] = nil
				level--
				if level < 0 {
					currentItem.posEnd = i + 1
					currentItem.itemBody = body[currentItem.posStart:currentItem.posEnd]
					r.Items = append(r.Items, currentItem)
					state = STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END
				} else {
					if stack[level] != nil {
						stack[level].valueEndPos = i + 1
						stack[level].Value = body[stack[level].valueStartPos:stack[level].valueEndPos]
					}
					if stackSign[level+levelDif] == FIELD_OBJECT {
						state = STATE_OBJECT_COMMA_OR_END_EXPECTED
					} else {
						state = STATE_ARRAY_COMMA_OR_END_EXPECTED
					}
				}
			case POST_STATE_OBJECT_CLOSE:
				stack[level] = nil
				level--
				if level < 0 {
					currentItem.posEnd = i + 1
					currentItem.itemBody = body[currentItem.posStart:currentItem.posEnd]
					r.Items = append(r.Items, currentItem)

					if levelDif == 0 {
						state = STATE_WAITING_END
					} else {
						state = STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END
					}
				} else {
					if stack[level] != nil {
						stack[level].valueEndPos = i + 1
						stack[level].Value = body[stack[level].valueStartPos:stack[level].valueEndPos]
					}
					if stackSign[level+levelDif] == FIELD_OBJECT {
						state = STATE_OBJECT_COMMA_OR_END_EXPECTED
					} else {
						state = STATE_ARRAY_COMMA_OR_END_EXPECTED
					}
				}
			}
			postState = POST_STATE_NONE
		}
	}
	if state != STATE_INITIAL && state != STATE_WAITING_END {
		r.Err = "Unexpected end of json with " + strconv.Itoa(state)
	}
	return r
}

func ConvertDvCrudItemToDvFieldInfo(src *DvCrudItem) (dst *DvFieldInfo) {
	dst = &DvFieldInfo{
		Name:          src.Name,
		Value:         src.Value,
		Kind:          src.Kind,
		Fields:        src.Fields,
		valueStartPos: src.valueStartPos,
		valueEndPos:   src.valueEndPos,
		FieldStatus:   src.FieldStatus,
		posStart:      src.posStart,
	}
	return
}

func ConvertDvCrudParsingInfoToDvFieldInfo(crudParsingInfo *DvCrudParsingInfo) *DvFieldInfo {
	res := &DvFieldInfo{
		Kind:  crudParsingInfo.Kind,
		Value: crudParsingInfo.Value,
	}
	switch res.Kind {
	case FIELD_ARRAY:
		n := len(crudParsingInfo.Items)
		res.Fields = make([]*DvFieldInfo, n)
		for i := 0; i < n; i++ {
			res.Fields[i] = ConvertDvCrudItemToDvFieldInfo(crudParsingInfo.Items[i])
		}
	case FIELD_OBJECT:
		res = ConvertDvCrudItemToDvFieldInfo(crudParsingInfo.Items[0])
	}
	return res
}

func JsonFullParser(body []byte) (*DvFieldInfo, error) {
	crudDetails := &DvCrudDetails{}
	highLevelObject := false
	parsed := JsonQuickParser(body, crudDetails, highLevelObject, OPTIONS_FIELDS_DETAILED)
	if parsed.Err != "" {
		return nil, errors.New(parsed.Err)
	}
	return ConvertDvCrudParsingInfoToDvFieldInfo(parsed), nil
}

const MaxInt = int64(^uint(0) >> 1)

func ConvertByteArrayToIntOrDouble(data []byte) (interface{}, bool) {
	n := len(data)
	i := 0
	for ; i < n && data[i] <= ' '; i++ {
	}
	positive := true
	if data[i] == '-' {
		positive = false
		i++
	}
	if data[i] == '+' {
		i++
	}
	v := int64(0)
	for ; i < n; i++ {
		c := data[i]
		if c >= '0' && c <= '9' {
			v = v*10 + int64(c) - 48
		} else if c == '.' || c == 'e' || c == 'E' {
			break
		} else {
			return 0, false
		}
	}
	if i == n {
		if !positive {
			v = -v
		}
		if v <= MaxInt && v >= -MaxInt {
			return int(v), true
		}
		return float64(v), true
	}
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return 0, false
	}
	vf := int(f)
	if float64(vf) == f {
		return vf, true
	}
	return f, true
}

func ConvertSimpleKindAndValueToInterface(kind int, data []byte) (interface{}, bool) {
	switch kind {
	case FIELD_EMPTY:
		return nil, true
	case FIELD_NULL:
		return nil, true
	case FIELD_NUMBER:
		return ConvertByteArrayToIntOrDouble(data)
	case FIELD_BOOLEAN:
		return len(data) != 0 && (data[0] == 't' || data[0] == 'T'), true
	case FIELD_STRING:
		return string(data), true
	}
	return nil, false
}

func (parseInfo *DvCrudParsingInfo) ConvertSimpleValueToInterface() (interface{}, bool) {
	return ConvertSimpleKindAndValueToInterface(parseInfo.Kind, parseInfo.Value)
}

func (item *DvFieldInfo) ConvertSimpleValueToInterface() (interface{}, bool) {
	return ConvertSimpleKindAndValueToInterface(item.Kind, item.Value)
}

func (item *DvFieldInfo) ReadSimpleStringMap(data map[string]string) error {
	if item.Kind != FIELD_OBJECT {
		return errors.New(string(item.Name) + " must be an object { }")
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		p := item.Fields[i]
		k := string(p.Name)
		v := string(p.Value)
		data[k] = v
	}
	return nil
}

func (item *DvFieldInfo) ReadSimpleStringList(data []string) ([]string, error) {
	if item.Kind != FIELD_ARRAY {
		return data, errors.New(string(item.Name) + " must be an object { }")
	}
	n := len(item.Fields)
	if data == nil {
		data = make([]string, 0, n)
	}
	for i := 0; i < n; i++ {
		p := item.Fields[i]
		v := string(p.Value)
		data = append(data, v)
	}
	return data, nil
}

func (item *DvFieldInfo) ReadSimpleString() (string, error) {
	if item.Kind == FIELD_OBJECT || item.Kind == FIELD_ARRAY {
		return "[]", errors.New(string(item.Name) + " must be a simple type")
	}
	return string(item.Value), nil
}

func (item *DvFieldInfo) ConvertValueToInterface() (interface{}, bool) {
	switch item.Kind {
	case FIELD_ARRAY:
		fields := item.Fields
		n := len(fields)
		data := make([]interface{}, n)
		var ok bool
		for i := 0; i < n; i++ {
			data[i], ok = fields[i].ConvertValueToInterface()
			if !ok {
				return nil, false
			}
		}
		return data, true
	case FIELD_OBJECT:
		fields := item.Fields
		n := len(fields)
		data := make(map[string]interface{}, n)
		var ok bool
		for i := 0; i < n; i++ {
			data[string(fields[i].Name)], ok = fields[i].ConvertValueToInterface()
			if !ok {
				return nil, false
			}
		}
		return data, true
	}
	return item.ConvertSimpleValueToInterface()
}

func (parseInfo *DvCrudParsingInfo) GetDvFieldInfoHierarchy() []*DvFieldInfo {
	n := len(parseInfo.Items)
	res := make([]*DvFieldInfo, n)
	for i := 0; i < n; i++ {
		res[i] = &parseInfo.Items[i].DvFieldInfo
	}
	return res
}

func (item *DvFieldInfo) ReadSimpleChild(fieldName string) *DvFieldInfo {
	n := len(item.Fields)
	if item.Kind == FIELD_ARRAY {
		k, err := strconv.Atoi(fieldName)
		if err != nil || k < 0 || k >= n {
			return nil
		}
		return item.Fields[k]
	}
	if item.Kind != FIELD_OBJECT {
		return nil
	}
	name := []byte(fieldName)
	for i := 0; i < n; i++ {
		if bytes.Equal(item.Fields[i].Name, name) {
			return item.Fields[i]
		}
	}
	return nil
}

func (item *DvFieldInfo) ReadSimpleChildValue(fieldName string) string {
	subItem := item.ReadSimpleChild(fieldName)
	if subItem == nil {
		return ""
	}
	return dvparser.GetUnquotedString(subItem.GetStringValue())
}

func (item *DvFieldInfo) ReadChildStringValue(fieldName string) string {
	subItem, err := item.ReadChild(fieldName, nil)
	if err != nil || subItem == nil {
		return ""
	}
	return dvparser.GetUnquotedString(subItem.GetStringValue())
}

func ExpressionEvaluation(expr string, resolver ExpressionResolver) (string, error) {
	expr = strings.TrimSpace(expr)
	n := len(expr)
	if n != 0 {
		if IsJsonNumber([]byte(expr)) {
			return expr, nil
		}
		if IsJsonString([]byte(expr)) {
			return expr[1 : n-1], nil
		}
	}
	if resolver == nil {
		return "", fmt.Errorf("Expression %s is complex and there is no expression resolver", expr)
	}
	return resolver(expr)
}

func (item *DvFieldInfo) ReadChild(childName string, resolver ExpressionResolver) (res *DvFieldInfo, err error) {
	n := len(childName)
	if n == 0 {
		return item, nil
	}
	strict := true
	pos := 0
	c := childName[pos]
	data := childName
	fn := ""
	if c == '.' {
		pos++
		if pos == n {
			return item, nil
		}
		c = childName[pos]
	}
	if c == '[' {
		endPos, err := dvparser.ReadInsideBrackets(childName, pos)
		if err != nil {
			return nil, err
		}
		data, err = ExpressionEvaluation(childName[pos+1:endPos], resolver)
		if err != nil {
			return nil, err
		}
		childName = childName[endPos+1:]
	} else {
		i := pos
		for ; i < n; i++ {
			c := childName[i]
			if c == '.' || c == '[' {
				break
			}
			if c == '(' {
				fn = strings.TrimSpace(childName[pos:i])
				if fn == "" {
					return nil, fmt.Errorf("Empty function name before ( in %s at %d", childName, i)
				}
				pos = i + 1
				var err error
				i, err = dvparser.ReadInsideBrackets(childName, i)
				if err != nil {
					return nil, err
				}
				break
			}
		}
		data = childName[pos:i]
		if fn != "" {
			i++
		}
		childName = childName[i:]
		n = len(data)
		if n > 0 && data[n-1] == '?' {
			data = data[:n-1]
			strict = false
		}
	}
	var current *DvFieldInfo
	if fn != "" {
		current, err = ExecuteProcessorFunction(fn, data, item)
		if err != nil {
			return nil, err
		}
	} else {
		current = item.ReadSimpleChild(data)
	}
	if childName == "" || current == nil && !strict {
		return current, nil
	}
	if current == nil {
		return nil, fmt.Errorf("Cannot read %s of undefined in %s", childName, data)
	}
	return current.ReadChild(childName, resolver)
}

func (item *DvFieldInfo) GetStringValue() string {
	if item == nil || item.Kind == FIELD_EMPTY {
		return ""
	}
	switch item.Kind {
	case FIELD_OBJECT:
		res := "{"
		subfields := item.Fields
		n := len(subfields)
		for i := 0; i < n; i++ {
			if i != 0 {
				res += ","
			}
			res += "\"" + string(subfields[i].Name) + "\":"
			data := subfields[i].GetStringValue()
			res += data
		}
		return res + "}"
	case FIELD_ARRAY:
		res := "["
		subfields := item.Fields
		n := len(subfields)
		for i := 0; i < n; i++ {
			if i != 0 {
				res += ","
			}
			data := subfields[i].GetStringValue()
			res += data
		}
		return res + "]"
	case FIELD_STRING:
		return QuoteEscapedJsonBytesToString(item.Value)
	case FIELD_NULL:
		return "null"
	}
	return string(item.Value)
}
