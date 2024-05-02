// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"os"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
)

func fileWebKindInit(tbl *dvcontext.DatabaseTable, db *dvcontext.DatabaseConfig) genTable {
	path := db.Root + "/" + tbl.Name + ".json"
	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		os.WriteFile(path, []byte(emptyArray), 0644)
	}
	webPath := db.WebRoot + tbl.Web
	keyFirst := evaluateKeyFirst(tbl)
	os.MkdirAll(webPath, 0755)
	ref := &fileWebTable{path, webPath, keyFirst}
	return ref
}

func (tbl *fileWebTable) ReadAll() interface{} {
	return readWholeFileAsJsonArray(tbl.path)
}

func (tbl *fileWebTable) ReadOne(key interface{}) interface{} {
	return findSingleEntryInJsonArray(tbl.path, key, tbl.keyFirst)
}
