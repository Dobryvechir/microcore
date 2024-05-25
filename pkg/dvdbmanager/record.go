// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func generateTableDoesNotExist(table string) string {
	s := ""
	for k, _ := range tableMap {
		if len(s) > 0 {
			s += ","
		}
		s += k
	}
	return "{\"error\":\"Table " + table + " does not exist [" + s + "]\"}"
}

func RecordBind(table string, items *dvevaluation.DvVariable, kind string, fields string) (res *dvevaluation.DvVariable, err error) {
	fieldList := dvtextutils.ConvertToNonEmptyList(fields)
	r, ok := tableMap[table]
	if !ok {
		err = errors.New(generateTableDoesNotExist(table))
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
		return generateTableDoesNotExist(table)
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
		return generateTableDoesNotExist(table)
	}
	d := r.DeleteKeys(ids)
	return d
}

func CreateOrUpdateByConditionsAndUpdateFields(table string, row *dvevaluation.DvVariable, conditions []string, fields []string) (*dvevaluation.DvVariable, error) {
	r, ok := tableMap[table]
	if !ok {
		return nil, errors.New(generateTableDoesNotExist(table))
	}
	if row == nil || len(conditions) == 0 || len(conditions) != len(fields) {
		return nil, errors.New("Error call to CreateOrUpdateByConditionsAndUpdateFields")
	}
	return r.CreateOrUpdateByConditionsAndUpdateFields(row, conditions, fields)
}

func RecordReadAll(table string) (*dvevaluation.DvVariable, error) {
	r, ok := tableMap[table]
	if !ok {
		return nil, errors.New(generateTableDoesNotExist(table))
	}
	d, err := r.ReadAll()
	return d, err
}

func RecordReadOne(table string, key interface{}) (*dvevaluation.DvVariable, error) {
	r, ok := tableMap[table]
	if !ok {
		return nil, errors.New(generateTableDoesNotExist(table))
	}
	d, err := r.ReadOne(key)
	return d, err
}

func RecordScan(table string, fields string) (res *dvevaluation.DvVariable, err error) {
	fieldList := dvtextutils.ConvertToNonEmptyList(fields)
	r, ok := tableMap[table]
	if !ok {
		err = errors.New(generateTableDoesNotExist(table))
		return
	}
	res, err = r.ReadFieldsForAll(fieldList)
	return
}

func RecordUpdate(table string, body string) interface{} {
	r, ok := tableMap[table]
	if !ok {
		return generateTableDoesNotExist(table)
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
