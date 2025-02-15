// Package dvconfig manages configuration for http server
// MicroCore Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
package dvconfig

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dobryvechir/microcore/pkg/dvsession"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"github.com/Dobryvechir/microcore/pkg/dvurl"

	"github.com/Dobryvechir/microcore/pkg/dvcom"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvmodules"
	"github.com/Dobryvechir/microcore/pkg/dvprocessors"
	"github.com/Dobryvechir/microcore/pkg/dvproviders"
)

func prepareComRewriteMap(rewrites []DvRewrite) dvcontext.RewriteMap {
	var rewriteResult = make(dvcontext.RewriteMap)
	for _, c := range rewrites {
		url := c.From
		if url != "" && url[0] == '/' {
			url = url[1:]
		}
		urls := strings.Split(url, "/")
		full := true
		k := urls[0]
		url = "/" + url
		if strings.HasSuffix(url, "*") {
			full = false
			url = url[:len(url)-1]
		}
		rewriteResult[k] = append(rewriteResult[k], dvcom.GetRewriteMapItem(url, full, c.To))
	}
	return rewriteResult
}

func prepareMapOfStringArrays(smap map[string]string) map[string][]string {
	r := make(map[string][]string)
	for k, v := range smap {
		r[k] = []string{v}
	}
	return r
}

func getLogLevelCode(level string) int {
	level = strings.ToLower(level)
	switch level {
	case "url":
		return dvcontext.LogLevelInputURL
	case "input":
		return dvcontext.LogLevelFullInput
	case "all":
		return dvcontext.LogLevelFull
	}
	return dvcontext.LogLevelNone
}

func prepareProxyInfo(servers []ProxyServerInfo) []*dvcontext.ProxyServerBlock {
	n := len(servers)
	if n == 0 {
		return nil
	}
	res := make([]*dvcontext.ProxyServerBlock, 0, n)
	for i := 0; i < n; i++ {
		url := dvcontext.GetPurePath(servers[i].Url)
		if url == "" || len(servers[i].Filter) == 0 {
			continue
		}
		filters := dvurl.PreparseMaskExpressions(servers[i].Filter)
		if len(filters) > 0 {
			pos := 0
			rewrite := servers[i].Rewrite
			if len(rewrite) > 0 {
				pos = strings.Index(servers[i].Filter, "*")
				if pos < 0 {
					pos = len(servers[i].Filter)
				}
			}
			res = append(res, &dvcontext.ProxyServerBlock{ServerUrl: url, FilterUrls: filters, RewritePoint: pos, RewritePrefix: rewrite})
		}
	}
	return res
}

