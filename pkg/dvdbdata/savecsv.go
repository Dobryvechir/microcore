// package dvdbdata orchestrates sql database access
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvcsv"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strconv"
	"strings"
)

func findDependantId(row []string, mapping map[string][]string, depCols []int, idCol int) (string, int) {
	n := len(depCols)
	for i := 0; i < n; i++ {
		k := depCols[i]
		id := row[k]
		if mapping[id] != nil && id != row[idCol] {
			return id, k
		}
	}
	return "", -1
}

func collectCycleInfo(pool [][]string, depCols []int, idCol int) string {
	mapping := make(map[string][]string)
	passed := make(map[string]bool)
	n := len(pool)
	if n == 0 {
		return "Pool of cycled objects is empty"
	}
	var id string
	for i := 0; i < n; i++ {
		p := pool[i]
		id = p[idCol]
		mapping[id] = p
	}
	info := id + " -> "
	for !passed[id] {
		passed[id] = true
		nId, nCol := findDependantId(mapping[id], mapping, depCols, idCol)
		if nCol < 0 {
			info += " System error - no cycling"
			break
		}
		info += nId + "(" + strconv.Itoa(nCol) + ") -> "
		id = nId
	}
	return info
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
				if depend != "" && mapping[depend] != nil && depend != item[idCol] {
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
			info := collectCycleInfo(pool[:poolSize], depCols, idCol)
			return nil, errors.New("cycled objects (" + info + ")")
		}
		for i := 0; i < m; i++ {
			delete(mapping, current[i][idCol])
		}
		res = append(res, current)
	}
	return res, nil
}

func appendSlashed(b []byte, v []byte, options int) []byte {
	if (options & SqlPostgresLike) != 0 {
		return appendSlashedPostgresLike(b, v)
	}
	return appendSlashedOracleLike(b, v)
}

func appendSlashedPostgresLike(b []byte, v []byte) []byte {
	n := len(v)
	for i := 0; i < n; i++ {
		d := v[i]
		if d == '\'' {
			b = append(b, '\'', '\'')
		} else {
			b = append(b, d)
		}
	}
	return b
}

func appendSlashedOracleLike(b []byte, v []byte) []byte {
	n := len(v)
	for i := 0; i < n; i++ {
		d := v[i]
		if d == '\'' || d == '\\' {
			b = append(b, '\\', d)
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
			pref = ComplexIdSeparator
		}
	}
	return s
}

func PlaceStringToSqlQuery(v string, tp string, b []byte, options int) ([]byte, error) {
	switch tp {
	case TypeDate:
		{
			if len(v) == 0 {
				b = append(b, NullStringAsBytes...)
			} else {
				b = append(b, '\'')
				b = appendSlashed(b, []byte(v), options)
				b = append(b, '\'')
			}
		}
	case TypeBool:
		{
			if len(v) == 0 {
				b = append(b, NullStringAsBytes...)
			} else if dvparser.IsDigitOnly(v) && len(v) == 1 {
				b = append(b, '\'')
				b = append(b, []byte(v)...)
				b = append(b, '\'')
			} else {
				return b, fmt.Errorf("not a one-digit number %s of type %s", v, tp)
			}
		}
	case TypeInt, TypeInt64:
		{
			if len(v) == 0 {
				b = append(b, NullStringAsBytes...)
			} else if dvparser.IsSignAndDigitsOnly(v) {
				b = append(b, '\'')
				b = append(b, []byte(v)...)
				b = append(b, '\'')
			} else {
				return b, fmt.Errorf("not a number %s of type %s", v, tp)
			}
		}
	default:
		{
			b = append(b, '\'')
			b = appendSlashed(b, []byte(v), options)
			b = append(b, '\'')
		}
	}
	return b, nil
}

