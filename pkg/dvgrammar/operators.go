/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

var CommonGrammarBaseDefinition = &GrammarBaseDefinition{
	Operators: map[string]*InterOperator{
		";": &InterOperator{
			Precedence: 0,
			Multi:      true,
		},
		",": &InterOperator{
			Precedence: 1,
			Multi:      true,
		},
		"=": &InterOperator{
			Precedence:  3,
			Multi:       true,
			RightToLeft: true,
			Assignment:  true,
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
		"?": &InterOperator{
			Precedence:  4,
			Multi:       true,
			RightToLeft: true,
		},
		":": &InterOperator{
			Precedence:  4,
			Multi:       true,
			RightToLeft: true,
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
	DefaultOperator: "",
}
