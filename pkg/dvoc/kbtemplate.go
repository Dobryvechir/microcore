/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
)

type MicroServiceSearch struct {
	MicroServices string `json:"services"`
    WorkFolder string `json:"folder"`
    Options string `json:"options"`
    Output string `json:"output"`
}

func MicroServiceTemplateInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &MicroServiceSearch{}
	if !dvaction.DefaultInitWithObject(command, config) {
		return nil, false
	}
	return nil, false
}

func MicroServiceTemplateRun(data []interface{}) bool {
	return false
}
