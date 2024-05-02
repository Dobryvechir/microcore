// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"os"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
)

func folderKindInit(tbl *dvcontext.DatabaseTable, db *dvcontext.DatabaseConfig) genTable {
	path := db.Root + "/" + tbl.Name
	os.MkdirAll(path, 0755)
	keyFirst := evaluateKeyFirst(tbl)
	ref := &folderTable{path, keyFirst}
	return ref
}

func (tbl *folderTable) ReadAll() interface{} {
	return readAllFolderItemsAsList(tbl.path)
}

func (tbl *folderTable) ReadOne(key interface{}) interface{} {
	return findSingleEntryInFolder(tbl.path, key, tbl.keyFirst)
}
