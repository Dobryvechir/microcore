/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import "strings"

const (
	UPDATE_MODE_REPLACE = iota
	UPDATE_MODE_APPEND
	UPDATE_MODE_ADD_BY_KEYS
	UPDATE_MODE_MERGE
)

func GetLengthOfAny(v interface{}) int {
	switch v.(type) {
	case *DvVariable:
		return v.(*DvVariable).GetLength()
	case string:
		return len(v.(string))
	case []string:
		return len(v.([]string))
	}
	return 0
}

func GetChildrenOfAnyByRange(src interface{}, startIndex int, endIndex int) interface{} {
	switch src.(type) {
	case *DvVariable:
		v := src.(*DvVariable).GetChildrenByRange(startIndex, endIndex-startIndex)
		return v
	case string:
		s := src.(string)[startIndex:endIndex]
		return s
	}
	return nil
}

func IsSimpleJsonPath(path string) bool {
	if strings.Contains(path, "[") {
		return false
	}
	if strings.Contains(path, "{") {
		return false
	}
	return true
}

func GetLatestKeyInJsonPath(path string) string {
	p1 := strings.LastIndex(path, ".")
	p2 := strings.LastIndex(path, "]")
	p3 := strings.LastIndex(path, "}")
	if p1 < p2 {
		p1 = p2
	}
	if p1 < p3 {
		p1 = p3
	}
	if p1 >= 0 {
		path = path[p1+1:]
	}
	return path
}

func UpdateAnyVariables(srcAny interface{}, valAny interface{}, path string, mode int,ids []string, env *DvObject) interface{} {
	src := AnyToDvVariable(srcAny)
	val := AnyToDvVariable(valAny)
	var parent *DvVariable = nil
	var node *DvVariable = src
	var err error = nil
	path=strings.Trim(path, ".")
	if path != "" {
		node, parent, err = src.ReadPath(path, true, env)
		if err != nil {
			if IsSimpleJsonPath(path) {
				src.InsertAtSimplePath(path, val)
			}
			return src
		}
	}
	if parent == nil {
		src.MergeOtherVariable(val, mode, ids)
	} else {
		n := parent.IndexOf(node)
		if n < 0 {
			key := GetLatestKeyInJsonPath(path)
			if key!="" {
				parent.InsertAtSimplePath(key, val)
			}
		} else {
			parent.MergeAtChild(n, val, mode)
		}
	}
	return src
}

func RemoveAnyVariable(srcAny interface{}, path string, env *DvObject) interface{} {
	path=strings.Trim(path, ".")
	if path == "" {
		return nil
	}
	src := AnyToDvVariable(srcAny)
	node, parent, err := src.ReadPath(path, true, env)
	if err == nil && parent != nil {
		parent.RemoveChild(node)
	}
	return src
}
