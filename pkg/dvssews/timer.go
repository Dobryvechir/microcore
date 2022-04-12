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
	"sync"
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

var pool = make([]*dvcontext.RequestContext, 0, 1024)
var cycling = false
var poolMutex sync.Mutex
var parallelTimeUnitMs time.Duration = 1000 * time.Millisecond

func PushSSEWSContext(ctx *dvcontext.RequestContext, parallelMode int) {
	if ctx == nil || ctx.Action == nil || ctx.Action.SseWs == nil || ctx.Writer == nil {
		return
	}
	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		dvlog.PrintlnError("Streaming unsupported")
		return
	}
	if ctx.Action.SseWs.HeartBeat<1 {
		ctx.Action.SseWs.HeartBeat = 1
	}
	ctx.ParallelExecution = &dvcontext.ParallelExecutionControl{
		HeartBitDownCounter: ctx.Action.SseWs.HeartBeat,
		IntervalDownCounter: ctx.Action.SseWs.Interval,
		TotalDownCounter:    ctx.Action.SseWs.TimeOut,
		Flusher:             flusher,
	}
	poolMutex.Lock()
	pool = append(pool, ctx)
	n := len(pool)
	cycling = n > 0
	poolMutex.Unlock()
	switch parallelMode {
	case dvcontext.PARALLEL_MODE_SSE:
		for k, v := range SseHeaders {
			ctx.Writer.Header().Set(k, v)
		}
		ctx.Writer.WriteHeader(200)
	}
	go requestStarter(ctx)
	if n == 1 {
		startChannel()
	}
}

func removeSSEWSContext(index int) {
	n := len(pool)
	if index >= n || index < 0 {
		return
	}
	if pool[index] != nil && pool[index].ExecutorFn != nil {
		go requestEnder(pool[index])
	}
	if index == n-1 {
		pool = pool[:index]
	} else {
		pool = append(pool[:index], pool[index+1:]...)
	}
	cycling = len(pool) > 0
}

func PokeSSEWSContext(ctx *dvcontext.RequestContext) {
	n := len(pool)
	for i := 0; i < n; i++ {
		if pool[i] == ctx {
			poolMutex.Lock()
			removeSSEWSContext(i)
			poolMutex.Unlock()
			break
		}
	}
}

func startChannel() {
	go func() {
		count := 0
		tick := time.Tick(parallelTimeUnitMs)
		for cycling {
			select {
			case <-tick:
				count++
				poolMutex.Lock()
				n := len(pool)
				for i := 0; i < n; i++ {
					ctx := pool[i]
					ctx.ParallelExecution.TotalDownCounter--
					if ctx.ParallelExecution.TotalDownCounter <= 0 {
						removeSSEWSContext(i)
						i--
						n--
						continue
					}
					if !ctx.ParallelExecution.IsBusy {
						startable := ctx.ParallelExecution.IntervalDownCounter == 0
						if ctx.ParallelExecution.IntervalDownCounter > 0 {
							ctx.ParallelExecution.IntervalDownCounter--
							startable = ctx.ParallelExecution.IntervalDownCounter == 0
						}
						ctx.ParallelExecution.HeartBitDownCounter--
						if startable {
							ctx.ParallelExecution.IntervalDownCounter = ctx.Action.SseWs.Interval
							ctx.ParallelExecution.HeartBitDownCounter = ctx.Action.SseWs.HeartBeat
							ctx.ParallelExecution.IsBusy = true
							go requestRunner(ctx, count, i)
						} else if ctx.ParallelExecution.HeartBitDownCounter<=0 {
							ctx.ParallelExecution.HeartBitDownCounter = ctx.Action.SseWs.HeartBeat
							SSESendHeartBeat(ctx)
						}
					}
				}
				poolMutex.Unlock()
			}
		}
	}()
}

func requestRunner(ctx *dvcontext.RequestContext, count int, index int) {
	if ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_COUNTER, count)
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_INDEX, index)
		ctx.PrimaryContextEnvironment.Set(REQUEST_SSE_EVENT, "v"+strconv.FormatInt(ctx.Id, 16)+"-"+strconv.Itoa(ctx.ParallelExecution.TotalDownCounter))
	}
	if ctx.ExecutorFn == nil {
		dvlog.PrintfError("executorFn not defined for %s", ctx.Url)
		ctx.ParallelExecution.TotalDownCounter = 1
		return
	}
	r := ctx.ExecutorFn(ctx, dvcontext.STAGE_MODE_MIDDLE)
	switch r.(type) {
	case bool:
		if !r.(bool) {
			ctx.ParallelExecution.TotalDownCounter = 1
		}
	default:
		dvlog.PrintfError("Incorrect executorFn (not bool) for %s", ctx.Url)
		ctx.ParallelExecution.TotalDownCounter = 1
	}
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
