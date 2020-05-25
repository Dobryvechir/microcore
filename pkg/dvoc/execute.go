/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ProcessFunction struct {
	Init  func(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool)
	Run   func([]interface{}) bool
	Async bool
}

const (
	CommandHttp                   = "http"
	CommandEnv                    = "env"
	CommandCopyToPod              = "copyToPod"
	CommandCopyFromPod            = "copyFromPod"
	CommandMicroServiceExec       = "microserviceExec"
	CommandMicroServiceUp         = "microserviceUp"
	CommandMicroServiceUpOnly     = "microserviceUpOnly"
	CommandMicroServiceDown       = "microserviceDown"
	CommandMicroServiceSave       = "microserviceSave"
	CommandMicroServiceRestore    = "microserviceRestore"
	CommandMicroServiceCacheClean = "microserviceCacheClean"
	CommandExpose                 = "expose"
	CommandNet                    = "net"
	CommandOs                     = "os"
	CommandPortForward            = "forward"
	CommandSql                    = "sql"
)

var processFunctions = map[string]ProcessFunction{
	CommandHttp:                   {Init: processNetInit, Run: processNetRun},
	CommandEnv:                    {Init: processEnvSettingInit, Run: processEnvSettingsRun},
	CommandCopyToPod:              {Init: copyToPodInit, Run: copyToPodRun},
	CommandCopyFromPod:            {Init: copyFromPodInit, Run: copyFromPodRun},
	CommandMicroServiceExec:       {Init: microServiceExecInit, Run: microServiceExecRun},
	CommandMicroServiceUp:         {Init: microServiceUpInit, Run: microServiceUpRun},
	CommandMicroServiceUpOnly:     {Init: microServiceUpInit, Run: microServiceUpOnlyRun},
	CommandMicroServiceDown:       {Init: microServiceDownInit, Run: microServiceDownRun},
	CommandMicroServiceSave:       {Init: microServiceSaveInit, Run: microServiceSaveRun},
	CommandMicroServiceRestore:    {Init: microServiceRestoreInit, Run: microServiceRestoreRun},
	CommandMicroServiceCacheClean: {Init: microServiceCacheCleanInit, Run: microServiceCacheCleanRun},
	CommandExpose:                 {Init: exposeMicroServiceInit, Run: exposeMicroServiceRun},
	CommandNet:                    {Init: SmartNetInit, Run: SmartNetRun},
	CommandOs:                     {Init: processOsInit, Run: processOsRun},
	CommandPortForward:            {Init: portForwardInit, Run: portForwardRun, Async: true},
	CommandSql:                    {Init: dvdbdata.SqlInit, Run: dvdbdata.SqlRun},
}

func AddProcessFunction(key string, processor ProcessFunction) {
	processFunctions[key] = processor
}

func getWaitKeys() string {
	res := ""
	for k := range processFunctions {
		if res == "" {
			res = k + ":"
		} else {
			res += " or " + k + ":"
		}
	}
	return res
}

func ExecuteProcessFunction(fn *ProcessFunction, pauseTime int, totalTime int, command string, group *sync.WaitGroup, ctx *dvcontext.RequestContext) bool {
	data, ok := fn.Init(command, ctx)
	if !ok {
		return false
	}
	if fn.Async {
		data = append(data, group)
	}
	for ; totalTime >= 0; totalTime -= pauseTime {
		if fn.Async && group != nil {
			group.Add(1)
		}
		if fn.Run(data) {
			return true
		}
		if pauseTime <= 0 {
			break
		}
		if totalTime > 0 {
			if Log >= LogInfo {
				dvlog.PrintfError("Waiting for %d / %d seconds", pauseTime, totalTime)
			}
			time.Sleep(time.Duration(pauseTime) * time.Second)
		}
	}
	dvlog.PrintfError("%s command finally failed", command)
	return false

}

func ExecuteSingleCommand(pauseTime int, totalTime int, prefix string, command string) bool {
	if prefix != "" {
		command = prefix + ":" + command
	}
	pos := strings.Index(command, ":")
	if pos < 0 {
		dvlog.PrintfError("At first, place the command with colon before %s", command)
	}
	kind := command[:pos]
	waiter, ok := processFunctions[kind]
	if !ok {
		dvlog.PrintfError("in %s the first parameter must start with either %s, not with %s", command, getWaitKeys(), kind)
		return false
	}
	return ExecuteProcessFunction(&waiter, pauseTime, totalTime, command, nil, nil)
}

