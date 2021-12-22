/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import (
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"net/http"
	"sync"
)

type ProcessorEndPointHandler func(request *RequestContext) bool

type ProcessorBlock struct {
	Name            string
	EndPointHandler ProcessorEndPointHandler
	Urls            []*MaskInfo
	Data            map[string]string
}

type MicroCoreHeaderAttribute struct {
	Kind  int
	Imap  map[string]int
	List  []string
	Plain string
}

type ProxyServerBlock struct {
	ServerUrl  string
	FilterUrls []*MaskInfo
}

type RewriteMapItem struct {
	Url    string
	UrlLen int
	Full   bool
	Src    string
}

type RewriteMap map[string][]*RewriteMapItem

type ServerSettings struct {
	MaxIdleConns          int  `json:"maxIdleConnections"`
	IdleConnTimeout       int  `json:"idleConnectionTimeout"`
	DisableCompression    bool `json:"disableCompression"`
	DisableKeepAlives     bool `json:"disableKeepAlives"`
	MaxIdleConnsPerHost   int  `json:"maxIdleConnectionsPerHost"`
	ResponseHeaderTimeout int  `json:"responseHeaderTimeout"`
	ExpectContinueTimeout int  `json:"expectContinueTimeout"`
}

type MicroCoreInfo struct {
	sync.RWMutex
	Client                    *http.Client
	BaseFolderUsed            bool
	ExtraStaticServer         bool
	ProxyServerHttp           bool
	HasProxyServers           bool
	BaseFolderUrl             string
	ProxyServerUrl            string
	ProxyServerSettings       ServerSettings
	ProxyName                 string
	DomainName                string
	ProxyServers              []*ProxyServerBlock
	ErrorPolicies	          map[string]*RequestErrorPolicy
	HeadersStatic             map[string][]string
	HeadersProxyServer        map[string][]string
	HeadersStaticOptions      map[string][]string
	HeadersProxyServerOptions map[string][]string
	HeadersSpecial            map[string]MicroCoreHeaderAttribute
	HeadersSpecialOptions     map[string]MicroCoreHeaderAttribute
	HeadersSpecialStatic      map[string]MicroCoreHeaderAttribute
	DirectoryIndex            []string
	BaseRewrite               RewriteMap
	BaseProcessorBlocks       []ProcessorBlock
	BaseProviderBlocks        []ProcessorBlock
	PostProcessorBlocks       []ProcessorBlock
	HostHeader                string
	ServerRewrite             RewriteMap
	ModuleHandler             HandlerFunc
	ActionHandler             HandlerFunc
	LogLevel                  int
}

type RequestContext struct {
	Extra                     map[string]interface{}
	PrimaryContextEnvironment *dvevaluation.DvObject
	LocalContextEnvironment   *dvevaluation.DvObject
	Params                    map[string]string
	UrlInlineParams           map[string]string
	Url                       string
	Urls                      []string
	UrlsLowerCase             []string
	FileName                  string
	DataType                  string
	Queries                   map[string]string
	Writer                    http.ResponseWriter
	Reader                    *http.Request
	Server                    *MicroCoreInfo
	Input                     []byte
	InputStr                  string
	InputJson                 interface{}
	Output                    []byte
	Headers                   map[string][]string
	Error                     error
	Action                    *DvAction
	StatusCode                int
}

type HandlerFunc func(request *RequestContext) bool

type RequestErrorPolicy struct {
	Name         string `json:"name"`
	Format       string `json:"format"`
	ContentType  string `json:"content_type"`
	FormatForced bool   `json:"format_forced"`
}

const (
	BODY_STRING        = "BODY_STRING"
	BODY_JSON          = "BODY_JSON"
	REQUEST_METHOD     = "REQUEST_METHOD"
	REQUEST_URI        = "REQUEST_URI"
	REQUEST_URL_PARAMS = "REQUEST_URL_PARAMS"
	AUTO_HEADER_SET_BY = "AUTO_HEADER_SET_BY"
	HEADERS_RESPONSE   = "HEADERS_RESPONSE"
)

var DefaultRequestErrorPolicy = &RequestErrorPolicy{
	Name:         "default",
	Format:       "{\"timestamp\":\"$$$TIMESTAMP\",\"status\":$$$STATUS,\"error\":\"$$$ERROR\",\"message\":\"$$$MESSAGE\",\"path\":\"$$$PATH\"}",
	ContentType:  "application/json",
	FormatForced: false,
}
