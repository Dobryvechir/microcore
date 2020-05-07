/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import (
	"database/sql"
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

const (
	SQL_KIND_UPDATE     = 0
	SQL_KIND_SINGLE     = 1
	SQL_KIND_ROW        = 2
	SQL_KIND_TABLE      = 3
	SQL_KIND_LIST       = 4
	SQL_KIND_ROW_TEXT   = 5
	SQL_KIND_TABLE_TEXT = 6
)

type SqlAction struct {
	Db           string   `json:"db"`
	Query        string   `json:"query"`
	QueryOracle  string   `json:"queryOracle"`
	QueryPostgre string   `json:"queryPostgre"`
	Result       string   `json:"result"`
	Kind         string   `json:"kind"`
	Columns      []string `json:"columns"`
	Empty        int      `json:"empty"`
	Error        string   `json:"error"`
	KindNo       int
}

func SqlInit(command string, ctx *dvmeta.RequestContext) ([]interface{}, bool) {
	p := strings.Index(command, ":")
	command = strings.TrimSpace(command[p+1:])
	sqlAction := &SqlAction{}
	err := json.Unmarshal([]byte(command), sqlAction)
	if err != nil {
		dvlog.PrintfError("Error %v", err)
		return nil, false
	}
	kind := strings.TrimSpace(strings.ToUpper(sqlAction.Kind))
	switch kind {
	case "", "UPDATE":
		sqlAction.KindNo = SQL_KIND_UPDATE
		break
	case "SINGLE":
		sqlAction.KindNo = SQL_KIND_SINGLE
		break
	case "ROW":
		sqlAction.KindNo = SQL_KIND_ROW
		break
	case "TABLE":
		sqlAction.KindNo = SQL_KIND_TABLE
		break
	case "LIST":
		sqlAction.KindNo = SQL_KIND_LIST
		break
	case "ROW_TEXT":
		sqlAction.KindNo = SQL_KIND_ROW_TEXT
		break
	case "TEXT":
		sqlAction.KindNo = SQL_KIND_TABLE_TEXT
		break
	default:
		dvlog.PrintfError("Unknown kind: %s", kind)
		return nil, false
	}
	return []interface{}{sqlAction, ctx}, true
}

func SqlRun(data []interface{}) bool {
	sqlAction := data[0].(*SqlAction)
	ctx := data[1].(*dvmeta.RequestContext)
	db, sqlName, err := GetDBConnection(dvparser.GlobalProperties, sqlAction.Db)
	if err != nil {
		dvlog.PrintfError("Connection to %s failed %v", sqlAction.Db, err)
		return false
	}
	query := sqlAction.Query
	switch strings.ToUpper(sqlName) {
	case "ORACLE":
		if sqlAction.QueryOracle != "" {
			query = sqlAction.QueryOracle
		}
	case "POSTGRE":
		if sqlAction.QueryPostgre != "" {
			query = sqlAction.QueryPostgre
		}
	}
	var res interface{} = nil
	kind := sqlAction.KindNo
	if kind == SQL_KIND_UPDATE {
		res, err = db.Exec(query)
	} else {
		var rs *sql.Rows
		rs, err = db.Query(query)
		if err == nil {
			switch kind {
			case SQL_KIND_SINGLE:
				if rs.Next() {
					err = rs.Scan(&res)
				}
				break
			case SQL_KIND_LIST:
				{
					data := make([]string, 0, 1024)
					for rs.Next() {
						var s string
						err = rs.Scan(&s)
						if err != nil {
							break
						}
						data = append(data, s)
					}
					res = data
				}
				break
			case SQL_KIND_ROW, SQL_KIND_ROW_TEXT:
				{
					var dataCol map[string]string = nil
					var dataText []string = nil
					columns := sqlAction.Columns
					if len(columns) == 0 {
						columns, err = rs.Columns()
						if err != nil {
							break
						}
					}
					n := len(columns)
					cols := make([]interface{}, n)
					if rs.Next() {
						r := make([]string, n)
						for i := 0; i < n; i++ {
							cols[i] = &r[i]
						}
						rs.Scan(cols...)
						if kind == SQL_KIND_ROW {
							m := make(map[string]string, n)
							for i := 0; i < n; i++ {
								m[columns[i]] = r[i]
							}
							dataCol = m
						} else {
							dataText = r
						}
					}
					if kind == SQL_KIND_ROW {
						res = dataCol
					} else {
						res = dataText
					}
				}
				break
			case SQL_KIND_TABLE, SQL_KIND_TABLE_TEXT:
				{
					var dataCol []map[string]string = nil
					var dataText [][]string = nil
					columns := sqlAction.Columns
					if len(columns) == 0 {
						columns, err = rs.Columns()
						if err != nil {
							break
						}
					}
					n := len(columns)
					cols := make([]interface{}, n)
					for rs.Next() {
						r := make([]string, n)
						for i := 0; i < n; i++ {
							cols[i] = &r[i]
						}
						rs.Scan(cols...)
						if kind == SQL_KIND_ROW {
							m := make(map[string]string, n)
							for i := 0; i < n; i++ {
								m[columns[i]] = r[i]
							}
							dataCol = append(dataCol, m)
						} else {
							dataText = append(dataText, r)
						}
					}
					if kind == SQL_KIND_ROW {
						res = dataCol
					} else {
						res = dataText
					}
				}
			}
		}
		rs.Close()
	}
	if err != nil {
		if sqlAction.Error != "" {
			ctx.ExtraAsDvObject.Set(sqlAction.Error, err.Error())
			return true
		} else {
			dvlog.PrintfError("Error %s: %v", query, err)
			return false
		}
	}
	if res == nil && sqlAction.Empty != 0 {
		if ctx == nil {
			dvlog.PrintfError("Empty result of %s %d", query, sqlAction.Empty)
		} else {
			ctx.SetHttpErrorCode(sqlAction.Empty, "")
		}
		return false
	}
	if sqlAction.Result != "" {
		ctx.ExtraAsDvObject.Set(sqlAction.Result, res)
	}
	return true
}
