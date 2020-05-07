/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func Init() {
	object_init()
	array_init()
	dvevaluation.Init()
	net_init()
	math_init()
	date_init()
	json_init()
}
