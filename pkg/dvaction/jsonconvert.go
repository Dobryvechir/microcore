/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

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
	Source      *JsonRead           `json:"source"`
	Result      string              `json:"result"`
	StorePrefix string              `json:"prefix"`
	Add       []JsonConvertModify `json:"add"`
	Remove     []JsonConvertModify `json:"remove"`
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
	if config.Source == nil || config.Source.Place == "" {
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
