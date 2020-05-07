/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

const (
	BYTE_ONLY_INSIDE_STRING = iota
	BYTE_CONTROL
	BYTE_DIGIT
	BYTE_LETTER
	BYTE_LIST_CLOSING
	BYTE_LIST_SEPARATOR
	BYTE_OPERATOR
	BYTE_QUOTATION
	BYTE_SPACE
)

type TokenBlock struct {
	Kind        int
	IsCommon    bool
	Controls    []map[string]string
	MaxWords    int
	MaxLength   int
	Recognizers []TypeRecognizer
}

type TokenMapType [256]*TokenBlock
type GrammarTable struct {
	TokenMap              TokenMapType
	IsSingleQuoteTypeChar bool
}
