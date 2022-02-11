/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strings"
)

func dynamicActionInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	cmd := strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if cmd == "" {
		log.Printf("Empty parameters in %s", command)
		return nil, false
	}
	env := GetEnvironment(ctx)
	v, ok := env.Get(cmd)
	if !ok {
		log.Printf("Unknown variable in %s", command)
		return nil, false
	}
	return []interface{}{v, ctx}, true
}

func dynamicActionRun(data []interface{}) bool {
	v := data[0]
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	r, err := dvjson.ParseAny(v)
	if err != nil {
		log.Printf("Cannot parse dynamic action %v", err)
		return true
	}
	return DynamicActionByConfig(r, ctx)
}

func DynamicActionByConfig(config *dvevaluation.DvVariable, ctx *dvcontext.RequestContext) bool {
	if config == nil {
		return true
	}
	env := GetEnvironment(ctx)
	name := "production"
	if dvparser.IsDevelopment() {
		name = "development"
	}
	props, ok := config.FindChildByKey(name)
	if ok && props != nil {
		DynamicSetProperties(props)
	}
	props, ok = config.FindChildByKey("properties")
	if ok && props != nil {
		DynamicSetProperties(props)
	}
	props, ok = config.FindChildByKey("json")
	if ok && props != nil {
		DynamicSetProperties(props)
	}
	props, ok = config.FindChildByKey("actions")
	if ok && props != nil {
		DynamicSetActions(props, env)
	}
	return true
}

func DynamicSetProperties(props *dvevaluation.DvVariable) {
	fields := props.Fields
	n := len(fields)
	for i := 0; i < n; i++ {
		field := fields[i]
		if field != nil {
			dvparser.SetGlobalPropertiesAnyValue(string(field.Name), field)
		}
	}
}

func ConvertToAction(v *dvevaluation.DvVariable) *dvcontext.DvAction {
	if v == nil {
		return nil
	}
	action := &dvcontext.DvAction{
		Name:        v.ReadChildStringValue("name"),
		Typ:         v.ReadChildStringValue("type"),
		Url:         v.ReadChildStringValue("url"),
		Method:      v.ReadChildStringValue("method"),
		QueryParams: v.ReadChildMapValue("query"),
		Body:        v.ReadChildMapValue("body"),
		Result:      v.ReadChildStringValue("result"),
		ResultMode:  v.ReadChildStringValue("mode"),
		Definitions: v.ReadChildMapValue("definitions"),
		InnerParams: v.ReadChildStringValue("params"),
		Conditions:  v.ReadChildMapValue("conditions"),
		Validations: readValidations(v),
		Roles:       v.ReadChildStringValue("roles"),
		Auth:        v.ReadChildStringValue("auth"),
	}
	if len(action.Name) < 5 || len(action.Url) < 5 {
		log.Printf("Too small name or url %v", action)
	}
	return action
}

func readValidations(v *dvevaluation.DvVariable) []*dvcontext.ValidatePattern {
	if v == nil {
		return nil
	}
	v = v.ReadSimpleChild("validations")
	if v == nil || v.Kind != dvevaluation.FIELD_ARRAY || len(v.Fields) == 0 {
		return nil
	}
	n := len(v.Fields)
	res := make([]*dvcontext.ValidatePattern, n)
	for i := 0; i < n; i++ {
		val := v.Fields[i]
		if val == nil {
			continue
		}
		valid := &dvcontext.ValidatePattern{
			Source:       v.ReadChildStringValue("source"),
			Message:      v.ReadChildStringValue("message"),
			EmptyMessage: v.ReadChildStringValue("empty"),
			Contains:     v.ReadChildStringValue("contains"),
			RegPattern:   v.ReadChildStringValue("pattern"),
		}
		res[i] = valid
	}
	return res
}

func DynamicSetActions(props *dvevaluation.DvVariable, env *dvevaluation.DvObject) {
	fields := props.Fields
	n := len(fields)
	for i := 0; i < n; i++ {
		field := fields[i]
		if field != nil {
			action := ConvertToAction(field)
			if action != nil {
				dvmodules.AddDynamicAction(action, env)
			}
		}
	}
}
