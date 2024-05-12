// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"sync"
)

const (
	emptyArray               = "[]"
	emptyObject              = "{}"
	defaultKeyFirst          = "id"
	defaultWebField          = "file"
	defaultWebFileName       = "fileName"
	defaultWebAllowedFormats = "i"
)

type genTable interface {
	ReadAll() interface{}
	ReadOne(key interface{}) interface{}
	ReadFieldsForIds(ids []*dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error)
	ReadFieldsForId(id *dvevaluation.DvVariable, fields []string) (*dvevaluation.DvVariable, error)
	ReadFieldsForAll(fields []string) (*dvevaluation.DvVariable, error)
	DeleteKeys(keys []string) interface{}
	CreateRecord(record *dvevaluation.DvVariable, newId string) (*dvevaluation.DvVariable, error)
	UpdateRecord(record *dvevaluation.DvVariable) (*dvevaluation.DvVariable, error)
}

var tableMap map[string]genTable

type fileTable struct {
	mu            sync.Mutex
	path          string
	keyFirst      string
	allowCustomId bool
	version       string
}

type folderTable struct {
	mu            sync.Mutex
	path          string
	keyFirst      string
	allowCustomId bool
	version       string
}

type fileWebTable struct {
	mu                sync.Mutex
	path              string
	allowCustomId     bool
	version           string
	webUrl            string
	webPath           string
	keyFirst          string
	webField          string
	webFileName       string
	webAllowedFormats string
}
