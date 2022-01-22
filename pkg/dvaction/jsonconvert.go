/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type JsonConvertConfig struct {
	Source  *JsonRead   `json:"source"`
	Result  string      `json:"result"`
	Add     []*JsonRead `json:"add"`
	Merge   []*JsonRead `json:"merge"`
	Replace []*JsonRead `json:"replace"`
	Update  []*JsonRead `json:"update"`
	Push    []*JsonRead `json:"push"`
	Concat  []*JsonRead `json:"concat"`
	Remove  []string    `json:"remove"`
}

func jsonConvertInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &JsonConvertConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Source == nil || config.Source.Var == "" {
		log.Printf("source must be present in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func jsonConvertRun(data []interface{}) bool {
	config := data[0].(*JsonConvertConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return JsonConvertRunByConfig(config, ctx)
}

func updateVariablesByConfig(config []*JsonRead, mode int, src interface{}, env *dvevaluation.DvObject) (interface{}, bool) {
	n := len(config)
	for i := 0; i < n; i++ {
		conf := config[i]
		v, err := JsonExtract(conf, env)
		if err != nil {
			dvlog.PrintlnError("Error in json extracting by " + conf.Var)
			return src, false
		}
		src = dvevaluation.UpdateAnyVariables(src, v, conf.Destination,
			mode, conf.Ids, env)
	}
	return src, true
}

func JsonConvertRunByConfig(config *JsonConvertConfig, ctx *dvcontext.RequestContext) bool {
	env := GetEnvironment(ctx)
	src, err := JsonExtract(config.Source, env)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Source.Var)
		return true
	}
	var ok bool
	src, ok = updateVariablesByConfig(config.Add, dvevaluation.UPDATE_MODE_ADD_BY_KEYS, src, env)
	if !ok {
		return true
	}
	src, ok = updateVariablesByConfig(config.Merge, dvevaluation.UPDATE_MODE_MERGE, src, env)
	if !ok {
		return true
	}
	src, ok = updateVariablesByConfig(config.Replace, dvevaluation.UPDATE_MODE_REPLACE, src, env)
	if !ok {
		return true
	}
	src, ok = updateVariablesByConfig(config.Update, dvevaluation.UPDATE_MODE_APPEND, src, env)
	if !ok {
		return true
	}
	n := len(config.Remove)
	for i := 0; i < n; i++ {
		v := config.Remove[i]
		src = dvevaluation.RemoveAnyVariable(src, v, env)
	}
	s := dvevaluation.AnyToDvVariable(src)
	n = len(config.Push)
	for i := 0; i < n; i++ {
		JsonConvertPush(config.Push[i], s, env)
	}
	n = len(config.Concat)
	for i := 0; i < n; i++ {
		JsonConvertConcat(config.Concat[i], s, env)
	}
	SaveActionResult(config.Result, s, ctx)
	return true
}

func JsonConvertPush(push *JsonRead, dst *dvevaluation.DvVariable, env *dvevaluation.DvObject) {
	src, err := JsonExtract(push, env)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + push.Var)
		return
	}
	if src != nil && dst != nil && dst.Kind == dvevaluation.FIELD_ARRAY {
		s := dvevaluation.AnyToDvVariable(src)
		dst.Fields = append(dst.Fields, s)
	}
}

func JsonConvertConcat(push *JsonRead, dst *dvevaluation.DvVariable, env *dvevaluation.DvObject) {
	src, err := JsonExtract(push, env)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + push.Var)
		return
	}
	if src != nil && dst != nil && dst.Kind == dvevaluation.FIELD_ARRAY {
		s := dvevaluation.AnyToDvVariable(src)
		if s!=nil && s.Kind==dvevaluation.FIELD_ARRAY {
			n:=len(s.Fields)
			for i:=0;i<n;i++ {
				dst.Fields = append(dst.Fields, s.Fields[i])
			}
		}
	}
}
