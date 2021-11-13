/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"os"
	"strings"
)

func SaveOpenshiftSecret(folder string, microServiceName string) bool {
	if folder != "" {
		c := folder[len(folder)-1]
		os.MkdirAll(folder, 0755)
		if c != '/' && c != '\\' {
			folder += "/"
		}
	}
	username, passwrd, ok := GetOpenshiftSecrets(microServiceName)
	if !ok {
		return false
	}
	err := ioutil.WriteFile(folder+"username", []byte(username), 0644)
	err1 := ioutil.WriteFile(folder+"password", []byte(passwrd), 0644)
	if err != nil {
		dvlog.Printf("Cannot save username for %s: %v", microServiceName, err)
		return false
	}
	if err1 != nil {
		dvlog.Printf("Cannot save password for %s: %v", microServiceName, err1)
		return false
	}
	return true
}

func CreateClientCredentials(user string, pw string, microserviceName string) bool {
	yaml := "apiVersion: v1\ndata:\n  password: >-\n    " +
		pw +
		"\n  username: " +
		user +
		"\nkind: Secret\nmetadata:\n  name: " +
		microserviceName +
		"-client-credentials" +
		"\n  namespace: " +
		dvparser.GlobalProperties["OPENSHIFT_NAMESPACE"] +
		"\ntype: Opaque\n"
	path := dvdir.GetTempPathSlashed() + "__dobryvechir__debug_fragments_secret.yaml"
	err := ioutil.WriteFile(path, []byte(yaml), 0664)
	if err != nil {
		dvlog.PrintfError("Cannot temporarily save a secret file %s", path)
		return false
	}
	cmdLine := "create -f " + path
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return false
	} else {
		dvlog.PrintlnError(info)
	}
	return true

}

func RegisterUserCredentialsWithIdentityProvider(user, pw, microServiceName string, onBehalfMicroserviceName string, tmpSuffix string) bool {
	cmdLine := strings.ReplaceAll(openShiftSecrets, "%1", "identity-provider")
	info, ok := RunOCCommand(cmdLine)
	if !ok {
		return false
	}
	pos := strings.Index(info, "  "+microServiceName+":")
	if pos > 0 {
		return true
	}
	pos = strings.Index(info, "  "+onBehalfMicroserviceName+":")
	if pos < 0 {
		dvlog.PrintlnError(info)
		dvlog.PrintlnError("could not find mui-platform here")
		return false
	}
	line := "  " + microServiceName + ": " + pw + "\n"
	newSecret := info[:pos] + line + info[pos:]
	path := dvdir.GetTempPathSlashed() + tmpSuffix + ".yaml"
	err := ioutil.WriteFile(path, []byte(newSecret), 0664)
	if err != nil {
		dvlog.PrintfError("Cannot temporarily save a secret file %s", path)
		return false
	}
	cmdLine = openShiftDeleteSecret + "identity-provider-client-credentials"
	info, ok = RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return false
	} else {
		dvlog.PrintlnError(info)
	}

	cmdLine = "create -f " + path
	info, ok = RunOCCommand(cmdLine)
	if !ok {
		dvlog.PrintfError("Failed to execute %s", cmdLine)
		return false
	} else {
		dvlog.PrintlnError(info)
	}
	return true
}

func provideKeycloakAutoUpdate() {

	//LATER TODO:  find the way to start autoupdate of keycloak

}
