/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

var LogicalGrammarBaseDefinition = &GrammarBaseDefinition{
	Operators: map[string]*InterOperator{
		"||": &InterOperator{
			Precedence:         0,
			Multi:              true,
			LazyLoadedOperands: true,
		},
		"&&": &InterOperator{
			Precedence:         1,
			Multi:              true,
			LazyLoadedOperands: true,
		},
	},
	UnaryOperators: map[string]*UnaryOperator{
		"!": &UnaryOperator{
			Post: false,
			Pre:  true,
		},
	},
	DefaultOperator: "&&",
}
