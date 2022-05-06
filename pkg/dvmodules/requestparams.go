/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"strings"
)

func collectRequestParameters(request *dvcontext.RequestContext) error {
	action := request.Action
	bodyParams := action.Body
	if action.Name != "" && request.PrimaryContextEnvironment != nil {
		logLevel := request.PrimaryContextEnvironment.GetString(action.Name + "_LOG")
		if logLevel != "" {
			request.LogLevel = dvlog.GetLogLevel(logLevel)
		}
	}
	method := request.Reader.Method
	if method != "GET" {
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
				if err == nil {
					request.PrimaryContextEnvironment.Set(dvcontext.BODY_JSON, request.InputJson)
					if bodyParams != nil {
						dvevaluation.CollectJsonVariables(request.InputJson, bodyParams, request.PrimaryContextEnvironment, true)
					}
				}
			} else {
				mp, err := dvtextutils.DecomposeUrlEncodedString(request.InputStr)
				if err != nil && request.LogLevel >= dvlog.LogInfo {
					dvlog.PrintfError("Body parsing error: %v %s", err, request.InputStr)
				}
				request.PrimaryContextEnvironment.Set(dvcontext.BODY_STRING_ARRAY_MAP, mp)
				m := dvtextutils.ConvertArrayStringMapIntoSimpleStringMap(mp)
				request.PrimaryContextEnvironment.Set(dvcontext.BODY_STRING_MAP, m)
				if len(m) > 0 {
					for k, _ := range bodyParams {
						v := m[k]
						key := dvcontext.BODY_PARAM_PREFIX + dvtextutils.ConvertToUpperAlphaDigital([]byte(k))
						request.PrimaryContextEnvironment.Set(key, v)
					}
				}
			}
		}
	}
	dvevaluation.CollectVariablesByStringMap(request.Queries, action.QueryParams, request.PrimaryContextEnvironment, true)
	return nil
}
