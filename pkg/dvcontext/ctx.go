/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"log"
)

const (
	UrlPrefix = "URL_"
)

func (ctx *RequestContext) SetHttpErrorCode(errorCode int, message string) {
	if ctx == nil {
		log.Printf("Error %d: %s", errorCode, message)
	} else {
		ctx.StatusCode = errorCode
		if message!="" {
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

func (ctx *RequestContext) SetUrlInlineParameters(params map[string]string) {
	ctx.UrlInlineParams = params
	ctx.PrimaryContextEnvironment.SetPropertiesWithPrefixFromString(UrlPrefix, params, dvevaluation.TransformUpperCase)
}
