/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvurl

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"github.com/Dobryvechir/microcore/pkg/dvtextutils"
	"log"
	"strings"
)

var logLevel = false

type UrlResolver struct {
	Handler  interface{}
	maskInfo *dvcontext.MaskInfo
	rewrite  *RewriteInfo
}

type UrlPool struct {
	dirs  map[string]*UrlPool
	files []*UrlResolver
}

type UrlVerifier func(context interface{}, resolver *UrlResolver, urlResultInfo *UrlResultInfo) bool

func GetUrlHandler() *UrlPool {
	return &UrlPool{dirs: make(map[string]*UrlPool), files: make([]*UrlResolver, 0, 16)}
}

func placeUrlResolverToHierarchy(urlPool *UrlPool, keys []string, start int, resolver *UrlResolver) {
	n := len(keys)
	if start < n {
		key := keys[start]
		nextUrlPool, ok := urlPool.dirs[key]
		if !ok {
			nextUrlPool = GetUrlHandler()
			urlPool.dirs[key] = nextUrlPool
		}
		placeUrlResolverToHierarchy(nextUrlPool, keys, start+1, resolver)
	} else {
		urlPool.files = append(urlPool.files, resolver)
	}
}

func addUrlResolverToPool(urlPool *UrlPool, resolver *UrlResolver) {
	firstFixedParts := strings.Split(resolver.maskInfo.FixedStart, "/")
	n := len(firstFixedParts) - 1
	placeUrlResolverToHierarchy(urlPool, firstFixedParts[:n], 0, resolver)
}

func (urlPool *UrlPool) RegisterRewrite(from string, to string, condition string, options string) {
	maskInfo := ConvertSingleMask(from)
	addUrlResolverToPool(urlPool, &UrlResolver{
		maskInfo: maskInfo,
		rewrite:  createRewriteInfo(to, condition, options, ScanNamedIds(maskInfo)),
	})
}

func (urlPool *UrlPool) Rewrite(url string) (string, bool) {
	//TODO
	return "", false
}

func matchUrlResolver(url string, resolver *UrlResolver, extraParams *dvevaluation.DvObject) (bool, *UrlResultInfo) {
	isCaseInsensitive := resolver.maskInfo.IsCaseInsensitive == MatchCaseInsensitive
	ok, urlResultInfo := matchMaskUrl(resolver.maskInfo, url, isCaseInsensitive, extraParams)
	if !ok {
		return false, nil
	}
	ConvertUrlResultInfoIntoStringMap(url, urlResultInfo)
	return true, urlResultInfo
}

func (urlPool *UrlPool) RegisterHandlerFunc(url string, handler interface{}) {
	maskInfo := PreparseMaskExpressions(url)
	n := len(maskInfo)
	for i := 0; i < n; i++ {
		addUrlResolverToPool(urlPool, &UrlResolver{
			maskInfo: maskInfo[i],
			Handler:  handler,
		})
	}
}

func hierarchySearch(context interface{}, urlPool *UrlPool, urls []string, fullUrl string, urlVerifier UrlVerifier, extraParams *dvevaluation.DvObject) (bool, *UrlResultInfo) {
	l := len(urls)
	if l > 0 {
		name := urls[0]
		if logLevel {
			log.Printf("dvurl search %s", name)
		}
		dir, ok := urlPool.dirs[name]
		if ok {
			if isOk, urlInfo := hierarchySearch(context, dir, urls[1:], fullUrl, urlVerifier, extraParams); isOk {
				return true, urlInfo
			}
		}
	}
	for _, v := range urlPool.files {
		ok, urlData := matchUrlResolver(fullUrl, v, extraParams)
		if ok {
			if urlVerifier(context, v, urlData) {
				return true, urlData
			}
		}
	}
	if logLevel {
		log.Printf("dvurl search %s not found", fullUrl)
	}
	return false, nil
}

func UrlSearch(context interface{}, urlPool *UrlPool, urls []string, urlVerifier UrlVerifier, extraParams *dvevaluation.DvObject) (bool, *UrlResultInfo) {
	fullUrl := strings.Join(urls, "/")
	return hierarchySearch(context, urlPool, urls, fullUrl, urlVerifier, extraParams)
}

func SingleSimplifiedUrlVerifier(context interface{}, resolver *UrlResolver, urlResultInfo *UrlResultInfo) bool {
	urlResultInfo.CustomObject = resolver.Handler
	return true
}

func SingleSimplifiedUrlSearch(urlPool *UrlPool, url string) (bool, *UrlResultInfo) {
	urls := dvtextutils.ConvertURLToList(strings.ToLower(url))
	extra := map[string]interface{}{"URL": url}
	extraParams := dvparser.GetProptertiesPrototypedToGlobalProperties(extra)
	return UrlSearch(url, urlPool, urls, SingleSimplifiedUrlVerifier, extraParams)
}
