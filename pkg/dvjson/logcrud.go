/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"strconv"
)

func logInfoForItem(item *DvCrudItem) (res string) {
	if item == nil {
		res = "<item is nil>"
	} else {
		res = "[" + string(item.itemBody) + "](Id:" + string(item.Id) + " fields:("
		for _, k := range item.Fields {
			res += string(k.Name) + ":" + string(k.Value) + " body: " + strconv.Itoa(k.posStart) + " start:" +
				strconv.Itoa(k.valueStartPos) + " end:" + strconv.Itoa(k.valueEndPos) + ","
		}
		res += "))"
	}
	return
}
