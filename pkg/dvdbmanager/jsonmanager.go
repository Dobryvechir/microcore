// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"bytes"
	"os"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
)

func readWholeFileAsJsonArray(path string) interface{} {
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return err
	}
	return d
}

func readWholeFileAsJson(path string) (*dvevaluation.DvVariable, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	js, err := dvjson.JsonFullParser(data)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func findSingleEntryInJsonArray(path string, key interface{}, keyFirst string) interface{} {
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return err
	}
	res := findInJsonArrayByKeyFirst(d, key, keyFirst)
	return res
}

func findInJsonArrayByKeyFirst(d *dvevaluation.DvVariable, key interface{}, keyFirst string) *dvevaluation.DvVariable {
	if d == nil || len(d.Fields) == 0 {
		return nil
	}
	keyValue := dvevaluation.AnyToByteArray(key)
	keyName := []byte(keyFirst)
	for _, item := range d.Fields {
		if item != nil && len(item.Fields) != 0 && findKeyNameValue(item.Fields, keyValue, keyName) {
			return item
		}
	}
	return nil
}

func findKeyNameValue(fields []*dvevaluation.DvVariable, name []byte, value []byte) bool {
	for _, field := range fields {
		if field != nil && bytes.Equal(field.Name, name) {
			return bytes.Equal(field.Value, value)
		}
	}
	return false
}
