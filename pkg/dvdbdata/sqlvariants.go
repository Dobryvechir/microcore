/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

func GetDateNowFunction(sqlType int) string {
	if (sqlType & SQL_ORACLE_LIKE) != 0 {
		return "CURRENT_TIMESTAMP"
	}
	return "NOW()"
}

func GetTimestampLessDay(sqlType int) string {
	if (sqlType & SQL_ORACLE_LIKE) != 0 {
		return "CURRENT_TIMESTAMP - 1"
	}
	return "NOW() - INTERVAL '1 DAY'"
}
