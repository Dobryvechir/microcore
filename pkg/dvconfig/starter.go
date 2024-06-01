/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvconfig

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Dobryvechir/microcore/pkg/dvaction"
	_ "github.com/Dobryvechir/microcore/pkg/dvaction/dvdynamic"
	"github.com/Dobryvechir/microcore/pkg/dvssews"

	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
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
		ctx := &dvcontext.RequestContext{
			Id:                        dvcontext.GetUniqueId(),
			PrimaryContextEnvironment: dvparser.GetGlobalPropertiesAsDvObject(),
		}
		dvaction.ExecuteSequence("EXECUTE_"+osargs2, ctx, nil)
		if ctx.StatusCode >= 400 {
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
	dvssews.SetParallelProcessingParameters(cf.ParallelProcessing)
	dvmodules.MakeModuleGlobalInitialization(cf.Modules)
	dvmodules.MakeHookGlobalInitialization(cf.Hooks)
	dvprocessors.MakeProcessorGlobalInitialization(cf.Processors)
	dvprocessors.MakePostProcessorGlobalInitialization(cf.PostProcessors)
	dvproviders.MakeProviderGlobalInitialization(cf.Providers)
	ProcessBaseFolder(&cf.Server, cf.HostServers)
	if len(cf.Dbs) != 0 {
		dvdbmanager.DbManagerInit(cf.Dbs)
	}
	dvcom.ResolveAdministrativeTasks()
	err := dvmodules.FireStartHookEvent([]interface{}{cf})
	if err != nil {
		log.Printf("Error on start: %v", err)
		return
	}
	if len(cf.Listen) == 0 {
		cf.Listen = []string{":80"}
	}
	n := len(cf.Listen)
	done := make(chan bool, n)
	for i := 0; i < n; i++ {
		startListenAndServe(cf.Listen[i], done)
	}
	for i := 0; i < n; i++ {
		<-done
	}
}

func startListenAndServe(listen string, done chan bool) {
	listen = strings.TrimSpace(listen)
	log.Printf("Started MicroCore%v on %s \n", dvlog.StartTime, listen)
	var err error
	pos := strings.Index(listen, ",")
	if pos >= 0 {
		urlPort := strings.TrimSpace(listen[:pos])
		if urlPort == "" {
			urlPort = ":443"
		}
		certificates := strings.TrimSpace(listen[pos+1:])
		certificate1 := certificates
		certificate2 := ""
		pos = strings.Index(certificates, ",")
		if pos >= 0 {
			certificate1 = strings.TrimSpace(certificates[:pos])
			certificate2 = strings.TrimSpace(certificates[pos+1:])
		}
		err = http.ListenAndServeTLS(urlPort, certificate1, certificate2, nil)
	} else {
		err = http.ListenAndServe(listen, nil)
	}
	log.Printf("Impossible to occupy %s : %v", listen, err)
	done <- true
}
