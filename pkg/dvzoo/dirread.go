/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvzoo

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/go-zookeeper/zk"
	"strings"
)

const (
	ZooSysFolder     = "/zookeeper"
	ZooSysFolderPref = "/zookeeper/"
)

var (
	path_st            = "path"
	children_st        = "children"
	value_st           = "value"
	path_bt            = []byte("path")
	children_bt        = []byte("children")
	value_bt           = []byte("value")
	zkCreatedNode      = int32(-30)
	zkAlreadySameValue = int32(-31)
)

func ReadWholeFolder(conn *zk.Conn, path string, includeSys bool, includeErr bool, fullPath bool) (r *dvevaluation.DvVariable, err error) {
	if (!includeSys && (path == ZooSysFolder || strings.HasPrefix(path, ZooSysFolderPref))) || path == "" || path[0] != '/' {
		return nil, nil
	}
	var v []string
	var d []byte
	v, _, err = conn.Children(path)
	if err != nil {
		return nil, err
	}
	d, _, err = conn.Get(path)
	if err != nil {
		return nil, err
	}
	n := len(v)
	fld := make([]*dvevaluation.DvVariable, 0, n)
	vpath := path
	if !fullPath {
		vpath = latestPath(vpath)
	}
	r = &dvevaluation.DvVariable{
		Kind: dvevaluation.FIELD_OBJECT,
		Fields: []*dvevaluation.DvVariable{
			{
				Kind:  dvevaluation.FIELD_STRING,
				Name:  path_bt,
				Value: []byte(vpath),
			},
			{
				Kind:  dvevaluation.FIELD_STRING,
				Name:  value_bt,
				Value: d,
			},
			{
				Kind:   dvevaluation.FIELD_ARRAY,
				Name:   children_bt,
				Fields: fld,
			},
		},
	}
	addPath := path
	if addPath[len(addPath)-1] != '/' {
		addPath = addPath + "/"
	}
	for i := 0; i < n; i++ {
		s := v[i]
		rC, errC := ReadWholeFolder(conn, addPath+s, includeSys, includeErr, fullPath)
		if errC != nil {
			if includeErr {
				rC = &dvevaluation.DvVariable{
					Kind:  dvevaluation.FIELD_STRING,
					Name:  []byte("Error"),
					Value: []byte(s + " : " + errC.Error()),
				}
				fld = append(fld, rC)
			}
			dvlog.PrintfError("Error at %s %s", addPath, s)
		} else if rC != nil {
			fld = append(fld, rC)
		} else {
			dvlog.PrintfError("not added because empty %s %s", addPath, s)
		}
	}
	r.Fields[2].Fields = fld
	return r, nil
}

func DeleteWholeFolder(conn *zk.Conn, path string, version int32, includeSys bool) (err error) {
	if (!includeSys && (path == ZooSysFolder || strings.HasPrefix(path, ZooSysFolderPref))) || path == "" {
		return nil
	}
	children, _, err := conn.Children(path)
	if err != nil {
		return err
	}
	n := len(children)
	addPath := path
	if addPath[len(addPath)-1] != '/' {
		addPath += "/"
	}
	for i := 0; i < n; i++ {
		s := addPath + children[i]
		err = DeleteWholeFolder(conn, s, version, includeSys)
		if err != nil {
			dvlog.PrintfError("Zoo cannot delete [%s] %v", s, err)
			return
		}
	}
	if path != "/" {
		_, stat, err1 := conn.Get(path)
		versionReal := version
		if err1 != nil || stat == nil {
			dvlog.PrintfError("Error reading %s %v %v", path, stat, err1)
		} else {
			versionReal = stat.Version
		}
		err = conn.Delete(path, versionReal)
		if err != nil {
			dvlog.PrintfError("Zoo cannot delete at [%s] %v level %d", path, err, version)
			return
		}
	}
	return
}

func SaveWholeFolder(conn *zk.Conn, path string, r *dvevaluation.DvVariable, version int32, includeSys bool) (err error) {
	if path == "" || path[0] != '/' || (!includeSys && (path == ZooSysFolder || strings.HasPrefix(path, ZooSysFolderPref))) {
		return nil
	}
	value := ""
	if r != nil && r.Kind == dvevaluation.FIELD_OBJECT {
		value = r.ReadChildStringValue(value_st)
	}
	s, err := EnsureZooPathValue(conn, path, version, value)
	if err != nil {
		return err
	}
	if s != "" {
		path = s
	}
	if r != nil && r.Kind == dvevaluation.FIELD_OBJECT {
		kids := r.ReadSimpleChild(children_st)
		if kids != nil && kids.Kind == dvevaluation.FIELD_ARRAY || len(kids.Fields) > 0 {
			n := len(kids.Fields)
			addPath := path
			if addPath[len(addPath)-1] != '/' {
				addPath = addPath + "/"
			}
			for i := 0; i < n; i++ {
				v := kids.Fields[i]
				if v == nil {
					continue
				}
				subPath := latestPath(v.ReadChildStringValue(path_st))
				if subPath == "" {
					continue
				}
				allPath := addPath + subPath
				err = SaveWholeFolder(conn, allPath, v, version, includeSys)
				if err != nil {
					return err
				}
			}
		}
	}
	return
}

func EnsureZooPath(conn *zk.Conn, path string, version int32, defValue string) (string, int32, error) {
	vl, stat, err1 := conn.Get(path)
	if err1 == nil && stat != nil {
		realVersion := stat.Version + 1
		if string(vl) == defValue {
			realVersion = zkAlreadySameValue
		}
		return "", realVersion, nil
	}
	s, err := conn.Create(path, []byte(defValue), version, zk.WorldACL(zk.PermAll))
	return s, zkCreatedNode, err
}

func EnsureZooPathValue(conn *zk.Conn, path string, version int32, defValue string) (string, error) {
	newPath, newVersion, err := EnsureZooPath(conn, path, version, defValue)
	if err != nil {
		return "", err
	}
	if newVersion >= 0 {
		_, err = conn.Set(path, []byte(defValue), newVersion)
	}
	return newPath, err
}

func latestPath(path string) string {
	p := strings.LastIndex(path, "/")
	if p < 0 {
		return path
	}
	return path[p+1:]
}
