/***********************************************************************
MicroCore
Copyright 2020-2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
)

var (
	TenantResolvePropertyPrefix     = "TENANT_RESOLVER"
	TenantResolveUrlTemplateDefault = "{public-gateway}/api/v2/tenant-manager/registration/tenants?dns={TENANT}"
	TenantResolveMethodDefault      = "GET"
	TenantResolveBodyDefault        = ""
)

var resolvedTenantIds = make(map[string]string)

func ResolveTenantIdByTenant(tenant string) (string, error) {
	tenantId := resolvedTenantIds[tenant]
	tenantId, err := ResolveUrlRequestByGlobalPropertiesAndDefaults(
		TenantResolvePropertyPrefix,
		TenantResolveMethodDefault,
		TenantResolveUrlTemplateDefault,
		TenantResolveBodyDefault,
		nil,
		map[string]string{
			"TENANT": tenant,
		},
	)
	if err != nil {
		return "", err
	}
	return dvparser.GetUnquotedString(tenantId), nil
}

/***********************************************************************
	Functions to be used with command line only
************************************************************************/

const (
	TenantIdProperty = "TENANT_ID"
	TenantProperty   = "TENANT"
)

func EnsureTenantIdInGlobalProperties() bool {
	if dvparser.GlobalProperties[TenantIdProperty] != "" {
		return true
	}
	tenant := dvparser.GlobalProperties[TenantProperty]
	if tenant == "" {
		dvlog.PrintlnError("Error: TENANT property is not specified")
		return false
	}
	tenantId, err := ResolveTenantIdByTenant(tenant)
	if err != nil {
		dvlog.PrintlnError(err.Error())
		return false
	}
	dvparser.GlobalProperties[TenantIdProperty] = tenantId
	return true
}
