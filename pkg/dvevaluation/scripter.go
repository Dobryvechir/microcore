/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"log"
	"strconv"
)

type DvScript struct {
	routines map[string]DvRoutine
	master   *DvVariable
}

type DvContext struct {
	root     *DvVariable
	routines map[string]DvRoutine
}

var registeredMasterObjects map[string]*DvVariable = make(map[string]*DvVariable)

func ParseScripts(scripts []string) (*DvScript, error) {
	engine := &DvScript{master: DvVariableGetNewObject()}
	for k, v := range registeredMasterObjects {
		engine.master.Fields[k] = v
	}
	return engine, nil
}

func RegisterMasterVariable(name string, variable *DvVariable) *DvVariable {
	registeredMasterObjects[name] = variable
	return variable
}

func RegisterMasterObject(name string, values map[string]*DvVariable, functions map[string]DvvFunction) *DvVariable {
	masterObject := DvVariableGetNewObject()
	if values != nil {
		masterObject.Fields = values
	}
	prototypeObject := DvVariableGetNewObject()
	for k, v := range functions {
		prototypeObject.Fields[k] = ConvertDvFunctionToDvVariable(v)
	}
	masterObject.Prototype = prototypeObject
	registeredMasterObjects[name] = masterObject
	return masterObject
}

func (engine *DvScript) GetNewContext() *DvContext {
	//TODO make a deep copy of engine.master instead of assigning
	return &DvContext{root: engine.master, routines: engine.routines}
}

func (context *DvContext) DumpContextMemory() {
	if context.root == nil {
		log.Print("null")
	} else {
		buf := context.root.JsonStringify()
		log.Print(buf)
	}
}

func (context *DvContext) ExecuteRoutines(routines []string, params []string) error {
	err := AssignArrayStringToVariable(context.root, "process.argv", params, true)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "While assingning params to process.argv to execute routines")
	}
	err = AssignArrayStringToVariable(context.root, "process.routines", routines, true)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "While assingning routines to process.routines to execute them")
	}
	err = AssignIntToVariable(context.root, "process.current.routine", 0, true)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "While assingning int 0 to process.current.routine to execute routines")
	}
	notYet := true
	step := 0
	for notYet {
		notYet, err = context.ExecuteCurrentRoutine()
		if err != nil {
			return dvgrammar.EnrichErrorStr(err, "While executing routines at step "+strconv.Itoa(step))
		}
		step++
	}
	return nil
}

func (context *DvContext) ExecuteCurrentRoutine() (bool, error) {
	currentRoutine, _ := GetIntFromVariableAndIncrementAfter(context.root, "process.current.routine", true)
	routineNumber, _ := GetIntFromVariable(context.root, "process.routines.length", true)
	if currentRoutine >= routineNumber {
		return false, nil
	}
	routineNamePlace := "process.routines[" + strconv.Itoa(currentRoutine) + "]"
	routineName, err := GetStringFromVariable(context.root, routineNamePlace, false)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "While extracting routine name from "+routineNamePlace)
		return false, err
	}
	if routine, ok := context.routines[routineName]; ok {
		err = context.ExecuteRoutine(&routine)
		if err != nil {
			err = dvgrammar.EnrichErrorStr(err, "At executing routine["+strconv.Itoa(currentRoutine)+"] "+routineName)
		}
	} else {
		return false, errors.New("Routine " + routineName + " not found")
	}
	return true, err
}

func (context *DvContext) SetGeneralArguments(params map[string]string) error {
	err := AssignMapStringToVariable(context.root, "process.env", params, true)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At setting environment variables to process.env")
	}
	return err
}
