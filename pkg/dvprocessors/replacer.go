/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvprocessors

import (
	"errors"
	"github.com/Dobryvechir/microcore/pkg/dvcontext"
	"github.com/Dobryvechir/microcore/pkg/dvparser"
	"log"
	"regexp"
	"strings"
)

func isSimpleKey(k string) bool {
	return strings.Index(k, "@@@") != 0
}

func ReplacerHandler(request *dvcontext.RequestContext) bool {
	for k, v := range request.Params {
		if isSimpleKey(k) {
			request.Output = MakeSimpleReplacement(request.Output, []byte(k), []byte(v))
		} else {
			request.Output = MakeComplexReplacement(request.Output, k, v)
		}
	}
	return false
}

func MakeSimpleReplacement(src []byte, old []byte, neue []byte) []byte {
	m := len(src)
	oldLen := len(old)
	newLen := len(neue)
	dif := newLen - oldLen
	changes := 0
	firstMatch := -1
	lastMatch := -1
	b := byte(0)
	if oldLen > 0 {
		b = old[0]
	LoopSimpleEvaluation:
		for i := 0; i < m; i++ {
			if src[i] != b {
				continue
			}
			for j := 1; j < oldLen; j++ {
				if old[j] != src[i+j] {
					continue LoopSimpleEvaluation
				}
			}
			lastMatch = i
			if firstMatch < 0 {
				firstMatch = i
			}
			changes += dif
			i += oldLen - 1
		}
	}
	dst := make([]byte, m+changes)
	if firstMatch < 0 {
		firstMatch = m
	}
	for i := 0; i < firstMatch; i++ {
		dst[i] = src[i]
	}
	if lastMatch >= 0 {
		pos := firstMatch
		for j := 0; j < newLen; j++ {
			dst[pos] = neue[j]
			pos++
		}
	LoopSimpleReplace:
		for i := firstMatch + oldLen; i < lastMatch; i++ {
			if src[i] != b {
				dst[pos] = src[i]
				pos++
				continue
			}
			for j := 1; j < oldLen; j++ {
				if old[j] != src[i+j] {
					dst[pos] = src[i]
					pos++
					continue LoopSimpleReplace
				}
			}
			i += oldLen - 1
			for j := 0; j < newLen; j++ {
				dst[pos] = neue[j]
				pos++
			}
		}
		if lastMatch > firstMatch {
			for j := 0; j < newLen; j++ {
				dst[pos] = neue[j]
				pos++
			}
		}
		for i := lastMatch + oldLen; i < m; i++ {
			dst[pos] = src[i]
			pos++
		}
	}
	return dst
}

var defReplaceParams = []string{"", "", "SP0,T", "SP0,T"}

const (
	moveKindPos     = 0
	moveKindWord    = 1
	moveKindLine    = 2
	searchKindPlain = 0
	searchKindRegex = 1
	searchDirFirst  = 0
	searchDirLast   = 1
)

type replacerOptions struct {
	moveKind            int
	movePos             int
	moveWhere           int
	searchKind          int
	searchCompiled      *regexp.Regexp
	postSearchMoveKind  int
	postSearchMovePos   int
	postSearchMoveWhere int
	crlf                []byte
}

func readMoveKindPosWhere(options string, defPos int, defWhere int) (kind int, pos int, where int) {
	pos = defPos
	where = defWhere
	n := len(options)
	i := 1
	for ; i < n; i++ {
		c := options[i]
		switch c {
		case '^':
			where = searchDirFirst
		case '$':
			where = searchDirLast
		case 'L':
			kind = moveKindLine
			pos = dvparser.TryReadInteger(options[i+1:], defPos)
			return
		case 'P':
			kind = moveKindPos
			pos = dvparser.TryReadInteger(options[i+1:], defPos)
			return
		case 'W':
			kind = moveKindWord
			pos = dvparser.TryReadInteger(options[i+1:], defPos)
			return
		default:
			log.Printf("Unknown option at %d in %s", i, options)
			return
		}
	}
	return
}

