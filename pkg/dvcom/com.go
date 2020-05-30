/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvcom

import (
	"bytes"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var LogServer bool
var LogFileServer bool
var LogHosts bool
var defaultDirectoryIndex = []string{"$.html", "index.html", "index.htm"}
var ServiceFolder = "..$$$"

func SetRequestUrl(request *dvcontext.RequestContext, url string) {
	request.Url = url
	request.UrlsLowerCase = dvparser.ConvertURLToList(strings.ToLower(url))
	request.Urls = dvparser.ConvertURLToList(url)
}

func CheckProcessorBlocks(blocks []dvcontext.ProcessorBlock, request *dvcontext.RequestContext) bool {
	n := len(blocks)
	urls := request.Urls
	for i := 0; i < n; i++ {
		if dvurl.MatchMasksForUrlParts(blocks[i].Urls, urls, request.ExtraAsDvObject) {
			request.Params = blocks[i].Data
			if blocks[i].EndPointHandler(request) {
				return true
			}
		}
	}
	return false
}

func GetRewriteMapItem(url string, full bool, src string) *dvcontext.RewriteMapItem {
	return &dvcontext.RewriteMapItem{Url: url, UrlLen: len(url), Full: full, Src: src}
}

func getMimeByExtension(ext string, defExtension string, defMime string) string {
	value, ok := MimeTypes[ext]
	if ok {
		return value
	}
	if defMime != "" {
		return defMime
	}
	value, ok = MimeTypes[defExtension]
	if ok {
		return value
	}
	return ""
}

func getContentTypeByName(name string, defExtension string, defMime string) string {
	lst := strings.LastIndex(name, ".")
	ext := defExtension
	if lst >= 0 {
		ext = name[lst+1:]
	}
	if ext == "" {
		ext = defExtension
	}
	ext = strings.ToLower(ext)
	return getMimeByExtension(ext, defExtension, defMime) + "; charset=utf-8"
}

func GetContentTypeByFileName(name string) string {
	return getContentTypeByName(name, "", "text/plain")
}

func HandleFromFile(request *dvcontext.RequestContext) {
	if request.DataType == "" {
		request.DataType = GetContentTypeByFileName(request.FileName)
	}
	request.Output, request.Error = ioutil.ReadFile(request.FileName)
	request.HandleCommunication()
}

func HandleFromFileWithProcessorCheck(request *dvcontext.RequestContext) {
	method := request.Reader.Method
	toLog := LogFileServer && dvlog.CurrentLogLevel >= dvlog.LogInfo && (method != "OPTIONS" || dvlog.CurrentLogLevel >= dvlog.LogDetail)
	if toLog {
		bodyIo := request.Reader.Body
		var err error
		var body []byte
		if bodyIo != nil {
			body, err = ioutil.ReadAll(bodyIo)
			if err != nil {
				message := err.Error()
				log.Print(message)
				body = []byte(message)
			}
		} else {
			body = []byte{}
		}
		dvlog.WriteRequestToLog(body, request.Reader)
	}
	if request.Server.BaseProcessorBlocks != nil &&
		CheckProcessorBlocks(request.Server.BaseProcessorBlocks, request) {
		return
	}
	HandleFromFile(request)
}

func HandleFromString(request *dvcontext.RequestContext, data string) {
	request.Output = []byte(data)
	request.HandleCommunication()
}

func HandlerWriteDirect(request *dvcontext.RequestContext) {
	if request.StatusCode <= 0 {
		request.StatusCode = 200
	}
	request.Writer.WriteHeader(request.StatusCode)
	if request.Output != nil {
		res, err := request.Writer.Write(request.Output)
		if err != nil {
			log.Printf("failed to write output to request: %s: %d %v", string(request.Output), res, err)
		}
	}
}

func tryLocalFileByMethodAndRequest(request *dvcontext.RequestContext, folder string) bool {
	//TODO if OPTIONS return good options (scanning folders for methods allowed),
	// otherwise look at /{{METHOD}} folder
	//
	return false
}

func tryLocalFile(request *dvcontext.RequestContext, folder string) bool {
	name := dvcontext.GetPurePath(folder + request.Url)
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		for _, c := range request.Server.DirectoryIndex {
			request.FileName = name + "/" + c
			fi, err = os.Stat(request.FileName)
			if err == nil && fi.Mode().IsRegular() {
				request.Url = request.Url + "/" + c
				HandleFromFileWithProcessorCheck(request)
				return true
			}
		}
		c := name[len(name)-1]
		if c == '/' || c == '\\' {
			name += ServiceFolder
		} else {
			name += "/" + ServiceFolder
		}
		fi, err = os.Stat(name)
		if err == nil && fi.IsDir() {
			return tryLocalFileByMethodAndRequest(request, name+"/")
		}

	case mode.IsRegular():
		request.FileName = name
		HandleFromFileWithProcessorCheck(request)
		return true
	}
	return false
}

