/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"io/ioutil"
	"log"
	"os"
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
}

func readFileActionInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ReadFileConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	if config.FileName == "" {
		log.Printf("File name must be specified in %s", command)
		return nil, false
	}
	if config.Kind == "" {
		config.Kind = "json"
	}
	if config.Kind != "json" {
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
	if _, err := os.Stat(config.FileName); err != nil {
		log.Printf("File not found %s", config.FileName)
		return false
	}
	dat, err1 := ioutil.ReadFile(config.FileName)
	if err1 != nil {
		log.Printf("Error reading file %s %v", config.FileName, err1)
		return false
	}
	var res interface{}
	switch config.Kind {
	case "json":
		var props *dvevaluation.DvObject = nil
		if ctx != nil {
			props = ctx.PrimaryContextEnvironment
		}
		res, err1 = ReadJsonTrimmed(dat, config.Path, config.NoReadOfUndefined, props)
		if err1 == nil && config.Filter != "" {
			res, err1 = dvjson.IterateFilterByExpression(res, config.Filter, ctx.LocalContextEnvironment, config.ErrorSignificant)
		}
		if err1 == nil && len(config.Sort) > 0 {
			res, err1 = dvjson.IterateSortByFields(res, config.Sort, ctx.LocalContextEnvironment)
		}
	}
	return ProcessSavingActionResult(config.Result, res, ctx, err1, "in file ", config.FileName)
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