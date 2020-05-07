/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

type JsonWriter struct {
	level      int
	f          *os.File
	err        error
	indent     int
	buf        []byte
	stack      []byte
	total      int
	cursor     int
	eol        []byte
	comma      bool
	justOpened bool
	CustomInfo interface{}
}

var jsonEscapeTable = map[byte]byte{
	'\\':     '\\',
	'"':      '"',
	byte(8):  'b',
	byte(9):  't',
	byte(10): 'r',
	byte(12): 'f',
	byte(13): 'n',
}

func IsJsonNumber(d []byte) bool {
	n := len(d)
	pnt := 0
	if n > 0 && (d[0] == '+' || d[0] == '-') {
		pnt++
	}
	if pnt >= n || n >= 24 {
		return false
	}
	b := d[pnt]
	if !(b >= '0' && b <= '9') {
		return false
	}
	for pnt < n && d[pnt] >= '0' && d[pnt] <= '9' {
		pnt++
	}
	if pnt == n {
		return true
	}
	if d[pnt] == '.' {
		for pnt++; pnt < n && d[pnt] >= '0' && d[pnt] <= '9'; pnt++ {
		}
		if pnt == n {
			return true
		}
	}
	b = d[pnt]
	pnt++
	if b != 'e' && b != 'E' || pnt == n {
		return false
	}
	b = d[pnt]
	pnt++
	if b == '+' || b == '-' {
		if pnt == n {
			return false
		}
		b = d[pnt]
		pnt++
	}
	if !(b >= '0' && b <= '9') || n-pnt > 2 {
		return false
	}
	for ; pnt < n; pnt++ {
		b = d[pnt]
		if !(b >= '0' && b <= '9') {
			return false
		}
	}
	return true
}

func IsJsonString(str []byte) bool {
	n := len(str)
	if n < 2 || str[0] != str[n-1] {
		return false
	}
	return true
}

func CreateJsonWriter(fileName string, indent int, bufSize int, eol int, customInfo interface{}) (*JsonWriter, error) {
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	var eolBuf []byte
	if eol == 0 {
		eolBuf = []byte{10}
	} else if eol < 0 {
		eolBuf = []byte{13}
	} else {
		eolBuf = []byte{13, 10}
	}
	return &JsonWriter{
		f:          f,
		indent:     indent,
		total:      bufSize,
		buf:        make([]byte, bufSize),
		stack:      make([]byte, 0, 20),
		eol:        eolBuf,
		CustomInfo: customInfo,
	}, nil
}

func (w *JsonWriter) Close() {
	if w.f != nil {
		w.Flush()
		err := w.f.Close()
		if err != nil && w.err == nil {
			w.err = err
		}
		w.f = nil
	}
}

func (w *JsonWriter) SetError(err error) {
	if w.err == nil {
		w.err = err
	}
}

func (w *JsonWriter) SetErrorMessage(message string) {
	if w.err == nil {
		w.err = errors.New(message)
	}
}

func (w *JsonWriter) Flush() {
	if w.cursor > 0 {
		_, err := w.f.Write(w.buf[:w.cursor])
		w.cursor = 0
		if err != nil {
			w.SetError(err)
		}
	}
}

func (w *JsonWriter) PutData(data ...byte) {
	n := len(data)
	if w.cursor+n > w.total {
		w.Flush()
		if n > w.total {
			w.f.Write(data)
			return
		}
	}
	b := w.buf[w.cursor:]
	for i := 0; i < n; i++ {
		b[i] = data[i]
	}
	w.cursor += n
}

func (w *JsonWriter) PutEscapedString(data []byte) {
	n := len(data)
	k := n*2 + 2
	if w.cursor+k > w.total {
		w.Flush()
	}
	b := w.buf[w.cursor:]
	m := w.total - w.cursor
	b[0] = '"'
	j := 1
	for i := 0; i < n; i++ {
		d := data[i]
		if c, ok := jsonEscapeTable[d]; ok {
			b[j] = '\\'
			j++
			if j == m {
				w.cursor = w.total
				w.Flush()
				b = w.buf
				m = w.total
				j = 0
			}
			d = c
		} else if d < ' ' {
			d = ' '
		}
		b[j] = d
		j++
		if j == m {
			w.cursor = w.total
			w.Flush()
			b = w.buf
			m = w.total
			j = 0
		}
	}
	b[j] = '"'
	w.cursor += j + 1
}

