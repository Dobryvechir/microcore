/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
)

var crudRegistrationConfig *RegistrationConfig = &RegistrationConfig{
	Name:              "crud",
	GlobalInitHandler: dvjson.CrudGlobalInitialization,
	GenerateHandlers:  dvjson.CrudGenerateHandlers,
}

func crudInit() bool {
	dvevaluation.RegisterToStringConverter(dvjson.DvFieldInfoToStringConverter)
	return RegisterModule(crudRegistrationConfig, false)
}

var crudInited bool = crudInit()
