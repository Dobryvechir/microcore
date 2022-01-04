/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"math"
)

const (
	BooleanFalse = "false"
	BooleanTrue  = "true"
)

const (
	EVALUATE_OPTION_UNDEFINED = 1 << (iota + 8)
)

var reservedWords = map[string]int{
	"null":      dvgrammar.TYPE_NULL,
	"undefined": dvgrammar.TYPE_UNDEFINED,
	"true":      dvgrammar.TYPE_BOOLEAN,
	"false":     dvgrammar.TYPE_BOOLEAN,
	"NaN":       dvgrammar.TYPE_NAN,
}

type DvObject struct {
	Value      interface{}
	Options    int
	Properties map[string]interface{}
	Prototype  *DvObject
}

var buildinTypes map[string]interface{} = map[string]interface{}{
	"true":      true,
	"false":     false,
	"undefined": nil,
	"NaN":       math.NaN(),
	"null":      DvObject_null,
	"":          "",
}

const (
	ConversionOptionJSLike     = 0
	ConversionOptionSimpleLike = 1
	ConversionOptionJsonLike   = 2
)

var nullValueVersion = map[int]string{
	ConversionOptionJSLike:     "undefined",
	ConversionOptionSimpleLike: "",
	ConversionOptionJsonLike:   "null",
}

type DvFunctionObject struct {
	SelfRef   interface{}
	Context   *dvgrammar.ExpressionContext
	Executor  *DvFunction
}
