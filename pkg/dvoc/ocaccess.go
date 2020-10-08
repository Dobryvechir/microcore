/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvoc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvdir"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	M2MTokenPrefix          = "M2M_TOKEN_"
	M2MTokenPath            = "M2MTOKEN_PATH"
	M2MAuthorizationRequest = "AuthorizationM2MService"
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

func CleanM2MTokenInCache(microserviceName string) {
	dvparser.RemoveGlobalPropertiesValue(M2MTokenPrefix + GetMicroServicePropertyName(microserviceName))
	path := dvparser.GlobalProperties[M2MTokenPath]
	if path == "" {
		return
	}
	if path[len(path)-1] == '/' || path[len(path)-1] == '\\' {
		path += microserviceName
	} else {
		path += "/" + microserviceName
	}
	path += "/token"
	_, err := os.Stat(path)
	if err != nil {
		return
	}
	os.Remove(path)
}

func GetM2MTokenFromFileSystem(microserviceName string) (token string, ok bool) {
	path := dvparser.GlobalProperties[M2MTokenPath]
	if path == "" {
		return
	}
	if path[len(path)-1] == '/' || path[len(path)-1] == '\\' {
		path += microserviceName
	} else {
		path += "/" + microserviceName
	}
	path += "/token"
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile(path)
	if err != nil || len(data) == 0 {
		return
	}
	var accessToken = &AccessToken{}
	err = json.Unmarshal(data, accessToken)
	if err != nil || accessToken.AccessToken == "" {
		return
	}
	difInSeconds := time.Now().Sub(info.ModTime()) / time.Second
	if difInSeconds >= time.Duration(accessToken.ExpiresIn) {
		CleanM2MTokenInCache(microserviceName)
		return
	}
	token = accessToken.TokenType + " " + accessToken.AccessToken
	ok = true
	dvparser.SetGlobalPropertiesValue(M2MTokenPrefix+GetMicroServicePropertyName(microserviceName), token)
	return
}

func SaveM2MTokenInFileSystem(microserviceName string, accessToken *AccessToken) {
	path := dvparser.GlobalProperties[M2MTokenPath]
	if path == "" {
		return
	}
	if path[len(path)-1] == '/' || path[len(path)-1] == '\\' {
		path += microserviceName
	} else {
		path += "/" + microserviceName
	}
	os.MkdirAll(path, 0755)
	path += "/token"
	data, err := json.Marshal(accessToken)
	if err != nil {
		return
	}
	ioutil.WriteFile(path, data, 0644)
}

func GetM2MToken(microserviceName string) (token string, okFinal bool) {
	token, okFinal = GetM2MTokenFromFileSystem(microserviceName)
	if okFinal {
		return
	}
	username, passwrd, ok := GetOpenshiftSecrets(microserviceName)
	if !ok {
		return
	}
	m2mTokenUrlRaw := dvparser.GlobalProperties["M2MTOKEN_URL"]

	if m2mTokenUrlRaw == "" {
		dvlog.PrintlnError("Specify M2MTOKEN_URL in .properties")
		return
	}
	m2mTokenUrl, err1 := dvparser.ConvertByteArrayByGlobalProperties([]byte(m2mTokenUrlRaw), "M2M TOKEN URL")
	if err1 != nil {
		dvlog.PrintfError("Make sure you specified all constants %s in .properties: %v", m2mTokenUrlRaw, err1)
		return
	}
	body := []string{"grant_type",
		"client_credentials",
		"client_id",
		username,
		"client_secret",
		passwrd,
	}
	headers := map[string]string{"cache-control": "no-cache", "Content-Type": "application/x-www-form-urlencoded"}
	var accessToken = &AccessToken{}
	err := dvnet.LoadStructFormUrlEncoded("POST", m2mTokenUrl, body, headers, accessToken, dvnet.AveragePersistentOptions)
	if accessToken.TokenType == "" || accessToken.AccessToken == "" {
		err = dvnet.LoadStructFormUrlEncoded("POST", m2mTokenUrl, body, headers, &accessToken, dvnet.AveragePersistentOptions)
		if accessToken.TokenType == "" || accessToken.AccessToken == "" {
			dvlog.PrintfError("Cannot get M2M Access Token for %s (%v)", microserviceName, err)
			return
		}
	}
	authorization := accessToken.TokenType + " " + accessToken.AccessToken
	dvparser.SetGlobalPropertiesValue(M2MTokenPrefix+GetMicroServicePropertyName(microserviceName), authorization)
	SaveM2MTokenInFileSystem(microserviceName, accessToken)
	return authorization, true
}

func GetMicroServicePropertyName(microServiceName string) string {
	return strings.ToUpper(strings.ReplaceAll(microServiceName, "-", "_"))
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

func NetRequest(method string, url string, body string, headers map[string]string, options map[string]interface{}) ([]byte, error) {
	if options[M2MAuthorizationRequest] == nil || options[M2MAuthorizationRequest] == "" {
		return dvnet.NewRequest(method, url, body, headers, options)
	}
	repeatOption := options["repeats"]
	options["repeats"] = 1
	microServiceName := options[M2MAuthorizationRequest].(string)
	m2mToken, ok := GetM2MToken(microServiceName)
	if !ok {
		return nil, errors.New("Cannot get M2M token for " + microServiceName)
	}
	headers[Authorization] = m2mToken
	data, err := dvnet.NewRequest(method, url, body, headers, options)
	options["repeats"] = repeatOption
	if err == nil {
		return data, nil
	}
	message := err.Error()
	repeatNumber := 0
	if repeatOption != nil {
		repeatNumber = repeatOption.(int)
	}
	if repeatNumber == 0 && method == "POST" {
		repeatNumber = 1
	}
	if repeatNumber != 1 || (len(message) >= 3 && (message[:3] == "401" || message[:3] == "403")) {
		CleanM2MTokenInCache(microServiceName)
		m2mToken, ok = GetM2MToken(microServiceName)
		if !ok {
			return nil, errors.New("Cannot get M2M token for " + microServiceName)
		}
		headers[Authorization] = m2mToken
		return dvnet.NewRequest(method, url, body, headers, options)
	}
	return data, err
}

func GetIdentityProviderClientCredentials(microServiceName string) (user string, pw string, ok bool) {
	//TODO: create client credentials in keycloak
	res := base64.StdEncoding.EncodeToString([]byte(microServiceName))
	return res, res, true
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
