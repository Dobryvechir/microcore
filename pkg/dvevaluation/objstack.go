/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

type ObjectStack struct {
	BaseLevel    *DvObject
	CurrentLevel *DvObject
}

func NewObjectStack(object *DvObject) *ObjectStack {
	objStack := &ObjectStack{
		BaseLevel:    object,
		CurrentLevel: object,
	}
	return objStack
}

func (obj *ObjectStack) Get(key string) (interface{}, bool) {
	return obj.CurrentLevel.Get(key)
}

func (obj *ObjectStack) Set(key string, value interface{}) {
	obj.CurrentLevel.Properties[key] = value
}

func (obj *ObjectStack) StackPush(option int) {
	obj.CurrentLevel = &DvObject{
		Options:    option,
		Prototype:  obj.CurrentLevel,
		Properties: make(map[string]interface{}),
	}
}

func (obj *ObjectStack) StackPop() {
	if obj.CurrentLevel != obj.BaseLevel {
		obj.CurrentLevel = obj.CurrentLevel.Prototype
	}
}

func (obj *ObjectStack) SetDeep(key string, value interface{}) {
	dvobj := obj.CurrentLevel
	for {
		_, ok := dvobj.Properties[key]
		if ok {
			dvobj.Properties[key] = value
			return
		}
		if dvobj == obj.BaseLevel {
			break
		}
		dvobj = dvobj.Prototype
	}
	obj.CurrentLevel.Properties[key] = value
}
