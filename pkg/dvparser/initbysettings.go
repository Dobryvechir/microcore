/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvparser

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation/dvfunctions"
)

type InitBySettingFunc func(parameters map[string]string, functionPool map[string]interface{})

var initialiazerPool []InitBySettingFunc = make([]InitBySettingFunc, 0, 16)

func RegisterInitBySettingFunc(initializer InitBySettingFunc) bool {
	initialiazerPool = append(initialiazerPool, initializer)
	return true
}

func CallInitBySettingFunc(parameters map[string]string, functionPool map[string]interface{}) {
	for _, v := range initialiazerPool {
		v(parameters, functionPool)
	}
}

func initializeRegisteredFunctions() {
	dvevaluation.AddToGlobalFunctionPool(dvfunctions.RegisteredFunctions)
}

func CallInitBySettingFuncDefault() {
	CallInitBySettingFunc(GlobalProperties, dvevaluation.GlobalFunctionPool)
}
