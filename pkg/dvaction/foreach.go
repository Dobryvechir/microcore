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
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strconv"
	"strings"
)

type AssignmentBlock struct {
	Source string `json:"source"`
	Field  string `json:"field"`
	Value  string `json:"value"`
}

type ConditionBlock struct {
	Condition       string             `json:"condition"`
	SetList         []string           `json:"set"`
	UnsetList       []string           `json:"unset"`
	ThenAssign      []*AssignmentBlock `json:"then_assign"`
	ElseAssign      []*AssignmentBlock `json:"else_assign"`
	ThenWholeAssign string             `json:"then_to_whole"`
	ElseWholeAssign string             `json:"else_to_whole"`
	ThenCollection  *Collection        `json:"then_collection"`
	ElseCollection  *Collection        `json:"else_collection"`
	Match           string             `json:"match"`
	Item            string             `json:"item"`
	setMap          map[string]int
	unsetMap        map[string]int
}

type ForEachBlock struct {
	PreCondition string            `json:"pre_condition"`
	WildCardPath string            `json:"path"`
	Blocks       []*ConditionBlock `json:"blocks"`
}

const (
	ACT_NORMAL = iota
	ACT_BREAK
	ACT_REMOVE
	ACT_REPLACE
	ACT_REPLACE_BREAK
)

func preinitializeCollections(blocks []*ConditionBlock, env *dvevaluation.DvObject) {
	n := len(blocks)
	for i := 0; i < n; i++ {
		b := blocks[i]
		if b != nil {
			if b.ThenCollection != nil {
				b.ThenCollection.Initialize(env)
			}
			if b.ElseCollection != nil {
				b.ElseCollection.Initialize(env)
			}
		}
	}
}

func (proc *ForEachBlock) ForEachProcessing(src *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) {
	if proc == nil || src == nil {
		return
	}
	preinitializeCollections(proc.Blocks, env)
	if proc.WildCardPath == "" {
		proc.ForEachProcessingWithoutPath(src, env, ctx)
		return
	}
	pathes := dvtextutils.SeparateChildExpression(proc.WildCardPath)
	forEachPathKeys := &dvevaluation.DvVariable{
		Kind:   dvevaluation.FIELD_ARRAY,
		Fields: make([]*dvevaluation.DvVariable, 0, 16),
	}
	env.Set("FOR_EACH_PATH_KEYS", forEachPathKeys)
	forEachPathValues := &dvevaluation.DvVariable{
		Kind:   dvevaluation.FIELD_ARRAY,
		Fields: make([]*dvevaluation.DvVariable, 0, 16),
	}
	env.Set("FOR_EACH_PATH_VALUES", forEachPathValues)
	proc.ForEachProcessingWithPath(forEachPathKeys, forEachPathValues, pathes, src, env, ctx)
}

