/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvoc

import (
	"encoding/json"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"errors"
	"io/ioutil"
	"strings"
)

func CopyToBaseFolder(baseFolder string, params string) error {
	n := len(baseFolder) - 1
	if n >= 0 && (baseFolder[n] == '/' || baseFolder[n] == '\\') {
		baseFolder = baseFolder[:n]
	}
	params = strings.TrimSpace(params)
	n = len(params)
	if n == 0 {
		return nil
	}
	if params[0] != '{' || params[n-1] != '}' {
		return errors.New("Expected {key:value, key1:value,...} but found " + params)
	}
	info := make(map[string]string)
	err := json.Unmarshal([]byte(params), &info)
	if err != nil {
		return errors.New("Files " + params + " must be a string map as follows:{\"key\":\"value\", \"key\":\"value\"}")
	}
	for k, v := range info {
		file, err := ioutil.ReadFile(v)
		if err != nil {
			return err
		}
		n = len(k)
		if n == 0 || (n == 1 && k == "/") {
			return errors.New("One of the destination file names has not been specified in " + params)
		}
		dst := baseFolder + k
		if k[0] != '/' && k[0] != '\\' {
			dst = baseFolder + "/" + k
		}
		dvdir.EnsureFolderForUsualFileToBeSaved(dst)
		err = ioutil.WriteFile(dst, file, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
