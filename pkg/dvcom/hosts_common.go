/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvcom

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"log"
	"strconv"
	"strings"
)

type IpList struct {
	Kind string `json:"kind"`
	Ip   string `json:"ip"`
	Urls string `json:"urls"`
}

type hostsConfigLine struct {
	value string
	keys  []string
}

type hostsConfig struct {
	readLines  ReadLinesPool
	lineConfig []hostsConfigLine
	keyConfig  map[string]int
	changed    bool
}

type administrativeTask struct {
	name    string
	options []string
}

var administrativeTasks = make([]administrativeTask, 0, 5)

func deleteWholeLine(lineNo int, conf *hostsConfig) {
	n := len(conf.readLines.lines)
	if lineNo < 0 || lineNo >= n {
		return
	}
	if lineNo == n-1 {
		conf.readLines.lines = conf.readLines.lines[:lineNo]
		conf.lineConfig = conf.lineConfig[:lineNo]
	} else {
		conf.readLines.lines = append(conf.readLines.lines[:lineNo], conf.readLines.lines[lineNo+1:]...)
		conf.lineConfig = append(conf.lineConfig[:lineNo], conf.lineConfig[lineNo+1:]...)
		for k, v := range conf.keyConfig {
			if v > lineNo {
				conf.keyConfig[k] = v - 1
			}
		}
	}
}

func makeWholeLine(line *hostsConfigLine) string {
	if len(line.keys) == 0 {
		return ""
	}
	return line.value + " " + strings.Join(line.keys, " ")
}

func updateWholeLine(lineNo int, conf *hostsConfig) {
	n := len(conf.readLines.lines)
	if lineNo < 0 || lineNo >= n {
		return
	}
	conf.changed = true
	conf.readLines.lines[lineNo] = makeWholeLine(&conf.lineConfig[lineNo])
}

func deleteKeyMap(keyMap map[string]bool, conf *hostsConfig) {
	keys := conf.keyConfig
	for k, yet := range keyMap {
		lineNo, ok := keys[k]
		if ok {
			delete(keys, k)
			if yet {
				cell := &conf.lineConfig[lineNo]
				cellLen := len(cell.keys)
				for i := cellLen - 1; i >= 0; i-- {
					if _, okey := keyMap[cell.keys[i]]; okey {
						keyMap[cell.keys[i]] = false
						if i == cellLen-1 {
							cell.keys = cell.keys[:i]
						} else {
							cell.keys = append(cell.keys[:i], cell.keys[i+1:]...)
						}
						cellLen--
					}
				}
				if cellLen == 0 {
					deleteWholeLine(lineNo, conf)
					conf.changed = true
				} else {
					updateWholeLine(lineNo, conf)
					conf.changed = true
				}
			}
		}
	}
}

func readHostsConfig() hostsConfig {
	fileName, err := getHostsConfigFileName()
	conf := hostsConfig{changed: false}
	if err != nil {
		conf.readLines.err = err
		return conf
	}
	conf.readLines = readLinesFromFile(fileName)
	if conf.readLines.err == nil {
		n := len(conf.readLines.lines)
		conf.lineConfig = make([]hostsConfigLine, n)
		conf.keyConfig = make(map[string]int)
		for k, v := range conf.readLines.lines {
			v = strings.TrimSpace(v)
			if v != "" && v[0] != '#' {
				ss := strings.Split(v, " ")
				ssLen := len(ss)
				if ssLen > 1 {
					changed := false
					conf.lineConfig[k].value = ss[0]
					removed := ""
					for i := 1; i < ssLen; i++ {
						t := ss[i]
						if t != "" {
							if _, ok := conf.keyConfig[t]; ok {
								changed = true
								if removed == "" {
									removed = t
								} else {
									removed += "," + t
								}
							} else {
								conf.lineConfig[k].keys = append(conf.lineConfig[k].keys, t)
								conf.keyConfig[t] = k
							}
						}
					}
					if changed {
						conf.changed = true
						log.Printf("Hosts config corrected: at line %d removed %s", k, removed)
						if len(conf.lineConfig[k].keys) == 0 {
							conf.readLines.lines[k] = ""
						} else {
							updateWholeLine(k, &conf)
						}
					}
				}
			}
		}
	}
	return conf
}

func writeHostsConfig(conf *hostsConfig) {
	if conf.readLines.err == nil {
		writeLinesToFile(&conf.readLines)
	}
}

