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
        allowCustomId:=tbl.AllowCustomId
        version:=tbl.Version
	ref := &folderTable{path, keyFirst, allowCustomId, version}
	return ref
}

func (tbl *folderTable) ReadAll() interface{} {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	return readAllFolderItemsAsList(tbl.path)
}

func (tbl *folderTable) ReadOne(key interface{}) interface{} {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	return findSingleEntryInFolder(tbl.path, key, tbl.keyFirst)
}

func (tbl *folderTable) ReadFieldsForIds(ids []*dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	return readFieldsForIdsInFolder(tbl.path, ids, fields, tbl.keyFirst)
}

func (tbl *folderTable) ReadFieldsForId(id *dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	return readFieldsForIdInFolder(tbl.path, id, fields, tbl.keyFirst)
}

func (tbl *folderTable) ReadFieldsForAll(fields []string) (*dvevaluation.DvVariable, error) {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	return readFieldsForAllInFolder(tbl.path, fields)
}

func (tbl *folderTable) DeleteKeys(keys []string) interface{} {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	deleteWebFiles(tbl.path, keys)
	return keys
}

func (tbl *folderTable) CreateRecord(record *dvevaluation.DvVariable, newId string) (*dvevaluation.DvVariable, error) {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
        var err error  
        newId, err = resolveCustomId(record, newId, tbl.allowCustomId, tbl.keyFirst, tbl.version)
	if err != nil {
		return nil, err
	}
	return createRecordInFolder(tbl.path, record, tbl.keyFirst, newId)
}

func (tbl *folderTable) UpdateRecord(record *dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
        tbl.mu.Lock()
        defer tbl.mu.Unlock()
	return updateRecordInFolder(tbl.path, record, tbl.keyFirst, tbl.version)
}
