/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"log"
	"strconv"
	"strings"
)

type ValidationItem struct {
	Data  string `json:"data"`
	Part  string `json:"part"`
	Kind  string `json:"kind"`
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type ValidationConfig struct {
	Check  []*ValidationItem `json:"check"`
	Errors []string          `json:"errors"`
	Code   int               `json:"code"`
}

func validationInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ValidationConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if len(config.Check) == 0 {
		log.Printf("check must be specified in %s", command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func validationRun(data []interface{}) bool {
	config := data[0].(*ValidationConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ValidationRunByConfig(config, ctx)
}

func ValidationRunByConfig(config *ValidationConfig, ctx *dvcontext.RequestContext) bool {
	n := len(config.Check)
	env := GetEnvironment(ctx)
	for i := 0; i < n; i++ {
		code, message, ok := ValidateItem(config.Check[i], env)
		if !ok {
			if code < 400 || code >= 600 {
				code = config.Code
				if code < 400 || code >= 600 {
					code = 500
				}
			}
			part:=config.Check[i].Part
			if message != "" && message[0] >= '0' && message[0] <= '9' {
				p, err := strconv.Atoi(message)
				if err == nil && p >= 0 && p < len(config.Errors) {
					message = config.Errors[p]
				}
			}
			if message == "" {
				if len(config.Errors) > 0 && config.Errors[0] != "" {
					message = config.Errors[0]
				} else {
					message = "Internal Server Error"
				}
			}
			if strings.Contains(message,"$part") {
				message = strings.Replace(message,"$part", part, -1)
			}
			ActionExternalException(code, []byte(message), ctx)
			return true
		}
	}
	return true
}

func ValidateExists(data string, env *dvevaluation.DvObject) bool {
	if strings.Contains(data, ".") {
		r, err := env.EvaluateAnyTypeExpression(data)
		return err != nil && r != nil
	}
	_, ok := env.Get(data)
	return ok
}

func ValidateNonEmpty(data string, env *dvevaluation.DvObject) bool {
	if strings.Contains(data, ".") {
		r, err := env.EvaluateAnyTypeExpression(data)
		return err != nil && dvevaluation.AnyToBoolean(r)
	}
	v, ok := env.Get(data)
	return ok && dvevaluation.AnyToBoolean(v)
}

func ValidateCondition(data string, env *dvevaluation.DvObject) (bool, bool) {
	r, err := env.EvaluateAnyTypeExpression(data)
	if err != nil {
		return false, true
	}
	return dvevaluation.AnyToBoolean(r), false
}

func ValidateItem(item *ValidationItem, env *dvevaluation.DvObject) (int, string, bool) {
	message := item.Error
	isError := false
	data := item.Data
	switch item.Kind {
	case "exists":
		isError = !ValidateExists(data, env)
	case "non-exists":
		isError = ValidateExists(data, env)
	case "empty":
		isError = ValidateNonEmpty(data, env)
	case "non-empty":
		isError = !ValidateNonEmpty(data, env)
	case "condition", "":
		isCond, isFatal := ValidateCondition(data, env)
		isError = isFatal || !isCond
	case "not":
		isCond, isFatal := ValidateCondition(data, env)
		isError = isFatal || isCond
	default:
		isError = true
		message = "Unknown kind " + item.Kind
	}
	if isError {
		if strings.Contains(message, "%s") {
			r, err := env.EvaluateAnyTypeExpression(data)
			if err != nil {
				data = err.Error()
			} else {
				data = dvevaluation.AnyToString(r)
			}
			message = fmt.Sprintf(message, data)
		}
		return item.Code, message, false
	}
	return 0, "", true
}
