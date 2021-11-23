/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"sort"
	"strings"
)

const (
	kubernetesConfiguration             = "kubectl.kubernetes.io/last-applied-configuration"
	optionSilentForNotFoundMicroService = 1
)

var TemplateObjectPriority = map[string]int{
	"deploymentconfig": 1,
	"imagestream":      2,
	"service":          3,
	"configmap":        4,
	"route":            5,
	"serviceaccount":   6,
	"secret":           7,
}

func GetKubernetesConfigurationPart(cmdLine string, kind string, mode int, fn ConfigTransformer, notCritical bool) (string, error) {
	var info string
	var ok bool
	silentForNotFound := (mode & optionSilentForNotFoundMicroService) != 0
	isWithEditor := strings.HasPrefix(cmdLine, "edit")
	if isWithEditor {
		info, ok = RunOCCommandWithEditor(cmdLine)
	} else {
		if silentForNotFound {
			var status int
			info, status = RunOCCommandFailureAllowed(cmdLine, []string{"NotFound"})
			ok = status >= 0
		} else {
			info, ok = RunOCCommand(cmdLine)
		}
	}
	if !ok {
		if !notCritical {
			dvlog.PrintfError("Failed to execute %s", cmdLine)
		}
		return "", errors.New("Failed to execute " + cmdLine)
	}
	original := info
	pos := strings.Index(info, kubernetesConfiguration)
	if pos < 0 {
		if fn == nil {
			return "", errors.New(kind + " does not contain " + kubernetesConfiguration + " (see " + dvdir.SaveToUniqueFile(original) + ")")
		}
		data, err := dvjson.ReadYamlAsDvFieldInfo([]byte(info))
		if err != nil {
			return "", err
		}
		data, err = fn(data)
		if err != nil {
			return "", err
		}
		info = string(dvjson.PrintToJson(data,4))
		return info, nil
	}
	info = strings.TrimSpace(info[pos+len(kubernetesConfiguration):])
	if info == "" || info[0] != ':' && info[0] != '=' {
		return "", errors.New(kind + " expected to contain : or = after " + kubernetesConfiguration + " (see " + dvdir.SaveToUniqueFile(original) + ")")
	}
	info = dvtextutils.GetNextNonEmptyPartInYaml(info[1:])
	pos = strings.Index(info, "\n")
	if info == "" || info[0] != '{' || pos < 0 {
		return "", errors.New("corrupted " + kind + " (see " + dvdir.SaveToUniqueFile(original) + ")")
	}
	info = strings.TrimSpace(info[:pos])
	if info[len(info)-1] != '}' {
		return "", errors.New("Corrupted " + kind + " (see " + dvdir.SaveToUniqueFile(original) + ")")
	}
	return info, nil
}

func GetLiveConfiguration(cmdLine string) (*dvevaluation.DvVariable, error) {
	info, ok := RunOCCommandWithEditor(cmdLine)
	if !ok {
		return nil, errors.New("Failed to execute " + cmdLine)
	}
	return dvjson.ReadYamlAsDvFieldInfo([]byte(info))
}

func GetLiveDeploymentConfiguration(microServiceName string) (*dvevaluation.DvVariable, error) {
	return GetLiveConfiguration("edit dc " + microServiceName)
}

func GetShortOpenShiftNameForObjectType(openShiftObjectType string) (string, error) {
	switch strings.ToLower(openShiftObjectType) {
	case "deploymentconfig":
		return "dc", nil
	case "secret", "secrets":
		return "secret", nil
	case "serviceaccount":
		return "sa", nil
	case "imagestream":
		return "is", nil
	case "service":
		return "svc", nil
	case "configmap":
		return "configmap", nil
	case "route", "routes":
		return "route", nil
	}
	return "", fmt.Errorf("unsupported openshift object type: %s", openShiftObjectType)
}

func GetConfigurationByOpenShiftObjectType(microServiceName string, openShiftObjectType string) (*dvevaluation.DvVariable, error) {
	switch openShiftObjectType {
	case "DeploymentConfig":
		return GetLiveDeploymentConfiguration(microServiceName)
	case "Secret":
	case "ServiceAccount":
	case "ImageStream":
	case "Service":
	case "ConfigMap":
	case "Route":
	}
	return nil, fmt.Errorf("Unimplemented configuration getter for %s in %s", openShiftObjectType, microServiceName)
}

