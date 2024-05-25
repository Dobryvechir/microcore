// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"os"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func collectAllFolderItemsAsList(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	n := len(entries)
	res := make([]string, 0, n)
	for _, e := range entries {
		name := e.Name()
		if e.Type().IsRegular() && strings.HasSuffix(name, ".json") {
			res = append(res, name[:len(name)-5])
		}
	}
	return res, nil
}

func readAllFolderItemsAsList(path string) (*dvevaluation.DvVariable, error) {
	items, err := collectAllFolderItemsAsList(path)
	if err != nil {
		return nil, err
	}
	n := len(items)
	res := &dvevaluation.DvVariable{
		Kind:   dvevaluation.FIELD_ARRAY,
		Fields: make([]*dvevaluation.DvVariable, 0, n),
	}
	for i := 0; i < n; i++ {
		data, err := os.ReadFile(path + "/" + items[i] + ".json")
		if err == nil && data != nil {
			js, err := dvjson.JsonFullParser(data)
			if err == nil && js != nil {
				res.Fields = append(res.Fields, js)
			}
		}
	}
	return res, nil
}

func getEntryName(path string, key interface{}) string {
	return path + "/" + dvevaluation.AnyToString(key) + ".json"
}

func findSingleEntryInFolder(path string, key interface{}, _ string) (*dvevaluation.DvVariable, error) {
	keyPath := getEntryName(path, key)
	d, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil, nil
	}
	return d, nil
}

func readFieldsForIdsInFolder(path string, ids []*dvevaluation.DvVariable, fieldNames []string, key string) (*dvevaluation.DvVariable, error) {
	n := len(ids)
	fields := make([]*dvevaluation.DvVariable, 0, n)
	fieldMap := convertDvVariableFieldsToMap(fields)
	for i := 0; i < n; i++ {
		p, err := readFieldMapForIdInFolder(path, ids[i], fieldMap)
		if err == nil && p != nil {
			fields = append(fields, p)
		}
	}
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: fields}
	return res, nil
}

func readFieldsForIdInFolder(path string, id *dvevaluation.DvVariable, fields []string, key string) (*dvevaluation.DvVariable, error) {
	fieldMap := convertFieldsToMap(fields)
	return readFieldMapForIdInFolder(path, id, fieldMap)
}

func readFieldMapForIdInFolder(path string, id interface{}, fieldMap map[string]int) (*dvevaluation.DvVariable, error) {
	keyPath := path + "/" + dvevaluation.AnyToString(id) + ".json"
	d, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil, err
	}
	v := reduceJsonToFields(d, fieldMap)
	return v, nil
}

func readFieldsForAllInFolder(path string, fieldNames []string) (*dvevaluation.DvVariable, error) {
	ids, err := collectAllFolderItemsAsList(path)
	if err != nil {
		return nil, err
	}
	n := len(ids)
	fields := make([]*dvevaluation.DvVariable, 0, n)
	fieldMap := convertFieldsToMap(fieldNames)
	for i := 0; i < n; i++ {
		p, err := readFieldMapForIdInFolder(path, ids[i], fieldMap)
		if err == nil && p != nil {
			fields = append(fields, p)
		}
	}
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: fields}
	return res, nil
}

func createRecordInFolder(path string, record *dvevaluation.DvVariable, keyFirst string, newId string) (*dvevaluation.DvVariable, error) {
	if !setFieldInJsonAsString(record, keyFirst, newId) {
		return nil, errors.New("Request body is not a JSON object")
	}
	keyPath := getEntryName(path, newId)
	err := writeWholeFileAsJson(keyPath, record)
	return record, err
}

func updateRecordInFolder(path string, record *dvevaluation.DvVariable, keyFirst string, version string) (*dvevaluation.DvVariable, error) {
	id, ok := readFieldInJsonAsString(record, keyFirst)
	if !ok || !checkIntId(id) {
		return nil, errors.New("object has no id")
	}
	keyPath := getEntryName(path, id)
	oldRecord, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil, err
	}
	resolveVersion(oldRecord, record, version)
	err = writeWholeFileAsJson(keyPath, record)
	return record, err
}

func CreateOrUpdateByConditionsAndUpdateFieldsForFolder(path string, record *dvevaluation.DvVariable, conditions []string, fields []string, keyFirst string, version string) (*dvevaluation.DvVariable, error) {
	id, ok := readFieldInJsonAsString(record, keyFirst)
	if !ok || !checkIntId(id) {
		return nil, errors.New("object has no id")
	}
	keyPath := getEntryName(path, id)
	previousRecord, err := readWholeFileAsJson(keyPath)
	if err != nil || previousRecord == nil {
		if dvtextutils.IsStringContainedInArray("NEW", conditions) {
			if len(version) > 0 {
				setFieldInJsonAsString(record, version, "1")
			}
			err = writeWholeFileAsJson(keyPath, record)
			return record, err
		}
		return nil, nil
	}
	n, err := findFirstMetCondition(previousRecord, record, conditions)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return previousRecord, nil
	}
	changed := updateRecordByFields(previousRecord, record, fields[n])
	if !changed {
		return previousRecord, nil
	}
	resolveVersion(previousRecord, record, version)
	err = writeWholeFileAsJson(keyPath, record)
	return record, err
}
