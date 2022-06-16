/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
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
	Joiner            string   `json:"joiner"`
	Append            int      `json:"append"`
	NoReadOfUndefined bool     `json:"noReadOfUndefined"`
	ErrorSignificant  bool     `json:"errorSignificant"`
	IsTemplate        bool     `json:"template"`
	EolJoiner         bool     `json:"eol_joiner"`
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
	if config.Kind != "json" && config.Kind != "string" && config.Kind != "text" && config.Kind != "remove" && config.Kind != "binary" && config.Kind != "mkdir" {
		log.Printf("Supported file kind is not supported %s (available kind options: json,string,text,binary,remove,mkdir)", command)
		return nil, false
	}
	if config.Result == "" && config.Kind != "remove" && config.Kind != "mkdir" {
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
	switch config.Kind {
	case "remove":
		k := dvdir.DeleteFilesIfExist(fileNames)
		SaveActionResult(config.Result, strconv.Itoa(k)+"/"+strconv.Itoa(n), ctx)
		return true
	case "mkdir":
		k := dvdir.MakeALlDirs(fileNames)
		SaveActionResult(config.Result, strconv.Itoa(k)+"/"+strconv.Itoa(n), ctx)
		return true
	}
	isMulti := n > 1
	if isMulti {
		if config.Append == 0 {
			DeleteActionResult(config.Result, ctx)
		}
		if config.Kind != "json" && config.Sort != nil {
			fileNames = dvtextutils.SortStringArray(fileNames, config.Sort)
		}
	} else {
		isMulti = config.Append != 0
	}
	for i := 0; i < n; i++ {
		ProcessSingleFile(fileNames[i], config, ctx, isMulti)
	}
	return true
}

func combineStrings(s string, isMulti bool, config *ReadFileConfig, ctx *dvcontext.RequestContext) string {
	if isMulti {
		prev, ok := ReadActionResult(config.Result, ctx)
		if ok && prev != nil {
			t := dvevaluation.AnyToString(prev)
			if config.Append < 0 {
				t, s = s, t
			}
			if config.Joiner != "" {
				t += config.Joiner
			}
			if config.EolJoiner {
				t += "\n"
			}
			s = t + s
		}
	}
	return s
}

func ProcessSingleFile(fileName string, config *ReadFileConfig, ctx *dvcontext.RequestContext, isMulti bool) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		log.Printf("File not found %s", config.FileName)
		return false
	}
	var dat []byte
	env := ctx.GetEnvironment()
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
		res = combineStrings(string(dat), isMulti, config, ctx)
	case "string":
		res = combineStrings(strings.TrimSpace(string(dat)), isMulti, config, ctx)
	case "binary":
		if isMulti {
			prev, ok := ReadActionResult(config.Result, ctx)
			if ok && prev != nil {
				t := dvevaluation.AnyToByteArray(prev)
				if config.Append < 0 {
					t, dat = dat, t
				}
				if config.Joiner != "" {
					t = append(t, []byte(config.Joiner)...)
				}
				if config.EolJoiner {
					t = append(t, byte(10))
				}
				res = append(t, dat...)
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
