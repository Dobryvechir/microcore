/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"sort"
)

func array_init() {
	dvevaluation.ArrayMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("concat"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_concat,
				},
			},
			{
				Name: []byte("copyWithin"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_copyWithin,
				},
			},
			{
				Name: []byte("entries"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_entries,
				},
			},
			{
				Name: []byte("every"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_every,
				},
			},
			{
				Name: []byte("fill"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_fill,
				},
			},
			{
				Name: []byte("filter"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_filter,
				},
			},
			{
				Name: []byte("find"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_find,
				},
			},
			{
				Name: []byte("findIndex"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_findIndex,
				},
			},
			{
				Name: []byte("findLast"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_findLast,
				},
			},
			{
				Name: []byte("findLastIndex"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_findLastIndex,
				},
			},
			{
				Name: []byte("flat"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("flatMap"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flatMap,
				},
			},
			{
				Name: []byte("forEach"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_foreach,
				},
			},
			{
				Name: []byte("from"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_from,
				},
			},
			{
				Name: []byte("includes"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_includes,
				},
			},
			{
				Name: []byte("indexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_indexOf,
				},
			},
			{
				Name: []byte("isArray"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_isArray,
				},
			},
			{
				Name: []byte("join"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_join,
				},
			},
			{
				Name: []byte("keys"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_keys,
				},
			},
			{
				Name: []byte("lastIndexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_lastIndexOf,
				},
			},
			{
				Name: []byte("map"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_map,
				},
			},
			{
				Name: []byte("of"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_of,
				},
			},
			{
				Name: []byte("pop"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_pop,
				},
			},
			{
				Name: []byte("push"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_push,
				},
			},
			{
				Name: []byte("reduce"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_reduce,
				},
			},
			{
				Name: []byte("reduceRight"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_reduceRight,
				},
			},
			{
				Name: []byte("revert"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_revert,
				},
			},
			{
				Name: []byte("shift"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_shift,
				},
			},
			{
				Name: []byte("slice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("some"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_some,
				},
			},
			{
				Name: []byte("sort"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_sort,
				},
			},
			{
				Name: []byte("splice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_splice,
				},
			},
			{
				Name: []byte("toLocaleString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_toLocaleString,
				},
			},
			{
				Name: []byte("toString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_toString,
				},
			},
			{
				Name: []byte("unshift"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_unshift,
				},
			},
			{
				Name: []byte("values"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_values,
				},
			},
			{
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        Array_length,
					Immediate: true,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func Array_push(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil {
		return nil, errors.New("Cannot convert null to object")
	}
	n := len(params)
	for i := 0; i < n; i++ {
		d := dvevaluation.AnyToDvVariable(params[i])
		v.Fields = append(v.Fields, d)
	}
	n = len(v.Fields)
	return n, nil
}

func Array_slice(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	n := len(v.Fields)
	beg := readBeginIndex(params, 0, n)
	end := readEndIndex(params, 1, n)
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
}

func Array_length(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	n := 0
	if v != nil {
		n = len(v.Fields)
	}
	return n, nil
}

func Array_reduce(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var result interface{} = nil
	n := len(params)
	if n >= 2 {
		result = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		m := len(v.Fields)
		for i := 0; i < m; i++ {
			fnParams := []interface{}{result, v.Fields[i], i, v}
			result, err = dvevaluation.ExecuteAnyFunction(context, fn, v, fnParams)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func Array_reduceRight(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var result interface{} = nil
	n := len(params)
	if n >= 2 {
		result = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		m := len(v.Fields)
		for i := m - 1; i >= 0; i-- {
			fnParams := []interface{}{result, v.Fields[i], i, v}
			result, err = dvevaluation.ExecuteAnyFunction(context, fn, v, fnParams)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func Array_foreach(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		for i := 0; i < len(v.Fields); i++ {
			fnParams := []interface{}{v.Fields[i], i, v}
			_, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func Array_fill(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return v, nil
	}
	m := len(params)
	n := len(v.Fields)
	var val interface{} = nil
	if m >= 1 {
		val = params[0]
	}
	beg := readBeginIndex(params, 1, n)
	end := readEndIndex(params, 2, n)
	if beg < end {
		value := dvevaluation.AnyToDvVariable(val)
		for ; beg < end; beg++ {
			v.Fields[beg] = value
		}
	}
	return v, nil
}

func pushArray(dst *dvevaluation.DvVariable, src interface{}) {
	v := dvevaluation.AnyToDvVariable(src)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return
	}
	dst.Fields = append(dst.Fields, v.Fields...)
}

func Array_concat(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	pushArray(res, thisVariable)
	n := len(params)
	for i := 0; i < n; i++ {
		pushArray(res, params[i])
	}
	return res, nil
}

func Array_copyWithin(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return thisVariable, nil
	}
	n := len(v.Fields)
	dst := readBeginIndex(params, 0, n)
	beg := readBeginIndex(params, 1, n)
	end := readEndIndex(params, 2, n)
	if dst < n && beg < n && dst != beg {
		for ; beg < end; beg++ {
			if dst >= 0 && dst < n {
				v.Fields[dst] = v.Fields[beg]
			}
			dst++
		}
	}
	return v, nil
}

func Array_every(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		for i := 0; i < len(v.Fields); i++ {
			fnParams := []interface{}{v.Fields[i], i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return nil, err
			}
			if !dvevaluation.AnyToBoolean(res) {
				return false, nil
			}
		}
	}
	return true, nil
}

func Array_map(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if v != nil && v.Kind == dvevaluation.FIELD_STRING {
		s:=dvtextutils.SeparateBytesToUTF8Chars(v.Value)
		m := len(s)
		u := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, m)}
		for i := 0; i < m; i++ {
			u.Fields[i] = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Value: s[i]}
		}
		v = u
	}
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v != nil && len(v.Fields) > 0 {
		nc := len(v.Fields)
		res.Fields = make([]*dvevaluation.DvVariable, nc)
		if n >= 1 {
			fn := params[0]
			var err error
			var vl interface{}
			for i := 0; i < nc; i++ {
				fnParams := []interface{}{v.Fields[i], i, v}
				vl, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
				if err != nil {
					return nil, err
				}
				res.Fields[i] = dvevaluation.AnyToDvVariable(vl)
			}
		} else {
			for i := 0; i < nc; i++ {
				res.Fields[i] = v.Fields[i]
			}
		}

	}
	return res, nil
}

func Array_some(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		for i := 0; i < len(v.Fields); i++ {
			fnParams := []interface{}{v.Fields[i], i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return nil, err
			}
			if dvevaluation.AnyToBoolean(res) {
				return true, nil
			}
		}
	}
	return false, nil
}

func Array_filter(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	result := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		result.Fields = make([]*dvevaluation.DvVariable, 0, len(v.Fields))
		for i := 0; i < len(v.Fields); i++ {
			el := v.Fields[i]
			fnParams := []interface{}{el, i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return nil, err
			}
			if dvevaluation.AnyToBoolean(res) {
				result.Fields = append(result.Fields, el)
			}
		}
	}
	return result, nil
}

func Array_find(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		for i := 0; i < len(v.Fields); i++ {
			el := v.Fields[i]
			fnParams := []interface{}{el, i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return nil, err
			}
			if dvevaluation.AnyToBoolean(res) {
				return el, nil
			}
		}
	}
	return nil, nil
}

