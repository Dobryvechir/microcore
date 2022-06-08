/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import (
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	HeadersAddToList       = iota
	HeadersSetOrigin       = iota
	HeadersSetOriginAlways = iota
)

const (
	LogLevelNone      = iota
	LogLevelInputURL  = iota
	LogLevelFullInput = iota
	LogLevelFull      = iota
)

const LogHandled = -2000

var isRecord = false

var recordNo int
var recordPath string
var uniqueId int64 = 0
var idMutex sync.Mutex

func (request *RequestContext) HandleCommunication() {
	if request.LogLevel == LogHandled || request.Reader == nil {
		return
	}
	if request.DataType == "" {
		request.DataType = "application/json"
	}
	postHeaders := request.Headers
	if postHeaders == nil {
		postHeaders = make(map[string][]string)
	}
	if _, ok := postHeaders["Content-Type"]; !ok {
		postHeaders["Content-Type"] = []string{request.DataType}
	}
	preHeaders := request.Server.HeadersStatic
	if request.Reader.Method == "OPTIONS" {
		preHeaders = request.Server.HeadersStaticOptions
	}
	origin := "*"
	if korig, isok := request.Reader.Header["Origin"]; isok && len(korig) > 0 {
		origin = korig[0]
	}
	ProvideHeaders(preHeaders, postHeaders, origin, request.Server.HeadersSpecialStatic, request.Writer)
	if request.StatusCode > 0 {
		request.Writer.WriteHeader(request.StatusCode)
	}
	if request.Error != nil {
		if len(request.Output) == 0 {
			request.Output = []byte(request.Error.Error())
		}
	}
	Send(request.Writer, request.Reader, request.Output)
	if request.PlaceInfo == "" {
		s := strconv.Itoa(request.StatusCode)
		if request.Action != nil && len(request.Action.LogPolicy) != 0 && request.Server != nil && request.Server.ActionPolicies != nil {
			policyType := request.Server.ActionPolicies[request.Action.LogPolicy]
			if policyType != nil && (policyType.LogFirstTimes > 0 || policyType.LogNextTime > 0) {
				request.PlaceInfo = getRequestLogInfo(s, policyType, request.Reader.Method, request.Url)
			} else {
				request.PlaceInfo = s
			}
		} else {
			request.PlaceInfo = s
		}
	}
	request.LogLevel = LogHandled
}

func HandlerError(w http.ResponseWriter, r *http.Request, message string) {
	log.Printf("Error %s: %s [%s]", message, r.URL.Path, r.Method)
	errCode := 0
	if len(message) >= 3 && message[0] >= '1' && message[0] <= '5' && message[1] >= '0' && message[1] <= '9' && message[2] >= '0' && message[2] <= '9' {
		errCode, _ = strconv.Atoi(message[:3])
	}
	if errCode > 0 {
		w.WriteHeader(errCode)
	} else {
		message = "403 - " + message
		w.WriteHeader(http.StatusForbidden)
	}
	HandlerOk(w, r, message)
}

func HandleError(request *RequestContext, message string) {
	HandlerError(request.Writer, request.Reader, message)
}

func HandlerOk(w http.ResponseWriter, r *http.Request, message string) {
	Send(w, r, []byte(message))
}

func Send(w http.ResponseWriter, r *http.Request, message []byte) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if message != nil {
		w.Write(message)
	}
	if isRecord {
		recordNo++
		name := recordPath + strconv.Itoa(recordNo) + "_" + getSystemSafeFileName(r.URL.Path) + "_"
		saveRequest(name+"r.txt", r)
		saveResponse(name+"s.txt", w, message)
	}
}

func saveRequest(filename string, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		body = []byte("!!!Error reading body!!!!!")
	}
	f, er := os.Create(filename)
	if er != nil {
		log.Printf("Error creating file:%s %v", filename, er)
		return
	}
	defer f.Close()
	s := r.Proto + " " + r.URL.Path + "\n"
	for k, v := range r.Header {
		s += k + ": " + presentStringArray(v) + "\n"
	}
	s += "\n"
	_, e1 := f.WriteString(s)
	if e1 != nil {
		log.Printf("Error writing file:%s %v", filename, e1)
		return
	}
	f.Write(body)
	f.Sync()
}

