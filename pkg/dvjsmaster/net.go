/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvjsmaster

import (
	"bytes"
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type NetServerInfo struct {
	sync.RWMutex
	client              *http.Client
	ExtraServerSettings ServerSettings
}

type ServerSettings struct {
	MaxIdleConns          int  `json:"maxIdleConnections"`
	IdleConnTimeout       int  `json:"idleConnectionTimeout"`
	DisableCompression    bool `json:"disableCompression"`
	DisableKeepAlives     bool `json:"disableKeepAlives"`
	MaxIdleConnsPerHost   int  `json:"maxIdleConnectionsPerHost"`
	ResponseHeaderTimeout int  `json:"responseHeaderTimeout"`
	ExpectContinueTimeout int  `json:"expectContinueTimeout"`
}

var NetMaster *dvevaluation.DvVariable
var netServerInfo *NetServerInfo = &NetServerInfo{}

func Net_GetText(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "GET", "PLAIN", "TEXT")
}

func Net_PostText(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "POST", "JSON", "TEXT")
}

func Net_PutText(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "PUT", "JSON", "TEXT")
}

func Net_DeleteText(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "DELETE", "JSON", "TEXT")
}

func Net_RequestText(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], params[3].GetStringValue(), "JSON", "TEXT")
}

func Net_Get(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "GET", "PLAIN", "JSON")
}

func Net_Post(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "POST", "JSON", "JSON")
}

func Net_Put(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "PUT", "JSON", "JSON")
}

func Net_Delete(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], "DELETE", "JSON", "JSON")
}

func Net_Request(context *dvevaluation.DvContext, thisVariable *dvevaluation.DvVariable, params []*dvevaluation.DvVariable) (*dvevaluation.DvVariable, error) {
	return net_request(params[0], params[1], params[2], params[3].GetStringValue(), "JSON", "JSON")
}

func net_request(url *dvevaluation.DvVariable, data *dvevaluation.DvVariable, headers *dvevaluation.DvVariable, method string, inputFormat string, outputFormat string) (*dvevaluation.DvVariable, error) {
	urlStr := url.GetStringValue()
	if urlStr == "" {
		return nil, errors.New("Cannot make net request with empty url")
	}
	var bodyStr []byte
	if inputFormat == "JSON" {
		bodyStr = data.JsonStringifyNonEmpty()
	} else {
		bodyStr = data.GetStringValueAsBytes()
	}
	bodyIo := ioutil.NopCloser(bytes.NewReader(bodyStr))
	req, err := http.NewRequest(method, urlStr, bodyIo)
	if err != nil {
		if dvlog.CurrentLogLevel >= dvlog.LogError {
			log.Print("Error making request %s: %s", urlStr, err.Error())
		}
		return nil, err
	}
	req.Header = headers.GetStringArrayMap()
	if netServerInfo.client == nil {
		createClientForNetServerInfo(netServerInfo)
	}
	resp, err1 := netServerInfo.client.Do(req)
	if err1 != nil {
		if dvlog.CurrentLogLevel >= dvlog.LogError {
			log.Printf("Error executing %s: %s", urlStr, err1.Error())
		}
		return nil, err1
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		if dvlog.CurrentLogLevel >= dvlog.LogError {
			log.Printf("Error reading body %s: %s", urlStr, err2.Error())
		}
		return nil, err2
	}
	var reply *dvevaluation.DvVariable
	if resp.StatusCode < 400 {
		if outputFormat == "JSON" {
			reply, err1 = JSON_parse_direct(body, "Net request "+urlStr)
		} else {
			reply = dvevaluation.DvVariableFromString(nil, string(body))
		}
	} else {
		message := "Error in net response to " + urlStr + " [" + string(body) + "]"
		return nil, errors.New(message)
	}
	return reply, nil
}

func createClientBySettings(settings *ServerSettings) *http.Client {
	tr := &http.Transport{
		DisableKeepAlives:     settings.DisableKeepAlives,
		MaxIdleConnsPerHost:   settings.MaxIdleConnsPerHost,
		ResponseHeaderTimeout: time.Duration(time.Duration(settings.ResponseHeaderTimeout) * time.Second),
		ExpectContinueTimeout: time.Duration(time.Duration(settings.ExpectContinueTimeout) * time.Second),
		MaxIdleConns:          settings.MaxIdleConns,
		IdleConnTimeout:       time.Duration(time.Duration(settings.IdleConnTimeout) * time.Second),
		DisableCompression:    settings.DisableCompression,
	}
	return &http.Client{Transport: tr, Timeout: time.Second * 60}
}

func createClientForNetServerInfo(server *NetServerInfo) {
	server.Lock()
	if server.client != nil {
		server.Unlock()
		return
	}
	server.client = createClientBySettings(&server.ExtraServerSettings)
	server.Unlock()
}

func net_init() {
	NetMaster = dvevaluation.RegisterMasterVariable("Net", &dvevaluation.DvVariable{
		Fields: make(map[string]*dvevaluation.DvVariable),
		Kind:   dvevaluation.FIELD_OBJECT,
		Prototype: &dvevaluation.DvVariable{
			Fields: map[string]*dvevaluation.DvVariable{
				"getText": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_GetText,
				},
				"postText": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_PostText,
				},
				"putText": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_PutText,
				},
				"deleteText": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_DeleteText,
				},
				"requestText": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_RequestText,
				},
				"get": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_Get,
				},
				"post": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_Post,
				},
				"put": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_Put,
				},
				"delete": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_Delete,
				},
				"request": {
					Kind: dvevaluation.FIELD_FUNCTION,
					Fn:   Net_Request,
				},
			},
			Kind: dvevaluation.FIELD_OBJECT,
		},
	})

}
