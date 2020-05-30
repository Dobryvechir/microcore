// package dvdbdata provides functions for sql query
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbdata

import (
	"database/sql"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"strings"
	"sync"
)

var MaxConnections = 20

type connectionGraph struct {
	perDB map[string]*connectionPool
	mux   sync.Mutex
}

type connectionPool struct {
	db     []*sql.DB
	amount int
	kind   string
	mux    sync.Mutex
}

var graph = &connectionGraph{perDB: make(map[string]*connectionPool)}

func GetDBConnectionDirect(props map[string]string, connName string) (*sql.DB, string, error) {
	param := "DB_CONNECTION_" + connName
	connParams := props[param]
	if connParams == "" {
		return nil, "", errors.New(param + " expected as the definition for connection " + connName)
	}
	sqlName := ""
	conn := ""
	p := strings.Index(connParams, ",")
	if p > 0 {
		sqlName = strings.TrimSpace(connParams[:p])
		conn = strings.TrimSpace(connParams[p+1:])
	}
	if sqlName == "" || conn == "" {
		return nil, "", errors.New(param + " must be <sql name>,<connection string>")
	}
	driverName := sqlName
	p = strings.Index(sqlName, "(")
	if p > 0 {
		driverName = strings.TrimSpace(sqlName[p+1:])
		sqlName = strings.TrimSpace(sqlName[:p])
		p = strings.Index(driverName, ")")
		if p > 0 {
			driverName = strings.TrimSpace(driverName[:p])
		}
		if driverName == "" || sqlName == "" {
			return nil, "", errors.New(connName + " must contain non empty driver name in brackets")
		}
	}
	db, err := sql.Open(driverName, conn)
	if err != nil {
		return nil, "", errors.New(err.Error() + " for " + connName + "(" + connParams + ")")
	}
	props[propertyDefaultKind+connName] = sqlName
	return db, sqlName, nil
}

func GetDefaultDbConnection() string {
	return dvparser.GlobalProperties[propertyDefaultDb]
}

func GetConnectionType(connName string) int {
	sqlName := dvparser.GlobalProperties[connName]
	p := strings.Index(sqlName, ",")
	if p > 0 {
		sqlName = strings.TrimSpace(sqlName[:p])
	}
	p = strings.Index(sqlName, "(")
	if p > 0 {
		sqlName = strings.TrimSpace(sqlName[:p])
	}
	return GetConnectionKindMask(sqlName)
}

func GetConnectionKindMask(kind string) int {
	r := 0
	switch kind {
	case "oracle":
		r = SqlOracleLike
		break
	case "postgres":
		r = SqlPostgresLike
	}
	return r
}

func GetDBConnection(connName string) (r *DBConnection, err error) {
	var db *sql.DB
	var kind string
	if connName == "" {
		connName = GetDefaultDbConnection()
		if connName == "" {
			return nil, errors.New("Connection not specified")
		}
	}
	pool := graph.perDB[connName]
	if pool != nil && pool.amount != 0 {
		pool.mux.Lock()
		if pool.amount != 0 {
			pool.amount--
			db = pool.db[pool.amount]
			kind = pool.kind
		}
		pool.mux.Unlock()
	}
	if db == nil {
		db, kind, err = GetDBConnectionDirect(dvparser.GlobalProperties, connName)
	}
	kind = strings.ToLower(kind)
	r = &DBConnection{Db: db, Kind: kind, Name: connName, KindMask: GetConnectionKindMask(kind)}
	return
}

func backToPool(db *DBConnection) {
	pool := graph.perDB[db.Name]
	if pool == nil {
		graph.mux.Lock()
		pool = graph.perDB[db.Name]
		if pool == nil {
			pool = &connectionPool{db: make([]*sql.DB, MaxConnections), kind: db.Kind}
			graph.perDB[db.Name] = pool
		}
		graph.mux.Unlock()
	}
	if pool.amount < MaxConnections {
		pool.mux.Lock()
		if pool.amount < MaxConnections {
			pool.db[pool.amount] = db.Db
			pool.amount++
			db.Db = nil
		}
		pool.mux.Unlock()
		if db.Db != nil {
			err := db.Db.Close()
			if logPreExecuteLevel >= dvlog.LogDebug {
				log.Printf("Error closing db connection: %v", err)
			}
		}
	} else {
		err := db.Db.Close()
		if logPreExecuteLevel >= dvlog.LogDebug {
			log.Printf("Error closing db connection: %v", err)
		}
	}
}

func (db *DBConnection) Close(forced bool) (err error) {
	if db != nil && db.Db != nil {
		if forced {
			db.Db.Close()
		} else {
			backToPool(db)
		}
		db.Db = nil
	}
	return err
}

func (db *DBConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Db.Query(query, args...)
}

func (db *DBConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.Db.Exec(query, args...)
}
