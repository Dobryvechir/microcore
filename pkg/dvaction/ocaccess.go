/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvaction

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type CredentialProvider func(string) (string, string, bool)

var credentialProviders = make([]CredentialProvider, 0, 3)

func RegisterCredentialProvider(provider CredentialProvider) {
	credentialProviders = append(credentialProviders, provider)
}

func resolveCredentials(name string) (string, string, bool) {
	n := len(credentialProviders)
	for i := 0; i < n; i++ {
		us, ps, ok := credentialProviders[i](name)
		if ok {
			return us, ps, true
		}
	}
	return "", "", false
}

func CleanM2MTokenInCache(microserviceName string) {
	dvparser.RemoveGlobalPropertiesValue(M2MTokenPrefix + GetMicroServicePropertyName(microserviceName))
	path := dvparser.GlobalProperties[M2MTokenPath]
	if path == "" || microserviceName == "" {
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

func getPathForMicroserviceInfoStorage(microserviceName string) string {
	path := dvparser.GlobalProperties[M2MTokenPath]
	if path == "" {
		return ""
	}
	if path[len(path)-1] == '/' || path[len(path)-1] == '\\' {
		if microserviceName != "" {
			path += microserviceName
		} else {
			path = path[:len(path)-1]
		}
	} else {
		if microserviceName != "" {
			path += "/" + microserviceName
		}
	}
	return path
}

func GetM2MTokenFromFileSystem(microserviceName string) (token string, ok bool) {
	path := getPathForMicroserviceInfoStorage(microserviceName)
	if path == "" {
		return
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
	path := getPathForMicroserviceInfoStorage(microserviceName)
	if path == "" {
		return
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
	var username, passwrd string
	if microserviceName == "" {
		token, okFinal = dvparser.GlobalProperties[M2MTokenPrefix]
		if okFinal {
			return
		}
		token, okFinal = GetM2MTokenFromFileSystem(microserviceName)
		if okFinal {
			return
		}
		username = dvparser.GlobalProperties[M2M_USERNAME]
		passwrd = dvparser.GlobalProperties[M2M_PASSWORD]
	} else {
		token, okFinal = dvparser.GlobalProperties[M2MTokenPrefix+GetMicroServicePropertyName(microserviceName)]
		if okFinal {
			return
		}
		token, okFinal = GetM2MTokenFromFileSystem(microserviceName)
		if okFinal {
			return
		}
		username, passwrd, okFinal = resolveCredentials(microserviceName)
		if !okFinal {
			return
		}
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
	body := []string{
		"grant_type",
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
	if microserviceName == "" {
		dvparser.GlobalProperties[M2MTokenPrefix] = authorization
	} else {
		dvparser.SetGlobalPropertiesValue(M2MTokenPrefix+GetMicroServicePropertyName(microserviceName), authorization)
	}
	SaveM2MTokenInFileSystem(microserviceName, accessToken)
	return authorization, true
}

func GetMicroServicePropertyName(microServiceName string) string {
	return strings.ToUpper(strings.ReplaceAll(microServiceName, "-", "_"))
}

func NetRequest(method string, url string, body string, headers map[string]string, options map[string]interface{},forceM2mSimple bool) ([]byte, error, http.Header, int) {
	_, m2mSimple := headers[M2M]
	microServiceName, m2mComplex := headers[Authorization]
	if m2mComplex {
		if strings.HasPrefix(microServiceName, M2M_) {
			microServiceName = microServiceName[len(M2M_):]
		} else {
			m2mComplex = false
		}
	}
	if !m2mSimple && !m2mComplex {
		return dvnet.NewRequest(method, url, body, headers, options)
	}
	if m2mSimple {
		microServiceName = ""
		delete(headers, M2M)
	}
	m2mToken, ok := GetM2MToken(microServiceName)
	if !ok {
		return nil, errors.New("Cannot get M2M token " + microServiceName), nil, 500
	}
	headers[Authorization] = m2mToken
	if forceM2mSimple {
		return dvnet.NewRequest(method, url, body, headers, options)
	}
	data, err, heads, stat := dvnet.NewRequestRepeatPause(method, url, body, headers, options, 1, 0)
	if err == nil && stat < 500 && stat != 401 && stat != 403 {
		return data, nil, heads, stat
	}
	CleanM2MTokenInCache(microServiceName)
	m2mToken, ok = GetM2MToken(microServiceName)
	if !ok {
		return data, errors.New("Cannot get M2M token for " + microServiceName), heads, stat
	}
	headers[Authorization] = m2mToken
	return dvnet.NewRequest(method, url, body, headers, options)
}

func GetIdentityProviderClientCredentials(microServiceName string) (user string, pw string, ok bool) {
	//TODO: create client credentials in keycloak
	res := base64.StdEncoding.EncodeToString([]byte(microServiceName))
	return res, res, true
}
