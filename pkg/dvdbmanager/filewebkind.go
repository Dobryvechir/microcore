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
	webUrl := tbl.Web
	webPath := db.WebRoot + tbl.Web
	keyFirst := evaluateKeyFirst(tbl)
	webField := evaluateWebField(tbl)
	webFileName := evaluateWebFileName(tbl)
	webAllowedFormats := evaluateWebFormats(tbl)
	os.MkdirAll(webPath, 0755)
	ref := &fileWebTable{path, webUrl, webPath, keyFirst, webField, webFileName, webAllowedFormats}
	return ref
}

func (tbl *fileWebTable) ReadAll() interface{} {
	return readWholeFileAsJsonArray(tbl.path)
}

func (tbl *fileWebTable) ReadOne(key interface{}) interface{} {
	return findSingleEntryInJsonArray(tbl.path, key, tbl.keyFirst)
}

func (tbl *fileWebTable) ReadFieldsForIds(ids []*dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
	return readFieldsForIdsInJson(tbl.path, ids, fields, tbl.keyFirst)
}

func (tbl *fileWebTable) ReadFieldsForId(id *dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error) {
	return readFieldsForIdInJson(tbl.path, id, fields, tbl.keyFirst)
}

func (tbl *fileWebTable) ReadFieldsForAll(fields []string) (*dvevaluation.DvVariable, error) {
	return readFieldsForAllInJson(tbl.path, fields)
}

func (tbl *fileWebTable) DeleteKeys(keys []string) interface{} {
	deleteWebFiles(tbl.webPath, keys)
	return deleteKeysInJson(tbl.path, keys, tbl.keyFirst)
}

func (tbl *fileWebTable) CreateRecord(record *dvevaluation.DvVariable, newId string) (*dvevaluation.DvVariable, error) {
	name, err := updateWebFiles(tbl.webPath, record, newId, tbl.webField, tbl.webFileName, tbl.webUrl, tbl.webAllowedFormats)
	if err != nil {
		return nil, err
	}
	js, err := createRecordInJson(tbl.path, record, tbl.keyFirst, newId)
	if err == nil {
		return js, nil
	}
	cleanWebFiles(name)
	return err
}

func (tbl *fileWebTable) UpdateRecord(record *dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	id, ok := readFieldInJsonAsString(record, tbl.keyFirst)
	if !ok || len(id) == 0 {
		return nil, errors.New(tbl.keyFirst + " field is missing")
	}
	_, err := updateWebFiles(tbl.webPath, record, id, tbl.webField, tbl.webFileName, tbl.webUrl, tbl.webAllowedFormats)
	if err != nil {
		return nil, err
	}
	return updateRecordInJson(tbl.path, record, tbl.keyFirst)
}
