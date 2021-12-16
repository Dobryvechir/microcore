/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"strconv"
	"strings"
)

const (
	dcEnvironmentVariablesPath = "spec.template.spec.containers[0].env"
)

func ReadPodEnvironmentAndDC(microServiceName string) (res map[string]string, dc *dvevaluation.DvVariable, err error) {
	res = make(map[string]string)
	dc, err = GetLiveDeploymentConfiguration(microServiceName)
	if err != nil {
		return
	}
	data, _, err := dc.ReadChild(dcEnvironmentVariablesPath, nil)
	if err != nil || data == nil {
		return
	}
	res, err = dvjson.ConvertFieldItemArrayIntoMap(data, "name", "value", "")
	return
}

func GetSpecificVariableAtServer(microServiceName string, key string) (data string, err error) {
	info, err := ExecuteCommandOnPod(microServiceName, "echo $"+key)
	if err != nil {
		return "", err
	}
	return info, nil
}

func ReadPodReadyEnvironment(microServiceName string) (res map[string]string, dc *dvevaluation.DvVariable, err error) {
	res, dc, err = ReadPodEnvironmentAndDC(microServiceName)
	if err != nil {
		return
	}
	for k, v := range res {
		if v != "" && v[0] == '{' && strings.Index(v, "from") >= 0 {
			data, err1 := GetSpecificVariableAtServer(microServiceName, k)
			if err1 != nil {
				err = err1
				return
			}
			res[k] = data
		}
	}
	return
}

func ReadTemplateParameters(templateData []byte) (res map[string]string, obj *dvevaluation.DvVariable, err error) {
	obj, err = dvjson.JsonFullParser(templateData)
	if err != nil {
		return
	}
	if obj == nil {
		err = errors.New("Empty openshift template")
		return
	}
	res = make(map[string]string)
	params := obj.ReadSimpleChild("parameters")
	if params == nil {
		return
	}
	if params.Kind != dvevaluation.FIELD_ARRAY {
		err = errors.New("parameters at first level must be an array")
		return
	}
	f := params.Fields
	n := len(f)
	for i := 0; i < n; i++ {
		if f[i] != nil && f[i].Kind == dvevaluation.FIELD_OBJECT && f[i].Fields != nil {
			fld := f[i]
			fldName := fld.ReadSimpleChild("name")
			fldValue := fld.ReadSimpleChild("value")
			if fldName == nil {
				err = errors.New("Corrupted entry's name in parameters at " + strconv.Itoa(i))
			} else {
				k := fldName.GetStringValue()
				v := ""
				if fldValue != nil {
					v = fldValue.GetStringValue()
				}
				res[k] = v
			}
		} else {
			err = errors.New("Incorrect entry in parameters at " + strconv.Itoa(i))
		}
	}
	return
}
