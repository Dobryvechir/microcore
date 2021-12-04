package dvevaluation

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strconv"
	"strings"
)

type ExpressionResolver func(string, map[string]interface{}, int) (string, error)

func (item *DvVariable) ReadChildOfAnyLevel(name string, props *DvObject) (res *DvVariable, err error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || item==nil {
		return item, nil
	}
	defProps:=make(map[string]interface{})
	if props==nil {
		return nil, errors.New("Props cannot be null")
	} else {
		props = NewObjectWithPrototype(defProps,props)
	}
	res, err = item.ReadChild(name, func(expr string, data map[string]interface{},options int) (string, error) {
		if data==nil {
			data = defProps
		}
		props.Properties = data
		return ParseForDvObjectString(expr, props)
	})
	return
}

func (field *DvVariable) AddStringField(key string, value string) bool {
	if field.Fields == nil {
		field.Fields = make([]*DvVariable, 0, 7)
	}
	field.Fields = append(field.Fields, &DvVariable{Kind: FIELD_STRING, Name: []byte(key), Value: []byte(value)})
	return true
}

func (field *DvVariable) AddField(item *DvVariable) bool {
	if field.Fields == nil {
		field.Fields = make([]*DvVariable, 0, 7)
	}
	field.Fields = append(field.Fields, item)
	return true
}

func (item *DvVariable) ContainsItemIn(v interface{}) bool {
	if item==nil {
		return false
	}
	s:=AnyToString(v)
	n:=len(item.Fields)
	switch item.Kind {
	case FIELD_OBJECT:
		for i:=0;i<n;i++ {
			if string(item.Fields[i].Name) == s {
				return true
			}
		}
	case FIELD_ARRAY:
		for i:=0;i<n;i++ {
			f:=item.Fields[i]
			if f.Kind!=FIELD_ARRAY && f.Kind!=FIELD_OBJECT && string(f.Value) == s {
				return true
			}
		}
	default:
		return strings.Contains(string(item.Value), s)
	}
	return false
}

func (item *DvVariable) FindDifferenceByQuickMap(other *DvVariable,
	fillAdded bool, fillRemoved bool, fillUpdated bool, fillUnchanged bool,
	fillUpdatedCounterpart bool, unchangedAsUpdated bool) (added *DvVariable, removed *DvVariable,
	updated *DvVariable, unchanged *DvVariable, counterparts *DvVariable) {
	if unchangedAsUpdated {
		fillUnchanged = false
	}
	kind := FIELD_ARRAY
	n1 := 0
	if item != nil {
		n1 = len(item.Fields)
		if item.Kind==FIELD_ARRAY || item.Kind==FIELD_OBJECT {
			kind = item.Kind
		}
	}
	n2 := 0
	if other != nil {
		n2 = len(other.Fields)
		if (other.Kind==FIELD_ARRAY || other.Kind==FIELD_OBJECT) && n1==0 {
			kind = other.Kind
		}
	}
	m := n1
	if m > n2 {
		m = n2
	}
	if fillAdded {
		added = &DvVariable{Fields: make([]*DvVariable, 0, n1), Kind: kind}
	}
	if fillRemoved {
		removed = &DvVariable{Fields: make([]*DvVariable, 0, n2), Kind: kind}
	}
	if fillUpdated {
		updated = &DvVariable{Fields: make([]*DvVariable, 0, m), Kind: kind}
	}
	if fillUnchanged {
		unchanged = &DvVariable{Fields: make([]*DvVariable, 0, m), Kind: kind}
	}
	if fillUpdatedCounterpart {
		counterparts = &DvVariable{Fields: make([]*DvVariable, 0, m), Kind: kind}
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

func (item *DvVariable) EvaluateDvFieldItem(expression string, env *DvObject) (bool, error) {
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
				v = StringToNumber(string(f.Value))
			case FIELD_NULL:
			case FIELD_UNDEFINED:
				v = nil
			}
			env.Set(k, v)
		}
	}
	res, err := env.EvaluateBooleanExpression(expression)
	return res, err
}

func (item *DvVariable) CompareWholeDvField(other *DvVariable) int {
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
	case FIELD_UNDEFINED:
	case FIELD_NULL:
		return 0
	case FIELD_NUMBER:
		if bytes.Equal(item.Value, other.Value) {
			return 0
		}
		n1 := AnyToNumber(string(item.Value))
		n2 := AnyToNumber(string(other.Value))
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

func (item *DvVariable) CompareDvFieldByFields(other *DvVariable, fields []string) int {
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

func (item *DvVariable) ReadSimpleChild(fieldName string) *DvVariable {
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

func (item *DvVariable) ReadSimpleChildValue(fieldName string) string {
	subItem := item.ReadSimpleChild(fieldName)
	if subItem == nil {
		return ""
	}
	return dvtextutils.GetUnquotedString(subItem.GetStringValue())
}

func (item *DvVariable) ReadChildStringValue(fieldName string) string {
	subItem, err := item.ReadChild(fieldName, nil)
	if err != nil || subItem == nil {
		return ""
	}
	return dvtextutils.GetUnquotedString(subItem.GetStringValue())
}

func (item *DvVariable) ReadPath(childName string, rejectChildOfUndefined bool, env *DvObject) (*DvVariable, error) {
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

func (item *DvVariable) ReadChild(childName string, resolver ExpressionResolver) (res *DvVariable, err error) {
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
	var current *DvVariable
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

func (item *DvVariable) GetChildrenByRange(startIndex int, count int) *DvVariable {
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
	res := &DvVariable{Kind: FIELD_ARRAY, Fields: subfields[startIndex : startIndex+count]}
	return res
}

func (item *DvVariable) GetStringValue() string {
	if item == nil || item.Kind == FIELD_UNDEFINED {
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
			data := subfields[i].GetStringValueJson()
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
			data := subfields[i].GetStringValueJson()
			res += data
		}
		return res + "]"
	case FIELD_STRING:
		return string(item.Value)
	case FIELD_NULL:
		return "null"
	}
	return string(item.Value)
}

func (item *DvVariable) GetStringValueJson() string {
	if item == nil || item.Kind == FIELD_UNDEFINED {
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
		return dvtextutils.QuoteEscapedJsonBytesToString(item.Value)
	case FIELD_NULL:
		return "null"
	}
	return string(item.Value)
}

func (item *DvVariable) ConvertSimpleValueToInterface() (interface{}, bool) {
	return ConvertSimpleKindAndValueToInterface(item.Kind, item.Value)
}

func (item *DvVariable) ReadSimpleStringMap(data map[string]string) error {
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

func (item *DvVariable) ReadSimpleStringList(data []string) ([]string, error) {
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

func (item *DvVariable) ReadSimpleString() (string, error) {
	if item.Kind == FIELD_OBJECT || item.Kind == FIELD_ARRAY {
		return "[]", errors.New(string(item.Name) + " must be a simple type")
	}
	return string(item.Value), nil
}

func (item *DvVariable) ConvertValueToInterface() (interface{}, bool) {
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

