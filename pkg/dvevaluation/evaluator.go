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

var complexEvaluationSupported = true

type EvaluateResult struct {
	FinalResult interface{}
	Err         error
	EndPos      int
}

type EvaluateBoolean struct {
	FinalResult bool
	Err         error
}

var GlobalFunctionPool = make(map[string]interface{})
var GlobalFunctionPrototype = NewObject(GlobalFunctionPool)

func AddToGlobalFunctionPool(properties map[string]interface{}) {
	if properties == nil {
		return
	}
	for k, v := range properties {
		GlobalFunctionPool[k] = v
	}
}

func makeErrorMessage(main string, row int, column int, place string) error {
	if len(place) > 0 {
		place = " in " + place
	}
	return errors.New(main + " at row " + strconv.Itoa(row) + " column " + strconv.Itoa(column) + place)
}

func NewDvObjectFrom2Maps(localMap map[string]string, globalMap map[string]string) *DvObject {
	rootObj := NewDvObjectWithGlobalPrototype(globalMap)
	if localMap != nil {
		rootObj = NewDvObject(localMap, rootObj)
	}
	return rootObj
}

func Parse(data []byte, globalMap map[string]string, localMap map[string]string, row int, column int, place string) (r *EvaluateResult) {
	l := len(data)
	pos := 0
	for pos < l && data[pos] <= 32 {
		if data[pos] == 10 {
			row++
			column = 1
		} else {
			column++
		}
		pos++
	}
	for l > pos && data[l-1] <= 32 {
		l--
	}
	data = data[pos:l]
	if len(data) == 0 {
		r = &EvaluateResult{FinalResult: ""}
		return r
	}
	key := string(data)
	if localMap != nil {
		if v1, ok := localMap[key]; ok {
			return &EvaluateResult{FinalResult: v1}
		}
	}
	if v2, ok := globalMap[key]; ok {
		r = &EvaluateResult{FinalResult: v2}
		return r
	}
	if !complexEvaluationSupported {
		r = &EvaluateResult{Err: makeErrorMessage("Complex expressions are not parsed yet", row, column, place)}
	} else {
		r = ParseForDvObject(data, NewDvObjectFrom2Maps(localMap, globalMap), row, column, place)
	}
	return r
}

func ParseForDvObject(data []byte, params *DvObject, row int, column int, place string) *EvaluateResult {
	visitorOptions := 0
	ref := &dvgrammar.SourceReference{
		Row:    row,
		Column: column,
		Place:  place,
	}
	ev, err := CalculatorEvaluator(data, params, ref, visitorOptions)
	var res interface{}
	if err == nil {
		res = ev.Value
	}
	return &EvaluateResult{FinalResult: res, Err: err}
}

func EvalAsBoolean(data []byte, globalMap map[string]string, localMap map[string]string, row int, column int, place string) (r EvaluateBoolean) {
	eval := Parse(data, globalMap, localMap, row, column, place)
	if eval.Err != nil {
		r = EvaluateBoolean{Err: eval.Err}
		return
	}
	val := AnyToBoolean(eval.FinalResult)
	r = EvaluateBoolean{FinalResult: val}
	return
}

func FindEndAndParse(data []byte, pos int, limit int, sequenceLimit int, globalMap map[string]string, localMap map[string]string, row int, column int, place string) *EvaluateResult {
	sequence := 0
	for i := pos; i < limit; i++ {
		if data[i] == '}' {
			sequence++
			if sequence == sequenceLimit {
				posFinish := i - sequenceLimit + 1
				endPos := i + 1
				for pos < posFinish && data[pos] == '{' && endPos < limit && data[endPos] == '}' {
					pos++
					endPos++
				}
				r := Parse(data[pos:posFinish], globalMap, localMap, row, column, place)
				r.EndPos = endPos
				return r
			}
		} else {
			sequence = 0
		}
	}

	r := &EvaluateResult{Err: makeErrorMessage("No script end "+strings.Repeat("}", sequenceLimit), row, column, place), EndPos: pos}
	return r
}

func IsDefined(data []byte, scope dvgrammar.ScopeInterface, row int, col int, place string, visitorOptions int) (int, error) {
	end := len(data)
	start := 0
	for end > 0 && data[end-1] <= 32 {
		end--
	}
	for start < end && data[start] <= 32 {
		start++
	}
	if start < end {
		data = data[start:end]
		key := string(data)
		if _, ok := scope.Get(key); ok {
			if (visitorOptions & EVALUATE_OPTION_UNDEFINED) == 0 {
				return 1, nil
			}
			return 0, nil
		}
		reference := &dvgrammar.SourceReference{Row: row, Column: col, Place: place}
		calc, err := CalculateDefined(data, scope, reference, visitorOptions)
		if err != nil {
			return 0, err
		}
		if AnyToBoolean(calc.Value) {
			return 1, nil
		}
	}
	return 0, nil
}