func createClientBySettings(settings *dvcontext.ServerSettings) *http.Client {
	tr := &http.Transport{
		DisableKeepAlives:     settings.DisableKeepAlives,
		MaxIdleConnsPerHost:   settings.MaxIdleConnsPerHost,
		ResponseHeaderTimeout: time.Duration(time.Duration(settings.ResponseHeaderTimeout) * time.Second),
		ExpectContinueTimeout: time.Duration(time.Duration(settings.ExpectContinueTimeout) * time.Second),
		MaxIdleConns:          settings.MaxIdleConns,
		IdleConnTimeout:       time.Duration(time.Duration(settings.IdleConnTimeout) * time.Second),
		DisableCompression:    settings.DisableCompression,
	}
	return &http.Client{Transport: tr}
}

func createClientForMicroCoreInfo(server *dvcontext.MicroCoreInfo) {
	server.Lock()
	if server.Client != nil {
		server.Unlock()
		return
	}
	server.Client = createClientBySettings(&server.ExtraServerSettings)
	server.Unlock()
}

func copyNonEmptyHeaders(src map[string][]string) (dst map[string][]string) {
	dst = make(map[string][]string)
	for k, v := range src {
		if k != "" && len(v) > 0 {
			dst[k] = v
		}
	}
	return
}

func extractHostFromUrl(url string) string {
	p := strings.Index(url, "//")
	if p >= 0 {
		url = url[p+2:]
	}
	p = strings.Index(url, "/")
	if p > 0 {
		url = url[:p]
	}
	return url
}

