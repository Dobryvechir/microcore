/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdbdata

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func SplitSqlSequences(data []byte) []string {
	n := bytes.Count(data, []byte(";")) + 1
	res := make([]string, 0, n)
	n = len(data)
	p := 0
	for i := 0; i < n; i++ {
		b := data[i]
		if b == '\'' {
			for i++; i < n; i++ {
				b = data[i]
				if b == '\\' {
					i++
				} else if b == '\'' {
					break
				}
			}
		} else if b == ';' {
			for p < i && data[p] <= ' ' {
				p++
			}
			e := i - 1
			for e >= p && data[e] <= ' ' {
				e--
			}
			if e >= p {
				res = append(res, string(data[p:e+1]))
			}
			p = i + 1
		} else if b < ' ' {
			data[i] = ' '
		}
	}
	for p < n && data[p] <= ' ' {
		p++
	}
	e := n - 1
	for e >= p && data[e] <= ' ' {
		e--
	}
	if e >= p {
		res = append(res, string(data[p:e+1]))
	}
	return res
}

func ExecuteSqlData(db *DBConnection, data []byte) error {
	queries := SplitSqlSequences(data)
	n := len(queries)
	if logPreExecuteLevel >= dvlog.LogDetail {
		log.Printf(" contains %d sql requests\n", n)
	}
	for i := 0; i < n; i++ {
		query := queries[i]
		if logPreExecuteLevel >= dvlog.LogDetail {
			log.Printf("%d) %s\n", i+1, query)
		}
		_, err := db.Exec(query)
		if err != nil {
			if logPreExecuteLevel >= dvlog.LogError {
				log.Printf("Failed [%s] of [%v] because of [%v]", query, queries, err)
			}
			return errors.New(err.Error() + " in " + query)
		}
	}
	return nil
}

func ExecuteSqlFromFile(db *DBConnection, fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Failed to read the content of %s: %v\n", fileName, err)
		}
		return err
	}
	if logPreExecuteLevel >= dvlog.LogInfo {
		log.Printf("Executing sql from %s of %d bytes", fileName, len(data))
	}
	err = ExecuteSqlData(db, data)
	if err != nil {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Failed to execute the content of %s: %v\n", fileName, err)
		}
		return errors.New(err.Error() + " :" + fileName)
	}
	return nil
}

func ExecuteSqlFromFolder(db *DBConnection, root string) error {
	var files []string
	mask := strings.ToLower(db.Kind)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if logPreExecuteLevel >= dvlog.LogError {
				log.Printf("error %s: %v", path, err)
			}
			return nil
		} else if info.IsDir() {
			if logPreExecuteLevel >= dvlog.LogTrace {
				log.Printf("folder %s", path)
			}
			return nil
		}
		s := strings.ToLower(path)
		if strings.HasSuffix(s, ".sql") || strings.HasSuffix(s, ".csv") {
			n := len(path)
			t := s[:n-4]
			p := strings.LastIndex(t, ".")
			if p >= 0 {
				t = t[p+1:]
				if dvparser.IsAlphabeticalLowCase(t) && t != mask {
					if logPreExecuteLevel >= dvlog.LogInfo {
						log.Printf("omitting %s because it is %s, not %s\n", path, t, mask)
					}
					return nil
				}
			}
			files = append(files, path)
		} else {
			if logPreExecuteLevel >= dvlog.LogInfo {
				log.Printf("omitting %s because it is neither sql nor csv\n", path)
			}
		}
		return nil
	})
	if err != nil {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Cannot read thru folder %s: %v", root, err)
		}
		return err
	}
	sort.Strings(files)
	csvOptions := 0
	switch db.Kind {
	case "oracle":
		csvOptions |= SqlOracleLike
		break
	case "postgres":
		csvOptions |= SqlPostgresLike
		break
	}
	if logPreExecuteLevel >= dvlog.LogInfo {
		log.Printf("%d files are found for execution", len(files))
	}
	for _, file := range files {
		ext := strings.ToLower(file[len(file)-4:])
		switch ext {
		case ".sql":
			err = ExecuteSqlFromFile(db, file)
			break
		case ".csv":
			err = PreExecuteCsvFile(db, file, csvOptions)
			break
		}
		if err != nil {
			if logPreExecuteLevel >= dvlog.LogError {
				log.Printf("Preexecution interrupted at %s of %v because of %v ", file, files, err)
			}
			return err
		}
	}
	return nil
}

