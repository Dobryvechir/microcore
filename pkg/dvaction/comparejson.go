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

type JsonRead struct {
	Var               string   `json:"var"`
	Path              string   `json:"path"`
	Filter            string   `json:"filter"`
	Sort              []string `json:"sort"`
	AfterPath         string   `json:"afterPath"`
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

func JsonExtract(info *JsonRead, ctx *dvcontext.RequestContext) (interface{}, error) {
	return JsonExtractExtended(info.Var, info.Path, info.Filter, info.Sort, info.AfterPath,
		info.NoReadOfUndefined, info.ErrorSignificant, info.Ids, info.Convert, ctx)
}

func JsonExtractExtended(place string, path string, filter string, sort []string,
	afterPath string, noReadOfUndefined bool, errorSignificant bool, ids []string, convert string,
	ctx *dvcontext.RequestContext) (interface{}, error) {
	val, ok := ReadActionResult(place, ctx)
	env := GetEnvironment(ctx)
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
	if afterPath != "" {
		item, _, err := dvjson.ReadPathOfAny(val, afterPath, noReadOfUndefined, env)
		if err != nil {
			return nil, err
		}
		val = item
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
	sample, err := JsonExtract(config.Sample, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Sample.Var)
		return true
	}
	ref, err := JsonExtract(config.Ref, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Ref.Var)
		return true
	}
	added, removed, updated, unchanged, counterparts := dvjson.FindDifferenceForAnyType(sample, ref,
		config.Added != "", config.Removed != "", config.Updated != "", config.Unchanged != "",
		config.UpdatedRef != "", config.UnchangedAsUpdated, false)
	if config.Added != "" {
		SaveActionResult(config.Added, added, ctx)
	}
	if config.Removed != "" {
		SaveActionResult(config.Removed, removed, ctx)
	}
	if config.Updated != "" {
		SaveActionResult(config.Updated, updated, ctx)
	}
	if config.Unchanged != "" {
		SaveActionResult(config.Unchanged, unchanged, ctx)
	}
	if config.UpdatedRef != "" {
		SaveActionResult(config.UpdatedRef, counterparts, ctx)
	}
	return true
}
