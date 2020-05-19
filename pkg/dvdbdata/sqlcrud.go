/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import (
	"errors"
	"strings"
)

func SqlSingleValueByConnectionName(connName string, query string) (string, bool, error) {
	db, err := GetDBConnection(connName)
	if err != nil {
		return "", false, err
	}
	return SqlSingleValueByConnection(db, query)
}

func SqlSingleValueByConnection(db *DBConnection, query string) (string, bool, error) {
	rs, err := db.Query(query)
	if err != nil {
		return "", false, err
	}
	if rs.Next() {
		var r string
		err = rs.Scan(&r)
		if err != nil {
			return "", false, err
		}
		return r, true, nil
	}
	return "", false, nil
}

func SqlUpdateByConnectionName(connName string, query string) error {
	db, err := GetDBConnection(connName)
	if err != nil {
		return err
	}
	return SqlUpdateByConnection(db, query)
}

func SqlUpdateByConnection(db *DBConnection, query string) error {
	_, err := db.Exec(query)
	return err
}

func GetSqlTableByIds(db *DBConnection, tableId string, ids []string) ([][]string, error) {
	metaInfo, err := ReadTableMetaDataFromGlobal(tableId)
	if err != nil {
		return nil, err
	}
	start, finish, cols, err := GetSqlQueryForGettingRowById(metaInfo, "")
	if err != nil {
		return nil, err
	}
	return ReadItemsInBatches(db, start, finish, ids, cols)
}

func GetSqlQueryForGettingRowById(metaInfo *TableMetaData, columns string) (start string, finish string, cols int, err error) {
	if columns == "" {
		columns = "*"
	}
	if metaInfo.MajorColumn < 0 {
		err = errors.New("Table " + metaInfo.Name + " has no major column")
		return
	}
	colName := metaInfo.Columns[metaInfo.MajorColumn]
	if metaInfo.QuoteColumns {
		colName = "\"" + colName + "\""
	}
	start = "SELECT " + columns + " FROM " + metaInfo.Name + " WHERE " + colName + " in ("
	finish = ")"
	cols = len(metaInfo.Columns)
	if columns != "*" {
		cols = strings.Count(columns, ",") + 1
	}
	return
}

func ReadItemsInBulk(db *DBConnection, query string, cols int) ([][]string, error) {
	return AddItemsToPool(db, query, cols, make([][]string, 0, 100))
}

func FindIdsPlaceholder(s string) (outerStart int, outerFinish int, innerStart int, innerFinish int) {
	innerFinish = strings.Index(s, IdsPlaceholderFinish)
	outerStart = strings.LastIndex(s[:innerFinish+1], IdsPlaceholderStart)
	if outerStart < 0 || innerFinish < outerStart {
		innerFinish = -1
		innerStart = -1
		outerStart = -1
		outerFinish = -1
	} else {
		outerFinish = innerFinish + len(IdsPlaceholderFinish)
		innerStart = outerStart + len(IdsPlaceholderStart)
	}
	return
}

func GetSqlTableByQuery(db *DBConnection, ids []string, query string) ([][]string, error) {
	posStart, posFinish, _, _ := FindIdsPlaceholder(query)
	cols := 0
	if posStart >= 0 {
		start := query[:posStart]
		finish := query[posFinish:]
		return ReadItemsInBatches(db, start, finish, ids, cols)
	}
	return ReadItemsInBulk(db, query, cols)
}

func CollectAllChildInfo(tableId string, data [][]string) (map[string][]string, error) {
	metaInfo, err := ReadTableMetaDataFromGlobal(tableId)
	if err != nil {
		return nil, err
	}
	r := make(map[string][]string)
	refs := metaInfo.References
	n := len(refs)
	m := len(data)
	if n > 0 && m > 0 && len(data[0]) > 0 {
		if n > len(data[0]) {
			n = len(data[0])
		}
		for i := 0; i < n; i++ {
			refName := refs[i]
			if refName != "" && refName[0] != '_' {
				pool := make([]string, m)
				for j := 0; j < m; j++ {
					pool[j] = data[j][i]
				}
				r[refName] = append(r[refName], pool...)
			}
		}
	}
	return r, nil
}
