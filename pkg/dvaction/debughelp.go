/***********************************************************************
MicroCore
Copyright 2022 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func DebugShowVariablesByEnvList(ctx *dvcontext.RequestContext, lstName string) {
	lst, ok := ctx.PrimaryContextEnvironment.Get(lstName)
	if !ok {
		return
	}
	switch lst.(type) {
	case string:
		DebugShowVariablesByList(ctx, lst.(string))
	}
}

func DebugShowVariablesByList(ctx *dvcontext.RequestContext, lst string) {
	list := dvtextutils.ConvertToNonEmptyList(lst)
	n := len(list)
	for i := 0; i < n; i++ {
		DebugShowVariable(ctx, list[i])
	}
}

func DebugShowVariable(ctx *dvcontext.RequestContext, name string) {
	v, ok := ReadActionResult(name, ctx)
	if !ok {
		dvlog.Printf("Absent %s", "Absent %s", name)
		return
	}
	s := dvevaluation.AnyToString(v)
	dvlog.PrintfError("%s=%s", name, s)
}
