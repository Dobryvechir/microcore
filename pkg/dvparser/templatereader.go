/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvparser

import (
	"bytes"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"log"
	"os"
)

var GlobalProperties map[string]string
var GlobalPropertiesAsDvObject *dvevaluation.DvObject
var GeneralFilePaths []string
var GeneralNamespaceFolderIndex int = -1
var GeneralMicroCoreFolderIndex int = -1

func LogVariables(name string, vars map[string]string) {
	log.Print("START OF " + name)
	for k, v := range vars {
		log.Print(k + "=" + v)
	}
	log.Print("END OF " + name)
}

func FindInGeneralPaths(name string) string {
	if _, err := os.Stat(name); err == nil {
		return name
	}
	for _, v := range GeneralFilePaths {
		filename := v + "/" + name
		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}
	return ""
}

func SmartReadTemplateLines(fileName string, numberOfBrackets int, properties map[string]string) (lines [][]byte, err error, lastOutputMap map[string]string) {
	lastOutputMap = make(map[string]string, 0)
	configInfo := ConfigInfo{InputMap: properties,
		OutputMap:        lastOutputMap,
		NumberOfBrackets: numberOfBrackets,
		Options:          CONFIG_IS_NOT_VARIABLES,
		FilePaths:        GeneralFilePaths}
	LinearSmartConfigFromFile(fileName, &configInfo, -1)
	err = configInfo.Err
	lines = configInfo.OutputLines
	return
}

func SmartReadFileAsString(fileName string, properties map[string]string) (string, error) {
	res, err, _ := SmartReadTemplateLines(fileName, 3, properties)
	if err != nil {
		return "", err
	}
	return string(bytes.Join(res, []byte{13, 10})), nil
}

func SmartReadTemplate(filename string, numberOfBrackets int, joiner byte) ([]byte, error) {
	res, err, _ := SmartReadTemplateLines(filename, numberOfBrackets, GlobalProperties)
	if err != nil {
		return nil, err
	}
	return bytes.Join(res, []byte{joiner}), nil
}

func CloneGlobalProperties() map[string]string {
	newMap := make(map[string]string)
	for k, v := range GlobalProperties {
		newMap[k] = v
	}
	return newMap
}

func SetGlobalPropertiesValue(key string, value string) {
	GlobalProperties[key] = value
}

func RemoveGlobalPropertiesValue(key string) {
	delete(GlobalProperties, key)
}

func GetGlobalPropertiesAsDvObject() *dvevaluation.DvObject {
	if GlobalPropertiesAsDvObject == nil {
		GlobalPropertiesAsDvObject = dvevaluation.NewDvObjectWithGlobalPrototype(GlobalProperties)
	}
	return GlobalPropertiesAsDvObject
}

func GetProptertiesPrototypedToGlobalProperties(localMap map[string]interface{}) *dvevaluation.DvObject {
	return dvevaluation.NewObjectWithPrototype(localMap, GetGlobalPropertiesAsDvObject())
}

func FindEol(data []byte) int {
	n := len(data)
	for i := 0; i < n; i++ {
		if data[i] == 13 || data[i] == 10 {
			return i
		}
	}
	return -1
}
