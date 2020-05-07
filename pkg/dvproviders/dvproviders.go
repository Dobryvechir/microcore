/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvproviders

import (
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
)

type MethodEndPointHandler func(*dvmeta.RequestContext) bool
type MethodGlobalInitHandler func(map[string]string) error
type MethodServerInitHandler func(params []string) (map[string]string, error)

type RegistrationConfig struct {
	Name              string
	EndPointHandler   dvmeta.ProcessorEndPointHandler
	GlobalInitHandler MethodGlobalInitHandler
	ServerInitHandler MethodServerInitHandler
}

type ProviderConfig struct {
	Name   string   `json:"name"`
	Urls   string   `json:"url"`
	Params []string `json:"params"`
}

var registeredConfigs map[string]*RegistrationConfig = make(map[string]*RegistrationConfig)

func RegisterProvider(config *RegistrationConfig, silent bool) bool {
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

func GetParametersName(name string) string {
	return "dvproviders." + name + ".params"
}
func GetRegisteredConfig(name string, silent bool) *RegistrationConfig {
	config, ok := registeredConfigs[name]
	if !silent && !ok {
		panic("Method with name " + name + " is not present")
	}
	return config
}

func createProviderBlock(config *RegistrationConfig, provider *ProviderConfig) *dvmeta.ProcessorBlock {
	f := config.EndPointHandler
	if f == nil {
		panic("EndPointHandler is obligatory field but not specified in provider " + config.Name)
	}
	var data map[string]string
	var err error
	if config.ServerInitHandler != nil {
		data, err = config.ServerInitHandler(provider.Params)
		if err != nil {
			panic("Incorrect parameters for provider " + config.Name + ": " + err.Error())
		}
	}
	urls := dvurl.PreparseMaskExpressions(provider.Urls)
	return &dvmeta.ProcessorBlock{
		Name:            GetParametersName(config.Name),
		EndPointHandler: f,
		Urls:            urls,
		Data:            data,
	}
}

func MakeProviderGlobalInitialization(providerInits map[string]map[string]string) {
	if providerInits == nil {
		return
	}
	for providerName, providerGlobalValues := range providerInits {
		config := GetRegisteredConfig(providerName, false)
		if config.GlobalInitHandler != nil {
			if err := config.GlobalInitHandler(providerGlobalValues); err != nil {
				panic("Incorrect provider " + providerName + "'s global parameters:" + err.Error())
			}
		} else if providerGlobalValues != nil {
			panic("Provider " + providerName + " has no global values to be initialized, but you specified them in the config")
		}
	}
}

func MakeProviderBlocks(configs []ProviderConfig) (res []dvmeta.ProcessorBlock) {
	if configs == nil {
		return nil
	}
	res = make([]dvmeta.ProcessorBlock, len(configs))
	for i, f := range configs {
		config := GetRegisteredConfig(f.Name, false)
		res[i] = *createProviderBlock(config, &f)
	}
	return
}

func PlaceProviderReferences(request *dvmeta.RequestContext) {
	providers := request.Server.BaseProviderBlocks
	n := len(providers)
	urls := request.Urls
	for i := 0; i < n; i++ {
		provider := providers[i]
		urlRules := provider.Urls
		if len(urlRules) != 0 && !dvurl.MatchMasksWithDefault(urlRules, urls, dvurl.MatchDefaultFalse, request.ExtraAsDvObject) {
			continue
		}
		if provider.Data != nil {
			request.Extra[provider.Name] = provider.Data
		}
		_ = provider.EndPointHandler(request)
	}
}