func tryHttpForward(request *dvcontext.RequestContext, url string) bool {
	if request.Server.Client == nil {
		createClientForMicroCoreInfo(request.Server)
	}
	method := request.Reader.Method
	var logFile string
	var body []byte
	var bodyIo = request.Reader.Body
	toLog := LogServer && dvlog.CurrentLogLevel >= dvlog.LogInfo && (method != "OPTIONS" || dvlog.CurrentLogLevel >= dvlog.LogDetail)
	if toLog {
		body, err := ioutil.ReadAll(bodyIo)
		if err != nil {
			message := err.Error()
			log.Print(message)
			request.Output = dvlog.FormErrorMessage(message)
			request.HandleInternalServerError()
			return false
		}
		if len(body) == 0 {
			bodyIo = nil
		} else {
			bodyIo = ioutil.NopCloser(bytes.NewReader(body))
		}
		logFile = dvlog.WriteRequestToLog(body, request.Reader)
	}
	finalUrl := url + request.Reader.URL.Path
	if request.Reader.URL.RawQuery != "" {
		finalUrl += "?" + request.Reader.URL.RawQuery
	}
	req, err := http.NewRequest(method, finalUrl, bodyIo)
	if err != nil {
		if dvlog.CurrentLogLevel >= dvlog.LogError {
			log.Printf("Error making request %s: %s", finalUrl, err.Error())
		}
		request.Output = dvlog.FormErrorMessage(err.Error())
		request.StatusCode = 400
		request.HandleCommunication()
		return false
	}
	origin := "*"
	if korig, isok := request.Reader.Header["Origin"]; isok && len(korig) > 0 {
		origin = korig[0]
	}
	req.Header = copyNonEmptyHeaders(request.Reader.Header)
	hostHeader := request.Server.HostHeader
	if hostHeader != "-" {
		if hostHeader == "" {
			hostHeader = extractHostFromUrl(url)
		}
		req.Header["Host"] = []string{hostHeader}
	}
	req.Host = ""
	if request.Server.ProxyName != "" {
		name := request.Server.ProxyName
		if len(req.Header["Origin"]) > 0 {
			origin := req.Header["Origin"][0]
			p := strings.Index(origin, "//")
			if p > 0 {
				origin = origin[:p+2] + name
			} else {
				origin = "http://" + name
			}
			req.Header["Origin"] = []string{origin}
		}
		if len(req.Header["Referer"]) > 0 {
			referer := req.Header["Referer"][0]
			proxyReferer := ""
			p := strings.Index(referer, "//") + 2
			n := strings.IndexByte(referer[p:], '/') + p
			if p > 2 {
				proxyReferer = referer[:p] + name
				if n > p {
					proxyReferer += referer[n:]
				}
			}
			req.Header["Referer"][0] = proxyReferer
		}
	}
	resp, err1 := request.Server.Client.Do(req)
	if err1 != nil {
		if dvlog.CurrentLogLevel >= dvlog.LogError {
			log.Printf("Error executing %s: %s", finalUrl, err1.Error())
		}
		if resp != nil && resp.StatusCode >= 300 {
			request.Output = dvlog.FormErrorMessage(err1.Error())
			request.StatusCode = resp.StatusCode
		} else {
			request.Output = dvlog.FormErrorMessage(err1.Error())
			request.StatusCode = 500
		}
		request.HandleCommunication()
		return false
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		if dvlog.CurrentLogLevel >= dvlog.LogError {
			log.Printf("Error reading body %s: %s", url+request.Reader.URL.Path, err2.Error())
		}
		request.Output = dvlog.FormErrorMessage(err2.Error())
		request.HandleInternalServerError()
		return false
	}
	if toLog {
		dvlog.WriteResponseToLog(logFile, resp, body)
	}
	if method != "OPTIONS" {
		dvcontext.ProvideHeaders(resp.Header, request.Server.HeadersExtraServer, origin, request.Server.HeadersSpecial, request.Writer)
	} else {
		dvcontext.ProvideHeaders(resp.Header, request.Server.HeadersExtraServerOptions, origin, request.Server.HeadersSpecialOptions, request.Writer)
	}
	request.StatusCode = resp.StatusCode
	request.Output = body
	if method != "OPTIONS" && request.Server.PostProcessorBlocks != nil {
		oldLen := len(body)
		if CheckProcessorBlocks(request.Server.PostProcessorBlocks, request) {
			request.HandleFileNotFound()
			return false
		}
		if len(request.Output) != oldLen {
			request.Writer.Header().Del("Content-Length")
		}
	}
	HandlerWriteDirect(request)
	return true
}

func PrepareProxyName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return name
	}
	p := strings.Index(name, "//")
	if p >= 0 {
		name = name[p+2:]
	}
	p = strings.IndexByte(name, '/')
	if p > 0 {
		name = name[:p]
	}
	return name
}

func rewriteComRewriteItem(url string, r []*dvcontext.RewriteMapItem) string {
	l := len(url)
	for _, c := range r {
		if l < c.UrlLen {
			continue
		}
		if c.Full {
			if l == c.UrlLen && c.Url == url {
				return c.Src
			}
		} else {
			if url[:c.UrlLen] == c.Url {
				return c.Src
			}
		}
	}
	return url
}

func prepareMicroCoreInfo(serverInfo *dvcontext.MicroCoreInfo) {
	serverInfo.BaseFolderUsed = false
	if serverInfo.BaseFolderUrl != "" {
		if _, err := os.Stat(serverInfo.BaseFolderUrl); err != nil {
			log.Printf("BaseFolder path %s is not correct %v\n", serverInfo.BaseFolderUrl, err.Error())
		} else {
			serverInfo.BaseFolderUsed = true
		}
	}
	if len(serverInfo.DirectoryIndex) == 0 {
		serverInfo.DirectoryIndex = defaultDirectoryIndex
	}
	serverInfo.ExtraServerHttp = false
	serverInfo.ExtraServerFile = false
	if serverInfo.ExtraServerUrl != "" {
		if len(serverInfo.ExtraServerUrl) > 5 && (serverInfo.ExtraServerUrl[:5] == "http:" || serverInfo.ExtraServerUrl[:6] == "https:") {
			serverInfo.ExtraServerHttp = true
		} else {
			if _, err := os.Stat(serverInfo.ExtraServerUrl); err != nil {
				log.Printf("extraServer path %s is not correct %v\n", serverInfo.ExtraServerUrl, err.Error())
			} else {
				serverInfo.ExtraServerFile = true
			}
		}
	}
}

