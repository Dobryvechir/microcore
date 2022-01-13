// package modules allows to extend the basic functionality of the server
//  thru creating various modules and registering it here.
// MicroCore Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)

package dvmodules

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvlog"
	"github.com/Dobryvechir/microcore/pkg/dvurl"
	"strings"
)

type ActionEndPointHandler func(ctx *dvcontext.RequestContext) bool

var registeredActionProcessors = make(map[string]ActionEndPointHandler)
var dynamicWaitingActions []*dvcontext.DvAction
var dynamicActionPool map[string]*dvurl.UrlPool

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

func ValidateRequest(validations []*dvcontext.ValidatePattern, environment *dvevaluation.DvObject) string {
	n := len(validations)
	for i := 0; i < n; i++ {
		validation := validations[i]
		src, err := environment.CalculateString(validation.Source)
		if err != nil {
			return "Error in expression " + err.Error()
		}
		res := validation.Match(src)
		if res != "" {
			return res
		}
	}
	return ""
}

func FireAction(action *dvcontext.DvAction, request *dvcontext.RequestContext) bool {
	request.Action = action
	proc, ok := registeredActionProcessors[action.Typ]
	if !ok {
		dvlog.Printf("Action %s url %s has incorrect type %s", action.Name, action.Url, action.Typ)
		return false
	}
	if request.Reader.Method == "OPTIONS" {
		request.HandleCommunication()
		return true
	}
	err := collectRequestParameters(request)
	if err != nil {
		dvlog.PrintlnError("Cannot load body " + err.Error())
		request.HandleInternalServerError()
		return true
	}
	if len(action.Validations) > 0 {
		res := ValidateRequest(action.Validations, request.PrimaryContextEnvironment)
		if res != "" {
			request.SetHttpErrorCode(400, res)
			return true
		}
	}
	if request.Server!=nil && request.Server.Session!=nil && action.Session!=nil {
		sessionId :=""
		if request.PrimaryContextEnvironment!=nil {
			sessionId = request.PrimaryContextEnvironment.FindFirstNotEmptyString(action.Session.Id)
		}
		request.Session, err = request.Server.Session.GetSessionStorage(request, action.Session, sessionId)
		if err!=nil {
			dvlog.PrintlnError("Cannot handle session " + err.Error())
			request.HandleInternalServerError()
			return true
		}
	}
	return proc(request)
}

func RegisterEndPointActions(actions []*dvcontext.DvAction) (dvcontext.HandlerFunc, map[string]*dvurl.UrlPool) {
	n := len(actions)
	if n == 0 {
		return nil, nil
	}
	base := make(map[string]*dvurl.UrlPool)
	for i := 0; i < n; i++ {
		action := actions[i]
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
	return getActionHandlerFunc(base), base
}

func urlActionVerifier(context interface{}, resolver *dvurl.UrlResolver, urlData *dvurl.UrlResultInfo) bool {
	requestContext := context.(*dvcontext.RequestContext)
	requestContext.SetUrlInlineParameters(urlData.UrlKeys)
	action := resolver.Handler.(*dvcontext.DvAction)
	return FireAction(action, requestContext)
}

func getActionHandlerFunc(base map[string]*dvurl.UrlPool) dvcontext.HandlerFunc {
	return func(context *dvcontext.RequestContext) bool {
		method := strings.ToUpper(context.Reader.Method)
		urlPool := base[method]
		if urlPool == nil {
			return false
		}
		urls := context.Urls
		ok, _ := dvurl.UrlSearch(context, urlPool, urls, urlActionVerifier, context.PrimaryContextEnvironment)
		return ok
	}
}

func DynamicRegisterEndPointActions(actions []*dvcontext.DvAction, isDynamic bool) dvcontext.HandlerFunc {
	if isDynamic && dynamicWaitingActions != nil {
		actions = append(actions, dynamicWaitingActions...)
	}
	fn, pool := RegisterEndPointActions(actions)
	if isDynamic {
		dynamicActionPool = pool
	}
	return fn
}

func AddDynamicAction(action *dvcontext.DvAction, env *dvevaluation.DvObject) {
	if action.Method == "" {
		action.Method = "GET"
	} else {
		action.Method = strings.ToUpper(action.Method)
	}
	if dynamicActionPool == nil {
		if dynamicWaitingActions == nil {
			dynamicWaitingActions = make([]*dvcontext.DvAction, 1, 128)
			dynamicWaitingActions[0] = action
		} else {
			dynamicWaitingActions = append(dynamicWaitingActions, action)
		}
	} else {
		urlPool := dynamicActionPool[action.Method]
		if urlPool == nil {
			urlPool = dvurl.GetUrlHandler()
			dynamicActionPool[action.Method] = urlPool
			urlPool.RegisterHandlerFunc(action.Url, action)
		} else {
			ok, res := dvurl.SingleSimplifiedUrlSearch(urlPool, action.Url)
			if !ok {
				urlPool.RegisterHandlerFunc(action.Url, action)
			} else {
				currentAction := res.CustomObject.(*dvcontext.DvAction)
				currentAction.CloneFrom(action)
			}
		}
	}
}
