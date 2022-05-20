/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

import (
	"net/url"
	"strings"
)

func EncodeURI(s string) string {
	s = url.PathEscape(s)
	if strings.Contains(s, "%2F") {
		s = strings.Replace(s, "%2F", "/", -1)
	}
	if strings.Contains(s, "%3F") {
		s = strings.Replace(s, "%3F", "?", -1)
	}
	return s
}

func DecomposeUrlEncodedString(s string) (r map[string][]string, err error) {
	r = make(map[string][]string, 16)
	p := strings.Index(s, "?")
	if p >= 0 {
		s = s[p+1:]
	}
	t := strings.Split(s, "&")
	for _, v := range t {
		p := strings.Index(v, "=")
		if p <= 0 {
			continue
		}
		key, er := url.QueryUnescape(strings.TrimSpace(v[:p]))
		if len(key) == 0 || er != nil {
			if er != nil {
				err = er
			}
			continue
		}
		val, er := url.QueryUnescape(strings.TrimSpace(v[p+1:]))
		if er != nil {
			err = er
			continue
		}
		mapVal, ok := r[key]
		if ok {
			r[key] = append(mapVal, val)
		} else {
			r[key] = []string{val}
		}
	}
	return
}

func ConvertArrayStringMapIntoSimpleStringMap(m map[string][]string) map[string]string {
	r := make(map[string]string, len(m))
	for k, v := range m {
		if len(v) > 0 {
			r[k] = v[0]
		}
	}
	return r
}
