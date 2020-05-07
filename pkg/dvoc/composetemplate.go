/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

var templateGeneral string = `
{
    "kind": "Template",
    "apiVersion": "v1",
    "metadata": {
        "name": "{{MICROSERVICE}}",
        "annotations": {
            "openshift.io/display-name": "{{MICROSERVICE}}",
            "description": "Template for cloud service",
            "tags": "backend",
            "iconClass": "icon-php"
        }
    },
    "labels": {
          "template": "{{MICROSERVICE}}"
    },
    "parameters": [
        {
            "name": "SERVICE",
            "value": "{{MICROSERVICE}}",
            "description": "Service Name. For example: salvation-eternal-life",
            "required": false
        }
    ],
    "objects": [
        {{OBJECTS}}
    ]
}
`

const (
	templateGeneralObjects = "{{OBJECTS}}"
	templateDebugSign      = "salvation-eternal-life"
)

var template string = `
{
    "kind": "Template",
    "apiVersion": "v1",
    "metadata": {
        "name": "{{MICROSERVICE}}",
        "annotations": {
            "description": "Template for cloud service",
            "tags": "backend",
            "iconClass": "icon-php"
        }
    },
    "parameters": [
        {
            "name": "SERVICE",
            "value": "{{MICROSERVICE}}",
            "description": "Service Name. For example: salvation-eternal-life",
            "required": false
        },
        {
            "name": "BRANCH",
            "description": "Which git application development branch should be used to deploy",
            "value": "master",
            "required": false
        },
        {
            "name": "TAG",
            "description": "Which docker image tag should be used to deploy",
            "value": "latest",
            "required": false
        },
        {
            "name": "OPENSHIFT_SERVICE_NAME",
            "value": "{{OPENSHIFT_SERVICE_NAME}}",
            "description": "Service Name. For example: itl-com",
            "required": false
        },
        {
            "name": "PUBLIC_GATEWAY_URL",
            "description": "Frontend Gateway endpoint url",
            "value": "http://public-gateway-{{OPENSHIFT_PROJECT}}",
            "required": true
        },
        {
            "name": "PRIVATE_GATEWAY_URL",
            "description": "Frontend Gateway endpoint url",
            "value": "http://private-gateway-{{OPENSHIFT_PROJECT}}",
            "required": true
        },
        {
            "name": "PUBLIC_IDENTITY_PROVIDER_URL",
            "description": "Identity Provider endpoint url",
            "value": "http://public-gateway-{{OPENSHIFT_PROJECT}}/api/v1/identity-provider",
            "required": true
        },
        {
            "name": "CERTIFICATE_BUNDLE_MD5SUM",
            "value": "d41d8cd98f00b204e9800998ecf8427e",
            "description": "SSL secret name",
            "required": false
        },
        {
            "name": "SSL_SECRET",
            "value": "defaultsslcertificate",
            "description": "SSL secret name",
            "required": false
        }
    ],
    "objects": [
        {
            "kind": "DeploymentConfig",
            "apiVersion": "v1",
            "metadata": {
                "name": "${SERVICE}",
                "labels": {
                    "name": "${SERVICE}"
                }
            },
            "spec": {
                "replicas": 1,
                "strategy": {
                    "type": "Rolling",
                    "rollingParams": {
                        "updatePeriodSeconds": 1,
                        "intervalSeconds": 1,
                        "timeoutSeconds": 600,
                        "maxUnavailable": "25%",
                        "maxSurge": "25%"
                    }
                },
                "template": {
                    "metadata": {
                        "labels": {
                            "name": "${SERVICE}"
                        }
                    },
                    "spec": {
                        "volumes": [
                            {
                                "name": "${SSL_SECRET}",
                                "secret": {
                                    "secretName": "${SSL_SECRET}"
                                }
                            }
                        ],
                        "containers": [
                            {
                                "name": "${SERVICE}",
                                "image": "{{TEMPLATE_IMAGE}}",
                                "volumeMounts": [
                                    {
                                        "name": "${SSL_SECRET}",
                                        "mountPath": "/tmp/cert/${SSL_SECRET}"
                                    }
                                ],
                                "ports": [
                                    {
                                        "containerPort": 8080,
                                        "protocol": "TCP"
                                    }
                                ],
                                "env": [
                                    {
                                        "name": "CERTIFICATE_BUNDLE_${SSL_SECRET}_MD5SUM",
                                        "value": "${CERTIFICATE_BUNDLE_MD5SUM}"
                                    },
                                    {
                                        "name": "PUBLIC_GATEWAY_URL",
                                        "value": "${PUBLIC_GATEWAY_URL}"
                                    },
                                    {
                                        "name": "PRIVATE_GATEWAY_URL",
                                        "value": "${PRIVATE_GATEWAY_URL}"
                                    },
                                    {
                                        "name": "PUBLIC_IDENTITY_PROVIDER_URL",
                                        "value": "${PUBLIC_IDENTITY_PROVIDER_URL}"
                                    },
                                    {
                                        "name": "GLOWROOT_CLUSTER",
                                        "value": "${GLOWROOT_CLUSTER}"
                                    }
                                ],
                                "resources": {
                                    "requests": {
                                        "cpu": "100m",
                                        "memory": "32Mi"
                                    },
                                    "limits": {
                                        "memory": "32Mi",
                                        "cpu": "4"
                                    }
                                }
                            }
                        ]
                    }
                },
                "triggers": [
                    {
                        "type": "ConfigChange"
                    }
                ]
            }
        },
        {
            "kind": "Service",
            "apiVersion": "v1",
            "metadata": {
                "name": "${OPENSHIFT_SERVICE_NAME}",
                "annotations": {
          			{{ANNOTATIONS}}
        		}
            },
            "spec": {
                "ports": [
                    {
                        "name": "web",
                        "port": 8080,
                        "targetPort": 8080
                    }
                ],
                "selector": {
                    "name": "${SERVICE}"
                }
            }
        },
        {
            "apiVersion": "v1",
            "kind": "Route",
            "metadata": {
                "name": "{{OPENSHIFT_ROUTE_NAME}}"
            },
            "spec": {
                "to": {
                    "kind": "Service",
                    "name": "${OPENSHIFT_SERVICE_NAME}"
                }
            }
        }
		{{MICROSERVICE_CONFIG_MAP}}
    ]
}
`

