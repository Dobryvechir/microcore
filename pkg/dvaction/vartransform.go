/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type VarTransformConfig struct {
	Transform map[string]string `json:"transform"`
}

func varTransformInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &VarTransformConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	if config.Transform == nil {
		log.Printf("transform must be specified in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func varTransformRun(data []interface{}) bool {
	config := data[0].(*VarTransformConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return VarTransformRunByConfig(config, ctx)
}

func VarTransformRunByConfig(config *VarTransformConfig, ctx *dvcontext.RequestContext) bool {
	for k, v := range config.Transform {
		r, err := ctx.LocalContextEnvironment.EvaluateAnyTypeExpression(v)
		if err != nil {
			dvlog.PrintlnError("Error in expression " + v + ":" + err.Error())
		} else {
			ctx.LocalContextEnvironment.Set(k, r)
		}
	}
	return true
}
