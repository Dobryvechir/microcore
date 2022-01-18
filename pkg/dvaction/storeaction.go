/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
	"os"
	"strings"
)

type StoreConfig struct {
	Storage string            `json:"storage"`
	Data    map[string]string `json:"data"`
	Format  string            `json:"format"`
	Before  string            `json:"before"`
	After   string            `json:"after"`
}

func storeInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &StoreConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Storage == "" {
		config.Storage = "fs"
	}
	if config.Storage != "fs" && config.Storage != "text" && config.Storage != "binary" {
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
	env := GetEnvironment(ctx)
	for k, v := range config.Data {
		value, ok := env.Get(v)
		if !ok {
			continue
		}
		err := StoreValue(k, value, config.Format, config.Before, config.After, config.Storage, env)
		if err != nil {
			dvlog.Printf("Error writing %s to %s: %v\n", v, k, err)
		}
	}
	return true
}

func StoreValue(path string, data interface{}, format string, before string, after string, storage string, env *dvevaluation.DvObject) (err error) {
	var s []byte
	switch strings.ToLower(format) {
	case "yaml":
		dv := dvevaluation.AnyToDvVariable(data)
		s = dvjson.PrintToYaml(dv, 2)
	case "yaml_array":
		dv := dvevaluation.AnyToDvVariable(data)
		s = dvjson.PrintToYamlPlainArray(dv, 2)
	default:
		s = dvevaluation.AnyToByteArray(data)
	}
	if before != "" {
		beforeMessage, ok := env.Get(before)
		if !ok {
			beforeMessage = before
		}
		s = append(dvevaluation.AnyToByteArray(beforeMessage), s...)
	}
	if after != "" {
		afterMessage, ok := env.Get(after)
		if !ok {
			afterMessage = after
		}
		s = append(s, dvevaluation.AnyToByteArray(afterMessage)...)
	}
	switch storage {
	case "fs":
		err = os.WriteFile(path, s, 0466)
	case "text":
		env.Set(path, string(s))
	case "binary":
		env.Set(path, s)
	}
	return
}
