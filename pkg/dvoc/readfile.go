/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"io/ioutil"
	"log"
	"os"
)

type ReadFileConfig struct {
	FileName          string `json:"name"`
	Kind              string `json:"kind"`
	Result            string `json:"result"`
	Trim              string `json:"trim"`
	NoReadOfUndefined bool   `json:"noReadOfUndefined"`
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
		res, err1 = ReadJsonTrimmed(dat, config.Trim, config.NoReadOfUndefined)
	}
	return ProcessSavingActionResult(config.Result, res, ctx, err1, "in file ", config.FileName)
}

func ReadJsonTrimmed(data []byte, trim string, noReadOfUndefined bool) (interface{}, error) {
	item, err:=dvjson.ReadJsonChild(data, trim, noReadOfUndefined)
	if err!=nil {
		return nil, err
	}
	if item==nil {
		return nil, nil
	}
	return item, nil
}
