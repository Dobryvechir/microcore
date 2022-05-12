// package dvoc orchestrates actions, executions
// MicroCore Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvsecurity"
	"github.com/Dobryvechir/microcore/pkg/dvsession"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	ActionPrefix = "ACTION_"
)

var Log = dvlog.LogError

func fireAction(ctx *dvcontext.RequestContext) bool {
	return fireActionByName(ctx, ctx.Action.Name, ctx.Action.Definitions, false)
}

func fireActionByName(ctx *dvcontext.RequestContext, name string,
	definitions map[string]string, omitResults bool) bool {
	if ctx.Action != nil && ctx.Action.ErrorPolicy != "" {
		ctx.PrimaryContextEnvironment.Set("ERROR_POLICY", ctx.Action.ErrorPolicy)
	}
	prefix := ActionPrefix + name
	if ctx.PrimaryContextEnvironment.GetString(prefix+"_1") == "" {
		ctx.StatusCode = 501
		ctx.HandleCommunication()
		return true
	}
	res := ExecuteSequence(prefix, ctx, definitions)
	if !omitResults {
		ActionProcessResult(ctx, res)
	}
	return true
}

func ActionProcessResult(ctx *dvcontext.RequestContext, res bool) {
	if !res && ctx.StatusCode < 400 {
		ctx.HandleInternalServerError()
	} else {
		ActionContextResult(ctx)
	}
}

func fireStaticAction(ctx *dvcontext.RequestContext) bool {
	ActionProcessResult(ctx, true)
	return true
}

func ActionContextResult(ctx *dvcontext.RequestContext) {
	if ctx.StatusCode >= 400 {
		ctx.HandleCommunication()
		return
	}
	action := ctx.Action
	if action != nil && action.Result != "" {
		res, err := ctx.PrimaryContextEnvironment.CalculateString(action.Result)
		if err != nil {
			ctx.Error = err
			ctx.HandleInternalServerError()
			return
		}
		switch action.ResultMode {
		case "file":
			ctx.Output, err = GetContextFileResult(ctx, res)
		case "var":
			ctx.Output, err = GetContextVarResult(ctx, res)
		default:
			ctx.Output = []byte(res)
		}
		if err != nil {
			ctx.Error = err
			ctx.HandleInternalServerError()
			return
		}
		ctx.Output = []byte(res)
	}
	setHeadersName := ctx.PrimaryContextEnvironment.GetString(dvcontext.AUTO_HEADER_SET_BY)
	if setHeadersName != "" {
		provideResponseHeaders(setHeadersName+"_", ctx)
	}
	ctx.HandleCommunication()
}

func provideResponseHeaders(pref string, ctx *dvcontext.RequestContext) {
	if ctx.Headers == nil {
		ctx.Headers = make(map[string][]string)
	}
	n := len(pref)
	var res []string
	for k, v := range ctx.PrimaryContextEnvironment.Properties {
		if strings.HasPrefix(k, pref) {
			key := k[n:]
			if key != "" {
				switch v.(type) {
				case string:
					res = []string{v.(string)}
				case []string:
					res = v.([]string)
				default:
					sv := dvevaluation.AnyToString(v)
					if sv == "" {
						res = nil
					} else {
						res = []string{sv}
					}
				}
				if res != nil {
					ctx.Headers[key] = res
				}
			}
		}
	}
}

func GetContextVarResult(ctx *dvcontext.RequestContext, varName string) ([]byte, error) {
	dat, ok := ctx.PrimaryContextEnvironment.Get(varName)
	if !ok {
		return nil, errors.New("Variable " + varName + " not set")
	}
	str := dvevaluation.AnyToString(dat)
	return []byte(str), nil
}

func GetContextFileResult(ctx *dvcontext.RequestContext, fileName string) ([]byte, error) {
	if fileName == "" {
		ctx.HandleFileNotFound()
		return nil, errors.New("Empty file name")
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Cannot read %s: %v", fileName, err)
		ctx.HandleInternalServerError()
		return nil, errors.New("File " + fileName + " not found")
	}
	if !bytes.Contains(data, []byte("{{")) {
		return data, nil
	}
	res, err := ctx.PrimaryContextEnvironment.CalculateString(string(data))
	return []byte(res), err
}

