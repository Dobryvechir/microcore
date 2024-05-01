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

func folderKindInit(tbl *dvcontext.DatabaseTable, db *DatabaseConfig) *genTable {
	path := db.Root + "/" + tbl.Name
	os.MkdirAll(path, 0755)
	ref := &folderTable{path}
	return ref
}
