/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strings"
)

type ContainInProcessor func(interface{}, interface{}) (bool, bool, error)

var listContainInProcessors = make([]ContainInProcessor, 0, 7)

func ContainInProcess(contained interface{}, containing interface{}) bool {
	switch containing.(type) {
	case []string:
		v := AnyToString(contained)
		return dvtextutils.IsStringContainedInArray(v, containing.([]string))
	case map[string]string:
		v := AnyToString(contained)
		return dvtextutils.IsStringContainedInStringToStringMap(v, containing.(map[string]string))
	case map[string]interface{}:
		v := AnyToString(contained)
		return dvtextutils.IsStringContainedInStringToAnyMap(v, containing.(map[string]interface{}))
	case string:
		v := AnyToString(contained)
		return strings.Contains(containing.(string), v)
	case *DvVariable:
		res := containing.(*DvVariable).ContainsItemIn(contained)
		return res
	default:
		n := len(listContainInProcessors)
		for i := 0; i < n; i++ {
			res, ok, err := listContainInProcessors[i](contained, containing)
			if !ok {
				continue
			}
			if err != nil {
				dvlog.PrintlnError("Error contain processing:" + err.Error())
				break
			}
			return res
		}
	}
	return false
}

func RegisterContainInProcessor(fn ContainInProcessor) {
	listContainInProcessors = append(listContainInProcessors, fn)
}
