/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strings"
)

const (
	SQL_ORACLE_LIKE      = 1
	SQL_POSTGRES_LIKE    = 2
	COMMON_MAX_BATCH     = 1000
	COMPLEX_ID_SEPARATOR = "_._"
)

type TableMetaData struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Dependencies []int    `json:"dependencies"`
	IdColumns    []int    `json:"idColumn"`
	MajorColumn  int      `json:"majorColumn"`
	Types        []string `json:"types"`
	Columns      []string `json:"columns"`
}

func OrderObjectsByHierarchy(objects [][]string, leftObjects map[string]bool,
	idCol int, depCols []int) ([][][]string, error) {
	n := len(objects)
	v := len(depCols)
	mapping := make(map[string][]string)
	pool := make([][]string, n)
	poolSize := 0
	for i := 0; i < n; i++ {
		item := objects[i]
		id := item[idCol]
		if leftObjects != nil && leftObjects[id] {
			continue
		}
		pool[poolSize] = item
		poolSize++
		mapping[id] = item
	}
	res := make([][][]string, 0, n)

	for poolSize > 0 {
		current := make([][]string, 0, 16)
		m := 0
		for i := 0; i < poolSize; i++ {
			item := pool[i]
			reliant := false
			for k := 0; k < v; k++ {
				depend := item[depCols[k]]
				if depend != "" && mapping[depend] != nil {
					reliant = true
					break
				}
			}
			if !reliant {
				current = append(current, item)
			} else {
				pool[m] = item
				m++
			}
		}
		poolSize = m
		m = len(current)
		if m == 0 {
			return nil, errors.New("Cycled objects")
		}
		for i := 0; i < m; i++ {
			delete(mapping, current[i][idCol])
		}
		res = append(res, current)
	}
	return res, nil
}

func appendSlashed(b []byte, v []byte) []byte {
	n := len(v)
	for i := 0; i < n; i++ {
		d := v[i]
		if d == '\'' {
			b = append(b, '\\', '\'')
		} else {
			b = append(b, d)
		}
	}
	return b
}

func GetComplexIdForItem(row []string, ids []int) string {
	n := len(row)
	m := len(ids)
	s := ""
	pref := ""
	for i := 0; i < m; i++ {
		k := ids[i]
		if k < n {
			s += pref + row[k]
			pref = COMPLEX_ID_SEPARATOR
		}
	}
	return s
}

func SavePortionOfItems(items [][]string, sqlTable string, conn *sql.DB, left map[string]bool,
	columnIds []int, options int, types []string) error {
	oracleLike := (options & SQL_ORACLE_LIKE) != 0
	postgresLike := (options & SQL_POSTGRES_LIKE) != 0
	maxBatch := 1
	if oracleLike {
		maxBatch = COMMON_MAX_BATCH
	} else if postgresLike {
		maxBatch = 5000
	}
	cols := len(types)
	sqlStart := "INSERT INTO " + sqlTable + " "
	sqlEnd := ""
	if oracleLike {
		sqlStart += "with dv_sql_db_oracle as ("
		sqlEnd = ") select * from dv_sql_db_oracle"
	} else {
		sqlStart += "VALUES"
	}
	n := len(items)
	cn := n
	if cn > maxBatch {
		cn = maxBatch
	}
	cn = cn * cols * 80
	b := make([]byte, 0, cn)
	count := 0
	for j := 0; j < n; j++ {
		s := items[j]
		id := GetComplexIdForItem(s, columnIds)
		if left != nil && left[id] {
			continue
		}
		if count != 0 {
			if oracleLike {
				b = append(b, []byte(" UNION ALL ")...)
			} else {
				b = append(b, ',')
			}
		}
		if oracleLike {
			b = append(b, []byte("SELECT ")...)
		} else {
			b = append(b, '(')
		}
		count++
		for i := 0; i < cols; i++ {
			if i != 0 {
				b = append(b, ',')
			}
			v := s[i]
			if v == "null" {
				b = append(b, []byte("null")...)
			} else {
				tp := types[i]
				if tp == "Date" {
				}
				b = append(b, '\'')
				b = appendSlashed(b, []byte(v))
				b = append(b, '\'')
			}
		}
		if oracleLike {
			b = append(b, []byte(" from dual")...)
		} else {
			b = append(b, ')')
		}
		if count >= maxBatch {
			sqlFull := sqlStart + string(b) + sqlEnd
			if logPreExecuteLevel >= dvlog.LogDebug {
				log.Printf("Sql execution: %s", sqlFull)
			}
			_, err := conn.Exec(sqlFull)
			b = b[0:0]
			count = 0
			if err != nil {
				if logPreExecuteLevel >= dvlog.LogDebug {
					log.Printf("Failed to execute sql %s: %v", sqlFull, err)
				}
				return err
			}
		}
	}
	if count > 0 {
		sqlFull := sqlStart + string(b) + sqlEnd
		if logPreExecuteLevel >= dvlog.LogDebug {
			log.Printf("Sql execution: %s", sqlFull)
		}
		_, err := conn.Exec(sqlFull)
		if err != nil {
			if logPreExecuteLevel >= dvlog.LogDebug {
				log.Printf("Failed to execute sql %s: %v", sqlFull, err)
			}
			return err
		}
	}
	return nil
}

