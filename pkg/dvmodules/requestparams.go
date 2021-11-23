/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"io/ioutil"
	"strings"
)

func collectRequestParameters(request *dvcontext.RequestContext) error {
	action := request.Action
	bodyParams := action.Body
	method := request.Reader.Method
	if bodyParams != nil && method != "GET" {
		body, err := ioutil.ReadAll(request.Reader.Body)
		if err != nil {
			return err
		}
		if len(body) > 0 {
			request.Input = body
			request.InputStr = strings.TrimSpace(string(body))
			request.PrimaryContextEnvironment.Set(dvcontext.BODY_STRING, request.InputStr)
			if len(request.InputStr) > 0 && (request.InputStr[0] == '[' || request.InputStr[0] == '{') {
				request.InputJson, err = dvjson.JsonFullParser(request.Input)
				if err != nil {
					request.PrimaryContextEnvironment.Set(dvcontext.BODY_JSON, request.InputJson)
					dvevaluation.CollectJsonVariables(request.InputJson, bodyParams, request.PrimaryContextEnvironment, true)
				}
			}
		}
	}
	dvevaluation.CollectVariablesByStringMap(request.Queries, action.QueryParams, request.PrimaryContextEnvironment, true)
	return nil
}
