// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import "github.com/Dobryvechir/microcore/pkg/dvevaluation"

func RecordBind(table string, items *dvevaluation.DvVariable, kind string, fields string) (*dvevaluation.DvVariable, error) {
	return nil, nil
}

func RecordCreate(table string) interface{} {
	return nil
}

func RecordDelete(table string) interface{} {
	return nil
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

func RecordScan(table string) interface{} {
	return nil
}

func RecordUpdate(table string) interface{} {
	return nil
}
