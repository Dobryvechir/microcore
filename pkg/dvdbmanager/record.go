// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func RecordBind(table string, items *dvevaluation.DvVariable, kind string, fields string) (res *dvevaluation.DvVariable, err error) {
	fieldList := dvtextutils.ConvertToNonEmptyList(fields)
	r, ok := tableMap[table]
	if !ok {
		err = errors.New("Table " + table + " does not exist")
		return
	}
	if items == nil {
		return
	}
	if kind == "array" {
		res, err = r.ReadFieldsForIds(items.Fields, fieldList)
	} else {
		res, err = r.ReadFieldsForId(items, fieldList)
	}
	return
}

func RecordCreate(table string, body string, newId string) interface{} {
	r, ok := tableMap[table]
	if !ok {
		return "Table " + table + " does not exist"
	}
	js, err := dvjson.JsonFullParser([]byte(body))
	if err != nil {
		return err
	}
	if js == nil || js.Kind != dvevaluation.FIELD_OBJECT {
		return "Empty object"
	}
	js, err = r.CreateRecord(js, newId)
	if err != nil {
		return err
	}
	return js
}

func RecordDelete(table string, keys string) interface{} {
	ids := dvtextutils.ConvertToNonEmptyList(keys)
	if len(ids) == 0 {
		return ""
	}
	r, ok := tableMap[table]
	if !ok {
		return "Table " + table + " does not exist"
	}
	d := r.DeleteKeys(ids)
	return d
}

func RecordReadAll(table string) interface{} {
	r, ok := tableMap[table]
	if !ok {
		return "Table " + table + " does not exist"
	}
	d := r.ReadAll()
	return d
}

func RecordReadOne(table string, key interface{}) interface{} {
	r, ok := tableMap[table]
	if !ok {
		return "Table " + table + " does not exist"
	}
	d := r.ReadOne(key)
	return d
}

func RecordScan(table string, fields string) (res *dvevaluation.DvVariable, err error) {
	fieldList := dvtextutils.ConvertToNonEmptyList(fields)
	r, ok := tableMap[table]
	if !ok {
		err = errors.New("Table " + table + " does not exist")
		return
	}
	res, err = r.ReadFieldsForAll(fieldList)
	return
}

func RecordUpdate(table string, body string) interface{} {
	r, ok := tableMap[table]
	if !ok {
		return "Table " + table + " does not exist"
	}
	js, err := dvjson.JsonFullParser([]byte(body))
	if err != nil {
		return err
	}
	if js == nil || js.Kind != dvevaluation.FIELD_OBJECT {
		return "Empty object"
	}
	js, err = r.UpdateRecord(js)
	if err != nil {
		return err
	}
	return js
}
