/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
)

type ProcessorGlobalInitHandler func(map[string]string) error
type ProcessorServerInitHandler func(params []string) (map[string]string, error)

type RegistrationConfig struct {
	Name              string
	EndPointHandler   dvcontext.ProcessorEndPointHandler
	GlobalInitHandler ProcessorGlobalInitHandler
	ServerInitHandler ProcessorServerInitHandler
}

type ProcessorConfig struct {
	Name   string   `json:"name"`
	Urls   string   `json:"urls"`
	Params []string `json:"params"`
}

var registeredConfigs map[string]*RegistrationConfig = make(map[string]*RegistrationConfig)

func RegisterProcessor(config *RegistrationConfig, silent bool) bool {
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
		panic("Processor with name " + name + " is not present")
	}
	return config
}

func InitializeProcessors(processorConfigs []ProcessorConfig) (blocks []dvcontext.ProcessorBlock) {
	n := len(processorConfigs)
	if n == 0 {
		return nil
	}
	blocks = make([]dvcontext.ProcessorBlock, n)
	for i := 0; i < n; i++ {
		name := processorConfigs[i].Name
		params := processorConfigs[i].Params
		urls := dvurl.PreparseMaskExpressions(processorConfigs[i].Urls)
		config := GetRegisteredConfig(name, false)
		f := config.EndPointHandler
		if f == nil {
			panic("EndPointHandler is obligatory field but not specified in processor " + name)
		}
		var data map[string]string
		var err error
		if config.ServerInitHandler != nil {
			data, err = config.ServerInitHandler(params)
			if err != nil {
				panic("Incorrect parameters for processor " + name + ": " + err.Error())
			}
		}
		if data == nil {
			data = make(map[string]string)
		}
		blocks[i] = dvcontext.ProcessorBlock{EndPointHandler: f, Urls: urls, Data: data}
	}
	return
}

func MakeProcessorGlobalInitialization(processorInits map[string]map[string]string) {
	if processorInits == nil {
		return
	}
	for processorName, processorGlobalValues := range processorInits {
		config := GetRegisteredConfig(processorName, false)
		if config.GlobalInitHandler != nil {
			if err := config.GlobalInitHandler(processorGlobalValues); err != nil {
				panic("Incorrect processor " + processorName + "'s global parameters:" + err.Error())
			}
		} else if processorGlobalValues != nil {
			panic("Processor " + processorName + " has no global values to be initialized, but you specified them in the config")
		}
	}
}
