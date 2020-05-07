package dvlicense

import (
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
)

func GetLicense(request *dvmeta.RequestContext) bool {
	dvcom.HandleFromString(request, "GetLicense")
	return true
}

var licenseRegistrationConfig *dvmodules.RegistrationConfig = &dvmodules.RegistrationConfig{
	Name:            "license",
	EndPointHandler: GetLicense,
	//GlobalInitHandler: MethodGlobalInitHandler
	//ServerInitHandler: MethodServerInitHandler
}

func Init() bool {
	return dvmodules.RegisterModule(licenseRegistrationConfig, false)
}

var inited = Init()
