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
	ExtraServerFile           bool
	ExtraServerHttp           bool
	BaseFolderUrl             string
	ExtraServerUrl            string
	ExtraServerSettings       ServerSettings
	ProxyName                 string
	HeadersStatic             map[string][]string
	HeadersExtraServer        map[string][]string
	HeadersStaticOptions      map[string][]string
	HeadersExtraServerOptions map[string][]string
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
}

type RequestContext struct {
	Extra           map[string]interface{}
	ExtraAsDvObject *dvevaluation.DvObject
	Params          map[string]string
	UrlInlineParams map[string]string
	Url             string
	Urls            []string
	UrlsLowerCase   []string
	FileName        string
	DataType        string
	Writer          http.ResponseWriter
	Reader          *http.Request
	Server          *MicroCoreInfo
	Output          []byte
	Error           error
	Action          *DvAction
	StatusCode      int
}

type HandlerFunc func(request *RequestContext) bool