func Array_findIndex(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		for i := 0; i < len(v.Fields); i++ {
			el := v.Fields[i]
			fnParams := []interface{}{el, i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return -1, err
			}
			if dvevaluation.AnyToBoolean(res) {
				return i, nil
			}
		}
	}
	return -1, nil
}

func Array_findLast(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		for i := len(v.Fields) - 1; i >= 0; i-- {
			el := v.Fields[i]
			fnParams := []interface{}{el, i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return nil, err
			}
			if dvevaluation.AnyToBoolean(res) {
				return el, nil
			}
		}
	}
	return nil, nil
}

func Array_findLastIndex(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	var thisArg interface{} = nil
	n := len(params)
	if n >= 2 {
		thisArg = params[1]
	}
	if n >= 1 && v != nil && len(v.Fields) > 0 {
		fn := params[0]
		var err error
		var res interface{}
		for i := len(v.Fields) - 1; i >= 0; i-- {
			el := v.Fields[i]
			fnParams := []interface{}{el, i, v}
			res, err = dvevaluation.ExecuteAnyFunction(context, fn, thisArg, fnParams)
			if err != nil {
				return -1, err
			}
			if dvevaluation.AnyToBoolean(res) {
				return el, nil
			}
		}
	}
	return -1, nil
}

func Array_flat(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	depth := 1
	if len(params) > 0 {
		n64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && n64 >= 0 {
			depth = int(n64)
		}
	}
	flattenArray(res, v, depth)
	return res, nil
}

