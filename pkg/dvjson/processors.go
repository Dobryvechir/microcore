/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"fmt"
	"strings"
)

type ProcessorFunction func(parent *DvFieldInfo, params string) (*DvFieldInfo, error)

var ProcessorFunctions = map[string]ProcessorFunction{
	"find": ProcessorFind,
}

func ExecuteProcessorFunction(name string, params string, parent *DvFieldInfo) (*DvFieldInfo, error) {
	name = strings.TrimSpace(name)
	pf, ok := ProcessorFunctions[name]
	if !ok {
		return nil, fmt.Errorf("Unknown function %s(%s)", name, params)
	}
	return pf(parent, params)
}
