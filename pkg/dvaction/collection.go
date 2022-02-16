/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import "github.com/Dobryvechir/microcore/pkg/dvevaluation"

type Collection struct {
	Source    string             `json:"source"`
	Unique    []string           `json:"unique"`
	Assign    []*AssignmentBlock `json:"assign"`
	MergeMode int                `json:"merge_mode"`
	Append    bool               `json:"append"`
}

func (collection *Collection) AddItem(env *dvevaluation.DvObject) {
	n := len(collection.Assign)
	item := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, n)}
	for i := 0; i < n; i++ {
		a := collection.Assign[i]
		item.AssignToSubField(a.Field, a.Value, env)
	}
	var s *dvevaluation.DvVariable = nil
	src, ok := env.Get(collection.Source)
	if ok {
		s = dvevaluation.AnyToDvVariable(src)
	}
	if s == nil || s.Kind != dvevaluation.FIELD_ARRAY {
		s = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, 64)}
		env.Set(collection.Source, s)
	}
	s.MergeItemIntoArraysByIds(item, collection.Unique, collection.MergeMode, true)
}

func (collection *Collection) Initialize(env *dvevaluation.DvObject) {
	item := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, 64)}
	if collection.Append {
		v, ok := env.Get(collection.Source)
		if ok && v != nil {
			d := dvevaluation.AnyToDvVariable(v)
			if d != nil && d.Kind == dvevaluation.FIELD_ARRAY {
				item = d
			}
		}
	}
	env.Set(collection.Source, item)
	item.CreateQuickInfoByKeys(collection.Unique)
}

func (collection *Collection) AddItemSecondary(env *dvevaluation.DvObject) {
	n := len(collection.Assign)
	item := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, n)}
	for i := 0; i < n; i++ {
		a := collection.Assign[i]
		item.AssignToSubField(a.Field, a.Value, env)
	}
	src, _ := env.Get(collection.Source)
	s := dvevaluation.AnyToDvVariable(src)
	s.MergeItemIntoArraysByIds(item, collection.Unique, collection.MergeMode, false)
}
