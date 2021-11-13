// package dvoc orchestrates actions, executions
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvsecurity"
	"io/ioutil"
	"log"
)

const (
	ActionPrefix = "ACTION_"
)

var Log = dvlog.LogError

func fireAction(ctx *dvcontext.RequestContext) bool {
	action := ctx.Action
	prefix := ActionPrefix + action.Name
	if ctx.ExtraAsDvObject.GetString(prefix+"_1") == "" {
		ctx.StatusCode = 501
		ctx.HandleCommunication()
		return true
	}
	res := ExecuteSequence(prefix, ctx)
	if !res {
		ctx.HandleInternalServerError()
	} else {
		ActionContextResult(ctx)
	}
	return true
}

func fireStaticAction(ctx *dvcontext.RequestContext) bool {
	ActionContextResult(ctx)
	return true
}

func ActionContextResult(ctx *dvcontext.RequestContext) {
	if ctx.StatusCode >= 400 {
		ctx.HandleCommunication()
		return
	}
	action := ctx.Action
	if action != nil && action.Result != "" {
		res, err := ctx.ExtraAsDvObject.CalculateString(action.Result)
		if err != nil {
			ctx.Error = err
			ctx.HandleInternalServerError()
			return
		}
		ctx.Output = []byte(res)
	}
	ctx.HandleCommunication()
}

func fireFileAction(ctx *dvcontext.RequestContext) bool {
	action := ctx.Action
	fileName := action.Result
	conditions := action.Conditions
	if conditions != nil {
		for k, v := range conditions {
			res, err := ctx.ExtraAsDvObject.EvaluateBooleanExpression(k)
			if err != nil {
				log.Printf("Failed to evaluate %s: %v", k, err)
				ctx.HandleInternalServerError()
				return true
			}
			if res {
				fileName = v
			}
		}
	}
	if fileName == "" {
		ctx.HandleFileNotFound()
		return true
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Cannot read %s: %v", fileName, err)
		ctx.HandleInternalServerError()
		return true
	}
	ctx.Output = data
	ctx.HandleCommunication()
	return true
}

func securityEndPointHandler(ctx *dvcontext.RequestContext) bool {
	res := dvsecurity.LoginByRequestEndPointHandler(ctx)
	if res {
		ActionContextResult(ctx)
	}
	return res
}
