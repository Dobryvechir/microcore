/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvssews

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"net/http"
	"strconv"
	"time"
)

const (
	REQUEST_SSE_COUNTER = "REQUEST_SSE_COUNTER"
	REQUEST_SSE_INDEX   = "REQUEST_SSE_INDEX"
	REQUEST_SSE_EVENT   = "REQUEST_SSE_EVENT"
)

var SseHeaders = map[string]string{
	"Access-Control-Allow-Origin": "*",
	"Content-Type":                "text/event-stream",
	"Cache-Control":               "no-cache",
	"Connection":                  "keep-alive",
	"Transfer-Encoding":           "chunked",
}

var parallelTimeUnitMs time.Duration = 1000 * time.Millisecond

func prepareSSEWSContext(ctx *dvcontext.RequestContext, parallelMode int) (ok bool, heartBeat int, totalDownCounter int, intervalDownCounter int) {
	if ctx == nil || ctx.Action == nil || ctx.Action.SseWs == nil || ctx.Writer == nil {
		return
	}
	var flusher http.Flusher
	flusher, ok = ctx.Writer.(http.Flusher)
	if !ok {
		dvlog.PrintlnError("Streaming unsupported")
		return
	}
	heartBeat = ctx.Action.SseWs.HeartBeat
	if heartBeat < 1 {
		heartBeat = 1
	}
	totalDownCounter = ctx.Action.SseWs.TimeOut
	if totalDownCounter < 1 {
		totalDownCounter = 180
	}
	intervalDownCounter = ctx.Action.SseWs.Interval
	if intervalDownCounter < 1 {
		intervalDownCounter = 5
	}
	ctx.ParallelExecution = &dvcontext.ParallelExecutionControl{
		Flusher: flusher,
	}
	switch parallelMode {
	case dvcontext.PARALLEL_MODE_SSE:
		for k, v := range SseHeaders {
			ctx.Writer.Header().Set(k, v)
		}
		ctx.Writer.WriteHeader(200)
	}
	requestStarter(ctx)
	return
}

func RunInSSEWSContext(ctx *dvcontext.RequestContext, parallelMode int) {
	ok, heartBeat, totalDownCounter, intervalDownCounter := prepareSSEWSContext(ctx, parallelMode)
	if !ok {
		return
	}
	tick := time.Tick(parallelTimeUnitMs)
	heart := heartBeat
	interval := intervalDownCounter
downCounter:
	for count := 0; count < totalDownCounter; count++ {
		select {
		case <-tick:
			heart--
			interval--
			if interval <= 0 {
				interval = intervalDownCounter
				heart = heartBeat
				if !requestRunner(ctx, count) {
					break downCounter
				}
			} else if heart <= 0 {
				heart = heartBeat
				SSESendHeartBeat(ctx)
			}
		}
	}
	requestEnder(ctx)
}

func requestRunner(ctx *dvcontext.RequestContext, count int) bool {
	if ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_COUNTER, count)
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_EVENT, "v"+strconv.FormatInt(ctx.Id, 16)+"-"+strconv.Itoa(count))
	}
	if ctx.ExecutorFn == nil {
		dvlog.PrintfError("executorFn not defined for %s", ctx.Url)
		return false
	}
	r := ctx.ExecutorFn(ctx, dvcontext.STAGE_MODE_MIDDLE)
	v := false
	switch r.(type) {
	case bool:
		if r.(bool) {
			v = true
		}
	default:
		dvlog.PrintfError("Incorrect executorFn (not bool) for %s", ctx.Url)
	}
	return v
}

func SetTimeUnitInSeconds(unitTime float32) {
	var v int64 = int64(unitTime * 1000)
	if v >= 10 {
		parallelTimeUnitMs = time.Duration(v) * time.Millisecond
	}
}

func SetParallelProcessingParameters(params *dvcontext.ParallelProcessing) {
	if params != nil {
		if params.IntervalTimeUnit > 0 {
			SetTimeUnitInSeconds(params.IntervalTimeUnit)
		}
	}
}

func requestEnder(ctx *dvcontext.RequestContext) {
	if ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_EVENT, "v"+strconv.FormatInt(ctx.Id, 16)+"-END")
	}
	ctx.ExecutorFn(ctx, dvcontext.STAGE_MODE_END)
}

func requestStarter(ctx *dvcontext.RequestContext) {
	if ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_EVENT, "v"+strconv.FormatInt(ctx.Id, 16)+"-START")
	}
	ctx.ExecutorFn(ctx, dvcontext.STAGE_MODE_START)
}
