// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"os"
	"strconv"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

const (
	KindFile    = "file"
	KindFolder  = "folder"
	KindWebFile = "fileweb"
)

var dbAmount = 0

func DbManagerInit(conf []*dvcontext.DatabaseConfig) {
	dbAmount = len(conf)
	if dbAmount == 0 {
		return
	}
	tableMap = make(map[string]genTable)
	for i := 0; i < dbAmount; i++ {
		db := conf[i]
		if db == nil || len(db.Root) == 0 {
			panic("Incorrect dbs initialization")
		}
		os.MkdirAll(db.Root, 0755)
		if len(db.WebRoot) != 0 {
			os.MkdirAll(db.WebRoot, 0755)
		}
		n := len(db.Tables)
		if len(db.Name) == 0 {
			db.Name = strconv.Itoa(i)
		}
		var tableRef genTable
		for j := 0; j < n; j++ {
			tbl := db.Tables[i]
			switch tbl.Kind {
			case KindFile:
				tableRef = fileKindInit(tbl, db)
			case KindFolder:
				tableRef = folderKindInit(tbl, db)
			case KindWebFile:
				tableRef = fileWebKindInit(tbl, db)
			default:
				panic("Incorrect table definitions")
			}
			if i == 0 {
				tableMap[tbl.Name] = tableRef
			}
			tableMap[db.Name+"."+tbl.Name] = tableRef
			dvlog.PrintlnError("Registered table " + tbl.Name)
		}
	}
}

func evaluateKeyFirst(tbl *dvcontext.DatabaseTable) string {
	res := tbl.KeyFirst
	if len(res) == 0 {
		res = defaultKeyFirst
	}
	return res
}

func evaluateWebField(tbl *dvcontext.DatabaseTable) string {
	res := tbl.WebField
	if len(res) == 0 {
		res = defaultWebField
	}
	return res
}

func evaluateWebFileName(tbl *dvcontext.DatabaseTable) string {
	res := tbl.WebFileName
	if len(res) == 0 {
		res = defaultWebFileName
	}
	return res
}

func evaluateWebFormats(tbl *dvcontext.DatabaseTable) string {
	res := tbl.WebFormats
	if len(res) == 0 {
		res = defaultWebAllowedFormats
	}
	return res
}

func resolveCustomId(record *dvevaluation.DvVariable, newId string, allowCustomId bool, keyFirst string, version string) (string, error) {
	if record == nil || record.Kind != dvevaluation.FIELD_OBJECT {
		return newId, errors.New("Body must contain json object")
	}
	if len(version) > 0 {
		setFieldInJsonAsString(record, version, "1")
	}
	if !allowCustomId {
		return newId, nil
	}
	id, ok := readFieldInJsonAsString(record, keyFirst)
	if !ok || len(id) == 0 {
		return newId, nil
	}
	return id, nil
}

func resolveVersion(oldRecord *dvevaluation.DvVariable, newRecord *dvevaluation.DvVariable, version string) {
	if len(version) == 0 || newRecord == nil {
		return
	}
	value := 0
	if oldRecord != nil && oldRecord.Kind == dvevaluation.FIELD_OBJECT {
		str, ok := readFieldInJsonAsString(oldRecord, version)
		if ok && len(str) > 0 {
			n, err := strconv.Atoi(str)
			if err == nil && n > 0 {
				value = n
			}
		}
	}
	setFieldInJsonAsString(newRecord, version, strconv.Itoa(value+1))
}
