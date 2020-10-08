/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type cacheEntry struct {
	value interface{}
	stamp time.Time
}

type CacheValueGetter func(key string, item string) (interface{}, error)

type cacheConfig struct {
	memoryDuration int
	fileDuration   int
	folderName     string
	valueGetter    CacheValueGetter
}

const CommonGroup = "dv_common"

var cacheEntries = make(map[string]*cacheEntry)
var cacheConfigs = make(map[string]*cacheConfig)
var cacheMainFolder string

func GetCacheLocalFileInfo(folder string, item string) (folderFile string) {
	if cacheMainFolder == "" {
		cacheMainFolder = GetTempPathSlashed() + "___dv__server__cache/"
	}
	folderFile = cacheMainFolder + folder
	os.MkdirAll(folderFile, 0755)
	if item == "" {
		return folderFile
	}
	item = GetSafeFileName(item)
	return folderFile + "/" + item
}

func GetGroupByKey(key string) (grp string, item string) {
	pos := strings.LastIndex(key, ".")
	if pos >= 0 {
		grp = key[:pos]
		item = key[pos+1:]
	} else {
		grp = CommonGroup
		item = key
	}
	if grp == "" {
		grp = CommonGroup
	}
	if item == "" {
		item = "________"
	}
	return
}

func getConfigItemFileNameByKey(key string) (item string, config *cacheConfig, fileName string) {
	grp, item := GetGroupByKey(key)
	config = cacheConfigs[grp]
	if config != nil && config.fileDuration > 0 {
		if config.folderName == "" {
			config.folderName = GetSafeFileName(grp)
		}
		fileName = GetCacheLocalFileInfo(config.folderName, item)
	}
	return
}

var errorCacheNotFound = errors.New("Not found in cache")

func SetCacheConfigWhole(groupName string, memoryDuration int, fileDuration int, folderName string, valueGetter CacheValueGetter) bool {
	cacheConfigs[groupName] = &cacheConfig{
		memoryDuration: memoryDuration,
		fileDuration:   fileDuration,
		folderName:     folderName,
		valueGetter:    valueGetter,
	}
	return true
}

func ReadConfigByDescription(config string) (duration1 int, duration2 int, name string, ok bool) {
	data := strings.Split(config, ",")
	duration1 = 3600
	n := len(data)
	ok = true
	if n > 0 {
		val := strings.TrimSpace(data[0])
		if val != "" {
			m, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Config duration time must be integer, but it is " + val)
				ok = false
			} else {
				duration1 = m
			}
		}

	}
	if n > 1 {
		val := strings.TrimSpace(data[1])
		if val != "" {
			m, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Config duration time must be integer, but it is " + val)
				ok = false
			} else {
				duration2 = m
			}
		}

	}
	if n > 2 {
		val := strings.TrimSpace(data[2])
		if val != "" {
			name = val
		}
	}
	return
}

func SetCacheConfigByDescription(groupName string, config string, valueGetter CacheValueGetter) bool {
	memoryDuration, fileDuration, folderName, ok := ReadConfigByDescription(config)
	if !ok {
		return false
	}
	return SetCacheConfigWhole(groupName, memoryDuration, fileDuration, folderName, valueGetter)
}

func GetCacheValue(key string) (interface{}, error) {
	entry := cacheEntries[key]
	item, config, fileName := getConfigItemFileNameByKey(key)
	if entry == nil {
		if config == nil {
			return nil, errorCacheNotFound
		}
		if fileName != "" {
			stat, err := os.Stat(fileName)
			if err == nil {
				if time.Now().Sub(stat.ModTime()) < time.Duration(config.fileDuration)*time.Second {
					value, err := ioutil.ReadFile(fileName)
					if err != nil {
						log.Printf("Error in cache: %v", err)
					} else {
						SetCacheValueAndSave(key, value, false)
						return value, nil
					}
				} else {
					os.Remove(fileName)
				}
			}
		}
	} else if config == nil || time.Now().Sub(entry.stamp) < time.Duration(config.memoryDuration)*time.Second {
		return entry.value, nil
	}
	delete(cacheEntries, key)
	if config.valueGetter != nil {
		value, err := config.valueGetter(key, item)
		if err != nil {
			return nil, err
		}
		SetCacheValue(key, value)
		return value, nil
	}
	return nil, errorCacheNotFound
}

