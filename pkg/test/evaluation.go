/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package main

import (
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"log"
	"strings"
	"time"
)

const (
	KindANY        = 'A'
	KindInteger    = 'I'
	KindString     = 'S'
	KindFloat      = 'F'
	KindBoolean    = 'B'
	KindDvVariable = 'V'
)

var env *dvevaluation.DvObject
var tested int
var successful int
var startTime time.Time

func evaluateKind(v interface{}, expected byte) byte {
	switch v.(type) {
	case int64:
		return KindInteger
	}
	return expected
}

func testEvaluationSingle(vars string, expr string, result string, resultKind byte) {
	putVarsToEnvironment(vars)
	tested++
	res, err := env.EvaluateAnyTypeExpression(expr)
	if err != nil {
		fmt.Printf("Error expr=[%s] var=[%s] exp=[%s] %v\n", expr, vars, result, err)
		return
	}
	s := dvevaluation.AnyToString(res)
	if s != result {
		fmt.Printf("Expected [%s] but [%s] expr=[%s] var=[%s]\n", result, s, expr, vars)
		return
	}
	if resultKind != KindANY {
		y := evaluateKind(res, resultKind)
		if y != resultKind && !(resultKind == KindFloat && y == KindInteger) {
			fmt.Printf("Wanted [%v] but [%v] expr=[%s] var=[%s] result=[%s]\n", resultKind, y, expr, vars, s)
			return
		}
	}
	successful++
}

func checkError(expr string, result string, mode int) {
	tested++
	_, err := env.EvaluateAnyTypeExpression(expr)
	if err == nil {
		fmt.Printf("Must be error [%s] but nothing detected in expr=[%s]\n", result, expr)
		return
	}
	s := err.Error()
	res := s == result
	if !res {
		if mode == 1 {
			res = strings.HasPrefix(s, result)
		}
	}
	if !res {
		fmt.Printf("Expected error [%s] but [%s] in expr=[%s]\n", result, s, expr)
		return
	}
	successful++
}

func checkErrorPref(expr string, result string) {
	checkError(expr, result, 1)
}

func initEnvironment() {
	tested = 0
	successful = 0
	env = dvparser.GetGlobalPropertiesAsDvObject()
	startTime = time.Now()
}

func showResume() {
	elapsed := time.Since(startTime)
	fmt.Printf("%d/%d success for %s\n", successful, tested, elapsed)
}

func putVarsToEnvironment(vars string) {
	vars = strings.TrimSpace(vars)
	if vars == "" {
		return
	}
	m := dvtextutils.LoadSimpleMapFromString(vars)
	if m == nil {
		return
	}
	for k, v := range m {
		if k != "" {
			var res interface{} = v
			var ok bool
			var err error
			kn := strings.ToLower(k)
			switch kn[0] {
			case 'i':
				res, ok = dvevaluation.AnyToNumberInt(v)
				if !ok {
					fmt.Printf("Error converting to number %s\n", v)
				}
			case 'f':
				res = dvevaluation.AnyToNumber(res)
			case 'v':
				res, err = dvjson.JsonFullParser([]byte(v))
				if err != nil {
					log.Panicf("Incorrect object %s: %v", v, err)
				}
			}
			env.Set(k, res)
		}
	}
}
