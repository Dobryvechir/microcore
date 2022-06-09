/***********************************************************************
MicroCore
Copyright 2017 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvlog

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LogFatal = -1
	LogError = iota
	LogWarning
	LogInfo
	LogDetail
	LogDebug
	LogTrace
)

var StartTime = time.Now().Unix() - 1654740000
var _startTimeMutex = &sync.Mutex{}
var logsFolder string
var CurrentLogLevel = LogError
var CurrentNamespace string
var CurrentRootFolder string

var logTable = map[string]int{
	"fatal":    LogFatal,
	"error":    LogError,
	"warning":  LogWarning,
	"info":     LogInfo,
	"detail":   LogDetail,
	"debug":    LogDebug,
	"trace":    LogTrace,
	"internal": LogTrace,
}

func FormErrorMessage(err string) []byte {
	return []byte("{\"errorMessage\":\"" + strings.Replace(err, "\"", "\\\"", -1) + "\"}")
}

func GetLogLevel(level string) int {
	v, ok := logTable[strings.ToLower(level)]
	if !ok {
		return LogError
	}
	return v
}

func GetPrincipalFolder(ensure bool) string {
	path := CurrentRootFolder
	if path == "" {
		path = os.Getenv("HOME")
	}
	if path == "" {
		path = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	if path == "" {
		path = os.Getenv("TEMP")
	}
	if path == "" {
		path = os.Getenv("TMP")
	}
	if path == "" {
		path = "/tmp"
	}
	c := path[len(path)-1]
	if c != '/' && c != '\\' {
		path += "/"
	}
	path += "MicroCore"
	if ensure {
		os.MkdirAll(path, os.ModePerm)
	}
	return path
}

func GetSubFolder(mainFolder string, subFolder string) string {
	if subFolder == "" {
		subFolder = "/"
	}
	c := subFolder[0]
	if c != '/' && c != '\\' {
		mainFolder += "/"
	}
	subFolder = strings.Replace(subFolder, ":", "_", -1)
	mainFolder += subFolder
	mainFolder = EnsureFolder(mainFolder)
	return mainFolder
}

func EnsureFolder(folder string) string {
	if len(folder) > 0 {
		c := folder[len(folder)-1]
		if c == '/' || c == '\\' {
			folder = folder[:len(folder)-1]
		}
	}
	if folder != "" {
		os.MkdirAll(folder, os.ModePerm)
	}
	return folder
}

func EnsureFolderCheck(folder string) (string, bool) {
	folder = EnsureFolder(folder)
	if folder == "" {
		return folder, false
	}
	if _, err := os.Stat(folder); err != nil {
		return folder, false
	}
	return folder, true
}

func EnsureSubFolderCheck(folderUser string, folderSub string) (folder string, res bool) {
	folder, res = EnsureFolderCheck(folderUser)
	if res {
		return
	}
	folder, res = EnsureFolderCheck(GetTaskSubFolder(folderSub))
	return
}

func EnsureSubFolderCritical(folderUser string, folderSub string) string {
	folder, res := EnsureSubFolderCheck(folderUser, folderSub)
	if res {
		return folder
	}
	panic("incorrect path for " + folderSub)
	return folder
}

func GetTaskFolder() string {
	return GetSubFolder(GetPrincipalFolder(true), CurrentNamespace)
}

func GetTaskSubFolder(subFolder string) string {
	return GetSubFolder(GetTaskFolder(), subFolder)
}

func WriteRequestToLog(body []byte, r *http.Request) string {
	return WriteNetRequestToLog(body, r.Method, r.URL.Path, CurrentLogLevel >= LogDetail, r.Host, "R", r.Header, nil)
}

func WriteNetRequestToLog(body []byte, method string, url string, moreDetail bool,
	host string, rpref string, headersDouble map[string][]string, headersSingle map[string]string) string {
	number := int(GetNextUniqueNumber())
	if logsFolder == "" {
		logsFolder = GetTaskSubFolder("LOGS")
	}
	logFile := logsFolder + "/" + rpref + strconv.Itoa(number) + method
	fileName := logFile + "R." + GetSafeFileName(url)
	ioutil.WriteFile(fileName, body, os.ModePerm)
	if moreDetail {
		headers := make([]byte, 0, 1024)
		headers = append(headers, []byte(method)...)
		headers = append(headers, ' ')
		headers = append(headers, []byte(url)...)
		headers = append(headers, 13, 10)
		count := 0
		if host != "" {
			headers = append(headers, []byte("Host")...)
			headers = append(headers, ':', ' ')
			headers = append(headers, []byte(host)...)
			headers = append(headers, 13, 10)
			count++
		}
		if headersDouble != nil {
			for name, headerValues := range headersDouble {
				for _, h := range headerValues {
					headers = append(headers, []byte(name)...)
					headers = append(headers, ':', ' ')
					headers = append(headers, []byte(h)...)
					headers = append(headers, 13, 10)
					count++
				}
			}
		}
		if headersSingle != nil {
			for name, headerValue := range headersSingle {
				headers = append(headers, []byte(name)...)
				headers = append(headers, ':', ' ')
				headers = append(headers, []byte(headerValue)...)
				headers = append(headers, 13, 10)
				count++
			}
		}
		fileName = logFile + "H.I" + strconv.Itoa(count)
		ioutil.WriteFile(fileName, headers, os.ModePerm)
	}
	return logFile
}

func WriteResponseToLog(logFile string, resp *http.Response, body []byte) {
	WriteNetResponseToLog(logFile, body, nil, resp.Header, resp.StatusCode, CurrentLogLevel >= LogDetail)
}

func WriteNetResponseToLog(logFile string, body []byte, err error, headsDouble map[string][]string,
	status int, moreDetail bool) {
	respStatus := strconv.Itoa(status)
	fileName := logFile + "S." + GetSafeFileName(respStatus)
	ioutil.WriteFile(fileName, body, os.ModePerm)
	if err != nil {
		fileName = logFile + "E." + GetSafeFileName(respStatus)
		ioutil.WriteFile(fileName, []byte(err.Error()), os.ModePerm)
	}
	if moreDetail {
		headers := make([]byte, 0, 1024)
		headers = append(headers, []byte(respStatus)...)
		headers = append(headers, 13, 10)
		count := 0
		for name, headerValues := range headsDouble {
			for _, h := range headerValues {
				headers = append(headers, []byte(name)...)
				headers = append(headers, ':', ' ')
				headers = append(headers, []byte(h)...)
				headers = append(headers, 13, 10)
				count++
			}
		}
		fileName = logFile + "H.O" + strconv.Itoa(count)
		ioutil.WriteFile(fileName, headers, os.ModePerm)
	}
}

func SetLogLevel(logLevel string) {
	CurrentLogLevel = GetLogLevelByDefinition(logLevel, CurrentLogLevel)
}

func GetLogLevelByDefinition(logLevel string, defaultLevel int) int {
	logLevel = strings.ToLower(strings.TrimSpace(logLevel))
	level, ok := logTable[logLevel]
	if ok {
		return level
	}
	if logLevel == "" {
		return defaultLevel
	}
	if defaultLevel >= LogError {
		log.Print("logLevel can be debug, detail, info, warning, error, and none, but you specified " + logLevel)
	}
	return defaultLevel
}

func GetSafeFileName(dat string) string {
	b := []byte(dat)
	for i, c := range b {
		if c <= 32 || c > 126 || c == '\\' || c == ':' || c == '>' || c == '<' || c == '"' || c == '|' || c == '/' {
			b[i] = '_'
		} else if c == '?' {
			b[i] = '~'
		} else if c == '*' {
			b[i] = '$'
		}
		if i > 200 {
			break
		}
	}
	if len(b) > 200 {
		b = b[:200]
	}
	return string(b)
}

func GetNextUniqueNumber() int64 {
	var number int64
	_startTimeMutex.Lock()
	StartTime++
	number = StartTime
	_startTimeMutex.Unlock()
	return number
}

func SetCurrentNamespace(namespace string) {
	namespace = strings.Replace(namespace, ":", "_", -1)
	namespace = strings.Replace(namespace, "\\", "_", -1)
	namespace = strings.Replace(namespace, "/", "_", -1)
	namespace = strings.Replace(namespace, " ", "_", -1)
	if namespace != "" {
		CurrentNamespace = namespace
	}
}

func CleanEOL(dat []byte) {
	for i, c := range dat {
		if c < 32 {
			dat[i] = 32
		}
	}
}

func StartingLogFile() {
	if CurrentLogLevel != LogFatal {
		logFileName := GetTaskSubFolder("LOGS") + "/log" + strconv.Itoa(int(StartTime)) + ".txt"
		logFile, err1 := os.Create(logFileName)
		if err1 != nil {
			log.Fatal("Failed to create a log file " + logFileName + ": " + err1.Error())
			return
		}
		defer logFile.Close()
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)

	}
}

func GetTemporaryFileName() string {
	nmb := GetNextUniqueNumber()
	name := GetTaskSubFolder("TMP") + "/tmp" + strconv.FormatInt(nmb, 10)
	return name
}

func FindConfigFile(folder string, name string) (fileName string, err error) {
	if len(name) == 0 {
		err = errors.New("File name is not specified")
		return
	}
	c := name[0]
	if c == '\\' || c == '/' || c == '.' && (len(name) == 1 || name[1] == '.' || name[1] == '/' || name[1] == '\\') {
		fileName = name
	} else {
		fileName = folder + "/" + name
	}
	if _, err1 := os.Stat(fileName); err1 != nil {
		err = err1
	}
	return
}

func EnsureDirForFileExists(fileName string) bool {
	if fileName == "" {
		return false
	}
	dir := filepath.Dir(fileName)
	if dir == "" || dir == "." || dir == ".." {
		return true
	}
	return os.MkdirAll(dir, os.ModePerm) == nil
}

func LogLevelToInt(levelWord string) (level int, ok bool) {
	level = GetLogLevelByDefinition(levelWord, -1)
	ok = level >= 0
	if !ok {
		level = LogError
	}
	return
}

func LogLevelFromProperties(parameters map[string]string, propertyName string, defLevel int) int {
	parameter := parameters[propertyName]
	CheckDvLogStreamToFile(parameters)
	if len(parameter) == 0 {
		return defLevel
	}
	level, ok := LogLevelToInt(parameter)
	if !ok {
		log.Printf("Unknown value %s of property %s (allowed: none, error, warning, info, detail, bigDetail)", propertyName, parameter)
		return defLevel
	}
	return level
}
