/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvsession

import (
	"log"
	"strconv"
)

const (
	SESSION_ENGINE_FATAL_IMPORTANCE = iota
	SESSION_ENGINE_NON_IMPORTANCE
	SESSION_ENGINE_CAN_BE_MEMORY
	SESSION_ENGINE_CAN_BE_MEMORY_OR_NONE
)
const DefaultRetentionTime = 3600 * 3

const MEMORY_SESSION_ENGINE = "MEMORY"
const SESSION_RETENTION_TIME = "RETENTION_TIME"

type SessionRequest interface {
	Init(id string, createOnly bool, updateOnly bool) (SessionStorage, error, bool)
}

type SessionEngine interface {
	Init(map[string]string) (SessionRequest, error)
	Close()
}

type SessionStorage interface {
	SetItem(key string, value interface{})
	GetItem(key string) interface{}
	RemoveItem(key string)
	Clear()
	Keys() []string
	Values() map[string]interface{}
}

var sessionEngines = make(map[string]SessionEngine, 3)

func RegisterSessionEngine(name string, engine SessionEngine) bool {
	sessionEngines[name] = engine
	return true
}

func GetSessionRequest(name string, params map[string]string, option int) SessionRequest {
	engine := GetSessionEngine(name, option)
	if engine == nil {
		return nil
	}
	request, err := engine.Init(params)
	if err == nil {
		return request
	}
	log.Printf("Session init failed %v", err)
	switch option {
	case SESSION_ENGINE_NON_IMPORTANCE:
		return nil
	case SESSION_ENGINE_CAN_BE_MEMORY:
		if name != MEMORY_SESSION_ENGINE {
			return GetSessionRequest(MEMORY_SESSION_ENGINE, params, option)
		}
	case SESSION_ENGINE_CAN_BE_MEMORY_OR_NONE:
		if name != MEMORY_SESSION_ENGINE {
			return GetSessionRequest(MEMORY_SESSION_ENGINE, params, option)
		}
		return nil
	}
	panic("SESSION IS Critical")
}

func GetRetentionTime(params map[string]string) int {
	v := params[SESSION_RETENTION_TIME]
	if v == "" {
		return DefaultRetentionTime
	}
	tm, err := strconv.Atoi(v)
	if err != nil || tm < 30 {
		log.Printf("Error in session retention time %s %v", v, err)
		return DefaultRetentionTime
	}
	return tm
}

func GetSessionEngine(name string, option int) SessionEngine {
	if name == "" {
		for k, v := range sessionEngines {
			if k != MEMORY_SESSION_ENGINE {
				return v
			}
		}
		name = MEMORY_SESSION_ENGINE
	}
	engine, ok := sessionEngines[name]
	if ok {
		return engine
	}
	switch option {
	case SESSION_ENGINE_NON_IMPORTANCE:
		return nil
	case SESSION_ENGINE_CAN_BE_MEMORY, SESSION_ENGINE_CAN_BE_MEMORY_OR_NONE:
		engine, ok = sessionEngines[MEMORY_SESSION_ENGINE]
		if ok {
			return engine
		}
		if option == SESSION_ENGINE_CAN_BE_MEMORY_OR_NONE {
			return nil
		}
	}
	panic("Session engine " + name + " not found")
}
