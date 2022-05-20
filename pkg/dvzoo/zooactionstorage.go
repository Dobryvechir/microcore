/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvzoo

import (
	"github.com/Dobryvechir/microcore/pkg/dvaction"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"github.com/go-zookeeper/zk"
	"strconv"
	"sync"
	"time"
)

var zkConn *zk.Conn
var zkConnMutex sync.Mutex

type ZooActionProvider struct{}

func GetZConnect() (*zk.Conn, error) {
	if zkConn != nil {
		return zkConn, nil
	}
	names := dvparser.GetByGlobalPropertiesOrDefault("ZOOKEEPER_ADDRESS", "127.0.0.1:2181")
	nameList := dvtextutils.ConvertToList(names)
	var err error
	zkConnMutex.Lock()
	if zkConn != nil {
		zkConn, _, err = zk.Connect(nameList, time.Second*10)
	}
	zkConnMutex.Unlock()
	return zkConn, err
}

func (provider *ZooActionProvider) PathSupported() bool {
	return false
}

func (provider *ZooActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool) {
	conn, err := GetZConnect()
	if err != nil {
		return nil, false
	}
	var res []byte
	res, _, err = conn.Get(key)
	if err != nil || len(res) == 0 {
		return nil, false
	}
	return string(res), true
}

func (provider *ZooActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	conn, err := GetZConnect()
	if err != nil {
		return err
	}
	var version int32 = 0
	if prefix != "" {
		i, err := strconv.Atoi(prefix)
		if err == nil {
			version = int32(i)
		}
	}
	err = conn.Delete(key, version)
	return err
}

func (provider *ZooActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	conn, err := GetZConnect()
	if err != nil {
		return err
	}
	var version int32 = 0
	if prefix != "" {
		i, err := strconv.Atoi(prefix)
		if err == nil {
			version = int32(i)
		}
	}
	data := []byte(dvevaluation.AnyToString(value))
	_, err = conn.Set(key, data, version)
	return err
}

func zooInit() bool {
	dvaction.RegisterStorageActionProvider("zoo", &ZooActionProvider{})
	dvaction.RegisterStorageActionProvider("zoonode", &ZooFolderActionProvider{})
	return true
}

var inited = zooInit()
