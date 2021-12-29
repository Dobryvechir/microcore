/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjsmaster"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

type PlaceRegExpression struct {
	RegExpr  string `json:"reg-expr"`
	Group    string `json:"group"`
	DefValue string `json:"def-value"`
	Source   string `json:"source"`
	IsAll    bool   `json:"is-all"`
	Count    int    `json:"count"`
}

type ReplaceRegExpression struct {
	RegExpr     string `json:"reg-expr"`
	Replacement string `json:"replacement"`
	Source      string `json:"source"`
	Literal     bool   `json:"literal"`
}

type JsonParseData struct {
	Var        string `json:"string"`
	Evaluation int    `json:"evaluation"`
}

type IncreaseVersionInfo struct {
	Var        string `json:"var"`
	Limit      int    `json:"limit"`
	DefVersion string `json:"def_version"`
}

type VarTransformConfig struct {
	JsonParse       map[string]*JsonParseData        `json:"parse"`
	Read            map[string]*JsonRead             `json:"read"`
	ToInteger       map[string]string                `json:"to_integer"`
	Transform       map[string]string                `json:"transform"`
	Clone           map[string]string                `json:"clone"`
	DefaultString   map[string]string                `json:"default_string"`
	DefaultAny      map[string]string                `json:"default_any"`
	FindRegExpr     map[string]*PlaceRegExpression   `json:"find"`
	ReplaceRegExpr  map[string]*ReplaceRegExpression `json:"replace"`
	IncreaseVersion map[string]*IncreaseVersionInfo  `json:"increase_version"`
	RemoveVars      []string                         `json:"remove_vars"`
	CreateObject    map[string]map[string]string     `json:"create_object"`
	CreateArray     map[string][]string              `json:"create_array"`
}

func varTransformInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &VarTransformConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func varTransformRun(data []interface{}) bool {
	config := data[0].(*VarTransformConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return VarTransformRunByConfig(config, ctx)
}

func VarTransformRunByConfig(config *VarTransformConfig, ctx *dvcontext.RequestContext) bool {
	env := GetEnvironment(ctx)
	if config.JsonParse != nil {
		for k, v := range config.JsonParse {
			j, _ := ReadActionResult(v.Var, ctx)
			var err error = nil
			if v.Evaluation > 0 {
				s := dvevaluation.AnyToString(j)
				s, err = env.CalculateStringWithBrackets([]byte(s), v.Evaluation)
				if err == nil {
					j = s
				}
			}
			if err == nil {
				j, err = dvjson.ParseAny(j)
			}
			if err != nil {
				dvlog.PrintlnError("Error in expression " + k + ":" + err.Error())
			} else {
				SaveActionResult(k, j, ctx)
			}
		}
	}
	if config.Read != nil {
		for k, v := range config.Read {
			r, err := JsonExtract(v, env)
			if err != nil {
				dvlog.PrintlnError("Error in expression " + k + ":" + err.Error())
			} else {
				SaveActionResult(k, r, ctx)
			}
		}
	}
	if config.ToInteger != nil {
		for k, v := range config.ToInteger {
			r, ok := ReadActionResult(v, ctx)
			if !ok {
				if Log >= dvlog.LogWarning {
					dvlog.PrintlnError("Variable not found " + v + " to be stored in " + k)
				}
			} else {
				r, ok = dvevaluation.AnyToNumberInt(r)
				if !ok {
					if Log >= dvlog.LogWarning {
						dvlog.PrintlnError("Variable not integer " + v + " to be stored in " + k)
					}
				} else {
					SaveActionResult(k, r, ctx)
				}
			}
		}
	}
	if config.Transform != nil {
		for k, v := range config.Transform {
			r, err := env.EvaluateAnyTypeExpression(v)
			if err != nil {
				dvlog.PrintlnError("Error in expression " + v + ":" + err.Error())
			} else {
				SaveActionResult(k, r, ctx)
			}
		}
	}
	if config.Clone != nil {
		for k, v := range config.Clone {
			r, ok := ReadActionResult(v, ctx)
			if ok {
				d := dvevaluation.AnyToDvVariable(r)
				d = d.Clone()
				SaveActionResult(k, d, ctx)
			}
		}
	}
	if config.DefaultString != nil {
		for k, v := range config.DefaultString {
			_, ok := ReadActionResult(k, ctx)
			if !ok {
				SaveActionResult(k, v, ctx)
			}
		}
	}
	if config.DefaultAny != nil {
		for k, v := range config.DefaultAny {
			_, ok := ReadActionResult(k, ctx)
			if !ok {
				r, err := dvjson.JsonFullParser([]byte(v))
				if err != nil {
					dvlog.PrintlnError("Error in json " + v + ":" + err.Error())
				}
				SaveActionResult(k, r, ctx)
			}
		}
	}
	if config.FindRegExpr != nil {
		for k, v := range config.FindRegExpr {
			r, ok := ReadActionResult(v.Source, ctx)
			if ok {
				src := dvevaluation.AnyToString(r)
				res, err := dvtextutils.FindByRegularExpression(src, v.RegExpr, v.Group, v.DefValue, v.IsAll, v.Count)
				if err != nil {
					dvlog.PrintlnError("Error in regular expression " + src + ":" + err.Error())
				}
				SaveActionResult(k, res, ctx)
			}
		}
	}
	if config.ReplaceRegExpr != nil {
		for k, v := range config.ReplaceRegExpr {
			r, ok := ReadActionResult(v.Source, ctx)
			if ok {
				src := dvevaluation.AnyToString(r)
				res, err := dvtextutils.ReplaceByRegularExpression(src, v.RegExpr, v.Replacement, v.Literal)
				if err != nil {
					dvlog.PrintlnError("Error in regular expression " + src + ":" + err.Error())
				}
				SaveActionResult(k, res, ctx)
			}
		}
	}
	if config.IncreaseVersion != nil {
		for k, v := range config.IncreaseVersion {
			r, ok := ReadActionResult(v.Var, ctx)
			if !ok {
				dvlog.PrintlnError("Version variable not found " + v.Var + " to be stored in " + k)
			} else {
				s := dvevaluation.AnyToString(r)
				s = dvjsmaster.MathIncreaseVersion(s, v.Limit, v.DefVersion)
				SaveActionResult(k, s, ctx)
			}
		}
	}
	if config.RemoveVars != nil {
		for _, v := range config.RemoveVars {
			DeleteActionResult(v, ctx)
		}
	}
	if config.CreateObject != nil {
		for kp, vp := range config.CreateObject {
			f := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, 32)}
			if vp != nil {
				for k, v := range vp {
					vs, ok := ReadActionResult(v, ctx)
					if ok {
						t := dvevaluation.AnyToDvVariable(vs)
						t.Name = []byte(k)
						f.Fields = append(f.Fields, t)
					} else {
						f.Fields = append(f.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_NULL, Name: []byte(k)})
					}
				}
			}
			SaveActionResult(kp, f, ctx)
		}
	}
	if config.CreateArray != nil {
		for kp, vp := range config.CreateArray {
			f := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, 32)}
			if vp != nil {
				for _, v := range vp {
					vs, ok := ReadActionResult(v, ctx)
					if ok {
						t := dvevaluation.AnyToDvVariable(vs)
						f.Fields = append(f.Fields, t)
					} else {
						f.Fields = append(f.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_NULL})
					}
				}
			}
			SaveActionResult(kp, f, ctx)
		}
	}
	return true
}
