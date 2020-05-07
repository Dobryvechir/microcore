/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvconfig

import (
	"encoding/json"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvjsmaster"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvprocessors"
	"github.com/Dobryvechir/microcore/pkg/dvproviders"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var LogConfig bool
var GlobalProperties map[string]string

const MICRO_CORE_CONFIG = "MicroCore.conf"
const MICRO_CORE_PROPERTIES = "MicroCore.properties"
const MICRO_CORE_PATH = "MICRO_CORE_PATH"
const MICRO_CORE_NAMESPACE = "MICRO_CORE_CURRENT_NAMESPACE"
const DV_CONFIG_DEBUG_WRITE = "DEBUG_CONFIG_SAVE_FILENAME"

type DvRewrite struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Condition string `json:"condition"`
	Options   string `json:"options"`
}

type DvHostServer struct {
	Hosts                         string                         `json:"hosts"`
	BaseFolder                    string                         `json:"baseFolder"`
	Actions                       []dvmeta.DvAction              `json:"actions"`
	Rewrites                      []DvRewrite                    `json:"rewrites"`
	ExtraServer                   string                         `json:"extraServer"`
	ExtraServerSettings           dvmeta.ServerSettings          `json:"extraServerSettings"`
	ServerRewrites                []DvRewrite                    `json:"serverRewrites"`
	ProxyName                     string                         `json:"proxyName"`
	CacheControl                  string                         `json:"cacheControl"`
	DirectoryIndex                string                         `json:"directoryIndex"`
	HeadersStatic                 map[string]string              `json:"headersStatic"`
	HeadersExtraServer            map[string]string              `json:"headersExtraServer"`
	HeadersStaticOptions          map[string]string              `json:"headersStaticOptions"`
	HeadersExtraServerOptions     map[string]string              `json:"headersExtraServerOptions"`
	AccessControlAllowOrigin      string                         `json:"accessControlAllowOrigin"`
	AccessControlAllowMethod      string                         `json:"accessControlAllowMethod"`
	AccessControlAllowHeaders     string                         `json:"accessControlAllowHeaders"`
	AccessControlMaxAge           string                         `json:"accessControlMaxAge"`
	AccessControlExposeHeaders    string                         `json:"accessControlExposeHeaders"`
	AccessControlAllowCredentials string                         `json:"accessControlAllowCredentials"`
	Modules                       []dvmodules.ModuleConfig       `json:"modules"`
	Processors                    []dvprocessors.ProcessorConfig `json:"processors"`
	PostProcessors                []dvprocessors.ProcessorConfig `json:"postProcessors"`
	Providers                     []dvproviders.ProviderConfig   `json:"providers"`
	HostHeader                    string                         `json:"hostHeader"`
}

type DvConfig struct {
	Namespace      string                       `json:"namespace"`
	Listen         string                       `json:"listen"`
	RootFolder     string                       `json:"rootFolder"`
	LogLevel       string                       `json:"logLevel"`
	LogModules     string                       `json:"logModules"`
	Hosts          []dvcom.IpList               `json:"hosts"`
	HostServers    []DvHostServer               `json:"hostServers"`
	SystemInfo     string                       `json:"systemInfo"`
	Server         DvHostServer                 `json:"server"`
	Modules        map[string]map[string]string `json:"modules"`
	Hooks          map[string]map[string]string `json:"hooks"`
	Processors     map[string]map[string]string `json:"processors"`
	PostProcessors map[string]map[string]string `json:"postProcessors"`
	Providers      map[string]map[string]string `json:"providers"`
}

var CurrentDir string

func SaveConfig(place string, cf *DvConfig) {
	if place == "" || place == "." {
		place = CurrentDir + "/" + MICRO_CORE_CONFIG
	} else if place == "#" {
		place = dvlog.GetTaskFolder() + "/" + MICRO_CORE_CONFIG
	} else {
		fmt.Printf("To configure, you can specify only . or #, not " + place)
		return
	}
	configStr, err := json.Marshal(cf)
	if err != nil {
		fmt.Printf("Error converting the config to json: %s", err.Error())
		return
	}
	err = ioutil.WriteFile(place, configStr, os.ModePerm)
	if err != nil {
		fmt.Printf("Error %s writing the config to file %s", err.Error(), place)
	} else {
		fmt.Printf("Config was successfully written to %s", place)
	}
}

