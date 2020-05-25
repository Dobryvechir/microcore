/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"io/ioutil"
	"strings"
)

const (
	StarKey = "StarKey"
)

func IsCorrectHash(s string) bool {
	n := len(s)
	if n == 0 {
		return false
	}
	for i := 0; i < n; i++ {
		c := s[i]
		if !(c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' || c >= '0' && c <= '9') {
			return false
		}
	}
	return true
}

func hashReplacerEndPointHandler(request *dvcontext.RequestContext) bool {
	url := request.Url
	finish := strings.LastIndex(url, ".")
	start := strings.LastIndex(url[:finish], ".") + 1
	if finish <= 0 || start <= 1 || !IsCorrectHash(url[start:finish]) || request.Server.BaseFolderUrl == "" {
		return false
	}
	dirName := ""
	prefix := url[:start]
	suffix := url[finish:]
	pos := strings.LastIndex(prefix, "/")
	if pos >= 0 {
		dirName := url[:pos]
		if dirName != "" && dirName[0] != '/' {
			dirName = "/" + dirName
		}
		prefix = prefix[pos+1:]
	}
	dirName = request.Server.BaseFolderUrl + dirName
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return false
	}
	for _, f := range files {
		name := f.Name()
		if strings.HasPrefix(name, prefix) && strings.HasSuffix(name, suffix) {
			request.FileName = dirName + "/" + name
			dvcom.HandleFromFile(request)
			return true
		}
	}
	return false
}

func hashReplacerServerInitHandler(params []string) (map[string]string, error) {
	starLen := len(params)
	if starLen != 1 {
		return nil, errors.New("module hashReplacer requires 1 parameter")
	}
	replacerName := params[0]
	if replacerName != "hash" {
		return nil, errors.New("Only 'hash' is supported, but you specified " + replacerName)
	}
	res := make(map[string]string)
	res[StarKey] = replacerName
	return res, nil
}

var hashReplacerConfig = &RegistrationConfig{
	Name:              "hashReplacer",
	EndPointHandler:   hashReplacerEndPointHandler,
	ServerInitHandler: hashReplacerServerInitHandler,
}

func registerHashReplacer() bool {
	RegisterModule(hashReplacerConfig, false)
	return true
}

var hashReplacerInited = registerHashReplacer()