func (proc *ForEachBlock) ForEachProcessingWithPath(forEachPathKeys *dvevaluation.DvVariable, forEachPathValues *dvevaluation.DvVariable, pathes []string, src *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) {
	if len(pathes) == 0 {
		proc.ForEachProcessingWithoutPath(src, env, ctx)
		return
	}
	path := strings.TrimSpace(pathes[0])
	rest := pathes[1:]
	if path == "" {
		proc.ForEachProcessingWithPath(forEachPathKeys, forEachPathValues, rest, src, env, ctx)
		return
	}
	if path == "*" {
		n := len(src.Fields)
		forEachPathKeys.Fields = append(forEachPathKeys.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING})
		forEachPathValues.Fields = append(forEachPathValues.Fields, nil)
		m := len(forEachPathKeys.Fields) - 1
		for i := 0; i < n; i++ {
			f := src.Fields[i]
			if f != nil {
				if src.Kind == dvevaluation.FIELD_OBJECT {
					forEachPathKeys.Fields[m].Value = f.Name
				} else {
					forEachPathKeys.Fields[m].Value = []byte(strconv.Itoa(i))
				}
				forEachPathValues.Fields[m] = f
				proc.ForEachProcessingWithPath(forEachPathKeys, forEachPathValues, rest, f, env, ctx)
			}
		}
		forEachPathKeys.Fields = forEachPathKeys.Fields[:m]
		forEachPathValues.Fields = forEachPathValues.Fields[:m]
	} else if path != "" && path[0] == '{' && path[len(path)-1] == '}' {
		n := len(src.Fields)
		forEachPathKeys.Fields = append(forEachPathKeys.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING})
		forEachPathValues.Fields = append(forEachPathValues.Fields, nil)
		m := len(forEachPathKeys.Fields) - 1
		path = path[1 : len(path)-1]
		for i := 0; i < n; i++ {
			f := src.Fields[i]
			if f != nil {
				if applyFilter(f, i, env, path) {
					if src.Kind == dvevaluation.FIELD_OBJECT {
						forEachPathKeys.Fields[m].Value = f.Name
					} else {
						forEachPathKeys.Fields[m].Value = []byte(strconv.Itoa(i))
					}
					forEachPathValues.Fields[m] = f
					proc.ForEachProcessingWithPath(forEachPathKeys, forEachPathValues, rest, f, env, ctx)
				}
			}
		}
		forEachPathKeys.Fields = forEachPathKeys.Fields[:m]
		forEachPathValues.Fields = forEachPathValues.Fields[:m]
	} else {
		f, _, err := src.ReadPath(path, false, env)
		if err != nil {
			dvlog.PrintfError("Error forEach %v", err)
			return
		}
		if f != nil {
			forEachPathKeys.Fields = append(forEachPathKeys.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING})
			forEachPathValues.Fields = append(forEachPathValues.Fields, nil)
			m := len(forEachPathKeys.Fields) - 1
			if src.Kind == dvevaluation.FIELD_OBJECT {
				forEachPathKeys.Fields[m].Value = f.Name
			} else {
				i := src.IndexOf(f)
				forEachPathKeys.Fields[m].Value = []byte(strconv.Itoa(i))
			}
			forEachPathValues.Fields[m] = f
			proc.ForEachProcessingWithPath(forEachPathKeys, forEachPathValues, rest, f, env, ctx)
			forEachPathKeys.Fields = forEachPathKeys.Fields[:m]
			forEachPathValues.Fields = forEachPathValues.Fields[:m]
		}
	}
}

func (proc *ForEachBlock) ForEachProcessingWithoutPath(src *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) {
	env.Set("_root", src)
	if proc.PreCondition != "" {
		b, err := env.EvaluateBooleanExpression(proc.PreCondition)
		if err != nil {
			dvlog.PrintfError("Error in evaluation %s : %v", proc.PreCondition, err)
			return
		}
		if !b {
			return
		}
	}
	n := len(proc.Blocks)
	m := len(src.Fields)
	for i := 0; i < m; i++ {
		env.Set("_index", i)
		f := src.Fields[i]
		if f == nil {
			continue
		}
		removeAct := false
		fields := dvjson.CreateLocalVariables(env, f)
		for j := 0; j < n; j++ {
			b := proc.Blocks[j]
			if b != nil {
				act := b.Process(f, env, ctx, fields)
				if act == ACT_BREAK {
					break
				}
				if act == ACT_REMOVE {
					removeAct = true
					break
				}
				if act == ACT_REPLACE || act == ACT_REPLACE_BREAK {
					v, ok := env.Get("this")
					if !ok {
						f = nil
					} else {
						f = dvevaluation.AnyToDvVariable(v)
					}
					src.Fields[i] = f
					if act == ACT_REPLACE_BREAK {
						break
					}
				}
			}
		}
		if removeAct {
			if i == m-1 {
				src.Fields = src.Fields[:i]
			} else {
				src.Fields = append(src.Fields[:i], src.Fields[i+1:]...)
			}
			i--
			m--
		}
		dvjson.RemoveLocalVariables(env, fields)
	}
}

func (b *ConditionBlock) Process(f *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext, fields []string) int {
	if b.Match == "" {
		return b.ProcessSingle(f, env, ctx, fields)
	}
	if b.Item == "" {
		b.Item = "_item"
	}
	indexName := b.Item + "_index"
	d, ok := env.Get(b.Match)
	if !ok || d == nil {
		return ACT_NORMAL
	}
	v := dvevaluation.AnyToDvVariable(d)
	if v == nil || v.Fields == nil {
		return ACT_NORMAL
	}
	n := len(v.Fields)
	for i := 0; i < n; i++ {
		env.Set(b.Item, v.Fields[i])
		env.Set(indexName, i)
		p := b.ProcessSingle(f, env, ctx, fields)
		if p != ACT_NORMAL {
			return p
		}
	}
	return ACT_NORMAL
}

