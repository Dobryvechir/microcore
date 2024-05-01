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

func fileWebKindInit(tbl *dvcontext.DatabaseTable, db *DatabaseConfig) *genTable {
	path := db.Root + "/" + tbl.Name + ".json"
	if _, err := os.Stat(path); err != nil && errors.Is(err, os.ErrNotExist) {
		ioutil.WriteFile(path, emptyArray, 0644)
	}
	webPath := db.WebRoot + tbl.web
	os.MkdirAll(webPath, 0755)
	ref := &fileWebTable{path, webPath}
	return ref

}
