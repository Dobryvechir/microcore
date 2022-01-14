/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvconfig

import (
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvprocessors"
	"github.com/Dobryvechir/microcore/pkg/dvproviders"
)

// ServerStart starts http server, the config and properties are read from the current folder or by other options
func ServerStart() {
	cf := ReadConfig()
	args := dvparser.GetCommandLine()
	argsLen := len(args)
	action := ""
	osargs2 := ""
	if argsLen > 0 {
		action = args[0]
	}
	if argsLen > 1 {
		osargs2 = args[1]
	}
	switch action {
	case "config":
		SaveConfig(osargs2, cf)
	case "finish":
		dvcom.ProcessHosts(cf.Hosts, true)
		dvcom.ResolveAdministrativeTasks()
	case "execute":
		if osargs2 == "" {
			osargs2 = "SERVER"
		} else {
			osargs2 = strings.ToUpper(osargs2)
		}
		data, err := json.Marshal(cf)
		if err != nil {
			log.Printf("Could not convert config %v", err)
		} else {
			dvparser.SetGlobalPropertiesValue("CURRENT_MICROCORE_CONFIG", string(data))
		}
		ProvideServerCommand()
		ctx:=&dvcontext.RequestContext{PrimaryContextEnvironment: dvparser.GetGlobalPropertiesAsDvObject()}
		dvaction.ExecuteSequence("EXECUTE_"+osargs2, ctx, nil)
		if ctx.StatusCode>=400 {
			log.Printf("Error %s ", string(ctx.Output))
		}
	default:
		serverStartByConfigDirect(cf)
	}
}

// ProvideServerCommand registers the http server as server for command execution purposes
func ProvideServerCommand() {
	dvaction.AddProcessFunction("server", dvaction.ProcessFunction{
		Init:  processServerInit,
		Run:   processServerRun,
		Async: true,
	})
}

func processServerInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	command = strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if command == "" || command[0] != '{' || command[len(command)-1] != '}' {
		log.Printf("Invalid execution server command, config expected {}")
		return nil, false
	}
	cf := &DvConfig{}
	err := json.Unmarshal([]byte(command), cf)
	if err != nil {
		log.Printf("Error in config %s: %v", command, err)
		return nil, false
	}
	return []interface{}{cf}, true
}

func processServerRun(data []interface{}) bool {
	config := data[0].(*DvConfig)
	wg := data[1].(*sync.WaitGroup)
	go func() {
		ServerStartByConfig(config)
		wg.Done()
	}()
	return true
}

// ServerStartByConfig starts http server by the config provided as parameters
func ServerStartByConfig(cf *DvConfig) {
	postConfigInit(cf)
	serverStartByConfigDirect(cf)
}

func serverStartByConfigDirect(cf *DvConfig) {
	dvlog.StartingLogFile()
	dvcom.ProcessHosts(cf.Hosts, false)
	dvmodules.MakeModuleGlobalInitialization(cf.Modules)
	dvmodules.MakeHookGlobalInitialization(cf.Hooks)
	dvprocessors.MakeProcessorGlobalInitialization(cf.Processors)
	dvprocessors.MakePostProcessorGlobalInitialization(cf.PostProcessors)
	dvproviders.MakeProviderGlobalInitialization(cf.Providers)
	ProcessBaseFolder(&cf.Server, cf.HostServers)
	dvcom.ResolveAdministrativeTasks()
	err := dvmodules.FireStartHookEvent([]interface{}{cf})
	if err != nil {
		log.Printf("Error on start: %v", err)
		return
	}
	if cf.Listen != "" {
		log.Printf("Started MicroCore%v on %s \n", dvlog.StartTime, cf.Listen)
		http.ListenAndServe(cf.Listen, nil)
		log.Print("Impossible to occupy " + cf.Listen)
	} else {
		log.Printf("listen server is not specified  %v\n", cf)
	}
}
