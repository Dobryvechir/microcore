// package dvoc orchestrates actions, executions
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdbmanager"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

/*****************************************************
ACTION_PICTURE_READ_1=recordreadall:{"table":"picture","result":"request:RESULT"}

ACTION_PICTURE_UPLOAD_1=recordcreate:{"table":"picture","result":"request:RESULT"}

ACTION_PICTURE_DELETE_1=recorddelete:{"table":"picture","key":"URL_PARAM_ID","result":"request:RESULT"}

ACTION_GROUP_ALL_1=recordreadone:{"table":"group","key":"URL_PARAM_PARENT","result":"request:RESULT"}
ACTION_GROUP_ALL_2=recordbind:{"table":"group","src":"groupIds","dst":"groups","root":"request:RESULT","kind":"array","fields":"id,title"}
ACTION_GROUP_ALL_3=recordscan:{"table":"group","fields":"id,title","result":"request.RESULT.allGroups"}

ACTION_GROUP_CREATE_1=recordcreate:{"table":"group","result":"request:RESULT"}

ACTION_GROUP_UPDATE_1=recordupdate:{"table":"group","result":"request:RESULT"}

ACTION_GROUP_DELETE_1=recorddelete:{"table":"group","key":"URL_PARAM_ID","result":"request:RESULT"}







	CommandRecordBind:    {Init: recordBindInit, Run: recordBindRun},
	CommandRecordCreate:  {Init: recordCreateInit, Run: recordCreateRun},
	CommandRecordDelete:  {Init: recordDeleteInit, Run: recordDeleteRun},
	CommandRecordReadAll: {Init: recordReadAllInit, Run: recordReadAllRun},
	CommandRecordReadOne: {Init: recordReadOneInit, Run: recordReadOneRun},
	CommandRecordScan:    {Init: recordScanInit, Run: recordScanRun},


*******************************************************/
/************** BIND ***************************************************/

type recordBindConfig struct {
	Table  string `json:"table"`
	Result string `json:"result"`
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordBind(config.Table)
	SaveActionResult(config.Result, r, ctx)
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordCreate(config.Table)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** Delete ***************************************************/

type recordDeleteConfig struct {
	Table  string `json:"table"`
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordDelete(config.Table)
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordReadAll(config.Table)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** ReadOne ***************************************************/

type recordReadOneConfig struct {
	Table  string `json:"table"`
	Result string `json:"result"`
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordReadOne(config.Table)
	SaveActionResult(config.Result, r, ctx)
	return true
}

/************** Scan ***************************************************/

type recordScanConfig struct {
	Table  string `json:"table"`
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordScan(config.Table)
	SaveActionResult(config.Result, r, ctx)
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
	env := GetEnvironment(ctx)
	r := dvdbmanager.RecordUpdate(config.Table)
	SaveActionResult(config.Result, r, ctx)
	return true
}
