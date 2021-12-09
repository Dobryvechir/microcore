/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ReadFileConfig struct {
	FileName          string   `json:"name"`
	Kind              string   `json:"kind"`
	Result            string   `json:"result"`
	Path              string   `json:"path"`
	Filter            string   `json:"filter"`
	Sort              []string `json:"sort"`
	NoReadOfUndefined bool     `json:"noReadOfUndefined"`
	ErrorSignificant  bool     `json:"errorSignificant"`
	IsTemplate        bool     `json:"template"`
}

func readFileActionInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ReadFileConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.FileName == "" {
		log.Printf("File name must be specified in %s", command)
		return nil, false
	}
	if config.Kind == "" {
		config.Kind = "json"
	}
	if config.Kind != "json" && config.Kind != "string" && config.Kind != "text" {
		log.Printf("Supported file kind is not supported %s (available kind options: json)", command)
		return nil, false
	}
	if config.Result == "" {
		log.Printf("Result name is not specified in command %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func readFileActionRun(data []interface{}) bool {
	config := data[0].(*ReadFileConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ReadFileByConfigKind(config, ctx)
}

func ReadFileByConfigKind(config *ReadFileConfig, ctx *dvcontext.RequestContext) bool {
	_, err := os.Stat(config.FileName)
	if err != nil {
		log.Printf("File not found %s", config.FileName)
		return false
	}
	var dat []byte
	env := GetEnvironment(ctx)
	if config.IsTemplate {
		dat, err = dvparser.SmartReadLikeJsonTemplate(config.FileName, 3, env)
	} else {
		dat, err = ioutil.ReadFile(config.FileName)
	}
	if err != nil {
		log.Printf("Error reading file %s %v", config.FileName, err)
		return false
	}
	var res interface{}
	switch config.Kind {
	case "json":
		res, err = ReadJsonTrimmed(dat, config.Path, config.NoReadOfUndefined, env)
		if err == nil && config.Filter != "" {
			res, err = dvjson.IterateFilterByExpression(res, config.Filter, env, config.ErrorSignificant)
		}
		if err == nil && len(config.Sort) > 0 {
			res, err = dvjson.IterateSortByFields(res, config.Sort, env)
		}
	case "text":
		res = string(dat)
	case "string":
		res = strings.TrimSpace(string(dat))
	}
	return ProcessSavingActionResult(config.Result, res, ctx, err, "in file ", config.FileName)
}

func ReadJsonTrimmed(data []byte, path string, noReadOfUndefined bool, env *dvevaluation.DvObject) (interface{}, error) {
	item, err := dvjson.ReadJsonChild(data, path, noReadOfUndefined, env)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	return item, nil
}