func prepareMicroCoreInfo(server *DvHostServer) *dvcontext.MicroCoreInfo {
	server.AccessControlMaxAge = strings.TrimSpace(server.AccessControlMaxAge)
	if server.AccessControlMaxAge != "" {
		if _, err := strconv.Atoi(server.AccessControlMaxAge); err != nil {
			log.Printf("accessControlMaxAge is expected to be either empty or number but it is %s \n", server.AccessControlMaxAge)
			server.AccessControlMaxAge = ""
		}
	}
	server.AccessControlAllowCredentials = strings.TrimSpace(server.AccessControlAllowCredentials)
	if server.AccessControlAllowCredentials != "" {
		if server.AccessControlAllowCredentials == "false" {
			server.AccessControlAllowCredentials = ""
		} else if server.AccessControlAllowCredentials != "true" {
			log.Printf("accessControlAllowCredentials is expected to be either empty or true or false but it is %s \n", server.AccessControlAllowCredentials)
			server.AccessControlAllowCredentials = ""
		}
	}
	dvServerInfo := &dvcontext.MicroCoreInfo{
		BaseFolderUrl:             dvcontext.GetPurePath(server.BaseFolder),
		ProxyServerUrl:            dvcontext.GetPurePath(server.DefaultProxyServer),
		ProxyServerSettings:       server.ProxyServerSettings,
		ProxyServerMap:            server.ProxyServerMap,
		ProxyName:                 dvcom.PrepareProxyName(server.ProxyName),
		DomainName:                dvcom.PrepareProxyName(server.DomainName),
		ProxyServers:              prepareProxyInfo(server.ProxyServers),
		BaseRewrite:               prepareComRewriteMap(server.Rewrites),
		ServerRewrite:             prepareComRewriteMap(server.ServerRewrites),
		HeadersReplace:            prepareMapOfStringArrays(server.HeadersReplace),
		HeadersStatic:             prepareMapOfStringArrays(server.HeadersStatic),
		HeadersStaticOptions:      prepareMapOfStringArrays(server.HeadersStaticOptions),
		HeadersProxyServer:        prepareMapOfStringArrays(server.HeadersProxyServer),
		HeadersProxyServerOptions: prepareMapOfStringArrays(server.HeadersProxyServerOptions),
		HeadersSpecial:            make(map[string]dvcontext.MicroCoreHeaderAttribute),
		HeadersSpecialOptions:     make(map[string]dvcontext.MicroCoreHeaderAttribute),
		HeadersSpecialStatic:      make(map[string]dvcontext.MicroCoreHeaderAttribute),
		DirectoryIndex:            dvtextutils.ConvertToNonEmptyList(server.DirectoryIndex),
		BaseProcessorBlocks:       dvprocessors.InitializeProcessors(server.Processors),
		BaseProviderBlocks:        dvproviders.MakeProviderBlocks(server.Providers),
		PostProcessorBlocks:       dvprocessors.InitializePostProcessors(server.PostProcessors),
		HostHeader:                strings.TrimSpace(server.HostHeader),
		LogLevel:                  getLogLevelCode(server.LogLevel),
		ActionPolicies:            server.ActionPolicies,
		SecurityInfo:              server.SecurityInfo,
		UrlIncludeOption:          server.UrlIncludeOption,
	}

	accessControlAllowOrigin := dvcom.PrepareAccessControlLists(server.AccessControlAllowOrigin)
	accessControlAllowMethods := dvcom.PrepareAccessControlLists(server.AccessControlAllowMethod)
	accessControlAllowHeaders := dvcom.PrepareAccessControlLists(server.AccessControlAllowHeaders)
	accessControlExposeHeaders := dvcom.PrepareAccessControlLists(server.AccessControlExposeHeaders)
	if server.AccessControlMaxAge != "" {
		dvServerInfo.HeadersStaticOptions["Access-Control-Max-Age"] = []string{server.AccessControlMaxAge}
		dvServerInfo.HeadersProxyServerOptions["Access-Control-Max-Age"] = []string{server.AccessControlMaxAge}
	}
	if server.AccessControlAllowCredentials != "" {
		dvServerInfo.HeadersStatic["Access-Control-Allow-Credentials"] = []string{"true"}
		dvServerInfo.HeadersStaticOptions["Access-Control-Allow-Credentials"] = []string{"true"}
		dvServerInfo.HeadersProxyServer["Access-Control-Allow-Credentials"] = []string{"true"}
		dvServerInfo.HeadersProxyServerOptions["Access-Control-Allow-Credentials"] = []string{"true"}
	}
	cacheControl := strings.TrimSpace(server.CacheControl)
	if cacheControl != "" {
		dvServerInfo.HeadersStatic["Cache-Control"] = []string{cacheControl}
		dvServerInfo.HeadersProxyServer["Cache-Control"] = []string{cacheControl}
	}
	if accessControlExposeHeaders.Kind >= 0 {
		dvServerInfo.HeadersStatic["Access-Control-Expose-Headers"] = []string{accessControlExposeHeaders.Plain}
		dvServerInfo.HeadersSpecial["Access-Control-Expose-Headers"] = accessControlExposeHeaders
	}
	if accessControlAllowHeaders.Kind >= 0 {
		dvServerInfo.HeadersStaticOptions["Access-Control-Allow-Headers"] = []string{accessControlAllowHeaders.Plain}
		dvServerInfo.HeadersSpecialOptions["Access-Control-Allow-Headers"] = accessControlAllowHeaders
	}
	if accessControlAllowMethods.Kind >= 0 {
		dvServerInfo.HeadersStaticOptions["Access-Control-Allow-Methods"] = []string{}
		dvServerInfo.HeadersSpecialOptions["Access-Control-Allow-Methods"] = accessControlAllowMethods
	}
	if accessControlAllowOrigin.Kind >= 0 {
		accessControlAllowOrigin.Kind = dvcontext.HeadersSetOrigin
		if _, ok := accessControlAllowOrigin.Imap["*"]; ok {
			accessControlAllowOrigin.Kind = dvcontext.HeadersSetOriginAlways
		}
		dvServerInfo.HeadersSpecial["Access-Control-Allow-Origin"] = accessControlAllowOrigin
		dvServerInfo.HeadersSpecialOptions["Access-Control-Allow-Origin"] = accessControlAllowOrigin
		dvServerInfo.HeadersSpecialStatic["Access-Control-Allow-Origin"] = accessControlAllowOrigin
	}
	dvServerInfo.ModuleHandler = dvmodules.RegisterEndPointHandlers(server.Modules)
	dvServerInfo.ActionHandler = dvmodules.DynamicRegisterEndPointActions(server.Actions, server.DynamicActions)
	if server.Session != nil {
		dvServerInfo.Session = dvsession.GetServerSessionProvider(server.Session.Name, server.Session.Option, server.Session.Params, server.Session.Urls, server.Session.Prefix)
	}
	return dvServerInfo
}

// ProcessBaseFolder sets up the http server using information for all requested servers
func ProcessBaseFolder(server *DvHostServer, hostServers []DvHostServer) {
	if server != nil || len(hostServers) > 0 {
		defaultServerInfo := prepareMicroCoreInfo(server)
		hostServerInfo := make(map[string]*dvcontext.MicroCoreInfo)
		for i, c := range hostServers {
			hosts := strings.Split(strings.TrimSpace(c.Hosts), " ")
			if len(hosts) == 0 || hosts[0] == "" {
				log.Printf("HostsServer[%d] is omitted because not filled", i)
				continue
			}
			serverInfo := prepareMicroCoreInfo(&c)
			for _, host := range hosts {
				if _, ok := hostServerInfo[host]; ok {
					log.Printf("HostServer host %s is duplicated", host)
					continue
				}
				hostServerInfo[host] = serverInfo
			}
		}
		http.HandleFunc("/", dvcom.MakeDefaultHandler(defaultServerInfo, hostServerInfo))
	}
}
