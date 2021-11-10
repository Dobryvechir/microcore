/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"encoding/json"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

type ConnectionPropertiesInfo struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Url      string `json:"url"`
	DbName   string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type DbaasInfo struct {
	ConnectionProperties ConnectionPropertiesInfo `json:"connectionProperties"`
}

const (
	OcPostgreSql                          = "postgresql"
	OcMongoDb                             = "mongodb"
	DbaaSServerUrlTemplateProperty        = "DBAAS_SERVER"
	DbaaSServerUrlDefault                 = "{dbaas-agent}/api/v1/dbaas/{{{OPENSHIFT_NAMESPACE}}}/databases"
	DbaaSRequestMethodProperty            = "DBAAS_REQUEST_METHOD"
	DbaaSRequestMethod                    = "PUT"
	DbaaSRequestBodyTenantAwareProperty   = "DBAAS_REQUEST_TENANT_AWARE"
	DbaaSRequestBodyTenantUnawareProperty = "DBAAS_REQUEST_TENANT_UNAWARE"
	DbaaSRequestBodyTenantAwareDefault    = "{\"classifier\":{ {{{DBAAS_EXTRA_REQUEST_BODY}}}\"tenantId\":\"{TENANT_ID}\",\"microserviceName\":\"{SERVICE}\"},\"type\":\"{DBTYPE}\",\"namePrefix\":null,\"dbName\":null,\"username\":null,\"password\":null,\"physicalDatabaseId\":null,\"initScriptIdentifiers\":null,\"backupDisabled\":null,\"settings\":null}"
	DbaaSRequestBodyTenantUnawareDefault  = "{\"classifier\":{ {{{DBAAS_EXTRA_REQUEST_BODY}}}\"microserviceName\":\"{SERVICE}\"},\"type\":\"{DBTYPE}\",\"namePrefix\":null,\"dbName\":null,\"username\":null,\"password\":null,\"physicalDatabaseId\":null,\"initScriptIdentifiers\":null,\"backupDisabled\":null,\"settings\":null}"
)

func GetDbaasProperties(microServiceName string, m2mToken string, database string, tenantId string) (*DbaasInfo, error) {
	var ok bool
	if m2mToken == "" {
		m2mToken, ok = GetM2MToken(microServiceName)
		if !ok {
			return nil, fmt.Errorf("Cannot get M2M token for %s", microServiceName)
		}
	}
	url, err := ResolveUrlTemplate(DbaaSServerUrlTemplateProperty, DbaaSServerUrlDefault)
	if err != nil {
		return nil, err
	}
	body := ""
	if tenantId == "" {
		body = dvparser.GetByGlobalPropertiesOrDefault(DbaaSRequestBodyTenantUnawareProperty, DbaaSRequestBodyTenantUnawareDefault)
	} else {
		body = dvparser.GetByGlobalPropertiesOrDefault(DbaaSRequestBodyTenantAwareProperty, DbaaSRequestBodyTenantAwareDefault)
	}
	body = strings.ReplaceAll(body, "{SERVICE}", microServiceName)
	body = strings.ReplaceAll(body, "{TENANT_ID}", tenantId)
	body = strings.ReplaceAll(body, "{DBTYPE}", database)
	headers := map[string]string{"Authorization": m2mToken}
	if tenantId != "" {
		headers["Tenant"] = tenantId
	}
	method := dvparser.GetByGlobalPropertiesOrDefault(DbaaSRequestMethodProperty, DbaaSRequestMethod)
	res, err := dvnet.NewJsonRequest(method, url, body, headers, dvnet.AveragePersistentOptions)
	if err != nil {
		return nil, err
	}
	dbaasInfo := &DbaasInfo{}
	if err = json.Unmarshal(res, dbaasInfo); err != nil {
		return nil, err
	}
	return dbaasInfo, nil
}
