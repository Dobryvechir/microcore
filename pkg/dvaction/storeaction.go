/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
	"os"
)

type StoreConfig struct {
	Storage string            `json:"storage"`
	Data    map[string]string `json:"data"`
}

func storeInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &StoreConfig{}
	if !DefaultInitWithObject(command, config) {
		return nil, false
	}
	if config.Storage == "" {
		log.Printf("storage must be specified in %s", command)
		return nil, false
	}
	if config.Storage != "fs" {
		log.Printf("unknown storage in %s", command)
		return nil, false
	}
	if config.Data == nil {
		log.Printf("no data to store in storage in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func storeRun(data []interface{}) bool {
	config := data[0].(*StoreConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return StoreRunByConfig(config, ctx)
}

func StoreRunByConfig(config *StoreConfig, ctx *dvcontext.RequestContext) bool {
	for k, v := range config.Data {
		value, ok := ctx.LocalContextEnvironment.Get(v)
		if !ok {
			continue
		}
		err := StoreValue(k, value)
		if err != nil {
			dvlog.Printf("Error writing %s to %s: %v\n", v, k, err)
		}
	}
	return true
}

func StoreValue(path string, data interface{}) error {
	s := dvevaluation.AnyToString(data)
	err := os.WriteFile(path, []byte(s), 0466)
	return err
}