func addToHosts(ips map[string]string, conf *hostsConfig) {
	pool := make(map[string]int)
	currentLineNo := len(conf.lineConfig)
	startLineNo := currentLineNo
	for url, ip := range ips {
		if lineNo, ok := pool[ip]; ok {
			conf.lineConfig[lineNo].keys = append(conf.lineConfig[lineNo].keys, url)
			conf.keyConfig[url] = lineNo
		} else {
			lineNo = currentLineNo
			conf.lineConfig = append(conf.lineConfig, hostsConfigLine{value: ip, keys: []string{url}})
			pool[ip] = lineNo
			conf.keyConfig[url] = lineNo
			currentLineNo++
		}
	}
	if currentLineNo != startLineNo {
		if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogDebug {
			log.Printf("Add new hosts: currentLineNo %d startLineNo %d content %q", currentLineNo, startLineNo, conf.lineConfig[startLineNo:])
		}
		conf.changed = true
		for startLineNo < currentLineNo {
			conf.readLines.lines = append(conf.readLines.lines, makeWholeLine(&conf.lineConfig[startLineNo]))
			startLineNo++
		}
	}
}

func makeUrlList(url string) []string {
	urls := dvtextutils.ConvertToList(url)
	count := 0
	for i := 0; i < len(urls); i++ {
		a := urls[i]
		if a == "" {
			continue
		}
		if p := strings.Index(a, "//"); p >= 0 {
			a = a[p+2:]
		}
		if p := strings.Index(a, "/"); p >= 0 {
			a = a[:p]
		}
		if a == "" {
			if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogError {
				log.Printf("Incorrect url in the list %s", urls[i])
			}
			continue
		}
		urls[count] = a
		count++
	}
	return urls[:count]
}
func AddToHosts(ipList []IpList) error {
	conf := readHostsConfig()
	if conf.readLines.err == nil {
		toAdd := make(map[string]string)
		toDelete := make(map[string]bool)
		for _, v := range ipList {
			ip := v.Ip
			if ip == "" {
				ip = "127.0.0.1"
			}
			urls := makeUrlList(v.Urls)
			if v.Kind != "" {
				if strings.ToLower(v.Kind) == "dns" {
					addToDNSSearchList(urls)
				} else {
					if dvlog.CurrentLogLevel >= dvlog.LogError {
						log.Printf("Only \"kind\":\"DNS\" is allowed in hosts")
					}
				}
				continue
			}
			for _, k := range urls {
				if i, ok := toAdd[k]; ok {
					if i != ip {
						log.Print("Conflicting hosts in config: " + k)
					} else {
						log.Print("Duplicating hosts in config: " + k)
					}
				} else {
					if line, okay := conf.keyConfig[k]; okay {
						if conf.lineConfig[line].value == ip {
							continue
						}
						if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogInfo {
							log.Printf("Have to change host %s from %s to %s at %d", k, conf.lineConfig[line].value, ip, line)
						}
						toDelete[k] = true
					}
					toAdd[k] = ip
					if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogDetail {
						log.Printf("Adding %s to %s", k, ip)
					}

				}
			}
		}
		if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogDebug {
			log.Printf("Hosts results: added %q deleted %q changed %q", toAdd, toDelete, conf.changed)
		}
		deleteKeyMap(toDelete, &conf)
		addToHosts(toAdd, &conf)
		if conf.changed {
			writeHostsConfig(&conf)
		}
	}
	if dvlog.CurrentLogLevel >= dvlog.LogError {
		if conf.readLines.err != nil {
			log.Printf("Hosts adding problem: %s", conf.readLines.err.Error())
		} else if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogWarning {
			if conf.changed {
				log.Print("Hosts adding was successful!")
			} else {
				log.Print("Hosts adding was not necessary")
			}
		}
	}
	return conf.readLines.err
}

func RemoveFromHosts(ipList []IpList) error {
	conf := readHostsConfig()
	if conf.readLines.err == nil {
		keyMap := make(map[string]bool)
		for _, v := range ipList {
			urls := makeUrlList(v.Urls)
			if v.Kind != "" {
				continue
			}
			for _, url := range urls {
				if url != "" {
					keyMap[url] = true
				}
			}
		}
		deleteKeyMap(keyMap, &conf)
		if conf.changed {
			writeHostsConfig(&conf)
		}
	}
	if dvlog.CurrentLogLevel >= dvlog.LogError {
		if conf.readLines.err != nil {
			log.Printf("Hosts removal problem: %s", conf.readLines.err.Error())
		} else if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogWarning {
			if conf.changed {
				log.Print("Hosts removal was successful!")
			} else {
				log.Print("Hosts removal was not necessary")
			}
		}
	}
	return conf.readLines.err
}

func ResolveAdministrativeTasks() {
	if len(administrativeTasks) > 0 {
		params := ""
		for _, task := range administrativeTasks {
			options, count := dvtextutils.PrepareAndMayQuoteParams(task.options)
			params += " " + task.name + "_" + strconv.Itoa(count) + " " + options
		}
		resolveAdministrativeTasks(administrativeTasks, params)
	}
}

func AddAdministrativeTask(name string, options []string) {
	administrativeTasks = append(administrativeTasks, administrativeTask{name: name, options: options})
}

func ProcessHosts(ipList []IpList, isRemoval bool) {
	if len(ipList) == 0 {
		return
	}
	if isRemoval {
		RemoveFromHosts(ipList)
	} else {
		AddToHosts(ipList)
	}
}
