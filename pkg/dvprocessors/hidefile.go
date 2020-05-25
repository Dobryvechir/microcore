package dvprocessors

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
)

func hideFileHandler(request *dvcontext.RequestContext) bool {
	request.HandleFileNotFound()
	return true
}

var hideFileConfig *RegistrationConfig = &RegistrationConfig{
	Name:            "hidefile",
	EndPointHandler: hideFileHandler,
}

var hideFileInited bool = RegisterProcessor(hideFileConfig, true)
