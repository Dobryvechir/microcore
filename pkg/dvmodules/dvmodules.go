/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
)

type MethodEndPointHandler func(request *dvmeta.RequestContext) bool
type MethodGlobalInitHandler func(map[string]string) error
type MethodServerInitHandler func(params []string) (map[string]string, error)
type MethodOwnHandlerGenerator func(url string, params []string, urlPool *dvurl.UrlPool) error

type RegistrationConfig struct {
	Name              string
	EndPointHandler   MethodEndPointHandler
	GlobalInitHandler MethodGlobalInitHandler
	ServerInitHandler MethodServerInitHandler
	GenerateHandlers  MethodOwnHandlerGenerator
}

type ModuleConfig struct {
	Name   string   `json:"name"`
	Url    string   `json:"url"`
	Params []string `json:"params"`
}

var registeredConfigs = make(map[string]*RegistrationConfig)

func RegisterModule(config *RegistrationConfig, silent bool) bool {
	name := config.Name
	if _, ok := registeredConfigs[name]; ok {
		if silent {
			return false
		}
		panic("Processor with name " + name + " already registered")
	}
	registeredConfigs[name] = config
	return true
}

func GetRegisteredConfig(name string, silent bool) *RegistrationConfig {
	config, ok := registeredConfigs[name]
	if !silent && !ok {
		panic("Method with name " + name + " is not present")
	}
	return config
}

func MakeModuleHandler(config *RegistrationConfig, params []string) dvmeta.HandlerFunc {
	f := config.EndPointHandler
	if f == nil {
		panic("EndPointHandler is obligatory field but not specified in module " + config.Name)
	}
	var data map[string]string
	var err error
	if config.ServerInitHandler != nil {
		data, err = config.ServerInitHandler(params)
		if err != nil {
			panic("Incorrect parameters for module " + config.Name + ": " + err.Error())
		}
	}
	if data == nil {
		data = make(map[string]string)
	}
	return func(request *dvmeta.RequestContext) bool {
		request.Params = data
		return f(request)
	}
}

func MakeModuleGlobalInitialization(moduleInits map[string]map[string]string) {
	if moduleInits == nil {
		return
	}
	for moduleName, moduleGlobalValues := range moduleInits {
		config := GetRegisteredConfig(moduleName, false)
		if config.GlobalInitHandler != nil {
			if err := config.GlobalInitHandler(moduleGlobalValues); err != nil {
				panic("Incorrect module " + moduleName + "'s global parameters:" + err.Error())
			}
		} else if moduleGlobalValues != nil {
			panic("Module " + moduleName + " has no global values to be initialized, but you specified them in the config")
		}
	}
}

func RegisterEndPointHandlers(configs []ModuleConfig) dvmeta.HandlerFunc {
	if len(configs) == 0 {
		return nil
	}
	base := dvurl.GetUrlHandler()
	for _, f := range configs {
		config := GetRegisteredConfig(f.Name, false)
		if config.GenerateHandlers != nil {
			err := config.GenerateHandlers(f.Url, f.Params, base)
			if err != nil {
				panic("Module registration error: " + err.Error())
			}
		} else {
			handler := MakeModuleHandler(config, f.Params)
			if f.Url != "" {
				base.RegisterHandlerFunc(f.Url, handler)
			}
		}
	}
	return getHandlerFunc(base)
}

func urlVerifier(context interface{}, resolver *dvurl.UrlResolver, urlData *dvurl.UrlResultInfo) bool {
	requestContext := context.(*dvmeta.RequestContext)
	requestContext.UrlInlineParams = urlData.UrlKeys
	handler := resolver.Handler.(dvmeta.HandlerFunc)
	return handler(requestContext)
}

func getHandlerFunc(urlPool *dvurl.UrlPool) dvmeta.HandlerFunc {
	return func(context *dvmeta.RequestContext) bool {
		urls := context.Urls
		ok, _ := dvurl.UrlSearch(context, urlPool, urls, urlVerifier, context.ExtraAsDvObject)
		return ok
	}
}
