/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

const (
	UrlPathPrefix = "URL_PATH_"
	ErrorPolicy   = "ERROR_POLICY"
)

func (ctx *RequestContext) SetHttpErrorCode(errorCode int, message string) {
	if ctx == nil {
		log.Printf("Error %d: %s", errorCode, message)
	} else {
		ctx.StatusCode = errorCode
		if message != "" {
			ctx.Error = errors.New(message)
		}
	}
}

func (ctx *RequestContext) SetErrorMessage(message string) {
	ctx.SetHttpErrorCode(500, message)
}

func (ctx *RequestContext) SetError(err error) {
	ctx.SetHttpErrorCode(500, err.Error())
}

func (ctx *RequestContext) SetHeaderUnique(key string, val string) {
	if ctx.Headers == nil {
		ctx.Headers = make(map[string][]string)
	}
	ctx.Headers[key] = []string{val}
}

func (ctx *RequestContext) SetHeaderArray(key string, val []string) {
	if ctx.Headers == nil {
		ctx.Headers = make(map[string][]string)
	}
	ctx.Headers[key] = val
}

func (ctx *RequestContext) SetUrlInlineParameters(params map[string]string) {
	ctx.UrlInlineParams = params
	ctx.PrimaryContextEnvironment.SetPropertiesWithPrefixFromString(UrlPathPrefix, params, dvevaluation.TransformUpperCase)
}

func (ctx *RequestContext) GetCurrentErrorPolicy() *RequestErrorPolicy {
	if ctx != nil && ctx.Server != nil && ctx.Server.ErrorPolicies != nil &&
		ctx.PrimaryContextEnvironment != nil {
		policy := ctx.Server.ErrorPolicies[ctx.PrimaryContextEnvironment.GetString(ErrorPolicy)]
		if policy != nil {
			if policy.Format == "" {
				policy.Format = DefaultRequestErrorPolicy.Format
				policy.ContentType = DefaultRequestErrorPolicy.ContentType
			} else if policy.ContentType == "" {
				policy.ContentType = DefaultRequestErrorPolicy.ContentType
			}
			return policy
		}
	}
	return DefaultRequestErrorPolicy
}

func (ctx *RequestContext) GetEnvironment() *dvevaluation.DvObject {
	if ctx == nil || ctx.PrimaryContextEnvironment == nil {
		return dvparser.GetGlobalPropertiesAsDvObject()
	}
	if ctx.LocalContextEnvironment != nil {
		return ctx.LocalContextEnvironment
	}
	return ctx.PrimaryContextEnvironment
}

func (ctx *RequestContext) StoreStoringSession() {
	if ctx.Session != nil && ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Set(ServerSessionStoringKey, ctx.Session)
	}
}

func (ctx *RequestContext) GetUrlParameter(param string) string {
	if ctx.UrlInlineParams == nil {
		return ""
	}
	return ctx.UrlInlineParams[param]
}