func GetConfigurationByOpenShiftObjectTypeAndName(objectName string, openShiftObjectType string) (*dvevaluation.DvVariable, error) {
	shortName, err := GetShortOpenShiftNameForObjectType(openShiftObjectType)
	if err != nil {
		return nil, err
	}
	return GetLiveConfiguration("edit " + shortName + " " + objectName)
}

func GetMicroServiceFullList() ([]string, error) {
	return GetObjectFullList("dc")
}

func GetObjectFullListByObjectType(openShiftObjectType string) ([]string, error) {
	shortName, err := GetShortOpenShiftNameForObjectType(openShiftObjectType)
	if err != nil {
		return nil, err
	}
	return GetObjectFullList(shortName)
}

func GetObjectFullList(shortName string) ([]string, error) {
	info, err := RunOCCommandOrCache("get " + shortName)
	if err != nil {
		return nil, err
	}
	items := strings.Split(info, "\n")
	n := len(items)
	res := make([]string, n-1)
	k := 0
	for i := 1; i < n; i++ {
		s := dvtextutils.GetNextWordInText(items[i])
		if s != "" {
			res[k] = s
			k++
		}
	}
	return res[:k], nil
}

func GetMicroServiceDeploymentConfigs(microServiceName string, notCritical bool) ([]string, error) {
	cmdLine := "edit dc " + microServiceName
	info, err := GetKubernetesConfigurationPart(cmdLine, "Descriptor config for "+microServiceName, 0, ConfigTransformerDc, notCritical)
	if err != nil {
		return nil, err
	}
	return []string{info}, nil
}

func GetMicroServiceServices(microServiceName string) ([]string, error) {
	info, err := RunOCCommandOrCache("describe svc")
	if err != nil {
		return nil, err
	}
	lookup := "name=" + microServiceName
	res := make([]string, 0, 1)
	pos := strings.Index(info, lookup)
	for pos >= 0 {
		start := pos - 400
		if start < 0 {
			start = 0
		}
		s := info[start:pos]
		namePos := strings.LastIndex(s, "Name:")
		if namePos < 0 {
			return nil, errors.New("Corrupted structure of service description " + s + lookup)
		}
		name := dvtextutils.GetNextWordInText(s[namePos+5:])
		if name == "" {
			return nil, errors.New("Corrupted structure of service description " + s + lookup)
		}
		res = append(res, name)
		newPos := strings.Index(info[pos+5:], lookup)
		if newPos < 0 {
			break
		}
		pos += newPos + 5
	}
	return res, nil
}

func GetMicroServiceServiceConfigs(services []string, notCritical bool) ([]string, error) {
	n := len(services)
	res := make([]string, n)
	for i := 0; i < n; i++ {
		cmdLine := "edit svc " + services[i]
		info, err := GetKubernetesConfigurationPart(cmdLine, "Service config for "+services[i], 0, ConfigTransformerSvc, notCritical)
		if err != nil {
			return nil, err
		}
		res[i] = info
	}
	return res, nil
}

func GetMicroServiceRoutes(services []string) ([]string, error) {
	info, err := RunOCCommandOrCache("get route")
	if err != nil {
		return nil, err
	}
	items := strings.Split(info, "\n")
	n := len(items)
	m := len(services)
	res := make([]string, 0, m)
	for i := 0; i < n; i++ {
		s := dvtextutils.ConvertToNonEmptyList(items[i])
		if len(s) > 2 {
			svc := s[2]
			for j := 0; j < m; j++ {
				if services[j] == svc {
					res = append(res, s[0])
					break
				}
			}
		}
	}
	return res, nil
}

func GetMicroServiceRouteConfigs(routes []string, notCritical bool) ([]string, error) {
	n := len(routes)
	res := make([]string, n)
	for i := 0; i < n; i++ {
		cmdLine := "describe route " + routes[i]
		info, err := GetKubernetesConfigurationPart(cmdLine, "Route config for "+routes[i], 0, ConfigTransformerRoute, notCritical)
		if err != nil {
			return nil, err
		}
		res[i] = info
	}
	return res, nil
}

