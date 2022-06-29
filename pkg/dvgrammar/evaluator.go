/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

func Compile(data []byte, context *ExpressionContext) ([]*BuildNode, error) {
	CheckCreateGrammarTable(context.Rules)
	tokens, err := Tokenize(context.Reference, data, context.Rules.Grammar)
	if err != nil {
		return nil, err
	}
	tree, err1 := buildExpressionTree(tokens, context.Rules.BaseGrammar)
	if err1 != nil {
		return nil, err1
	}
	return tree, nil
}

func CompileOrCache(data []byte, context *ExpressionContext, cache bool) ([]*BuildNode, error) {
	if !cache {
		return Compile(data, context)
	}
	key := string(data)
	if context.Rules.cache == nil {
		context.Rules.cache = make(map[string][]*BuildNode)
	}
	tree, ok := context.Rules.cache[key]
	if ok {
		return tree, nil
	}
	var err error
	tree, err = Compile(data, context)
	if err != nil {
		return nil, err
	}
	context.Rules.cache[key] = tree
	return tree, nil
}

func FastEvaluation(data []byte, context *ExpressionContext) (*ExpressionValue, error) {
	cache := (context.VisitorOptions & VISITOR_OPTION_CASHED) != 0
	forest, err := CompileOrCache(data, context, cache)
	if err != nil {
		return nil, err
	}
	var value *ExpressionValue
	value, err = BuildNodeExecution(forest, context)
	n:=len(forest)
	if !cache {
		for i := 0; i < n; i++ {
			tree := forest[i]
			fullTreeClean(tree)
		}
	}
	return value, err
}

func BuildNodeExecution(nodes []*BuildNode, context *ExpressionContext) (value *ExpressionValue,err error) {
	flow := 0
	n := len(nodes)
forestTrack:
	for i := 0; i < n; i++ {
		tree := nodes[i]
		flow, value, err = tree.ExecuteExpression(context)
		if err != nil {
			break
		}
		switch flow {
		case FLOW_RETURN:
			break forestTrack
		}
	}
	return value, err
}
