/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvevaluation

import (
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
)

const (
	BOOLEAN_FALSE = "false"
	BOOLEAN_TRUE  = "true"
)

const (
	UNARY_BOOLEAN_NOT = 1 << iota
	UNARY_LOGICAL_NOT = 1 << iota
	UNARY_MINUS       = 1 << iota
	UNARY_PLUS        = 1 << iota
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