func PreExecuteForNewerVersions(props map[string]string, db *DBConnection, folder string) error {
	if _, err := os.Stat(folder); err != nil {
		if logPreExecuteLevel >= dvlog.LogInfo {
			log.Printf("Preexecution is omitted because folder %s does not exist", folder)
		}
		return nil
	}
	pos := strings.LastIndex(folder, "/")
	key := folder
	if pos >= 0 {
		key = folder[pos+1:]
	}
	key = "DB_PREEXECUTE_VERSION_" + key
	version, _ := ReadGlobalDBProperty(props, db, key, "0")
	versionIndex := dvparser.ReadHexValue(version)
	if versionIndex < 0 {
		versionIndex = 0
	}
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Preexecution interrupted because of the failure to read the %s folder content: %v", folder, err)
		}
		return err
	}
	version = dvparser.GetCanonicalVersion(versionIndex)
	if logPreExecuteLevel >= dvlog.LogInfo {
		log.Printf("Sql preexecution started in folder %s, sql - %s, version older than %s", folder, db.Kind, version)
	}
	commonSql := make([]string, 0, 3)
	versionSql := make([]string, 0, 3)
	for _, file := range files {
		if file.IsDir() {
			nm := file.Name()
			if nm != "" && (nm[0] == 'v' || nm[0] == 'V') {
				versionNmb := dvparser.GetVersionIndex(nm[1:])
				if versionNmb > versionIndex {
					versionSql = append(versionSql, dvparser.Int64ToFullHex(versionNmb)+nm)
					if logPreExecuteLevel > dvlog.LogTrace {
						log.Printf("File %s added for processing", nm)
					}
				} else if logPreExecuteLevel > dvlog.LogInfo {
					log.Printf("Folder %s omitted because its version %s is not older than %s", file.Name(), dvparser.GetCanonicalVersion(versionNmb), version)
				}
			} else if strings.HasPrefix(strings.ToLower(nm), "common") {
				commonSql = append(commonSql, nm)
				if logPreExecuteLevel > dvlog.LogTrace {
					log.Printf("File %s added for processing", nm)
				}
			} else if logPreExecuteLevel > dvlog.LogInfo {
				log.Printf("Folder %s omitted because is neither common nor version", file.Name())
			}
		} else if logPreExecuteLevel >= dvlog.LogInfo {
			log.Printf("File %s omitted because only folders are used at this level", file.Name())
		}
	}
	sort.Strings(commonSql)
	sort.Strings(versionSql)
	if logPreExecuteLevel >= dvlog.LogInfo {
		log.Printf("Found %d version folders and %d common folders for processing", len(versionSql), len(commonSql))
	}
	for _, file := range versionSql {
		name := folder + "/" + file[16:]
		if logPreExecuteLevel >= dvlog.LogWarning {
			log.Printf("Processing %s, sql - %s", name, db.Kind)
		}
		err := ExecuteSqlFromFolder(db, name)
		if err != nil {
			return err
		}
	}
	for _, file := range commonSql {
		name := folder + "/" + file
		if logPreExecuteLevel >= dvlog.LogWarning {
			log.Printf("Processing %s, sql - %s", name, db.Kind)
		}
		err := ExecuteSqlFromFolder(db, name)
		if err != nil {
			return err
		}
	}
	n := len(versionSql) - 1
	if n >= 0 {
		name := versionSql[n][:16]
		if logPreExecuteLevel >= dvlog.LogInfo {
			log.Printf("Highest version stored is %s", dvparser.GetCanonicalVersionFromHexName(name))
		}
		err := WriteGlobalDBProperty(props, db, key, name)
		if err != nil {
			if logPreExecuteLevel >= dvlog.LogError {
				log.Printf("Cannot write highest version: %v", err)
			}
			return err
		}
	}
	return nil
}

func PreExecute(properties map[string]string) error {
	inits := strings.TrimSpace(properties["DB_CONNECTIONS_INIT"])
	if inits == "" {
		if logPreExecuteLevel >= dvlog.LogInfo {
			log.Printf("Pre-execution is omitted because DB_CONNECTIONS_INIT is empty")
		}
		return nil
	}
	folder := strings.TrimSpace(properties["DB_ROOT_PREEXECUTION_FOLDER"])
	if folder == "" {
		if logPreExecuteLevel >= dvlog.LogError {
			log.Printf("Pre-execution is omitted and error fired because DB_ROOT_PREEXECUTION_FOLDER is empty")
		}
		return errors.New("Specify DB_CONNECTIONS_DIR where preexecution sqls are stored")
	}
	c := folder[len(folder)-1]
	if c != '/' && c != '\\' {
		folder += "/"
	}
	connections := dvparser.ConvertToNonEmptyList(inits)
	for _, connection := range connections {
		if logPreExecuteLevel >= dvlog.LogInfo {
			log.Printf("Preexecution started for connection %s", connection)
		}
		db, err := GetDBConnection(connection)
		if err != nil {
			if logPreExecuteLevel >= dvlog.LogError {
				log.Printf("Preexecution is interrupted because of error of opening the connection %s: %v", connection, err)
			}
			return err
		}
		if logPreExecuteLevel >= dvlog.LogTrace {
			log.Printf("Type of connection %s is %s", connection, db.Kind)
		}
		err = PreExecuteForNewerVersions(properties, db, folder+connection)
		err1 := db.Close(err!=nil)
		if err != nil {
			return err
		}
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func preExecuteStart(eventName string, data []interface{}) error {
	logPreExecuteLevel = dvlog.GetLogLevelByDefinition(dvparser.GlobalProperties["DVLOG_PREEXECUTION_LEVEL"], logPreExecuteLevel)
	if logPreExecuteLevel >= dvlog.LogTrace {
		log.Printf("Event %s with data %v fired", eventName, data)
	}
	return PreExecute(dvparser.GlobalProperties)
}

var hookConfig = &dvmodules.HookRegistrationConfig{
	Name:            "dvDbDataPreInit",
	HookEventMapper: map[string]dvmodules.HookMethodEndPointHandler{dvmodules.HookStartEvent: preExecuteStart},
}

func hookRegister() bool {
	dvmodules.SubscribeForEvents(hookConfig, true)
	return true
}

var registeredHookConfig = hookRegister()
