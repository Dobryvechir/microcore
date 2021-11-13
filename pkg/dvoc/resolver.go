/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

const (
	MicroserviceUrlProperty           = "MICROSERVICE_PATH_"
	ApplicationInsideCloud            = "APPLICATION_INSIDE_CLOUD"
	MicroServiceExposedRoute          = "http://{SERVICE}-{{{OPENSHIFT_NAMESPACE}}}.{{{OPENSHIFT_SERVER}}}.{{{OPENSHIFT_DOMAIN}}}"
	MicroServiceInternalRoute         = "http://{SERVICE}:8080"
	MicroServiceExposedRouteTemplate  = "MICROSERVICE_EXPOSED_ROUTE_TEMPLATE"
	MicroServiceInternalRouteTemplate = "MICROSERVICE_INTERNAL_ROUTE_TEMPLATE"
)

var resolvedMicroServiceUrls = make(map[string]string)

func GetUrlByGlobalPropertiesAndService(globalUrlTemplatePropertyName, defaultUrlTemplate, serviceName string) (string, error) {
	template := dvparser.GetByGlobalPropertiesOrDefault(globalUrlTemplatePropertyName, defaultUrlTemplate)
	urlTemplate, err := dvparser.ConvertStringByGlobalProperties(template, globalUrlTemplatePropertyName)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(urlTemplate, "{SERVICE}", serviceName), nil
}

func ResolveMicroServiceUrl(microServiceName string) (string, error) {
	url := resolvedMicroServiceUrls[microServiceName]
	if url != "" {
		return url, nil
	}
	url = dvparser.GlobalProperties[MicroserviceUrlProperty+dvaction.GetMicroServicePropertyName(microServiceName)]
	if url != "" {
		resolvedMicroServiceUrls[microServiceName] = url
		return url, nil
	}
	if dvparser.GlobalProperties[ApplicationInsideCloud] == "true" {
		url, err := GetUrlByGlobalPropertiesAndService(MicroServiceInternalRouteTemplate, MicroServiceInternalRoute, microServiceName)
		if err != nil {
			return "", err
		}
		resolvedMicroServiceUrls[microServiceName] = url
		return url, nil
	}
	url, err := GetUrlByGlobalPropertiesAndService(MicroServiceExposedRouteTemplate, MicroServiceExposedRoute, microServiceName)
	if err != nil {
		return "", err
	}
	resolvedMicroServiceUrls[microServiceName] = url
	return url, nil
}

func ResolveUrlTemplate(globalPropertiesTemplate string, defaultTemplate string) (string, error) {
	template := dvparser.GetByGlobalPropertiesOrDefault(globalPropertiesTemplate, defaultTemplate)
	if template == "" {
		return "", errors.New("Empty url template")
	}
	p := strings.Index(template, "}")
	if template[0] == '{' && p > 0 {
		service := template[1:p]
		if service == "" {
			return "", fmt.Errorf("Empty service name in url %s", template)
		}
		url, err := ResolveMicroServiceUrl(service)
		if err != nil {
			return "", err
		}
		return url + template[p+1:], nil
	}
	return template, nil
}

func ResolveUrlRequestByGlobalPropertiesAndDefaults(prefix string, defaultMethod string, defaultUrl string, defaultBody string, headers map[string]string, replacer map[string]string) (string, error) {
	url, err := ResolveUrlTemplate(prefix+"_URL_TEMPLATE", defaultUrl)
	if err != nil {
		return "", err
	}
	body := dvparser.GetByGlobalPropertiesOrDefault(prefix+"_BODY_TEMPLATE", defaultBody)
	method := dvparser.GetByGlobalPropertiesOrDefault(prefix+"_METHOD", defaultMethod)
	if replacer != nil {
		for k, v := range replacer {
			s := "{" + k + "}"
			url = strings.ReplaceAll(url, s, v)
			body = strings.ReplaceAll(body, s, v)
			method = strings.ReplaceAll(method, s, v)
		}
	}
	res, err := dvnet.NewJsonRequest(method, url, body, headers, dvnet.AveragePersistentOptions)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
