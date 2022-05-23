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
	"time"
)

var zooServerList []string

type ZooActionProvider struct{}

func GetZooConnect() (*zk.Conn, error) {
	if zooServerList == nil {
		names := dvparser.GetByGlobalPropertiesOrDefault("ZOOKEEPER_ADDRESS", "127.0.0.1:2181")
		zooServerList = dvtextutils.ConvertToList(names)
	}
	zkConn, _, err := zk.Connect(zooServerList, time.Second*600)
	if err != nil {
		return nil, err
	}
	return zkConn, nil
}

// TODO provide the mechanism for caching connections
func LeaveZooConnect(zkConn *zk.Conn) {
	if zkConn != nil {
		zkConn.Close()
	}
}

func (provider *ZooActionProvider) PathSupported() bool {
	return false
}

func (provider *ZooActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	var res []byte
	conn, err := GetZooConnect()
	if err != nil {
		return nil, false, err
	}
	defer LeaveZooConnect(conn)
	res, _, err = conn.Get(key)
	if err != nil {
		return nil, false, err
	}
	if len(res) == 0 {
		return nil, false, nil
	}
	return string(res), true, nil
}

func (provider *ZooActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	conn, err := GetZooConnect()
	if err != nil {
		return err
	}
	defer LeaveZooConnect(conn)
	var version int32 = 0
	if prefix != "" {
		i, err := strconv.Atoi(prefix)
		if err == nil {
			version = int32(i)
		}
	}
	_, stat, err1 := conn.Get(key)
	if err1 == nil {
		if stat != nil && version == 0 {
			version = stat.Version
		}
		err = conn.Delete(key, version)
	}
	return err
}

func (provider *ZooActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	conn, err := GetZooConnect()
	if err != nil {
		return err
	}
	defer LeaveZooConnect(conn)
	var version int32 = 0
	if prefix != "" {
		i, err := strconv.Atoi(prefix)
		if err == nil {
			version = int32(i)
		}
	}
	data := dvevaluation.AnyToString(value)
	_, err = EnsureZooPathValue(conn, key, version, data)
	return err
}

func zooInit() bool {
	dvaction.RegisterStorageActionProvider("zoo", &ZooActionProvider{})
	dvaction.RegisterStorageActionProvider("zoonode", &ZooFolderActionProvider{})
	return true
}

var inited = zooInit()
