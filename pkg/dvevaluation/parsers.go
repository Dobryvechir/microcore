/*
**********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
***********************************************************************
*/
package dvevaluation

import (
	"math"
	"strconv"

	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

func StringToNumber(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func GetNumberKindFromString(v []byte, n int) (int, int, int64, float64) {
	var i int
	for i = 0; i < n && v[i] <= ' '; i++ {
	}
	if i == n {
		return -1, n, 0, 0
	}
	c := v[i]
	neg := false
	if c == '+' || c == '-' {
		i++
		if i >= n {
			return -1, n, 0, 0
		}
		if c == '-' {
			neg = true
		}
		c = v[i]
	}
	if !(c >= '0' && c <= '9') {
		return -1, i, 0, 0
	}
	var resInt int64 = int64(c) - 48
	for i++; i < n && v[i] >= '0' && v[i] <= '9'; i++ {
		resInt = resInt*10 + (int64(v[i]) - 48)
	}
	if i == n {
		c = ' '
	} else {
		c = v[i]
	}
	if c != '.' && c != 'e' && c != 'E' {
		if neg {
			resInt = -resInt
		}
		return dvgrammar.TYPE_NUMBER_INT, i, resInt, 0
	}
	var resFloat float64 = float64(resInt)
	if c == '.' {
		var fraction float64 = 1
		for i++; i < n && v[i] >= '0' && v[i] <= '9'; i++ {
			fraction = fraction * 0.1
			resFloat += fraction * float64(int(v[i])-48)
		}
		if neg {
			resFloat = -resFloat
		}
		if i == n || (v[i] != 'e' && v[i] != 'E') {
			return dvgrammar.TYPE_NUMBER, i, resInt, resFloat
		}
	} else {
		if neg {
			resFloat = -resFloat
		}
	}
	i++
	if i == n {
		return dvgrammar.TYPE_NUMBER, i, resInt, resFloat
	}
	c = v[i]
	neg = c == '-'
	if c == '+' || neg {
		i++
		if i == n {
			return dvgrammar.TYPE_NUMBER, i, resInt, resFloat
		}
		c = v[i]
	}
	power := int(c) - 48
	if !(power >= 0 && power <= 9) {
		return dvgrammar.TYPE_NUMBER, i, resInt, resFloat
	}
	for i++; i < n && v[i] >= '0' && v[i] <= '9'; i++ {
		power = power*10 + (int(v[i]) - 48)
	}
	if neg {
		power = -power
	}
	if power != 0 {
		resFloat = resFloat * math.Pow10(power)
	}
	return dvgrammar.TYPE_NUMBER, i, resInt, resFloat
}

func collectKeysToMap(keys []string) map[string]int {
	res := make(map[string]int)
	if keys == nil {
		return res
	}
	for _, v := range keys {
		res[v] = 1
	}
	return res
}
