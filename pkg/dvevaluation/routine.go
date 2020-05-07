/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strconv"
	"strings"
)

type DvCall struct {
	Input  []string `json:"input"`
	Name   string   `json:"name"`
	Output string   `json:"output"`
}

type DvRoutine struct {
	Name         string   `json:"name"`
	Calls        []DvCall `json:"calls"`
	Dependencies []string `json:"dependencies"`
}

func (engine *DvScript) AddRoutines(routines []DvRoutine) error {
	if engine.routines == nil {
		engine.routines = make(map[string]DvRoutine)
	}
	for _, routine := range routines {
		if _, ok := engine.routines[routine.Name]; ok {
			return errors.New("Routine " + routine.Name + " is duplicated")
		}
		engine.routines[routine.Name] = routine
	}
	return nil
}

func (engine *DvScript) VerifyRoutines(routines []string) error {
	n := len(routines)
	for i := 0; i < n; i++ {
		if _, ok := engine.routines[routines[i]]; !ok {
			return errors.New("Routine " + routines[i] + " does not exist")
		}
	}
	return nil
}

func (context *DvContext) ExecuteRoutine(routine *DvRoutine) error {
	n := len(routine.Calls)
	if n == 0 {
		return nil
	}
	err := AssignIntToVariable(context.root, "process.current.block", 0, true)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At assigning int 0 to process.current.block in executing routine "+routine.Name)
		return err
	}
	err = AssignIntToVariable(context.root, "process.current.blocks", n, true)
	if err != nil {
		err = dvgrammar.EnrichErrorStr(err, "At assigning int "+strconv.Itoa(n)+" to process.current.blocks in executing routine "+routine.Name)
		return err
	}
	notYet := true
	step := 0
	for notYet {
		notYet, err = context.CurrentRoutineCallExecute(routine)
		if err != nil {
			err = dvgrammar.EnrichErrorStr(err, "At executing step "+strconv.Itoa(step)+" of "+strconv.Itoa(n)+" in routine "+routine.Name)

		}
		step++
	}
	return nil
}

func (context *DvContext) CurrentRoutineCallShift(delta int) (int, error) {
	return GetIntFromVariableAndModify(context.root, "process.current.block", delta, false, true)
}

func (context *DvContext) CurrentRoutineCallBreak() error {
	return AssignIntToVariable(context.root, "process.current.block", -1, true)
}

func (context *DvContext) CurrentRoutineCallExecute(routine *DvRoutine) (bool, error) {
	currentCall, _ := GetIntFromVariableAndIncrementAfter(context.root, "process.current.block", true)
	callNumber, _ := GetIntFromVariable(context.root, "process.current.blocks", true)
	if currentCall < 0 || currentCall >= callNumber {
		return false, nil
	}
	dvCall := &routine.Calls[currentCall]
	if dvCall == nil {
		return false, errors.New("Empty call in routine " + routine.Name)
	}
	return true, context.RoutineCallExecute(dvCall)
}

func (context *DvContext) RoutineCallExecute(dvCall *DvCall) error {
	params, err := QuickVariableArrayEvaluation(context.root, dvCall.Input)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "Before executing routine call "+dvCall.Name+" to be saved to "+dvCall.Output+" due to processing params "+strings.Join(dvCall.Input, ", "))
	}
	res, err := context.FunctionCallByVariableDefinition(context.root, dvCall.Name, params, nil)
	if err != nil {
		return dvgrammar.EnrichErrorStr(err, "At executing routine call "+dvCall.Name+" to be saved to "+dvCall.Output+" with processing params "+strings.Join(dvCall.Input, ", "))
	}
	if dvCall.Output != "" {
		err = AssignVariableToVariable(context.root, dvCall.Output, res, true)
		if err != nil {
			return dvgrammar.EnrichErrorStr(err, "After executing routine call "+dvCall.Name+" at saving to "+dvCall.Output+" with processing params "+strings.Join(dvCall.Input, ", "))
		}
	}
	return nil
}