func saveResponse(filename string, w http.ResponseWriter, body []byte) {
	f, er := os.Create(filename)
	if er != nil {
		log.Printf("Error creating file:%s %v", filename, er)
		return
	}
	defer f.Close()
	s := ""
	for k, v := range w.Header() {
		s += k + ": " + presentStringArray(v) + "\n"
	}
	s += "\n"
	_, e1 := f.WriteString(s)
	if e1 != nil {
		log.Printf("Error writing file:%s %v", filename, e1)
		return
	}
	f.Write(body)
	f.Sync()
}

func ProvideHeaders(preHeaders map[string][]string, postHeaders map[string][]string, origin string, specialHeaders map[string]MicroCoreHeaderAttribute, w http.ResponseWriter) {
	for nm, hd := range preHeaders {
		if _, ok := specialHeaders[nm]; ok {
			continue
		}
		if _, okey := postHeaders[nm]; okey {
			continue
		}
		for _, h := range hd {
			w.Header().Set(nm, h)
		}
	}
	for nm, hd := range postHeaders {
		if _, ok := specialHeaders[nm]; ok {
			continue
		}
		for _, h := range hd {
			w.Header().Set(nm, h)
		}
	}
	for nm, hd := range specialHeaders {
		pre, ok := preHeaders[nm]
		post, okey := postHeaders[nm]
		if !okey && ok {
			post = pre
			okey = true
		}
		switch hd.Kind {
		case HeadersAddToList:
			if okey && len(post) > 0 {
				s := strings.TrimSpace(post[0])
				oldList := dvtextutils.ConvertToNonEmptyList(s)
				s = dvtextutils.AddNonRepeatingWords(s, oldList, hd.List, hd.Imap, hd.Plain, ", ")
				w.Header().Set(nm, s)
			} else {
				w.Header().Set(nm, hd.Plain)
			}
		case HeadersSetOrigin:
			if _, ok = hd.Imap[origin]; ok {
				w.Header().Set(nm, origin)
			} else if okey && len(post) > 0 {
				w.Header().Set(nm, post[0])
			}
		case HeadersSetOriginAlways:
			w.Header().Set(nm, origin)
		}
	}
}

func getSystemSafeFileName(name string) string {
	var a [100]byte
	len := copy(a[:], name)
	for i, c := range a {
		if !(c >= 48 && c < 58 || c >= 65 && c < 91 || c >= 97 && c < 123) {
			a[i] = 95
		}
	}
	return string(a[:len])
}

//TODO: not effective, need to be rewritten
func presentStringArray(dat []string) string {
	s := "["
	if dat != nil {
		for k, v := range dat {
			if k != 0 {
				s += "]["
			}
			s += v
		}
	}
	return s + "]"
}

func GetPurePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	lastSym := path[len(path)-1]
	if lastSym == '/' || lastSym == '\\' {
		path = path[:len(path)-1]
	}
	return path
}

func SetRecordMode(path string) {
	path = GetPurePath(path)
	if path == "" {
		return
	}
	if _, err := os.Stat(path); err != nil {
		log.Printf("Record path is not correct %s\n", err.Error())
		return
	}
	isRecord = true
	recordNo = 0
	recordPath = path + "/_" + strconv.FormatInt(dvlog.StartTime, 36) + "_"
	log.Printf("Recording started at %v\n", recordPath)
}

func (ctx *RequestContext) HandleInternalServerError() {
	if ctx.StatusCode < 400 {
		ctx.StatusCode = 500
	}
	if ctx.Reader == nil {
		dvlog.PrintlnError("Exit by critical error")
		os.Exit(1)
	}
	ctx.HandleCommunication()
}

func (ctx *RequestContext) HandleHttpError(statusCode int) {
	ctx.StatusCode = statusCode
	ctx.HandleCommunication()
}

func (ctx *RequestContext) HandleFileNotFound() {
	ctx.StatusCode = 404
	ctx.HandleCommunication()
}

func GetUniqueId() int64 {
	idMutex.Lock()
	v := uniqueId
	uniqueId++
	idMutex.Unlock()
	return v
}
