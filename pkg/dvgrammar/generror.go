/***********************************************************************
MicroCore
Copyright 2020 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import (
	"errors"
	"strconv"
)

func errorMessage(message string, token *Token) error {
	err := message + " at " + token.Value + " in " + token.Place + " (" + strconv.Itoa(token.Row) + ":" + strconv.Itoa(token.Column) + ")"
	return errors.New(err)
}

func errorMessageAtEnd(message string, token *Token) error {
	err := message + " in " + token.Place + " (after " + strconv.Itoa(token.Row) + ":" + strconv.Itoa(token.Column) + ")"
	return errors.New(err)
}

func errorMessageForNode(message string, node *BuildNode, context *ExpressionContext, atEnd bool) error {
	if node == nil {
		return errorMessageForContext(message, context)
	}
	if node.Operator != "" && !atEnd {
		message = message + " at " + node.Operator
	}
	if node.Value != nil {
		if atEnd {
			return errorMessageAtEnd(message, node.Value)
		}
		return errorMessage(message, node.Value)
	}
	l := len(node.Children)
	if l > 0 {
		if atEnd {
			return errorMessageForNode(message, node.Children[l-1], context, true)
		}
		return errorMessageForNode(message, node.Children[0], context, true)
	}
	return errorMessageForContext(message, context)
}

func ErrorMessageForNode(message string, node *BuildNode, context *ExpressionContext) error {
	return errorMessageForNode(message, node, context, false)
}

func errorMessageForContext(message string, context *ExpressionContext) error {
	err := message + " in " + context.Reference.Place + " (" + strconv.Itoa(context.Reference.Row) + ":" + strconv.Itoa(context.Reference.Column) + ")"
	return errors.New(err)
}

func EnrichErrorStr(err error, info string) error {
	if err == nil {
		return err
	}
	return errors.New(err.Error() + "\n   " + info)
}