func fireSwitchAction(ctx *dvcontext.RequestContext) bool {
	action := ctx.Action
	actionName := action.Result
	conditions := action.Conditions
	if nil != conditions {
		for k, v := range conditions {
			res, err := ctx.PrimaryContextEnvironment.EvaluateBooleanExpression(k)
			if err != nil {
				log.Printf("Failed to evaluate %s: %v", k, err)
				ctx.HandleInternalServerError()
				return true
			}
			if res {
				actionName = v
				break
			}
		}
	}
	return fireActionByName(ctx, actionName, action.Definitions, false)
}

func securityEndPointHandler(ctx *dvcontext.RequestContext) bool {
	res := dvsecurity.LoginByRequestEndPointHandler(ctx)
	if res {
		ActionContextResult(ctx)
	}
	return res
}

func GetEnvironment(ctx *dvcontext.RequestContext) *dvevaluation.DvObject {
	if ctx == nil || ctx.PrimaryContextEnvironment == nil {
		return dvparser.GetGlobalPropertiesAsDvObject()
	}
	if ctx.LocalContextEnvironment != nil {
		return ctx.LocalContextEnvironment
	}
	return ctx.PrimaryContextEnvironment
}

func SaveActionResult(result string, data interface{}, ctx *dvcontext.RequestContext) {
	env := GetEnvironment(ctx)
	if result != "" {
		level, varName, path := GetLevelMainPath(result)
		if level == "log" || level == "error" {
			s := result[strings.Index(result, ":")+1:]
			v := dvevaluation.AnyToString(data)
			log.Printf(s, v)
			if level == "error" {
				ActionInternalException(0, s, v, ctx)
			}
			return
		}
		if path != "" {
			name := varName
			if level != "" {
				name = level + ":" + name
			}
			dat, ok := ReadActionResult(name, ctx)
			if !ok {
				data = dvevaluation.CreateDvVariableByPathAndData(path, data, nil)
			} else {
				dvevaluation.UpdateAnyVariables(dat, data, path, dvevaluation.UPDATE_MODE_REPLACE, nil, env)
				return
			}
		}
		isMap := strings.HasPrefix(level, "map_")
		if level == "session" || level == "session?" {
			isErrorFatal := level == "session"
			if ctx.Session == nil {
				if isErrorFatal {
					dvlog.PrintlnError("No session available")
					ctx.HandleInternalServerError()
				}
				return
			}
			ctx.Session.SetItem(varName, data)
		} else if ctx != nil && level != "global" && !isMap {
			if ctx.LocalContextEnvironment != nil && level != "request" {
				if level != "" && level[0] >= '1' && level[0] <= '9' {
					levelVal, err := strconv.Atoi(level)
					if err == nil {
						ctx.LocalContextEnvironment.SetAtParent(varName, data, levelVal)
					} else {
						dvlog.PrintfError("Unknown level %s", level)
					}
				} else {
					ctx.LocalContextEnvironment.Set(varName, data)
				}
			} else {
				ctx.PrimaryContextEnvironment.Set(varName, data)
			}
		} else {
			if isMap {
				mapName := level[4:]
				dvsession.GlobalMapWrite(mapName, varName, data)
			} else {
				dvparser.SetGlobalPropertiesAnyValue(varName, data)
			}
		}
	}
}

