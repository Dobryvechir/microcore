/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvsession"
)

type LogActionProvider struct{}
type ErrorActionProvider struct{}
type SessionActionProvider struct{}
type MapActionProvider struct{}
type GlobalActionProvider struct{}
type RequestActionProvider struct{}

func (provider *LogActionProvider) PathSupported() bool {
	return false
}

func (provider *LogActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	return nil, false, nil
}

func (provider *LogActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	return nil
}

func (provider *LogActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	if prefix != "" {
		dvlog.PrintfError(prefix+" "+key, value)
	} else {
		dvlog.PrintfError(key, value)
	}
	return nil
}

//********************************************************************************************************

func (provider *ErrorActionProvider) PathSupported() bool {
	return false
}

func (provider *ErrorActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	return nil, false, nil
}

func (provider *ErrorActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	return nil
}

func (provider *ErrorActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	if prefix != "" {
		dvlog.PrintfError(prefix+" "+key, value)
	} else {
		dvlog.PrintfError(key, value)
	}
	v := dvevaluation.AnyToString(value)
	ActionInternalException(0, prefix, v, ctx)
	return nil
}

//********************************************************************************************************

func (provider *SessionActionProvider) PathSupported() bool {
	return true
}

func (provider *SessionActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	if ctx.Session == nil {
		return nil, false, errors.New("No session available for "+ key)
	}
	res := ctx.Session.GetItem(key)
	return res, res!=nil, nil
}

func (provider *SessionActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	if ctx.Session == nil {
		return errors.New("No session available")
	}
	ctx.Session.RemoveItem(key)
	return nil
}

func (provider *SessionActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	if ctx.Session == nil {
		return errors.New("No session provided")
	}
	ctx.Session.SetItem(key, value)
	return nil
}

//********************************************************************************************************
func (provider *MapActionProvider) PathSupported() bool {
	return true
}

func (provider *MapActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	res, ok := dvsession.GlobalMapRead(prefix, key)
	return res, ok, nil
}

func (provider *MapActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	dvsession.GlobalMapDelete(prefix, key)
	return nil
}

func (provider *MapActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	dvsession.GlobalMapWrite(prefix, key, value)
	return nil
}

//********************************************************************************************************
func (provider *GlobalActionProvider) PathSupported() bool {
	return true
}

func (provider *GlobalActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	res, ok := dvparser.ReadGlobalPropertiesAny(key)
	return res, ok, nil
}

func (provider *GlobalActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	dvparser.RemoveGlobalPropertiesValue(key)
	return nil
}

func (provider *GlobalActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	dvparser.SetGlobalPropertiesAnyValue(key, value)
	return nil
}

//********************************************************************************************************
func (provider *RequestActionProvider) PathSupported() bool {
	return true
}

func (provider *RequestActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	if ctx!=nil && ctx.PrimaryContextEnvironment!=nil {
		res, ok:=ctx.PrimaryContextEnvironment.Get(key)
		return res, ok, nil
	}
	return nil, false, nil
}

func (provider *RequestActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	if ctx!=nil && ctx.PrimaryContextEnvironment!=nil {
		ctx.PrimaryContextEnvironment.Delete(key)
	}
	return nil
}

func (provider *RequestActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	if ctx!=nil && ctx.PrimaryContextEnvironment!=nil {
		ctx.PrimaryContextEnvironment.Set(key, value)
	}
	return nil
}

//********************************************************************************************************

func initActionProvider() bool {
	StorageProviders["log"] = &LogActionProvider{}
	StorageProviders["error"] = &ErrorActionProvider{}
	StorageProviders["session"] = &SessionActionProvider{}
	StorageProviders["map"] = &MapActionProvider{}
	StorageProviders["global"] = &GlobalActionProvider{}
	StorageProviders["request"] = &RequestActionProvider{}
	return true
}

var actionProviderInited = initActionProvider()
