/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvcrypt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strings"
)

func string_init() {
	dvevaluation.StringMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("charAt"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_charAt,
				},
			},
			{
				Name: []byte("charCodeAt"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_charCodeAt,
				},
			},
			{
				Name: []byte("codePointAt"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_codePointAt,
				},
			},
			{
				Name: []byte("concat"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_concat,
				},
			},
			{
				Name: []byte("endsWith"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_endsWith,
				},
			},
			{
				Name: []byte("fromCharCode"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_fromCharCode,
				},
			},
			{
				Name: []byte("fromCodePoint"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_fromCodePoint,
				},
			},
			{
				Name: []byte("includes"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("indexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_indexOf,
				},
			},
			{
				Name: []byte("lastIndexOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_lastIndexOf,
				},
			},
			{
				Name: []byte("localeCompare"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_localeCompare,
				},
			},
			{
				Name: []byte("match"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_match,
				},
			},
			{
				Name: []byte("matchAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("normalize"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("padEnd"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("padStart"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("raw"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("repeat"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("replace"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("replaceAll"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("search"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("slice"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_includes,
				},
			},
			{
				Name: []byte("split"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_split,
				},
			},
			{
				Name: []byte("startsWith"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("substring"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toLocaleLowerCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toLocaleUpperCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toLowerCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("toString"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_toString,
				},
			},
			{
				Name: []byte("toUpperCase"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("trim"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("trimEnd"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("trimStart"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("valueOf"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_startsWith,
				},
			},
			{
				Name: []byte("length"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn:        String_length,
					Immediate: true,
				},
			},
			{
				Name: []byte("isValidUUID"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: String_isValidUUID,
				},
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func String_includes(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	b := strings.Contains(s, s1)
	return b, nil
}

func String_charAt(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	m := len(s)
	if m == 0 {
		return "", nil
	}
	n := len(params)
	p := 0
	if n != 0 && params[0] != nil {
		p64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok {
			if p64 >= int64(m) || p64 < 0 {
				return "", nil
			}
			p = int(p64)
		}
	}
	return s[p : p+1], nil
}

func String_length(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	n := len(s)
	return n, nil
}

func String_endsWith(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	if n >= 2 && params[1] != nil {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && p >= 0 && p < int64(len(s)) {
			m := int(p)
			s = s[:m]
		}
	}
	b := strings.HasSuffix(s, s1)
	return b, nil
}

func String_startsWith(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	if n >= 2 && params[1] != nil {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && p >= 0 && p < int64(len(s)) {
			m := int(p)
			s = s[:m]
		}
	}
	b := strings.HasPrefix(s, s1)
	return b, nil
}

func String_split(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 || params[0] == nil {
		return true, nil
	}
	s1 := dvevaluation.AnyToString(params[0])
	limit := 0
	if n >= 2 && params[1] != nil {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok && p >= 0 && p < int64(len(s)) {
			limit = int(p)
		}
	}
	var b []string
	if limit > 0 {
		b = strings.SplitN(s, s1, limit)
	} else {
		b = strings.Split(s, s1)
	}
	return b, nil
}

func String_isValidUUID(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	b := dvcrypt.IsValidUUID(s)
	return b, nil
}

func String_charCodeAt(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	m := len(s)
	if m == 0 {
		return nil, nil
	}
	n := len(params)
	p := 0
	if n != 0 && params[0] != nil {
		p64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok {
			if p64 >= int64(m) || p64 < 0 {
				return nil, nil
			}
			p = int(p64)
		}
	}
	return int(s[p]), nil
}

func String_codePointAt(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	if thisVariable == nil {
		return false, nil
	}
	s := dvevaluation.AnyToString(thisVariable)
	b := dvtextutils.SeparateBytesToUTF8Chars([]byte(s))
	m := len(b)
	if m == 0 {
		return nil, nil
	}
	n := len(params)
	p := 0
	if n != 0 && params[0] != nil {
		p64, ok := dvevaluation.AnyToNumberInt(params[0])
		if ok {
			if p64 >= int64(m) || p64 < 0 {
				return nil, nil
			}
			p = int(p64)
		}
	}
	res := dvtextutils.GetCodePoint(b[p])
	return res, nil
}

func String_concat(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	for i := 0; i < n; i++ {
		s += dvevaluation.AnyToString(params[i])
	}
	return s, nil
}

func String_fromCharCode(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	buf := make([]byte, 0, 2*n)
	for i := 0; i < n; i++ {
		b, ok := dvevaluation.AnyToNumberInt(params[i])
		if !ok {
			continue
		}
		b = b & 0xffff
		if b < 256 {
			buf = append(buf, byte(b))
		} else {
			buf = append(buf, dvtextutils.GetBytesFromPointCode(int(b))...)
		}
	}
	return string(buf), nil
}

func String_fromCodePoint(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	buf := make([]byte, 0, n<<2)
	for i := 0; i < n; i++ {
		b, ok := dvevaluation.AnyToNumberInt(params[i])
		if !ok {
			continue
		}
		b = b & 0xffffffff
		sub := dvtextutils.GetBytesFromPointCode(int(b))
		buf = append(buf, sub...)
	}
	return string(buf), nil
}

func String_indexOf(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 {
		return -1, nil
	}
	t := dvevaluation.AnyToString(params[0])
	pos := 0
	m := len(s)
	if n > 1 {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok {
			if p >= int64(m) {
				return -1, nil
			}
			if p >= 0 {
				pos = int(p)
			}
		}
	}
	res := strings.Index(s[pos:], t)
	if res >= 0 {
		res += pos
	}
	return res, nil
}

func String_lastIndexOf(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	n := len(params)
	if n == 0 {
		return -1, nil
	}
	t := dvevaluation.AnyToString(params[0])
	m := len(s)
	pos := m
	if n > 1 {
		p, ok := dvevaluation.AnyToNumberInt(params[1])
		if ok {
			if p <= 0 {
				return -1, nil
			}
			if p < int64(m) {
				pos = int(p)
			}
		}
	}
	res := strings.LastIndex(s[:pos], t)
	return res, nil
}

func String_localeCompare(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	t := ""
	if len(params) > 0 {
		t = dvevaluation.AnyToString(params[0])
	}
	b := strings.Compare(s, t)
	return b, nil
}

func String_match(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	s := dvevaluation.AnyToString(thisVariable)
	t := ""
	if len(params) > 0 {
		t = dvevaluation.AnyToString(params[0])
	}
	b := strings.Compare(s, t)
	return b, nil
}
