/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
)

type UpsertJsonConfig struct {
	Sample *JsonRead `json:"change"`
	Ref    *JsonRead `json:"stored"`
}

func upsertJsonInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &UpsertJsonConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Sample == nil || config.Sample.Var == "" {
		log.Printf("sample.place must be specified in %s", command)
		return nil, false
	}
	if config.Ref == nil || config.Ref.Var == "" {
		log.Printf("reference.place must be present in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func upsertJsonRun(data []interface{}) bool {
	config := data[0].(*UpsertJsonConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return UpsertJsonByConfig(config, ctx)
}

func UpsertJsonByConfig(config *UpsertJsonConfig, ctx *dvcontext.RequestContext) bool {
	sample, err := JsonExtract(config.Sample, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Sample.Var)
		return true
	}
	ref, err := JsonExtract(config.Ref, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Ref.Var)
		return true
	}
	res := dvevaluation.AnyToDvVariable(ref)
	added, _, updated, _, _ := dvjson.FindDifferenceForAnyType(sample, res,
		true, false, true, false,
		false, false, true)
	if updated != nil {
		n := len(updated.Fields)
		for i := 0; i < n; i++ {
			f := updated.Fields[i]
			ind := f.Extra.(int)
			res.Fields[ind] = f
		}
	}
	if added != nil {
		n := len(added.Fields)
		for i := 0; i < n; i++ {
			res.Fields = append(res.Fields, added.Fields[i])
		}
	}
	SaveActionResult(config.Ref.Destination, res, ctx)
	return true
}
