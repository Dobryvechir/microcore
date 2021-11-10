/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"log"
)

type JsonConvertModify struct {
	Source      string `json:"src"`
	Path        string `json:"path"`
	Destination string `json:"dst"`
	Expression  string `json:"expr"`
}

type JsonConvertConfig struct {
	Source      string              `json:"source"`
	Path        string              `json:"path"`
	Result      string              `json:"result"`
	StorePrefix string              `json:"prefix"`
	SortField   string              `json:"sort"`
	Added       []JsonConvertModify `json:"added"`
	Removed     []JsonConvertModify `json:"removed"`
}

func jsonConvertInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &JsonConvertConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	if config.StorePrefix == "" {
		log.Printf("prefix must be specified in %s", command)
		return nil, false
	}
	if config.Source == "" {
		log.Printf("source must be present in %s", command)
		return nil, false
	}
	if config.Result == "" {
		log.Printf("Result name is not specified in command %s", command)
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

func JsonConvertRunByConfig(config *JsonConvertConfig, ctx *dvcontext.RequestContext) bool {

	return true
}
