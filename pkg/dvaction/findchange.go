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

type FindChangeConfig struct {
	Sample       *JsonRead `json:"sample"`
	Ref          *JsonRead `json:"ref"`
	Result       string    `json:"result"`
	Algorithm    string    `json:"algorithm"`
	AddFieldName string    `json:"add_field_name"`
}

func findChangeInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &FindChangeConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Sample == nil || config.Sample.Var == "" {
		log.Printf("sample.place must be specified in %s", command)
		return nil, false
	}
	if config.Ref == nil || config.Ref.Var == "" {
		log.Printf("reference.place must be present in %s", command)
		return nil, false
	}
	if config.Result == "" {
		log.Printf("result must be present and positive in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func findChangeRun(data []interface{}) bool {
	config := data[0].(*FindChangeConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return FindChangeByConfig(config, ctx)
}

func FindChangeByConfig(config *FindChangeConfig, ctx *dvcontext.RequestContext) bool {
	src, err := JsonExtract(config.Sample, ctx.LocalContextEnvironment)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Sample.Var)
		return true
	}
	ref, err := JsonExtract(config.Ref, ctx.LocalContextEnvironment)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Ref.Var)
		return true
	}
	env := ctx.LocalContextEnvironment
	res, err := dvjson.FindChangeAny(src, ref, config.AddFieldName, config.Algorithm, env)
	if err != nil {
		dvlog.PrintlnError("Error in find change " + err.Error())
		return true
	}
	env.Set(config.Result, res)
	return true
}
