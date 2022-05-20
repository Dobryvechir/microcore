/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvzoo

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/go-zookeeper/zk"
	"strings"
)

const (
	ZooSysFolder     = "/zookeeper"
	ZooSysFolderPref = "/zookeeper/"
)

var (
	path_st     = "path"
	children_st = "children"
	value_st    = "value"
	path_bt     = []byte("path")
	children_bt = []byte("children")
	value_bt    = []byte("value")
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
		} else if rC != nil {
			fld = append(fld, rC)
		}
	}
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
		err = DeleteWholeFolder(conn, addPath+children[i], version, includeSys)
		if err != nil {
			return err
		}
	}
	if path != "/" {
		err = conn.Delete(path, version)
	}
	return
}

func SaveWholeFolder(conn *zk.Conn, path string, r *dvevaluation.DvVariable, version int32, includeSys bool) (err error) {
	if path == "" || path[0] != '/' || (!includeSys && (path == ZooSysFolder || strings.HasPrefix(path, ZooSysFolderPref))) {
		return nil
	}
	s, err := EnsureZooPath(conn, path, version, "")
	if err != nil {
		return err
	}
	if s != "" {
		path = s
	}
	if r != nil && r.Kind == dvevaluation.FIELD_OBJECT {
		value := r.ReadChildStringValue(value_st)
		_, err = conn.Set(path, []byte(value), version)
		if err != nil {
			return err
		}
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

func EnsureZooPath(conn *zk.Conn, path string, version int32, defValue string) (string, error) {
	_, _, err1 := conn.Get(path)
	if err1 == nil {
		return "", nil
	}
	s, err := conn.Create(path, []byte(defValue), version, nil)
	return s, err
}

func latestPath(path string) string {
	p := strings.LastIndex(path, "/")
	if p < 0 {
		return path
	}
	return path[p+1:]
}
