// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"os"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
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

func readAllFolderItemsAsList(path string) interface{} {
	items, err := collectAllFolderItemsAsList(path)
	if err != nil {
		return err
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
	return res
}

func getEntryName(path string, key interface{}) string {
	return path + "/" + dvevaluation.AnyToString(key) + ".json"
}
func findSingleEntryInFolder(path string, key interface{}, _ string) interface{} {
	keyPath := getEntryName(path, key)
	d, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil
	}
	return d
}

func readFieldsForIdsInFolder(path string, ids []*dvevaluation.DvVariable, fieldNames []string, key string) (*dvevaluation.DvVariable, error) {
	n := len(ids)
	fields := make([]*dvevaluation.DvVariable, 0, n)
	fieldMap := convertFieldsToMap(fields)
	for i := 0; i < n; i++ {
		p := readFieldMapForIdInFolder(path, ids[i], fieldMap, key)
		if p != nil {
			fields = append(fields, p)
		}
	}
	res := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: fields}
	return res, nil
}

func readFieldsForIdInFolder(path string, id *dvevaluation.DvVariable, fields []string, key string) (*dvevaluation.DvVariable, error) {
	fieldMap := convertFieldsToMap(fields)
	return readFieldMapForIdInFolder(path, id, fieldMap, key)
}

func readFieldMapForIdInFolder(path string, id interface{}, fieldMap map[string]int, key string) (*dvevaluation.DvVariable, error) {
	keyPath := path + "/" + dvevaluation.AnyToString(id) + ".json"
	d, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil
	}
	v := reduceJsonToFields(d, fieldMap)
	return v
}

func readFieldsForAllInFolder(path string, fields []string) (*dvevaluation.DvVariable, error) {
	ids, err := collectAllFolderItemsAsList(path)
	if err != nil {
		return nil, err
	}
	n := len(ids)
	fields := make([]*dvevaluation.DvVariable, 0, n)
	fieldMap := convertFieldsToMap(fields)
	for i := 0; i < n; i++ {
		p := readFieldMapForIdInFolder(path, ids[i], fieldMap, key)
		if p != nil {
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
	keyPath := getEntryName(path, id)
	err := writeWholeFileAsJson(keyPath, record)
	return record, err
}

func updateRecordInFolder(path string, record *dvevaluation.DvVariable, keyFirst string) (*dvevaluation.DvVariable, error) {
	id, ok := readFieldInJsonAsString(record, keyFirst)
	if !ok || !checkIntId(id) == 0 {
		return nil, errors.New("object has no id")
	}
	keyPath := getEntryName(path, id)
	_, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil, err
	}
	err := writeWholeFileAsJson(keyPath, record)
	return record, err
}
