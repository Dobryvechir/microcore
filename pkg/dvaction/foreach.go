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
	"strings"
)

type AssignmentBlock struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type ConditionBlock struct {
	Condition          string             `json:"condition"`
	SetList            []string           `json:"set"`
	UnsetList          []string           `json:"unset"`
	AssignSuccess      []*AssignmentBlock `json:"assign_success"`
	AssignFailure      []*AssignmentBlock `json:"assign_failure"`
	WholeAssignSuccess string             `json:"whole_success"`
	WholeAssignFailure string             `json:"whole_failure"`
	setMap             map[string]int
	unsetMap           map[string]int
}

type ForEachBlock struct {
	PreCondition string            `json:"pre_condition"`
	WildCardPath string            `json:"wild_card_path"`
	Blocks       []*ConditionBlock `json:"blocks"`
}

const (
	ACT_NORMAL = iota
	ACT_BREAK
	ACT_REMOVE
	ACT_REPLACE
	ACT_REPLACE_BREAK
)

func (proc *ForEachBlock) ForEachProcessing(src *dvevaluation.DvVariable, env *dvevaluation.DvObject, ctx *dvcontext.RequestContext) {
	if proc == nil || src == nil {
		return
	}
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
	success := b.EvaluateSuccess(env, ctx, fields)
	if success {
		AssignVariables(f, env, ctx, b.AssignSuccess)
		return ActCalculation(env, ctx, b.WholeAssignSuccess)
	}
	AssignVariables(f, env, ctx, b.AssignFailure)
	return ActCalculation(env, ctx, b.WholeAssignFailure)
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
	for i := 0; i < n; i++ {
		a := assigns[i]
		err := f.AssignToSubField(a.Field, a.Value, env)
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
