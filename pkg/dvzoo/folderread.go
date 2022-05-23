/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvzoo

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strconv"
	"strings"
)

type ZooFolderActionProvider struct{}

func (provider *ZooFolderActionProvider) PathSupported() bool {
	return false
}

func (provider *ZooFolderActionProvider) Read(ctx *dvcontext.RequestContext, prefix string, key string) (interface{}, bool, error) {
	conn, err := GetZooConnect()
	if err != nil {
		return nil, false, err
	}
	defer LeaveZooConnect(conn)
	includeSys := false
	includeErr := false
	fullPath := false
	if strings.HasPrefix(prefix, "S") {
		includeSys = true
		prefix = prefix[1:]
	}
	if strings.HasPrefix(prefix, "E") {
		includeErr = true
		prefix = prefix[1:]
	}
	if strings.HasPrefix(prefix, "P") {
		fullPath = true
		prefix = prefix[1:]
	}
	v, err := ReadWholeFolder(conn, key, includeSys, includeErr, fullPath)
	if err != nil {
		return nil, false, err
	}
	return v, true, nil
}

func (provider *ZooFolderActionProvider) Delete(ctx *dvcontext.RequestContext, prefix string, key string) error {
	conn, err := GetZooConnect()
	if err != nil {
		return err
	}
	defer LeaveZooConnect(conn)
	var version int32 = 0
	includeSys := false
	if prefix != "" {
		if strings.HasPrefix(prefix, "S") {
			includeSys = true
			prefix = prefix[1:]
		}
		i, err := strconv.Atoi(prefix)
		if err == nil {
			version = int32(i)
		}
	}
	err = DeleteWholeFolder(conn, key, version, includeSys)
	return err
}

func (provider *ZooFolderActionProvider) Save(ctx *dvcontext.RequestContext, prefix string, key string, value interface{}) error {
	conn, err := GetZooConnect()
	if err != nil {
		return err
	}
	defer LeaveZooConnect(conn)
	var version int32 = 0
	includeSys := false
	if prefix != "" {
		if strings.HasPrefix(prefix, "S") {
			includeSys = true
			prefix = prefix[1:]
		}
		i, err := strconv.Atoi(prefix)
		if err == nil {
			version = int32(i)
		}
	}
	data := dvevaluation.AnyToDvVariable(value)
	err = SaveWholeFolder(conn, key, data, version, includeSys)
	return err
}
