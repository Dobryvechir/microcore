/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func object_init() {
	dvevaluation.ObjectMaster.Prototype = &dvevaluation.DvVariable{
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("keys"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra:   Object_keys,
			},
			{
				Name: []byte("entries"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra:   Array_slice,
			},
		},
		Kind: dvevaluation.FIELD_OBJECT,
	}
}

func Object_keys(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return nil, nil
}

func Object_entries(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return nil, nil
}
