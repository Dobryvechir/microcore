/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

var tempPath string
var tmpPrefix = "dv_"

func SetTempPath(path string, prefix string) {
	tempPath = path
	tmpPrefix = prefix
}

func GetTempPathSlashed() string {
	path := GetTempPath()
	if path == "" {
		return ""
	}
	if path[len(path)-1] != '/' && path[len(path)-1] != '\\' {
		path += "/"
	}
	return path
}

func GetTemporaryFileName() string {
	path := GetTempPathSlashed()
	if path == "" {
		return ""
	}
	for i := 0; i < 2000000; i++ {
		fileName := path + tmpPrefix + strconv.Itoa(i)
		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			return fileName
		}
	}
	log.Print("Your temporary folder is not accessible, please define a good temporary folder in TEMP environment variable")
	return ""
}

func GetUniqueFileName(dir string, prefix string, suffix string) string {
	needSlash := dir != "" && dir[len(dir)-1] == '/' || dir[len(dir)-1] == '\\'
	if needSlash {
		if prefix != "" && (prefix[0] == '/' || prefix[0] == '\\') {
			dir += prefix
		} else {
			dir += "/" + prefix
		}
	} else {
		if prefix != "" && (prefix[0] == '/' || prefix[0] == '\\') {
			dir += prefix[1:]
		} else {
			dir += prefix
		}
	}
	inter := ""
	for i := 1; i < 2000000; i++ {
		fileName := dir + inter + suffix
		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			return fileName
		}
		inter = strconv.Itoa(i)
	}
	return ""
}

func GetSafeFileName(src string) string {
	data := []byte(src)
	n := len(data)
	for i := 0; i < n; i++ {
		c := data[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-') {
			data[i] = '_'
		}
	}
	return string(data)
}

func GetTempPath() string {
	if tempPath != "" {
		return tempPath
	}
	tempPath = os.Getenv("TEMP")
	if tempPath != "" {
		if _, err := os.Stat(tempPath); err == nil {
			return tempPath
		}
	}
	tempPath = os.Getenv("TMP")
	if tempPath != "" {
		if _, err := os.Stat(tempPath); err == nil {
			return tempPath
		}
	}
	tempPath = "/tmp"
	if _, err := os.Stat(tempPath); err == nil {
		return tempPath
	}
	tempPath = "/temp"
	if _, err := os.Stat(tempPath); err == nil {
		return tempPath
	}
	log.Print("temporary folder is not available (define it in TEMP environment variable)")
	return ""
}

func CleanTmpFolder(tmpFolder string) bool {
	if tmpFolder == "" || tmpFolder == "/" {
		return true
	}
	os.RemoveAll(tmpFolder + "/")
	return true
}

func GetUniqueTmpFolder() string {
	dirName := GetTemporaryFileName()
	if dirName == "" {
		return dirName
	}
	os.MkdirAll(dirName+"/", 0777)
	return dirName
}

func GetFirstExistingPath(paths ...string) string {
	n := len(paths)
	for i := 0; i < n; i++ {
		s := paths[i]
		if _, err := os.Stat(s); err == nil {
			return s
		}
	}
	return ""
}

func GetFolderNameWithoutLastSlash(name string) string {
	if name == "" {
		return name
	}
	n := len(name) - 1
	if name[n] == '/' {
		name = name[:n]
	}
	return name
}

func GetFolderLastPartName(name string) string {
	name = GetFolderNameWithoutLastSlash(name)
	pos := strings.LastIndex(name, "/")
	if pos < 0 {
		return name
	}
	return name[pos+1:]
}

func GetFolderNameWithoutLastPart(name string) string {
	name = GetFolderNameWithoutLastSlash(name)
	pos := strings.LastIndex(name, "/")
	if pos < 0 {
		return ""
	}
	return name[:pos]
}

func SaveToUniqueFile(data string) string {
	fileName := GetTemporaryFileName()
	if fileName == "" {
		return ""
	}
	err := ioutil.WriteFile(fileName, []byte(data), 0664)
	if err != nil {
		log.Printf("Failed to save %s: %v", data, err)
	}
	return fileName
}
