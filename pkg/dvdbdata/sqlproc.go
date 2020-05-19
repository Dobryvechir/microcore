// package dvdbdata provides functions for sql query
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbdata

import (
	"database/sql"
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"strings"
)

const (
	SqlKindUpdate    = 0
	SqlKindSingle    = 1
	SqlKindRow       = 2
	SqlKindTable     = 3
	SqlKindList      = 4
	SqlKindRowText   = 5
	SqlKindTableText = 6
)

type SqlAction struct {
	Db            string   `json:"db"`
	Query         string   `json:"query"`
	QueryOracle   string   `json:"queryOracle"`
	QueryPostgres string   `json:"queryPostgres"`
	Result        string   `json:"result"`
	Kind          string   `json:"kind"`
	Columns       []string `json:"columns"`
	Empty         int      `json:"empty"`
	Error         string   `json:"error"`
	KindNo        int
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
		sqlAction.KindNo = SqlKindUpdate
		break
	case "SINGLE":
		sqlAction.KindNo = SqlKindSingle
		break
	case "ROW":
		sqlAction.KindNo = SqlKindRow
		break
	case "TABLE":
		sqlAction.KindNo = SqlKindTable
		break
	case "LIST":
		sqlAction.KindNo = SqlKindList
		break
	case "ROW_TEXT":
		sqlAction.KindNo = SqlKindRowText
		break
	case "TEXT":
		sqlAction.KindNo = SqlKindTableText
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
	db, err := GetDBConnection(sqlAction.Db)
	if err != nil {
		dvlog.PrintfError("Connection to %s failed %v", sqlAction.Db, err)
		return false
	}
	query := sqlAction.Query
	switch db.KindMask {
	case SqlOracleLike:
		if sqlAction.QueryOracle != "" {
			query = sqlAction.QueryOracle
		}
	case SqlPostgresLike:
		if sqlAction.QueryPostgres != "" {
			query = sqlAction.QueryPostgres
		}
	}
	var res interface{} = nil
	kind := sqlAction.KindNo
	if kind == SqlKindUpdate {
		res, err = db.Exec(query)
	} else {
		var rs *sql.Rows
		rs, err = db.Query(query)
		if err == nil {
			switch kind {
			case SqlKindSingle:
				if rs.Next() {
					err = rs.Scan(&res)
				}
				break
			case SqlKindList:
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
			case SqlKindRow, SqlKindRowText:
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
						err = rs.Scan(cols...)
						if kind == SqlKindRow {
							m := make(map[string]string, n)
							for i := 0; i < n; i++ {
								m[columns[i]] = r[i]
							}
							dataCol = m
						} else {
							dataText = r
						}
					}
					if kind == SqlKindRow {
						res = dataCol
					} else {
						res = dataText
					}
				}
				break
			case SqlKindTable, SqlKindTableText:
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
						err = rs.Scan(cols...)
						if kind == SqlKindRow {
							m := make(map[string]string, n)
							for i := 0; i < n; i++ {
								m[columns[i]] = r[i]
							}
							dataCol = append(dataCol, m)
						} else {
							dataText = append(dataText, r)
						}
					}
					if kind == SqlKindRow {
						res = dataCol
					} else {
						res = dataText
					}
				}
			}
		}
		if rs != nil {
			rs.Close()
		}
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
