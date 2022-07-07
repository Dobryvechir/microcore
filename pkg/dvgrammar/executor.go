/***********************************************************************
MicroCore
Copyright 2017 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import (
	"errors"
)

func (tree *BuildNode) ExecuteExpression(context *ExpressionContext) (int, *ExpressionValue, error) {
	var value *ExpressionValue
	var err error
	var flow = 0
	l := len(tree.Children)
	var lastVarName string
	var lastParent *ExpressionValue
	if tree.Operator != "" {
		visitor, ok := context.Rules.Visitors[tree.Operator]
		operator, ok1 := context.Rules.BaseGrammar.Operators[tree.Operator]
		if !ok || !ok1 {
			langVisitor, ok := context.Rules.LanguageOperator[tree.Operator]
			if ok {
				flow, value, err = langVisitor(tree, context)
				return flow, value, err
			}
			return flow, nil, ErrorMessageForNode("Unexpected operator "+tree.Operator, tree, context)
		}
		var v []*ExpressionValue
		if !operator.LazyLoadedOperands {
			v = make([]*ExpressionValue, l)
			for i := 0; i < l; i++ {
				flow, v[i], err = tree.Children[i].ExecuteExpression(context)
				if err != nil {
					return flow, nil, err
				}
			}
		}
		value, err = visitor(v, tree, context, tree.Operator)
		if err != nil {
			return flow, nil, ErrorMessageForNode(err.Error(), tree, context)
		}
	} else if tree.Value != nil {
		value = nil
		useParent := (context.VisitorOptions & EVALUATE_OPTION_PARENT) != 0
		hasNoParent := tree.Value.DataType == TYPE_FUNCTION
		if !hasNoParent {
			lastVarName = tree.Value.Value
			value, err = context.Rules.DataGetter(tree.Value, context)
			if err != nil {
				if useParent && l > 0 && value != nil {
					value.Parent = ErrorExpressionValue
				}
				return flow, value, err
			}
		}
		if l > 0 {
			needName := false
			if (context.VisitorOptions & EVALUATE_OPTION_NAME) == 0 {
				needName = hasAssigningOperators(tree.PostAttributes)
				if !needName {
					needName = hasAssigningOperators(tree.PreAttributes)
				}
			}
			if needName {
				context.VisitorOptions |= EVALUATE_OPTION_NAME
			}
			value, lastParent, err = ExecuteBracketExpression(value, hasNoParent, tree.Children, context)
			if needName {
				context.VisitorOptions ^= EVALUATE_OPTION_NAME
			}
			if err != nil {
				return flow, value, err
			}
			if useParent && lastParent != nil && value != nil {
				value.Parent = lastParent
			}
		}
	}
	if len(tree.PostAttributes) != 0 {
		if lastParent!=nil && value!=nil && value.Name!="" {
			lastVarName = value.Name
		}
		for _, vl := range tree.PostAttributes {
			v, err := context.Rules.UnaryPostVisitors[vl](value, tree, context, vl, lastVarName, lastParent)
			if err != nil {
				return flow, nil, err
			}
			value = v
		}
	}
	if len(tree.PreAttributes) != 0 {
		if lastParent!=nil && value!=nil && value.Name!="" {
			lastVarName = value.Name
		}
		for _, vl := range tree.PreAttributes {
			v, err := context.Rules.UnaryPreVisitors[vl](value, tree, context, vl, lastVarName, lastParent)
			if err != nil {
				return flow, nil, err
			}
			value = v
		}
	}
	return flow, value, nil
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
	_, val, err := tree.Children[childNo].ExecuteExpression(context)
	return val, err
}

func (tree *BuildNode) GetChildrenNumber() int {
	return len(tree.Children)
}

func hasAssigningOperators(attrs []string) bool {
	n := len(attrs)
	for i := 0; i < n; i++ {
		if attrs[i] == "++" || attrs[i] == "--" {
			return true
		}
	}
	return false
}
