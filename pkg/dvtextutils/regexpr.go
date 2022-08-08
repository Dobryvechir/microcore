/***********************************************************************
MicroCore
Copyright 2020 -2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

import (
	"regexp"
	"strconv"
	"strings"
)

type RegExpession struct {
	Pattern       string
	Flags         string
	Compiled      *regexp.Regexp
	ResultWord    string
	ResultIndices [][]int
	ResultCount   int
	GlobalSearch  bool
}

func NewRegExpression(pattern string, flags string) (*RegExpession, error) {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	v := &RegExpession{
		Pattern:      pattern,
		Flags:        flags,
		Compiled:     compiled,
		GlobalSearch: strings.Contains(pattern, "g") || strings.Contains(pattern, "y"),
	}
	return v, nil
}

func FindByRegularExpression(source string, regexpr string, group string, defvalue string, isall bool, count int) (interface{}, error) {
	reg, err := regexp.Compile(regexpr)
	if err != nil {
		return defvalue, err
	}
	grp := 0
	isIndex := false
	if group != "" {
		if group[0] >= '0' && group[0] <= '9' {
			grp, err = strconv.Atoi(group)
			if err == nil {
				isIndex = true
			}
		}
	}
	if isall {
		if isIndex {
			res := reg.FindAllString(source, grp)
			if count > 0 {
				if count <= len(res) {
					return res[count-1], nil
				}
				return defvalue, nil
			}
			return res, nil
		}
		res := reg.FindAllStringSubmatch(source, -1)
		if count > 0 {
			if count <= len(res) {
				return res[count-1], nil
			}
			return defvalue, nil
		}
		return res, nil
	}
	if isIndex {
		res := reg.FindStringSubmatch(source)
		if grp < len(res) {
			return res[grp], nil
		}
		return defvalue, nil
	}
	return reg.FindString(source), nil
}

func ReplaceByRegularExpression(source string, regexpr string, replacement string, isLiteral bool) (string, error) {
	reg, err := regexp.Compile(regexpr)
	if err != nil {
		return source, err
	}
	if isLiteral {
		return reg.ReplaceAllLiteralString(source, replacement), nil
	}
	return reg.ReplaceAllString(source, replacement), nil
}
