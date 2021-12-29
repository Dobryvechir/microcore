// package dvevaluation manages expressions, functions using a grammar
// MicroCore Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvevaluation

import (
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func placeQuotedStringToJson(buf []byte, s []byte) []byte {
	buf = append(buf, '"')
	n := len(s)
	for i := 0; i < n; i++ {
		b := s[i]
		c, ok := dvtextutils.JsonEscapeTable[b]
		if ok {
			buf = append(buf, '\\', c)
		} else {
			buf = append(buf, b)
		}
	}
	buf = append(buf, '"')
	return buf
}

func ConvertAnyTypeToJson(buf []byte, v interface{}) []byte {
	n := len(buf)
	if n == 0 {
		buf = make([]byte, 0, 4096)
	} else if cap(buf)-n < 1024 {
		buf = append(buf, buf...)
		buf = buf[:n]
	}
	var ok bool
	buf, ok, _  = ConvertSimpleTypeToBuf(buf, v)
	if ok {
		return buf
	}
	switch v.(type) {
	case nil:
		buf = append(buf, []byte("null")...)
	case string:
		{
			p := v.(string)
			buf = placeQuotedStringToJson(buf, []byte(p))
		}
	case []interface{}:
		{
			p := v.([]interface{})
			if p == nil {
				buf = append(buf, []byte("null")...)
			} else {
				buf = append(buf, '[')
				isNext := false
				for _, val := range p {
					if isNext {
						buf = append(buf, ',')
					} else {
						isNext = true
					}
					buf = ConvertAnyTypeToJson(buf, val)
				}
				buf = append(buf, ']')
			}
		}
	case map[string]interface{}:
		{
			p := v.(map[string]interface{})
			if p == nil {
				buf = append(buf, []byte("null")...)
			} else {
				buf = append(buf, '{')
				isNext := false
				for key, val := range p {
					if isNext {
						buf = append(buf, ',')
					} else {
						isNext = true
					}
					buf = placeQuotedStringToJson(buf, []byte(key))
					buf = append(buf, ':')
					buf = ConvertAnyTypeToJson(buf, val)
				}
				buf = append(buf, '}')
			}
		}
	default:
		{
			p, err := json.Marshal(v)
			if err != nil {
				buf = append(buf, []byte("{\"error\":")...)
				buf = placeQuotedStringToJson(buf, []byte(err.Error()))
				buf = append(buf, '}')
			} else {
				buf = append(buf, p...)
			}
		}
	}
	return buf
}

func ConvertAnyTypeToJsonString(data interface{}) string {
	return string(ConvertAnyTypeToJson(nil, data))
}
