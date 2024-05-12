// package dvdbmanager provides functions for database query
// MicroCore Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvdbmanager

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

const (
	TRANSFORM_BASE64 = iota
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

func saveUserFile(fileName string, position int, transform int, data string) (int64, error) {
	data = data[position:]
	var buf []byte = nil
	var err error = nil
	switch transform {
	case TRANSFORM_BASE64:
		buf, err = base64.StdEncoding.DecodeString(data)
	}
	if err != nil {
		return 0, err
	}
	if buf != nil {
		err := os.WriteFile(fileName, buf, 0644)
		return int64(len(buf)), err
	}
	return 0, errors.New("Unsupported transformation")
}

func analyzeUserFile(data string) (extension string, position int, prefix string, transorm int) {
	n := len(FormatDefinitions)
	for i := 0; i < n; i++ {
		v := FormatDefinitions[i]
		if strings.HasPrefix(data, v.DataPrefix) {
			return v.Extension, len(v.DataPrefix), v.FormatLetter, v.Transform
		}
	}
	return
}
