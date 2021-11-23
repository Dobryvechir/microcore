/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)



func (parseInfo *DvCrudParsingInfo) ConvertSimpleValueToInterface() (interface{}, bool) {
	return dvevaluation.ConvertSimpleKindAndValueToInterface(parseInfo.Kind, parseInfo.Value)
}

func DvFieldArrayToBytes(val []*dvevaluation.DvVariable) []byte {
	buf := make([]byte, 1, 102400)
	buf[0] = '['
	n := len(val)
	for i := 0; i < n; i++ {
		if i != 0 {
			buf = append(buf, ',')
		}
		b := PrintToJson(val[i], 2)
		buf = append(buf, b...)
	}
	buf = append(buf, ']')
	return buf
}

func DvFieldInfoToStringConverter(v interface{}) (string, bool) {
	switch v.(type) {
	case *dvevaluation.DvVariable:
		return string(PrintToJson(v.(*dvevaluation.DvVariable), 2)), true
	case []*dvevaluation.DvVariable:
		return string(DvFieldArrayToBytes(v.([]*dvevaluation.DvVariable))), true
	}
	return "", false
}
