/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import (
	"strconv"
	"sync"
)

type logInfoCounters struct {
	firstCounter int
	nextCounter  int
}

var logInfoMap = make(map[string]*logInfoCounters, 32)
var logInfoSync sync.Mutex

func getRequestLogInfo(message string, policyType *ActionPolicy, method string, url string) string {
	key := method + "$" + url
	omit := 0
	logInfoSync.Lock()
	item := logInfoMap[key]
	if item == nil {
		item = &logInfoCounters{
			firstCounter: policyType.LogFirstTimes,
			nextCounter:  policyType.LogNextTime,
		}
		logInfoMap[key] = item
	}
	if item.firstCounter > 0 {
		omit = 1
		item.firstCounter--
	} else if item.nextCounter > 0 {
		item.nextCounter--
	} else if policyType.LogNextTime > 0 {
		omit = policyType.LogNextTime
		item.nextCounter = policyType.LogNextTime
	}
	logInfoSync.Unlock()
	if omit == 0 {
		return DoNotShowPlaceInfo
	}
	if omit == 1 {
		return message
	}
	if message == "0" {
		message = ""
	}
	return "(" + strconv.Itoa(omit) + ") " + message
}
