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
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strconv"
	"strings"
)

type StorageActionProvider interface {
	PathSupported() bool
	Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool)
	Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error
	Delete(ctx *dvcontext.RequestContext, prefix string, key string) error
}

var StorageProviders map[string]StorageActionProvider = make(map[string]StorageActionProvider, 32)

func RegisterStorageActionProvider(name string,provider StorageActionProvider) {
	StorageProviders[name] = provider
}

func GetLevelMainPath(data string) (level string, subLevel string, result string, path string, prePath string, isFatal bool, provider StorageActionProvider) {
	pos := strings.Index(data, ":")
	result = data
	pathIsProcessed := true
	isFatal = true
	if pos >= 0 {
		level = strings.TrimSpace(result[:pos])
		result = result[pos+1:]
		if len(level) > 0 && level[0] == '~' {
			level = level[1:]
			pathIsProcessed = false
		}
		p := strings.Index(level, "_")
		if p > 0 {
			subLevel = strings.TrimSpace(level[p+1:])
			level = strings.TrimSpace(level[:p])
		}
		p = len(level)
		if p > 0 && level[p-1] == '?' {
			isFatal = false
			level = level[:p]
		}
		level = strings.ToLower(level)
		provider = StorageProviders[level]
		if provider != nil && !provider.PathSupported() {
			pathIsProcessed = false
		}
	}
	dot := strings.Index(result, ".")
	if pathIsProcessed && dot >= 0 {
		path = result[dot+1:]
		result = result[:dot]
		prePath = data[:len(data)-len(path)-1]
	}
	return
}

func SaveActionResult(result string, data interface{}, ctx *dvcontext.RequestContext) {
	if result != "" {
		env := GetEnvironment(ctx)
		level, subLevel, varName, path, prePath, isFatal, provider := GetLevelMainPath(result)
		if path != "" {
			dat, ok := ReadActionResult(prePath, ctx)
			if !ok {
				data = dvevaluation.CreateDvVariableByPathAndData(path, data, nil)
			} else {
				dvevaluation.UpdateAnyVariables(dat, data, path, dvevaluation.UPDATE_MODE_REPLACE, nil, env)
				data = dat
			}
		}
		if provider != nil {
			err := provider.Save(ctx, subLevel, varName, data)
			if err != nil {
				dvlog.PrintfError("%s", err.Error())
				if isFatal {
					ctx.HandleInternalServerError()
				}
			}
			return
		}
		if ctx != nil && ctx.PrimaryContextEnvironment != nil {
			if level != "" && level[0] >= '1' && level[0] <= '9' {
				levelVal, err := strconv.Atoi(level)
				if err == nil {
					ctx.LocalContextEnvironment.SetAtParent(varName, data, levelVal)
				} else {
					dvlog.PrintfError("Unknown level %s", level)
				}
			} else {
				env.Set(varName, data)
			}
		} else {
			dvparser.SetGlobalPropertiesAnyValue(varName, data)
		}
	}
}

func DeleteActionResult(result string, ctx *dvcontext.RequestContext) {
	if result != "" {
		env := GetEnvironment(ctx)
		level, subLevel, varName, path, prePath, isFatal, provider := GetLevelMainPath(result)
		if path != "" {
			res, ok := ReadActionResult(prePath, ctx)
			var err error
			if ok {
				res, _, err = dvjson.RemovePathOfAny(res, path, false, env)
				if err != nil {
					log.Printf("Cannot remove %s : %v", path, err)
				} else {
					SaveActionResult(prePath, res, ctx)
				}
			}
			return
		}
		if provider != nil {
			err := provider.Delete(ctx, subLevel, varName)
			if err != nil {
				dvlog.PrintfError("%s", err.Error())
				if isFatal {
					ctx.HandleInternalServerError()
				}
			}
			return
		}
		if ctx != nil && ctx.PrimaryContextEnvironment != nil {
			if ctx.LocalContextEnvironment != nil && level != "" && level[0] >= '1' && level[0] <= '9' {
				levelVal, err := strconv.Atoi(level)
				if err == nil {
					ctx.LocalContextEnvironment.DeleteAtParent(varName, levelVal)
				} else {
					dvlog.PrintfError("Unknown level %s", level)
				}
			} else {
				env.Delete(varName)
			}
		}
	}
}

func ReadActionResult(result string, ctx *dvcontext.RequestContext) (res interface{}, ok bool) {
	env := GetEnvironment(ctx)
	if result != "" {
		if result[0] == '\'' && result[len(result)-1] == '\'' && len(result) >= 2 {
			return result[1 : len(result)-1], true
		}
		level, subLevel, varName, path, _, _, provider := GetLevelMainPath(result)
		if provider != nil {
			res, ok = provider.Read(ctx, subLevel, varName)
		} else {
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
				ok = true
			case "$":
				rs, err := env.EvaluateAnyTypeExpression(result[2:])
				if err != nil {
					dvlog.PrintfError("Failed to evaluate %s: %v", result, err)
					rs = nil
					ok = false
				}
				return rs, true
			default:
				if ctx.LocalContextEnvironment != nil {
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
					res, ok = env.Get(varName)
				}
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