func SetCacheValueAndSave(key string, value interface{}, tosave bool) {
	cacheEntries[key] = &cacheEntry{value: value, stamp: time.Now()}
	if !tosave {
		return
	}
	_, _, fileName := getConfigItemFileNameByKey(key)
	if fileName == "" {
		return
	}
	v, err := ConvertInterfaceToByteArrayValue(key, value)
	if err != nil {
		return
	}
	ioutil.WriteFile(fileName, v, 0664)
}

func SetCacheValue(key string, value interface{}) {
	SetCacheValueAndSave(key, value, true)
}

func ResetCacheItems(items ...string) {
	for _, key := range items {
		delete(cacheEntries, key)
		_, _, fileName := getConfigItemFileNameByKey(key)
		if fileName != "" {
			os.RemoveAll(fileName)
		}
	}
}

func ResetCacheGroups(groups ...string) {
	for _, group := range groups {
		for key := range cacheEntries {
			grp, _ := GetGroupByKey(key)
			if grp == group {
				delete(cacheEntries, key)
			}
		}
		config := cacheConfigs[group]
		var fileName string
		if config != nil {
			if config.folderName == "" {
				config.folderName = GetSafeFileName(group)
			}
			fileName = GetCacheLocalFileInfo(config.folderName, "")
		} else {
			fileName = GetCacheLocalFileInfo(GetSafeFileName(group), "")
		}
		if fileName != "" {
			os.RemoveAll(fileName)
		}
	}
}

func ResetAllLocalFileCache() {
	if cacheMainFolder == "" {
		cacheMainFolder = GetTempPathSlashed() + "___dv__server__cache/"
	}
	os.RemoveAll(cacheMainFolder)
}

func ResetTotalCache() {
	ResetAllLocalFileCache()
	cacheEntries = make(map[string]*cacheEntry)
}

func GetCacheStringValue(key string) (string, error) {
	v, err := GetCacheValue(key)
	if err != nil {
		return "", err
	}
	switch v.(type) {
	case string:
		return v.(string), nil
	case []byte:
		return string(v.([]byte)), nil
	case int:
		return strconv.Itoa(v.(int)), nil
	case int64:
		return strconv.FormatInt(v.(int64), 10), nil
	case float64:
		return strconv.FormatFloat(v.(float64), 'E', -1, 64), nil
	}
	return "", errors.New("Different type in cache for " + key)
}

func GetCacheByteArrayValue(key string) ([]byte, error) {
	v, err := GetCacheValue(key)
	if err != nil {
		return nil, err
	}
	return ConvertInterfaceToByteArrayValue(key, v)
}

func ConvertInterfaceToByteArrayValue(key string, v interface{}) ([]byte, error) {
	switch v.(type) {
	case []byte:
		return v.([]byte), nil
	case string:
		return []byte(v.(string)), nil
	case int:
		return []byte(strconv.Itoa(v.(int))), nil
	case int64:
		return []byte(strconv.FormatInt(v.(int64), 10)), nil
	case float64:
		return []byte(strconv.FormatFloat(v.(float64), 'E', -1, 64)), nil
	}
	return nil, errors.New("Different type in cache for " + key)
}

func GetCacheInt64Value(key string) (int64, error) {
	v, err := GetCacheValue(key)
	if err != nil {
		return 0, err
	}
	switch v.(type) {
	case int64:
		return v.(int64), nil
	case int:
		return int64(v.(int)), nil
	case float64:
		return int64(v.(float64)), nil
	case string:
		return strconv.ParseInt(v.(string), 10, 64)
	case []byte:
		return strconv.ParseInt(string(v.([]byte)), 10, 64)
	}
	return 0, errors.New("Different type in cache for " + key)
}

func GetCacheFloat64Value(key string) (float64, error) {
	v, err := GetCacheValue(key)
	if err != nil {
		return 0, err
	}
	switch v.(type) {
	case float64:
		return v.(float64), nil
	case int64:
		return float64(v.(int64)), nil
	case int:
		return float64(v.(int)), nil
	case string:
		return strconv.ParseFloat(v.(string), 64)
	case []byte:
		return strconv.ParseFloat(string(v.([]byte)), 64)
	}
	return 0, errors.New("Different type in cache for " + key)
}