func (b *ConditionBlock) ProcessSingle(f *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext, fields []string) int {
	success := b.EvaluateSuccess(env, ctx, fields)
	if success {
		AssignVariables(f, env, ctx, b.ThenAssign)
		if b.ThenCollection != nil {
			b.ThenCollection.AddItemSecondary(env)
		}
		return ActCalculation(env, ctx, b.ThenWholeAssign)
	}
	AssignVariables(f, env, ctx, b.ElseAssign)
	if b.ElseCollection != nil {
		b.ElseCollection.AddItemSecondary(env)
	}
	return ActCalculation(env, ctx, b.ElseWholeAssign)
}

func (b *ConditionBlock) EvaluateSuccess(env *dvevaluation.DvObject, ctx *dvcontext.RequestContext, fields []string) bool {
	ns := len(b.SetList)
	nu := len(b.UnsetList)
	useSet := ns > 0
	useUnset := nu > 0
	if useSet && b.setMap == nil {
		b.setMap = make(map[string]int, ns)
		for i := 0; i < ns; i++ {
			s := b.SetList[i]
			if _, ok := b.setMap[s]; ok {
				if i == ns-1 {
					b.SetList = b.SetList[:i]
				} else {
					b.SetList = append(b.SetList[:i], b.SetList[i+1:]...)
				}
				i--
				ns--
			} else {
				b.setMap[s] = 1
			}
		}
	}
	if useUnset && b.unsetMap == nil {
		b.unsetMap = make(map[string]int, nu)
		for i := 0; i < nu; i++ {
			s := b.UnsetList[i]
			if _, ok := b.unsetMap[s]; ok {
				if i == nu-1 {
					b.UnsetList = b.UnsetList[:i]
				} else {
					b.UnsetList = append(b.UnsetList[:i], b.UnsetList[i+1:]...)
				}
				i--
				nu--
			} else {
				b.unsetMap[s] = 1
			}
		}
	}
	if useUnset || useSet {
		n := len(fields)
		setCount := ns
		for i := 1; i < n; i++ {
			fld := fields[i]
			if useUnset {
				if _, ok := b.unsetMap[fld]; ok {
					return false
				}
			}
			if useSet {
				if _, ok := b.setMap[fld]; ok {
					setCount--
				}
			}
		}
		if setCount > 0 {
			return false
		}
	}
	if b.Condition != "" {
		r, err := env.EvaluateBooleanExpression(b.Condition)
		if err != nil {
			dvlog.PrintfError("Error in expr %s: %v", b.Condition, err)
		}
		return r
	}
	return true
}

func AssignVariables(f *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext, assigns []*AssignmentBlock) {
	n := len(assigns)
	var err error
	var d *dvevaluation.DvVariable
	for i := 0; i < n; i++ {
		a := assigns[i]
		if a.Source == "" {
			err = f.AssignToSubField(a.Field, a.Value, env)
		} else {
			v, ok := ReadActionResult(a.Source, ctx)
			if !ok || v == nil {
				d = nil
			} else {
				d = dvevaluation.AnyToDvVariable(v)
			}
			if d == nil {
				d = &dvevaluation.DvVariable{
					Kind:   dvevaluation.FIELD_OBJECT,
					Fields: make([]*dvevaluation.DvVariable, 0, 7),
				}
			}
			err = d.AssignToSubField(a.Field, a.Value, env)
			SaveActionResult(a.Source, d, ctx)
		}
		if err != nil {
			dvlog.PrintfError("Error %s : %s: %v", a.Field, a.Value, err)
		}
	}
}

func ActCalculation(env *dvevaluation.DvObject, ctx *dvcontext.RequestContext, wholeAssign string) int {
	switch wholeAssign {
	case "":
		return ACT_NORMAL
	case "break":
		return ACT_BREAK
	case "delete":
		return ACT_REMOVE
	}
	act := ACT_REPLACE
	c := ";break"
	if strings.HasSuffix(wholeAssign, c) {
		wholeAssign = wholeAssign[:len(wholeAssign)-len(c)]
		act = ACT_REPLACE_BREAK
	}
	v, err := env.EvaluateAnyTypeExpression(wholeAssign)
	if err != nil {
		dvlog.PrintfError("Error in expression %s: %v", wholeAssign, err)
		return ACT_NORMAL
	}
	env.Set("this", v)
	return act
}

func applyFilter(f *dvevaluation.DvVariable, ind int, env *dvevaluation.DvObject, filter string) bool {
	env.Set("__index", ind)
	nm := string(f.Name)
	if nm != "" {
		env.Set("__key", nm)
	}
	fields := dvjson.CreateLocalVariables(env, f)
	r, err := env.EvaluateBooleanExpression(filter)
	dvjson.RemoveLocalVariables(env, fields)
	if err == nil && r {
		return true
	}
	return false
}
