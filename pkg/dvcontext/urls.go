/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvcontext

import "regexp"

type MaskKind int

const (
	MaskWord MaskKind = iota
	MaskSlashAware
	MaskSlashUnaware
	MaskRegExp
	MaskCondition
)

type MaskInfoPart struct {
	Min       int
	Max       int
	Kind      MaskKind
	Regex     *regexp.Regexp
	Data      string
	Condition string
}

type MaskInfo struct {
	FixedStart        string
	FixedEnd          string
	Middle            []*MaskInfoPart
	IsNegative        bool
	IsCaseInsensitive int
}
