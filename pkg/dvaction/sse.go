/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvssews"
)

func fireSseAction(ctx *dvcontext.RequestContext) bool {
	action := ctx.Action
	if action == nil || action.SseWs == nil {
		dvlog.PrintfError("SseWs is not defined for %s", ctx.Url)
		return false
	}
	ctx.ExecutorFn = SseExecutor
	dvssews.RunInSSEWSContext(ctx, dvcontext.PARALLEL_MODE_SSE)
	return true
}

func SseExecutor(ctx *dvcontext.RequestContext, v interface{}) interface{} {
	mode := v.(int)
	change := ctx.Action.SseWs.Change
	checkChange := change != nil && change.ActionCheck != ""
	switch mode {
	case dvcontext.STAGE_MODE_START:
		if ctx.Action.SseWs.Start != nil {
			return SseServe(ctx, ctx.Action.SseWs.Start)
		} else if ctx.Action.SseWs.ServeMidAtStart && ctx.Action.SseWs.Mid != nil {
			return SseServe(ctx, ctx.Action.SseWs.Mid)
		}
		if checkChange {
			fireActionByName(ctx, change.ActionCheck, ctx.Action.Definitions, true)
			s, ok := ReadActionResult(change.ActionCheckVar, ctx)
			if !ok {
				s = nil
			}
			ctx.ParallelExecution.Value = s
			if s != nil && ctx.Action.SseWs.ServeMidAtStart && ctx.Action.SseWs.Mid == nil {
				if ctx.LogLevel >= dvlog.LogDetail {
					str := dvevaluation.AnyToString(s)
					dvlog.PrintfError("Start change %s ", str)
				}
				serveSseChange(ctx, change)
			} else {
				if ctx.LogLevel >= dvlog.LogDetail {
					str := dvevaluation.AnyToString(s)
					dvlog.PrintfError("Start (not sent) change %s ", str)
				}
			}
		}
	case dvcontext.STAGE_MODE_MIDDLE:
		if ctx.Action.SseWs.Mid != nil {
			return SseServe(ctx, ctx.Action.SseWs.Mid)
		}
		if checkChange {
			fireActionByName(ctx, change.ActionCheck, ctx.Action.Definitions, true)
			v, ok := ReadActionResult(ctx.Action.Result, ctx)
			if !ok {
				v = nil
			}
			changed := !deltaCompare(ctx, v, ctx.ParallelExecution.Value, change.Places)
			ctx.ParallelExecution.Value = v
			if changed {
				if ctx.LogLevel >= dvlog.LogDetail {
					str := dvevaluation.AnyToString(v)
					dvlog.PrintfError("Change %s ", str)
				}
				serveSseChange(ctx, change)
			} else if ctx.LogLevel >= dvlog.LogDebug {
				dvlog.PrintlnError("No change")
			}
		}
	case dvcontext.STAGE_MODE_END:
		if ctx.Action.SseWs.End != nil {
			SseServe(ctx, ctx.Action.SseWs.End)
		}
	}
	return true
}

func SseServe(ctx *dvcontext.RequestContext, control *dvcontext.Stage) bool {
	if control.Action != "" {
		fireActionByName(ctx, control.Action, ctx.Action.Definitions, true)
	}
	var err error
	res := true
	if control.Condition != "" {
		res, err = ctx.PrimaryContextEnvironment.EvaluateBooleanExpression(control.Condition)
		if err != nil {
			dvlog.PrintfError("Error in condition %s %v", control.Condition, err)
			return false
		}
	}
	if control.Result != "" {
		if control.Result == ":" {
			dvssews.SSESendHeartBeat(ctx)
		} else {
			s, err := ctx.LocalContextEnvironment.CalculateString(control.Result)
			if err != nil && s != "" {
				dvssews.SSEMessageInterface(ctx, s)
			} else {
				dvssews.SSESendHeartBeat(ctx)
				if err != nil && ctx.LogLevel >= dvlog.LogInfo {
					dvlog.PrintfError("Error in %s: %s", control.Result, err.Error())
				}
			}
		}
	}
	return res
}

func serveSseChange(ctx *dvcontext.RequestContext, change *dvcontext.SSEChange) {
	var res interface{} = ctx.ParallelExecution.Value
	ok := true
	if change.ActionFull != "" {
		fireActionByName(ctx, change.ActionFull, ctx.Action.Definitions, true)
		res, ok = ReadActionResult(change.ActionFullResult, ctx)
	} else if change.ActionFullResult != "" {
		res, ok = ReadActionResult(change.ActionFullResult, ctx)
	}
	if ok && res != nil {
		dvssews.SSEMessageInterface(ctx, res)
	} else {
		dvssews.SSESendHeartBeat(ctx)
	}
}

func deltaCompare(ctx *dvcontext.RequestContext, newVal interface{}, oldVal interface{}, places []string) bool {
	newStr := dvevaluation.AnyToString(newVal)
	oldStr := dvevaluation.AnyToString(oldVal)
	if newStr == oldStr {
		if ctx.LogLevel >= dvlog.LogTrace {
			dvlog.PrintfError("Completely unchanged %s", newStr)
		}
		return true
	}
	if ctx.LogLevel >= dvlog.LogTrace {
		dvlog.PrintfError("Comparing %s with %s", newStr, oldStr)
	}
	newDv := dvevaluation.AnyToDvVariable(newVal)
	oldDv := dvevaluation.AnyToDvVariable(oldVal)
	n := len(places)
	if n == 0 {
		return pointwiseDeltaCompare(ctx.PrimaryContextEnvironment, newDv, oldDv, "")
	}
	for i := 0; i < n; i++ {
		if !pointwiseDeltaCompare(ctx.PrimaryContextEnvironment, newDv, oldDv, places[i]) {
			return false
		}
	}
	return true
}

func pointwiseDeltaCompare(env *dvevaluation.DvObject, newV *dvevaluation.DvVariable, oldV *dvevaluation.DvVariable, path string) bool {
	if path != "" {
		res, _, err := dvjson.ReadPathOfAny(newV, path, false, env)
		if err != nil {
			newV = nil
		} else {
			newV = dvevaluation.AnyToDvVariable(res)
		}
		res, _, err = dvjson.ReadPathOfAny(oldV, path, false, env)
		if err != nil {
			oldV = nil
		} else {
			oldV = dvevaluation.AnyToDvVariable(res)
		}
	}
	r := newV.CompareWholeDvField(oldV)
	return r == 0
}
