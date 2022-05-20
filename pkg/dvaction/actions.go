// package dvoc orchestrates actions, executions
// MicroCore Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvsecurity"
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
