/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvurl

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"strings"
)

const (
	MatchDefaultFalse = iota
	MatchDefaultTrue
	MatchDefaultFalseIfNoneOrNotNegationAtFirst
	MatchDefaultTrueIfNoneOrNegationAtFirst
	MatchDefaultMask     = 0x1f
	MatchCaseInsensitive = 32
)

type UrlResultInfo struct {
	namedIds     map[string]string
	posInfo      []int
	Custom       string
	Conditions   map[string]string
	UrlKeys      map[string]string
	CustomObject interface{}
}

func StringMapClone(src map[string]string) (dst map[string]string) {
	dst = make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (urlResultInfo *UrlResultInfo) HalfClone() *UrlResultInfo {
	return &UrlResultInfo{
		namedIds:     StringMapClone(urlResultInfo.namedIds),
		posInfo:      urlResultInfo.posInfo[:],
		Conditions:   urlResultInfo.Conditions,
		Custom:       urlResultInfo.Custom,
		CustomObject: urlResultInfo.CustomObject,
	}
}

func MatchMasksForUrlParts(masks []*dvcontext.MaskInfo, url []string, extraParams *dvevaluation.DvObject) bool {
	return MatchMasksWithDefault(masks, url, MatchDefaultFalse, extraParams)
}

func MatchMasksWithDefault(masks []*dvcontext.MaskInfo, url []string, defaultMatch int, extraParams *dvevaluation.DvObject) bool {
	urlFull := strings.Join(url, "/")
	n := len(masks)
	isCaseInsensitive := (defaultMatch & MatchCaseInsensitive) != 0
	defaultMatch &= MatchDefaultMask
	match := defaultMatch == MatchDefaultTrueIfNoneOrNegationAtFirst || defaultMatch == MatchDefaultTrue
	if n != 0 && (defaultMatch == MatchDefaultTrueIfNoneOrNegationAtFirst || defaultMatch == MatchDefaultFalseIfNoneOrNotNegationAtFirst) {
		match = masks[0].IsNegative
	}
	for i := 0; i < n; i++ {
		negative := masks[i].IsNegative
		if match == negative {
			insensitive := isCaseInsensitive
			if masks[i].IsCaseInsensitive == MatchCaseInsensitive {
				insensitive = true
			} else if masks[i].IsCaseInsensitive == -MatchCaseInsensitive {
				insensitive = false
			}
			if ok, _ := matchMaskUrl(masks[i], urlFull, insensitive, extraParams); ok {
				match = !match
			}
		}
	}
	return match
}
