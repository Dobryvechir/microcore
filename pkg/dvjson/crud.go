/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var configPathGeneral, dataPathGeneral string

type DvCrudInfo struct {
	Entities []DvCrud `json:"entities"`
}

type DvCrud struct {
	Url      string `json:"url"`
	IdName   string `json:"idName"`
	IdType   string `json:"idType"`
	Options  string `json:"options"`
	IdPrefix string `json:"idPrefix"`
}

var LogCrud, LogJson, LogProcessors bool

func crudSaveTable(crudInfo *DvCrudDetails) string {
	crudInfo.table.Kind = FIELD_ARRAY
	data := jsonWholeItems(crudInfo.table)
	err := ioutil.WriteFile(crudInfo.storeFile, []byte(data), os.ModePerm)
	if err != nil {
		log.Printf("Error writing %s: %s", crudInfo.storeFile, err.Error())
		return err.Error()
	}
	return ""
}

func crudGetAll(crudInfo *DvCrudDetails) ([]byte, int) {
	b, err := ioutil.ReadFile(crudInfo.storeFile)
	if err != nil {
		message := err.Error()
		info := strings.Index(message, "system cannot find") > 0
		if (dvlog.CurrentLogLevel >= dvlog.LogError && !info) || (LogCrud && dvlog.CurrentLogLevel >= dvlog.LogDetail) {
			log.Printf("Error reading %s: %s", crudInfo.storeFile, message)
		}
		return []byte("[]"), 0
	}
	isEmpty := true
	for _, d := range b {
		if d > 32 {
			isEmpty = false
			break
		}
	}
	if isEmpty {
		return []byte{'[', ']'}, 0
	}
	return b, 0
}

func crudGetAllTableData(crudInfo *DvCrudDetails) {
	body, statusCode := crudGetAll(crudInfo)
	if statusCode >= 300 {
		log.Printf("Crud get all failed: %d", statusCode)
	} else if LogCrud && dvlog.CurrentLogLevel >= dvlog.LogDebug {
		log.Print("info: crudGetAllTableData read " + crudInfo.storeFile + " result: " + string(body))
	}
	data := jsonConvertToItems(body, crudInfo, OPTIONS_QUICK)
	if data.Err != "" {
		log.Printf("Error in table! %s", data.Err)
	}
	crudInfo.table = data
}

func crudGetAllTableDataWithIdStatistics(crudInfo *DvCrudDetails) {
	crudGetAllTableData(crudInfo)
	makeIdStatistics(crudInfo)
}

func crudGetSingleItem(crudInfo *DvCrudDetails, id string) (*DvCrudItem, int) {
	crudGetAllTableData(crudInfo)
	items := crudInfo.table.Items
	l := len(items)
	idData := []byte(id)
	for i := 0; i < l; i++ {
		if bytes.Equal(items[i].Id, idData) {
			return items[i], 0
		}
	}
	return nil, 404 //Not Found
}

func crudGetSingle(crudInfo *DvCrudDetails, id string) ([]byte, int) {
	dat, statusCode := crudGetSingleItem(crudInfo, id)
	if statusCode >= 300 {
		return []byte{}, 404 //Not Found
	}
	return dat.itemBody, 0
}

func crudCreate(crudInfo *DvCrudDetails, body []byte) ([]byte, int) {
	data := jsonConvertToItems(body, crudInfo, OPTIONS_QUICK)
	if data.Err != "" {
		return dvlog.FormErrorMessage(data.Err), 400 //Bad Request
	}
	l := len(data.Items)
	if l == 0 {
		return []byte{}, 0
	}
	if LogCrud && dvlog.CurrentLogLevel >= dvlog.LogTrace {
		log.Printf("CrudBefore %s", logInfoForItem(data.Items[0]))
	}
	crudGetAllTableDataWithIdStatistics(crudInfo)
	if LogCrud && dvlog.CurrentLogLevel >= dvlog.LogTrace {
		log.Printf("CrudAfter %s", logInfoForItem(data.Items[0]))
	}
	for i := 0; i < l; i++ {
		itemId := data.Items[i].Id
		id := string(itemId)
		if crudInfo.userIdPossible && len(itemId) > 0 {
			if _, ok := crudInfo.idMap[id]; ok {
				return dvlog.FormErrorMessage("Id " + id + " already exists"), 406 //Not Acceptable
			}
		} else {
			id = crudInfo.autoIdPrefix + strconv.Itoa(crudInfo.nextId)
			crudInfo.nextId++
			if LogCrud && dvlog.CurrentLogLevel >= dvlog.LogTrace {
				log.Printf("CrudChange %s %s", id, logInfoForItem(data.Items[i]))
			}
			changeItemId(data.Items[i], []byte(id), crudInfo)
			if LogCrud && dvlog.CurrentLogLevel >= dvlog.LogTrace {
				log.Printf("CrudChanged %s %s", id, logInfoForItem(data.Items[i]))
			}
		}
	}
	crudInfo.table.Items = append(crudInfo.table.Items, data.Items...)
	err := crudSaveTable(crudInfo)
	if err != "" {
		return dvlog.FormErrorMessage(err), 500 //Internal Server Error"
	}
	return jsonWholeItems(data), 0
}

