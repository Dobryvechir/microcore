/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"regexp"
	"strings"
)

func ReadFileList(s string) []string {
	lst := dvtextutils.ConvertToNonEmptySemicolonList(s)
	n := len(lst)
	res := make([]string, 0, n*2)
	for i := 0; i < n; i++ {
		res = CollectFileListForSingleName(lst[i], res)
	}
	return res
}

func CollectFileListForSingleName(s string, pool []string) []string {
	pos := strings.Index(s, "^")
	pos1 := strings.LastIndex(s, "$")
	if pos >= 0 && pos1 > pos {
		pref := s[:pos]
		suffix := s[pos1+1:]
		reg := s[pos : pos1+1]
		return CollectFileListForRegularExpression(pref, suffix, reg, pool)
	}
	path, name := ParsePathName(s)
	if name == "" {
		name = "*"
	}
	if strings.Contains(name, "*") {
		return CollectFileListForWildCard(path, name, pool)
	}
	pool = append(pool, s)
	return pool
}

func ParsePathName(s string) (string, string) {
	pos := strings.LastIndex(s, "/")
	pos1 := strings.LastIndex(s, "\\")
	if pos < pos1 {
		pos = pos1
	}
	path := "."
	name := s
	if pos >= 0 {
		path = s[:pos]
		name = s[pos+1:]
		if path == "" {
			path = "."
		}
	}
	return path, name
}

func CollectFileListForRegularExpression(pref string, suffix string, reg string, pool []string) []string {
	path, name := ParsePathName(pref)
	pref = name
	exp, err := regexp.Compile(reg)
	if err != nil {
		dvlog.PrintfError("Failed to parse regular expression %s", reg)
		return pool
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		dvlog.PrintfError("Failed to read dir %v ", err)
		return pool
	}
	prefLen := len(pref)
	sufLen := len(suffix)
	for _, file := range files {
		if file.Mode().IsRegular() {
			flNm := file.Name()
			_, nm := ParsePathName(flNm)
			if sufLen != 0 {
				if !strings.HasSuffix(nm, suffix) {
					continue
				}
				nm = nm[:len(nm)-sufLen]
			}
			if prefLen != 0 {
				if !strings.HasPrefix(nm, pref) {
					continue
				}
				nm = nm[prefLen:]
			}
			if !exp.MatchString(nm) {
				continue
			}
			pool = append(pool, flNm)
		}
	}
	return pool
}
func CollectRegularFileList(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		dvlog.PrintfError("Failed to read dir %v ", err)
		return nil
	}
	n := len(files)
	res := make([]string, 0, n)
	for i := 0; i < n; i++ {
		if files[i].Mode().IsRegular() {
			res = append(res, files[i].Name())
		}
	}
	return res
}

func CollectFileListForWildCard(path string, name string, pool []string) []string {
	files := CollectRegularFileList(path)
	if len(files) == 0 {
		return pool
	}
	list := WildMaskFullPathScan(files, name)
	n := len(list)
	if n > 0 {
		if path != "" && path[len(path)-1] != '\\' && path[len(path)-1] != '/' {
			path += "/"
		}
		for i := 0; i < n; i++ {
			pool = append(pool, path+list[i])
		}
	}
	return pool
}
