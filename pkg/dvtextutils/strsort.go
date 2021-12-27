/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

import (
	"sort"
	"strings"
)

func SortStringArray(src []string,options []string) []string {
	n:=len(options)
	for i:=0;i<n;i++ {
		src = SortStringArrayByOption(src,options[i])
	}
	return src
}

func SortStringArrayByOption(src []string,option string) []string {
	option=strings.TrimSpace(option)
	s:=strings.ToLower(option)
	switch s {
	case "asc":
		sort.Strings(src)
	case "desc":
		return SortStringArrayByOrder(src, true)
	}
	return src
}

func SortStringArrayByOrder(src []string, desc bool) []string {
	sort.Strings(src)
	if desc {
		src = ReverseStringArray(src)
	}
	return src
}

func ReverseStringArray(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