func DeleteActionResult(result string, ctx *dvcontext.RequestContext) {
	if result != "" {
		level, varName, path := GetLevelMainPath(result)
		var res interface{}
		var ok bool
		var env *dvevaluation.DvObject = nil
		isSession := false
		isMap := strings.HasPrefix(level, "map_")
		if ctx != nil && level != "global" && !isMap {
			if level == "session" || level == "session?" {
				isErrorFatal := level == "session"
				if ctx.Session == nil {
					if isErrorFatal {
						dvlog.PrintlnError("No session available")
						ctx.HandleInternalServerError()
					}
					return
				}
				if path == "" {
					ctx.Session.RemoveItem(varName)
				} else {
					isSession = true
					env = &dvevaluation.DvObject{Properties: ctx.Session.Values()}
				}
			} else if ctx.LocalContextEnvironment != nil && level != "request" {
				if level != "" && level[0] >= '1' && level[0] <= '9' {
					levelVal, err := strconv.Atoi(level)
					if err == nil {
						if path == "" {
							ctx.LocalContextEnvironment.DeleteAtParent(varName, levelVal)
						} else {
							res, ok = ctx.LocalContextEnvironment.ReadAtParent(varName, levelVal)
							if ok {
								res, _, err = dvjson.RemovePathOfAny(res, path, false, ctx.LocalContextEnvironment)
								if err != nil {
									log.Printf("Cannot remove %s : %v", path, err)
								} else {
									ctx.LocalContextEnvironment.SetAtParent(varName, res, levelVal)
								}
							}
						}
					} else {
						dvlog.PrintfError("Unknown level %s", level)
					}
				} else {
					if path == "" {
						ctx.LocalContextEnvironment.Delete(varName)
					} else {
						env = ctx.LocalContextEnvironment
					}
				}
			} else {
				if path == "" {
					ctx.PrimaryContextEnvironment.Delete(varName)
				} else {
					env = ctx.PrimaryContextEnvironment
				}
			}
		} else {
			if isMap {
				mapName := level[4:]
				dvsession.GlobalMapDelete(mapName, varName)
			} else {
				dvparser.RemoveGlobalPropertiesValue(varName)
			}
		}
		if path != "" && env != nil {
			res, ok = env.Get(varName)
			var err error
			if ok {
				res, _, err = dvjson.RemovePathOfAny(res, path, false, env)
				if err != nil {
					log.Printf("Cannot remove %s : %v", path, err)
				} else {
					if isSession {
						ctx.Session.SetItem(varName, res)
					} else {
						env.Set(varName, res)
					}
				}
			}
		}
	}
}

func GetLevelMainPath(result string) (string, string, string) {
	p := strings.Index(result, ":")
	level := ""
	path := ""
	pathIsProcessed := true
	if p >= 0 {
		level = strings.ToLower(result[:p])
		result = result[p+1:]
		if len(level) > 0 && level[0] == '~' {
			level = level[1:]
			pathIsProcessed = false
		}
	}
	p = strings.Index(result, ".")
	if pathIsProcessed && p >= 0 {
		path = result[p+1:]
		result = result[:p]
	}
	return level, result, path
}

func ReadActionResult(result string, ctx *dvcontext.RequestContext) (res interface{}, ok bool) {
	env := GetEnvironment(ctx)
	if result != "" {
		if result[0] == '\'' && result[len(result)-1] == '\'' && len(result) >= 2 {
			return result[1 : len(result)-1], true
		}
		level, varName, path := GetLevelMainPath(result)
		isMap := strings.HasPrefix(level, "map_")
		if ctx != nil && level != "global" && !isMap {
			switch level {
			case "":
				res, ok = env.Get(varName)
			case "_":
				j, err := dvjson.ParseAny(result[2:])
				if err != nil {
					dvlog.PrintfError("Failed to convert to json %s: %v", result, err)
					j = nil
					ok = false
				}
				res = j
			case "$":
				rs, err := env.EvaluateAnyTypeExpression(result[2:])
				if err != nil {
					dvlog.PrintfError("Failed to evaluate %s: %v", result, err)
					rs = nil
					ok = false
				}
				return rs, true
			case "session", "session?":
				isErrorFatal := level == "session"
				if ctx.Session == nil {
					if isErrorFatal {
						dvlog.PrintlnError("No session available")
						ctx.HandleInternalServerError()
					}
					return
				}
				res = ctx.Session.GetItem(varName)
				if path != "" {
					env = &dvevaluation.DvObject{Properties: ctx.Session.Values()}
				}
			default:
				if ctx.LocalContextEnvironment != nil && level != "request" {
					if level[0] >= '1' && level[0] <= '9' {
						levelVal, err := strconv.Atoi(level)
						if err == nil {
							res, ok = ctx.LocalContextEnvironment.ReadAtParent(varName, levelVal)
						} else {
							dvlog.PrintfError("Unknown level %s", level)
						}
					} else {
						res, ok = ctx.LocalContextEnvironment.Get(varName)
					}
				} else {
					res, ok = ctx.PrimaryContextEnvironment.Get(varName)
				}
			}
		} else {
			if isMap {
				mapName := level[4:]
				res, ok = dvsession.GlobalMapRead(mapName, varName)
			} else {
				res, ok = dvparser.GlobalProperties[varName]
			}
		}
		if ok && path != "" {
			var err error
			res, _, err = dvjson.ReadPathOfAny(res, path, false, env)
			if err != nil {
				res = nil
				ok = false
			}
		}
	}
	return
}

