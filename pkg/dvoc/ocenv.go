/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"strings"
	"time"
)

func OpenShiftConvertListToMap(list string) (res map[string][]string, ok bool) {
	items := strings.Split(list, ";;")
	res = make(map[string][]string)
	n := len(items)
	ok = true
	for i := 0; i < n; i++ {
		s := strings.TrimSpace(items[i])
		if s == "" {
			continue
		}
		pos := strings.Index(s, ":")
		if pos <= 0 {
			dvlog.PrintlnError("Format for environment is as follows:")
			dvlog.PrintlnError("microservice:envName=envName;;microservice:envName=envValue")
			dvlog.PrintfError("You missed : at %d in %s ", i, list)
			return nil, false
		}
		name := strings.TrimSpace(s[:pos])
		rest := strings.TrimSpace(s[pos+1:])
		if strings.Index(name, " ") >= 0 || len(name) == 0 {
			dvlog.PrintlnError("Format for environment is as follows:")
			dvlog.PrintlnError("microservice:envName=envName;;microservice:envName=envValue")
			dvlog.PrintfError("You have extra space in microservice name at %d in %s ", i, list)
			return nil, false
		}
		pos = strings.Index(rest, "=")
		if pos <= 0 {
			dvlog.PrintlnError("Format for environment is as follows:")
			dvlog.PrintlnError("microservice:envName=envName;;microservice:envName=envValue")
			dvlog.PrintfError("You did not specify = at %d in %s ", i, list)
			return nil, false
		}
		k := strings.TrimSpace(rest[:pos])
		v := strings.TrimSpace(rest[pos+1:])
		if k == "" || strings.Index(k, " ") >= 0 || v == "" || strings.Index(v, " ") >= 0 {
			dvlog.PrintlnError("Format for environment is as follows:")
			dvlog.PrintlnError("microservice:envName=envName;;microservice:envName=envValue")
			if k == "" {
				dvlog.PrintfError("You did not specify envName at %d in %s ", i, list)
			}
			if strings.Index(k, " ") >= 0 {
				dvlog.PrintfError("You specified envName with has extra space at %d in %s ", i, list)
			}
			if v == "" {
				dvlog.PrintfError("You did not specify envValue at %d in %s ", i, list)
			}
			if strings.Index(v, " ") >= 0 {
				dvlog.PrintfError("You specified envValue with has extra space at %d in %s ", i, list)
			}
			return nil, false
		}
		k = k + "=" + v
		if res[name] == nil {
			res[name] = make([]string, 1, n)
			res[name][0] = k
		} else {
			res[name] = append(res[name], k)
		}
	}
	return
}

/***
For environment variable changes:
When an old value is set, openshift automatically does nothing
When a new value is set, openshift makes a change to the deployment config
and an old pod is destroyed and a new pod is created and started
*/
func OpenShiftSetEnv(list string) bool {
	envs, ok := OpenShiftConvertListToMap(list)
	if !ok {
		return false
	}
	res := true
	for k, v := range envs {
		if !OpenShiftSetEnvironment(k, v) {
			res = false
		}
	}
	return res
}

func OpenShiftScale(list string, replicas int) bool {
	microServices := dvtextutils.ConvertToNonEmptyList(list)
	n := len(microServices)
	res := true
	for i := 0; i < n; i++ {
		if !OpenShiftScaleToReplicas(microServices[i], replicas) {
			res = false
		}
	}
	return res
}

func OpenShiftResetUp(service string) {
	time.Sleep(60 * time.Second)
	OpenShiftScaleToReplicas(service, 1)
}

func OpenShiftReset(list string) bool {
	microServices := dvtextutils.ConvertToNonEmptyList(list)
	n := len(microServices)
	res := true
	for i := 0; i < n; i++ {
		name := microServices[i]
		if !OpenShiftScaleToReplicas(name, 0) {
			res = false
		}
		OpenShiftResetUp(name)
	}
	return res
}

func processEnvSettingInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	if !strings.HasPrefix(command, "env:") {
		dvlog.PrintfError("command %s must start with env:", command)
		return nil, false
	}
	command = strings.TrimSpace(command[4:])
	n := len(command) - 1
	if n < 0 || command[0] != '{' || command[n] != '}' {
		dvlog.PrintfError("env format is env:{microservice:key=value;;microservice:key=value} but it is %s", command)
		return nil, false
	}
	command = command[1:n]
	return []interface{}{command}, true
}

func processEnvSettingsRun(data []interface{}) bool {
	command := data[0].(string)
	return OpenShiftSetEnv(command)
}
