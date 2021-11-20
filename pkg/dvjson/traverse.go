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

func (item *DvFieldInfo) ReadPath(childName string, rejectChildOfUndefined bool, env *dvevaluation.DvObject) (*DvFieldInfo, error) {
	item, err := item.ReadChildOfAnyLevel(childName, env)
	if err != nil {
		if strings.HasPrefix(err.Error(), "undefined ") {
			if rejectChildOfUndefined {
				return nil, nil
			}
			return nil, errors.New("Read of undefined at " + childName)
		}
		return nil, err
	}
	return item, nil

}

func ReadPathOfAny(item interface{}, childName string, rejectChildOfUndefined bool, env *dvevaluation.DvObject) (*DvFieldInfo, error) {
	switch item.(type) {
	case *DvFieldInfo:
		return item.(*DvFieldInfo).ReadPath(childName, rejectChildOfUndefined, env)
	}
	return nil, nil
}

func ReadJsonChild(data []byte, childName string, rejectChildOfUndefined bool, env *dvevaluation.DvObject) (*DvFieldInfo, error) {
	item, err := JsonFullParser(data)
	if err != nil {
		return nil, err
	}
	return item.ReadPath(childName, rejectChildOfUndefined, env)
}

func (item *DvFieldInfo) EvaluateDvFieldItem(expression string, env *dvevaluation.DvObject) (bool, error) {
	env.Set("this", item)
	if item != nil && len(item.Fields) > 0 {
		n := len(item.Fields)
		var v interface{}
		for i := 0; i < n; i++ {
			f := item.Fields[i]
			v = f
			k := string(f.Name)
			switch f.Kind {
			case FIELD_BOOLEAN:
				v = len(f.Value) == 4 && f.Value[0] == 't'
			case FIELD_STRING:
				v = string(f.Value)
			case FIELD_NUMBER:
				v = dvevaluation.StringToNumber(string(f.Value))
			case FIELD_NULL:
			case FIELD_EMPTY:
				v = nil
			}
			env.Set(k, v)
		}
	}
	res, err := env.EvaluateBooleanExpression(expression)
	return res, err
}
func (item *DvFieldInfo) CompareWholeDvField(other *DvFieldInfo) int {
	if item == nil {
		if other == nil {
			return 0
		}
		return -1
	}
	if other == nil {
		return 1
	}
	dif := item.Kind - other.Kind
	if dif != 0 {
		return dif
	}
	switch item.Kind {
	case FIELD_BOOLEAN:
		n1 := 0
		if len(item.Value) == 4 && item.Value[0] == 't' {
			n1 = 1
		}
		n2 := 0
		if len(other.Value) == 4 && other.Value[0] == 't' {
			n2 = 1
		}
		return n1 - n2
	case FIELD_STRING:
		return strings.Compare(string(item.Value), string(other.Value))
	case FIELD_EMPTY:
	case FIELD_NULL:
		return 0
	case FIELD_NUMBER:
		if bytes.Equal(item.Value, other.Value) {
			return 0
		}
		n1 := dvevaluation.AnyToNumber(string(item.Value))
		n2 := dvevaluation.AnyToNumber(string(other.Value))
		if n1 == n2 {
			return 0
		}
		if n1 < n2 {
			return -1
		}
		return 1
	case FIELD_ARRAY:
		n := len(item.Fields)
		dif := n - len(other.Fields)
		if dif != 0 {
			return dif
		}
		for i := 0; i < n; i++ {
			dif = item.Fields[i].CompareWholeDvField(other.Fields[i])
			if dif != 0 {
				return dif
			}
		}
		return 0
	case FIELD_OBJECT:
		n := len(item.Fields)
		dif := n - len(other.Fields)
		if dif != 0 {
			return dif
		}
		if item.QuickSearch == nil || len(item.QuickSearch.Looker) != n {
			item.CreateQuickInfoForObjectType()
		}
		if other.QuickSearch == nil || len(other.QuickSearch.Looker) != n {
			other.CreateQuickInfoForObjectType()
		}
		for k, v := range item.QuickSearch.Looker {
			v1, ok := other.QuickSearch.Looker[k]
			if !ok || v1 == nil {
				if v == nil {
					return 0
				}
				return 1
			}
			if v == nil {
				return -1
			}
			dif = v.CompareWholeDvField(v1)
			if dif != 0 {
				return dif
			}
		}
	}
	return 0
}