func Array_flatMap(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	resMap, err := Array_map(context, thisVariable, params)
	if err != nil {
		return nil, err
	}
	v := dvevaluation.AnyToDvVariable(resMap)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	depth := 1
	if len(params) > 2 {
		n64, ok := dvevaluation.AnyToNumberInt(params[2])
		if ok && n64 >= 0 {
			depth = int(n64)
		}
	}
	flattenArray(res, v, depth)
	return res, nil
}

func flattenArray(dst *dvevaluation.DvVariable, src *dvevaluation.DvVariable, depth int) {
	if src != nil && src.Kind == dvevaluation.FIELD_ARRAY {
		n := len(src.Fields)
		if n > 0 {
			if dst.Fields == nil {
				dst.Fields = make([]*dvevaluation.DvVariable, 0, n*2)
			}
			for i := 0; i < n; i++ {
				v := src.Fields[i]
				if v != nil && v.Kind != dvevaluation.FIELD_NULL {
					if v.Kind == dvevaluation.FIELD_ARRAY && depth > 0 {
						flattenArray(dst, v, depth-1)
					} else {
						dst.Fields = append(dst.Fields, v)
					}
				}
			}
		}
	}
}

func Array_isArray(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if len(params) == 0 {
		return false, nil
	}
	v := dvevaluation.AnyToDvVariable(params[0])
	res := v != nil && v.Kind == dvevaluation.FIELD_ARRAY
	return res, nil
}

func Array_indexOf(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || len(v.Fields) == 0 {
		return -1, nil
	}
	n := len(v.Fields)
	m := len(params)
	var el interface{} = nil
	if m >= 1 {
		el = params[0]
	}
	beg := readBeginIndex(params, 1, n)
	if beg >= n {
		return -1, nil
	}
	elSearch := dvevaluation.AnyToDvVariable(el)
	for ; beg < n; beg++ {
		elOther := v.Fields[beg]
		if elOther.CompareWholeDvField(elSearch) == 0 {
			return beg, nil
		}
	}
	return -1, nil
}

func Array_includes(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || len(v.Fields) == 0 {
		return false, nil
	}
	n := len(v.Fields)
	m := len(params)
	var el interface{} = nil
	if m >= 1 {
		el = params[0]
	}
	beg := readBeginIndex(params, 1, n)
	if beg >= n {
		return -1, nil
	}
	elSearch := dvevaluation.AnyToDvVariable(el)
	for ; beg < n; beg++ {
		elOther := v.Fields[beg]
		if elOther.CompareWholeDvField(elSearch) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func Array_lastIndexOf(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || len(v.Fields) == 0 {
		return -1, nil
	}
	n := len(v.Fields)
	m := len(params)
	var el interface{} = nil
	if m >= 1 {
		el = params[0]
	}
	end := readEndIndex(params, 1, n)
	if end > n-1 {
		end = n - 1
	}
	elSearch := dvevaluation.AnyToDvVariable(el)
	for ; end >= 0; end-- {
		elOther := v.Fields[end]
		if elOther.CompareWholeDvField(elSearch) == 0 {
			return end, nil
		}
	}
	return -1, nil
}

func readBeginIndex(params []interface{}, index int, n int) int {
	m := len(params)
	if index >= m {
		return 0
	}
	beg64, ok := dvevaluation.AnyToNumberInt(params[index])
	beg := 0
	if ok && beg64 > int64(-n) {
		if beg64 >= int64(n) {
			return n
		}
		beg = int(beg64)
		if beg < 0 {
			beg += n
		}
	}
	return beg
}

func readEndIndex(params []interface{}, index int, n int) int {
	m := len(params)
	if index >= m {
		return n
	}
	end := n
	end64, ok := dvevaluation.AnyToNumberInt(params[index])
	if ok && end64 < int64(n) {
		if end64 <= int64(-n) {
			return -1
		}
		end = int(end64)
		if end < 0 {
			end += n
		}
	}
	return end
}

func Array_unshift(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY {
		return 0, nil
	}
	n := len(v.Fields)
	m := len(params)
	newN := n + m
	fld := v.Fields
	v.Fields = make([]*dvevaluation.DvVariable, newN)
	for i := 0; i < n; i++ {
		v.Fields[i+m] = fld[i]
	}
	for i := 0; i < m; i++ {
		v.Fields[i] = dvevaluation.AnyToDvVariable(params[i])
	}
	return newN, nil
}

func Array_shift(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return nil, nil
	}
	res := v.Fields[0]
	v.Fields = v.Fields[1:]
	return res, nil
}

func Array_pop(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return nil, nil
	}
	n := len(v.Fields) - 1
	res := v.Fields[n]
	v.Fields = v.Fields[:n]
	return res, nil
}

