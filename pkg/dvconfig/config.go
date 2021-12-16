// Package dvconfig provides configuration management for the http server
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
package dvconfig

import (
	"encoding/json"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvjsmaster"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvprocessors"
	"github.com/Dobryvechir/microcore/pkg/dvproviders"
)

// LogConfig specifies whether message from package dvconfig must be logged
var LogConfig bool

// GlobalProperties stores initial properties from environment and properties files
var GlobalProperties map[string]string

// MicroCoreConfig is the config file name for http server, which can be customized by application
var MicroCoreConfig = "MicroCore.conf"

// MicroCorePath is an environment variable name to specify a folder where config and other info is stored
var MicroCorePath = "MICRO_CORE_PATH"

// MicroCoreNameSpace is an environment variable name to specify the namespace, which is used for creation a folder inside the global work folder
var MicroCoreNameSpace = "MICRO_CORE_CURRENT_NAMESPACE"

// debugWriteName is an environment variable specifying where to store config file in case of failure to parse it for debug purposes
const debugWriteName = "DEBUG_CONFIG_SAVE_FILENAME"

// DvRewrite describes the url rewriting parameters
type DvRewrite struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Condition string `json:"condition"`
	Options   string `json:"options"`
}
type ProxyServerInfo struct {
	Filter string `json:"filter"`
	Url    string `json:"url"`
	Kind   string `json:"kind"`
}

// DvHostServer collects all parameters for a specific host server
type DvHostServer struct {
	Hosts                         string                         `json:"hosts"`
	BaseFolder                    string                         `json:"baseFolder"`
	Actions                       []*dvcontext.DvAction           `json:"actions"`
	Rewrites                      []DvRewrite                    `json:"rewrites"`
	DefaultProxyServer            string                         `json:"defaultProxy"`
	ProxyServers                  []ProxyServerInfo              `json:"proxyServers"`
	ProxyServerSettings           dvcontext.ServerSettings       `json:"proxyServerSettings"`
	ServerRewrites                []DvRewrite                    `json:"serverRewrites"`
	ProxyName                     string                         `json:"proxyName"`
	DomainName                    string                         `json:"domain"`
	CacheControl                  string                         `json:"cacheControl"`
	DirectoryIndex                string                         `json:"directoryIndex"`
	HeadersStatic                 map[string]string              `json:"headersStatic"`
	HeadersProxyServer            map[string]string              `json:"headersProxyServer"`
	HeadersStaticOptions          map[string]string              `json:"headersStaticOptions"`
	HeadersProxyServerOptions     map[string]string              `json:"headersProxyServerOptions"`
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
	LogLevel                      string                         `json:"logLevel"`
	DynamicActions				  bool						     `json:"dynamic"`
}

// DvConfig is a full structure of the config for http server
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

// CurrentDir is a current folder where the application started
var CurrentDir string

// SetApplicationName allows to customize all config names, property file names and variable prefixes
func SetApplicationName(name string) {
	dvparser.SetPrefixesByApplicationName(name)
	upName := strings.ToUpper(name)
	MicroCoreConfig = name + ".conf"
	MicroCorePath = upName + "_PATH"
	MicroCoreNameSpace = upName + "_CURRENT_NAMESPACE"
}

// SaveConfig saves the already-read config for the test purposes
func SaveConfig(place string, cf *DvConfig) {
	if place == "" || place == "." {
		place = CurrentDir + "/" + MicroCoreConfig
	} else if place == "#" {
		place = dvlog.GetTaskFolder() + "/" + MicroCoreConfig
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

// ResetNamespaceFolder is used to change the current namespace
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
	if namespace, isOk := dvparser.GlobalProperties[MicroCoreNameSpace]; isOk {
		dvlog.SetCurrentNamespace(namespace)
	}
	ResetNamespaceFolder()
	if path, ok := dvparser.GlobalProperties[MicroCorePath]; ok {
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

// FindAndReadConfigs finds the config by config name and properties file starting with the current folder
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

// ReadConfig reads the config and properties by default names for http server
func ReadConfig() *DvConfig {
	filename := FindAndReadConfigs(MicroCoreConfig, dvparser.MicroCorePropertiesInCurrentFolderFileName)
	cf := &DvConfig{}
	if filename == "" {
		cf.Namespace = dvlog.CurrentNamespace
		cf.Listen = ":80"
		cf.Server = DvHostServer{BaseFolder: "."}
	} else {
		data, err := dvparser.SmartReadTemplate(filename, 3, byte(' '))
		if err == nil {
			dvlog.CleanEOL(data)
			if saveConfig, okSave := dvparser.GlobalProperties[debugWriteName]; okSave {
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
		logMods := dvtextutils.ConvertToList(logModules)
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
