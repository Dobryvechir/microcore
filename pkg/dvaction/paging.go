/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type PagingConfig struct {
	Source      *JsonRead `json:"source"`
	Result      string    `json:"result"`
	StorePrefix string    `json:"prefix"`
	PageSize    int       `json:"pageSize"`
	CurrentPage int       `json:"currentPage"`
}

func pagingInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &PagingConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Source == nil || config.Source.Var == "" {
		log.Printf("source.var must be specified in %s", command)
		return nil, false
	}
	if config.PageSize <= 0 {
		config.PageSize = 3
	}
	if config.Result == "" {
		log.Printf("Result name is not specified in command %s", command)
		return nil, false
	}
	if config.StorePrefix == "" {
		config.StorePrefix = "PAGING"
	}
	if config.CurrentPage <= 0 {
		config.CurrentPage = 1
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
	src, err := JsonExtract(config.Source, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Source.Var)
		return true
	}
	itemAmount := dvevaluation.GetLengthOfAny(src)
	pageSize := config.PageSize
	pages := (itemAmount + pageSize) / pageSize
	currentPage := config.CurrentPage
	if currentPage > pages {
		currentPage = 1
	}
	startIndex := (currentPage - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > itemAmount {
		endIndex = itemAmount
	}
	var items interface{} = nil
	size := endIndex - startIndex
	if size > 0 {
		items = dvevaluation.GetChildrenOfAnyByRange(src, startIndex, endIndex)
	}
	pref := config.StorePrefix + "_"
	env := GetEnvironment(ctx)
	env.Set(pref+"PAGE_SIZE", pageSize)
	env.Set(pref+"OFFSET", startIndex)
	env.Set(pref+"PAGE", currentPage)
	env.Set(pref+"TOTAL_SIZE", itemAmount)
	env.Set(pref+"TOTAL_PAGES", pages)
	env.Set(config.Result, items)
	return true
}
