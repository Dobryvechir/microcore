/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import "errors"

func (tree *BuildNode) ExecuteExpression(context *ExpressionContext) (*ExpressionValue, error) {
	var value *ExpressionValue
	var err error
	l := len(tree.Children)
	if tree.Operator != "" {
		visitor, ok := context.Rules.Visitors[tree.Operator]
		operator, ok1 := context.Rules.BaseGrammar.Operators[tree.Operator]
		if !ok || !ok1 {
			return nil, ErrorMessageForNode("Unexpected operator "+tree.Operator, tree, context)
		}
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
		value = nil
		hasNoParent := tree.Value.DataType == TYPE_FUNCTION
		if !hasNoParent {
			value, err = context.Rules.DataGetter(tree.Value, context)
			if err != nil {
				return nil, err
			}
		}
		if l > 0 {
			value, _, err = ExecuteBracketExpression(value, hasNoParent, tree.Children, context)
		}
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

func ExecuteBracketExpression(parent *ExpressionValue, hasNoParent bool, nodes []*BuildNode, context *ExpressionContext) (*ExpressionValue, *ExpressionValue, error) {
	n := len(nodes)
	var child *ExpressionValue = parent
	var toStop bool
	var err error
	for i := 0; i < n; i++ {
		node := nodes[i]
		parent = child
		key := node.Operator
		if !hasNoParent {
			key = "*" + key
		}
		bracketVisitor := context.Rules.BracketVisitor[key]
		if bracketVisitor == nil {
			return parent, child, errors.New("Bracket operator " + key + " is not supported")
		}
		child, parent, toStop, err, hasNoParent = bracketVisitor(parent, node, context, nodes[i+1:])
		if err != nil {
			return child, parent, err
		}
		if toStop {
			break
		}
	}
	return child, parent, nil
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
