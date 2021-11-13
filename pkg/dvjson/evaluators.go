/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

type ExpressionResolver func(string, map[string]interface{}, int) (string, error)

const (
	EXPRESSION_RESOLVER_CACHE = 1 << iota
)

func (item *DvFieldInfo) ReadChildOfAnyLevel(name string, props *dvevaluation.DvObject) (res *DvFieldInfo, err error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 || item==nil {
		return item, nil
	}
	defProps:=make(map[string]interface{})
	if props==nil {
		props = dvparser.GetPropertiesPrototypedToGlobalProperties(defProps);
	} else {
		props = dvevaluation.NewObjectWithPrototype(defProps,props)
	}
	res, err = item.ReadChild(name, func(expr string, data map[string]interface{},options int) (string, error) {
		if data==nil {
			data = defProps
		}
		props.Properties = data
		return dvevaluation.ParseForDvObjectString(expr, props)
	})
	return
}