func processReplaceOptions(pattern string, defStartWhere int) (res *replacerOptions) {
	res = &replacerOptions{moveWhere: defStartWhere}
	options := dvparser.ConvertToNonEmptyList(pattern)
	n := len(options)
	for i := 0; i < n; i++ {
		s := options[i]
		switch s[0] {
		case 'S':
			res.moveKind, res.movePos, res.moveWhere = readMoveKindPosWhere(s, res.movePos, res.moveWhere)
		case 'F':
			res.postSearchMoveKind, res.postSearchMovePos, res.postSearchMoveWhere = readMoveKindPosWhere(s, res.postSearchMovePos, res.postSearchMoveWhere)
		case 'T':
			res.searchKind = searchKindPlain
		case 'R':
			res.searchKind = searchKindRegex
		case 'C':
			if s == "CR" {
				res.crlf = []byte{13}
			} else if s == "CRLF" {
				res.crlf = []byte{13, 10}
			}
		case 'L':
			if s == "LF" {
				res.crlf = []byte{10}
			}
		default:
			log.Printf("Uknown option %s", s)
		}
	}
	return
}

func MoveLinesInText(src []byte, lim int, pos int, lines int, dir int) int {
	if dir == searchDirFirst {
		for pos > 0 && src[pos-1] != 10 && src[pos-1] != 13 {
			pos--
		}
		lines = -lines
	} else {
		for pos < lim && src[pos] != 10 && src[pos-1] != 13 {
			pos++
		}
		if pos < lim {
			pos++
			if pos < lim && src[pos] == 10 && src[pos-1] == 13 {
				pos++
			}
		}
	}
	if lines < 0 {
		for pos > 0 && lines < 0 {
			pos--
			if pos > 0 {
				c := src[pos-1]
				if c == 10 {
					lines++
				} else if c == 13 && src[pos] != 10 {
					lines++
				}
			}
		}
	} else {
		for pos < lim && lines > 0 {
			c := src[pos]
			if c == 13 {
				if pos+1 < lim && src[pos+1] != 10 {
					lines--
				}
			} else if c == 10 {
				lines--
			}
			pos++
		}
	}
	return pos
}

func MoveWordsInText(src []byte, lim int, pos int, words int, dir int) int {
	if dir == searchDirFirst {
		for pos > 0 && src[pos-1] > ' ' {
			pos--
		}
		words = -words
	} else {
		for pos < lim && src[pos] > ' ' {
			pos++
		}
	}
	if words < 0 {
		for pos > 0 && words < 0 {
			pos--
			for pos > 0 && src[pos] <= ' ' {
				pos--
			}
			for pos > 0 && src[pos-1] > ' ' {
				pos--
			}
			words++
		}
	} else {
		for pos < lim && words > 0 {
			for pos < lim && src[pos] <= ' ' {
				pos++
			}
			for pos < lim && src[pos] > ' ' {
				pos++
			}
			words--
		}
	}
	return pos
}

func MoveInTextByMoveKind(src []byte, lim int, pos int, posLast int, points int, kind int, where int) int {
	if where == searchDirFirst {
		pos = posLast
	} else {
		points = -points
	}
	switch kind {
	case moveKindPos:
		pos -= points
		if pos < 0 {
			pos = 0
		}
		if pos > lim {
			pos = lim
		}
	case moveKindWord:
		pos = MoveWordsInText(src, lim, pos, points, where)
	case moveKindLine:
		pos = MoveLinesInText(src, lim, pos, points, where)
	}
	return pos
}

func SearchTextBySearchKind(src []byte, startLimit int, endLimit int, search []byte, dir int, options *replacerOptions) (start int, end int, ok bool) {
	m := len(search)
	if m == 0 {
		return startLimit, endLimit, true
	}
	k := m - 1
	switch options.searchKind {
	case searchKindPlain:
		if dir == searchDirFirst {
			b := search[0]
			n := endLimit - k
		searchPlainFirst:
			for ; startLimit < n; startLimit++ {
				if src[startLimit] != b {
					continue
				}
				for j := 1; j < m; j++ {
					if search[j] != src[startLimit+j] {
						continue searchPlainFirst
					}
				}
				return startLimit, startLimit + m, true
			}
		} else {
			b := search[m-1]
			n := startLimit + k
		searchPlainLast:
			for endLimit--; endLimit >= n; endLimit-- {
				if src[endLimit] != b {
					continue
				}
				for j := 0; j < k; j++ {
					if search[j] != src[endLimit-k+j] {
						continue searchPlainLast
					}
				}
				return endLimit - k, endLimit + 1, true
			}
		}
	case searchKindRegex:
		if options.searchCompiled == nil {
			r, err := regexp.Compile(string(search))
			if err != nil {
				log.Printf("Incorrect regexp %s: %v", search, err)
			}
			return -1, -1, false
			options.searchCompiled = r
		}
		if dir == searchDirFirst {
			res := options.searchCompiled.FindIndex(src[startLimit:endLimit])
			if len(res) == 2 {
				return startLimit + res[0], startLimit + res[1], true
			}
		} else {
			res := options.searchCompiled.FindAllIndex(src[startLimit:endLimit], -1)
			n := len(res)
			if n > 0 {
				return startLimit + res[n-1][0], startLimit + res[n-1][1], true
			}
		}
	}
	return endLimit, endLimit, false
}

