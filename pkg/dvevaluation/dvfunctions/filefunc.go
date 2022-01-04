/***********************************************************************
MicroCore
Copyright 2020-2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvfunctions

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"io/ioutil"
)

func FileRead(context *dvgrammar.ExpressionContext,sels interface{}, args []interface{}) (interface{}, error) {
	n := len(args)
	var fileName string
	if n > 0 {
		fileName = dvevaluation.AnyToString(args[0])
	}
	if fileName == "" {
		return nil, errors.New("File name is not specified")
	}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

var filePrototypeMap = map[string]interface{}{
	"read": dvevaluation.NewDvFunction("read", []string{"fileName"}, FileRead),
}

var fileObject = dvevaluation.NewObject(filePrototypeMap)
