/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvaction

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdbdata"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
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
	CommandCall        = "call"
	CommandCompareJson = "compare"
	CommandConvert     = "convert"
	CommandFile        = "file"
	CommandFor         = "for"
	CommandHttp        = "http"
	CommandIf          = "if"
	CommandIfArray     = "ifarray"
	CommandIfEmpty     = "ifempty"
	CommandJs          = "js"
	CommandNet         = "net"
	CommandOs          = "os"
	CommandPaging      = "paging"
	CommandPortForward = "forward"
	CommandProxy       = "proxy"
	CommandRange       = "range"
	CommandReturn      = "return"
	CommandStore       = "store"
	CommandSql         = "sql"
	CommandSwitch      = "switch"
	CommandUpsert      = "upsert"
	CommandValidate    = "validate"
	CommandVar         = "var"
	CommandVersion     = "version"
	CommandVoid        = "void"
)

var processFunctions = map[string]ProcessFunction{
	CommandHttp:        {Init: processNetInit, Run: processNetRun},
	CommandOs:          {Init: processOsInit, Run: processOsRun},
	CommandPortForward: {Init: portForwardInit, Run: portForwardRun, Async: true},
	CommandNet:         {Init: SmartNetInit, Run: SmartNetRun},
	CommandProxy:       {Init: ProxyNetInit, Run: ProxyNetRun},
	CommandFile:        {Init: readFileActionInit, Run: readFileActionRun},
	CommandJs:          {Init: jsInit, Run: jsRun},
	CommandConvert:     {Init: jsonConvertInit, Run: jsonConvertRun},
	CommandCompareJson: {Init: compareJsonInit, Run: compareJsonRun},
	CommandVar:         {Init: varTransformInit, Run: varTransformRun},
	CommandSql:         {Init: dvdbdata.SqlInit, Run: dvdbdata.SqlRun},
	CommandStore:       {Init: storeInit, Run: storeRun},
	CommandValidate:    {Init: validationInit, Run: validationRun},
	CommandUpsert:      {Init: upsertJsonInit, Run: upsertJsonRun},
}

var logicProcessFunctions = map[string]ProcessFunction{
	CommandPaging:  {Init: pagingInit, Run: pagingRun},
	CommandVersion: {Init: versionInit, Run: versionRun},
}

const (
	ExSeqPrefix                = "EXSEQ_"
	ExSeqIP                    = "IP"
	ExSeqActionName            = "ACTION_NAME"
	ExSeqLevel                 = ExSeqPrefix + "LEVEL"
	ExSeqSuffix                = "_"
	ExSeqReturn                = "ACTION_RETURN"
	ExSeqIgnoreReturn          = "_"
	ExSeqCurrentLevel          = ExSeqPrefix + "CURRENT_LEVEL"
	ExSeqReturnAllDefaultNames = ""
	ExSeqReturnSingleDefault   = "RETURN"
	ExActionLog                = "__ACTION_LOG__"
)

func AddProcessFunction(key string, processor ProcessFunction) {
	processFunctions[key] = processor
}

func AddProcessFunctions(pf map[string]ProcessFunction) {
	for key, processor := range pf {
		processFunctions[key] = processor
	}
}

var ActionLog bool = false

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

