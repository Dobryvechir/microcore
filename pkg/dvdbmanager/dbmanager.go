// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"os"
	"strconv"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
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
