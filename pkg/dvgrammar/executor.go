/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

func (tree *BuildNode) ExecuteExpression(context *ExpressionContext) (*ExpressionValue, error) {
	var value *ExpressionValue
	var err error
	if tree.Operator != "" {
		visitor, ok := context.Rules.Visitors[tree.Operator]
		operator, ok1 := context.Rules.BaseGrammar.Operators[tree.Operator]
		if !ok || !ok1 {
			return nil, ErrorMessageForNode("Unexpected operator "+tree.Operator, tree, context)
		}
		l := len(tree.Children)
		var v []*ExpressionValue
		if !operator.LazyLoadedOperands {
			v = make([]*ExpressionValue, l)
			for i := 0; i < l; i++ {
				v[i], err = tree.Children[i].ExecuteExpression(context)
				if err != nil {
					return nil, err
				}
			}
		}
		value, err = visitor(v, tree, context, tree.Operator)
		if err != nil {
			return nil, ErrorMessageForNode(err.Error(), tree, context)
		}
	} else if tree.Value != nil {
		v, err := context.Rules.DataGetter(tree.Value, context)
		if err != nil {
			return nil, err
		}
		value = v
	}
	if len(tree.PostAttributes) != 0 {
		for _, vl := range tree.PostAttributes {
			v, err := context.Rules.UnaryPostVisitors[vl](value, tree, context, vl)
			if err != nil {
				return nil, err
			}
			value = v
		}
	}
	if len(tree.PreAttributes) != 0 {
		for _, vl := range tree.PreAttributes {
			v, err := context.Rules.UnaryPreVisitors[vl](value, tree, context, vl)
			if err != nil {
				return nil, err
			}
			value = v
		}
	}
	return value, nil
}

func (tree *BuildNode) GetChildrenExpressionValue(childNo int, context *ExpressionContext) (*ExpressionValue, error) {
	if childNo < 0 || childNo >= len(tree.Children) {
		return nil, ErrorMessageForNode("Children no is out of range", tree, context)
	}
	return tree.Children[childNo].ExecuteExpression(context)
}

func (tree *BuildNode) GetChildrenNumber() int {
	return len(tree.Children)
}
