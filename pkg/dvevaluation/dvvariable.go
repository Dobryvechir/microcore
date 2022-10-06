package dvevaluation

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strconv"
	"strings"
)

type ExpressionResolver func(string, map[string]interface{}, int) (string, error)

func (item *DvVariable) ReadChildOfAnyLevel(name string, props *DvObject) (res *DvVariable, parent *DvVariable, err error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || item == nil {
		return item, nil, nil
	}
	defProps := make(map[string]interface{})
	if props == nil {
		return nil, nil, errors.New("Props cannot be null")
	} else {
		props = NewObjectWithPrototype(defProps, props)
	}
	res, parent, err = item.ReadChild(name, func(expr string, data map[string]interface{}, options int) (string, error) {
		if data == nil {
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
	if item == nil {
		return false
	}
	s := AnyToString(v)
	n := len(item.Fields)
	switch item.Kind {
	case FIELD_OBJECT, FIELD_FUNCTION:
		for i := 0; i < n; i++ {
			if item.Fields[i] != nil {
				b := item.Fields[i].Name
				if len(b) > 0 && string(b) == s {
					return true
				}
			}
		}
	case FIELD_ARRAY:
		for i := 0; i < n; i++ {
			f := item.Fields[i]
			if f.Kind != FIELD_ARRAY && f.Kind != FIELD_OBJECT && string(f.Value) == s {
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
	fillUpdatedCounterpart bool, unchangedAsUpdated bool, useIndex bool) (added *DvVariable, removed *DvVariable,
	updated *DvVariable, unchanged *DvVariable, counterparts *DvVariable) {
	if unchangedAsUpdated {
		fillUnchanged = false
	}
	kind := FIELD_ARRAY
	n1 := 0
	if item != nil {
		n1 = len(item.Fields)
		if item.Kind == FIELD_ARRAY || item.Kind == FIELD_OBJECT {
			kind = item.Kind
		}
	}
	n2 := 0
	if other != nil {
		n2 = len(other.Fields)
		if (other.Kind == FIELD_ARRAY || other.Kind == FIELD_OBJECT) && n1 == 0 {
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
	if useIndex && n2 > 0 {
		for i := 0; i < n2; i++ {
			f := other.Fields[i]
			if f == nil {
				f = &DvVariable{Kind: FIELD_NULL}
				other.Fields[i] = f
			}
			f.Extra = i
		}
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
		if v == nil {
			v = &DvVariable{Kind: FIELD_NULL}
		}
		if ok {
			dif := 1
			if !unchangedAsUpdated {
				dif = v.CompareWholeDvField(v1)
			}
			if dif == 0 {
				if fillUnchanged {
					unchanged.Fields = append(unchanged.Fields, v1)
				}
			} else {
				if fillUpdated {
					if useIndex {
						v.Extra = v1.Extra
					}
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
	case FIELD_OBJECT, FIELD_FUNCTION:
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
	if item == nil {
		return nil
	}
	n := len(item.Fields)
	if item.Kind == FIELD_ARRAY {
		k, err := strconv.Atoi(fieldName)
		if err != nil || k < 0 || k >= n {
			if item.Prototype != nil {
				return item.Prototype.ReadSimpleChild(fieldName)
			}
			return nil
		}
		return item.Fields[k]
	}
	if item.Kind != FIELD_OBJECT && item.Kind != FIELD_FUNCTION {
		return nil
	}
	name := []byte(fieldName)
	for i := 0; i < n; i++ {
		r := item.Fields[i]
		if r == nil {
			continue
		}
		if bytes.Equal(r.Name, name) {
			return r
		}
	}
	if item.Prototype != nil {
		return item.Prototype.ReadSimpleChild(fieldName)
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
	subItem, _, err := item.ReadChild(fieldName, nil)
	if err != nil || subItem == nil {
		return ""
	}
	return dvtextutils.GetUnquotedString(subItem.GetStringValue())
}

func (item *DvVariable) ReadChildStringArrayValue(fieldName string) []string {
	subItem, _, err := item.ReadChild(fieldName, nil)
	if err != nil || subItem == nil || subItem.Fields == nil {
		return nil
	}
	n := len(subItem.Fields)
	res := make([]string, n)
	for i := 0; i < n; i++ {
		s := subItem.Fields[i]
		if s == nil {
			res[i] = ""
		} else {
			res[i] = dvtextutils.GetUnquotedString(s.GetStringValue())
		}
	}
	return res
}

func (item *DvVariable) ReadChildMapValue(fieldName string) map[string]string {
	subItem, _, err := item.ReadChild(fieldName, nil)
	if err != nil || subItem == nil || subItem.Kind != FIELD_OBJECT && subItem.Kind != FIELD_FUNCTION || len(subItem.Fields) == 0 {
		return nil
	}
	res := make(map[string]string)
	n := len(subItem.Fields)
	for i := 0; i < n; i++ {
		field := subItem.Fields[i]
		if field == nil || len(field.Name) == 0 {
			continue
		}
		key := string(field.Name)
		val := dvtextutils.GetUnquotedString(field.GetStringValue())
		res[key] = val
	}
	return res
}

func (item *DvVariable) ReadPath(childName string, rejectChildOfUndefined bool, env *DvObject) (*DvVariable, *DvVariable, error) {
	item, parent, err := item.ReadChildOfAnyLevel(childName, env)
	if err != nil {
		if strings.HasPrefix(err.Error(), "undefined ") {
			if rejectChildOfUndefined {
				return nil, parent, nil
			}
			return nil, nil, errors.New("Read of undefined at " + childName)
		}
		return nil, nil, err
	}
	return item, parent, nil
}

func (item *DvVariable) ReadChild(childName string, resolver ExpressionResolver) (res *DvVariable, parent *DvVariable, err error) {
	n := len(childName)
	parent = item
	if n == 0 {
		return item, nil, nil
	}
	strict := true
	pos := 0
	c := childName[pos]
	data := childName
	fn := ""
	if c == '.' {
		pos++
		if pos == n {
			return item, nil, nil
		}
		c = childName[pos]
	}
	if c == '[' || c == '{' {
		endPos, err := dvtextutils.ReadInsideBrackets(childName, pos)
		if err != nil {
			return nil, nil, err
		}
		insideStr := childName[pos+1 : endPos]
		endPos++
		if endPos < n && childName[endPos] == '?' {
			strict = false
			endPos++
		}
		if c == '[' {
			data, err = ExpressionEvaluation(insideStr, resolver)
		} else {
			item, err = FindItemByExpression(insideStr, resolver, item, strict)
			fn = fnResolvedSign
		}
		if err != nil {
			return nil, nil, err
		}
		childName = childName[endPos:]
	} else {
		i := pos
		for ; i < n; i++ {
			c := childName[i]
			if c == '.' || c == '[' || c == '{' {
				break
			}
			if c == '(' {
				fn = strings.TrimSpace(childName[pos:i])
				if fn == "" {
					return nil, nil, fmt.Errorf("Empty function name before ( in %s at %d", childName, i)
				}
				pos = i + 1
				var err error
				i, err = dvtextutils.ReadInsideBrackets(childName, i)
				if err != nil {
					return nil, nil, err
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
				return nil, nil, err
			}
		}
	} else {
		current = item.ReadSimpleChild(data)
	}
	if childName == "" || current == nil && !strict {
		return current, parent, nil
	}
	if current == nil {
		return nil, nil, fmt.Errorf("Cannot read %s of undefined in %s", childName, data)
	}
	return current.ReadChild(childName, resolver)
}

func (item *DvVariable) GetChildrenByRange(startIndex int, count int) *DvVariable {
	if item == nil || (item.Kind != FIELD_ARRAY && item.Kind != FIELD_OBJECT && item.Kind != FIELD_FUNCTION) {
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
	case FIELD_OBJECT, FIELD_FUNCTION:
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
	case FIELD_OBJECT, FIELD_FUNCTION:
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
	if item.Kind != FIELD_OBJECT && item.Kind != FIELD_FUNCTION {
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
	case FIELD_OBJECT, FIELD_FUNCTION:
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

func (item *DvVariable) GetLength() int {
	if item == nil {
		return 0
	}
	switch item.Kind {
	case FIELD_ARRAY, FIELD_OBJECT, FIELD_FUNCTION:
		return len(item.Fields)
	case FIELD_STRING:
		return len(item.Value)
	}
	return 0
}

func (item *DvVariable) IndexOf(child *DvVariable) int {
	if item == nil {
		return -1
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		sub := item.Fields[i]
		if sub == child {
			return i
		}
	}
	if child == nil {
		child = &DvVariable{Kind: FIELD_NULL}
	}
	for i := 0; i < n; i++ {
		sub := item.Fields[i]
		if sub.Kind == child.Kind && bytes.Equal(sub.Name, child.Name) {
			if len(sub.Name) > 0 {
				return i
			}
			if sub.Kind != FIELD_OBJECT && sub.Kind != FIELD_ARRAY && sub.Kind != FIELD_FUNCTION {
				if bytes.Equal(sub.Value, child.Value) {
					return i
				}
			}
		}
	}
	return -1
}

func (item *DvVariable) RemoveChild(child *DvVariable) {
	n := item.IndexOf(child)
	if n >= 0 {
		item.Fields = append(item.Fields[:n], item.Fields[n+1:]...)
	}
}

func (item *DvVariable) FindIndex(key string) int {
	if item == nil {
		return -1
	}
	n := len(item.Fields)
	keyBytes := []byte(key)
	for i := 0; i < n; i++ {
		if item.Fields[i] != nil && bytes.Equal(item.Fields[i].Name, keyBytes) {
			return i
		}
	}
	return -1
}

func (item *DvVariable) InsertAtSimplePath(path string, child *DvVariable) bool {
	if item == nil || (item.Kind != FIELD_OBJECT && item.Kind != FIELD_ARRAY && item.Kind != FIELD_NULL && item.Kind != FIELD_UNDEFINED) {
		return false
	}
	if child == nil {
		child = &DvVariable{Kind: FIELD_NULL}
	}
	path = strings.Trim(path, ".")
	n := strings.Index(path, ".")
	current := path
	next := ""
	isLast := n < 0
	if !isLast {
		current = path[:n]
		next = path[n+1:]
	}
	if item.Fields == nil {
		item.Fields = make([]*DvVariable, 0, 1)
	}
	n = len(item.Fields)
	if item.Kind == FIELD_OBJECT {
		k := item.FindIndex(current)
		if k >= 0 {
			if isLast {
				item.Fields[k].CloneExceptKey(child, false)
			} else {
				return item.Fields[k].InsertAtSimplePath(next, child)
			}
		} else {
			itemAdd := &DvVariable{
				Kind: FIELD_NULL,
			}
			item.Fields = append(item.Fields, itemAdd)
			if isLast {
				itemAdd.CloneExceptKey(child, false)
				itemAdd.Name = []byte(current)
			} else {
				itemAdd.Name = []byte(current)
				return itemAdd.InsertAtSimplePath(next, child)
			}
		}
	} else {
		pos := dvtextutils.TryReadInteger(current, -1)
		if pos >= 0 && pos < n {
			if isLast {
				item.Fields[pos] = child
			} else {
				if item.Fields[pos] == nil {
					item.Fields[pos] = &DvVariable{Kind: FIELD_NULL}
				}
				return item.Fields[pos].InsertAtSimplePath(next, child)
			}
		} else {
			if item.Kind != FIELD_ARRAY {
				if pos == 0 {
					item.Kind = FIELD_ARRAY
				} else {
					item.Kind = FIELD_OBJECT
				}
			}
			if isLast {
				item.Fields = append(item.Fields, child)
				return true
			}
			nextItem := &DvVariable{Kind: FIELD_NULL}
			if item.Kind == FIELD_OBJECT {
				nextItem.Name = []byte(current)
			}
			item.Fields = append(item.Fields, nextItem)
			return nextItem.InsertAtSimplePath(next, child)
		}
	}
	return true
}

func (item *DvVariable) MergeAtChild(index int, child *DvVariable, mode int) bool {
	if item == nil || index < 0 || index >= len(item.Fields) {
		return false
	}
	if item.Fields[index] == nil {
		item.Fields[index] = child
		return true
	}
	item.Fields[index].CloneExceptKey(child, false)
	return true
}

func (item *DvVariable) MergeOtherVariable(other *DvVariable, mode int, ids []string) *DvVariable {
	if item == nil {
		return other
	}
	if other == nil {
		return item
	}
	switch item.Kind {
	case FIELD_NULL, FIELD_UNDEFINED:
		item.CloneExceptKey(other, false)
	case FIELD_ARRAY:
		if other.Kind == FIELD_ARRAY {
			if mode == UPDATE_MODE_REPLACE {
				item.Fields = other.Fields
			} else if mode == UPDATE_MODE_ADD_BY_KEYS || mode == UPDATE_MODE_MERGE {
				item.MergeArraysByIds(other, ids, mode)
			} else {
				item.Fields = append(item.Fields, other.Fields...)
			}
		} else {
			if mode == UPDATE_MODE_REPLACE {
				item.Fields = append(item.Fields[:0], other)
			} else {
				item.Fields = append(item.Fields, other)
			}
		}
	case FIELD_OBJECT:
		if other.Kind == FIELD_OBJECT {
			item.MergeObjectIntoObject(other, mode == UPDATE_MODE_REPLACE || mode == UPDATE_MODE_MERGE,
				mode == UPDATE_MODE_MERGE)
		}
	}
	return item
}

func (item *DvVariable) CloneExceptKey(other *DvVariable, deep bool) *DvVariable {
	if item == nil {
		item = &DvVariable{}
	}
	item.Kind = other.Kind
	item.Value = other.Value
	item.Extra = other.Extra
	item.Prototype = other.Prototype
	item.QuickSearch = nil
	if item.Kind == FIELD_ARRAY || item.Kind == FIELD_OBJECT || item.Kind == FIELD_FUNCTION {
		fld := other.Fields
		if deep {
			n := len(fld)
			newFields := make([]*DvVariable, n)
			item.Fields = newFields
			for i := 0; i < n; i++ {
				oldField := fld[i]
				if oldField == nil {
					continue
				}
				field := &DvVariable{Name: oldField.Name}
				field.CloneExceptKey(oldField, true)
				newFields[i] = field
			}
		} else {
			item.Fields = fld
		}

	} else {
		item.Fields = nil
	}
	return item
}

func (item *DvVariable) CloneWithKey(other *DvVariable, deep bool) *DvVariable {
	if item == nil {
		item = &DvVariable{}
	}
	item.CloneExceptKey(other, deep)
	item.Name = other.Name
	return item
}

func (item *DvVariable) MergeObjectIntoObject(other *DvVariable, replace bool, deep bool) {
	if item == nil || other == nil || len(other.Fields) == 0 {
		return
	}
	item.CreateQuickInfoForObjectType()
	fields := other.Fields
	n := len(fields)
	lookup := item.QuickSearch.Looker
	for i := 0; i < n; i++ {
		field := fields[i]
		name := string(field.Name)
		if m, ok := lookup[name]; ok {
			if replace {
				if deep && m.Kind == FIELD_OBJECT {
					m.MergeObjectIntoObject(field, true, true)
				} else {
					m.CloneExceptKey(field, false)
				}
			}
		} else {
			item.Fields = append(item.Fields, field)
			lookup[name] = field
		}
	}
}

func (item *DvVariable) Clone() *DvVariable {
	if item == nil {
		return nil
	}
	other := &DvVariable{Name: item.Name}
	other.CloneExceptKey(item, true)
	return other
}

func (item *DvVariable) MergeArraysByIds(other *DvVariable, ids []string, mode int) {
	if item == nil || other == nil {
		return
	}
	item.CreateQuickInfoByKeys(ids)
	other.CreateQuickInfoByKeys(ids)
	lookerMain := item.QuickSearch.Looker
	lookerSecond := other.QuickSearch.Looker
	for k, v := range lookerSecond {
		if m, ok := lookerMain[k]; ok {
			switch mode {
			case UPDATE_MODE_MERGE:
				m.MergeObjectIntoObject(v, true, true)
			case UPDATE_MODE_ADD_BY_KEYS:
				m.CloneExceptKey(v, true)
			}
		} else {
			item.Fields = append(item.Fields, v)
			lookerMain[k] = v
		}
	}
}

func (item *DvVariable) MergeItemIntoArraysByIds(other *DvVariable, ids []string, mode int, init bool) {
	if item == nil || other == nil {
		return
	}
	if init {
		item.CreateQuickInfoByKeys(ids)
	}
	key := other.CreateQuickInfoForSingleItemByKeys(ids)
	lookerMain := item.QuickSearch.Looker
	if m, ok := lookerMain[key]; ok {
		switch mode {
		case UPDATE_MODE_REPLACE:
			m.CloneExceptKey(other, true)
		case UPDATE_MODE_APPEND, UPDATE_MODE_MERGE_MAX, UPDATE_MODE_MERGE_MIN, UPDATE_MODE_ADD_BY_KEYS, UPDATE_MODE_MERGE:
			n := len(other.Fields)
			looker := m.QuickSearch.Looker
			for i := 0; i < n; i++ {
				f := other.Fields[i]
				if f == nil {
					continue
				}
				name := string(f.Name)
				p := looker[name]
				if p == nil {
					m.Fields = append(m.Fields, f)
				} else {
					v1 := string(p.Value)
					v2 := string(f.Value)
					switch mode {
					case UPDATE_MODE_ADD_BY_KEYS:
						p.CloneExceptKey(f, true)
					case UPDATE_MODE_APPEND:
						p.Value = []byte(v1 + "; " + v2)
					case UPDATE_MODE_MERGE_MAX:
						if v2 > v1 {
							p.Value = []byte(v2)
						}
					case UPDATE_MODE_MERGE_MIN:
						if v2 < v1 {
							p.Value = []byte(v2)
						}
					}
				}
			}
		}
	} else {
		item.Fields = append(item.Fields, other)
		lookerMain[key] = other
	}
}

func (item *DvVariable) ToDvGrammarExpressionValue() *dvgrammar.ExpressionValue {
	if item == nil || item.Kind == FIELD_NULL || item.Kind == FIELD_UNDEFINED {
		return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NULL}
	}
	switch item.Kind {
	case FIELD_BOOLEAN:
		return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_BOOLEAN, Value: len(item.Value) == 4 && item.Value[0] == 't'}
	case FIELD_STRING:
		return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_STRING, Value: string(item.Value)}
	case FIELD_NUMBER:
		s := string(item.Value)
		f, ok := AnyToNumberInt(s)
		if !ok || strings.Contains(s, ".") {
			return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NUMBER, Value: AnyToNumber(s)}
		} else {
			return &dvgrammar.ExpressionValue{DataType: dvgrammar.TYPE_NUMBER_INT, Value: f}
		}
	}
	return &dvgrammar.ExpressionValue{Value: item, DataType: dvgrammar.TYPE_OBJECT}
}

func (item *DvVariable) AssignToSubField(field string, value string, env *DvObject) error {
	if item == nil || item.Fields == nil {
		return nil
	}
	n := len(item.Fields)
	var err error
	if strings.HasPrefix(field, "$:") {
		field, err = env.EvaluateStringTypeExpression(field[2:])
		if err != nil {
			return err
		}
	}
	name := []byte(field)
	p := item.IndexOfByKey(name)
	if p < 0 {
		if value == "delete" {
			return nil
		}
		p = n
		item.Fields = append(item.Fields, &DvVariable{Kind: FIELD_NULL, Name: name})
	}
	if value == "delete" {
		if p == n-1 {
			item.Fields = item.Fields[:p]
		} else {
			item.Fields = append(item.Fields[:p], item.Fields[p+1:]...)
		}
	} else {
		var r *DvVariable
		if value == "" {
			r = &DvVariable{Kind: FIELD_NULL}
		} else {
			v, err := env.EvaluateAnyTypeExpression(value)
			if err != nil {
				return err
			}
			r = AnyToDvVariable(v)
			if r == nil {
				r = &DvVariable{Kind: FIELD_NULL}
			} else {
				other := &DvVariable{}
				other.CloneExceptKey(r, true)
				r = other
			}
		}
		r.Name = item.Fields[p].Name
		item.Fields[p] = r
	}
	return nil
}

func (item *DvVariable) IndexOfByKey(field []byte) int {
	if item == nil {
		return -1
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		if item.Fields[i] != nil && bytes.Equal(item.Fields[i].Name, field) {
			return i
		}
	}
	return -1
}
