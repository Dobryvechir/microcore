/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

//bitwise values of EvaluateOptions
const (
	NON_LATIN_LETTERS_ONLY_IN_STRINGS = 1 << iota
	SINGLE_QUOTE_HAS_TYPE_CHAR
	EXCLUDE_RECOGNIZER_BYTE_LETTER
	EXCLUDE_RECOGNIZER_BYTE_DIGIT
	EXCLUDE_RECOGNIZER_BYTE_QUOTATION
	EXCLUDE_RECOGNIZER_BYTE_CONTROL
	EXCLUDE_RECOGNIZER_BYTE_LIST_CONTROL
)

//bitwise values of VisitorOptions
const (
	VISITOR_OPTION_CASHED = 1 << iota
)
const (
	TYPE_CONTROL    = iota
	TYPE_OPERATOR   = iota
	TYPE_STRING     = iota
	TYPE_CHAR       = iota
	TYPE_NUMBER     = iota
	TYPE_NUMBER_INT = iota
	TYPE_BOOLEAN    = iota
	TYPE_NULL       = iota
	TYPE_UNDEFINED  = iota
	TYPE_NAN        = iota
	TYPE_FUNCTION   = iota
	TYPE_DATA       = iota
	TYPE_OBJECT     = iota
	TYPE_MASK       = 0x7ff
)

type Token struct {
	DataType int
	Row      int
	Column   int
	Place    string
	Value    string
}

type InterOperator struct {
	Precedence         int
	Multi              bool
	RightToLeft        bool
	LazyLoadedOperands bool
	Assignment         bool
	Opening            string
	Closing            string
}

type UnaryOperator struct {
	Post bool
	Pre  bool
}

type ScopeInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

type GrammarBaseDefinition struct {
	Operators       map[string]*InterOperator
	UnaryOperators  map[string]*UnaryOperator
	DefaultOperator string
}

type BuildNode struct {
	Children       []*BuildNode
	Parent         *BuildNode
	Operator       string
	Value          *Token
	PreAttributes  []string
	PostAttributes []string
	Group          int
	closed         bool
}

type InterOperatorVisitor func([]*ExpressionValue, *BuildNode, *ExpressionContext, string) (*ExpressionValue, error)
type BracketOperatorVisitor func(*ExpressionValue,*BuildNode, *ExpressionContext, []*BuildNode) (*ExpressionValue, *ExpressionValue, bool, error, bool)

type UnaryVisitor func(*ExpressionValue, *BuildNode, *ExpressionContext, string) (*ExpressionValue, error)

type ExpressionValue struct {
	DataType int
	Value    interface{}
}

type TypeRecognizer struct {
	StartByteTypes []int
	Processor      func(grammar *GrammarTable, data []byte, pos int, n int) (*Token, int, error)
	exceptFlag     int
}

type SourceReference struct {
	Row    int
	Column int
	Place  string
}

type GrammarRuleDefinitions struct {
	BaseGrammar       *GrammarBaseDefinition
	GrammarAliases    map[string]string
	Recognizers       []TypeRecognizer
	EvaluateOptions   int
	Visitors          map[string]InterOperatorVisitor
	BracketVisitor    map[string]BracketOperatorVisitor
	UnaryPostVisitors map[string]UnaryVisitor
	UnaryPreVisitors  map[string]UnaryVisitor
	DataGetter        func(*Token, *ExpressionContext) (*ExpressionValue, error)
	cache             map[string][]*BuildNode
	Grammar           *GrammarTable
}

type ExpressionContext struct {
	Scope          ScopeInterface
	Reference      *SourceReference
	Rules          *GrammarRuleDefinitions
	VisitorOptions int
}
