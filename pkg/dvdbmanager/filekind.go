// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"os"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
)

func fileKindInit(tbl *dvcontext.DatabaseTable, db *dvcontext.DatabaseConfig) genTable {
	path := db.Root + "/" + tbl.Name + ".json"
	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		os.WriteFile(path, []byte(emptyArray), 0644)
	}
	keyFirst := evaluateKeyFirst(tbl)
	ref := &fileTable{path, keyFirst}
	return ref
}

func (tbl *fileTable) ReadAll() interface{} {
	return readWholeFileAsJsonArray(tbl.path)
}

func (tbl *fileTable) ReadOne(key interface{}) interface{} {
	return findSingleEntryInJsonArray(tbl.path, key, tbl.keyFirst)
}
