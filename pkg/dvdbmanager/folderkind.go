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

func (tbl *folderTable) ReadFieldsForIds(ids []*dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
	return readFieldsForIdsInFolder(tbl.path, ids, fields, tbl.keyFirst)
}

func (tbl *folderTable) ReadFieldsForId(id *dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
	return readFieldsForIdInFolder(tbl.path, id, fields, tbl.keyFirst)
}

func (tbl *folderTable) ReadFieldsForAll(fields []string) (*dvevaluation.DvVariable, error) {
	return readFieldsForAllInFolder(tbl.path, fields)
}

func (tbl *folderTable) DeleteKeys(keys []string) interface{} {
	deleteWebFiles(tbl.path, keys)
	return keys
}

func (tbl *folderTable) CreateRecord(record *dvevaluation.DvVariable, newId string) (*dvevaluation.DvVariable, error) {
	return createRecordInFolder(tbl.path, record, tbl.keyFirst, newId)
}

func (tbl *folderTable) UpdateRecord(record *dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return updateRecordInFolder(tbl.path, record, tbl.keyFirst)
}
