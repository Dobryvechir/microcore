/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type ExecCallConfig struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
	Result string            `json:"result"`
}

var execStatementProcessFunctions = map[string]ProcessFunction{
	CommandCall:        {Init: execCallInit, Run: execCallRun},
	CommandIf:          {Init: execIfInit, Run: execIfRun},
	CommandFor:         {Init: execForInit, Run: execForRun},
	CommandRange:       {Init: execRangeInit, Run: execRangeRun},
	CommandSwitch:      {Init: execSwitchInit, Run: execSwitchRun},
	CommandReturn:      {Init: execReturnInit, Run: execReturnRun},
}

func execCallInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ExecCallConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	if config.Action == "" {
		log.Printf("action must be specified in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func execCallRun(data []interface{}) bool {
	config := data[0].(*ExecCallConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ExecCall(config, ctx)
}

func ExecCall(config *ExecCallConfig, ctx *dvcontext.RequestContext) bool {
	if config.Action != "" {
		ExecuteAddSubsequence(ctx, config.Action, config.Params, config.Result)
	}
	return true
}

type ExecIfConfig struct {
	Condition string          `json:"condition"`
	Then      *ExecCallConfig `json:"then"`
	Else      *ExecCallConfig `json:"else"`
}

func execIfInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ExecIfConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	if config.Condition == "" {
		log.Printf("condition must be specified in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func execIfRun(data []interface{}) bool {
	config := data[0].(*ExecIfConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ExecIf(config, ctx)
}

func ExecIf(config *ExecIfConfig, ctx *dvcontext.RequestContext) bool {
	res, err := ctx.LocalContextEnvironment.EvaluateBooleanExpression(config.Condition)
	if err != nil {
		dvlog.PrintlnError("Error in " + config.Condition)
		return true
	}
	if res {
		ExecCall(config.Then, ctx)
	} else {
		ExecCall(config.Else, ctx)
	}
	return true
}

type ExecForConfig struct {
	Params     map[string]string `json:"initial"`
	Condition  string            `json:"condition"`
	Next       map[string]string `json:"next"`
	BodyAction string            `json:"body"`
	Result     string            `json:"result"`
}

func execForInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ExecForConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func execForRun(data []interface{}) bool {
	config := data[0].(*ExecForConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ExecFor(config, ctx)
}

func ExecuteCalculate(ctx *dvcontext.RequestContext, data map[string]string) {
	if data != nil {
		for k, v := range data {
			if v == "" {
				ctx.LocalContextEnvironment.Set(k, v)
			}
			value, err := ctx.LocalContextEnvironment.EvaluateAnyTypeExpression(v)
			if err != nil {
				dvlog.PrintlnError("Error evaluating " + v + " (" + err.Error() + ")")
				return
			}
			ctx.LocalContextEnvironment.Set(k, value)
		}
	}
}

func ExecFor(config *ExecForConfig, ctx *dvcontext.RequestContext) bool {
	ExecuteAddSubsequence(ctx, config.BodyAction, config.Params, config.Result)
	level := ctx.PrimaryContextEnvironment.GetInt(ExSeqLevel)
	env := ctx.LocalContextEnvironment
	for true {
		if config.Condition != "" {
			res, err := ctx.LocalContextEnvironment.EvaluateBooleanExpression(config.Condition)
			if err != nil {
				dvlog.PrintlnError("Error evaluating for condition: " + config.Condition)
			}
			if !res {
				break
			}
		}
		ExecuteSequenceCycle(ctx, level)
		ctx.LocalContextEnvironment = env
		ctx.PrimaryContextEnvironment.Set(ExSeqLevel, level)
		ExecuteCalculate(ctx, config.Next)
	}
	ExecReturnShort(config.Result, ctx)
	return true
}

type ExecRangeConfig struct {
	Params     map[string]string `json:"params"`
	MapObject  string            `json:"range"`
	BodyAction string            `json:"body"`
	Result     string            `json:"result"`
}

func execRangeInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ExecRangeConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func execRangeRun(data []interface{}) bool {
	config := data[0].(*ExecRangeConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ExecRange(config, ctx)
}

func ExecRange(config *ExecRangeConfig, ctx *dvcontext.RequestContext) bool {
	ExecuteAddSubsequence(ctx, config.BodyAction, config.Params, config.Result)
	level := ctx.PrimaryContextEnvironment.GetInt(ExSeqLevel)
	mapObj, err := ctx.LocalContextEnvironment.EvaluateAnyTypeExpression(config.MapObject)
	env := ctx.LocalContextEnvironment
	if err != nil {
		dvlog.PrintlnError("Error evaluating " + config.MapObject + " (" + err.Error() + ")")
	} else {
		dvjson.IterateOnAnyType(mapObj, func(key string, val interface{}, index int, previous interface{}) (interface{}, bool) {
			env.Set(key, val)
			env.Set("FOR_INDEX", index)
			ExecuteSequenceCycle(ctx, level)
			ctx.LocalContextEnvironment = env
			ctx.PrimaryContextEnvironment.Set(ExSeqLevel, level)
			return nil, true
		}, nil)
	}
	ExecReturnShort(config.Result, ctx)
	return true
}

type ExecSwitchConfig struct {
	DefaultAction string            `json:"defaultAction"`
	Params        map[string]string `json:"params"`
	Cases         map[string]string `json:"cases"`
	Result        string            `json:"result"`
}

func execSwitchInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ExecSwitchConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func execSwitchRun(data []interface{}) bool {
	config := data[0].(*ExecSwitchConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ExecSwitch(config, ctx)
}

func ExecSwitch(config *ExecSwitchConfig, ctx *dvcontext.RequestContext) bool {
	actionName := config.DefaultAction
	if config.Cases != nil {
		for k, v := range config.Cases {
			res, err := ctx.PrimaryContextEnvironment.EvaluateBooleanExpression(k)
			if err != nil {
				log.Printf("Failed to evaluate %s: %v", k, err)
				ctx.HandleInternalServerError()
				return true
			}
			if res {
				actionName = v
				break
			}
		}
	}
	ExecCall(&ExecCallConfig{Action: actionName, Params: config.Params, Result: config.Result}, ctx)
	return true
}

type ExecReturnConfig struct {
	Result string `json:"result"`
}

func execReturnInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ExecReturnConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func execReturnRun(data []interface{}) bool {
	config := data[0].(*ExecReturnConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ExecReturn(config, ctx)
}

func ExecReturn(config *ExecReturnConfig, ctx *dvcontext.RequestContext) bool {
	return ExecReturnShort(config.Result, ctx)
}

func ExecReturnShort(result string, ctx *dvcontext.RequestContext) bool {
	ExecuteReturnSubsequence(ctx, result)
	return true
}
