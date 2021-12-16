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

type JsonRead struct {
	Var               string   `json:"var"`
	Path              string   `json:"path"`
	Sort              []string `json:"sort"`
	Filter            string   `json:"filter"`
	NoReadOfUndefined bool     `json:"noReadOfUndefined"`
	ErrorSignificant  bool     `json:"errorSignificant"`
	Convert           string   `json:"convert"`
	Ids               []string `json:"ids"`
	Destination       string   `json:"destination"`
}

type CompareJsonConfig struct {
	Sample             *JsonRead `json:"sample"`
	Ref                *JsonRead `json:"ref"`
	Added              string    `json:"added"`
	Removed            string    `json:"removed"`
	Updated            string    `json:"updated"`
	UpdatedRef         string    `json:"updatedRef"`
	Unchanged          string    `json:"unchanged"`
	UnchangedAsUpdated bool      `json:"unchangedAsUpdated"`
}

func compareJsonInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &CompareJsonConfig{}
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
	return []interface{}{config, ctx}, true
}

func compareJsonRun(data []interface{}) bool {
	config := data[0].(*CompareJsonConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return CompareJsonByConfig(config, ctx)
}

func JsonExtract(info *JsonRead, env *dvevaluation.DvObject) (interface{}, error) {
	return JsonExtractExtended(info.Var, info.Path, info.Sort, info.Filter,
		info.NoReadOfUndefined, info.ErrorSignificant, info.Ids, info.Convert, env)
}

func JsonExtractExtended(place string, path string, sort []string, filter string,
	noReadOfUndefined bool, errorSignificant bool, ids []string, convert string,
	env *dvevaluation.DvObject) (interface{}, error) {
	val, ok := env.Get(place)
	if !ok {
		return nil, nil
	}
	if path != "" {
		item, _, err := dvjson.ReadPathOfAny(val, path, noReadOfUndefined, env)
		if err != nil {
			return nil, err
		}
		val = item
	}
	if filter != "" {
		res, err := dvjson.IterateFilterByExpression(val, filter, env, errorSignificant)
		if err != nil {
			return nil, err
		}
		val = res
	}
	if len(sort) > 0 {
		res, err := dvjson.IterateSortByFields(val, sort, env)
		if err != nil {
			return nil, err
		}
		val = res
	}
	dvjson.CreateQuickInfoByKeysForAny(val, ids)
	if convert != "" {
		env.Set("arg", val)
		res, err := env.EvaluateAnyTypeExpression(convert)
		return res, err
	}
	return val, nil
}

func CompareJsonByConfig(config *CompareJsonConfig, ctx *dvcontext.RequestContext) bool {
	sample, err := JsonExtract(config.Sample, ctx.LocalContextEnvironment)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Sample.Var)
		return true
	}
	ref, err := JsonExtract(config.Ref, ctx.LocalContextEnvironment)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Ref.Var)
		return true
	}
	added, removed, updated, unchanged, counterparts := dvjson.FindDifferenceForAnyType(sample, ref,
		config.Added != "", config.Removed != "", config.Updated != "", config.Unchanged != "",
		config.UpdatedRef != "", config.UnchangedAsUpdated)
	env := ctx.LocalContextEnvironment.Properties
	if added != nil {
		env[config.Added] = added
	}
	if removed != nil {
		env[config.Removed] = removed
	}
	if updated != nil {
		env[config.Updated] = updated
	}
	if unchanged != nil {
		env[config.Unchanged] = unchanged
	}
	if counterparts != nil {
		env[config.UpdatedRef] = counterparts
	}
	return true
}