func Array_join(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return "", nil
	}
	joiner := ","
	if len(params) > 0 {
		joiner = dvevaluation.AnyToString(params[0])
	}
	res := ArrayJoinWith(v, joiner)
	return res, nil
}

func Array_splice(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY {
		return res, nil
	}
	m := len(params)
	n := len(v.Fields)
	start := readBeginIndex(params, 0, n)
	deleteCount := readEndIndex(params, 1, n)
	if start > n || start < 0 {
		return res, nil
	}
	if deleteCount > n-start {
		deleteCount = n - start
	}
	if deleteCount < 0 {
		deleteCount = 0
	}
	m -= 2
	if m < 0 {
		m = 0
	}
	if deleteCount > 0 {
		res.Fields = make([]*dvevaluation.DvVariable, deleteCount)
		for i := 0; i < deleteCount; i++ {
			res.Fields[i] = v.Fields[start+i]
		}
	}
	if m > deleteCount {
		fld := v.Fields
		newN := n + m - deleteCount
		v.Fields = make([]*dvevaluation.DvVariable, newN)
		copy(v.Fields, fld[:start])
		copy(v.Fields[start+m:], fld[start+deleteCount:])
	} else if m < deleteCount {
		if start+deleteCount < n {
			v.Fields = append(v.Fields[:start+m], v.Fields[start+deleteCount:]...)
		} else {
			v.Fields = v.Fields[start+m:]
		}
	}
	for i := 0; i < m; i++ {
		v.Fields[start+i] = dvevaluation.AnyToDvVariable(params[i+2])
	}
	return res, nil
}

func Array_revert(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY {
		return thisVariable, nil
	}
	n := len(v.Fields)
	m := n >> 1
	mn := n - 1
	for i := 0; i <= m; i++ {
		v.Fields[i], v.Fields[mn-i] = v.Fields[mn-i], v.Fields[i]
	}
	return v, nil
}

func Array_sort(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	n := len(params)
	var fn interface{} = nil
	if n >= 1 {
		fn = params[0]
	}
	var commonError error = nil
	if v != nil && len(v.Fields) > 0 {
		if fn != nil {
			sort.SliceStable(v.Fields, func(i int, j int) bool {
				fnParams := []interface{}{v.Fields[i], v.Fields[j]}
				res, err := dvevaluation.ExecuteAnyFunction(context, fn, v, fnParams)
				if err != nil {
					commonError = err
					return false
				}
				nres, ok := dvevaluation.AnyToNumberInt(res)
				if ok && nres < 0 {
					return true
				}
				return false
			})
		} else {
			sort.SliceStable(v.Fields, func(i int, j int) bool {
				a := v.Fields[i]
				b := v.Fields[j]
				if a == nil || a.Kind == dvevaluation.FIELD_NULL {
					return false
				}
				if b == nil || a.Kind == dvevaluation.FIELD_NULL {
					return true
				}
				res := a.CompareWholeDvField(b)
				return res < 0
			})
		}
		return v, commonError
	}
	return thisVariable, nil
}

func Array_from(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	m:=len(params)
	if m==0 {
		return nil, errors.New("Array.from requires parameters")
	}
	self:=params[0]
	params = params[1:]
	v, err:=Array_map(context, self, params)
	return v, err
}

func Array_of(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	m:=len(params)
	v:=&dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, m)}
	for i:=0;i<m;i++ {
		v.Fields[i] = dvevaluation.AnyToDvVariable(params[i])
	}
	return v, nil
}