func PrepareAccessControlLists(data string) dvcontext.MicroCoreHeaderAttribute {
	res := dvcontext.MicroCoreHeaderAttribute{Kind: dvcontext.HeadersAddToList, Imap: make(map[string]int), List: dvparser.ConvertToNonEmptyList(data)}
	if len(res.List) == 0 {
		res.Kind = -1
		return res
	}
	res.Plain = strings.Join(res.List, ", ")
	for i, c := range res.List {
		res.Imap[c] = i
	}
	return res
}

func calculateRequestContextParameters(r *http.Request) (res map[string]interface{}) {
	res = make(map[string]interface{})
	headers := r.Header
	res["HTTP_ACCEPT"] = headers.Get("Accept")
	res["HTTP_COOKIE"] = headers.Get("Cookie")
	res["HTTP_FORWARDED"] = headers.Get("Forwarded")
	if res["HTTP_FORWARDED"] == "" {
		res["HTTP_FORWARDED"] = headers.Get("X-Forwarded-For")
		if res["HTTP_FORWARDED"] == "" {
			res["HTTP_FORWARDED"] = headers.Get("X-Forwarded-Host")
			if res["HTTP_FORWARDED"] == "" {
				res["HTTP_FORWARDED"] = headers.Get("X-Forwarded-Proto")
			}
		}
	}
	res["HTTP_HOST"] = headers.Get("Host")
	res["HTTP_REFERER"] = headers.Get("Referer")
	res["HTTP_USER_AGENT"] = headers.Get("User-Agent")
	res["REQUEST_METHOD"] = r.Method
	res["REQUEST_URI"] = r.URL.Path
	res["REQUEST_SCHEME"] = "http"
	res["REQUEST_URL_PARAMS"] = r.URL.RawQuery
	if r.TLS != nil {
		res["REQUEST_SCHEME"] = "https"
	}
	res["SERVER_NAME"] = r.Host
	for key, value := range headers {
		if len(value) > 0 {
			res["H_"+dvparser.ConvertToUpperAlphaDigital([]byte(key))] = value[0]
		}
	}
	m, _ := url.ParseQuery(r.URL.RawQuery)
	for key, value := range m {
		if len(value) > 0 {
			res["G_"+dvparser.ConvertToUpperAlphaDigital([]byte(key))] = value[0]
			res["GM_"+key] = value
		}
	}
	return
}

func MakeDefaultHandler(defaultServerInfo *dvcontext.MicroCoreInfo, hostServerInfo map[string]*dvcontext.MicroCoreInfo) http.HandlerFunc {
	prepareMicroCoreInfo(defaultServerInfo)
	for _, c := range hostServerInfo {
		prepareMicroCoreInfo(c)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		currentServer := defaultServerInfo
		if d, okey := hostServerInfo[r.Host]; okey {
			currentServer = d
		}
		request := &dvcontext.RequestContext{
			Extra:  calculateRequestContextParameters(r),
			Server: currentServer,
			Writer: w,
			Reader: r,
		}
		request.ExtraAsDvObject = dvparser.GetProptertiesPrototypedToGlobalProperties(request.Extra)
		SetRequestUrl(request, r.URL.Path)
		firstUrl := ""
		if len(request.Urls) >= 1 {
			firstUrl = request.Urls[0]
		}
		rew, ok := currentServer.BaseRewrite[firstUrl]
		if ok {
			SetRequestUrl(request, rewriteComRewriteItem(request.Url, rew))
		}
		if currentServer.BaseFolderUsed && tryLocalFile(request, currentServer.BaseFolderUrl) {
		} else if currentServer.ActionHandler != nil && currentServer.ActionHandler(request) {
		} else if currentServer.ModuleHandler != nil && currentServer.ModuleHandler(request) {
		} else {
			rew, ok := currentServer.ServerRewrite[firstUrl]
			if ok {
				SetRequestUrl(request, rewriteComRewriteItem(request.Url, rew))
			}
			if currentServer.ExtraServerFile && tryLocalFile(request, currentServer.ExtraServerUrl) {
			} else if currentServer.ExtraServerHttp {
				tryHttpForward(request, currentServer.ExtraServerUrl)
			} else {
				request.HandleFileNotFound()
			}
		}
	}

}
