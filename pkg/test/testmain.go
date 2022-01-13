/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package main

import (
	_ "github.com/Dobryvechir/microcore/pkg/dvjsmaster"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
)

func testCurrent() {
	initEnvironment()
	testEvaluationSingle("S=Good-Word", "S.length","9", KindInteger)
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
