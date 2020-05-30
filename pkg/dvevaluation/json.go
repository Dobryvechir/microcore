// package dvevaluation manages expressions, functions using agrammar
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvevaluation

import (
	"encoding/json"
	"fmt"
	"strconv"
)

var jsonEscapeTable = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	13:   'n',
	10:   'r',
	12:   'f',
	9:    't',
	8:    'b',
}

func placeQuotedStringToJson(buf []byte, s []byte) []byte {
	buf = append(buf, '"')
	n := len(s)
	for i := 0; i < n; i++ {
		b := s[i]
		c, ok := jsonEscapeTable[b]
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
	switch v.(type) {
	case int:
		buf = append(buf, strconv.Itoa(v.(int))...)
	case int64:
		buf = append(buf, strconv.FormatInt(v.(int64), 10)...)
	case float32:
		buf = append(buf, fmt.Sprintf("%f", float64(v.(float32)))...)
	case float64:
		buf = append(buf, fmt.Sprintf("%f", v.(float64))...)
	case nil:
		buf = append(buf, []byte("null")...)
	case bool:
		{
			if v.(bool) {
				buf = append(buf, []byte("true")...)
			} else {
				buf = append(buf, []byte("false")...)
			}
		}
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
