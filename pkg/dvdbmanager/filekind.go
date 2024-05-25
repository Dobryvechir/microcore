// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"os"

	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

func fileKindInit(tbl *dvcontext.DatabaseTable, db *dvcontext.DatabaseConfig) genTable {
	path := db.Root + "/" + tbl.Name + ".json"
	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		os.WriteFile(path, []byte(emptyArray), 0644)
	}
	keyFirst := evaluateKeyFirst(tbl)
	allowCustomId := tbl.AllowCustomId
	version := tbl.Version
	ref := &fileTable{path: path, keyFirst: keyFirst, allowCustomId: allowCustomId, version: version}
	return ref
}

func (tbl *fileTable) ReadAll() (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return readWholeFileAsJsonArray(tbl.path)
}

func (tbl *fileTable) ReadOne(key interface{}) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return findSingleEntryInJsonArray(tbl.path, key, tbl.keyFirst)
}

func (tbl *fileTable) ReadFieldsForIds(ids []*dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return readFieldsForIdsInJson(tbl.path, ids, fields, tbl.keyFirst)
}

func (tbl *fileTable) ReadFieldsForId(id *dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return readFieldsForIdInJson(tbl.path, id, fields, tbl.keyFirst)
}

func (tbl *fileTable) ReadFieldsForAll(fields []string) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return readFieldsForAllInJson(tbl.path, fields)
}

func (tbl *fileTable) DeleteKeys(keys []string) interface{} {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return deleteKeysInJson(tbl.path, keys, tbl.keyFirst)
}

func (tbl *fileTable) CreateRecord(record *dvevaluation.DvVariable, newId string) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	var err error
	newId, err = resolveCustomId(record, newId, tbl.allowCustomId, tbl.keyFirst, tbl.version)
	if err != nil {
		return nil, err
	}
	return createRecordInJson(tbl.path, record, tbl.keyFirst, newId)
}

func (tbl *fileTable) UpdateRecord(record *dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return updateRecordInJson(tbl.path, record, tbl.keyFirst, tbl.version)
}

func (tbl *fileTable) CreateOrUpdateByConditionsAndUpdateFields(record *dvevaluation.DvVariable, conditions []string, fields []string) (*dvevaluation.DvVariable, error) {
	tbl.mu.Lock()
	defer tbl.mu.Unlock()
	return CreateOrUpdateByConditionsAndUpdateFieldsForJson(tbl.path, record, conditions, fields, tbl.keyFirst, tbl.version)
}
