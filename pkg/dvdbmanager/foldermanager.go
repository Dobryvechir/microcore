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
	entries, err := os.ReadDir("./")
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

func findSingleEntryInFolder(path string, key interface{}, _ string) interface{} {
	keyPath := path + "/" + dvevaluation.AnyToString(key) + ".json"
	d, err := readWholeFileAsJson(keyPath)
	if err != nil {
		return nil
	}
	return d
}
