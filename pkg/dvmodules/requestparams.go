/***********************************************************************
MicroCore
Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

const (
	multiPartFormData = "multipart/form-data;"
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
	var err error = nil
	if method != "GET" {
		if strings.HasPrefix(request.GetContentType(), multiPartFormData) {
			err = processMultipartFormDataBody(request, bodyParams)
		} else {
			err = processSimpleBody(request, bodyParams)
		}
		if err != nil {
			return err
		}
	}
	dvevaluation.CollectVariablesByStringMap(request.Queries, action.QueryParams, request.PrimaryContextEnvironment, true, "", true)
	return nil
}

func processMultipartFormDataBody(request *dvcontext.RequestContext, bodyParams map[string]string) error {
	err := request.Reader.ParseMultipartForm(200 << 20)
	if err != nil {
		dvlog.PrintlnError(err.Error())
	}
	body := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, 16)}
	for k, v := range request.Reader.Form {
		val := dvtextutils.FindNonEmptyLastString(v)
		if bodyParams != nil {
			tp := bodyParams[k]
			if len(tp) > 0 {
				key := dvcontext.BODY_PARAM_PREFIX + dvtextutils.ConvertToUpperAlphaDigital([]byte(k))
				request.PrimaryContextEnvironment.Set(key, val)
			}
		}
		body.Fields = append(body.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Name: []byte(k), Value: []byte(val)})
	}
	for k, v := range request.Reader.MultipartForm.File {
		f := v[0]
		origName := f.Filename
		tmpName := generateTmpFileName(strconv.FormatInt(request.Id, 10), origName)
		tmpName, err = saveFile(f, tmpName)
		if err != nil {
			return err
		}
		r := "file:" + pureName(origName) + ":" + tmpName
		request.AddTempFile(tmpName)
		body.Fields = append(body.Fields, &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Name: []byte(k), Value: []byte(r)})
	}
	request.PrimaryContextEnvironment.Set(dvcontext.BODY_STRING, body)
	return nil
}

func saveFile(f *multipart.FileHeader, fileName string) (string, error) {
	fil, err := f.Open()
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(fil)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(fileName, data, 0644)
	return fileName, err
}

func generateTmpFileName(name string, ext string) string {
	name = name + dvtextutils.GetLowCaseExtension(ext)
	fullName := dvdir.GetTempPathSlashed() + name
	return fullName
}

func pureName(name string) string {
	p := strings.LastIndex(name, ":")
	if p >= 0 {
		name = name[p+1:]
	}
	p = strings.LastIndex(name, "/")
	if p >= 0 {
		name = name[p+1:]
	}
	p = strings.LastIndex(name, "\\")
	if p >= 0 {
		name = name[p+1:]
	}
	return name
}

func processSimpleBody(request *dvcontext.RequestContext, bodyParams map[string]string) error {
	body, err := io.ReadAll(request.Reader.Body)
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
					dvevaluation.CollectJsonVariables(request.InputJson, bodyParams, request.PrimaryContextEnvironment, true, dvcontext.BODY_PARAM_PREFIX, true)
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
	return nil
}
