/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strings"
)

func (item *DvFieldInfo) ContainsItemIn(v interface{}) bool {
	if item==nil {
		return false
	}
	s:=dvevaluation.AnyToString(v)
	n:=len(item.Fields)
	switch item.Kind {
	case dvevaluation.FIELD_OBJECT:
		for i:=0;i<n;i++ {
			if string(item.Fields[i].Name) == s {
				return true
			}
		}
	case dvevaluation.FIELD_ARRAY:
		for i:=0;i<n;i++ {
			f:=item.Fields[i]
			if f.Kind!=dvevaluation.FIELD_ARRAY && f.Kind!=dvevaluation.FIELD_OBJECT && string(f.Value) == s {
				return true
			}
		}
	default:
		return strings.Contains(string(item.Value), s)
	}
	return false
}

func EvaluationContainInProcessor(contained interface{}, containing interface{}) (bool, bool, error) {
	switch containing.(type) {
	case *DvFieldInfo:
		res := containing.(*DvFieldInfo).ContainsItemIn(contained)
		return res, true, nil
	}
	return false, false, nil
}

func evaluationRegistrations() bool {
	dvevaluation.RegisterContainInProcessor(EvaluationContainInProcessor)
	return true
}

var inited = evaluationRegistrations()
