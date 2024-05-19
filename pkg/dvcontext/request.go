/***********************************************************************
MicroCore
Copyright 2020 - 2024 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import (
	"net/http"
	"sync"

	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
)

const (
	DoNotShowPlaceInfo = "~DO NOT SHOW PLACE~"
)
const (
	UrlIncludeOption_HostFirst          = 1
	UrlIncludeOption_PredefinedHostOnly = 2
	UrlIncludeOption_HostWithoutPort    = 4
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
	ServerUrl     string
	FilterUrls    []*MaskInfo
	RewritePoint  int
	RewritePrefix string
}

type RewriteMapItem struct {
	Url    string
	UrlLen int
	Full   bool
	Src    string
}

type SecurityServerInfo struct {
	RolePrefix      string   `json:"role_prefix"`
	SuperAdminRoles []string `json:"super_admin_roles"`
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

type ParallelProcessing struct {
	IntervalTimeUnit float32 `json:"interval_time_unit"`
}

type ParallelExecutionControl struct {
	Flusher http.Flusher
	Value   interface{}
}

type ActionPolicy struct {
	LogFirstTimes int `json:"log_first_times"`
	LogNextTime   int `json:"log_next_time"`
}

type DatabaseTable struct {
	Name          string `json:"name"`
	Kind          string `json:"kind"`
	Web           string `json:"web"`
	KeyFirst      string `json:"keyFirst"`
	WebField      string `json:"webField"`
	WebFileName   string `json:"webFileName"`
	WebFormats    string `json:"webFormats"`
	AllowCustomId bool   `json:"customId"`
	Version       string `json:"version"`
}

type DatabaseConfig struct {
	Name    string           `json:"name"`
	Root    string           `json:"root"`
	WebRoot string           `json:"webroot"`
	Tables  []*DatabaseTable `json:"tables"`
}

type MicroCoreInfo struct {
	sync.RWMutex
	Client                    *http.Client
	BaseFolderUsed            bool
	ExtraStaticServer         bool
	ProxyServerHttp           bool
	HasProxyServers           bool
	BaseFolderUrl             string
	ProxyServerMap            map[string]string
	ProxyServerUrl            string
	ProxyServerSettings       ServerSettings
	ProxyName                 string
	DomainName                string
	ProxyServers              []*ProxyServerBlock
	ErrorPolicies             map[string]*RequestErrorPolicy
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
	Session                   ServerSessionProvider
	SecurityInfo              *SecurityServerInfo
	ModuleHandler             HandlerFunc
	ActionPolicies            map[string]*ActionPolicy
	ActionHandler             HandlerFunc
	LogLevel                  int
	UrlIncludeOption          int
}

type InterfaceExecutor func(*RequestContext, interface{}) interface{}

type RequestContext struct {
	Id                        int64
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
	Session                   RequestSession
	Input                     []byte
	InputStr                  string
	InputJson                 interface{}
	Output                    []byte
	Headers                   map[string][]string
	Error                     error
	Action                    *DvAction
	StatusCode                int
	ParallelExecution         *ParallelExecutionControl
	ExecutorFn                InterfaceExecutor
	LogLevel                  int
	PlaceInfo                 string
	TempFiles                 []string
}

type HandlerFunc func(request *RequestContext) bool

type RequestErrorPolicy struct {
	Name         string `json:"name"`
	Format       string `json:"format"`
	ContentType  string `json:"content_type"`
	FormatForced bool   `json:"format_forced"`
}

const (
	BODY_STRING           = "BODY_STRING"
	BODY_JSON             = "BODY_JSON"
	BODY_STRING_ARRAY_MAP = "BODY_STRING_ARRAY_MAP"
	BODY_STRING_MAP       = "BODY_STRING_MAP"
	BODY_PARAM_PREFIX     = "BODY_PARAM_"
	REQUEST_METHOD        = "REQUEST_METHOD"
	REQUEST_URI           = "REQUEST_URI"
	REQUEST_URL_PARAMS    = "REQUEST_URL_PARAMS"
	AUTO_HEADER_SET_BY    = "AUTO_HEADER_SET_BY"
	HEADERS_RESPONSE      = "HEADERS_RESPONSE"
)

var DefaultRequestErrorPolicy = &RequestErrorPolicy{
	Name:         "default",
	Format:       "{\"timestamp\":\"$$$TIMESTAMP\",\"status\":$$$STATUS,\"error\":\"$$$ERROR\",\"message\":\"$$$MESSAGE\",\"path\":\"$$$PATH\"}",
	ContentType:  "application/json",
	FormatForced: false,
}
