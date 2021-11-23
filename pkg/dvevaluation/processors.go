/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"fmt"
	"strings"
)

type ProcessorFunction func(parent *DvVariable, params string) (*DvVariable, error)
type ProcessByteArray func([]byte) (*DvVariable, error)

var ProcessorFunctions = map[string]ProcessorFunction{
	"find": ProcessorFind,
}
var JsonFullParser ProcessByteArray

func ExecuteProcessorFunction(name string, params string, parent *DvVariable) (*DvVariable, error) {
	name = strings.TrimSpace(name)
	pf, ok := ProcessorFunctions[name]
	if !ok {
		return nil, fmt.Errorf("Unknown function %s(%s)", name, params)
	}
	return pf(parent, params)
}

func RegisterJsonFullParser(fn ProcessByteArray) bool {
	JsonFullParser = fn
	return true
}
