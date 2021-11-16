/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"os"
)

type DebugConfig struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

type DebugActionFn func(*DebugConfig, *dvcontext.RequestContext) error

var debugActionList = map[string]DebugActionFn{
	"ALL_VARIABLES":    RunReadAllVariables,
	"GET_ONE_VARIABLE": RunGetOneVariable,
	"SET_ONE_VARIABLE": RunSetOneVariable,
	"EXIT":             RunExit,
}

func debugInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &DebugConfig{}
	key, okay := dvparser.GlobalProperties["DEBUG_KEY"]
	if !okay || key == "" || !DefaultInitWithObject(command, config) {
		return nil, false
	}
	currentKey, okey := ctx.PrimaryContextEnvironment.Properties["URL_PARAM_KEY"]
	if !okey {
		currentKey, okey = ctx.PrimaryContextEnvironment.Properties["BODY_PARAM_KEY"]
		if !okey {
			return nil, false
		}
	}
	switch currentKey.(type) {
	case string:
		if currentKey.(string) != key {
			return nil, false
		}
	default:
		return nil, false
	}
	fn, ok := debugActionList[config.Action]
	if !ok {
		log.Printf("action must be specified properly in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx, fn}, true
}

func debugRun(data []interface{}) bool {
	config := data[0].(*DebugConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	fn := data[2].(DebugActionFn)
	err := fn(config, ctx)
	if err != nil {
		log.Printf("Error %v", err)
		return false
	}
	return true
}

func RunReadAllVariables(config *DebugConfig, ctx *dvcontext.RequestContext) error {
	result := config.Data
	if result == "" {
		result = "DEBUG_ALL_VARIABLES"
	}
	props := dvparser.GlobalProperties
	ctx.PrimaryContextEnvironment.Properties[result] = props
	return nil
}

func RunGetOneVariable(config *DebugConfig, ctx *dvcontext.RequestContext) error {
	result := config.Data
	if result == "" {
		result = "DEBUG_ONE_VARIABLE"
	}
	nameStr, ok := ctx.PrimaryContextEnvironment.Properties["URL_PARAM_NAME"]
	if !ok {
		return nil
	}
	var name string
	switch nameStr.(type) {
	case string:
		name = nameStr.(string)
	default:
		return nil
	}
	val := dvparser.GlobalProperties[name]
	ctx.PrimaryContextEnvironment.Properties[result] = val
	return nil
}

func RunSetOneVariable(config *DebugConfig, ctx *dvcontext.RequestContext) error {
	result := config.Data
	if result == "" {
		result = "DEBUG_ONE_VARIABLE"
	}
	nameStr, ok := ctx.PrimaryContextEnvironment.Properties["BODY_PARAM_NAME"]
	valueStr, ok1 := ctx.PrimaryContextEnvironment.Properties["BODY_PARAM_VALUE"]
	if !ok || !ok1 {
		return nil
	}
	var name, value string
	switch nameStr.(type) {
	case string:
		name = nameStr.(string)
	default:
		return nil
	}
	switch valueStr.(type) {
	case string:
		value = valueStr.(string)
	default:
		return nil
	}
	val := dvparser.GlobalProperties[name]
	ctx.PrimaryContextEnvironment.Properties[result] = val
	dvparser.GlobalProperties[name] = value
	return nil
}

func RunExit(config *DebugConfig, ctx *dvcontext.RequestContext) error {
	log.Printf("Exit command fired")
	os.Exit(7)
	return nil
}