func ReadTableMetaData(table string, props map[string]string) (*TableMetaData, error) {
	jsonDef := strings.TrimSpace(props["DB_SQL_META_"+table])
	if len(jsonDef) == 0 {
		return nil, errors.New("No table definition for " + table)
	}
	meta := &TableMetaData{}
	err := json.Unmarshal([]byte(jsonDef), meta)
	if err != nil {
		return nil, err
	}
	meta.Id = table
	if len(meta.IdColumns) == 1 {
		meta.MajorColumn = meta.IdColumns[0]
	} else if len(meta.IdColumns) == 0 {
		meta.MajorColumn = -1
	} else if meta.MajorColumn >= 0 {
		if FindIntInIntArray(meta.MajorColumn, meta.IdColumns) < 0 {
			return nil, errors.New("majorColumn for " + table + " must be inside idColumns")
		}
	}
	return meta, nil
}

func GetIdsFromItems(meta *TableMetaData, items [][]string) [][]string {
	idIndices := meta.IdColumns
	n := len(items)
	m := len(idIndices)
	if m == 0 {
		return nil
	}
	ids := make([][]string, n)
	for i := 0; i < n; i++ {
		ids[i] = make([]string, m)
		for j := 0; j < m; j++ {
			ids[i][j] = items[i][idIndices[j]]
		}
	}
	return ids
}

func GetColumnListFromMetaByIndices(meta *TableMetaData, indices []int) string {
	n := len(indices)
	s := ""
	for i := 0; i < n; i++ {
		if i != 0 {
			s += ","
		}
		s += meta.Columns[indices[i]]
	}
	return s
}

func FindIntInIntArray(val int, data []int) int {
	n := len(data)
	for i := 0; i < n; i++ {
		if data[i] == val {
			return i
		}
	}
	return -1
}

func GetSingleValuesFromString(data [][]string, column int) []string {
	n := len(data)
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = data[i][column]
	}
	return r
}

func makeSingleValueArray(items [][]string, separ string) []string {
	n := len(items)
	r := make([]string, n)
	for i := 0; i < n; i++ {
		r[i] = strings.Join(items[i], separ)
	}
	return r
}

