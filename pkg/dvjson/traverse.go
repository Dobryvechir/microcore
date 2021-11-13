/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strconv"
	"strings"
)

func (item *DvFieldInfo) GetChildrenByRange(startIndex int, count int) *DvFieldInfo {
	if item == nil || (item.Kind != FIELD_ARRAY && item.Kind != FIELD_OBJECT) {
		return nil
	}
	subfields := item.Fields
	n := len(subfields)
	if n == 0 {
		return nil
	}
	for startIndex < 0 {
		startIndex += n
	}
	if startIndex > n {
		return nil
	}
	if count > n-startIndex {
		count = n - startIndex
	}
	res := &DvFieldInfo{Kind: FIELD_ARRAY, Fields: subfields[startIndex : startIndex+count]}
	return res
}

func (parseInfo *DvCrudParsingInfo) GetDvFieldInfoHierarchy() []*DvFieldInfo {
	n := len(parseInfo.Items)
	res := make([]*DvFieldInfo, n)
	for i := 0; i < n; i++ {
		res[i] = &parseInfo.Items[i].DvFieldInfo
	}
	return res
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
	if c == '[' || c == '{' {
		endPos, err := dvtextutils.ReadInsideBrackets(childName, pos)
		if err != nil {
			return nil, err
		}
		endPos++
		if endPos < n && childName[endPos] == '?' {
			strict = false
			endPos++
		}
		if c == '[' {
			data, err = ExpressionEvaluation(childName[pos+1:endPos], resolver)
		} else {
			item, err = FindItemByExpression(childName[pos+1:endPos], resolver, item, strict)
			fn = fnResolvedSign
		}
		if err != nil {
			return nil, err
		}
		childName = childName[endPos:]
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
				i, err = dvtextutils.ReadInsideBrackets(childName, i)
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
		if fn == fnResolvedSign {
			current = item
		} else {
			current, err = ExecuteProcessorFunction(fn, data, item)
			if err != nil {
				return nil, err
			}
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

func MeetItemExpression(expr string, item *DvFieldInfo, resolver ExpressionResolver, options int, index int) (bool, error) {
	props := ConvertDvFieldInfoToProperties(item, index)
	res, err := resolver(expr, props, options)
	if err != nil {
		return false, err
	}
	return res == "true", nil
}

func FindItemByExpression(expr string, resolver ExpressionResolver, item *DvFieldInfo, strict bool) (*DvFieldInfo, error) {
	if item == nil || item.Kind != FIELD_OBJECT && item.Kind != FIELD_ARRAY || item.Fields == nil || len(item.Fields) == 0 {
		if item == nil && strict {
			return nil, errors.New(expr + " of undefined")
		}
		return nil, nil
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		current := item.Fields[i]
		ok, err := MeetItemExpression(expr, current, resolver, EXPRESSION_RESOLVER_CACHE, i)
		if err != nil {
			return nil, err
		}
		if ok {
			return current, nil
		}
	}
	return nil, nil
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
	return resolver(expr, nil, 0)
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
	return dvtextutils.GetUnquotedString(subItem.GetStringValue())
}

func (item *DvFieldInfo) ReadChildStringValue(fieldName string) string {
	subItem, err := item.ReadChild(fieldName, nil)
	if err != nil || subItem == nil {
		return ""
	}
	return dvtextutils.GetUnquotedString(subItem.GetStringValue())
}


func ReadJsonChild(data []byte, childName string, rejectChildOfUndefined bool,props *dvevaluation.DvObject) (*DvFieldInfo, error) {
	item, err := JsonFullParser(data)
	if err != nil {
		return nil, err
	}
	item, err = item.ReadChildOfAnyLevel(childName, props)
	if err!=nil {
		if strings.HasPrefix(err.Error(),"undefined ") {
			if rejectChildOfUndefined {
				return nil, nil
			}
			return nil, errors.New("Read of undefined at " + childName)
		}
		return nil, err
	}
	return item, nil
}
