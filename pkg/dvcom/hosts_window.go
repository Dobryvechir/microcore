// +build windows

/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcom

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"log"
	"os"
	"os/exec"
	"strings"
)

var EOL_BYTES = []byte{13, 10}

func getHostsConfigFileName() (string, error) {
	fileName := os.Getenv("windir")
	if fileName == "" {
		fileName = "c:/windows"
	}
	fileName += "/System32/drivers/etc/hosts"
	return fileName, nil
}

func resolveAdministrativeTasks(tasks []administrativeTask, params string) {
	params = "/C start dvadmin.exe" + params
	if LogHosts && dvlog.CurrentLogLevel >= dvlog.LogDetail {
		log.Print("Starting cmd.exe: ", params)
	}
	executeOtherProgram("cmd.exe", params)
}

func executeOtherProgram(name string, params string) {
	cmd := exec.Command(name, params)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print("Error: " + err.Error())
	}
	log.Printf("%s\n", stdoutStderr)
}

func addToDNSSearchList(urls []string) {
	list := strings.Join(urls, ",")
	regKey1 := `HKEY_LOCAL_MACHINE\SYSTEM\ControlSet001\services\Tcpip\Parameters~SZ~SearchList`
	regKey2 := `HKEY_LOCAL_MACHINE\SYSTEM\ControlSet002\services\Tcpip\Parameters~SZ~SearchList`
	AddAdministrativeTask("regadd", []string{regKey1, list, regKey2, list})
}
