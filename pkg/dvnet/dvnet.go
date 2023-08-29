/***********************************************************************
MicroCore
Copyright 2020 - 2023 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvnet

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: time.Duration(180 * time.Second),
	}).Dial,
	TLSHandshakeTimeout: time.Duration(180 * time.Second),
}

var netClientCommon = &http.Client{
	Timeout:   time.Duration(180 * time.Second),
	Transport: netTransport,
}

const (
	LogFatal   = 0
	LogError   = 1
	LogWarning = 2
	LogInfo    = 3
	LogDetail  = 4
	LogDebug   = 5
	LogTrace   = 6
)

var Log = LogError

var AveragePersistentOptions = map[string]interface{}{
	"repeats": 0,
	"pause":   5,
}

func GetAveragePersistentOptions() map[string]interface{} {
	options := make(map[string]interface{})
	for k, v := range AveragePersistentOptions {
		options[k] = v
	}
	return options
}

func GetIntegerOption(options map[string]interface{}, name string, defValue int) int {
	if options == nil {
		return defValue
	}
	value, ok := options[name]
	if !ok {
		return defValue
	}
	switch value.(type) {
	case int:
		return value.(int)
	case string:
		s := strings.TrimSpace(value.(string))
		if s == "" || !(s[0] >= '0' && s[0] <= '9' || s[0] == '+' || s[0] == '-') {
			return defValue
		}
		res, err := strconv.Atoi(s)
		if err != nil {
			return defValue
		}
		return res
	}
	return defValue
}

func NewRequest(method string, url string, body string, headers map[string]string, options map[string]interface{}) ([]byte, error, http.Header, int) {
	repeats := GetIntegerOption(options, "repeats", 0)
	pause := GetIntegerOption(options, "pause", 5)
	return NewRequestRepeatPause(method, url, body, headers, options, repeats, pause)
}

func NewRequestRepeatPause(method string, url string, body string, headers map[string]string, options map[string]interface{}, repeats int, pause int) ([]byte, error, http.Header, int) {
	insecureSkipVerify := true
	if Log >= LogInfo {
		if Log >= LogDetail {
			dvlog.Println(method+" "+url, fmt.Sprintf("%s %s [%s] %v", method, url, body, headers))
		} else {
			bodyStuff := body
			if len(bodyStuff) > 70 {
				bodyStuff = bodyStuff[:70] + "..."
			}
			bodyStuff = dvparser.CleanLogString(bodyStuff)
			dvlog.Println(method+" "+url, fmt.Sprintf("%s %s [%s] %v", method, url, bodyStuff, headers))
		}
	}
	if repeats <= 0 {
		repeats = 1
		if method != "POST" {
			repeats = 7
		}
	}
	httpHeaders := http.Header{}
	if headers != nil {
		for k, v := range headers {
			httpHeaders[k] = []string{v}
		}
	}
	var response *http.Response
	var err error = nil
	var request *http.Request
	var buf []byte
	var statusCode int
	var respHeaders http.Header
	for ; repeats > 0; repeats-- {
		bodyIo := ioutil.NopCloser(bytes.NewReader([]byte(body)))
		request, err = http.NewRequest(method, url, bodyIo)
		request.Header = httpHeaders
		if err != nil {
			if Log >= LogDebug {
				dvlog.PrintfFullOnly("NewRequest error [" + CreateRequestInfoForLog(method, url, body, httpHeaders) + "]")
			}
			return nil, err, nil, 400 // bad request
		}
		tr := http.DefaultTransport.(*http.Transport)
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureSkipVerify}
		} else {
			tr.TLSClientConfig.InsecureSkipVerify = insecureSkipVerify
		}
		tr.IdleConnTimeout = time.Duration(240 * time.Second)
		tr.TLSHandshakeTimeout = time.Duration(180 * time.Second)
		netClient := &http.Client{Timeout: 360 * time.Second, Transport: tr}
		response, err = netClient.Do(request)
		if err == nil {
			statusCode = response.StatusCode
			buf, err = ioutil.ReadAll(response.Body)
		} else {
			statusCode = 502 // Bad Gateway
		}
		if response != nil {
			response.Body.Close()
			respHeaders = response.Header
		}
		if (err != nil || statusCode >= 400) && Log >= LogDebug {
			dvlog.Printf("Http Error %s %v %d", CreateRequestInfoForLog(method, url, body, httpHeaders), err, statusCode)
		}
		if statusCode < 500 {
			break
		}
		if Log >= LogDebug {
			message := ""
			if err != nil {
				message = err.Error()
			} else {
				message = string(buf)
			}
			dvlog.PrintfFullOnly("Waiting %d s/%d because %d %s", pause, repeats, statusCode, message)
		}
		if pause > 0 && repeats > 1 {
			time.Sleep(time.Duration(pause) * time.Second)
		}
	}
	erroneous := statusCode >= 400
	if Log >= LogInfo || (Log >= LogError && erroneous) {
		var bufStuff string
		if Log < LogDebug && len(bufStuff) > 70 {
			bufStuff = dvparser.CleanLogByteArray(buf[:70]) + "..."
		} else {
			bufStuff = dvparser.CleanLogByteArray(buf)
		}
		s := fmt.Sprintf("Response: %d %s %v", statusCode, bufStuff, err)
		if erroneous {
			dvlog.PrintlnError(s)
		} else {
			dvlog.Println("", s)
		}
	}
	if err != nil {
		message := ""
		if response != nil {
			message = response.Status
		} else {
			message = err.Error()
		}
		return nil, errors.New(message + " " + string(buf)), nil, statusCode
	}
	return buf, nil, respHeaders, statusCode
}

func CreateRequestInfoForLog(method string, url string, body string, headers http.Header) string {
	return "(method=[" + method + "] url=[" + url + "] body=[" + body + "] headers=[" + CreateHeadersForLog(headers) + "]"
}

func CreateHeadersForLog(headers http.Header) string {
	s := ""
	for k, v := range headers {
		s += k + "=[" + CreateStringArrayForLog(v) + "]"
	}
	return s
}

func CreateStringArrayForLog(v []string) string {
	n := len(v)
	s := ""
	for i := 0; i < n; i++ {
		if i != 0 {
			s += ","
		}
		s += "'" + v[i] + "'"
	}
	return s
}

func NewJsonRequest(method string, url string, body string, headers map[string]string, options map[string]interface{}) ([]byte, error, http.Header, int) {
	if headers == nil {
		headers = make(map[string]string)
	}
	if headers["Content-Type"] == "" {
		headers["Content-Type"] = "application/json;charset=UTF-8"
	}
	if headers["Accept"] == "" {
		headers["Accept"] = "application/json"
	}
	return NewRequest(method, url, body, headers, options)
}

func UrlEncodedPart(data interface{}) string {
	params := url.Values{}
	switch data.(type) {
	case map[string]string:
		dataMap := data.(map[string]string)
		for k, v := range dataMap {
			params.Add(k, v)
		}
	case []string:
		key := ""
		dataArray := data.([]string)
		for _, v := range dataArray {
			if key == "" {
				key = v
			} else {
				params.Add(key, v)
				key = ""
			}
		}
	case string:
		return data.(string)
	}
	return params.Encode()
}

func LoadStruct(method string, url string, body string, headers map[string]string, v interface{}, options map[string]interface{}) error {
	buf, err, _, stat := NewRequest(method, url, body, headers, options)
	if err != nil {
		return err
	}
	if stat >= 400 {
		return errors.New("Failed net " + method + " " + url + " " + strconv.Itoa(stat))
	}
	return json.Unmarshal(buf, v)
}

func LoadStructFormUrlEncoded(method string, url string, body interface{}, headers map[string]string, v interface{}, options map[string]interface{}) error {
	data := UrlEncodedPart(body)
	if method == "GET" {
		if strings.Index("url", "?") > 0 {
			url += "&" + data
		} else {
			url += "?" + data
		}
		data = ""
	}
	return LoadStruct(method, url, data, headers, v, options)
}

func UpdatePropertiesThruNetRequest(request string) error {
	return errors.New("Not supported yet")
}

func InitBySettings(parameters map[string]string, functionPool map[string]interface{}) {
	Log = dvlog.LogLevelFromProperties(parameters, dvparser.MicroCorePrexix+"DVNET_LOG", Log)
}

var initedBySettings = dvparser.RegisterInitBySettingFunc(InitBySettings)
