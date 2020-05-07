/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"strings"
)

const (
	contentSecurityPolicyHeader = "Content-Security-Policy"
	addHeaderInConfig           = "add_header"
)

func UpdateContentSecurityPolicyOnPods(pods string, hosts string) bool {
	podList := dvparser.ConvertToNonEmptyList(pods)
	hosts = strings.Replace(hosts, "http://", "", -1)
	hosts = strings.Replace(hosts, "https://", "", -1)
	hosts = strings.Replace(hosts, "ws://", "", -1)
	hosts = strings.Replace(hosts, "wss://", "", -1)
	hosts = strings.TrimSpace(hosts)
	if hosts == "" {
		return true
	}
	n := len(podList)
	for i := 0; i < n; i++ {
		pod := podList[i]
		dat, err := OpenShiftReadTextFile(pod, StrategyReadWriteSingleFileBest)
		if err != nil {
			dvlog.PrintfError("Cannot read %s : %v", pod, err)
			return false
		}
		if strings.Index(dat, "location") < 0 {
			dvlog.PrintfError("Cannot read %s nginx config (no location): %v", pod, dat)
			return false
		}
		dat = IntroduceContentSecurityPolicyInNginxConfig(dat, hosts)
		err = OpenShiftWriteTextFile(pod, dat, StrategyReadWriteSingleFileThruDir)
		if err != nil {
			dvlog.PrintfError("Cannot write %s : %v", pod, err)
			return false
		}
		err = OpenShiftNginxRestart(pod)
		if err != nil {
			dvlog.PrintfError("Failed to restart nginx at %s : %v", pod, err)
			return false
		}
	}
	return true
}

func FindAddHeaderInNginxConfig(data string, headerName string, afterPos int) (posStart int, posEnd int) {
	n := len(data)
	posEnd = -1
	for afterPos < n {
		pos := strings.Index(data[afterPos:], headerName) + afterPos
		if pos < afterPos {
			break
		}
		if !strings.HasSuffix(strings.TrimSpace(data[:pos]), addHeaderInConfig) {
			afterPos = pos + 1
			continue
		}
		afterPos = pos + len(headerName)
		for ; afterPos < n && data[afterPos] <= ' '; afterPos++ {
		}
		if afterPos == n {
			break
		}
		if data[afterPos] != '"' {
			continue
		}
		posStart = afterPos + 1
		for pos = posStart; pos < n && data[pos] != '"'; pos++ {
		}
		if pos < n {
			posEnd = pos
			return
		}
	}
	return
}

func FixContentSecurityPolicyLine(originalLine string, source string) string {
	line := strings.Split(originalLine, ";")
	n := len(line)
	for i := 0; i < n; i++ {
		s := line[i]
		if strings.TrimSpace(s) == "" {
			continue
		}
		if strings.Index(s, source) < 0 {
			line[i] = s + " " + source
		}
	}
	return strings.Join(line, ";")
}

func IntroduceContentSecurityPolicyInNginxConfig(data string, hosts string) string {
	for pos := 0; pos < len(data); pos++ {
		start, end := FindAddHeaderInNginxConfig(data, contentSecurityPolicyHeader, pos)
		if end < 0 {
			return data
		}
		fix := FixContentSecurityPolicyLine(data[start:end], hosts)
		data = data[:start] + fix + data[end:]
		pos = start + len(fix)
	}
	return data
}

func OpenShiftNginxRestart(name string) error {
	info, err := ExecuteCommandOnPod(name, "nginx -s reload")
	if err != nil {
		return err
	}
	if strings.Index(info, "error") >= 0 {
		return errors.New(info)
	}
	return nil
}
