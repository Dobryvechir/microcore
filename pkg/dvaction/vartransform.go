/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type VarTransformConfig struct {
	Transform     map[string]string    `json:"transform"`
	Clone         map[string]string    `json:"clone"`
	JsonParse     map[string]string    `json:"parse"`
	Read          map[string]*JsonRead `json:"read"`
	DefaultString map[string]string    `json:"default_string"`
	DefaultAny    map[string]string    `json:"default_any"`
}

func varTransformInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &VarTransformConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Transform == nil && config.Read == nil {
		log.Printf("transform or read must be specified in %s", command)
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
	if config.JsonParse != nil {
		for k, v := range config.JsonParse {
            j,_:=ctx.LocalContextEnvironment.Get(v)
			r, err := dvjson.ParseAny(j)
			if err != nil {
				dvlog.PrintlnError("Error in expression " + k + ":" + err.Error())
			} else {
				SaveActionResult(k, r, ctx)
			}
		}
	}
	if config.Read != nil {
		for k, v := range config.Read {
			r, err := JsonExtract(v, ctx.LocalContextEnvironment)
			if err != nil {
				dvlog.PrintlnError("Error in expression " + k + ":" + err.Error())
			} else {
				SaveActionResult(k, r, ctx)
			}
		}
	}
	if config.Transform != nil {
		for k, v := range config.Transform {
			r, err := ctx.LocalContextEnvironment.EvaluateAnyTypeExpression(v)
			if err != nil {
				dvlog.PrintlnError("Error in expression " + v + ":" + err.Error())
			} else {
				SaveActionResult(k, r, ctx)
			}
		}
	}
	if config.Clone != nil {
		for k, v := range config.Clone {
			r, ok := ctx.LocalContextEnvironment.Get(v)
			if ok {
				d := dvevaluation.AnyToDvVariable(r)
				d = d.Clone()
				SaveActionResult(k, d, ctx)
			}
		}
	}
	if config.DefaultString != nil {
		for k, v := range config.DefaultString {
			_, ok := ctx.LocalContextEnvironment.Get(k)
			if !ok {
				SaveActionResult(k, v, ctx)
			}
		}
	}
	if config.DefaultAny != nil {
		for k, v := range config.DefaultAny {
			_, ok := ctx.LocalContextEnvironment.Get(k)
			if !ok {
				r, err := dvjson.JsonFullParser([]byte(v))
				if err != nil {
					dvlog.PrintlnError("Error in json " + v + ":" + err.Error())
				}
				SaveActionResult(k, r, ctx)
			}
		}
	}
	return true
}
