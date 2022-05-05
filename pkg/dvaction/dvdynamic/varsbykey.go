/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdynamic

import (
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

type KeyVariablesConfig struct {
	Prefix     string   `json:"prefix"`
	Values     []string `json:"values"`
	SuccessVar string   `json:"success_var"`
	KeyParam   string   `json:"key_param"`
}

func KeyVariablesInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &KeyVariablesConfig{}
	if !dvaction.DefaultInitWithObject(command, config, dvaction.GetEnvironment(ctx)) {
		return nil, false
	}
	if config.KeyParam == "" {
		config.KeyParam = "URL_PARAM_KEY"
	}
	if config.SuccessVar == "" {
		config.SuccessVar = "RESULT"
	}
	return []interface{}{config, ctx}, true
}

func KeyVariablesRun(data []interface{}) bool {
	config := data[0].(*KeyVariablesConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	KeyVariablesRunByConfig(config, ctx)
	return true
}

func KeyVariablesRunByConfig(config *KeyVariablesConfig, ctx *dvcontext.RequestContext) {
	dvaction.SaveActionResult(config.SuccessVar, false, ctx)
	v := ctx.LocalContextEnvironment.GetString(config.KeyParam)
	s := ctx.LocalContextEnvironment.GetString(config.Prefix + v)
	if v == "" || s == "" {
		return
	}
	data := dvtextutils.SmartReadStringList(s, false)
	m := len(data)
	n := len(config.Values)
	for i := 0; i < n; i++ {
		r := ""
		if i < m {
			r = data[i]
		}
		k := config.Values[i]
		dvaction.SaveActionResult(k, r, ctx)
	}
	if n > 0 && m > 0 {
		dvaction.SaveActionResult(config.SuccessVar, true, ctx)
	}
}
