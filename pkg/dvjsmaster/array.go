/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func array_init() {
	dvevaluation.ArrayMaster.Prototype = &dvevaluation.DvVariable{
		Fields: map[string]*dvevaluation.DvVariable{
			"push": {
				Kind: dvevaluation.FIELD_FUNCTION,
				Fn:   Array_push,
			},
			"slice": {
				Kind: dvevaluation.FIELD_FUNCTION,
				Fn:   Array_slice,
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func Array_push(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return nil, nil
}

func Array_slice(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return nil, nil
}
