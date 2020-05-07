/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strings"
)

type DvFunc func(interface{}, []interface{}) (interface{}, error)
type DvFunction struct {
	Name string
	Args []string
	Fn   DvFunc
}

var dvFunctionPrototypeMap = map[string]interface{}{
	"bind":  buildDvFunction("bind", []string{"this", "args"}, BindDvFunction, nil),
	"call":  buildDvFunction("call", []string{"this", "args"}, CallDvFunction, nil),
	"apply": buildDvFunction("apply", []string{"this", "args"}, ApplyDvFunction, nil),
}

var DvFunctionPrototype *DvObject

func BindDvFunction(self interface{}, args []interface{}) (interface{}, error) {
	funcObj, ok := self.(*DvObject)
	if !ok || !IsFunction(funcObj) {
		return nil, errors.New(" is not a function in bind")
	}
	if len(args) == 0 {
		return funcObj, nil
	}
	n := len(args) - 1
	var proc = func(self1 interface{}, args1 []interface{}) (interface{}, error) {
		l := len(args1)
		argsTotal := make([]interface{}, n+l)
		for i := 0; i < n; i++ {
			argsTotal[i] = args[i+1]
		}
		for i := 0; i < l; i++ {
			argsTotal[i+n] = args1[i]
		}
		return ExecFunction(funcObj, args[0], argsTotal)
	}
	newArgs := GetDvFunctionArguments(funcObj)
	if len(newArgs) <= n {
		newArgs = nil
	} else {
		newArgs = newArgs[n:]
	}
	name, okey := funcObj.Get("name")
	if !okey {
		name = "bound"
	}
	return NewDvFunction(AnyToString(name), newArgs, proc), nil
}

func GetDvFunctionArguments(funcObj *DvObject) []string {
	args, ok := funcObj.Get("arguments")
	if !ok || args == nil || args == "" {
		return nil
	}
	return strings.Split(AnyToString(args), ",")
}

func buildDvFunction(name string, args []string, fn DvFunc, proto *DvObject) *DvObject {
	value := &DvFunction{
		Name: name,
		Args: args,
		Fn:   fn,
	}
	props := map[string]interface{}{
		"name":      name,
		"length":    len(args),
		"arguments": strings.Join(args, ","),
	}
	obj := NewDvObjectWithSpecialValues(value, dvgrammar.TYPE_FUNCTION, proto, props)
	return obj
}

func NewDvFunction(name string, args []string, fn DvFunc) *DvObject {
	return buildDvFunction(name, args, fn, DvFunctionPrototype)
}

func (fn *DvFunction) ToString() string {
	res := "function " + fn.Name + "("
	n := len(fn.Args)
	for i := 0; i < n; i++ {
		if i == 0 {
			res += fn.Args[i]
		} else {
			res += "," + fn.Args[i]
		}
	}
	return res + ")"
}

func IsFunction(v interface{}) bool {
	switch v.(type) {
	case *DvFunction:
		return true
	case *DvObject:
		return IsFunction(v.(*DvObject).Value)
	}
	return false
}

func ExecFunction(fn interface{}, self interface{}, arguments []interface{}) (interface{}, error) {
	switch fn.(type) {
	case *DvFunction:
		dvFunction, _ := fn.(DvFunction)
		v, err := dvFunction.Fn(self, arguments)
		if err != nil {
			return v, errors.New(err.Error() + " in " + dvFunction.ToString())
		}
		return v, nil
	case *DvObject:
		return ExecFunction(fn.(*DvObject).Value, self, arguments)
	}
	return nil, errors.New(" is not a function")
}

func CallDvFunction(self interface{}, args []interface{}) (interface{}, error) {
	//TODO
	return nil, nil
}

func ApplyDvFunction(self interface{}, args []interface{}) (interface{}, error) {
	//TODO
	return nil, nil
}

func dvfunction_init() {
	DvFunctionPrototype = NewObject(dvFunctionPrototypeMap)
}
