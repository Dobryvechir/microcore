// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
)

const (
	TRANSFORM_BASE64 = iota
	TRANSFORM_MOVE_FILE
)

const (
	filePrefix = "file:"
)

type FormatDescription struct {
	DataPrefix   string
	Extension    string
	FormatLetter string
	Transform    int
}

var FormatDefinitions = []*FormatDescription{
	{
		DataPrefix:   "data:image/png;base64,",
		Extension:    ".png",
		FormatLetter: "i",
		Transform:    TRANSFORM_BASE64,
	},
	{
		DataPrefix:   "data:image/jpeg;base64,",
		Extension:    ".jpg",
		FormatLetter: "i",
		Transform:    TRANSFORM_BASE64,
	},
	{
		DataPrefix:   "data:image/gif;base64,",
		Extension:    ".gif",
		FormatLetter: "i",
		Transform:    TRANSFORM_BASE64,
	},
	{
		DataPrefix:   "data:video/mp4;base64,",
		Extension:    ".mp4",
		FormatLetter: "v",
		Transform:    TRANSFORM_BASE64,
	},
	{
		DataPrefix:   "data:video/ogg;base64,",
		Extension:    ".ogg",
		FormatLetter: "v",
		Transform:    TRANSFORM_BASE64,
	},
	{
		DataPrefix:   "data:video/webm;base64,",
		Extension:    ".webm",
		FormatLetter: "v",
		Transform:    TRANSFORM_BASE64,
	},
}

func saveUserFile(fileName string, position int, transform int, data string, src string) (int64, error) {
	var buf []byte = nil
	var err error = nil
	switch transform {
	case TRANSFORM_BASE64:
		data = data[position:]
		buf, err = base64.StdEncoding.DecodeString(data)
	case TRANSFORM_MOVE_FILE:
		err = dvdir.RenameOrMoveFile(src, fileName, false)
		if err != nil {
			return 0, err
		}
		fi, err := os.Stat(fileName)
		if err != nil {
			return 0, err
		}
		return fi.Size(), nil
	}
	if err != nil {
		return 0, err
	}
	if buf != nil {
		err = os.WriteFile(fileName, buf, 0644)
		return int64(len(buf)), err
	}
	return 0, errors.New("unsupported transformation")
}

func analyzeUserFile(data string) (extension string, position int, prefix string, transorm int, src string) {
	n := len(FormatDefinitions)
	if strings.HasPrefix(data, filePrefix) {
		data = data[len(filePrefix):]
		p := strings.Index(data, ":")
		if p <= 0 {
			return
		}
		extension = dvtextutils.GetLowCaseExtension(data[:p])
		for i := 0; i < n; i++ {
			v := FormatDefinitions[i]
			if v.Extension == extension {
				return v.Extension, 0, v.FormatLetter, TRANSFORM_MOVE_FILE, data[p+1:]
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		v := FormatDefinitions[i]
		if strings.HasPrefix(data, v.DataPrefix) {
			return v.Extension, len(v.DataPrefix), v.FormatLetter, v.Transform, ""
		}
	}
	return
}
