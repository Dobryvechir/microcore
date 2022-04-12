/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvssews

import (
	"bytes"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strings"
)

var (
	SSELineEnd      = []byte("\n")
	SSELineFullEnd  = []byte("\n\n")
	SSEIdPrefix     = []byte("id: ")
	SSEHeartBeat    = []byte(": heartbeat\n\n")
	SSEEventPrefix  = []byte("event: ")
	SSEEventDefault = []byte("message")
	SSEDataPrefix   = []byte("data: ")
)

func SSEMessageInterface(ctx *dvcontext.RequestContext, v interface{}) {
	s := dvevaluation.AnyToString(v)
	s = strings.TrimSpace(s)
	if s != "" {
		SSEMessageString(ctx, s)
	} else {
		SSESendHeartBeat(ctx)
	}
}

func SSEMessageString(ctx *dvcontext.RequestContext, s string) {
	SSEMessageBytes(ctx, nil, nil, []byte(s))
}

func cleanEOF(b []byte) {
	n := len(b)
	for i := 0; i < n; i++ {
		if b[i] < ' ' {
			b[i] = ' '
		}
	}
}

func SSEMessageBytes(ctx *dvcontext.RequestContext, id []byte, event []byte, data []byte) {
	buf := bytes.NewBuffer(make([]byte, 256+len(data)))
	if event == nil {
		event = SSEEventDefault
	}
	if id == nil {
		id = []byte(ctx.PrimaryContextEnvironment.GetString(REQUEST_SSE_EVENT))
	}
	cleanEOF(id)
	cleanEOF(event)
	cleanEOF(data)
	buf.Write(SSEIdPrefix)
	buf.Write(id)
	buf.Write(SSELineEnd)
	buf.Write(SSEEventPrefix)
	buf.Write(event)
	buf.Write(SSELineEnd)
	buf.Write(SSEDataPrefix)
	buf.Write(data)
	buf.Write(SSELineFullEnd)
	ctx.Writer.Write(buf.Bytes())
	ctx.ParallelExecution.Flusher.Flush()
}

func SSESendPortion(ctx *dvcontext.RequestContext, data []byte) {
	ctx.Writer.Write(data)
	ctx.ParallelExecution.Flusher.Flush()
}

func SSESendHeartBeat(ctx *dvcontext.RequestContext) {
	SSESendPortion(ctx, SSEHeartBeat)
}