func SavePortionOfItems(items [][]string, sqlTable string, conn *DBConnection, left map[string]bool,
	columnIds []int, options int, types []string) (err error) {
	oracleLike := (options & SqlOracleLike) != 0
	postgresLike := (options & SqlPostgresLike) != 0
	maxBatch := 1
	if oracleLike {
		maxBatch = CommonMaxBatch
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
	if logPreExecuteLevel >= dvlog.LogTrace {
		log.Printf("\nInserting %s(%d): %d by %d\n", sqlTable, cols, n, cn)
	}
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
		currCols := len(s)
		for i := 0; i < cols; i++ {
			if i != 0 {
				b = append(b, ',')
			}
			var v string
			if i < currCols {
				v = s[i]
			} else {
				v = ""
			}
			if v == "null" {
				b = append(b, NullStringAsBytes...)
			} else {
				tp := types[i]
				b, err = PlaceStringToSqlQuery(v, tp, b, options)
				if err != nil {
					return
				}
				if logPreExecuteLevel >= dvlog.LogTrace {
					log.Printf("[%d:%s=%s]", i, tp, v)
				}
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

func ReadTableMetaDataFromGlobal(tableId string) (*TableMetaData, error) {
	return ReadTableMetaData(tableId, dvparser.GlobalProperties)
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

func makeSingleValueArray(items [][]interface{}, separ string) []string {
	n := len(items)
	r := make([]string, n)
	for i := 0; i < n; i++ {
		s := dvevaluation.ConvertInterfaceListToStringList(items[i], dvevaluation.ConversionOptionSimpleLike)
		r[i] = strings.Join(s, separ)
	}
	return r
}

func GetExistingItems(meta *TableMetaData, ids [][]string, db *DBConnection) ([]string, error) {
	if ids == nil || meta.MajorColumn < 0 {
		return nil, nil
	}
	idIndices := meta.IdColumns
	n := len(idIndices)
	var sqlStart string
	sqlEnd := ")"
	col := 0
	idCol := meta.Columns[meta.MajorColumn]
	if meta.QuoteColumns {
		idCol = "\"" + idCol + "\""
	}
	if n == 1 {
		sqlStart = "SELECT " + idCol + " FROM " + meta.Name + " WHERE " + idCol + " in ("
	} else {
		name := GetColumnListFromMetaByIndices(meta, idIndices)
		if meta.QuoteColumns {
			name = "\"" + name + "\""
		}
		sqlStart = "SELECT " + name + " FROM " + meta.Name + " WHERE " + idCol + " in ("
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
	return makeSingleValueArray(items, ComplexIdSeparator), nil
}

func ConvertListToBooleanMap(ids []string) map[string]bool {
	n := len(ids)
	m := make(map[string]bool)
	for i := 0; i < n; i++ {
		m[ids[i]] = true
	}
	return m
}

func GetMetaInfo(meta *TableMetaData) string {
	n := len(meta.Columns)
	s := meta.Name + "(" + meta.Id + ":" + strconv.Itoa(n) + ")["
	if len(meta.Types) != n {
		k := len(meta.Types)
		s += " Error - different lengths of column names (" + strconv.Itoa(n) + " and types (" + strconv.Itoa(k) + ")"
		if n > k {
			n = k
		}
	}
	for i := 0; i < n; i++ {
		s += "(" + strconv.Itoa(i) + "-" + meta.Columns[i] + ":" + meta.Types[i] + ")"
	}
	s += "]"
	return s
}

func PreExecuteCsvFile(conn *DBConnection, name string, options int) error {
	data, err := dvcsv.ReadCsvFromFile(name)
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
	isSingleRequests := props["DB_PRELOAD_SINGLE_REQUEST"] == "true"
	for _, table := range tables {
		items := data[table]
		n := len(items)
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
			if logPreExecuteLevel >= dvlog.LogDetail {
				log.Printf("Into table %s %d items will be saved\n", table, n)
				if logPreExecuteLevel >= dvlog.LogTrace {
					log.Printf("ids: %v, leftList: %v", ids, leftList)
				}
			}
			if err != nil {
				if logPreExecuteLevel >= dvlog.LogError {
					log.Printf("For table %s failed to find existing items: %v\n", table, err)
				}
				return err
			}
			if !isSingleRequests && (len(meta.Dependencies) == 0 || len(meta.IdColumns) != 1) {
				if logPreExecuteLevel >= dvlog.LogDetail {
					log.Printf("Single blocks saved: %s", GetMetaInfo(meta))
				}
				err = SavePortionOfItems(items, meta.Name, conn, left, meta.IdColumns, options, meta.Types)
				if err != nil {
					if logPreExecuteLevel >= dvlog.LogError {
						log.Printf("For table %s failed to write existing items in single block: %v\n", table, err)
					}
					return err
				}
			} else {
				var groups [][][]string
				if isSingleRequests {
					groups = makeSingleRequests(items)
				} else {
					groups, err = OrderObjectsByHierarchy(items, left, meta.IdColumns[0], meta.Dependencies)
				}
				if err != nil {
					if logPreExecuteLevel >= dvlog.LogError {
						log.Printf("For table %s failed to divide items into groups: %v\n", table, err)
					}
					return err
				}
				n = len(groups)
				if logPreExecuteLevel >= dvlog.LogDetail {
					log.Printf("%n blocks saved", n, GetMetaInfo(meta))
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
		} else if logPreExecuteLevel >= dvlog.LogDetail {
			log.Printf("(No items for %s)", table)
		}
	}
	return nil
}

func makeSingleRequests(items [][]string) [][][]string {
	n := len(items)
	groups := make([][][]string, n)
	for i := 0; i < n; i++ {
		groups[i] = [][]string{items[i]}
	}
	return groups
}
