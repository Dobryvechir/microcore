/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package main

import (
	_ "github.com/Dobryvechir/microcore/pkg/dvjsmaster"
	_ "github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

func testCurrent() {
	initEnvironment()
	testEvaluationSingle("", "B=Array.of(1,2,3,4);B[2]", "3", KindANY)
	showResume()
}

func testInit() {
	err := dvparser.ReadPropertiesFileWithEnvironmentVariablesInCurrentDirectory("test.properties")
	if err != nil {
		log.Panicf("Cannot read properties %v", err)
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