func GetExistingItems(meta *TableMetaData, ids [][]string, db *sql.DB) ([]string, error) {
	if ids == nil || meta.MajorColumn < 0 {
		return nil, nil
	}
	idIndices := meta.IdColumns
	n := len(idIndices)
	var sqlStart string
	sqlEnd := ")"
	col := 0
	idCol := meta.Columns[meta.MajorColumn]
	if n == 1 {
		sqlStart = "SELECT " + idCol + " FROM " + meta.Name + " WHERE " + idCol + " in ("
	} else {
		sqlStart = "SELECT " + GetColumnListFromMetaByIndices(meta, idIndices) + " FROM " + meta.Name + " WHERE " + idCol + " in ("
		col = FindIntInIntArray(meta.MajorColumn, idIndices)
	}
	singleIds := GetSingleValuesFromString(ids, col)
	if logPreExecuteLevel >= dvlog.LogDetail {
		log.Printf("Getting existing items %s?%s for %v", sqlStart, sqlEnd, singleIds)
	}
	items, err := ReadItemsInBatches(db, sqlStart, sqlEnd, singleIds, n)
	if err != nil {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Failed to read items in batches %s?%s: %v", sqlStart, sqlEnd, err)
		}
		return nil, err
	}
	return makeSingleValueArray(items, COMPLEX_ID_SEPARATOR), nil
}

func ConvertListToBooleanMap(ids []string) map[string]bool {
	n := len(ids)
	m := make(map[string]bool)
	for i := 0; i < n; i++ {
		m[ids[i]] = true
	}
	return m
}

func PreExecuteCsvFile(conn *sql.DB, name string, options int) error {
	data, err := ReadCsvFromFile(name)
	if err != nil {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Failed to parse csv file %s : %v", name, err)
		}
		return err
	}
	props := dvparser.GlobalProperties
	tables := dvparser.ConvertToNonEmptyList(props["DB_CSV_TABLE_ALIASES"])
	if logPreExecuteLevel >= dvlog.LogDetail {
		log.Printf("csv %s has %d tables, target tables are %d: %v\n", name, len(data), len(tables), tables)
	}
	for _, table := range tables {
		items := data[table]
		n := len(items)
		if logPreExecuteLevel >= dvlog.LogDetail {
			log.Printf("Into table %s %n items will be saved\n", table, n)
		}
		if n > 0 {
			meta, err := ReadTableMetaData(table, props)
			if err != nil {
				if logPreExecuteLevel >= dvlog.LogError {
					log.Printf("For table %s metadata is wrong: %v\n", table, err)
				}
				return err
			}
			ids := GetIdsFromItems(meta, items)
			leftList, err := GetExistingItems(meta, ids, conn)
			left := ConvertListToBooleanMap(leftList)
			if err != nil {
				if logPreExecuteLevel >= dvlog.LogError {
					log.Printf("For table %s failed to find existing items: %v\n", table, err)
				}
				return err
			}
			if len(meta.Dependencies) == 0 || len(meta.IdColumns) != 1 {
				if logPreExecuteLevel >= dvlog.LogDetail {
					log.Printf("Single blocks of %s are saved", table)
				}
				err = SavePortionOfItems(items, meta.Name, conn, left, meta.IdColumns, options, meta.Types)
				if err != nil {
					if logPreExecuteLevel >= dvlog.LogError {
						log.Printf("For table %s failed to write existing items in single block: %v\n", table, err)
					}
					return err
				}
			} else {
				groups, err := OrderObjectsByHierarchy(items, left, meta.IdColumns[0], meta.Dependencies)
				if err != nil {
					if logPreExecuteLevel >= dvlog.LogError {
						log.Printf("For table %s failed to divide items into groups: %v\n", table, err)
					}
					return err
				}
				n = len(groups)
				if logPreExecuteLevel >= dvlog.LogDetail {
					log.Printf("%n blocks of %s are saved", n, table)
				}
				for i := 0; i < n; i++ {
					err = SavePortionOfItems(groups[i], meta.Name, conn, left, meta.IdColumns, options, meta.Types)
					if err != nil {
						if logPreExecuteLevel >= dvlog.LogError {
							log.Printf("For table %s failed to write %dth block: %v\n", table, n, err)
						}
						return err
					}
				}
			}
			// TODO add updates for the left
		}
	}
	return nil
}
