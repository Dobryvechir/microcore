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
	"strings"
)

type FilterOutBlock struct {
	Condition string   `json:"condition"`
	Set       []string `json:"set"`
	Unset     []string `json:"unset"`
	NoEmpty   bool     `json:"no_empty"`
}

type ConvertReplaceBlock struct {
	FindObject    string `json:"find_object"`
	FindCondition string `json:"find_condition"`
	IsLast        bool   `json:"is_last"`
	ReplaceWith   string `json:"replace_with"`
	ReplaceMode   string `json:"replace_mode"`
	InDefault     string `json:"in_default"`
}

type JsonConvertConfig struct {
	Source       *JsonRead            `json:"source"`
	Result       string               `json:"result"`
	Add          []*JsonRead          `json:"add"`
	Merge        []*JsonRead          `json:"merge"`
	Replace      []*JsonRead          `json:"replace"`
	Update       []*JsonRead          `json:"update"`
	Push         []*JsonRead          `json:"push"`
	Concat       []*JsonRead          `json:"concat"`
	Remove       []string             `json:"remove"`
	ForEach      *ForEachBlock        `json:"for_each"`
	Filter       *FilterOutBlock      `json:"filter"`
	ReplaceBlock *ConvertReplaceBlock `json:"replace_block"`
}

func jsonConvertInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &JsonConvertConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Source == nil || config.Source.Var == "" {
		log.Printf("source must be present in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func jsonConvertRun(data []interface{}) bool {
	config := data[0].(*JsonConvertConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return JsonConvertRunByConfig(config, ctx)
}

func updateVariablesByConfig(config []*JsonRead, mode int, src interface{}, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) (interface{}, bool) {
	n := len(config)
	for i := 0; i < n; i++ {
		conf := config[i]
		v, err := JsonExtract(conf, ctx)
		if err != nil {
			dvlog.PrintlnError("Error in json extracting by " + conf.Var)
			return src, false
		}
		src = dvevaluation.UpdateAnyVariables(src, v, conf.Destination,
			mode, conf.Ids, env)
	}
	return src, true
}

func JsonConvertRunByConfig(config *JsonConvertConfig, ctx *dvcontext.RequestContext) bool {
	env := GetEnvironment(ctx)
	src, err := JsonExtract(config.Source, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + config.Source.Var)
		return true
	}
	var ok bool
	src, ok = updateVariablesByConfig(config.Add, dvevaluation.UPDATE_MODE_ADD_BY_KEYS, src, env, ctx)
	if !ok {
		return true
	}
	src, ok = updateVariablesByConfig(config.Merge, dvevaluation.UPDATE_MODE_MERGE, src, env, ctx)
	if !ok {
		return true
	}
	src, ok = updateVariablesByConfig(config.Replace, dvevaluation.UPDATE_MODE_REPLACE, src, env, ctx)
	if !ok {
		return true
	}
	src, ok = updateVariablesByConfig(config.Update, dvevaluation.UPDATE_MODE_APPEND, src, env, ctx)
	if !ok {
		return true
	}
	n := len(config.Remove)
	for i := 0; i < n; i++ {
		v := config.Remove[i]
		src = dvevaluation.RemoveAnyVariable(src, v, env)
	}
	s := dvevaluation.AnyToDvVariable(src)
	n = len(config.Push)
	for i := 0; i < n; i++ {
		JsonConvertPush(config.Push[i], s, ctx)
	}
	n = len(config.Concat)
	for i := 0; i < n; i++ {
		JsonConvertConcat(config.Concat[i], s, ctx)
	}
	FilterOutByConditionSetUnset(env, s, config.Filter)
	if config.ForEach != nil {
		config.ForEach.ForEachProcessing(s, env, ctx)
	}
	if config.ReplaceBlock != nil {
		ok = replaceBlockConvert(config.ReplaceBlock, s, env, ctx)
		if !ok {
			return true
		}
	}
	SaveActionResult(config.Result, s, ctx)
	return true
}

func FilterOutByConditionSetUnset(env *dvevaluation.DvObject, src *dvevaluation.DvVariable, filter *FilterOutBlock) {
	if filter != nil && src != nil && src.Fields != nil {
		n := len(src.Fields)
		mu := len(filter.Unset)
		ms := len(filter.Set)
		mst := ms
		checkSet := ms > 0
		checkUnset := mu > 0
		checkFields := checkUnset || checkSet
		var unsetMap map[string]int
		var setMap map[string]int
		if checkUnset {
			unsetMap = make(map[string]int, mu)
			for j := 0; j < mu; j++ {
				unsetMap[filter.Unset[j]] = 1
			}
		}
		if checkSet {
			setMap = make(map[string]int, ms)
			for j := 0; j < ms; j++ {
				h := filter.Set[j]
				if _, ok := setMap[h]; ok {
					mst--
				} else {
					setMap[h] = 1
				}
			}
		}
		condition := filter.Condition
		for i := 0; i < n; i++ {
			f := src.Fields[i]
			rm := false
			if f == nil || f.Kind == dvevaluation.FIELD_NULL || f.Kind == dvevaluation.FIELD_UNDEFINED {
				rm = filter.NoEmpty
			} else {
				fields := dvjson.CreateLocalVariables(env, src.Fields[i])
				if checkFields {
					fn := len(fields)
					cnt := mst
					for j := 1; j < fn; j++ {
						s := fields[j]
						if checkUnset {
							if _, ok := unsetMap[s]; ok {
								rm = true
								break
							}
						}
						if checkSet {
							if _, ok := setMap[s]; ok {
								cnt--
							}
						}
					}
					if !rm && cnt > 0 {
						rm = true
					}
				}
				if !rm && condition != "" {
					c, err := env.EvaluateBooleanExpression(condition)
					if err != nil {
						dvlog.PrintfError("Error in expression %s %v", condition, err)
						break
					}
					rm = !c
				}
				dvjson.RemoveLocalVariables(env, fields)
			}
			if rm {
				if i == n-1 {
					src.Fields = src.Fields[:i]
				} else {
					src.Fields = append(src.Fields[:i], src.Fields[i+1:]...)
				}
				n--
				i--
			}
		}
	}
}

func JsonConvertPush(push *JsonRead, dst *dvevaluation.DvVariable, ctx *dvcontext.RequestContext) {
	src, err := JsonExtract(push, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + push.Var)
		return
	}
	if src != nil && dst != nil && dst.Kind == dvevaluation.FIELD_ARRAY {
		s := dvevaluation.AnyToDvVariable(src)
		dst.Fields = append(dst.Fields, s)
	}
}