func crudUpdate(crudInfo *DvCrudDetails, body []byte) ([]byte, int) {
	data := jsonConvertToItems(body, crudInfo, OPTIONS_QUICK)
	if data.Err != "" {
		return dvlog.FormErrorMessage(data.Err), 400 //Bad Request
	}
	l := len(data.Items)
	if l == 0 {
		return []byte{}, 0
	}
	crudGetAllTableDataWithIdStatistics(crudInfo)
	changed := false
	for _, item := range data.Items {
		id := string(item.Id)
		if n, ok := crudInfo.idMap[id]; ok {
			if !bytes.Equal(crudInfo.table.Items[n].itemBody, item.itemBody) {
				changed = true
				crudInfo.table.Items[n] = item
			}
		} else if LogCrud && dvlog.CurrentLogLevel >= dvlog.LogInfo {
			log.Printf("Update element is not present for id=%s data=%v", item.Id, item)
		}
	}
	if changed {
		err := crudSaveTable(crudInfo)
		if err != "" {
			return dvlog.FormErrorMessage(err), 500 //Internal Server Error
		}
	}
	return jsonWholeItems(data), 0
}

func crudDeleteByArray(crudInfo *DvCrudDetails, ids []string) ([]byte, int) {
	crudGetAllTableData(crudInfo)
	items := crudInfo.table.Items
	l := len(items)
	n := len(ids)
	if n == 0 || l == 0 {
		return []byte{}, 0
	}
	m := make(map[string]bool)
	for i := 0; i < n; i++ {
		m[ids[i]] = true
	}
	k := 0
	for i := 0; i < l; i++ {
		if !m[string(items[i].Id)] {
			if k != i {
				items[k] = items[i]
			}
			k++
		}
	}
	if k != l {
		crudInfo.table.Items = items[:k]
		err := crudSaveTable(crudInfo)
		if err != "" {
			return dvlog.FormErrorMessage(err), 500 //Internal Server Error
		}
	}
	return []byte{}, 0
}

func crudDeleteByStrings(crudInfo *DvCrudDetails, ids string) ([]byte, int) {
	ids = strings.TrimSpace(ids)
	if ids == "" {
		return []byte{}, 0
	}
	idArray := strings.Split(ids, ",")
	return crudDeleteByArray(crudInfo, idArray)
}

func MakeCrudHandler(crudInfo *DvCrudDetails) dvmeta.HandlerFunc {
	return func(request *dvmeta.RequestContext) bool {
		urls := strings.Split(request.Url, "/")
		urlsLen := len(urls)
		dataUrl := ""
		if crudInfo.dataUrlIndex < urlsLen {
			dataUrl = urls[crudInfo.dataUrlIndex]
		}
		mes := []byte{'[', ']'}
		stat := 403 //method forbidden
		var logFile string
		body, err := ioutil.ReadAll(request.Reader.Body)
		if err != nil {
			message := err.Error()
			log.Print(message)
			mes = dvlog.FormErrorMessage(message)
			stat = 400 //bad request
		} else {
			if LogCrud {
				if dvlog.CurrentLogLevel >= dvlog.LogInfo && (request.Reader.Method != "OPTIONS" || dvlog.CurrentLogLevel >= dvlog.LogDetail) {
					logFile = dvlog.WriteRequestToLog(body, request.Reader)
				}
			}
			switch request.Reader.Method {
			case "OPTIONS":
				stat = 200 //ok
			case "GET":
				if dataUrl == "" {
					mes, stat = crudGetAll(crudInfo)
				} else {
					mes, stat = crudGetSingle(crudInfo, dataUrl)
				}
			case "POST":
				mes, stat = crudCreate(crudInfo, body)
			case "PUT":
				mes, stat = crudUpdate(crudInfo, body)
			case "DELETE":
				mes, stat = crudDeleteByStrings(crudInfo, dataUrl)
			}
		}
		if logFile != "" {
			fileName := logFile + "S." + strconv.Itoa(stat)
			ioutil.WriteFile(fileName, mes, os.ModePerm)
		}
		if stat == 403 {
			return false
		}
		request.Output = mes
		request.StatusCode = stat
		dvcom.HandleRequestContext(request)
		return true
	}
}

