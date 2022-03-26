/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdynamic

import (
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type DynamicCacheConfig struct {
	SessionId     string                          `json:"session_id"`
	SessionParams *dvcontext.SessionActionRequest `json:"params"`
	Write         map[string]string               `json:"write"`
	Read          []string                        `json:"read"`
}

func DynamicCacheInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &DynamicCacheConfig{}
	if !dvaction.DefaultInitWithObject(command, config, dvaction.GetEnvironment(ctx)) {
		return nil, false
	}
	if config.SessionId == "" {
		log.Printf("sessionId must be specified in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func DynamicCacheRun(data []interface{}) bool {
	config := data[0].(*DynamicCacheConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return DynamicCacheByConfig(config, ctx)
}

func DynamicCacheByConfig(config *DynamicCacheConfig, ctx *dvcontext.RequestContext) bool {
	if config == nil {
		return true
	}
	request, err := ctx.Server.Session.GetSessionStorage(ctx, config.SessionParams, config.SessionId)
	if err != nil {
		dvlog.PrintfError("Error in session request %v", err)
		return true
	}
	if config.Write != nil {
		for k, v := range config.Write {
			request.SetItem(k, v)
		}
	}
	res := &dvevaluation.DvVariable{
		Kind:   dvevaluation.FIELD_OBJECT,
		Fields: make([]*dvevaluation.DvVariable, 0, 16),
	}
	if config.Read != nil {
		n := len(config.Read)
		for i := 0; i < n; i++ {
			k := config.Read[i]
			v := request.GetItem(k)
			d := dvevaluation.AnyToDvVariable(v)
			d.Name = []byte(k)
			res.Fields = append(res.Fields, d)
		}
	}
	dvaction.SaveActionResult("request:RESULT", res, ctx)
	return true
}