func JsonConvertConcat(push *JsonRead, dst *dvevaluation.DvVariable, ctx *dvcontext.RequestContext) {
	src, err := JsonExtract(push, ctx)
	if err != nil {
		dvlog.PrintlnError("Error in json extracting by " + push.Var)
		return
	}
	if src != nil && dst != nil && dst.Kind == dvevaluation.FIELD_ARRAY {
		s := dvevaluation.AnyToDvVariable(src)
		if s != nil && s.Kind == dvevaluation.FIELD_ARRAY {
			n := len(s.Fields)
			for i := 0; i < n; i++ {
				dst.Fields = append(dst.Fields, s.Fields[i])
			}
		}
	}
}

func findObjectOrCondition(obj string, condition string, isLast bool, v *dvevaluation.DvVariable, ctx *dvcontext.RequestContext, env *dvevaluation.DvObject) int {
	if v == nil {
		return -1
	}
	n := len(v.Fields)
	if n == 0 {
		return -1
	}
	if obj != "" {
		d, ok := ReadActionResult(obj, ctx)
		if !ok || d == nil {
			return -1
		}
		dv := dvevaluation.AnyToDvVariable(d)
		if dv == nil {
			return -1
		}
		for i := 0; i < n; i++ {
			if v.Fields[i] == dv {
				return i
			}
		}
		return -1
	}
	if condition == "" {
		return -1
	}
	res := -1
	for i := 0; i < n; i++ {
		fields := dvjson.CreateLocalVariables(env, v.Fields[i])
		r, err := env.EvaluateBooleanExpression(condition)
		dvjson.RemoveLocalVariables(env, fields)
		if err == nil && r {
			res = i
			if !isLast {
				break
			}
		}
	}
	return res
}

func replaceBlockConvert(block *ConvertReplaceBlock, v *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) bool {
	ind := findObjectOrCondition(block.FindObject, block.FindCondition, block.IsLast, v, ctx, env)
	if ind < 0 {
		return inDefaultReplace(block.InDefault, block.ReplaceWith, v, env, ctx)
	}
	return replaceByMode(block.ReplaceWith, block.ReplaceMode, v, env, ind)
}

func inDefaultReplace(inDefault string, replaceWith string, v *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) bool {
	var atEnd bool
	switch strings.ToLower(strings.TrimSpace(inDefault)) {
	case "ignore":
		return true
	case "push":
		atEnd = true
	case "unshift":
		atEnd = false
	default:
		dvlog.PrintlnError("Not resolved replacement " + replaceWith)
		ctx.HandleInternalServerError()
		return false
	}
	r, ok := env.Get(replaceWith)
	if !ok || v == nil {
		return true
	}
	d := dvevaluation.AnyToDvVariable(r)
	if atEnd || len(v.Fields) == 0 {
		v.Fields = append(v.Fields, d)
	} else {
		n := len(v.Fields)
		v.Fields = append(v.Fields, d)
		m := v.Fields[0]
		v.Fields[0] = d
		for i := 1; i <= n; i++ {
			d = v.Fields[i]
			v.Fields[i] = m
			m = d
		}
	}
	return true
}

func replaceByMode(replaceWith string, replaceMode string, v *dvevaluation.DvVariable, env *dvevaluation.DvObject, ind int) bool {
	r, ok := env.Get(replaceWith)
	if !ok || v == nil {
		return true
	}
	d := dvevaluation.AnyToDvVariable(r)
	switch strings.ToLower(strings.TrimSpace(replaceMode)) {
	case "at":
		v.Fields[ind] = d
		return true
	case "before":
	default:
		ind++
	}
	n := len(v.Fields)
	if ind >= n {
		v.Fields = append(v.Fields, d)
		return true
	}
	v.Fields = append(v.Fields, d)
	m := v.Fields[ind]
	v.Fields[ind] = d
	for ind++; ind <= n; ind++ {
		d = v.Fields[ind]
		v.Fields[ind] = m
		m = d
	}
	return true
}
