/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvjson

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
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

type DvFieldInfoExtra struct {
	valueStartPos int
	valueEndPos   int
	FieldStatus   int
	posStart      int
}

type DvCrudItem struct {
	dvevaluation.DvVariable
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

const (
	OPTIONS_QUICK           = iota
	OPTIONS_FIELDS_ENTRIES  = iota
	OPTIONS_FIELDS_DETAILED = iota
	OPTIONS_ANY_OBJECT      = 1 << 7
	OPTIONS_LOW_MASK        = 1<<7 - 1
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

func updateGroupOffset(details *dvevaluation.DvVariable, startOffset int, dif int) {
	if details.Extra.(*DvFieldInfoExtra).posStart > startOffset {
		details.Extra.(*DvFieldInfoExtra).posStart += dif
	}
	if details.Extra.(*DvFieldInfoExtra).valueStartPos > startOffset {
		details.Extra.(*DvFieldInfoExtra).valueStartPos += dif
	}
	if details.Extra.(*DvFieldInfoExtra).valueEndPos > startOffset {
		details.Extra.(*DvFieldInfoExtra).valueEndPos += dif
	}
	for _, v := range details.Fields {
		updateGroupOffset(v, startOffset, dif)
	}
}

func updateItemOffsets(item *DvCrudItem, startOffset int, dif int) {
	if dif == 0 {
		return
	}
	if item.Extra.(*DvFieldInfoExtra).posStart > startOffset {
		item.Extra.(*DvFieldInfoExtra).posStart += dif
	}
	if item.posEnd > startOffset {
		item.posEnd += dif
	}
	for _, v := range item.Fields {
		updateGroupOffset(v, startOffset, dif)
	}
}

func changeItemId(item *DvCrudItem, fieldValue []byte, crudInfo *DvCrudDetails) {
	knd := dvevaluation.FIELD_STRING
	if crudInfo.userIdInteger {
		knd = dvevaluation.FIELD_NUMBER
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
				dif = v.Extra.(*DvFieldInfoExtra).valueEndPos - v.Extra.(*DvFieldInfoExtra).valueStartPos - fieldValueLen
				res := make([]byte, 0, len(item.itemBody)+dif)
				res = append(res, item.itemBody[:v.Extra.(*DvFieldInfoExtra).valueStartPos]...)
				res = append(res, fieldValue...)
				res = append(res, item.itemBody[v.Extra.(*DvFieldInfoExtra).valueEndPos:]...)
				item.itemBody = res
				updateItemOffsets(item, v.Extra.(*DvFieldInfoExtra).valueStartPos, dif)
				return dif
			}
		}
	}
	bracket := "\""
	if fieldKind != dvevaluation.FIELD_STRING {
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
	valueStartPos := item.Extra.(*DvFieldInfoExtra).posStart + pos + subPos
	updateItemOffsets(item, valueStartPos, dif)
	newField := &dvevaluation.DvVariable{
		Name:  fieldName,
		Value: fieldValue,
		Kind:  fieldKind,
		Extra: &DvFieldInfoExtra{
			posStart:      strings.LastIndex(insertion, "\":\"") + valueStartPos - 1,
			valueStartPos: valueStartPos,
			valueEndPos:   valueStartPos + valueEndOffset,
		},
	}
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
	extra := info.Kind == dvevaluation.FIELD_ARRAY || l != 1
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
		return getExpectedButFound(sample, 255, pos, str), maxOffset, maxOffset, maxOffset, dvevaluation.FIELD_UNDEFINED
	}
	for i := 0; i < l; i++ {
		c1 := str[i+pos]
		c2 := sample[i]
		if c1 != c2 && c1 != c2+32 {
			return getExpectedButFound(sample, c1, pos, str), maxOffset, maxOffset, maxOffset, dvevaluation.FIELD_UNDEFINED
		}
	}
	return "", pos, pos + l, pos + l, kind
}

