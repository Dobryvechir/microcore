/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
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
	FEATURE_ROUND_BRACKET          = 1 << iota
	FEATURE_CURLY_BRACKETS         = 1 << iota
	FEATURE_CURLY_BRACKETS_OR_CODE = 1 << iota
	FEATURE_FINISH                 = 1 << iota
	FEATURE_FINISH_OR_ELSE         = 1 << iota
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
	TYPE_ERROR      = iota
	TYPE_MASK       = 0x7ff
)

const (
	FLOW_NORMAL = iota
	FLOW_RETURN
	FLOW_BREAK
	FLOW_CONTINUE
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

type LanguageOperator struct {
	AlwaysFirst                  bool
	CanHaveArgument              bool
	MustHaveArgument             bool
	ParenthesesFollow            bool
	CurlyBracesFollowParentheses bool
	FeatureOptions               int
}

type UnaryOperator struct {
	Post bool
	Pre  bool
}

type ScopeInterface interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	StackPush(option int)
	StackPop()
	SetDeep(key string, value interface{})
}

type GrammarBaseDefinition struct {
	Operators       map[string]*InterOperator
	UnaryOperators  map[string]*UnaryOperator
	VoidOperators   map[string]int
	Language        map[string]*LanguageOperator
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
type BracketOperatorVisitor func(*ExpressionValue, *BuildNode, *ExpressionContext, []*BuildNode) (*ExpressionValue, *ExpressionValue, bool, error, bool)
type LanguageOperatorVisitor func(*BuildNode, *ExpressionContext) (int, *ExpressionValue, error)

type UnaryVisitor func(*ExpressionValue, *BuildNode, *ExpressionContext, string, string, *ExpressionValue) (*ExpressionValue, error)

type ExpressionValue struct {
	DataType int
	Value    interface{}
	Name     string
	Parent   *ExpressionValue
}

var ErrorExpressionValue = &ExpressionValue{DataType: TYPE_ERROR}

const (
	EVALUATE_OPTION_UNDEFINED = 1 << (iota + 8)
	EVALUATE_OPTION_PARENT = 1 << (iota + 8)
	EVALUATE_OPTION_NAME = 1 << (iota +8)
)

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
	LanguageOperator  map[string]LanguageOperatorVisitor
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
