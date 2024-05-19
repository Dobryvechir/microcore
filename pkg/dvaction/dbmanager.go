// package dvoc orchestrates actions, executions
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvaction

import (
	"strconv"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

/************** BIND ***************************************************/

type recordBindConfig struct {
	Table      string `json:"table"`
	SrcField   string `json:"src"`
	DstField   string `json:"dst"`
	RootObject string `json:"root"`
	Kind       string `json:"kind"`
	Fields     string `json:"fields"`
}

func recordBindInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordBindConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordBindRun(data []interface{}) bool {
	config := data[0].(*recordBindConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordBindRunByConfig(config, ctx)
}

func recordBindRunByConfig(config *recordBindConfig, ctx *dvcontext.RequestContext) bool {
	root, ok := ReadActionResult(config.RootObject, ctx)
	if !ok {
		return true
	}
	dv := dvevaluation.AnyToDvVariable(root)
	if dv == nil || dv.Kind != dvevaluation.FIELD_OBJECT || len(dv.Fields) == 0 {
		return true
	}
	item, ok := dv.FindChildByKey(config.SrcField)
	if !ok {
		return true
	}
	r, err := dvdbmanager.RecordBind(config.Table, item, config.Kind, config.Fields)
	if err != nil {
		return true
	}
	resItem, ok := dv.FindChildByKey(config.DstField)
	if ok && resItem != nil {
		resItem.Fields = r.Fields
		resItem.Kind = r.Kind
	} else {
		r.Name = []byte(config.DstField)
		dv.Fields = append(dv.Fields, r)
	}
	SaveActionResult(config.RootObject, dv, ctx)
	return true
}

/************** Create ***************************************************/

type recordCreateConfig struct {
	Table  string `json:"table"`
	Result string `json:"result"`
}

func recordCreateInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordCreateConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordCreateRun(data []interface{}) bool {
	config := data[0].(*recordCreateConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordCreateRunByConfig(config, ctx)
}

func recordCreateRunByConfig(config *recordCreateConfig, ctx *dvcontext.RequestContext) bool {
	body := ctx.PrimaryContextEnvironment.GetString(dvcontext.BODY_STRING)
	id := strconv.FormatInt(ctx.Id, 10)
	r := dvdbmanager.RecordCreate(config.Table, body, id)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** Delete ***************************************************/

type recordDeleteConfig struct {
	Table  string `json:"table"`
	Key    string `json:"key"`
	Result string `json:"result"`
}

func recordDeleteInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordDeleteConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordDeleteRun(data []interface{}) bool {
	config := data[0].(*recordDeleteConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordDeleteRunByConfig(config, ctx)
}

func recordDeleteRunByConfig(config *recordDeleteConfig, ctx *dvcontext.RequestContext) bool {
	key, ok := ReadActionResult(config.Key, ctx)
	if !ok {
		SaveActionResult(config.Result, "Error key "+config.Key+" is not provided", ctx)
		return true
	}
	v := dvevaluation.AnyToString(key)
	r := dvdbmanager.RecordDelete(config.Table, v)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** ReadAll ***************************************************/

type recordReadAllConfig struct {
	Table  string `json:"table"`
	Result string `json:"result"`
}

func recordReadAllInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordReadAllConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordReadAllRun(data []interface{}) bool {
	config := data[0].(*recordReadAllConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordReadAllRunByConfig(config, ctx)
}

func recordReadAllRunByConfig(config *recordReadAllConfig, ctx *dvcontext.RequestContext) bool {
	r := dvdbmanager.RecordReadAll(config.Table)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** ReadOne ***************************************************/

type recordReadOneConfig struct {
	Table  string `json:"table"`
	Result string `json:"result"`
	Key    string `json:"key"`
}

func recordReadOneInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordReadOneConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordReadOneRun(data []interface{}) bool {
	config := data[0].(*recordReadOneConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordReadOneRunByConfig(config, ctx)
}

func recordReadOneRunByConfig(config *recordReadOneConfig, ctx *dvcontext.RequestContext) bool {
	key, ok := ReadActionResult(config.Key, ctx)
	if !ok {
		SaveActionResult(config.Result, "Error key "+config.Key+" is not provided", ctx)
		return true
	}
	r := dvdbmanager.RecordReadOne(config.Table, key)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** Scan ***************************************************/

type recordScanConfig struct {
	Table  string `json:"table"`
	Fields string `json:"fields"`
	Result string `json:"result"`
}

func recordScanInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordScanConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordScanRun(data []interface{}) bool {
	config := data[0].(*recordScanConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordScanRunByConfig(config, ctx)
}

func recordScanRunByConfig(config *recordScanConfig, ctx *dvcontext.RequestContext) bool {
	r, err := dvdbmanager.RecordScan(config.Table, config.Fields)
	if err != nil {
		SaveActionResult(config.Result, err.Error(), ctx)
	} else {
		SaveActionResult(config.Result, r, ctx)
	}
	return true
}

/************** Update ***************************************************/

type recordUpdateConfig struct {
	Table  string `json:"table"`
	Result string `json:"result"`
}

func recordUpdateInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &recordUpdateConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func recordUpdateRun(data []interface{}) bool {
	config := data[0].(*recordUpdateConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return recordUpdateRunByConfig(config, ctx)
}

func recordUpdateRunByConfig(config *recordUpdateConfig, ctx *dvcontext.RequestContext) bool {
	body := ctx.PrimaryContextEnvironment.GetString(dvcontext.BODY_STRING)
	r := dvdbmanager.RecordUpdate(config.Table, body)
	SaveActionResult(config.Result, r, ctx)
	return true
}