func ExecuteProcessFunction(fn *ProcessFunction, pauseTime int, totalTime int, command string, group *sync.WaitGroup, ctx *dvcontext.RequestContext, debugList string) bool {
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
			if debugList != "" {
				DebugShowVariablesByEnvList(ctx, debugList)
			}
			return true
		}
		if pauseTime <= 0 {
			break
		}
		if totalTime > 0 {
			if Log >= dvlog.LogInfo {
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
	return ExecuteProcessFunction(&waiter, pauseTime, totalTime, command, nil, nil, "")
}

func ExecuteReturnSubsequence(ctx *dvcontext.RequestContext, retValue string) {
	level := ctx.PrimaryContextEnvironment.GetInt(ExSeqLevel)
	namePref := ExSeqPrefix + strconv.Itoa(level) + ExSeqSuffix
	param := ctx.LocalContextEnvironment.GetString(namePref + ExSeqReturn)
	var returnKeys []string
	var returnValues []interface{}
	returnNumber := 0
	if param != ExSeqIgnoreReturn && retValue != "" {
		returnKeys = dvtextutils.ConvertToNonEmptyList(retValue)
		returnNumber = len(returnKeys)
		returnValues = make([]interface{}, returnNumber)
		var ok bool
		for i := 0; i < returnNumber; i++ {
			returnValues[i], ok = ctx.LocalContextEnvironment.Get(returnKeys[i])
			if !ok {
				if returnKeys[i] != ExSeqReturnSingleDefault {
					dvlog.PrintlnError("Value " + returnKeys[i] + " was not defined")
				}
				returnValues[i] = ""
			}
		}
		if param != "" {
			returnKeys = dvtextutils.ConvertToNonEmptyList(param)
			if len(returnKeys) < returnNumber {
				returnNumber = len(returnKeys)
			}
		}
	}
	level--
	ctx.PrimaryContextEnvironment.Set(ExSeqLevel, level)
	if level < 0 {
		ctx.LocalContextEnvironment = ctx.PrimaryContextEnvironment
	} else {
		ctx.LocalContextEnvironment = ctx.LocalContextEnvironment.Prototype
	}
	if returnNumber != 0 {
		for i := 0; i < returnNumber; i++ {
			key := returnKeys[i]
			if key != "_" {
				ctx.LocalContextEnvironment.Set(key, returnValues[i])
			}
		}
	}
}

func pushSubsequence(ctx *dvcontext.RequestContext, actionName string,
	returnKey string, paramStr map[string]string, level int) {
	params := make(map[string]interface{})
	if paramStr != nil && len(paramStr) > 0 {
		for k, v := range paramStr {
			if len(v) > 1 && v[0] == '@' {
				vl, ok := ReadActionResult(v[1:], ctx)
				if !ok {
					vl = ""
				}
				params[k] = vl
			} else {
				params[k] = v
			}
		}
	}
	putSubsequence(ctx, actionName, returnKey, level, params)
}

func putSubsequence(ctx *dvcontext.RequestContext, actionName string,
	returnKey string, level int, params map[string]interface{}) {
	ctx.PrimaryContextEnvironment.Set(ExSeqLevel, level)
	if params == nil {
		params = make(map[string]interface{})
	}
	previousEnvironment := ctx.PrimaryContextEnvironment
	if level > 0 {
		previousEnvironment = ctx.LocalContextEnvironment
	}
	ctx.LocalContextEnvironment = dvevaluation.NewObjectWithPrototype(params, previousEnvironment)
	namePref := ExSeqPrefix + strconv.Itoa(level) + ExSeqSuffix
	ctx.LocalContextEnvironment.Set(namePref+ExSeqIP, 0)
	ctx.LocalContextEnvironment.Set(namePref+ExSeqActionName, actionName)
	ctx.LocalContextEnvironment.Set(namePref+ExSeqReturn, returnKey)
	ctx.LocalContextEnvironment.Set(ExSeqCurrentLevel, level)
}

func ExecuteAddSubsequence(ctx *dvcontext.RequestContext, actionName string,
	paramStr map[string]string, returnKey string) {
	level := ctx.PrimaryContextEnvironment.GetInt(ExSeqLevel)
	level++
	pushSubsequence(ctx, actionName, returnKey, paramStr, level)
}

func ExecuteAddSubsequenceShort(ctx *dvcontext.RequestContext, actionName string) {
	if actionName != "" {
		ExecuteAddSubsequence(ctx, actionName, nil, "")
		ExecuteSequenceCycle(ctx, -1)
	}
}

func ExecuteSequence(startActionName string, ctx *dvcontext.RequestContext, initialParams map[string]string) bool {
	if ctx == nil {
		ctx = &dvcontext.RequestContext{
			Id:                        dvcontext.GetUniqueId(),
			PrimaryContextEnvironment: dvparser.GetGlobalPropertiesAsDvObject(),
		}
	}
	debug, ok := ctx.PrimaryContextEnvironment.Get(startActionName + "_LOG")
	if ok {
		ActionLog = true
		debugMode := 0
		switch debug.(type) {
		case string:
			if strings.ToLower(debug.(string)) == "debug" {
				debugMode = 1
			}
		}
		ctx.PrimaryContextEnvironment.Set(ExActionLog, debugMode)
	}
	if ActionLog {
		dvlog.Printf(startActionName, "ExecuteSequence %s\n", startActionName)
	}
	pushSubsequence(ctx, startActionName, ExSeqReturnAllDefaultNames, initialParams, 0)
	return ExecuteSequenceCycle(ctx, 0)
}

func ExecuteSequenceCycle(ctx *dvcontext.RequestContext, cycleLevel int) bool {
	var wg sync.WaitGroup
	var waitCommand string
	var err error
	if cycleLevel < 0 {
		cycleLevel = ctx.PrimaryContextEnvironment.GetInt(ExSeqLevel)
	}
	debugMode := -1
	if ActionLog {
		debug, ok := ctx.PrimaryContextEnvironment.Get(ExActionLog)
		if ok {
			ActionLog = true
			switch debug.(type) {
			case int:
				debugMode = debug.(int)
			}
		}
	}
	for true {
		level := ctx.PrimaryContextEnvironment.GetInt(ExSeqLevel)
		if level < cycleLevel {
			wg.Wait()
			return true
		}
		namePrefix := ExSeqPrefix + strconv.Itoa(level) + ExSeqSuffix
		ip := ctx.LocalContextEnvironment.GetInt(namePrefix + ExSeqIP)
		if ip < 0 {
			ExecuteReturnSubsequence(ctx, ExSeqReturnSingleDefault)
			continue
		}
		ip++
		ctx.LocalContextEnvironment.Set(namePrefix+ExSeqIP, ip)
		p := ctx.LocalContextEnvironment.GetString(namePrefix+ExSeqActionName) + "_" + strconv.Itoa(ip)
		waitCommandRaw := strings.TrimSpace(ctx.LocalContextEnvironment.GetString(p))
		if waitCommandRaw == "" {
			ExecuteReturnSubsequence(ctx, ExSeqReturnSingleDefault)
			continue
		}
		debugList := ""
		if debugMode >= 0 {
			dvlog.Printf(p, "Executing %s %s", p, waitCommandRaw)
			if debugMode >= 1 {
				DebugShowVariablesByEnvList(ctx, p+"_SHOW_BEFORE")
				debugList = p + "_SHOW_AFTER"
			}
		}
		if ctx == nil {
			waitCommand, err = dvparser.ConvertByteArrayByGlobalProperties([]byte(waitCommandRaw), waitCommandRaw)
		} else {
			waitCommand, err = ctx.LocalContextEnvironment.CalculateString(waitCommandRaw)
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
			if Log >= dvlog.LogInfo {
				dvlog.PrintfError("idle waiting for %d seconds before %s", idleTime, waitCommand)
			}
			time.Sleep(time.Duration(idleTime) * time.Second)
		}
		if Log >= dvlog.LogInfo {
			dvlog.PrintfError("starting waiting %d seconds for %s", totalTime, waitCommand)
		}
		if parallel {
			wg.Add(1)
			ok = true
			go func() {
				defer wg.Done()
				ExecuteProcessFunction(&waiter, pauseTime, totalTime, waitCommand, &wg, ctx, debugList)
			}()
		} else {
			ok = ExecuteProcessFunction(&waiter, pauseTime, totalTime, waitCommand, &wg, ctx, debugList)
		}
		if !ok && strict {
			return false
		}
	}
	return true
}

func ocExecutorStartByEvent(eventName string, data []interface{}) error {
	prefix := "EXECUTE_" + strings.ToUpper(eventName)
	if !ExecuteSequence(prefix, nil, nil) {
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
	AddProcessFunctions(execStatementProcessFunctions)
	AddProcessFunctions(logicProcessFunctions)
	dvmodules.RegisterActionProcessor("", fireAction, false)
	dvmodules.RegisterActionProcessor("single", fireSingleAction, false)
	dvmodules.RegisterActionProcessor("static", fireStaticAction, false)
	dvmodules.RegisterActionProcessor("switch", fireSwitchAction, false)
	dvmodules.RegisterActionProcessor("security", securityEndPointHandler, false)
	dvmodules.RegisterActionProcessor("sse", fireSseAction, false)
	return dvmodules.SubscribeForEvents(ocExecutorRegistrationConfig, false)
}

var registeredOC = RegisterOC()
