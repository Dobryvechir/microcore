// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"encoding/json"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	emptyArray  = "[]"
	emptyObject = "{}"
)

type genTable interface {
	//      ReadAll()
	//      ReadSingle()
	//      UpdateSingle()
	//      DeleteBatch()
}

var tableMap map[string]*genTable

type fileTable struct {
	path string
}

type folderTable struct {
	path string
}

type fileWebTable struct {
	path    string
	webPath string
}