var templateConfigMap = `
		,
        {
            "kind": "ConfigMap",
            "apiVersion": "v1",
            "metadata": {
                "name": "${SERVICE}.monitoring-config"
            },
            "data": {
                "url.health": "http://%(ip)s:8080/health"
            }
        }
`

const templateRequired = "??????????????????????"

var templateComposeDefaults = map[string]string{
	"MICROSERVICE":            templateRequired,
	"OPENSHIFT_SERVICE_NAME":  templateRequired,
	"OPENSHIFT_ROUTE_NAME":    templateRequired,
	"TEMPLATE_IMAGE":          templateRequired,
	"OPENSHIFT_PROJECT":       templateRequired,
	"ANNOTATIONS":             "",
	"MICROSERVICE_CONFIG_MAP": "",
}

func ComposeOpenShiftJsonTemplateBySample(sample string, requiredParams map[string]string, params map[string]string) ([]byte, error) {
	r := strings.TrimSpace(sample)
	for k, v := range requiredParams {
		value := params[k]
		if value == "" {
			value = dvparser.GlobalProperties[k]
		}
		if value == "" {
			if v == templateRequired {
				return nil, errors.New("Required parameter " + k + " was not provided for " + params[MicroServiceProperty])
			}
		}
		if k == "MICROSERVICE_CONFIG_MAP" && value == "true" {
			value = templateConfigMap
		}
		r = strings.ReplaceAll(r, "{{"+k+"}}", value)
	}
	return []byte(r), nil
}

func ComposeOpenShiftJsonTemplate(params map[string]string) ([]byte, error) {
	return ComposeOpenShiftJsonTemplateBySample(template, templateComposeDefaults, params)
}

func GetStartEndPartsOfGeneralTemplate(microServiceName string) (string, string, map[string]string) {
	r := strings.TrimSpace(strings.ReplaceAll(templateGeneral, "{{MICROSERVICE}}", microServiceName))
	pos := strings.Index(r, templateGeneralObjects)
	replaceMap := map[string]string{
		"\"name\": \"" + microServiceName + "\"": "\"name\": \"${SERVICE}\"",
		"\"name\":\"" + microServiceName + "\"":  "\"name\": \"${SERVICE}\"",
	}
	return r[:pos], r[pos+len(templateGeneralObjects):], replaceMap
}
