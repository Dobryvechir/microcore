/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvsession

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
)

var sessionLocalStorageFns = []*dvevaluation.DvVariable{
	{
		Name: []byte("localStorage"),
		Kind: dvevaluation.FIELD_OBJECT,
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("setItem"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: localStorage_setItem,
				},
			},
			{
				Name: []byte("getItem"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: localStorage_getItem,
				},
			},
			{
				Name: []byte("removeItem"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: localStorage_removeItem,
				},
			},
			{
				Name: []byte("clear"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: localStorage_clear,
				},
			},
		},
	},
	{
		Name: []byte("sessionStorage"),
		Kind: dvevaluation.FIELD_OBJECT,
		Fields: []*dvevaluation.DvVariable{
			{
				Name: []byte("setItem"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: sessionStorage_setItem,
				},
			},
			{
				Name: []byte("getItem"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: sessionStorage_getItem,
				},
			},
			{
				Name: []byte("removeItem"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: sessionStorage_removeItem,
				},
			},
			{
				Name: []byte("clear"),
				Kind: dvevaluation.FIELD_FUNCTION,
				Extra: &dvevaluation.DvFunction{
					Fn: sessionStorage_clear,
				},
			},
		},
	},
}

var mapOfAddedGlobalKeys = make(map[string]int)

func registerSessionLocalStorage() {
	dvevaluation.AddListToGlobalFunctionPool(sessionLocalStorageFns)
}

func localStorage_setItem(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n > 0 {
		key := dvevaluation.AnyToString(params[0])
		var value interface{} = nil
		if n > 1 {
			value = params[1]
		}
		mapOfAddedGlobalKeys[key] = 1
		dvparser.SetGlobalPropertiesAnyValue(key, value)
	}
	return nil, nil
}

func localStorage_getItem(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return nil, nil
	}
	key := dvevaluation.AnyToString(params[0])
	v, ok := dvparser.ReadGlobalPropertiesAny(key)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func localStorage_removeItem(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	n := len(params)
	if n == 0 {
		return nil, nil
	}
	key := dvevaluation.AnyToString(params[0])
	dvparser.RemoveGlobalPropertiesValue(key)
	delete(mapOfAddedGlobalKeys, key)
	return nil, nil
}

func localStorage_clear(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	for key, _ := range mapOfAddedGlobalKeys {
		dvparser.RemoveGlobalPropertiesValue(key)
	}
	mapOfAddedGlobalKeys = make(map[string]int)
	return nil, nil
}

func GetStorageSession(context *dvgrammar.ExpressionContext) (dvcontext.RequestSession, error) {
	v, ok := context.Scope.Get(dvcontext.ServerSessionStoringKey)
	if !ok {
		return nil, errors.New("Session is not available")
	}
	session, ok := v.(dvcontext.RequestSession)
	if !ok {
		return nil, errors.New("Session is corrupted")
	}
	if session == nil {
		return nil, errors.New("Session has not been initialized")
	}
	return session, nil
}

func sessionStorage_setItem(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	session, err := GetStorageSession(context)
	if err != nil {
		return nil, err
	}
	n := len(params)
	if n > 0 {
		key := dvevaluation.AnyToString(params[0])
		var value interface{} = nil
		if n > 1 {
			value = params[1]
		}
		session.SetItem(key, value)
	}
	return nil, nil
}

func sessionStorage_getItem(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	session, err := GetStorageSession(context)
	if err != nil {
		return nil, err
	}
	n:=len(params)
	if n == 0 {
		return nil, nil
	}
	key := dvevaluation.AnyToString(params[0])
	v := session.GetItem(key)
	return v, nil
}

func sessionStorage_removeItem(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	session, err := GetStorageSession(context)
	if err != nil {
		return nil, err
	}
	n := len(params)
	if n == 0 {
		return nil, nil
	}
	key := dvevaluation.AnyToString(params[0])
	session.RemoveItem(key)
	return nil, nil
}

func sessionStorage_clear(context *dvgrammar.ExpressionContext, thisVariable interface{}, params []interface{}) (interface{}, error) {
	session, err := GetStorageSession(context)
	if err != nil {
		return nil, err
	}
	session.Clear()
	return nil, nil
}
