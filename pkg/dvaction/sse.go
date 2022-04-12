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
	dvssews.PushSSEWSContext(ctx, dvcontext.PARALLEL_MODE_SSE)
	return true
}

func SseExecutor(ctx *dvcontext.RequestContext, v interface{}) interface{} {
	mode := v.(int)
	delta := ctx.Action.SseWs.Delta
	isDelta := delta != nil && delta.ActionDelta != ""
	switch mode {
	case dvcontext.STAGE_MODE_START:
		if ctx.Action.SseWs.Start != nil {
			return SseServe(ctx, ctx.Action.SseWs.Start)
		} else if ctx.Action.SseWs.ServeMidAtStart && ctx.Action.SseWs.Mid != nil {
			return SseServe(ctx, ctx.Action.SseWs.Mid)
		}
		if isDelta {
			fireActionByName(ctx, delta.ActionDelta, ctx.Action.Definitions, true)
			s, ok := ReadActionResult(ctx.Action.Result, ctx)
			if !ok {
				s = nil
			}
			ctx.ParallelExecution.Value = s
			if s != nil && ctx.Action.SseWs.ServeMidAtStart && ctx.Action.SseWs.Mid == nil {
				serveDelta(ctx, delta)
			}
		}
	case dvcontext.STAGE_MODE_MIDDLE:
		if ctx.Action.SseWs.Mid != nil {
			return SseServe(ctx, ctx.Action.SseWs.Mid)
		}
		if isDelta {
			fireActionByName(ctx, delta.ActionDelta, ctx.Action.Definitions, true)
			v, ok := ReadActionResult(ctx.Action.Result, ctx)
			if !ok {
				v = nil
			}
			changed := deltaCompare(ctx, v, ctx.ParallelExecution.Value, delta.Places)
			ctx.ParallelExecution.Value = v
			if changed {
				serveDelta(ctx, delta)
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
			s, ok := ReadActionResult(control.Result, ctx)
			if ok && s != nil {
				dvssews.SSEMessageInterface(ctx, s)
			} else {
				dvssews.SSESendHeartBeat(ctx)
			}
		}
	}
	return res
}

func serveDelta(ctx *dvcontext.RequestContext, delta *dvcontext.SSEDelta) {
	var res interface{} = ctx.ParallelExecution.Value
	ok := true
	if delta.ActionFull != "" {
		fireActionByName(ctx, delta.ActionFull, ctx.Action.Definitions, true)
		res, ok = ReadActionResult(delta.ActionFullResult, ctx)
	}
	if ok && res != nil {
		dvssews.SSEMessageInterface(ctx, res)
	} else {
		dvssews.SSESendHeartBeat(ctx)
	}
}

func deltaCompare(ctx *dvcontext.RequestContext, newVal interface{}, oldVal interface{}, places []string) bool {
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
