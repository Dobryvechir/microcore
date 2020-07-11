package dvcontext

import (
	"regexp"
	"sync"
)

type ValidatePattern struct {
	Source       string `json:"source"`
	Message      string `json:"message"`
	EmptyMessage string `json:"empty"`
	Contains     string `json:"contains"`
	RegPattern   string `json:"pattern"`
	containMap   map[byte]int
	regular      *regexp.Regexp
	inited       bool
	mutex        sync.Mutex
}

var validatePatternContainMapCache = make(map[string]map[byte]int)
var validatePatternRegularExpCache = make(map[string]*regexp.Regexp)
var containMapMutex sync.Mutex
var regularExprMutex sync.Mutex

func PrepareValidatePattern(pattern *ValidatePattern, defContainMap map[byte]int, defRegPattern *regexp.Regexp) error {
	contains := []byte(pattern.Contains)
	n := len(contains)
	if n == 0 {
		pattern.containMap = nil
	} else if defContainMap != nil {
		pattern.containMap = defContainMap
	} else {
		m := make(map[byte]int)
		pattern.containMap = m
		for i := 0; i < n; i++ {
			b := contains[i]
			if b == '\\' {
				i++
				if i == n {
					break
				}
				b = contains[i]
				switch b {
				case 'w':
					for i := 0; i < 26; i++ {
						m[byte(i+65)] = 1
						m[byte(i+97)] = 1
					}
					for i := 0; i < 10; i++ {
						m[byte(i+48)] = 1
					}
					m['_'] = 1
					break
				case 'd':
					for i := 0; i < 10; i++ {
						m[byte(i+48)] = 1
					}
					break
				case 'D':
					for i := 0; i < 10; i++ {
						m[byte(i+48)] = 1
					}
					m[','] = 1
					break
				case 'h':
					for i := 0; i < 10; i++ {
						m[byte(i+48)] = 1
					}
					for i := 0; i < 6; i++ {
						m[byte(i+65)] = 1
						m[byte(i+97)] = 1
					}
					m['-'] = 1
					break
				case 'W':
					for i := 0; i < 26; i++ {
						m[byte(i+65)] = 1
						m[byte(i+97)] = 1
					}
					for i := 0; i < 10; i++ {
						m[byte(i+48)] = 1
					}
					m['_'] = 1
					m[' '] = 1
					m[','] = 1
					m['-'] = 1
					break
				default:
					m[b] = 1
				}
			} else if b == '!' {
				i++
				if i == n {
					break
				}
				b = contains[i]
				if _, ok := m[b]; ok {
					delete(m, b)
				}
			} else {
				m[b] = 1
			}
		}
	}
	var err error = nil
	if pattern.RegPattern == "" {
		pattern.regular = nil
	} else if defRegPattern != nil {
		pattern.regular = defRegPattern
	} else {
		pattern.regular, err = regexp.Compile(pattern.RegPattern)
	}
	pattern.inited = err == nil
	return err
}

func MatchByContainMap(src []byte, containMap map[byte]int) bool {
	n := len(src)
	for i := 0; i < n; i++ {
		if containMap[src[i]] == 0 {
			return false
		}
	}
	return true
}

func (pattern *ValidatePattern) MatchByCompiledPattern(src string) string {
	if len(src) == 0 {
		return pattern.EmptyMessage
	}
	if pattern.containMap != nil {
		res := MatchByContainMap([]byte(src), pattern.containMap)
		if !res {
			return pattern.Message
		}
	}
	if pattern.regular != nil {
		res := pattern.regular.MatchString(src)
		if !res {
			return pattern.Message
		}
	}
	return ""
}

func (pattern *ValidatePattern) Match(src string) string {
	if pattern == nil {
		return ""
	}
	if !pattern.inited {
		defContainMap := validatePatternContainMapCache[pattern.Contains]
		defRegular := validatePatternRegularExpCache[pattern.RegPattern]
		pattern.mutex.Lock()
		err := PrepareValidatePattern(pattern, defContainMap, defRegular)
		pattern.mutex.Unlock()
		if err != nil {
			return "500 Internal System Error:" + err.Error()
		}
		if defContainMap == nil && pattern.containMap != nil {
			containMapMutex.Lock()
			validatePatternContainMapCache[pattern.Contains] = pattern.containMap
			containMapMutex.Unlock()
		}
		if defRegular == nil && pattern.regular != nil {
			regularExprMutex.Lock()
			validatePatternRegularExpCache[pattern.RegPattern] = pattern.regular
			regularExprMutex.Unlock()
		}
	}
	return pattern.MatchByCompiledPattern(src)
}
