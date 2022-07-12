/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
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
					Fn: Array_flat,
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
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("includes"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("indexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("isArray"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("join"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("keys"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("lastIndexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("map"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("of"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_flat,
				},
			},
			{
				Name: []byte("pop"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_push,
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
					Fn: Array_reduce,
				},
			},
			{
				Name: []byte("revert"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_reduce,
				},
			},
			{
				Name: []byte("shift"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_reduce,
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
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("sort"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("splice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("toLocaleString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("toString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("unshift"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
				},
			},
			{
				Name: []byte("values"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: Array_slice,
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
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
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
	beg := 0
	end := n
	if m >= 2 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return v, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 3 {
		end64, ok := dvevaluation.AnyToNumberInt(params[2])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return v, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		value := dvevaluation.AnyToDvVariable(val)
		for ; beg < end; beg++ {
			v.Fields[beg] = value
		}
	}
	return v, nil
}

func Array_concat(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
}

func Array_copyWithin(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
}

func Array_entries(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
}

func Array_every(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
}

func Array_filter(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	v := dvevaluation.AnyToDvVariable(thisVariable)
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return res, nil
	}
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
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
	m := len(params)
	n := len(v.Fields)
	beg := 0
	end := n
	if m >= 1 {
		beg64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok && beg64 > int64(-n) {
			if beg64 >= int64(n) {
				return res, nil
			}
			beg = int(beg64)
			if beg < 0 {
				beg += n
			}
		}
	}
	if m >= 2 {
		end64, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && end64 < int64(n) {
			if end64 <= int64(-n) {
				return res, nil
			}
			end = int(end64)
			if end < 0 {
				end += n
			}
		}
	}
	if beg < end {
		res.Fields = v.Fields[beg:end]
	}
	return res, nil
}