func ExecuteSequence(prefix string, ctx *dvcontext.RequestContext) bool {
	var wg sync.WaitGroup
	var waitCommand string
	var err error
	for n := 1; n < 1000000; n++ {
		p := prefix + "_" + strconv.Itoa(n)
		waitCommandRaw := strings.TrimSpace(dvparser.GlobalProperties[p])
		if waitCommandRaw == "" {
			wg.Wait()
			return true
		}
		if ctx == nil {
			waitCommand, err = dvparser.ConvertByteArrayByGlobalProperties([]byte(waitCommandRaw), waitCommandRaw)
		} else {
			waitCommand, err = ctx.ExtraAsDvObject.CalculateString(waitCommandRaw)
		}
		if err != nil {
			dvlog.PrintfError("Make sure you specified all constants in %s .properties file: %v", waitCommandRaw, err)
			return false
		}
		k := 0
		strict := true
		parallel := false
	LookAtControlChars:
		for ; k < len(waitCommand); k++ {
			c := waitCommand[k]
			switch c {
			case '?':
				strict = false
			case '!':
				parallel = true
			default:
				if c > ' ' {
					break LookAtControlChars
				}
			}
		}
		if k > 0 {
			waitCommand = waitCommand[k:]
		}
		if waitCommand != "" && !(waitCommand[0] >= '0' && waitCommand[0] <= '9') {
			waitCommand = "0,0,0," + waitCommand
		}
		pos := strings.Index(waitCommand, ",")
		if pos <= 0 {
			dvlog.PrintfError("in %s the first parameter must be the first idle time in seconds followed by comma", p)
			return false
		}
		idleTime, err := strconv.Atoi(waitCommand[:pos])
		if err != nil {
			dvlog.PrintfError("in %s the first parameter must be the first idle time (integer) in seconds followed by comma", p)
			return false
		}
		waitCommand = strings.TrimSpace(waitCommand[pos+1:])
		pos = strings.Index(waitCommand, ",")
		if pos <= 0 {
			dvlog.PrintfError("in %s the second parameter must be the pause time in seconds followed by comma", p)
			return false
		}
		pauseTime, err := strconv.Atoi(waitCommand[:pos])
		if err != nil {
			dvlog.PrintfError("in %s the second parameter must be the pause time (integer) in seconds followed by comma", p)
			return false
		}
		waitCommand = strings.TrimSpace(waitCommand[pos+1:])
		pos = strings.Index(waitCommand, ",")
		if pos <= 0 {
			dvlog.PrintfError("in %s the third parameter must be the total wait time in seconds followed by comma", p)
			return false
		}
		totalTime, err := strconv.Atoi(waitCommand[:pos])
		if err != nil {
			dvlog.PrintfError("in %s the third parameter must be the total wait time (integer) in seconds followed by comma", p)
			return false
		}
		waitCommand = strings.TrimSpace(waitCommand[pos+1:])
		kindPos := strings.Index(waitCommand, ":")
		if kindPos <= 0 {
			dvlog.PrintfError("in %s the fourth parameter must start with either http: or other command: (: is mandatory)", p)
			return false

		}
		kind := waitCommand[:kindPos]
		waiter, ok := processFunctions[kind]
		if !ok {
			dvlog.PrintfError("in %s the third parameter must start with either %s, not with %s", p, getWaitKeys(), kind)
			return false
		}
		if waiter.Async {
			parallel = false
		}
		if idleTime > 0 {
			if Log >= LogInfo {
				dvlog.PrintfError("idle waiting for %d seconds before %s", idleTime, waitCommand)
			}
			time.Sleep(time.Duration(idleTime) * time.Second)
		}
		if Log >= LogInfo {
			dvlog.PrintfError("starting waiting %d seconds for %s", totalTime, waitCommand)
		}
		if parallel {
			wg.Add(1)
			ok = true
			go func() {
				defer wg.Done()
				ExecuteProcessFunction(&waiter, pauseTime, totalTime, waitCommand, &wg, ctx)
			}()
		} else {
			ok = ExecuteProcessFunction(&waiter, pauseTime, totalTime, waitCommand, &wg, ctx)
		}
		if !ok && strict {
			return false
		}
	}
	wg.Wait()
	return true
}

func ocExecutorStartByEvent(eventName string, data []interface{}) error {
	prefix := "EXECUTE_" + strings.ToUpper(eventName)
	if !ExecuteSequence(prefix, nil) {
		return errors.New("Failed to execute " + prefix)
	}
	return nil
}

var ocExecutorRegistrationConfig = &dvmodules.HookRegistrationConfig{
	Name: "ocExecutor",
	HookEventMapper: map[string]dvmodules.HookMethodEndPointHandler{
		dvmodules.HookStartEvent: ocExecutorStartByEvent,
	},
	//GlobalInitHandler: MethodGlobalInitHandler
	//ServerInitHandler: MethodServerInitHandler
}

func RegisterOC() bool {
	dvmodules.RegisterActionProcessor("", fireAction, false)
	return dvmodules.SubscribeForEvents(ocExecutorRegistrationConfig, false)
}

var registeredOC = RegisterOC()
