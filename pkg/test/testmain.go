/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package main

import (
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

func testCurrent() {
	initEnvironment()
	testEvaluationSingle( "i=2", "i<10 ? i<7? 1: 2 : 3", "1", KindInteger)
	showResume()
}

func testInit() {
	err:=dvparser.ReadPropertiesFileWithEnvironmentVariablesInCurrentDirectory("test.properties")
	if err!=nil {
		log.Panicf("Cannot read properties %v",err)
	}
}

func testAll() {
	testEvaluation()
}

func main() {
	testInit()
	testCurrent()
	testAll()
}
