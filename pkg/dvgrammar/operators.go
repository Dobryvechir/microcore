/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

var CommonGrammarBaseDefinition = &GrammarBaseDefinition{
	Operators: map[string]*InterOperator{
		";": &InterOperator{
			Precedence: 0,
			Multi:      true,
		},
		"in": &InterOperator{
			Precedence: 1,
		},
		"of": &InterOperator{
			Precedence: 1,
		},
		"else": &InterOperator{
			Precedence: 1,
			Multi:      true,
		},
		",": &InterOperator{
			Precedence: 1,
			Multi:      true,
		},
		"=": &InterOperator{
			Precedence:         3,
			Multi:              true,
			RightToLeft:        true,
			Assignment:         true,
			LazyLoadedOperands: true,
		},
		"+=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"-=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"**=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"*=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"/=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"%=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"<<=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		">>=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		">>>=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"&=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"^=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"|=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"&&=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"||=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"??=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
		},
		"?": &InterOperator{
			Precedence:         4,
			Multi:              false,
			RightToLeft:        true,
			LazyLoadedOperands: true,
		},
		":": &InterOperator{
			Precedence:         4,
			Multi:              false,
			RightToLeft:        true,
			LazyLoadedOperands: true,
		},
		"||": &InterOperator{
			Precedence:         5,
			Multi:              true,
			LazyLoadedOperands: true,
		},
		"&&": &InterOperator{
			Precedence:         6,
			Multi:              true,
			LazyLoadedOperands: true,
		},
		"??": &InterOperator{
			Precedence:         6,
			Multi:              true,
			LazyLoadedOperands: true,
		},
		"|": &InterOperator{
			Precedence: 7,
			Multi:      true,
		},
		"^": &InterOperator{
			Precedence: 8,
			Multi:      true,
		},
		"&": &InterOperator{
			Precedence: 9,
			Multi:      true,
		},
		"!==": &InterOperator{
			Precedence: 10,
		},
		"===": &InterOperator{
			Precedence: 10,
		},
		"==": &InterOperator{
			Precedence: 10,
		},
		"!=": &InterOperator{
			Precedence: 10,
		},
		"<": &InterOperator{
			Precedence: 11,
		},
		">": &InterOperator{
			Precedence: 11,
		},
		"<=": &InterOperator{
			Precedence: 11,
		},
		">=": &InterOperator{
			Precedence: 11,
		},
		"<<": &InterOperator{
			Precedence: 12,
		},
		">>": &InterOperator{
			Precedence: 12,
		},
		">>>": &InterOperator{
			Precedence: 12,
		},
		"+": &InterOperator{
			Precedence: 13,
			Multi:      true,
		},
		"-": &InterOperator{
			Precedence: 13,
			Multi:      true,
		},
		"*": &InterOperator{
			Precedence: 14,
		},
		"/": &InterOperator{
			Precedence: 14,
		},
		"%": &InterOperator{
			Precedence: 14,
		},
		"**": &InterOperator{
			Precedence: 15,
		},
		"IN": &InterOperator{
			Precedence: 16,
		},
		"NOT IN": &InterOperator{
			Precedence: 16,
		},
		"=>": &InterOperator{
			Precedence: 18,
			Multi:      true,
		},
		"(": &InterOperator{
			Precedence: 19,
			Closing:    ")",
		},
		")": &InterOperator{
			Precedence: 19,
			Opening:    "(",
		},
		"[": &InterOperator{
			Precedence: 19,
			Closing:    "]",
		},
		"]": &InterOperator{
			Precedence: 19,
			Opening:    "[",
		},
		".": &InterOperator{
			Precedence: 19,
			Multi:      true,
		},
	},
	UnaryOperators: map[string]*UnaryOperator{
		"!": &UnaryOperator{
			Post: false,
			Pre:  true,
		},
		"~": &UnaryOperator{
			Post: false,
			Pre:  true,
		},
		"+": &UnaryOperator{
			Post: false,
			Pre:  true,
		},
		"-": &UnaryOperator{
			Post: false,
			Pre:  true,
		},
		"++": &UnaryOperator{
			Post: true,
			Pre:  true,
		},
		"--": &UnaryOperator{
			Post: true,
			Pre:  true,
		},
		"...": &UnaryOperator{
			Post: true,
			Pre:  false,
		},
	},
	VoidOperators: map[string]int{
		"var":        1,
		"const":      1,
		"let":        1,
		"try":        1,
		"catch":      1,
		"public":     1,
		"private":    1,
		"protected":  1,
		"await":      1,
		"class":      1,
		"debugger":   1,
		"enum":       1,
		"export":     1,
		"extends":    1,
		"finally":    1,
		"implements": 1,
		"import":     1,
		"instanceof": 1,
		"interface":  1,
		"new":        1,
		"package":    1,
		"super":      1,
		"static":     1,
		"void":       1,
		"with":       1,
		"yield":      1,
	},
	Language: map[string]*LanguageOperator{
		"return": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: false,
		},
		"for": {
			AlwaysFirst:                  true,
			CanHaveArgument:              false,
			MustHaveArgument:             false,
			ParenthesesFollow:            true,
			CurlyBracesFollowParentheses: true,
			FeatureOptions:               FEATURE_ROUND_BRACKET | FEATURE_CURLY_BRACKETS_OR_CODE | FEATURE_FINISH,
		},
		"if": {
			AlwaysFirst:                  true,
			CanHaveArgument:              false,
			MustHaveArgument:             false,
			ParenthesesFollow:            true,
			CurlyBracesFollowParentheses: true,
			FeatureOptions:               FEATURE_ROUND_BRACKET | FEATURE_CURLY_BRACKETS_OR_CODE | FEATURE_FINISH_OR_ELSE,
		},
		"else": {
			AlwaysFirst:      false,
			CanHaveArgument:  true,
			MustHaveArgument: false,
		},
		"break": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: false,
		},
		"continue": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: false,
		},
		"switch": {
			AlwaysFirst:                  true,
			CanHaveArgument:              false,
			MustHaveArgument:             false,
			ParenthesesFollow:            true,
			CurlyBracesFollowParentheses: true,
		},
		"case": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: true,
		},
		"default": {
			AlwaysFirst:      true,
			CanHaveArgument:  false,
			MustHaveArgument: false,
		},
		"function": {
			AlwaysFirst:                  false,
			CanHaveArgument:              true,
			MustHaveArgument:             false,
			ParenthesesFollow:            true,
			CurlyBracesFollowParentheses: true,
		},
		"while": {
			AlwaysFirst:                  true,
			CanHaveArgument:              false,
			MustHaveArgument:             false,
			ParenthesesFollow:            true,
			CurlyBracesFollowParentheses: true,
		},
		"typeof": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: true,
		},
		"do": {
			AlwaysFirst:      true,
			CanHaveArgument:  false,
			MustHaveArgument: false,
		},
		"throw": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: true,
		},
		"delete": {
			AlwaysFirst:      true,
			CanHaveArgument:  true,
			MustHaveArgument: true,
		},
	},
	DefaultOperator: "",
}
