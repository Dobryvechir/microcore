/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmeta

import (
	"errors"
	"log"
	"strconv"
)

func (ctx *RequestContext) SetHttpErrorCode(errorCode int, message string) {
	if ctx == nil {
		log.Printf("Error %d: %s", errorCode, message)
	} else {
		if message != "" {
			ctx.Output = []byte(message)
		}
		ctx.Error = errors.New(strconv.Itoa(errorCode) + " " + message)
	}
}

func (ctx *RequestContext) SetErrorMessage(message string) {
	ctx.SetHttpErrorCode(500, message)
}

func (ctx *RequestContext) SetError(err error) {
	ctx.SetHttpErrorCode(500, err.Error())
}
