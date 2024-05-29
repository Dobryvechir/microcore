// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func readWholeFileAsJsonArray(path string) (*dvevaluation.DvVariable, error) {
	d, err := readWholeFileAsJson(path)
	return d, err
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

func findSingleEntryInJsonArray(path string, key interface{}, keyFirst string) (*dvevaluation.DvVariable, error) {
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return nil, err
	}
	res := findInJsonArrayByKeyFirst(d, key, keyFirst)
	return res, nil
}

func findInJsonArrayByKeyFirst(d *dvevaluation.DvVariable, key interface{}, keyFirst string) *dvevaluation.DvVariable {
	if d == nil || len(d.Fields) == 0 {
		return nil
	}
	keyValue := dvevaluation.AnyToByteArray(key)
	keyName := []byte(keyFirst)
	for _, item := range d.Fields {
		if item != nil && len(item.Fields) != 0 && findKeyNameValue(item.Fields, keyName, keyValue) {
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

func readFieldsForIdsInJson(path string, ids []*dvevaluation.DvVariable, fields []string, keyFirst string) (*dvevaluation.DvVariable, error) {
	n := len(ids)
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return nil, err
	}
	if d == nil || len(d.Fields) == 0 || d.Kind != dvevaluation.FIELD_ARRAY {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, 0)}, nil
	}
	oldFields := d.Fields
	m := len(oldFields)
	fieldMap := convertFieldsToMap(fields)
	newFields := make([]*dvevaluation.DvVariable, 0, n)
	key := []byte(keyFirst)
	idMap := convertIdsToMap(ids)
	for i := 0; i < m; i++ {
		p := oldFields[i]
		if isKeyInMap(p, key, idMap) {
			v := reduceJsonToFields(p, fieldMap)
			newFields = append(newFields, v)
		}
	}
	d.Fields = newFields
	return d, nil

}

func readFieldsForIdInJson(path string, id *dvevaluation.DvVariable, fields []string, keyFirst string) (*dvevaluation.DvVariable, error) {
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return nil, err
	}
	res := findInJsonArrayByKeyFirst(d, id, keyFirst)
	if res == nil {
		return nil, nil
	}
	fieldMap := convertFieldsToMap(fields)
	r := reduceJsonToFields(res, fieldMap)
	return r, nil
}

func readFieldsForAllInJson(path string, fields []string) (*dvevaluation.DvVariable, error) {
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return nil, err
	}
	if d == nil || len(d.Fields) == 0 || d.Kind != dvevaluation.FIELD_ARRAY {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, 0)}, nil
	}
	n := len(d.Fields)
	fieldMap := convertFieldsToMap(fields)
	for i := 0; i < n; i++ {
		d.Fields[i] = reduceJsonToFields(d.Fields[i], fieldMap)
	}
	return d, nil
}

func reduceJsonToFields(d *dvevaluation.DvVariable, names map[string]int) *dvevaluation.DvVariable {
	if d == nil || len(d.Fields) == 0 || d.Kind != dvevaluation.FIELD_OBJECT {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: make([]*dvevaluation.DvVariable, 0, 0)}
	}
	oldFields := d.Fields
	m := len(oldFields)
	fields := make([]*dvevaluation.DvVariable, 0, m)
	for i := 0; i < m; i++ {
		p := oldFields[i]
		if p != nil {
			_, ok := names[string(p.Name)]
			if ok {
				fields = append(fields, p)
			}
		}
	}
	d.Fields = fields
	return d
}

func convertIdsToMap(ids []*dvevaluation.DvVariable) map[string]int {
	n := len(ids)
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		s := dvevaluation.AnyToString(ids[i])
		m[s] = i
	}
	return m
}

func convertStringIdsToMap(ids []string) map[string]int {
	n := len(ids)
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		s := ids[i]
		m[s] = i
	}
	return m
}

func convertFieldsToMap(fields []string) map[string]int {
	n := len(fields)
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		m[fields[i]] = i
	}
	return m
}

func convertDvVariableFieldsToMap(fields []*dvevaluation.DvVariable) map[string]int {
	n := len(fields)
	m := make(map[string]int, n)
	for i := 0; i < n; i++ {
		s := dvevaluation.AnyToString(fields[i])
		m[s] = i
	}
	return m
}

