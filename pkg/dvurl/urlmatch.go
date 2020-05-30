package dvurl

import (
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvevaluation"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"regexp"
	"strconv"
	"strings"
)

const (
	urlMatchLenOffset = 16
)

func matchMaskPartForRegularExpression(url string, regex *regexp.Regexp, start int, isCaseInsensitive bool, track *UrlResultInfo) (bool, *UrlResultInfo) {
	res := regex.FindStringSubmatchIndex(url)
	n := len(res) >> 1
	if n == 0 {
		return false, nil
	}
	for k := 1; k < n; k++ {
		pos := k << 1
		track.posInfo = append(track.posInfo, res[pos]+start|(res[pos+1]-res[pos])<<urlMatchLenOffset)
	}
	return true, track
}

func matchMaskMiddleForWildOrRegularExpressions(url string, parts []*dvcontext.MaskInfoPart, start int, isCaseInsensitive bool, track *UrlResultInfo) (bool, *UrlResultInfo) {
	n := len(parts)
	urlLen := len(url)
	switch n {
	case 0:
		if url == "" {
			return true, track
		} else {
			return false, nil
		}
	case 1:
		if parts[0].Min > urlLen || parts[0].Max < urlLen {
			return false, nil
		}
		switch parts[0].Kind {
		case dvcontext.MaskRegExp:
			return matchMaskPartForRegularExpression(url, parts[0].Regex, start, isCaseInsensitive, track)
		case dvcontext.MaskSlashAware:
			if strings.Index(url, "/") >= 0 {
				return false, nil
			}
			if parts[0].Data != "" {
				track.namedIds[parts[0].Data] = strconv.Itoa(len(track.posInfo) + 1)
			}
		}
		track.posInfo = append(track.posInfo, start|urlLen<<urlMatchLenOffset)
		return true, track
	}
	pos := parts[0].Min
	posMax := urlLen
	for i := 1; i < n; i++ {
		posMax -= parts[i].Min
	}
	if parts[0].Kind == dvcontext.MaskSlashAware {
		slashPos := strings.Index(url, "/")
		if slashPos >= 0 && slashPos < posMax {
			posMax = slashPos
		}
	}
	for ; pos < posMax; pos++ {
		trackNew := track.HalfClone()
		ok := true
		if parts[0].Kind == dvcontext.MaskRegExp {
			ok, trackNew = matchMaskPartForRegularExpression(url[:pos], parts[0].Regex, start, isCaseInsensitive, trackNew)
		} else {
			if parts[0].Kind == dvcontext.MaskSlashAware && parts[0].Data != "" {
				trackNew.namedIds[parts[0].Data] = strconv.Itoa(len(track.posInfo) + 1)
			}
			trackNew.posInfo = append(trackNew.posInfo, start|(pos<<urlMatchLenOffset))
		}
		if ok {
			ok, trackNew := matchMaskMiddleForWildOrRegularExpressions(url[:pos], parts[1:], start+pos, isCaseInsensitive, trackNew)
			if ok {
				return true, trackNew
			}
		}
	}
	return false, nil
}

func matchMaskMiddle(url string, parts []*dvcontext.MaskInfoPart, start int, isCaseInsensitive bool, track *UrlResultInfo) (bool, *UrlResultInfo) {
	n := len(parts)
	if n == 0 {
		if url == "" {
			return true, track
		} else {
			return false, nil
		}
	}
	i := 0
	for ; i < n; i++ {
		kind := parts[i].Kind
		if kind == dvcontext.MaskWord || kind == dvcontext.MaskCondition {
			break
		}
	}
	if i == n {
		return matchMaskMiddleForWildOrRegularExpressions(url, parts, start, isCaseInsensitive, track)
	}
	word := parts[i].Data
	if parts[i].Kind == dvcontext.MaskCondition {
		word = track.Conditions[parts[i].Condition]
	}
	if isCaseInsensitive {
		word = strings.ToLower(word)
	}
	m := len(word)
	pos := strings.Index(url, word)
	if pos < 0 {
		return false, nil
	}
	if i == 0 {
		if pos != 0 {
			return false, nil
		}
		if parts[i].Kind == dvcontext.MaskCondition {
			track.posInfo = append(track.posInfo, start|(m<<urlMatchLenOffset))
		}
		return matchMaskMiddle(url[m:], parts[1:], start+m, isCaseInsensitive, track)
	}
	for {
		res, trackNew := matchMaskMiddleForWildOrRegularExpressions(url[:pos], parts[:i], start, isCaseInsensitive, track.HalfClone())
		if res {
			if parts[i].Kind == dvcontext.MaskCondition {
				trackNew.posInfo = append(trackNew.posInfo, pos+start|(m<<urlMatchLenOffset))
			}
			res, trackNew = matchMaskMiddle(url[start+pos+m:], parts[i+1:], start+pos+m, isCaseInsensitive, trackNew)
			if res {
				return true, trackNew
			}
		}
		newPos := strings.Index(url[pos+1:], word)
		if newPos < 0 {
			break
		}
		pos += newPos + 1
	}
	return false, nil
}

func matchMaskUrl(mask *dvcontext.MaskInfo, url string, isCaseInsensitive bool, extraParams *dvevaluation.DvObject) (bool, *UrlResultInfo) {
	if isCaseInsensitive {
		url = strings.ToLower(url)
	}
	n := len(mask.FixedStart)
	e := len(mask.FixedEnd)
	u := len(url)
	if u < n+e {
		return false, nil
	}
	if n != 0 {
		s := mask.FixedStart
		if isCaseInsensitive {
			s = strings.ToLower(s)
		}
		if url[:n] != s {
			return false, nil
		}
		url = url[n:]
		u -= n
	}
	if e != 0 {
		u -= e
		s := mask.FixedEnd
		if isCaseInsensitive {
			s = strings.ToLower(s)
		}
		if url[u:] != s {
			return false, nil
		}
		url = url[:u]
	}
	middleLength := len(mask.Middle)
	track := &UrlResultInfo{posInfo: make([]int, 0, 16), namedIds: make(map[string]string)}
	if u > 0 || middleLength > 0 {
		for i := 0; i < middleLength; i++ {
			if mask.Middle[i].Kind == dvcontext.MaskCondition {
				if track.Conditions == nil {
					track.Conditions = make(map[string]string)
				}
				condition := mask.Middle[i].Condition
				if _, ok := track.Conditions[condition]; !ok {
					track.Conditions[condition] = dvparser.StringEvaluatorAsStringWithErrorLoggedAsWarning(condition, extraParams, "~~__error__~~")
				}
			}
		}
		return matchMaskMiddle(url, mask.Middle, n, isCaseInsensitive, track)
	}
	return true, track
}

func ConvertUrlResultInfoIntoStringMap(url string, track *UrlResultInfo) {
	posInfo := track.posInfo
	n := len(posInfo)
	res := make(map[string]string)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i + 1)
		val := posInfo[i]
		valOffset := val & (1<<urlMatchLenOffset - 1)
		valLength := val >> urlMatchLenOffset
		valStr := url[valOffset : valOffset+valLength]
		res[key] = valStr
	}
	for k, v := range track.namedIds {
		res[k] = res[v]
	}
	track.UrlKeys = res
}
