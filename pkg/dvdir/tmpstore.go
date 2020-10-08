/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func getSaveFileFolder(dirName string) string {
	pos := strings.Index(dirName, "/")
	if pos >= 0 {
		os.MkdirAll(dirName, 0755)
		return dirName
	}
	dir := GetTempPathSlashed() + dirName
	os.MkdirAll(dir, 0755)
	return dir
}

func getSaveSpecificFileFolder(dirName string, fileName string) string {
	dir := getSaveFileFolder(dirName) + "/" + fileName
	os.MkdirAll(dir, 0755)
	return dir + "/"
}

func GetNextSaveFileName(dirName string, fileName string) string {
	dir := getSaveSpecificFileFolder(dirName, fileName)
	for i := 0; i < 10000000; i++ {
		s := dir + strconv.Itoa(i)
		if !Exists(s) {
			return s
		}
	}
	return ""
}

func GetLastSavedFileName(dirName string, fileName string) string {
	dir := getSaveSpecificFileFolder(dirName, fileName)
	p := ""
	for i := 0; i < 10000000; i++ {
		s := dir + strconv.Itoa(i)
		if !Exists(s) {
			return p
		} else {
			p = s
		}
	}
	return ""
}

func DeleteLastSavedFile(dirName string, fileName string, count int) bool {
	if count == 0 {
		return DeleteAllSavedFiles(dirName, fileName)
	}
	dir := getSaveSpecificFileFolder(dirName, fileName)
	i := 0
	for ; i < 10000000; i++ {
		s := dir + strconv.Itoa(i)
		if !Exists(s) {
			break
		}
	}
	if count < 0 {
		if i <= 1 {
			return true
		}
		count = 1
	}
	if i <= count {
		os.RemoveAll(dir)
		return true
	}
	for i--; count > 0; count-- {
		os.Remove(dir + strconv.Itoa(i))
		i--
	}
	return true
}

func DeleteAllSavedFiles(dirName string, fileName string) bool {
	dir := getSaveSpecificFileFolder(dirName, fileName)
	os.RemoveAll(dir)
	return true
}

func DeleteAllSavedFilesExceptZero(dirName string, fileName string) bool {
	dir := getSaveSpecificFileFolder(dirName, fileName)
	data, err := ioutil.ReadFile(dir + "0")
	os.RemoveAll(dir)
	if err == nil {
		os.MkdirAll(dir, 0755)
		ioutil.WriteFile(dir+"0", data, 0664)
	}
	return true
}
