// package dvdbdata provides functions for sql query
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbdata

import (
	"database/sql"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
)

const (
	propertyDefaultKind            = "DB_KIND_"
	propertyDefaultDb              = "DB_CONNECTIONS_DEFAULT"
	propertyTableDefinitionName    = "DB_LOCAL_STORAGE_TABLE"
	propertyTableDefinitionDefault = "R_STORAGE_LOCAL(id varchar(255) primary,name varchar(4000))"
)

const (
	SqlOracleLike      = 1
	SqlPostgresLike    = 2
	CommonMaxBatch     = 1000
	ComplexIdSeparator = "_._"
)

const (
	TypeDate   = "Date"
	TypeInt    = "int"
	TypeInt64  = "int64"
	TypeString = "string"
	TypeBool   = "bool"
)

type TableMetaData struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Dependencies []int    `json:"dependencies"`
	IdColumns    []int    `json:"idColumns"`
	MajorColumn  int      `json:"majorColumn"`
	Types        []string `json:"types"`
	Columns      []string `json:"columns"`
	References   []string `json:"references"`
	QuoteColumns bool     `json:"quoteColumns"`
}

type DBConnection struct {
	Db       *sql.DB
	KindMask int
	Kind     string
	Name     string
}

const (
	IdsPlaceholderStart  = "${"
	IdsPlaceholderFinish = "IDS}"
)

var logPreExecuteLevel = dvlog.LogFatal
var NullStringAsBytes = []byte("NULL")
