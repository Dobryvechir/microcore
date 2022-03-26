/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdynamic

import (
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

type DynamicServerConfig struct {
	Url string `json:"url"`
	Log int    `json:"log"`
}

func DynamicServerInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &DynamicServerConfig{}
	if !dvaction.DefaultInitWithObject(command, config, dvaction.GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func DynamicServerRun(data []interface{}) bool {
	config := data[0].(*DynamicServerConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return DynamicServerByConfig(config, ctx)
}

func DynamicServerByConfig(config *DynamicServerConfig, ctx *dvcontext.RequestContext) bool {
	ctx.Server.ProxyServerUrl = config.Url
	if config.Url != "" {
		ctx.Server.ProxyServerHttp = true
		if config.Log > 0 {
			dvcom.LogServer = true
			dvlog.CurrentLogLevel = config.Log
		}
	} else {
		ctx.Server.ProxyServerHttp = false
		dvcom.LogServer = false
	}
	return true
}