func (w *JsonWriter) RepeatData(data byte, n int) {
	m := w.total - w.cursor
	if m > n {
		m = n
	}
	b := w.buf[w.cursor:]
	for i := 0; i < m; i++ {
		b[i] = data
	}
	n -= m
	w.cursor += m
	if n != 0 {
		w.Flush()
		m = w.total
		if m > n {
			m = n
		}
		b = w.buf
		for i := 0; i < m; i++ {
			b[i] = data
		}
		if n == m {
			w.cursor = m
		} else {
			for n > 0 {
				w.f.Write(b)
				n -= m
				if n < m {
					m = n
					b = b[:n]
				}
			}
		}
	}
}

func (w *JsonWriter) StartArray() {
	if w.level == len(w.stack) {
		w.stack = append(w.stack, ']')
	} else {
		w.stack[w.level] = ']'
	}
	w.level++
	if w.comma {
		w.PutData(',')
	}
	w.PutData('[')
	w.comma = false
	w.justOpened = true
}

func (w *JsonWriter) EndArray() {
	if w.level == 0 {
		w.SetErrorMessage("not matching closing ]")
		return
	}
	w.level--
	if w.stack[w.level] != ']' {
		w.SetErrorMessage("expected } but found ]")
		return
	}
	if !w.justOpened {
		w.PutData(w.eol...)
		w.RepeatData(' ', w.level*w.indent)
	}
	w.PutData(']')
	w.comma = true
	w.justOpened = false
}

func (w *JsonWriter) StartObject() {
	if w.level == len(w.stack) {
		w.stack = append(w.stack, '}')
	} else {
		w.stack[w.level] = '}'
	}
	w.level++
	if w.comma {
		w.PutData(',')
	}
	w.PutData('{')
	w.comma = false
	w.justOpened = true
}

func (w *JsonWriter) EndObject() {
	if w.level == 0 {
		w.SetErrorMessage("not matching closing ]")
		return
	}
	w.level--
	if w.stack[w.level] != '}' {
		w.SetErrorMessage("expected ] but found }")
		return
	}
	if !w.justOpened {
		w.PutData(w.eol...)
		w.RepeatData(' ', w.level*w.indent)
	}
	w.PutData('}')
	w.comma = true
	w.justOpened = false
}

func (w *JsonWriter) PrintKey(key string) {
	if w.comma {
		w.PutData(',')
	}
	w.PutData(w.eol...)
	w.RepeatData(' ', w.level*w.indent)
	w.PutEscapedString([]byte(key))
	w.PutData(':', ' ')
	w.comma = false
	w.justOpened = false
}

func (w *JsonWriter) printPreValue() {
	if w.comma {
		w.PutData(',')
		w.PutData(w.eol...)
		w.RepeatData(' ', w.level*w.indent)
	} else if w.justOpened {
		w.PutData(w.eol...)
		w.RepeatData(' ', w.level*w.indent)
	}
	w.comma = true
	w.justOpened = false
}

func (w *JsonWriter) PrintValueSimple(value string) {
	w.printPreValue()
	w.PutData([]byte(value)...)
}

func (w *JsonWriter) PrintValueInteger(value int) {
	w.PrintValueSimple(strconv.Itoa(value))
}

func (w *JsonWriter) PrintValueString(value string) {
	w.printPreValue()
	w.PutEscapedString([]byte(value))
}

func (w *JsonWriter) PrintValueSmart(value interface{}) bool {
	f := ""
	switch value.(type) {
	case string:
		w.PrintValueString(value.(string))
		return true
	case int:
		f = strconv.Itoa(value.(int))
	case float64:
		f = fmt.Sprintf("%g", value.(float64))
	case nil:
		f = "null"
	case bool:
		if value.(bool) {
			f = "true"
		} else {
			f = "false"
		}
	case int64:
		f = strconv.FormatInt(value.(int64), 10)
	case uint:
		f = strconv.FormatInt(int64(value.(uint)), 10)
	case int8:
		f = strconv.Itoa(int(value.(int8)))
	case int16:
		f = strconv.Itoa(int(value.(int16)))
	case int32:
		f = strconv.Itoa(int(value.(int32)))
	case uint8:
		f = strconv.Itoa(int(value.(uint8)))
	case float32:
		f = fmt.Sprintf("%g", value.(float32))
	}
	if f != "" {
		w.PrintValueSimple(f)
		return true
	}
	return false
}

