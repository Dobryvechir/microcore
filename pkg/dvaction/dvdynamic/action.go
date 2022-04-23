/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdynamic

import "github.com/Dobryvechir/microcore/pkg/dvaction"

const (
	CommandDebug         = "debug"
	CommandDynamic       = "dynamic"
	CommandDynamicServer = "dserver"
	CommandDynamicCache  = "dcache"
	CommandKeyVariables  = "key_variables"
)

var processFunctions = map[string]dvaction.ProcessFunction{
	CommandDebug:         {Init: DebugInit, Run: DebugRun},
	CommandDynamic:       {Init: DynamicActionInit, Run: DynamicActionRun},
	CommandDynamicServer: {Init: DynamicServerInit, Run: DynamicServerRun},
	CommandDynamicCache:  {Init: DynamicCacheInit, Run: DynamicCacheRun},
	CommandKeyVariables:  {Init: KeyVariablesInit, Run: KeyVariablesRun},
}

func Init() bool {
	dvaction.AddProcessFunctions(processFunctions)
	return true
}

var inited = Init()
