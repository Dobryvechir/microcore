/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvsession

import (
	"time"
)

type MemorySessionStorage struct {
	id         string
	accessTime int64
	values     map[string]interface{}
}

func (storage *MemorySessionStorage) UpdateAccessTime() {
	storage.accessTime = time.Now().UnixNano()
}

func (storage *MemorySessionStorage) SetItem(key string, value interface{}) {
	storage.values[key] = value
	storage.UpdateAccessTime()
}

func (storage *MemorySessionStorage) GetItem(key string) interface{} {
	storage.UpdateAccessTime()
	r := storage.values[key]
	return r
}
func (storage *MemorySessionStorage) RemoveItem(key string) {
	storage.UpdateAccessTime()
	delete(storage.values, key)
}
func (storage *MemorySessionStorage) Clear() {
	storage.UpdateAccessTime()
	storage.values = make(map[string]interface{})
}

func (storage *MemorySessionStorage) Keys() []string {
	storage.UpdateAccessTime()
	n := len(storage.values)
	keys := make([]string, n)
	i := 0
	for k := range storage.values {
		keys[i] = k
		i++
		if i == n {
			break
		}
	}
	return keys
}

func (storage *MemorySessionStorage) Values() map[string]interface{} {
	storage.UpdateAccessTime()
	return storage.values
}

func (storage *MemorySessionStorage) GetId() string {
	return storage.id
}

type MemorySessionRequest struct {
	retentionTime   int
	garbageCount    int
	sessionStorages map[string]*MemorySessionStorage
}

func (req *MemorySessionRequest) Init(id string) (SessionStorage, error, bool) {
	storage, ok := req.sessionStorages[id]
	if ok {
		storage.UpdateAccessTime()
		return storage, nil, true
	}
	storage = &MemorySessionStorage{}
	storage.Clear()
	req.sessionStorages[id] = storage
	if req.garbageCount > 10 {
		req.garbageCount = 0
		req.CollectGarbage()
	} else {
		req.garbageCount++
	}
	return storage, nil, false
}

func (req *MemorySessionRequest) CollectGarbage() {
	t := time.Now().UnixNano() - int64(req.retentionTime)*1000000000
	for k, v := range req.sessionStorages {
		if v.accessTime < t {
			delete(req.sessionStorages, k)
		}
	}
}

type MemorySessionEngine struct {
}

func (engine *MemorySessionEngine) Init(params map[string]string) (SessionRequest, error) {
	retentionTime := GetRetentionTime(params)
	return &MemorySessionRequest{
		retentionTime:   retentionTime,
		sessionStorages: make(map[string]*MemorySessionStorage),
	}, nil
}

func (engine *MemorySessionEngine) Close() {

}

var inited = RegisterSessionEngine(MEMORY_SESSION_ENGINE, &MemorySessionEngine{})