func IsValueSimple(v interface{}) bool {
	switch v.(type) {
	case string, int, int64, bool, float64, nil, uint, float32, int32, int16, int8, uint8:
		return true
	}
	return false
}

func (w *JsonWriter) GetError() error {
	return w.err
}

func ConvertToDouble(v interface{}) float64 {
	switch v.(type) {
	case float64:
		return v.(float64)
	case float32:
		return float64(v.(float32))
	case int:
		return float64(v.(int))
	case int64:
		return float64(v.(int64))
	case uint:
		return float64(v.(uint))
	case int32:
		return float64(v.(int32))
	case int16:
		return float64(v.(int16))
	case int8:
		return float64(v.(int8))
	case uint8:
		return float64(v.(uint8))
	}
	return math.NaN()
}

func ConvertToFloat(v interface{}) float32 {
	f, ok := v.(float32)
	if ok {
		return f
	}
	return float32(ConvertToDouble(v))
}

func ConvertToUnsignedLong(v interface{}) uint64 {
	switch v.(type) {
	case float64:
		return uint64(v.(float64))
	case float32:
		return uint64(v.(float32))
	case int:
		return uint64(v.(int))
	case int64:
		return uint64(v.(int64))
	case uint:
		return uint64(v.(uint))
	case int32:
		return uint64(v.(int32))
	case int16:
		return uint64(v.(int16))
	case int8:
		return uint64(v.(int8))
	case uint8:
		return uint64(v.(uint8))
	}
	return uint64(0)
}

func QuoteEscapedJsonBytes(b []byte) []byte {
	n := len(b)
	dstLen := n + 2
	res := make([]byte, dstLen)
	dst := 1
	res[0] = '"'
	for i := 0; i < n; i++ {
		d := b[i]
		if c, ok := jsonEscapeTable[d]; ok {
			if dst < dstLen {
				res[dst] = '\\'
				dst++
			} else {
				res = append(res, '\\')
			}
			d = c
		} else if d < ' ' {
			d = ' '
		}
		if dst < dstLen {
			res[dst] = d
			dst++
		} else {
			res = append(res, d)
		}
	}
	if dst < dstLen {
		res[dst] = '"'
		dst++
		if dst != dstLen {
			res = res[:dst]
		}
	} else {
		res = append(res, '"')
	}
	return res
}

func QuoteEscapedJsonBytesToString(b []byte) string {
	return string(QuoteEscapedJsonBytes(b))
}

func QuoteEscapedJsonString(s string) string {
	return QuoteEscapedJsonBytesToString([]byte(s))
}

func ConvertSimpleStringMapToJson(data map[string]string, includeBrackets bool) string {
	if data == nil {
		return ""
	}
	res := ""
	for k, v := range data {
		if res != "" {
			res += ","
		}
		res += QuoteEscapedJsonString(k) + ":" + QuoteEscapedJsonString(v)
	}
	if includeBrackets {
		res = "{" + res + "}"
	}
	return res
}

func GetKeyContentPecularities(key string) int {
	n := len(key)
	if n == 0 {
		return 1
	}
	level := 0
	for i := 0; i < n; i++ {
		c := key[i]
		if c == '"' || c == '\\' {
			return 2
		}
		if c == '.' || c == '/' || c == '[' || c == '{' || c == ']' || c == '}' {
			level = 1
		}
	}
	return level
}

func GetNextPathPartByKey(key string) string {
	switch GetKeyContentPecularities(key) {
	case 0:
		return "." + key
	case 1:
		return "[\"" + key + "\"]"
	}
	return "[" + string(QuoteEscapedJsonBytes([]byte(key))) + "]"
}
