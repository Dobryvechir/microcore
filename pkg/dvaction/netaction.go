/***********************************************************************
MicroCore
Copyright 2020 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvaction

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvjson"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvnet"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	Authorization = "Authorization"
)

const (
	M2MTokenPrefix = "M2M_TOKEN_"
	M2MTokenPath   = "M2MTOKEN_PATH"
	M2M            = "M2M"
	M2M_USERNAME   = "M2M_USERNAME"
	M2M_PASSWORD   = "M2M_PASSWORD"
	M2M_           = "M2M_"
)

type SmartNetConfigTemplate struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type SmartNetConfig struct {
	Url           string                 `json:"url"`
	Method        string                 `json:"method"`
	Headers       string                 `json:"headers"`
	Template      SmartNetConfigTemplate `json:"template"`
	ContentType   string                 `json:"type"`
	Result        string                 `json:"result"`
	DefaultResult string                 `json:"default"`
	Body          string                 `json:"body"`
}

type ProxyNetConfig struct {
	Url               string `json:"url"`
	Method            string `json:"method"`
	Headers           string `json:"headers"`
	ContentType       string `json:"type"`
	Result            string `json:"result"`
	Body              string `json:"body"`
	NotProxyHeaders   bool   `json:"not_proxy_headers"`
	NotProxyUrlParams bool   `json:"not_proxy_url_params"`
	NotAddUrlPath     bool   `json:"not_add_url_path"`
	NotProxyBody      bool   `json:"not_proxy_body"`
	NotReturnHeaders  bool   `json:"not_return_headers"`
}

func convertToHeader(list []string) (res map[string]string) {
	n := len(list)
	res = make(map[string]string)
	for i := 0; i < n; i++ {
		s := strings.TrimSpace(list[i])
		if s == "" {
			continue
		}
		pos := strings.Index(s, "=")
		if pos <= 0 {
			dvlog.PrintfError("Incorrect header %s (no =)", s)
			continue
		}
		k := strings.TrimSpace(s[:pos])
		v := strings.TrimSpace(s[pos+1:])
		res[k] = v
	}
	return
}

func processNetInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	params := dvtextutils.ConvertToNonEmptyList(command)
	url := params[0]
	headers := convertToHeader(params[1:])
	return []interface{}{url, headers}, true
}

func processNetRun(data []interface{}) bool {
	url := data[0].(string)
	headers := data[1].(map[string]string)
	_, err, _ := dvnet.NewRequest("GET", url, "", headers, dvnet.AveragePersistentOptions)
	return err == nil
}

func SmartNetInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &SmartNetConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Url == "" {
		log.Printf("Url must be specified in %s", command)
		return nil, false
	}
	switch config.Method {
	case "":
		config.Method = "GET"
	case "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE":
	default:
		log.Printf("Unknown method %s in %s", config.Method, command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func SmartNetRun(data []interface{}) bool {
	config := data[0].(*SmartNetConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return SmartNetRunByConfig(config, ctx)
}

func ProcessSavingActionResult(result string, data interface{}, ctx *dvcontext.RequestContext, err error, message1 string, message2 string) bool {
	if err != nil {
		log.Printf("%s %s %v", message1, message2, err)
		return false
	}
	SaveActionResult(result, data, ctx)
	return true
}

func smartNetProcessBody(body string, ctx *dvcontext.RequestContext) (string, bool) {
	if len(body) > 2 && body[0] == '@' && body[1] == '@' {
		k := body[2:]
		n := len(k)
		strict := true
		filled := false
		if n > 1 && k[n-1] == '?' {
			strict = false
			if k[n-2] == '?' {
				filled = true
				k = k[:n-2]
			} else {
				k = k[:n-1]
			}
		}
		env := GetEnvironment(ctx)
		v, ok := env.Get(k)
		if !ok {
			if strict {
				log.Printf("Empty #{k} in net request")
				if ctx != nil {
					ctx.SetHttpErrorCode(400, "Unset "+k)
				}
				return "", false
			}
			v = nil
		}
		if filled {
			if dvjson.IsEmptyAny(v) {
				return "", false
			}
		}
		body = dvevaluation.AnyToString(v)
	}
	return body, true
}

func smartNetProcessHeaders(headers map[string]string, newHeaders string, method string, body string) map[string]string {
	if newHeaders != "" {
		dvtextutils.PutDescribedAttributesToMapFromCommaSeparatedList(dvparser.GlobalProperties, headers, newHeaders)
	}
	if _, ok := headers["accept"]; !ok {
		headers["accept"] = "application/json"
	}
	if _, ok := headers["Content-Type"]; !ok && method != "GET" && IsLikeJson(body) {
		headers["Content-Type"] = "application/json"
	}
	return headers
}

func SaveHeaderResult(result string, heads http.Header, env *dvevaluation.DvObject) {
	env.Set("MC_NET_RESULT", result)
	result = "HEADERS_" + result + "_"
	for k, v := range heads {
		key := result + k
		if len(v) > 0 {
			env.Set(key, v[0])
		}
	}
}

func SmartNetRunByConfig(config *SmartNetConfig, ctx *dvcontext.RequestContext) bool {
	options := dvnet.GetAveragePersistentOptions()
	body, ok := smartNetProcessBody(config.Body, ctx)
	if !ok {
		return true
	}
	env := GetEnvironment(ctx)
	headers := make(map[string]string)
	headers = smartNetProcessHeaders(headers, config.Headers, config.Method, body)
	res, err, heads := NetRequest(config.Method, config.Url, body, headers, options)
	if err != nil {
		log.Println(res)
		log.Printf("%s %s failed: %v", config.Method, config.Url, err)
		if config.DefaultResult=="" {
			return false
		} else {
			res = []byte(config.DefaultResult)
		}
	}
	var result interface{}
	switch config.ContentType {
	case "json":
		result, err = dvjson.JsonFullParser(res)
		if err != nil {
			result = string(res)
		}
	default:
		result = string(res)
	}
	SaveHeaderResult(config.Result, heads, env)
	SaveActionResult(config.Result, result, ctx)
	if config.Template.Dst != "" {
		if config.Template.Src != "" {
			s, err := dvparser.ConvertByteArrayByGlobalPropertiesRuntime([]byte(config.Template.Src), "net step")
			if err != nil {
				log.Printf("Cannot convert %s: %v", config.Template.Src, err)
				return false
			}
			isFile := len(s) < 300 && strings.Index(s, "<") < 0 && strings.Index(s, ">") < 0 && strings.Index(s, "{") < 0
			if isFile {
				if _, err = os.Stat(s); err != nil {
					isFile = false
				}
			}
			if isFile {
				res, err = ioutil.ReadFile(s)
				if err != nil {
					log.Printf("Cannot read file %s: %v", s, err)
					return false
				}
			} else {
				res = []byte(s)
			}
		}
		result, err := dvparser.ConvertByteArrayByGlobalPropertiesRuntime(res, config.Url)
		if err != nil {
			log.Printf("Cannot convert for file %s: %v", config.Template.Dst, err)
			ioutil.WriteFile(config.Template.Dst, res, 0664)
			return false
		}
		err = ioutil.WriteFile(config.Template.Dst, []byte(result), 0664)
		if err != nil {
			log.Printf("Cannot write file %s: %v", config.Template.Dst, err)
			return false
		}
	}
	return true
}

func ProxyNetInit(command string, ctx *dvcontext.RequestContext) ([]interface{}, bool) {
	config := &ProxyNetConfig{}
	if !DefaultInitWithObject(command, config, GetEnvironment(ctx)) {
		return nil, false
	}
	if config.Url == "" {
		log.Printf("Url must be specified in %s", command)
		return nil, false
	}
	switch config.Method {
	case "", "GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "CONNECT", "TRACE":
	default:
		log.Printf("Unknown method %s in %s", config.Method, command)
		return nil, false
	}
	return []interface{}{config, ctx}, true
}

func ProxyNetRun(data []interface{}) bool {
	config := data[0].(*ProxyNetConfig)
	var ctx *dvcontext.RequestContext = nil
	if data[1] != nil {
		ctx = data[1].(*dvcontext.RequestContext)
	}
	return ProxyNetRunByConfig(config, ctx)
}

func ProxyNetRunByConfig(config *ProxyNetConfig, ctx *dvcontext.RequestContext) bool {
	options := dvnet.GetAveragePersistentOptions()
	var body string
	var ok bool
	if config.NotProxyBody {
		body, ok = smartNetProcessBody(config.Body, ctx)
		if !ok {
			return true
		}
	} else {
		body = ctx.PrimaryContextEnvironment.GetString(dvcontext.BODY_STRING)
	}
	method := config.Method
	if method == "" {
		method = ctx.PrimaryContextEnvironment.GetString(dvcontext.REQUEST_METHOD)
	}
	headers := make(map[string]string)
	if !config.NotProxyHeaders {
		smartProvideProxyHeaders(headers, ctx)
	}
	headers = smartNetProcessHeaders(headers, config.Headers, method, body)
	url := config.Url
	if !config.NotAddUrlPath {
		url = smartNetProcessUrl(url, ctx)
	}
	if !config.NotAddUrlPath {
		url = smartNetAddUrlParams(url, ctx)
	}
	res, err, heads := NetRequest(method, url, body, headers, options)
	if err != nil {
		log.Println(res)
		log.Printf("%s %s failed: %v", method, url, err)
		return false
	}
	var result interface{}
	switch config.ContentType {
	case "json":
		result, err = dvjson.JsonFullParser(res)
		if err != nil {
			result = string(res)
		}
	default:
		result = string(res)
	}
	if !config.NotReturnHeaders && ctx.PrimaryContextEnvironment != nil {
		SaveHeaderResult("RESPONSE", heads, ctx.PrimaryContextEnvironment)
		ctx.PrimaryContextEnvironment.Set(dvcontext.AUTO_HEADER_SET_BY, dvcontext.HEADERS_RESPONSE)
	}
	SaveActionResult(config.Result, result, ctx)
	return true
}

func smartProvideProxyHeaders(headers map[string]string, ctx *dvcontext.RequestContext) {
	origHeaders := ctx.Reader.Header
	for k, v := range origHeaders {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
}

func smartNetProcessUrl(url string, ctx *dvcontext.RequestContext) string {
	n := len(url) - 1
	if n >= 0 && (url[n] == '/' || url[n] == '\\') {
		url = url[:n]
	}
	env := GetEnvironment(ctx)
	s := env.GetString(dvcontext.REQUEST_URI)
	if s != "" {
		url += s
	}
	return url
}

func smartNetAddUrlParams(url string, ctx *dvcontext.RequestContext) string {
	env := GetEnvironment(ctx)
	s := env.GetString(dvcontext.REQUEST_URL_PARAMS)
	if s != "" {
		url += "?" + s
	}
	return url
}