func (item *DvFieldInfo) CompareDvFieldByFields(other *DvFieldInfo, fields []string) int {
	item.CreateQuickInfoForObjectType()
	other.CreateQuickInfoForObjectType()
	n := len(fields)
	for i := 0; i < n; i++ {
		field := fields[i]
		if len(field) == 0 {
			continue
		}
		asc := true
		if field[0] == '~' {
			asc = false
			field = field[1:]
		}
		df1 := item.QuickSearch.Looker[field]
		df2 := other.QuickSearch.Looker[field]
		dif := df1.CompareWholeDvField(df2)
		if dif != 0 {
			if asc {
				return -dif
			}
			return dif
		}
	}
	return 0
}

func CountChildren(val interface{}) int {
	if val == nil {
		return 0
	}
	switch val.(type) {
	case *DvFieldInfo:
		return len(val.(*DvFieldInfo).Fields)
	}
	return 0
}

func (item *DvFieldInfo) FindDifferenceByQuickMap(other *DvFieldInfo,
	fillAdded bool, fillRemoved bool, fillUpdated bool, fillUnchanged bool,
	fillUpdatedCounterpart bool, unchangedAsUpdated bool) (added *DvFieldInfo, removed *DvFieldInfo,
	updated *DvFieldInfo, unchanged *DvFieldInfo, counterparts *DvFieldInfo) {
	n1 := 0
	if item != nil {
		n1 = len(item.Fields)
	}
	n2 := 0
	if other != nil {
		n2 = len(other.Fields)
	}
	m := n1
	if m > n2 {
		m = n2
	}
	if unchangedAsUpdated {
		fillUnchanged = false
	}
	if fillAdded {
		added = &DvFieldInfo{Fields: make([]*DvFieldInfo, 0, n1), Kind: FIELD_ARRAY}
	}
	if fillRemoved {
		removed = &DvFieldInfo{Fields: make([]*DvFieldInfo, 0, n2), Kind: FIELD_ARRAY}
	}
	if fillUpdated {
		updated = &DvFieldInfo{Fields: make([]*DvFieldInfo, 0, m), Kind: FIELD_ARRAY}
	}
	if fillUnchanged {
		unchanged = &DvFieldInfo{Fields: make([]*DvFieldInfo, 0, m), Kind: FIELD_ARRAY}
	}
	if fillUpdatedCounterpart {
		counterparts = &DvFieldInfo{Fields: make([]*DvFieldInfo, 0, m), Kind: FIELD_ARRAY}
	}
	if n2 == 0 {
		if n1 == 0 {
			return
		}
		added.Fields = append(added.Fields, item.Fields...)
		return
	}
	if n1 == 0 {
		removed.Fields = append(removed.Fields, other.Fields...)
		return
	}
	for k, v := range item.QuickSearch.Looker {
		v1, ok := other.QuickSearch.Looker[k]
		if ok {
			dif := 1
			if !unchangedAsUpdated {
				dif = v.CompareWholeDvField(v1)
			}
			if dif == 0 {
				if fillUnchanged {
					unchanged.Fields = append(unchanged.Fields, v)
				}
			} else {
				if fillUpdated {
					updated.Fields = append(updated.Fields, v)
				}
				if fillUpdatedCounterpart {
					counterparts.Fields = append(counterparts.Fields, v1)
				}
			}
		} else if fillAdded {
			added.Fields = append(added.Fields, v)
		}
	}
	if fillRemoved {
		for k, v := range other.QuickSearch.Looker {
			_, ok := item.QuickSearch.Looker[k]
			if !ok {
				removed.Fields = append(removed.Fields, v)
			}
		}
	}
	return
}

func FindDifferenceForAnyType(itemAny interface{}, otherAny interface{},
	fillAdded bool, fillRemoved bool, fillUpdated bool, fillUnchanged bool,
	fillUpdatedCounterpart bool, unchangedAsUpdated bool) (added *DvFieldInfo, removed *DvFieldInfo,
	updated *DvFieldInfo, unchanged *DvFieldInfo, counterparts *DvFieldInfo) {
	var item, other *DvFieldInfo
	switch itemAny.(type) {
	case *DvFieldInfo:
		item = itemAny.(*DvFieldInfo)
	}
	switch otherAny.(type) {
	case *DvFieldInfo:
		other = otherAny.(*DvFieldInfo)
	}
	if item.QuickSearch == nil {
		item.CreateQuickInfoForObjectType()
	}
	if other.QuickSearch == nil {
		other.CreateQuickInfoForObjectType()
	}
	return item.FindDifferenceByQuickMap(other, fillAdded, fillRemoved,
		fillUpdated, fillUnchanged, fillUpdatedCounterpart, unchangedAsUpdated)
}
