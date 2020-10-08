/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvdir

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

type DirScanListGetter func(path string) ([]string, error)
type DirScanList struct {
	list []string
	path string
}

func (d *DirScanList) findInList(s string) int {
	n := len(d.list)
	for i := 0; i < n; i++ {
		if s == d.list[i] {
			return i
		}
	}
	return -1
}

func (d *DirScanList) addToList(s string, negative bool) {
	n := d.findInList(s)
	if negative {
		if n >= 0 {
			d.list = append(d.list[:n], d.list[n+1:]...)
		}
	} else {
		if n < 0 {
			d.list = append(d.list, s)
		}
	}
}

func (d *DirScanList) addArrayToList(s []string, negative bool) {
	n := len(s)
	for i := 0; i < n; i++ {
		d.addToList(s[i], negative)
	}
}

func DirScanProvider(data map[string][]string, currentDir string, dirListProc DirScanListGetter) ([]*DirScanList, error) {
	res := make([]*DirScanList, 0, 3)
	for k, v := range data {
		var info []string
		var err error
		k = strings.TrimSpace(k)
		if len(k) == 0 || k == "*" {
			k = currentDir
		}
		r := &DirScanList{list: make([]string, 0, 30), path: k}
		for _, i := range v {
			i = strings.TrimSpace(i)
			if len(i) == 0 {
				continue
			}
			isNegative := i[0] == '!'
			if isNegative {
				i = strings.TrimSpace(i[1:])
			}
			if len(i) == 0 {
				continue
			}
			if strings.Index(i, "*") >= 0 {
				if info == nil {
					info, err = dirListProc(k)
					if err!=nil {
						return nil, err
					}
				}
				lst := WildMaskStringScan(info, i)
				r.addArrayToList(lst, isNegative)
			} else {
				r.addToList(i, isNegative)
			}
		}
		if len(r.list) != 0 {
			res = append(res, r)
		}
	}
	return res, nil
}

func WildMaskStringScan(lst []string, mask string) []string {
	mask = strings.TrimSpace(mask)
	n := len(lst)
	p := len(mask)
	if n == 0 || p == 0 {
		return nil
	}
	beginAny := mask[0] == '*'
	if beginAny && p == 1 {
		return lst
	}
	endAny := mask[p-1] == '*'
	m := dvparser.ReduceListToNonEmptyList(strings.Split(mask, "*"))
	ml := len(m)
	if ml == 0 {
		return lst
	}
	start := m[0]
	res := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s := lst[i]
		k := strings.Index(s, start)
		if k >= 0 && WildMaskInString(m, beginAny, endAny, s) {
			res = append(res, s)
		}
	}
	return res
}

func WildMaskInString(pattern []string, beginAny bool, endAny bool, s string) bool {
	n := len(pattern)
	if n == 0 {
		return beginAny || endAny
	}
	if !beginAny {
		if !strings.HasPrefix(s, pattern[0]) {
			return false
		}
		s = s[len(pattern[0]):]
		pattern = pattern[1:]
		n--
		if n == 0 {
			return endAny || len(s) == 0
		}
	}
	if !endAny {
		if !strings.HasSuffix(s, pattern[n-1]) {
			return false
		}
		s = s[:len(pattern[n-1])]
		pattern = pattern[:n]
		n--
	}
	for i := 0; i < n; i++ {
		p := strings.Index(s, pattern[i])
		if p < 0 {
			return false
		}
		s = s[p+len(pattern[i]):]
	}
	return true
}

func DirScanProviderByString(data string, currentDir string, dirListProc DirScanListGetter) ([]*DirScanList, error) {
	m := make(map[string][]string)
	data = strings.TrimSpace(data)
	n := len(data)
	if n == 0 {
		return nil, nil
	}
	if data[0] == '{' && data[n-1] == '}' {
		data = strings.TrimSpace(data[1 : n-1])
	}
	if n == 0 {
		return nil, nil
	}
	p := strings.Index(data, "[")
	if p < 0 {
		m["*"] = strings.Split(data, ",")
	} else {
		d := strings.Split(data, "]")
		n = len(d)
		for i := 0; i < n; i++ {
			s := strings.TrimSpace(d[i])
			if len(s) == 0 {
				continue
			}
			if strings.HasPrefix(s, ",") {
				s = strings.TrimSpace(s[1:])
				if len(s) == 0 {
					continue
				}
			}
			p = strings.Index(s, "[")
			if p < 0 {
				return nil, errors.New("Corrupted json structure in " + data)
			}
			pk := strings.TrimSpace(s[:p])
			pv := strings.Split(strings.TrimSpace(s[p+1:]), ",")
			m[pk] = pv
		}
	}
	return DirScanProvider(m, currentDir, dirListProc)
}
