/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"log"
)

type PagingConfig struct {
	Source      string `json:"source"`
	Path		string `json:"path"`
	Result      string `json:"result"`
	StorePrefix string `json:"prefix"`
	SortField   string `json:"sort"`
	PageSize    int    `json:"pageSize"`
	CurrentPage int    `json:"currentPage"`
}

func pagingInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &PagingConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.StorePrefix == "" {
		log.Printf("prefix must be specified in %s", command)
		return nil, false
	}
	if config.PageSize <= 0 {
		log.Printf("pageSize must be present and positive in %s", command)
		return nil, false
	}
	if config.Result == "" {
		log.Printf("Result name is not specified in command %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func pagingRun(data []interface{}) bool {
	config := data[0].(*PagingConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return PagingRunByConfig(config, ctx)
}

func PagingRunByConfig(config *PagingConfig, ctx *dvcontext.RequestContext) bool {

	return true
}
