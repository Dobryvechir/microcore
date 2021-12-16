/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func (parseInfo *DvCrudParsingInfo) GetDvFieldInfoHierarchy() []*dvevaluation.DvVariable {
	n := len(parseInfo.Items)
	res := make([]*dvevaluation.DvVariable, n)
	for i := 0; i < n; i++ {
		res[i] = &parseInfo.Items[i].DvVariable
	}
	return res
}

func ReadPathOfAny(item interface{}, childName string, rejectChildOfUndefined bool, env *dvevaluation.DvObject) (*dvevaluation.DvVariable, *dvevaluation.DvVariable, error) {
	switch item.(type) {
	case *dvevaluation.DvVariable:
		return item.(*dvevaluation.DvVariable).ReadPath(childName, rejectChildOfUndefined, env)
	}
	return nil, nil, nil
}

func ReadJsonChild(data []byte, childName string, rejectChildOfUndefined bool, env *dvevaluation.DvObject) (*dvevaluation.DvVariable, error) {
	item, err := JsonFullParser(data)
	if err != nil {
		return nil, err
	}
	res, _, err := item.ReadPath(childName, rejectChildOfUndefined, env)
	return res, err
}

func CountChildren(val interface{}) int {
	if val == nil {
		return 0
	}
	switch val.(type) {
	case *dvevaluation.DvVariable:
		return len(val.(*dvevaluation.DvVariable).Fields)
	}
	return 0
}

func FindDifferenceForAnyType(itemAny interface{}, otherAny interface{},
	fillAdded bool, fillRemoved bool, fillUpdated bool, fillUnchanged bool,
	fillUpdatedCounterpart bool, unchangedAsUpdated bool) (added *dvevaluation.DvVariable, removed *dvevaluation.DvVariable,
	updated *dvevaluation.DvVariable, unchanged *dvevaluation.DvVariable, counterparts *dvevaluation.DvVariable) {
	var item, other *dvevaluation.DvVariable
	switch itemAny.(type) {
	case *dvevaluation.DvVariable:
		item = itemAny.(*dvevaluation.DvVariable)
	}
	switch otherAny.(type) {
	case *dvevaluation.DvVariable:
		other = otherAny.(*dvevaluation.DvVariable)
	}
	if item.QuickSearch == nil {
		item.CreateQuickInfoForObjectType()
	}
	if other.QuickSearch == nil {
		other.CreateQuickInfoForObjectType()
	}
	return item.FindDifferenceByQuickMap(other, fillAdded, fillRemoved,
		fillUpdated, fillUnchanged, fillUpdatedCounterpart, unchangedAsUpdated)
}

func CreateQuickInfoByKeysForAny(data interface{}, ids []string) {
	n := len(ids)
	if n > 0 {
		switch data.(type) {
		case *dvevaluation.DvVariable:
			data.(*dvevaluation.DvVariable).CreateQuickInfoByKeys(ids)
		}
	} else {
		switch data.(type) {
		case *dvevaluation.DvVariable:
			data.(*dvevaluation.DvVariable).CreateQuickInfoForObjectType()
		}
	}
}