func GetMicroServiceConfigMaps(microServiceName string, notCritical bool) ([]string, []string, error) {
	name := microServiceName + ".monitoring-config"
	cmdLine := "describe configmap " + name
	info, err := GetKubernetesConfigurationPart(cmdLine, name, optionSilentForNotFoundMicroService, ConfigTransformerConfigMap, notCritical)
	if err != nil {
		return nil, nil, nil
	}
	return []string{name}, []string{info}, nil
}

func GetExistingFullOpenShiftTemplate(microServiceName string, notCritical bool) (deployment string, deleteInfo []string, services []string, routes []string, configMaps []string, errFinal error) {
	start, end, replacer := GetStartEndPartsOfGeneralTemplate(microServiceName)
	deleteInfo = make([]string, 0, 5)
	list := make([]string, 0, 5)
	deploymentConfigs, err := GetMicroServiceDeploymentConfigs(microServiceName, notCritical)
	if err != nil {
		errFinal = err
	} else {
		list = append(list, deploymentConfigs...)
	}
	deleteInfo = append(deleteInfo, "delete all -l name="+microServiceName)
	services, err = GetMicroServiceServices(microServiceName)
	if err != nil {
		if errFinal == nil {
			errFinal = err
		}
		return
	}
	for _, service := range services {
		deleteInfo = append(deleteInfo, "delete svc "+service)
	}
	svcConfigs, err := GetMicroServiceServiceConfigs(services, notCritical)
	if err != nil {
		if errFinal == nil {
			errFinal = err
		}
	} else {
		list = append(list, svcConfigs...)
	}
	routes, err = GetMicroServiceRoutes(services)
	if err != nil {
		if errFinal == nil {
			errFinal = err
		}
	} else {
		for _, route := range routes {
			deleteInfo = append(deleteInfo, "delete routes "+route)
		}
		routeConfigs, err := GetMicroServiceRouteConfigs(routes, notCritical)
		if err != nil {
			if errFinal == nil {
				errFinal = err
			}
		} else {
			list = append(list, routeConfigs...)
		}
	}
	configMaps, configMapConfigs, err := GetMicroServiceConfigMaps(microServiceName, notCritical)
	if err != nil {
		if errFinal == nil {
			errFinal = err
		}
		return
	}
	for _, configMap := range configMaps {
		deleteInfo = append(deleteInfo, "delete configmap "+configMap)
	}
	list = append(list, configMapConfigs...)
	n := len(list)
	for i := 0; i < n; i++ {
		s := list[i]
		for k, v := range replacer {
			s = strings.ReplaceAll(s, k, v)
		}
		if i != 0 {
			start += ","
		}
		start += "\n     " + s
	}
	deploymentBytes, err := dvjson.ReformatJson([]byte(start+end), 4, [][]byte{[]byte("(MISSING)")}, 2)
	if err != nil {
		if errFinal == nil {
			errFinal = err
		}
		return
	}
	deployment = string(deploymentBytes)
	return
}

func OrderTemplateObjectsByDependencies(objects []*dvevaluation.DvVariable, silent bool) (res []*dvevaluation.DvVariable, err error) {
	n := len(objects)
	res = make([]*dvevaluation.DvVariable, 0, n)
	for i := 0; i < n; i++ {
		p := objects[i]
		if p == nil {
			continue
		}
		v, ok := TemplateObjectPriority[strings.ToLower(p.ReadSimpleChildValue("kind"))]
		if !ok {
			if silent {
				dvlog.Printf("Unknown openshift object type %v", "Unknown openshift object type %v", p)
			} else {
				err = fmt.Errorf("Unknown openshift object type %v", p)
				return
			}
			continue
		}
		p.Extra.(*dvjson.DvFieldInfoExtra).FieldStatus = v
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Extra.(*dvjson.DvFieldInfoExtra).FieldStatus < res[j].Extra.(*dvjson.DvFieldInfoExtra).FieldStatus
	})
	return
}

func ResolveMostSimilarObjectByMicroserviceNameAndObjectType(microServiceName string, objectType string) (name string, ok bool) {
	list, err := GetObjectFullListByObjectType(objectType)
	if err != nil {
		return "", false
	}
	choice := ""
	choiceRate := 0
	n := len(list)
	for i := 0; i < n; i++ {
		s := list[i]
		if s == microServiceName {
			return s, true
		}
		rate := dvtextutils.EvaluateDifferenceRate(s, microServiceName)
		if rate > choiceRate {
			choiceRate = rate
			choice = s
		}
	}
	return choice, choiceRate > 50
}
