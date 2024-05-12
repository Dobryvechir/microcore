// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

func deleteWebFiles(webPath string, keys []string) {
	keyMap := convertFieldsToMap(keys)
	entries, err := os.ReadDir(webPath)
	if err != nil {
		return
	}
	for _, e := range entries {
		name := e.Name()
		if e.Type().IsRegular() {
			p := strings.Index(name, ".")
			if p > 0 {
				s := name[:p]
				_, ok := keyMap[s]
				if ok {
					nm := webPath + "/" + name
					os.Remove(nm)
				}
			}
		}
	}
}

func cleanWebFiles(name string) {
	if len(name) > 0 {
		os.Remove(name)
	}
}

func findResourceForId(webPath string, id string) string {
	entries, err := os.ReadDir(webPath)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		name := e.Name()
		if e.Type().IsRegular() {
			p := strings.Index(name, ".")
			if p > 0 {
				s := name[:p]
				if s == id {
					return name
				}
			}
		}
	}
	return ""
}

func updateWebFiles(webPath string, record *dvevaluation.DvVariable, id string, webField string, webFileName string, webUrl string, formats string) (string, error) {
	if !checkIntId(id) {
		return "", errors.New("Wrong id " + id)
	}
	userFile, ok := readFieldInJsonAsString(record, webField)
	oldResource := findResourceForId(webPath, id)
	if len(oldResource) > 0 {
		if ok && userFile == webUrl+"/"+oldResource {
			return "", nil
		}
		cleanWebFiles(webPath + "/" + oldResource)
	}
	if !ok || len(userFile) == 0 {
		return "", nil
	}
	extension, position, prefix, transform := analyzeUserFile(userFile)
	if len(prefix) == 0 {
		return "", nil
	}
	if !strings.Contains(formats, prefix) {
		return "", errors.New("Unsupported format")
	}
	saveName := webPath + "/" + id + extension
	urlName := webUrl + "/" + id + extension
	fileSize, err := saveUserFile(saveName, position, transform, userFile)
	if err != nil {
		return "", err
	}
	setFieldInJsonAsString(record, webField, urlName)
	sum := dvtextutils.CalculateStringCheckSum(userFile)
	if sum < 0 {
		sum = -sum
	}
	uniqueName := prefix + id + "_" + strconv.FormatInt(sum, 10) + "-" + strconv.FormatInt(fileSize, 10) + extension
	setFieldInJsonAsString(record, webFileName, uniqueName)
	return saveName, nil
}
