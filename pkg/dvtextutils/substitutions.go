/***********************************************************************
MicroCore
Copyright 2020 -2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

import (
	"bytes"
	"fmt"
	"strings"
)

var SubstitutionStartParameter = []byte{'$', '{'}
var SubstitutionEndParameter = []byte{'}'}

func SubstitutionGetList(data []byte) (list []string, rest []string, ok bool) {
	m := bytes.Index(data, SubstitutionStartParameter)
	if m < 0 {
		return
	}
	sub := data[:m]
	data = data[m+2:]
	m = bytes.Index(data, SubstitutionEndParameter)
	if m < 0 {
		return
	}
	ok = true
	list = make([]string, 0, 1)
	rest = make([]string, 0, 2)
	for m >= 0 {
		listValue := string(data[:m])
		data = data[m+1:]
		if len(listValue) > 0 && listValue[0] == '{' && len(data) > 0 && data[0] == '}' {
			listValue = listValue[1:]
			data = data[1:]
		}
		rest = append(rest, string(sub))
		list = append(list, listValue)
		m = bytes.Index(data, SubstitutionStartParameter)
		if m < 0 {
			break
		}
		sub = data[:m]
		subData := data[m+2:]
		m = bytes.Index(subData, SubstitutionEndParameter)
		if m >= 0 {
			data = subData
		}
	}
	rest = append(rest, string(data))
	return
}

func MatchParametersByVariableSubstitution(model string, list []string, rest []string, occurrence int, strict bool) (params map[string]string, doubt bool, ok bool) {
	n := len(rest)
	modelLen := len(model)
	if n < 2 {
		return
	}
	start := len(rest[0])
	if start > 0 && !strings.HasPrefix(model, rest[0]) {
		return
	}
	end := modelLen - len(rest[n-1])
	if end < modelLen && !strings.HasSuffix(model, rest[n-1]) {
		return
	}
	restLen := n - 2
	model = model[start:end]
	if restLen == 0 {
		params = make(map[string]string)
		params[list[0]] = model
		ok = true
		return
	}
	positions := make([]int, restLen)
	rest = rest[1 : n-1]
	if occurrence >= 0 {
		prev := 0
		for i := 0; i < restLen; i++ {
			ind := strings.Index(model[prev:], rest[i]) + prev
			if ind < prev {
				return
			}
			positions[i] = ind
			prev = ind + len(rest[i])
		}
	} else {
		// looking from the last position
		ending := len(model)
		for i := restLen - 1; i >= 0; i-- {
			ind := strings.LastIndex(model[:ending], rest[i])
			if ind < 0 {
				return
			}
			positions[i] = ind
			ending = ind
		}
	}
	ok = true
	positionCopy := make([]int, restLen)
	if occurrence == 0 || occurrence == 1 {
		params = SubstitutionsGetParamsByPositions(list, rest, positions, model)
		if occurrence == 0 {
			if SubstitutionsGetNextPosition(list, rest, positions, model) {
				doubt = true
			}
		}
	} else if occurrence > 0 {
		for occurrence--; occurrence > 0; occurrence-- {
			for i := 0; i < restLen; i++ {
				positionCopy[i] = positions[i]
			}
			if !SubstitutionsGetNextPosition(list, rest, positions, model) {
				if strict {
					ok = false
				} else {
					params = SubstitutionsGetParamsByPositions(list, rest, positionCopy, model)
				}
				return
			}
		}
		params = SubstitutionsGetParamsByPositions(list, rest, positions, model)
	} else {
		for occurrence++; occurrence < 0; occurrence++ {
			for i := 0; i < restLen; i++ {
				positionCopy[i] = positions[i]
			}
			if !SubstitutionsGetNextPositionFromLast(list, rest, positions, model) {
				if strict {
					ok = false
				} else {
					params = SubstitutionsGetParamsByPositions(list, rest, positionCopy, model)
				}
				return
			}
		}
		params = SubstitutionsGetParamsByPositions(list, rest, positions, model)
	}
	return
}

func SubstitutionsGetParamsByPositions(list []string, rest []string, positions []int, model string) map[string]string {
	pos := 0
	params := make(map[string]string)
	n := len(list)
	maxPos := len(model)
	var endPos, nextPos int
	for i := 0; i < n; i++ {
		endPos = maxPos
		if i != n-1 {
			endPos = positions[i]
			nextPos = endPos + len(rest[i])
		}
		params[list[i]] = model[pos:endPos]
		pos = nextPos
	}
	return params
}

func SubstitutionsGetNextPosition(list []string, rest []string, positions []int, model string) bool {
	n := len(rest)
	for i := n - 1; i >= 0; i-- {
		pos := positions[i] + 1
		j := i
		for ; j < n; j++ {
			newPos := strings.Index(model[pos:], rest[j]) + pos
			if newPos < pos {
				break
			}
			positions[j] = newPos
			pos = newPos + len(rest[j])
		}
		if j == n {
			return true
		}
	}
	return false
}

func SubstitutionsGetNextPositionFromLast(list []string, rest []string, positions []int, model string) bool {
	//TODO rewrite for last
	n := len(rest)
	for i := 0; i < n; i++ {
		pos := positions[i]
		j := i
		for ; j >= 0; j-- {
			newPos := strings.LastIndex(model[:pos], rest[j])
			if newPos < 0 {
				break
			}
			positions[j] = newPos
			pos = newPos
		}
		if j < 0 {
			return true
		}
	}
	return false
}

func SubstitutionMatchModelByPattern(model string, pattern string, hintOccurrence int, hintStrict bool) (params map[string]string, doubt bool, ok bool) {
	list, rest, ok := SubstitutionGetList([]byte(pattern))
	if !ok {
		return
	}
	params, doubt, ok = MatchParametersByVariableSubstitution(model, list, rest, hintOccurrence, hintStrict)
	return
}

func UpdateModelByParamGetter(pattern string, paramGetter func(string) (string, bool)) (string, []string) {
	list, rest, ok := SubstitutionGetList([]byte(pattern))
	if !ok {
		return pattern, nil
	}
	paramMap := make(map[string]string)
	n := len(list)
	var unfoundParams []string
	for i := 0; i < n; i++ {
		k, ok := paramGetter(list[i])
		if ok {
			paramMap[list[i]] = k
		} else {
			if unfoundParams == nil {
				unfoundParams = make([]string, 1)
				unfoundParams[0] = list[i]
			} else {
				unfoundParams = append(unfoundParams, list[i])
			}
		}
	}
	model, _ := UpdateModelByListRestParams(pattern, list, rest, paramMap, true)
	return model, unfoundParams
}

func UpdateModelByListRestParams(pattern string, list []string, rest []string, paramMap map[string]string, silent bool) (string, error) {
	model := rest[0]
	n := len(list)
	for i := 0; i < n; i++ {
		k, ok := paramMap[list[i]]
		if !ok {
			if !silent {
				return pattern, fmt.Errorf("Unknown parameter %s in %s", list[i], pattern)
			}
			k = "${" + list[i] + "}"
		}
		model += k + rest[i+1]
	}
	return model, nil
}

func UpdateModelByParams(pattern string, paramMap map[string]string) (string, error) {
	list, rest, ok := SubstitutionGetList([]byte(pattern))
	if !ok {
		return pattern, nil
	}
	return UpdateModelByListRestParams(pattern, list, rest, paramMap, false)
}

func UpdateModelByParamsSilent(pattern string, paramMap map[string]string) (string, error) {
	list, rest, ok := SubstitutionGetList([]byte(pattern))
	if !ok {
		return pattern, nil
	}
	return UpdateModelByListRestParams(pattern, list, rest, paramMap, true)
}

func GetUnquotedString(data string) string {
	n := len(data)
	if n < 2 || data[0] != '"' || data[n-1] != data[0] {
		return data
	}
	return data[1 : n-1]
}

func EvaluateDifferenceRate(model string, pattern string) int {
	n := len(pattern)
	if n == 0 {
		return 0
	}
	i := 0
	for ; i < n; i++ {
		if model[i] != pattern[i] {
			break
		}
	}
	rate := i * 100 / n
	increase := len(model) - len(pattern)
	if increase > 0 {
		rate -= increase
	}
	return rate
}

func ReduceListRestByKnownValueAtIndex(list []string, rest []string, index int, value string) ([]string, []string, int) {
	n := len(list)
	if index < 0 || index >= n {
		return list, rest, index
	}
	list = append(list[:index], list[index+1:]...)
	restIndex := index + 1
	rest[index] += value + rest[restIndex]
	rest = append(rest[:restIndex], rest[restIndex+1:]...)
	return list, rest, index - 1
}
