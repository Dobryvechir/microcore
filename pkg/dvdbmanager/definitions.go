// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

const (
	emptyArray      = "[]"
	emptyObject     = "{}"
	defaultKeyFirst = "id"
)

type genTable interface {
	ReadAll() interface{}
	ReadOne(key interface{}) interface{}
	//      UpdateSingle()
	//      DeleteBatch()
}

var tableMap map[string]genTable

type fileTable struct {
	path     string
	keyFirst string
}

type folderTable struct {
	path     string
	keyFirst string
}

type fileWebTable struct {
	path     string
	webPath  string
	keyFirst string
}
