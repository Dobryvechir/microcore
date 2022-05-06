/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvsession

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"sync"
)

var globalMapSessionStorage SessionStorage
var syncArrayManipulation sync.Mutex

func checkGlobalMapSessionStorage() bool {
	syncArrayManipulation.Lock()
	if globalMapSessionStorage != nil {
		syncArrayManipulation.Unlock()
		return true
	}
	engine := GetSessionRequest("", nil, SESSION_ENGINE_CAN_BE_MEMORY)
	stor, err, _ := engine.Init("GlobalMap", true, true)
	if err != nil {
		dvlog.PrintfError("Error memory storage creating %v", err)
		syncArrayManipulation.Unlock()
		return false
	}
	globalMapSessionStorage = stor
	syncArrayManipulation.Unlock()
	return true
}

func GlobalMapWrite(mapName string, key string, val interface{}) {
	if globalMapSessionStorage == nil {
		if !checkGlobalMapSessionStorage() {
			dvlog.PrintfError("Cannot process global map %s", mapName)
			return
		}
	}
	syncArrayManipulation.Lock()
	mp := globalMapSessionStorage.GetItem(mapName)
	iMap, ok := mp.(map[string]interface{})
	if !ok || iMap == nil {
		iMap = make(map[string]interface{})
	}
	iMap[key] = val
	globalMapSessionStorage.SetItem(mapName, iMap)
	syncArrayManipulation.Unlock()
}

func GlobalMapRead(mapName string, key string) (interface{}, bool) {
	if globalMapSessionStorage == nil {
		if !checkGlobalMapSessionStorage() {
			dvlog.PrintfError("Cannot process global map %s", mapName)
			return nil, false
		}
	}
	mp := globalMapSessionStorage.GetItem(mapName)
	iMap, ok := mp.(map[string]interface{})
	if !ok || iMap == nil {
		return nil, false
	}
	syncArrayManipulation.Lock()
	res, isOk := iMap[key]
	syncArrayManipulation.Unlock()
	return res, isOk
}

func GlobalMapDelete(mapName string, key string) bool {
	if globalMapSessionStorage == nil {
		if !checkGlobalMapSessionStorage() {
			dvlog.PrintfError("Cannot process global map %s", mapName)
			return false
		}
	}
	syncArrayManipulation.Lock()
	mp := globalMapSessionStorage.GetItem(mapName)
	iMap, ok := mp.(map[string]interface{})
	if !ok || iMap == nil {
		syncArrayManipulation.Unlock()
		return false
	}
	delete(iMap, key)
	if len(iMap) == 0 {
		globalMapSessionStorage.RemoveItem(mapName)
	} else {
		globalMapSessionStorage.SetItem(mapName, iMap)
	}
	syncArrayManipulation.Unlock()
	return true
}
