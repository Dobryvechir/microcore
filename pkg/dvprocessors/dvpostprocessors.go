/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
)

var registeredPostConfigs map[string]*RegistrationConfig = make(map[string]*RegistrationConfig)

func RegisterPostProcessor(config *RegistrationConfig, silent bool) bool {
	name := config.Name
	if _, ok := registeredPostConfigs[name]; ok {
		if silent {
			return false
		}
		panic("Post Processor with name " + name + " already registered")
	}
	registeredPostConfigs[name] = config
	return true
}

func GetRegisteredConfigForPost(name string, silent bool) *RegistrationConfig {
	config, ok := registeredPostConfigs[name]
	if !silent && !ok {
		panic("Post Processor with name " + name + " is not present")
	}
	return config
}

func InitializePostProcessors(processorConfigs []ProcessorConfig) (blocks []dvmeta.ProcessorBlock) {
	n := len(processorConfigs)
	if n == 0 {
		return nil
	}
	blocks = make([]dvmeta.ProcessorBlock, n)
	for i := 0; i < n; i++ {
		name := processorConfigs[i].Name
		params := processorConfigs[i].Params
		urls := dvurl.PreparseMaskExpressions(processorConfigs[i].Urls)
		config := GetRegisteredConfigForPost(name, false)
		f := config.EndPointHandler
		if f == nil {
			panic("EndPointHandler is obligatory field but not specified in post processor " + name)
		}
		var data map[string]string
		var err error
		if config.ServerInitHandler != nil {
			data, err = config.ServerInitHandler(params)
			if err != nil {
				panic("Incorrect parameters for post processor " + name + ": " + err.Error())
			}
		}
		if data == nil {
			data = make(map[string]string)
		}
		blocks[i] = dvmeta.ProcessorBlock{EndPointHandler: f, Urls: urls, Data: data}
	}
	return
}

func MakePostProcessorGlobalInitialization(processorInits map[string]map[string]string) {
	if processorInits == nil {
		return
	}
	for processorName, processorGlobalValues := range processorInits {
		config := GetRegisteredConfigForPost(processorName, false)
		if config.GlobalInitHandler != nil {
			if err := config.GlobalInitHandler(processorGlobalValues); err != nil {
				panic("Incorrect post processor " + processorName + "'s global parameters:" + err.Error())
			}
		} else if processorGlobalValues != nil {
			panic("Post processor " + processorName + " has no global values to be initialized, but you specified them in the config")
		}
	}
}
