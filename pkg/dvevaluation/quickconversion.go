// package dvevaluation covers expression calculations
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvevaluation

import "strconv"

var voidBuf = make([]byte, 128)
var bytesTrue = []byte(BooleanTrue)
var bytesFalse = []byte(BooleanFalse)

func ConvertIntToBuf(buf []byte, v int) []byte {
	n := v / 10
	k := v % 10
	sign := 0
	if v < 0 {
		sign = 1
		n = -n
		k = -k
	}
	c := 1
	p := n
	for p != 0 {
		p = p / 10
		c++
	}
	t := c + sign
	l := len(buf)
	buf = append(buf, voidBuf[:t]...)
	if sign != 0 {
		buf[l] = '-'
		l++
	}
	l += c - 1
	buf[l] = byte(k + 48)
	for n != 0 {
		l--
		buf[l] = byte(n%10 + 48)
		n /= 10
	}
	return buf
}

func ConvertInt64ToBuf(buf []byte, v int64) []byte {
	n := v / 10
	k := v % 10
	sign := 0
	if v < 0 {
		sign = 1
		n = -n
		k = -k
	}
	c := 1
	p := n
	for p != 0 {
		p = p / 10
		c++
	}
	t := c + sign
	l := len(buf)
	buf = append(buf, voidBuf[:t]...)
	if sign != 0 {
		buf[l] = '-'
		l++
	}
	l += c - 1
	buf[l] = byte(k + 48)
	for n != 0 {
		l--
		buf[l] = byte(n%10 + 48)
		n /= 10
	}
	return buf
}

func ConvertUnsignedIntToBuf(buf []byte, v uint) []byte {
	n := v / 10
	k := v % 10
	c := 1
	p := n
	for p != 0 {
		p = p / 10
		c++
	}
	l := len(buf)
	buf = append(buf, voidBuf[:c]...)
	l += c - 1
	buf[l] = byte(k + 48)
	for n != 0 {
		l--
		buf[l] = byte(n%10 + 48)
		n /= 10
	}
	return buf
}

func ConvertUnsignedInt64ToBuf(buf []byte, v uint64) []byte {
	n := v / 10
	k := v % 10
	c := 1
	p := n
	for p != 0 {
		p = p / 10
		c++
	}
	l := len(buf)
	buf = append(buf, voidBuf[:c]...)
	l += c - 1
	buf[l] = byte(k + 48)
	for n != 0 {
		l--
		buf[l] = byte(n%10 + 48)
		n /= 10
	}
	return buf
}

func ConvertSimpleTypeToBuf(buf []byte, v interface{}) ([]byte, bool) {
	switch v.(type) {
	case int:
		return ConvertIntToBuf(buf, v.(int)), true
	case int64:
		return ConvertInt64ToBuf(buf, v.(int64)), true
	case int8:
		return ConvertInt64ToBuf(buf, int64(v.(int8))), true
	case int16:
		return ConvertInt64ToBuf(buf, int64(v.(int16))), true
	case int32:
		return ConvertInt64ToBuf(buf, int64(v.(int32))), true
	case uint:
		return ConvertUnsignedIntToBuf(buf, v.(uint)), true
	case uint64:
		return ConvertUnsignedInt64ToBuf(buf, v.(uint64)), true
	case uint8:
		return ConvertUnsignedInt64ToBuf(buf, uint64(v.(uint8))), true
	case uint16:
		return ConvertUnsignedInt64ToBuf(buf, uint64(v.(uint16))), true
	case uint32:
		return ConvertUnsignedInt64ToBuf(buf, uint64(v.(uint32))), true
	case bool:
		if v.(bool) {
			buf = append(buf, bytesTrue...)
		} else {
			buf = append(buf, bytesFalse...)
		}
		return buf, true
	case float32:
		buf = append(buf, strconv.FormatFloat(float64(v.(float32)), 'g', -1, 32)...)
		return buf, true
	case float64:
		buf = append(buf, strconv.FormatFloat(v.(float64), 'g', -1, 64)...)
		return buf, true
	}
	return buf, false
}

func ConvertSimpleTypeToString(v interface{}) (string, bool) {
	buf, ok := ConvertSimpleTypeToBuf(make([]byte, 0, 100), v)
	return string(buf), ok
}