func isKeyInMap(p *dvevaluation.DvVariable, key []byte, idMap map[string]int) bool {
	if p == nil || len(p.Fields) == 0 {
		return false
	}
	fields := p.Fields
	n := len(fields)
	for i := 0; i < n; i++ {
		e := fields[i]
		if e != nil && bytes.Equal(e.Name, key) {
			s := dvevaluation.AnyToString(e)
			_, ok := idMap[s]
			return ok
		}

	}
	return false
}

func deleteKeysInJson(path string, ids []string, keyFirst string) interface{} {
	n := len(ids)
	d, err := readWholeFileAsJson(path)
	if err != nil {
		return err
	}
	if d == nil || len(d.Fields) == 0 || d.Kind != dvevaluation.FIELD_ARRAY {
		return &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY, Fields: make([]*dvevaluation.DvVariable, 0, 0)}
	}
	oldFields := d.Fields
	m := len(oldFields)
	newFields := make([]*dvevaluation.DvVariable, 0, n)
	key := []byte(keyFirst)
	idMap := convertStringIdsToMap(ids)
	for i := 0; i < m; i++ {
		p := oldFields[i]
		if !isKeyInMap(p, key, idMap) {
			newFields = append(newFields, p)
		}
	}
	d.Fields = newFields
	err = writeWholeFileAsJson(path, d)
	return err
}

func writeWholeFileAsJson(path string, d *dvevaluation.DvVariable) error {
	s := dvevaluation.AnyToString(d)
	b := []byte(s)
	err := os.WriteFile(path, b, 0644)
	return err
}

func readFieldInJsonAsString(record *dvevaluation.DvVariable, key string) (string, bool) {
	if record == nil || record.Kind != dvevaluation.FIELD_OBJECT {
		return "", false
	}
	n := len(record.Fields)
	keyBytes := []byte(key)
	for i := 0; i < n; i++ {
		v := record.Fields[i]
		if v != nil && bytes.Equal(keyBytes, v.Name) {
			return dvevaluation.AnyToString(v), true
		}
	}
	return "", false
}

func setFieldInJsonAsString(record *dvevaluation.DvVariable, key string, value string) bool {
	if record == nil || record.Kind != dvevaluation.FIELD_OBJECT {
		return false
	}
	n := len(record.Fields)
	keyBytes := []byte(key)
	if n == 0 {
		record.Fields = make([]*dvevaluation.DvVariable, 0, 1)
	} else {
		for i := 0; i < n; i++ {
			v := record.Fields[i]
			if v != nil && bytes.Equal(keyBytes, v.Name) {
				v.Value = []byte(value)
				v.Kind = dvevaluation.FIELD_STRING
				return true
			}
		}
	}
	p := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_STRING, Name: keyBytes, Value: []byte(value)}
	record.Fields = append(record.Fields, p)
	return true
}

func createRecordInJson(path string, record *dvevaluation.DvVariable, keyFirst string, newId string) (*dvevaluation.DvVariable, error) {
	if !setFieldInJsonAsString(record, keyFirst, newId) {
		return nil, errors.New("Request body is not a JSON object")
	}
	pool, err := readWholeFileAsJson(path)
	if pool == nil || pool.Kind != dvevaluation.FIELD_ARRAY {
		pool = &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_ARRAY}
	}
	if pool.Fields == nil {
		pool.Fields = make([]*dvevaluation.DvVariable, 0, 3)
	}
	pool.Fields = append(pool.Fields, record)
	err = writeWholeFileAsJson(path, pool)
	return record, err
}

func retrieveUpdateInfo(path string, record *dvevaluation.DvVariable, keyFirst string) (*dvevaluation.DvVariable, int, error) {
	id, ok := readFieldInJsonAsString(record, keyFirst)
	if !ok || !checkIntId(id) {
		return nil, 0, errors.New("object has no id")
	}
	pool, err := readWholeFileAsJson(path)
	if err != nil {
		return nil, 0, err
	}
	if pool == nil || pool.Kind != dvevaluation.FIELD_ARRAY {
		return nil, 0, errors.New("table is not active yet")
	}
	i := findIndexOfJsonObject(pool.Fields, keyFirst, id)
	return pool, i, err
}

