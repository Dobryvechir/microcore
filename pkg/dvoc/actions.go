/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import "github.com/Dobryvechir/microcore/pkg/dvaction"

const (
	CommandCopyToPod              = "copyToPod"
	CommandCopyFromPod            = "copyFromPod"
	CommandEnv                    = "env"
	CommandExpose                 = "expose"
	CommandMicroServiceCacheClean = "microserviceCacheClean"
	CommandMicroServiceDown       = "microserviceDown"
	CommandMicroServiceExec       = "microserviceExec"
	CommandMicroServiceSave       = "microserviceSave"
	CommandMicroServiceTemplate   = "microserviceTemplate"
	CommandMicroServiceRestore    = "microserviceRestore"
	CommandMicroServiceUp         = "microserviceUp"
	CommandMicroServiceUpOnly     = "microserviceUpOnly"
)

var processFunctions = map[string]dvaction.ProcessFunction{
	CommandCopyToPod:              {Init: CopyToPodInit, Run: CopyToPodRun},
	CommandCopyFromPod:            {Init: CopyFromPodInit, Run: CopyFromPodRun},
	CommandEnv:                    {Init: ProcessEnvSettingInit, Run: ProcessEnvSettingsRun},
	CommandExpose:                 {Init: ExposeMicroServiceInit, Run: ExposeMicroServiceRun},
	CommandMicroServiceCacheClean: {Init: MicroServiceCacheCleanInit, Run: MicroServiceCacheCleanRun},
	CommandMicroServiceExec:       {Init: MicroServiceExecInit, Run: MicroServiceExecRun},
	CommandMicroServiceDown:       {Init: MicroServiceDownInit, Run: MicroServiceDownRun},
	CommandMicroServiceSave:       {Init: MicroServiceSaveInit, Run: MicroServiceSaveRun},
	CommandMicroServiceTemplate:   {Init: MicroServiceTemplateInit, Run: MicroServiceTemplateRun},
	CommandMicroServiceRestore:    {Init: MicroServiceRestoreInit, Run: MicroServiceRestoreRun},
	CommandMicroServiceUp:         {Init: MicroServiceUpInit, Run: MicroServiceUpRun},
	CommandMicroServiceUpOnly:     {Init: MicroServiceUpInit, Run: MicroServiceUpOnlyRun},
}

const (
	openShiftEnsureRoutes = "OPENSHIFT_ENSURE_ROUTES"
)

func actionInit() bool {
	dvaction.AddProcessFunctions(processFunctions)
	dvaction.RegisterCredentialProvider(GetOpenshiftSecrets)
	return true
}

var actionInited = actionInit()
