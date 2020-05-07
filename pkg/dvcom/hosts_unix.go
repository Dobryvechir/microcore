// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcom

import (
	"log"
	"os/exec"
)

var EOL_BYTES = []byte{10}

func getHostsConfigFileName() (string, error) {
	return "/etc/hosts", nil
}

func resolveAdministrativeTasks(tasks []administrativeTask, params string) {
	//TODO: to be developed
}

func executeOtherProgram(name string, params string) {
	cmd := exec.Command(name, params)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Print("Error: " + err.Error())
	}
	log.Printf("%s\n", stdoutStderr)
}

func addToDNSSearchList(list []string) {
	//TODO: to be developed
}
