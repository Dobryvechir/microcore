/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

const (
	COMPARE_NOT = iota
	COMPARE_AS_NUMBERS
	COMPARE_AS_INTEGERS
	COMPARE_AS_STRINGS
)

func CompareAsIntegers(value1 interface{}, value2 interface{}) int {
	a, oka := AnyToNumberInt(value1)
	b, okb := AnyToNumberInt(value2)
	if !oka || !okb {
		return CompareAsNumbers(value1, value2)
	}
	res := 0
	if a != b {
		if a > b {
			res = 1
		} else {
			res = -1
		}
	}
	return res
}

func CompareAsStrings(value1 interface{}, value2 interface{}) int {
	a := AnyToString(value1)
	b := AnyToString(value2)
	res := 0
	if a != b {
		if a > b {
			res = 1
		} else {
			res = -1
		}
	}
	return res
}

func CompareAsNumbers(value1 interface{}, value2 interface{}) int {
	a := AnyToNumber(value1)
	b := AnyToNumber(value2)
	res := 0
	if a != b {
		if a > b {
			res = 1
		} else {
			res = -1
		}
	}
	return res
}

func CompareByComparisonType(compareType int, value1 interface{}, value2 interface{}) int {
	res := -2
	switch compareType {
	case COMPARE_AS_INTEGERS:
		res = CompareAsIntegers(value1, value2)
	case COMPARE_AS_NUMBERS:
		res = CompareAsNumbers(value1, value2)
	case COMPARE_AS_STRINGS:
		res = CompareAsStrings(value1, value2)
	}
	return res
}

func getCompareType(kind int) int {
	res := COMPARE_NOT
	switch kind {
	case dvgrammar.TYPE_STRING:
		res = COMPARE_AS_STRINGS
	case dvgrammar.TYPE_NUMBER, dvgrammar.TYPE_NAN:
		res = COMPARE_AS_NUMBERS
	case dvgrammar.TYPE_CHAR,
		dvgrammar.TYPE_NUMBER_INT,
		dvgrammar.TYPE_BOOLEAN,
		dvgrammar.TYPE_NULL:
		res = COMPARE_AS_INTEGERS
	}
	return res
}

//returns 1-greater, -1-lesser, 0-equal, -2 not applied
func AnyCompareAnyWithTypes(kind1 int, value1 interface{}, kind2 int, value2 interface{}) int {
	if value1 == value2 {
		return 0
	}
	compareType := getCompareType(kind1)
	compareType1 := getCompareType(kind2)
	if compareType1 < compareType {
		compareType = compareType1
	}
	res := CompareByComparisonType(compareType, value1, value2)
	if res == 0 && kind1 != kind2 {
		if kind1 == dvgrammar.TYPE_NULL || kind2 == dvgrammar.TYPE_NULL {
			res = -2
		}
	}
	return res
}
