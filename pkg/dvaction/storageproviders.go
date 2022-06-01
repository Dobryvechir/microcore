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
	"io/ioutil"
	"os"
	"strconv"
)

type LogActionProvider struct{}
type ErrorActionProvider struct{}
type SessionActionProvider struct{}
type MapActionProvider struct{}
type GlobalActionProvider struct{}
type RequestActionProvider struct{}
type FileActionProvider struct{}
type FileSizeActionProvider struct{}

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
		return nil, false, errors.New("No session available for " + key)
	}
	res := ctx.Session.GetItem(key)
	return res, res != nil, nil
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
	if ctx != nil && ctx.PrimaryContextEnvironment != nil {
		res, ok := ctx.PrimaryContextEnvironment.Get(key)
		return res, ok, nil
	}
	return nil, false, nil
}

func (provider *RequestActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	if ctx != nil && ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Delete(key)
	}
	return nil
}

func (provider *RequestActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	if ctx != nil && ctx.PrimaryContextEnvironment != nil {
		ctx.PrimaryContextEnvironment.Set(key, value)
	}
	return nil
}

//********************************************************************************************************

func (provider *FileActionProvider) PathSupported() bool {
	return false
}

func actionFileName(key string) (string, error) {
	if len(key) == 0 {
		return "", errors.New("File name is not specified")
	}
	if key[0] != '/' && key[0] != '\\' && (len(key) == 1 || key[1] == ':') {
		key = "/tmp/" + key
	}
	return key, nil
}

func (provider *FileActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	name, err := actionFileName(key)
	if err != nil {
		return nil, false, err
	}
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, false, err
	}
	return string(data), true, nil
}

func (provider *FileActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	name, err := actionFileName(key)
	if err != nil {
		return err
	}
	os.Remove(name)
	return nil
}

func (provider *FileActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	v := []byte(dvevaluation.AnyToString(value))
	name, err := actionFileName(key)
	if err != nil {
		if ctx.LogLevel >= dvlog.LogError {
			dvlog.PrintfError("Error saved %s as %s size %d err=%v", key, name, len(v), err)
		}
		return err
	}
	err = ioutil.WriteFile(name, v, 0766)
	if ctx.LogLevel >= dvlog.LogInfo {
		dvlog.PrintfError("Saved %s as %s size %d err=%v", key, name, len(v), err)
	}
	return err
}

//********************************************************************************************************

func (provider *FileSizeActionProvider) PathSupported() bool {
	return false
}

func (provider *FileSizeActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	name, err := actionFileName(key)
	if err != nil {
		return nil, false, err
	}
	stat, err := os.Stat(name)
	if err != nil {
		return err.Error(), true, nil
	}
	if stat == nil {
		return "No statistics", true, nil
	}
	r := stat.Size()
	return strconv.FormatInt(r, 16), true, nil
}

func (provider *FileSizeActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	return nil
}

func (provider *FileSizeActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
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
	StorageProviders["file"] = &FileActionProvider{}
	StorageProviders["file"] = &FileSizeActionProvider{}
	return true
}

var actionProviderInited = initActionProvider()
