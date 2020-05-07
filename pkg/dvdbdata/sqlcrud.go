/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import "database/sql"

func SqlSingleValueByConnectionName(connName string, query string) (string, bool, error) {
	db, _, err := GetDB(connName)
	if err != nil {
		return "", false, err
	}
	return SqlSingleValueByConnection(db, query)
}

func SqlSingleValueByConnection(db *sql.DB, query string) (string, bool, error) {
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
	db, _, err := GetDB(connName)
	if err != nil {
		return err
	}
	return SqlUpdateByConnection(db, query)
}

func SqlUpdateByConnection(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}
