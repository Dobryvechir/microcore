// package dvoc orchestrates actions, executions
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvaction

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
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
	return fireActionByName(ctx, ctx.Action.Name, ctx.Action.Definitions)
}

func fireActionByName(ctx *dvcontext.RequestContext, name string,
	definitions map[string]string) bool {
	prefix := ActionPrefix + name
	if ctx.PrimaryContextEnvironment.GetString(prefix+"_1") == "" {
		ctx.StatusCode = 501
		ctx.HandleCommunication()
		return true
	}
	res := ExecuteSequence(prefix, ctx, definitions)
	ActionProcessResult(ctx, res)
	return true
}

func ActionProcessResult(ctx *dvcontext.RequestContext, res bool) {
	if !res {
		ctx.HandleInternalServerError()
	} else {
		ActionContextResult(ctx)
	}
}

func fireStaticAction(ctx *dvcontext.RequestContext) bool {
	ActionProcessResult(ctx, true)
	return true
}

func ActionContextResult(ctx *dvcontext.RequestContext) {
	if ctx.StatusCode >= 400 {
		ctx.HandleCommunication()
		return
	}
	action := ctx.Action
	if action != nil && action.Result != "" {
		res, err := ctx.PrimaryContextEnvironment.CalculateString(action.Result)
		if err != nil {
			ctx.Error = err
			ctx.HandleInternalServerError()
			return
		}
		switch action.ResultMode {
		case "file":
			ctx.Output, err = GetContextFileResult(ctx, res)
		case "var":
			ctx.Output, err = GetContextVarResult(ctx, res)
		default:
			ctx.Output = []byte(res)
		}
		if err != nil {
			ctx.Error = err
			ctx.HandleInternalServerError()
			return
		}
		ctx.Output = []byte(res)
	}
	ctx.HandleCommunication()
}

func GetContextVarResult(ctx *dvcontext.RequestContext, varName string) ([]byte, error) {
	dat, ok := ctx.PrimaryContextEnvironment.Get(varName)
	if !ok {
		return nil, errors.New("Variable " + varName + " not set")
	}
	str := dvevaluation.AnyToString(dat)
	return []byte(str), nil
}

func GetContextFileResult(ctx *dvcontext.RequestContext, fileName string) ([]byte, error) {
	if fileName == "" {
		ctx.HandleFileNotFound()
		return nil, errors.New("Empty file name")
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Cannot read %s: %v", fileName, err)
		ctx.HandleInternalServerError()
		return nil, errors.New("File " + fileName + " not found")
	}
	if !bytes.Contains(data, []byte("{{")) {
		return data, nil
	}
	res, err := ctx.PrimaryContextEnvironment.CalculateString(string(data))
	return []byte(res), err
}

func fireSwitchAction(ctx *dvcontext.RequestContext) bool {
	action := ctx.Action
	actionName := action.Result
	conditions := action.Conditions
	if nil != conditions {
		for k, v := range conditions {
			res, err := ctx.PrimaryContextEnvironment.EvaluateBooleanExpression(k)
			if err != nil {
				log.Printf("Failed to evaluate %s: %v", k, err)
				ctx.HandleInternalServerError()
				return true
			}
			if res {
				actionName = v
				break
			}
		}
	}
	return fireActionByName(ctx, actionName, action.Definitions)
}

func securityEndPointHandler(ctx *dvcontext.RequestContext) bool {
	res := dvsecurity.LoginByRequestEndPointHandler(ctx)
	if res {
		ActionContextResult(ctx)
	}
	return res
}
