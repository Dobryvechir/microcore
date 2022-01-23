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

type ObjectByArrayInfo struct {
	Src         *JsonRead `json:"src"`
	Dst         string    `json:"dst"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	KeyPolicy   int       `json:"key_policy"`
	ValuePolicy int       `json:"value_policy"`
}

type ObjectByObjectInfo struct {
	Src         *JsonRead `json:"src"`
	Dst         string    `json:"dst"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	KeyPolicy   int       `json:"key_policy"`
	ValuePolicy int       `json:"value_policy"`
}

type RemoveByKeysInfo struct {
	Src  *JsonRead `json:"src"`
	Dst  string    `json:"dst"`
	Keys []string  `json:"keys"`
}

type ReplaceTextInfo struct {
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	Rules string `json:"rules"`
}

type ConcatObjectInfo struct {
	Sources []string `json:"sources"`
	Dst     string   `json:"dst"`
}

type AssignInfo struct {
	Var       string `json:"var"`
	Condition string `json:"condition"`
	IfNotSet  string `json:"if_not_set"`
	IfEmpty   string `json:"if_empty"`
}

type VarTransformConfig struct {
	JsonParse       map[string]*JsonParseData        `json:"parse"`
	Read            map[string]*JsonRead             `json:"read"`
	ToInteger       map[string]string                `json:"to_integer"`
	Transform       map[string]string                `json:"transform"`
	Clone           map[string]string                `json:"clone"`
	DefaultString   map[string]string                `json:"default_string"`
	Assign          map[string]*AssignInfo           `json:"assign"`
	DefaultAny      map[string]string                `json:"default_any"`
	FindRegExpr     map[string]*PlaceRegExpression   `json:"find"`
	ReplaceRegExpr  map[string]*ReplaceRegExpression `json:"replace"`
	IncreaseVersion map[string]*IncreaseVersionInfo  `json:"increase_version"`
	RemoveVars      []string                         `json:"remove_vars"`
	CreateObject    map[string]map[string]string     `json:"create_object"`
	CreateArray     map[string][]string              `json:"create_array"`
	ConcatObjects   *ConcatObjectInfo                `json:"concat_objects"`
	ObjectByArray   *ObjectByArrayInfo               `json:"object_by_array"`
	ObjectByObject  *ObjectByObjectInfo              `json:"object_by_object"`
	RemoveByKeys    *RemoveByKeysInfo                `json:"remove_by_keys"`
	ReplaceText     *ReplaceTextInfo                 `json:"replace_text"`
	ErrorMessage    string                           `json:"error_message"`
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
			r, err := JsonExtract(v, ctx)
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
	if config.Assign != nil {
		for k, v := range config.Assign {
			if v != nil {
				name := v.Condition
				if name == "" {
					name = v.Var
				}
				r, ok := ReadActionResult(name, ctx)
				if !ok {
					if v.IfNotSet != "" {
						r, ok = ReadActionResult(v.IfNotSet, ctx)
					}
					if !ok && v.IfEmpty != "" {
						r, ok = ReadActionResult(v.IfEmpty, ctx)
						SaveActionResult(k, r, ctx)
						continue
					}
				}
				if ok {
					if v.IfEmpty != "" && !dvevaluation.AnyToBoolean(r) {
						b, ok := ReadActionResult(v.IfEmpty, ctx)
						if ok {
							r = b
						}
					} else if v.Condition != "" && v.Var != "" {
						b, ok := ReadActionResult(v.Var, ctx)
						if ok {
							r = b
						}
					}
					SaveActionResult(k, r, ctx)
				}
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
	if config.ReplaceText != nil {
		v, ok := ReadActionResult(config.ReplaceText.Src, ctx)
		if ok && v != nil {
			s := dvevaluation.AnyToString(v)
			mp, ok := ReadActionResult(config.ReplaceText.Rules, ctx)
			if ok && mp != nil {
				rules := dvevaluation.AnyToDvVariable(mp)
				s = dvjson.ReplaceTextByObjectMap(s, rules)
			}
			SaveActionResult(config.ReplaceText.Dst, s, ctx)
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
						if t == nil {
							t = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_NULL}
						}
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
	if config.ConcatObjects != nil {
		n := len(config.ConcatObjects.Sources)
		dst := &dvevaluation.DvVariable{
			Kind:   dvevaluation.FIELD_OBJECT,
			Fields: make([]*dvevaluation.DvVariable, 0, 128),
		}
		for i := 0; i < n; i++ {
			v, ok := ReadActionResult(config.ConcatObjects.Sources[i], ctx)
			if !ok || v == nil {
				continue
			}
			src := dvevaluation.AnyToDvVariable(v)
			dvjson.ConcatObjects(dst, src)
		}
		SaveActionResult(config.ConcatObjects.Dst, dst, ctx)
	}
	if config.ObjectByArray != nil {
		r, err := JsonExtract(config.ObjectByArray.Src, ctx)
		if err != nil {
			dvlog.PrintlnError("Error in expression " + config.ObjectByArray.Src.Var + ":" + err.Error())
			ActionExceptionByError(config.ErrorMessage, err, ctx)
			return true
		} else {
			d := dvevaluation.AnyToDvVariable(r)
			v, err := dvjson.CreateObjectByArray(d, config.ObjectByArray.Key, config.ObjectByArray.Value, env, config.ObjectByArray.KeyPolicy, config.ObjectByArray.ValuePolicy)
			if err != nil {
				ActionExceptionByError(config.ErrorMessage, err, ctx)
				return true
			} else {
				SaveActionResult(config.ObjectByArray.Dst, v, ctx)
			}
		}
	}
	if config.ObjectByObject != nil {
		r, err := JsonExtract(config.ObjectByObject.Src, ctx)
		if err != nil {
			dvlog.PrintlnError("Error in expression " + config.ObjectByObject.Src.Var + ":" + err.Error())
			ActionExceptionByError(config.ErrorMessage, err, ctx)
			return true
		} else {
			d := dvevaluation.AnyToDvVariable(r)
			v, err := dvjson.CreateObjectByArray(d, config.ObjectByObject.Key, config.ObjectByObject.Value, env, config.ObjectByObject.KeyPolicy, config.ObjectByObject.ValuePolicy)
			if err != nil {
				ActionExceptionByError(config.ErrorMessage, err, ctx)
				return true
			} else {
				SaveActionResult(config.ObjectByObject.Dst, v, ctx)
			}
		}
	}
	if config.RemoveByKeys != nil {
		r, err := JsonExtract(config.RemoveByKeys.Src, ctx)
		if err != nil {
			dvlog.PrintlnError("Error in expression " + config.RemoveByKeys.Src.Var + ":" + err.Error())
			ActionExceptionByError(config.ErrorMessage, err, ctx)
			return true
		} else {
			d := dvevaluation.AnyToDvVariable(r)
			v, err := dvjson.RemoveByKeys(d, config.RemoveByKeys.Keys, env)
			if err != nil {
				ActionExceptionByError(config.ErrorMessage, err, ctx)
				return true
			} else {
				SaveActionResult(config.RemoveByKeys.Dst, v, ctx)
			}
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