func jsonReadValue(str []byte, pos int, maxOffset int) (string, int, int, int, int) {
	for pos < maxOffset && str[pos] <= 32 {
		pos++
	}
	if pos >= maxOffset {
		return "", pos, pos, pos, dvevaluation.FIELD_UNDEFINED
	}
	c := str[pos]
	startPos := pos
	if c == '"' {
		startPos++
		pos++
		for pos < maxOffset {
			c = str[pos]
			if c == '"' {
				return "", startPos, pos, pos + 1, dvevaluation.FIELD_STRING
			}
			if c == '\\' {
				pos++
			}
			pos++
		}
		return getExpectedButFound("end quote", 255, pos, str), maxOffset, maxOffset, maxOffset, dvevaluation.FIELD_UNDEFINED
	} else if c == 'n' || c == 'N' {
		return jsonCheckWord(str, pos, maxOffset, "null", dvevaluation.FIELD_NULL)
	} else if c == 'f' || c == 'F' {
		return jsonCheckWord(str, pos, maxOffset, "false", dvevaluation.FIELD_BOOLEAN)
	} else if c == 't' || c == 'T' {
		return jsonCheckWord(str, pos, maxOffset, "true", dvevaluation.FIELD_BOOLEAN)
	} else if c == '+' || c == '-' || c == '.' || c >= '0' && c <= '9' {
		for pos < maxOffset {
			c = str[pos]
			if c == '+' || c == '-' || c == '.' || c >= '0' && c <= '9' || c == 'e' || c == 'E' {
				pos++
			} else {
				break
			}
		}
		return "", startPos, pos, pos, dvevaluation.FIELD_NUMBER
	}
	return getExpectedButFound("value", c, pos, str), maxOffset, maxOffset, maxOffset, dvevaluation.FIELD_UNDEFINED
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

func presentCurrentState(place string, stack []*dvevaluation.DvVariable, stackSign []int, level int, levelDif int, state int, postState int, i int, c byte, currentItem *DvCrudItem) {
	log.Printf("%s level: %d dif: %d state: %s postState: %d i:%d c:%c(%d) StackSign: %v Stack: %v %s", place, level, levelDif, _stateDebugInfo[state], postState, i, int(c), int(c), stackSign[:level+levelDif+1], stack[:level+1], logInfoForItem(currentItem))
}

func JsonQuickParser(body []byte, crudDetails *DvCrudDetails, highLevelObject bool, mode int) *DvCrudParsingInfo {
	stack := make([]*dvevaluation.DvVariable, 20, 20)
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
					currentField := &dvevaluation.DvVariable{
						Name: key,
						Extra: &DvFieldInfoExtra{
							posStart: nxtPos,
						},
					}
					if isCurrentFieldId {
						currentField.Extra.(*DvFieldInfoExtra).FieldStatus |= FIELD_STATUS_IS_ID
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
				currentItem = &DvCrudItem{
					DvVariable: dvevaluation.DvVariable{
						Kind: dvevaluation.FIELD_OBJECT,
						Extra: &DvFieldInfoExtra{
							posStart: i,
						},
					}}
				level = 0
				stackSign[levelDif] = dvevaluation.FIELD_OBJECT
			} else {
				if highLevelObject {
					r.Err = getExpectedButFound("] or { ", c, i, body)
					return r
				} else if c == '[' {
					state = STATE_VALUE_STARTED
					currentItem = &DvCrudItem{
						DvVariable: dvevaluation.DvVariable{
							Kind: dvevaluation.FIELD_ARRAY,
							Extra: &DvFieldInfoExtra{
								posStart: i,
							},
						}}
					level = 0
					stackSign[levelDif] = dvevaluation.FIELD_ARRAY
				} else {
					err, startPos, endPos, nxtPos, kind := jsonReadValue(body, i, l)
					if err != "" {
						r.Err = err
					} else {
						val := body[startPos:endPos]
						currentItem = &DvCrudItem{
							DvVariable: dvevaluation.DvVariable{
								Kind:  kind,
								Value: val,
								Extra: &DvFieldInfoExtra{
									posStart: startPos,
								},
							},
							itemBody: val,
							Id:       val,
							posEnd:   endPos,
						}
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
			if previousKind == dvevaluation.FIELD_ARRAY && c != ']' && (mode >= OPTIONS_FIELDS_DETAILED || mode == OPTIONS_FIELDS_ENTRIES && level == 0) {
				indexCurrent := 0
				currentField := &dvevaluation.DvVariable{
					Extra: &DvFieldInfoExtra{
						posStart: i,
					},
				}
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
					knd = dvevaluation.FIELD_OBJECT
				} else {
					knd = dvevaluation.FIELD_ARRAY
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
					stack[level].Extra.(*DvFieldInfoExtra).valueStartPos = i
					stack[level].Extra.(*DvFieldInfoExtra).valueEndPos = i
					stack[level].Kind = knd
				}
				level++
				stack[level] = nil
				stackSign[level+levelDif] = knd
			} else if c == ']' && stackSign[level+levelDif] == dvevaluation.FIELD_ARRAY {
				postState = POST_STATE_ARRAY_CLOSE
			} else {
				err, startPos, endPos, nxtPos, kind := jsonReadValue(body, i, l)
				if err != "" {
					r.Err = err
					return r
				}
				if stack[level] != nil {
					stack[level].Extra.(*DvFieldInfoExtra).valueStartPos = startPos
					stack[level].Extra.(*DvFieldInfoExtra).valueEndPos = endPos
					stack[level].Kind = kind
					stack[level].Value = body[startPos:endPos]
					if isCurrentFieldId && level == 0 {
						currentItem.Id = stack[level].Value
					}
				}
				i = nxtPos - 1
				if previousKind == dvevaluation.FIELD_OBJECT {
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
				r.Kind = dvevaluation.FIELD_ARRAY
				levelDif = 1
				stackSign[0] = dvevaluation.FIELD_ARRAY
			} else if c == '{' {
				state = STATE_OBJECT_STARTED
				r.Kind = dvevaluation.FIELD_OBJECT
				currentItem = &DvCrudItem{
					DvVariable: dvevaluation.DvVariable{
						Kind: dvevaluation.FIELD_OBJECT,
						Extra: &DvFieldInfoExtra{
							posStart: i,
						},
					},
				}
				level = 0
				levelDif = 0
				stackSign[0] = dvevaluation.FIELD_OBJECT
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
					currentItem.itemBody = body[currentItem.Extra.(*DvFieldInfoExtra).posStart:currentItem.posEnd]
					r.Items = append(r.Items, currentItem)
					state = STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END
				} else {
					if stack[level] != nil {
						stack[level].Extra.(*DvFieldInfoExtra).valueEndPos = i + 1
						stack[level].Value = body[stack[level].Extra.(*DvFieldInfoExtra).valueStartPos:stack[level].Extra.(*DvFieldInfoExtra).valueEndPos]
					}
					if stackSign[level+levelDif] == dvevaluation.FIELD_OBJECT {
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
					currentItem.itemBody = body[currentItem.Extra.(*DvFieldInfoExtra).posStart:currentItem.posEnd]
					r.Items = append(r.Items, currentItem)

					if levelDif == 0 {
						state = STATE_WAITING_END
					} else {
						state = STATE_INITIAL_ARRAY_WAITING_COMMA_OR_END
					}
				} else {
					if stack[level] != nil {
						stack[level].Extra.(*DvFieldInfoExtra).valueEndPos = i + 1
						stack[level].Value = body[stack[level].Extra.(*DvFieldInfoExtra).valueStartPos:stack[level].Extra.(*DvFieldInfoExtra).valueEndPos]
					}
					if stackSign[level+levelDif] == dvevaluation.FIELD_OBJECT {
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

func ConvertDvCrudItemToDvFieldInfo(src *DvCrudItem) (dst *dvevaluation.DvVariable) {
	dst = &dvevaluation.DvVariable{
		Name:   src.Name,
		Value:  src.Value,
		Kind:   src.Kind,
		Fields: src.Fields,
		Extra:  src.Extra,
	}
	return
}

func ConvertDvCrudParsingInfoToDvFieldInfo(crudParsingInfo *DvCrudParsingInfo) *dvevaluation.DvVariable {
	res := &dvevaluation.DvVariable{
		Kind:  crudParsingInfo.Kind,
		Value: crudParsingInfo.Value,
	}
	switch res.Kind {
	case dvevaluation.FIELD_ARRAY:
		n := len(crudParsingInfo.Items)
		res.Fields = make([]*dvevaluation.DvVariable, n)
		for i := 0; i < n; i++ {
			res.Fields[i] = ConvertDvCrudItemToDvFieldInfo(crudParsingInfo.Items[i])
		}
	case dvevaluation.FIELD_OBJECT:
		res = ConvertDvCrudItemToDvFieldInfo(crudParsingInfo.Items[0])
	}
	return res
}

func JsonFullParser(body []byte) (*dvevaluation.DvVariable, error) {
	crudDetails := &DvCrudDetails{}
	highLevelObject := false
	parsed := JsonQuickParser(body, crudDetails, highLevelObject, OPTIONS_FIELDS_DETAILED)
	if parsed.Err != "" {
		return nil, errors.New(parsed.Err)
	}
	return ConvertDvCrudParsingInfoToDvFieldInfo(parsed), nil
}

var jsonFullParserInited = dvevaluation.RegisterJsonFullParser(JsonFullParser)