func IsLikeJson(s string) bool {
	t := strings.TrimSpace(s)
	n := len(t)
	return n >= 2 && (t[0] == '{' && t[n-1] == '}' || t[0] == '[' && t[n-1] == ']')
}

func DefaultInitWithObject(command string, result interface{}, env *dvevaluation.DvObject) bool {
	cmd := strings.TrimSpace(command[strings.Index(command, ":")+1:])
	if cmd == "" {
		log.Printf("Empty parameters in %s", command)
		return false
	}
	cmdDat := []byte(cmd)
	if cmd[0] != '{' || cmd[len(cmd)-1] != '}' {
		if cmd[0] == '@' && len(cmd) > 1 {
			if cmd[1] == '@' && len(cmd) > 2 {
				dat, ok := env.Get(cmd[2:])
				if !ok {
					log.Printf("Empty parameter %s", cmd)
					return false
				}
				s := strings.TrimSpace(dvevaluation.AnyToString(dat))
				if strings.Contains(s, "{{") {
					ss, err := env.CalculateString(s)
					if err != nil {
						log.Printf("Wrong expression in %s: %v", s, err)
						return false
					}
					s = ss
				}
				if len(s) == 0 || s[0] != '{' || s[len(s)-1] != '}' {
					log.Printf("Wrong object described by %s", cmd)
					return false
				}
				cmdDat = []byte(s)
			} else {
				dat, err := dvparser.SmartReadLikeJsonTemplate(cmd[1:], 3, env)
				if err != nil {
					log.Printf("Bad file in %s %v", command, err)
					return false
				}
				dat = bytes.TrimSpace(dat)
				if len(dat) < 2 || dat[0] != '{' || dat[len(dat)-1] != '}' {
					log.Printf("Empty file in %s", command)
					return false
				}
				cmdDat = dat
			}
		} else {
			log.Printf("Empty parameters in %s", command)
			return false
		}
	}
	err := json.Unmarshal(cmdDat, result)
	if err != nil {
		log.Printf("Error converting parameters: %v in %s", err, command)
		return false
	}
	return true
}

func ActionFinalException(status int, body []byte, ctx *dvcontext.RequestContext) {
	ctx.StatusCode = status
	ctx.Output = body
	ctx.PrimaryContextEnvironment.Set(ExSeqLevel, -2)
}

func ActionInternalException(status int, shortMessage string, longMessage string, ctx *dvcontext.RequestContext) {
	policy := ctx.GetCurrentErrorPolicy()
	res := policy.Format
	t := time.Now()
	stamp := t.Format(time.RFC850)
	res = strings.Replace(res, "$$$TIMESTAMP", stamp, -1)
	res = strings.Replace(res, "$$$STATUS", strconv.Itoa(status), -1)
	res = strings.Replace(res, "$$$PATH", ctx.Url, -1)
	res = strings.Replace(res, "$$$MESSAGE", shortMessage, 1)
	res = strings.Replace(res, "$$$ERROR", longMessage, 1)
	ActionFinalException(status, []byte(res), ctx)
}

func ActionExternalException(status int, body []byte, ctx *dvcontext.RequestContext) {
	policy := ctx.GetCurrentErrorPolicy()
	if policy.FormatForced {
		shortMessage := strings.Replace(ctx.PrimaryContextEnvironment.GetString("ERROR_POLICY_STANDARD_ERROR"), "$$$STATUS", strconv.Itoa(status), -1)
		if shortMessage == "" {
			shortMessage = string(body)
		}
		ActionInternalException(status, shortMessage, string(body), ctx)
		return
	}
	ActionFinalException(status, body, ctx)
}

func ActionExceptionByError(comment string, err error, ctx *dvcontext.RequestContext) {
	longMessage := err.Error()
	if comment == "" {
		comment = "Internal System Error"
	}
	ActionInternalException(500, comment, longMessage, ctx)
}
