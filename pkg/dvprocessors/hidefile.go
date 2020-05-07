package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
)

func hideFileHandler(request *dvmeta.RequestContext) bool {
	dvcom.HandleError(request, "404 File Not Found")
	return true
}

var hideFileConfig *RegistrationConfig = &RegistrationConfig{
	Name:            "hidefile",
	EndPointHandler: hideFileHandler,
}

var hideFileInited bool = RegisterProcessor(hideFileConfig, true)