func MakeComplexReplacement(src []byte, k string, neue string) []byte {
	params := strings.Split(k, "@@@")
	if len(params) < 2 {
		log.Printf("Complex replacement must have at least 1 parameter")
		return src
	}
	m := len(src)
	old := []byte(params[1])
	rest := dvparser.GetStringArrayWithDefaults(params[2:], defReplaceParams)
	startSearch := []byte(rest[0])
	endSearch := []byte(rest[1])
	whereStart := processReplaceOptions(rest[2], searchDirFirst)
	whereEnd := processReplaceOptions(rest[3], searchDirLast)
	oldLen := len(old)
	neuBuf := []byte(neue)
	b := old[0]
	bend := m - oldLen + 1
LoopEqualLen:
	for i := 0; i < bend; i++ {
		if src[i] != b {
			continue
		}
		for j := 1; j < oldLen; j++ {
			if old[j] != src[i+j] {
				continue LoopEqualLen
			}
		}
		posStart := MoveInTextByMoveKind(src, m, i, i+oldLen, whereStart.movePos, whereStart.moveKind, whereStart.moveWhere)
		posEnd := MoveInTextByMoveKind(src, m, i, i+oldLen, whereEnd.movePos, whereEnd.moveKind, whereEnd.moveWhere)
		if posStart > posEnd {
			posStart = posEnd
		}
		start, end, ok := SearchTextBySearchKind(src, posStart, posEnd, startSearch, searchDirFirst, whereStart)
		if !ok {
			continue
		}
		posStart = MoveInTextByMoveKind(src, m, start, end, whereStart.postSearchMovePos, whereStart.postSearchMoveKind, whereStart.postSearchMoveWhere)
		start, end, ok = SearchTextBySearchKind(src, posStart, posEnd, endSearch, searchDirLast, whereEnd)
		if !ok {
			continue
		}
		posEnd = MoveInTextByMoveKind(src, m, start, end, whereEnd.postSearchMovePos, whereEnd.postSearchMoveKind, whereEnd.postSearchMoveWhere)
		if posStart > posEnd {
			continue
		}
		var dif int
		src, dif = dvparser.InsertTextIntoBuffer(src, posStart, posEnd, whereStart.crlf, neuBuf, whereEnd.crlf)
		i += dif
		reqMin := posStart + len(neue) - 1
		if i < reqMin {
			i = reqMin
		}
	}
	return src
}

func replaceInitHandler(params []string) (map[string]string, error) {
	n := len(params)
	if n == 0 {
		return nil, errors.New("replacer PostProcessor must have parameters (specified as Param)")
	}
	if n&1 != 0 {
		return nil, errors.New("replacer PostProcessor parameters must go in pairs")
	}
	m := n >> 1
	res := make(map[string]string)
	for i := 0; i < m; i++ {
		k := strings.TrimSpace(params[i<<1])
		v := strings.TrimSpace(params[i<<1|1])
		if _, ok := res[k]; ok {
			return nil, errors.New("In replacer PostProcessor in each pair the first parameter must be unique")
		}
		res[k] = v
	}
	return res, nil
}

var replacerFileConfig *RegistrationConfig = &RegistrationConfig{
	Name:              "replacer",
	EndPointHandler:   ReplacerHandler,
	ServerInitHandler: replaceInitHandler,
}

var replacerFileInited bool = RegisterPostProcessor(replacerFileConfig, true)
