/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strconv"
	"strings"
)

const MapKeySeparator = "~^~!~"
const fnResolvedSign = "_$_"
const MaxInt = int64(^uint(0) >> 1)
const (
	EXPRESSION_RESOLVER_CACHE = 1 << iota
)

func (item *DvVariable) CreateQuickInfoByKeys(ids []string) {
	if item != nil {
		if item.QuickSearch == nil {
			item.QuickSearch = &QuickSearchInfo{Looker: make(map[string]*DvVariable)}
		}
		if item.QuickSearch.Looker == nil {
			item.QuickSearch.Looker = make(map[string]*DvVariable)
		}
		n := len(item.Fields)
		m := len(ids)
		if m == 0 {
			ids = []string{"id"}
			m = 1
		}
		for i := 0; i < n; i++ {
			f := item.Fields[i]
			if f != nil {
				f.CreateQuickInfoForObjectType()
				key := ""
				for j := 0; j < m; j++ {
					id := ids[j]
					fc := f.QuickSearch.Looker[id]
					v := ""
					if fc != nil {
						v = string(fc.Value)
					}
					if j == 0 {
						key = v
					} else {
						key = key + MapKeySeparator + v
					}
				}
				f.QuickSearch.Key = key
				item.QuickSearch.Looker[key] = f
			}
		}
	}
}

func (item *DvVariable) CreateQuickInfoForObjectType() {
	if item != nil {
		if item.QuickSearch == nil {
			item.QuickSearch = &QuickSearchInfo{Looker: make(map[string]*DvVariable)}
		}
		if item.QuickSearch.Looker == nil {
			item.QuickSearch.Looker = make(map[string]*DvVariable)
		}
		n := len(item.Fields)
		for i := 0; i < n; i++ {
			f := item.Fields[i]
			if f != nil {
				name := string(f.Name)
				if f.QuickSearch == nil {
					f.QuickSearch = &QuickSearchInfo{Key: name}
				} else {
					f.QuickSearch.Key = name
				}
				item.QuickSearch.Looker[name] = f
			}
		}
	}
}

func MeetItemExpression(expr string, item *DvVariable, resolver ExpressionResolver, options int, index int) (bool, error) {
	props := ConvertDvFieldInfoToProperties(item, index)
	res, err := resolver(expr, props, options)
	if err != nil {
		return false, err
	}
	return res == "true", nil
}

func FindItemByExpression(expr string, resolver ExpressionResolver, item *DvVariable, strict bool) (*DvVariable, error) {
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
		if dvtextutils.IsJsonNumber([]byte(expr)) {
			return expr, nil
		}
		if dvtextutils.IsJsonString([]byte(expr)) {
			return expr[1 : n-1], nil
		}
	}
	if resolver == nil {
		return "", fmt.Errorf("Expression %s is complex and there is no expression resolver", expr)
	}
	return resolver(expr, nil, 0)
}

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
	case FIELD_UNDEFINED:
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

func ConvertDvFieldInfoToProperties(item *DvVariable, index int) map[string]interface{} {
	res := map[string]interface{}{"_index": index}
	if item == nil {
		return res
	}
	if len(item.Name) != 0 {
		res["_name"] = string(item.Name)
	}
	if len(item.Value) != 0 {
		res["_value"] = string(item.Value)
	}
	n := len(item.Fields)
	for i := 0; i < n; i++ {
		current := item.Fields[i]
		if current == nil {
			res["_"+strconv.Itoa(i)] = ""
		} else {
			var name string
			if len(current.Name) == 0 {
				name = "_" + strconv.Itoa(i)
			} else {
				name = string(current.Name)
			}
			var value interface{}
			if current.Kind == FIELD_ARRAY || current.Kind == FIELD_OBJECT {
				value = current
			} else {
				value = string(current.Value)
			}
			res[name] = value
		}
	}
	return res
}
