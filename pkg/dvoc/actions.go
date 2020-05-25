/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import "github.com/Dobryvechir/microcore/pkg/dvcontext"

const (
	ActionPrefix = "ACTION_"
)

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
