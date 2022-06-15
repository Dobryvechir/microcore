/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"io/ioutil"
	"log"
)

type JsConfig struct {
	File   string `json:"file"`
	Result string `json:"result"`
}

func jsInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &JsConfig{}
	r, s := DefaultOrSimpleInitWithObject(command, config, GetEnvironment(ctx))
	if !r {
		return nil, false
	}
	if s != "" {
		config.File = s
	}
	if config.File == "" {
		log.Printf("js.file must be specified in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func jsRun(data []interface{}) bool {
	config := data[0].(*JsConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return JsRunByConfig(config, ctx)
}

func JsRunByConfig(config *JsConfig, ctx *dvcontext.RequestContext) bool {
	data, err := ioutil.ReadFile(config.File)
	if err != nil {
		message := "Error in reading " + config.File + ": " + err.Error()
		ActionInternalException(500, message, message, ctx)
	}
	v, err := GetEnvironment(ctx).EvaluateAnyTypeExpression(string(data))
	if err != nil {
		message := "Error in reading " + config.File + ": " + err.Error()
		ActionInternalException(500, message, message, ctx)
	}
	if config.Result != "" {
		SaveActionResult(config.Result, v, ctx)
	}
	return true
}
