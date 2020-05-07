/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

const (
	HookStartEvent = "START"
)

type HookMethodEndPointHandler func(eventName string, data []interface{}) error

type HookRegistrationConfig struct {
	Name              string
	HookEventMapper   map[string]HookMethodEndPointHandler
	GlobalInitHandler MethodGlobalInitHandler
	ServerInitHandler MethodServerInitHandler
}

var hookRegisteredConfigs = make(map[string]*HookRegistrationConfig)
var hookRegisteredEvents = make(map[string][]*HookRegistrationConfig)

func SubscribeForEvents(config *HookRegistrationConfig, silent bool) bool {
	name := config.Name
	if _, ok := hookRegisteredConfigs[name]; ok {
		if silent {
			return false
		}
		panic("Hook with name " + name + " already registered")
	}
	hookRegisteredConfigs[name] = config
	if config.HookEventMapper != nil {
		for eventName, handler := range config.HookEventMapper {
			if handler != nil {
				if hookRegisteredEvents[eventName] == nil {
					hookRegisteredEvents[eventName] = make([]*HookRegistrationConfig, 1, 3)
					hookRegisteredEvents[eventName][0] = config
				} else {
					hookRegisteredEvents[eventName] = append(hookRegisteredEvents[eventName], config)
				}
			}
		}
	}
	return true
}

func GetRegisteredHookConfig(name string, silent bool) *HookRegistrationConfig {
	config, ok := hookRegisteredConfigs[name]
	if !silent && !ok {
		panic("Hook with name " + name + " is not present")
	}
	return config
}

func MakeHookGlobalInitialization(moduleInits map[string]map[string]string) {
	if moduleInits == nil {
		return
	}
	for moduleName, moduleGlobalValues := range moduleInits {
		config := GetRegisteredHookConfig(moduleName, false)
		if config.GlobalInitHandler != nil {
			if err := config.GlobalInitHandler(moduleGlobalValues); err != nil {
				panic("Incorrect module " + moduleName + "'s global parameters:" + err.Error())
			}
		} else if moduleGlobalValues != nil {
			panic("Hook " + moduleName + " has no global values to be initialized, but you specified them in the config")
		}
	}
}

func FireHookEvent(eventName string, data []interface{}) error {
	hooks := hookRegisteredEvents[eventName]
	n := len(hooks)
	for i := 0; i < n; i++ {
		if hooks[i].HookEventMapper != nil {
			handler := hooks[i].HookEventMapper[eventName]
			if handler != nil {
				err := handler(eventName, data)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func FireStartHookEvent(data []interface{}) error {
	return FireHookEvent(HookStartEvent, data)
}
