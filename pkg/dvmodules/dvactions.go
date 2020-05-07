/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvmodules

import (
	"github.com/Dobryvechir/microcore/pkg/dvmeta"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
	"strings"
)

type ActionEndPointHandler func(request *dvmeta.RequestContext) bool

var registeredActionProcessors = make(map[string]ActionEndPointHandler)

func RegisterActionProcessor(name string, proc ActionEndPointHandler, silent bool) bool {
	if _, ok := registeredActionProcessors[name]; ok {
		if silent {
			return false
		}
		panic("Processor with name " + name + " already registered")
	}
	registeredActionProcessors[name] = proc
	return true
}

func FireAction(action *dvmeta.DvAction, request *dvmeta.RequestContext) bool {
	request.Action = action
	proc, ok := registeredActionProcessors[action.Typ]
	if !ok {
		return false
	}
	return proc(request)
}

func RegisterEndPointActions(actions []dvmeta.DvAction) dvmeta.HandlerFunc {
	if len(actions) == 0 {
		return nil
	}
	base := make(map[string]*dvurl.UrlPool)
	for _, action := range actions {
		method := strings.ToUpper(strings.TrimSpace(action.Method))
		if method == "" {
			method = "GET"
		}
		pool := base[method]
		if pool == nil {
			pool = dvurl.GetUrlHandler()
			base[method] = pool
		}
		pool.RegisterHandlerFunc(action.Url, action)
	}
	return getActionHandlerFunc(base)
}

func urlActionVerifier(context interface{}, resolver *dvurl.UrlResolver, urlData *dvurl.UrlResultInfo) bool {
	requestContext := context.(*dvmeta.RequestContext)
	requestContext.UrlInlineParams = urlData.UrlKeys
	action := resolver.Handler.(*dvmeta.DvAction)
	return FireAction(action, requestContext)
}

func getActionHandlerFunc(base map[string]*dvurl.UrlPool) dvmeta.HandlerFunc {
	return func(context *dvmeta.RequestContext) bool {
		method := strings.ToUpper(context.Reader.Method)
		urlPool := base[method]
		if urlPool == nil {
			return false
		}
		urls := context.Urls
		ok, _ := dvurl.UrlSearch(context, urlPool, urls, urlActionVerifier, context.ExtraAsDvObject)
		return ok
	}
}
