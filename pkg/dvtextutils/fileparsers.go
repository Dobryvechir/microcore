/***********************************************************************
MicroCore
Copyright 2020 -2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvtextutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

func LoadSimpleMapFromByteArray(data []byte, separator byte) map[string]string {
	n := len(data)
	res := make(map[string]string)
	for i := 0; i < n; i++ {
		for ; i < n && data[i] <= 32; i++ {
		}
		startPos := i
		separPos := -1
		for ; i < n && data[i] != 13 && data[i] != 10; i++ {
			if separPos < 0 && data[i] == separator {
				separPos = i
			}
		}
		if separPos >= 0 {
			key := strings.TrimSpace(string(data[startPos:separPos]))
			value := strings.TrimSpace(string(data[separPos+1 : i]))
			res[key] = value
		}
	}
	return res
}

func LoadSimpleMapFromFile(fileName string, separator byte) (res map[string]string, err error) {
	var data []byte
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	return LoadSimpleMapFromByteArray(data, separator), nil
}

func LookInDescriptionByteArray(data []byte, lookNames map[string]string, mustPresent bool) (res map[string]string, err error) {
	n := len(data)
	res = make(map[string]string)
	for i := 0; i < n; i++ {
		for ; i < n && data[i] <= 32; i++ {
		}
		startPos := i
		separatorPosition := -1
		for ; i < n && data[i] != 13 && data[i] != 10; i++ {
			if separatorPosition < 0 && data[i] == ':' {
				separatorPosition = i
			}
		}
		if separatorPosition >= 0 {
			key := strings.TrimSpace(string(data[startPos:separatorPosition]))
			value := strings.TrimSpace(string(data[separatorPosition+1 : i]))
			if newKey, ok := lookNames[key]; ok {
				if newKey == "" {
					newKey = key
				}
				res[newKey] = value
			}
		}
	}
	if mustPresent && len(lookNames) != len(res) {
		keys := "Values"
		for k, _ := range lookNames {
			if _, okay := res[k]; !okay {
				keys += " " + k
			}
		}
		err = errors.New(keys + " not found in description")
	}
	return
}

func LookInDescriptionFile(fileName string, lookNames map[string]string, mustPresent bool) (res map[string]string, err error) {
	var data []byte
	data, err = ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	return LookInDescriptionByteArray(data, lookNames, mustPresent)
}

func PresentMapAsProperties(pool map[string]string, outputFileName string) error {
	res := ""
	for k, v := range pool {
		res += k + "=" + v + "\n"
	}
	if outputFileName == "" {
		fmt.Printf("%s", res)
		return nil
	}
	return ioutil.WriteFile(outputFileName, []byte(res), 0644)
}

func SeparateToKeyValue(line string) (key string, value string) {
	pos := strings.Index(line, ":")
	if pos <= 0 {
		pos = strings.Index(line, "=")
	}
	if pos > 0 {
		key = strings.TrimSpace(line[:pos])
		value = strings.TrimSpace(line[pos+1:])
	}
	return
}

func PutDescribedAttributesToMapFromList(params map[string]string, pool map[string]string, list []string) error {
	for _, line := range list {
		if line == "" {
			continue
		}
		k, v := SeparateToKeyValue(line)
		if k == "" {
			c := line[len(line)-1]
			if c == '?' {
				line = line[:len(line)-1]
				if params[line] == "" {
					continue
				}
			}
			l := params[line]
			if l != "" {
				k, v = SeparateToKeyValue(l)
			}
		}
		if k == "" {
			return errors.New("Parameter " + line + " is not determined")
		}
		pool[k] = v
	}
	return nil
}

func PutDescribedAttributesToMapFromCommaSeparatedList(params map[string]string, pool map[string]string, list string) error {
	return PutDescribedAttributesToMapFromList(params, pool, ConvertToNonEmptyList(list))
}

func ComposeParametersInTemplate(template string, brackets int, value string) (res map[string]string, ok bool) {
	res = make(map[string]string)
	ok = true
	n := len(template)
	m := 0
	for i := 0; i < n; i++ {
		if template[i] == '{' {
			m++
		} else {
			if m >= brackets {
				pos := i
				m = 0
				for ; i < n; i++ {
					if template[i] == '}' {
						m++
						if m == brackets {
							k := strings.TrimSpace(template[pos : i-m+1])
							res[k] = value
						}
					} else {
						m = 0
					}
				}
				if m < brackets && i == n {
					ok = false
				}
			}
			m = 0
		}
	}
	return
}
