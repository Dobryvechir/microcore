/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	if config.Kind != "json" && config.Kind != "string" && config.Kind != "text" && config.Kind != "remove" && config.Kind!="binary" {
		log.Printf("Supported file kind is not supported %s (available kind options: json)", command)
		return nil, false
	}
	if config.Result == "" && config.Kind != "remove" {
		log.Printf("Result is not specified in command %s", command)
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
	fileNames := dvdir.ReadFileList(config.FileName)
	n := len(fileNames)
	if config.Kind == "remove" {
		k := dvdir.DeleteFilesIfExist(fileNames)
		SaveActionResult(config.Result, strconv.Itoa(k)+"/"+strconv.Itoa(n), ctx)
		return true
	}
	isMulti := n > 1
	if isMulti {
		DeleteActionResult(config.Result, ctx)
	}
	for i := 0; i < n; i++ {
		ProcessSingleFile(fileNames[i], config, ctx, isMulti)
	}
	return true
}

func ProcessSingleFile(fileName string, config *ReadFileConfig, ctx *dvcontext.RequestContext, isMulti bool) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		log.Printf("File not found %s", config.FileName)
		return false
	}
	var dat []byte
	env := GetEnvironment(ctx)
	if config.IsTemplate {
		dat, err = dvparser.SmartReadLikeJsonTemplate(fileName, 3, env)
	} else {
		dat, err = ioutil.ReadFile(fileName)
	}
	if err != nil {
		log.Printf("Error reading file %s %v", fileName, err)
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
		s := strings.TrimSpace(string(dat))
		if isMulti {
			prev, ok := ReadActionResult(config.Result, ctx)
			if ok && prev != nil {
				res = append(prev.([]byte), dat...)
			} else {
				res = dat
			}
		} else {
			res = s
		}
	case "binary":
		if isMulti {
			prev, ok := ReadActionResult(config.Result, ctx)
			if ok && prev != nil {
				res = append(prev.([]byte), dat...)
			} else {
				res = dat
			}
		} else {
			res = dat
		}
	}
	return ProcessSavingActionResult(config.Result, res, ctx, err, "in file ", fileName)
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
