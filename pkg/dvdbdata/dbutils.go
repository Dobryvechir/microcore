/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
	"strings"
)

var readPropertySql, readPropertySqlEnd string
var createPropertySql []string = nil
var updatePropertySql []string = nil


func GetTableNameColumnsFromDefinition(def string) (table string, columns []string, colDef []string, err error) {
	first := strings.Index(def, "(")
	last := strings.LastIndex(def, ")")
	if first <= 0 || last < first {
		err = errors.New("Columns must be in round brackets: " + def)
		return
	}
	table = strings.TrimSpace(def[:first])
	cols := strings.TrimSpace(def[first+1 : last])
	if table == "" || cols == "" {
		err = errors.New("Tables and columns cannot be empty: " + def)
		return
	}
	colDef = strings.Split(cols, ",")
	n := len(colDef)
	columns = make([]string, n)
	for i := 0; i < n; i++ {
		s := strings.TrimSpace(colDef[i])
		colDef[i] = s
		p := strings.Index(s, " ")
		if p > 0 {
			s = s[:p]
		}
		columns[i] = s
	}
	return
}

func GetPropertyGlobalDefinition(props map[string]string) string {
	r := props[propertyTableDefinitionName]
	if r == "" {
		r = propertyTableDefinitionDefault
	}
	return r
}

func CreateTableByDefinition(db *DBConnection, def string) error {
	table, columns, colDefs, err := GetTableNameColumnsFromDefinition(def)
	if err != nil {
		return err
	}
	primary := ""
	query := "CREATE TABLE " + table + "("
	query1 := query
	for i, col := range colDefs {
		if i != 0 {
			query += ","
			query1 += ","
		}
		if strings.HasSuffix(col, "primary") {
			query += col + " key"
			query1 += col[:len(col)-7]
			if primary == "" {
				primary = columns[i]
			} else {
				primary += "," + columns[i]
			}
		} else {
			query += col
			query1 += col
		}
	}
	query += ")"
	if primary != "" {
		query1 += ",CONSTRAINT " + table + "_PK PRIMARY KEY (" + primary + ")"
	}
	query1 += ")"
	_, err = db.Exec(query)
	if err == nil {
		return nil
	}
	_, err1 := db.Exec(query1)
	if err1 == nil {
		return nil
	}
	return errors.New(err.Error() + " or (" + err1.Error() + ")")
}

func ReadGlobalDBProperty(props map[string]string, db *DBConnection, name string, defValue string) (string, error) {
	if readPropertySql == "" {
		table, columns, _, err := GetTableNameColumnsFromDefinition(GetPropertyGlobalDefinition(props))
		if err != nil {
			return "", err
		}
		if len(columns) < 2 {
			return "", errors.New("Global Property table must have at least 2 columns")
		}
		readPropertySql = "SELECT " + columns[1] + " FROM " + table + " WHERE " + columns[0] + "='"
		readPropertySqlEnd = "'"
	}
	query := readPropertySql + name + readPropertySqlEnd
	rs, err := db.Query(query)
	if err != nil {
		err1 := CreateTableByDefinition(db, GetPropertyGlobalDefinition(props))
		if err1 == nil {
			rs, err = db.Query(query)
		}
	}
	if err != nil {
		return "", err
	}
	if rs.Next() {
		var name string
		err = rs.Scan(&name)
		rs.Close()
		if err != nil {
			return "", err
		}
		return name, nil
	}
	rs.Close()
	return defValue, nil
}

func WriteGlobalDBProperty(props map[string]string, db *DBConnection, name string, value string) error {
	if createPropertySql == nil {
		table, columns, _, err := GetTableNameColumnsFromDefinition(GetPropertyGlobalDefinition(props))
		if err != nil {
			return err
		}
		if len(columns) < 2 {
			return errors.New("Global Property table must have at least 2 columns")
		}
		createPropertySql = make([]string, 3)
		createPropertySql[0] = "INSERT INTO " + table + "(" + columns[0] + "," + columns[1] + ") VALUES('"
		createPropertySql[1] = "','"
		createPropertySql[2] = "')"
		updatePropertySql = make([]string, 3)
		updatePropertySql[0] = "UPDATE " + table + " SET " + columns[1] + "='"
		updatePropertySql[1] = "' WHERE " + columns[0] + "='"
		updatePropertySql[2] = "'"
	}
	query := createPropertySql[0] + name + createPropertySql[1] + value + createPropertySql[2]
	_, err := db.Exec(query)
	if err != nil {
		err1 := CreateTableByDefinition(db, GetPropertyGlobalDefinition(props))
		if err1 != nil {
			query = updatePropertySql[0] + value + updatePropertySql[1] + name + updatePropertySql[2]
		}
		_, err = db.Exec(query)
	}
	return err
}

func AddItemsToPool(db *DBConnection, sql string, cols int, pool [][]string) ([][]string, error) {
	if logPreExecuteLevel >= dvlog.LogTrace {
		log.Printf("Add sql rows to pool: %s (%d columns)", sql, cols)
	}
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if cols <= 0 {
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		cols = len(columns)
	}
	data := make([]interface{}, cols)
	for rows.Next() {
		items := make([]string, cols)
		for i := 0; i < cols; i++ {
			data[i] = &items[i]
		}
		if err := rows.Scan(data...); err != nil {
			log.Fatal(err)
		}
		pool = append(pool, items)
	}
	return pool, nil
}

func ReadItemsInBatches(db *DBConnection, start string, finish string, ids []string, cols int) ([][]string, error) {
	pool := make([][]string, 0, 1024)
	n := len(ids)
	i := 0
	var err error
	for n > 0 {
		m := n
		if m > CommonMaxBatch {
			m = CommonMaxBatch
		}
		n -= m
		sqlQuery := start + strings.Join(ids[i:i+m], ",") + finish
		pool, err = AddItemsToPool(db, sqlQuery, cols, pool)
		if err != nil {
			return nil, err
		}
	}
	return pool, nil
}