func updateRecordInJson(path string, record *dvevaluation.DvVariable, keyFirst string, version string) (*dvevaluation.DvVariable, error) {
	pool, i, err := retrieveUpdateInfo(path, record, keyFirst)
	if err != nil {
		return nil, err
	}
	if i < 0 {
		return nil, nil
	}
	resolveVersion(pool.Fields[i], record, version)
	pool.Fields[i] = record
	err = writeWholeFileAsJson(path, pool)
	return record, err
}

func findIndexOfJsonObject(fields []*dvevaluation.DvVariable, key string, id string) int {
	n := len(fields)
	for i := 0; i < n; i++ {
		v, ok := readFieldInJsonAsString(fields[i], key)
		if ok && v == id {
			return i
		}
	}
	return -1
}

func checkIntId(id string) bool {
	n := len(id)
	if n == 0 || n > 19 {
		return false
	}
	for i := 0; i < n; i++ {
		c := id[i]
		if !(c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func CreateEnvironmentForPreviousCurrent(previousRecord *dvevaluation.DvVariable, currentRecord *dvevaluation.DvVariable) *dvevaluation.DvObject {
	data := make(map[string]interface{})
        data["previous"]=previousRecord
        data["current"]=currentRecord
	env := dvparser.GetPropertiesPrototypedToGlobalProperties(data)
	return env
}

func findFirstMetCondition(previousRecord *dvevaluation.DvVariable, currentRecord *dvevaluation.DvVariable, conditions []string) (int, error) {
	n := len(conditions)
	var env *dvevaluation.DvObject = nil
	for i := 0; i < n; i++ {
		condition := conditions[i]
		if condition == "NEW" {
			continue
		}
		if condition == "DEFAULT" {
			return i, nil
		}
		if env == nil {
			env = CreateEnvironmentForPreviousCurrent(previousRecord, currentRecord)
		}
		r, err := env.EvaluateAnyTypeExpression(condition)
		if err != nil {
			return 0, errors.New("Error in expression " + condition + ":" + err.Error())
		}
		if dvevaluation.AnyToBoolean(r) {
			return i, nil
		}
	}
	return -1, nil
}

func updateRecordByFields(previousRecord *dvevaluation.DvVariable, currentRecord *dvevaluation.DvVariable, fields string) bool {
	fields = strings.TrimSpace(fields)
	if len(fields) == 0 {
		return true
	}
	c := fields[0]
	cleaning := false
	if c == '!' {
		cleaning = true
		fields = fields[1:]
		if len(fields) == 0 {
			return false
		}
	} else if c == '^' {
		fields = fields[1:]
		newPrevious := &dvevaluation.DvVariable{Kind: dvevaluation.FIELD_OBJECT, Fields: currentRecord.Fields}
		currentRecord.Fields = previousRecord.Fields
		previousRecord = newPrevious
		if len(fields) == 0 {
			return false
		}
	}
	fieldList := dvtextutils.ConvertToNonEmptyList(fields)
	if cleaning {
		currentRecord.Fields = previousRecord.Fields
		currentRecord.CleanFields(fieldList)
	} else {
		currentRecord.CopyFieldsFromOther(fieldList, previousRecord)
	}
	return true
}

func CreateOrUpdateByConditionsAndUpdateFieldsForJson(path string, record *dvevaluation.DvVariable, conditions []string, fields []string, keyFirst string, version string) (*dvevaluation.DvVariable, error) {
	pool, i, err := retrieveUpdateInfo(path, record, keyFirst)
	if err != nil {
		return nil, err
	}
	if i < 0 {
		if dvtextutils.IsStringContainedInArray("NEW", conditions) {
			if len(version) > 0 {
				setFieldInJsonAsString(record, version, "1")
			}
			pool.Fields = append(pool.Fields, record)
			err = writeWholeFileAsJson(path, pool)
			return record, err
		}
		return nil, nil
	}
	n, err := findFirstMetCondition(pool.Fields[i], record, conditions)
	if err != nil {
		return nil, err
	}
	if n < 0 {
		return pool.Fields[i], nil
	}
	changed := updateRecordByFields(pool.Fields[i], record, fields[n])
	if !changed {
		return pool.Fields[i], nil
	}
	resolveVersion(pool.Fields[i], record, version)
	pool.Fields[i] = record
	err = writeWholeFileAsJson(path, pool)
	return record, err
}