func crudProcessEntities(baseUrl string, crudInfo *DvCrudInfo, urlPool *dvurl.UrlPool) error {
	crudCheckGlobalPaths()
	crudFolder := dataPathGeneral
	for _, crud := range crudInfo.Entities {
		url := crud.Url
		if len(url) > 0 {
			c := url[0]
			if c == '/' {
				url = baseUrl + url
			} else {
				url = baseUrl + "/" + url
			}
		} else {
			url = baseUrl
		}
		pos := strings.Index(url, "{")
		if pos > 0 {
			url = url[:pos]
		}
		if crud.IdName == "" {
			crud.IdName = "id"
		}
		if crud.IdType == "" {
			crud.IdType = "string"
		}
		storePath := dvlog.GetSubFolder(crudFolder, url)
		storeFile := storePath + "/data.json"
		options := make(map[string]string)
		optionData := crud.Options
		optionDataLen := len(optionData)
		for i := 0; i < optionDataLen; {
			key := optionData[i : i+1]
			i++
			k := i
			for k < optionDataLen && (optionData[k] < 'A' || optionData[k] > 'Z') {
				k++
			}
			val := optionData[i:k]
			if val == "" {
				val = "true"
			}
			options[key] = val
		}
		urlParts := strings.Split(url, "/")
		urlPartsLen := len(urlParts)
		if urlParts[urlPartsLen-1] == "" {
			urlPartsLen--
			urlParts = urlParts[:urlPartsLen]
		}
		crudInfoSelf := DvCrudDetails{storeFile: storeFile, initCrud: crud, userIdPossible: options["C"] == "true", lookedUrl: url,
			urlParts: urlParts, dataUrlIndex: urlPartsLen,
			userIdInteger: crud.IdType != "string", userIdName: []byte(crud.IdName), autoIdPrefix: crud.IdPrefix}
		handler := MakeCrudHandler(&crudInfoSelf)
		urlPool.RegisterHandlerFunc(url, handler)
	}
	return nil
}

func crudGlobalInitialization(data map[string]string) error {
	for k, v := range data {
		switch k {
		case "configPath":
			configPathGeneral = v
		case "dataPath":
			dataPathGeneral = v
		default:
			return errors.New("Unknown crud global initialization parameter: " + k)
		}
	}
	return nil
}

func crudCheckGlobalPaths() {
	configPathGeneral = dvlog.EnsureSubFolderCritical(configPathGeneral, "CRUD")
	dataPathGeneral = dvlog.EnsureSubFolderCritical(dataPathGeneral, "CRUD")
}

func crudGenerateHandlers(url string, params []string, urlPool *dvurl.UrlPool) error {
	if len(params) == 0 || params[0] == "" {
		return errors.New("Crud config name is not specified")
	}
	fileName, err := dvlog.FindConfigFile(configPathGeneral, params[0])
	if err != nil {
		return err
	}
	for url != "" {
		c := url[len(url)-1]
		if c == '/' || c == '\\' || c == ':' {
			url = url[:len(url)-1]
		} else {
			break
		}
	}
	data, err2 := dvparser.SmartReadTemplate(fileName, 3, byte(' '))
	if err2 != nil {
		return err2
	}
	dvlog.CleanEOL(data)
	cf := &DvCrudInfo{}
	err = json.Unmarshal(data, cf)
	if err != nil {
		return err
	}
	return crudProcessEntities(url, cf, urlPool)
}

var crudRegistrationConfig *dvmodules.RegistrationConfig = &dvmodules.RegistrationConfig{
	Name:              "crud",
	GlobalInitHandler: crudGlobalInitialization,
	GenerateHandlers:  crudGenerateHandlers,
}

var crudInited bool = dvmodules.RegisterModule(crudRegistrationConfig, false)