func ResetNamespaceFolder() {
	namespaceFolder := dvlog.GetPrincipalFolder(false) + "/" + dvlog.CurrentNamespace
	if _, err4 := os.Stat(namespaceFolder); err4 == nil {
		if dvparser.GeneralNamespaceFolderIndex < 0 {
			dvparser.GeneralNamespaceFolderIndex = len(dvparser.GeneralFilePaths)
			dvparser.GeneralFilePaths = append(dvparser.GeneralFilePaths, namespaceFolder)
		} else {
			dvparser.GeneralFilePaths[dvparser.GeneralNamespaceFolderIndex] = namespaceFolder
		}
	}

}

func setFilePaths() error {
	if namespace, isOk := dvparser.GlobalProperties[MICRO_CORE_NAMESPACE]; isOk {
		dvlog.SetCurrentNamespace(namespace)
	}
	ResetNamespaceFolder()
	if path, ok := dvparser.GlobalProperties[MICRO_CORE_PATH]; ok {
		if _, err := os.Stat(path); err == nil {
			if dvparser.GeneralMicroCoreFolderIndex < 0 {
				dvparser.GeneralMicroCoreFolderIndex = len(dvparser.GeneralFilePaths)
				dvparser.GeneralFilePaths = append(dvparser.GeneralFilePaths, path)
			} else {
				dvparser.GeneralFilePaths[dvparser.GeneralMicroCoreFolderIndex] = path
			}
		}
	}
	return nil
}

func readGlobalProperties(currentDir string, propertiesName string) {
	err := dvparser.ReadPropertiesFileWithEnvironmentVariables(currentDir, propertiesName, setFilePaths)
	if err != nil {
		panic("\n" + err.Error())
	}
}

func FindAndReadConfigs(configName string, propertiesName string) string {
	var err error
	CurrentDir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
		return ""
	}
	dvjsmaster.Init()
	readGlobalProperties(CurrentDir, propertiesName)
	return dvparser.FindInGeneralPaths(configName)
}

func ReadConfig() *DvConfig {
	filename := FindAndReadConfigs(MICRO_CORE_CONFIG, MICRO_CORE_PROPERTIES)
	cf := &DvConfig{}
	if filename == "" {
		cf.Namespace = dvlog.CurrentNamespace
		cf.Listen = ":80"
		cf.Server = DvHostServer{BaseFolder: "."}
	} else {
		data, err := dvparser.SmartReadTemplate(filename, 3, byte(' '))
		if err == nil {
			dvlog.CleanEOL(data)
			if saveConfig, okSave := dvparser.GlobalProperties[DV_CONFIG_DEBUG_WRITE]; okSave {
				err2 := ioutil.WriteFile(saveConfig, data, 0644)
				if err2 != nil {
					log.Print("Cannot write resulted config to " + saveConfig + ": " + err2.Error())
				}
			}
			err = json.Unmarshal(data, cf)
		}
		if err != nil {
			err2 := ioutil.WriteFile(CurrentDir+"/debug_microcore_conf.json", data, 0644)
			if err2 != nil {
				log.Print("Cannot write ./debug_microcore_conf.json: " + err2.Error())
			}
			panic("\nError: (see debug_microcore_conf.json)Incorrect json in " + filename + ": " + err.Error())
		}
	}
	postConfigInit(cf)
	return cf
}

func postConfigInit(cf *DvConfig) {
	dvparser.DvParserLog = false
	dvlog.SetCurrentNamespace(cf.Namespace)
	ResetNamespaceFolder()
	if cf.RootFolder != "" {
		dvlog.CurrentRootFolder = cf.RootFolder
	}
	if cf.LogLevel != "" {
		dvlog.SetLogLevel(cf.LogLevel)
	}
	logModules := strings.TrimSpace(cf.LogModules)
	if logModules != "" {
		logMods := dvparser.ConvertToList(logModules)
		for _, logModule := range logMods {
			if logModule == "" {
				continue
			}
			switch logModule {
			case "config":
				LogConfig = true
				if dvlog.CurrentLogLevel >= dvlog.LogError {
					dvparser.DvParserLog = true
				}
			case "processors":
				dvjson.LogProcessors = true
			case "db":
				dvjson.LogCrud = true
			case "json":
				dvjson.LogJson = true
			case "server":
				dvcom.LogServer = true
			case "file-server":
				dvcom.LogFileServer = true
			case "hosts":
				dvcom.LogHosts = true
			default:
				log.Print("Available modules for logging are only server, file-server, config, processors, db and json, not " + logModule)
			}
		}
	}
	if LogConfig && dvlog.CurrentLogLevel >= dvlog.LogInfo {
		dvparser.LogVariables("GLOBAL VARIABLES", dvparser.GlobalProperties)
	}
}
