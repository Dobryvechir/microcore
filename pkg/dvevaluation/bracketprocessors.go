/***********************************************************************
MicroCore
Copyright 2020 - 2022 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/

package dvevaluation

import (
	"errors"
	"fmt"
	"github.com/Dobryvechir/microcore/pkg/dvgrammar"
	"strconv"
)

var BracketProcessors = map[string]dvgrammar.BracketOperatorVisitor{
	"(":  ParentheseNoParentProcessor,
	"*(": ParentheseParentProcessor,
	"[":  SquareBracketNoParentProcessor,
	"*[": SquareBracketParentProcessor,
	"{":  CurlyBraceNoParentProcessor,
	"*{": CurlyBraceParentProcessor,
}

func SquareBracketNoParentProcessor(parent *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, rest []*dvgrammar.BuildNode) (value *dvgrammar.ExpressionValue, parentValue *dvgrammar.ExpressionValue, toStop bool, err error, noNextParent bool) {
	n := len(tree.Children)
	val := &DvVariable{Kind: FIELD_ARRAY, Fields: make([]*DvVariable, n)}
	for i := 0; i < n; i++ {
		t := tree.Children[i]
		if t != nil {
			_, r, err := t.ExecuteExpression(context)
			if err != nil {
				return nil, nil, false, err, false
			}
			if r == nil {
				continue
			}
			vl := AnyToDvVariable(r.Value)
			val.Fields[i] = vl
		}
	}
	value = &dvgrammar.ExpressionValue{Value: val, DataType: dvgrammar.TYPE_OBJECT}
	parentValue = parent
	return
}

func SquareBracketParentProcessor(parent *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, rest []*dvgrammar.BuildNode) (value *dvgrammar.ExpressionValue, parentValue *dvgrammar.ExpressionValue, toStop bool, err error, noNextParent bool) {
	n := len(tree.Children) - 1
	parentValue = parent
	for i := 0; i < n; i++ {
		t := tree.Children[i]
		if t != nil {
			_, _, err = t.ExecuteExpression(context)
			if err != nil {
				return
			}
		}
	}
	if n < 0 {
		return nil, nil, false, errors.New("No arguments in []"), false
	}
	node := tree.Children[n]
	node1, node2, ok := GetColumnSubNodes(node)
	if ok {
		_, value, err = node1.ExecuteExpression(context)
		if err != nil {
			return
		}
		value1, ok := AnyToNumberInt(value)
		if !ok {
			err = errors.New("First argument in [:] is not int: " + AnyToString(value))
		}
		_, value, err = node2.ExecuteExpression(context)
		if err != nil {
			return
		}
		value2, ok := AnyToNumberInt(value)
		if !ok {
			err = errors.New("Second argument in [:] is not int: " + AnyToString(value))
		}
		value, err = GetExpressionValueRange(parent, int(value1), int(value2))
		return
	}
	_, value, err = node.ExecuteExpression(context)
	if err != nil {
		return
	}
	value, err = GetExpressionValueChild(parent, value, context)
	return
}

func CurlyBraceNoParentProcessor(parent *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, rest []*dvgrammar.BuildNode) (value *dvgrammar.ExpressionValue, parentValue *dvgrammar.ExpressionValue, toStop bool, err error, noNextParent bool) {
	n := len(tree.Children)
	if n == 0 {
		d := &DvVariable{
			Kind:   FIELD_OBJECT,
			Fields: make([]*DvVariable, 0, 16),
		}
		value = AnyToDvGrammarExpressionValue(d)
		return
	}
	isObject := true
	for i := 0; i < n; i++ {
		if !isObjectLike(tree.Children[i]) {
			isObject = false
			break
		}
	}
	if isObject {
		d := &DvVariable{
			Kind:   FIELD_OBJECT,
			Fields: make([]*DvVariable, 0, n),
		}
		for i := 0; i < n; i++ {
			v, ok, err1 := ConvertToObjectKeyPair(tree.Children[i], context)
			if err1 != nil {
				err = err1
				return
			}
			if ok {
				d.Fields = append(d.Fields, v)
			}
		}
		value = AnyToDvGrammarExpressionValue(d)
		return
	}
	for i := 0; i < n; i++ {
		_, value, err = tree.Children[i].ExecuteExpression(context)
		if err != nil {
			break
		}
	}
	noNextParent = true
	parent = nil
	return
}

func isSimpleKey(node *dvgrammar.BuildNode) string {
	if node.Operator == "" && len(node.Children) == 0 && node.Value != nil {
		if node.Value.DataType == dvgrammar.TYPE_STRING || node.Value.DataType == dvgrammar.TYPE_DATA {
			return node.Value.Value
		}
	}
	return ""
}

func isObjectLike(node *dvgrammar.BuildNode) bool {
	if node.Operator == ":" && len(node.Children) == 2 {
		k := isSimpleKey(node.Children[0])
		if k != "" {
			return true
		}
	} else {
		k := isSimpleKey(node)
		if k != "" {
			return true
		}
	}
	return false
}

func ConvertToObjectKeyPair(node *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext) (*DvVariable, bool, error) {
	if node.Operator == ":" && len(node.Children) == 2 {
		k := isSimpleKey(node.Children[0])
		if k != "" {
			_, v, err := node.Children[1].ExecuteExpression(context)
			if err != nil {
				return nil, false, err
			}
			d := AnyToDvVariable(v)
			if d == nil {
				d = &DvVariable{Kind: FIELD_NULL}
			}
			d.Name = []byte(k)
			return d, true, nil
		}
	} else {
		k := isSimpleKey(node)
		if k != "" {
			v, ok := context.Scope.Get(k)
			if !ok {
				return nil, false, errors.New("Undefined " + k)
			}
			d := AnyToDvVariable(v)
			if d == nil {
				d = &DvVariable{Kind: FIELD_NULL}
			}
			d.Name = []byte(k)
			return d, true, nil
		}
	}
	return nil, false, nil
}

func CurlyBraceParentProcessor(parent *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, rest []*dvgrammar.BuildNode) (value *dvgrammar.ExpressionValue, parentValue *dvgrammar.ExpressionValue, toStop bool, err error, noNextParent bool) {
	return nil, nil, false, errors.New("Unexpected expression in {}"), true
}

func ParentheseNoParentProcessor(parent *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, rest []*dvgrammar.BuildNode) (value *dvgrammar.ExpressionValue, parentValue *dvgrammar.ExpressionValue, toStop bool, err error, noNextParent bool) {
	n := len(tree.Children)
	for i := 0; i < n; i++ {
		_, value, err = tree.Children[i].ExecuteExpression(context)
		if err != nil {
			break
		}
	}
	return
}

func ParentheseParentProcessor(parent *dvgrammar.ExpressionValue, tree *dvgrammar.BuildNode, context *dvgrammar.ExpressionContext, rest []*dvgrammar.BuildNode) (value *dvgrammar.ExpressionValue, parentValue *dvgrammar.ExpressionValue, toStop bool, err error, noNextParent bool) {
	if parent == nil || parent.Value == nil {
		return nil, nil, false, errors.New("Cannot execute function of null"), false
	}
	if dv, isDv := parent.Value.(*DvVariable); isDv && dv != nil && dv.Kind == FIELD_FUNCTION && dv.Extra != nil {
		if dvf, isDvf := dv.Extra.(*DvFunctionObject); isDvf && dvf != nil && dvf.Executor != nil && dvf.Executor.Special {
			value, toStop, err = dvf.Executor.FnSpecial(tree.Children, context, rest)
			return
		}
	}
	args, err := CalculateAllNodeParams(tree.Children, context)
	if err != nil {
		return
	}
	var val interface{}
	val, err = ExecuteAnyFunction(context, parent.Value, parent.Parent, args)
	if err == nil {
		value = AnyToDvGrammarExpressionValue(val)
		parentValue = parent
		return
	}
	return nil, nil, false, err, false
}

func GetExpressionValueRange(value *dvgrammar.ExpressionValue, indexFrom int, indexTo int) (*dvgrammar.ExpressionValue, error) {
	if value == nil {
		return nil, errors.New("Cannot get [:] from undefined")
	}
	switch value.DataType {
	case dvgrammar.TYPE_BOOLEAN, dvgrammar.TYPE_NULL:
		return nil, errors.New("Cannot get [:] from undefined")
	case dvgrammar.TYPE_NAN:
		return value, nil
	case dvgrammar.TYPE_OBJECT:
		v := AnyToDvVariable(value)
		if v == nil || v.Kind == FIELD_UNDEFINED {
			return nil, errors.New("Cannot get [:] from undefined")
		}
		res := v.GetChildrenByRange(indexFrom, indexTo-indexFrom)
		return &dvgrammar.ExpressionValue{
			Value:    res,
			DataType: dvgrammar.TYPE_OBJECT,
		}, nil
	case dvgrammar.TYPE_STRING:
		s := AnyToString(value.Value)
		res := ""
		if indexTo < 0 {
			indexTo = 0
		}
		if indexFrom > len(s) {
			indexFrom = len(s)
		}
		if indexTo > indexFrom && indexTo < len(s) {
			res = s[indexFrom:indexTo]
		}
		return &dvgrammar.ExpressionValue{Value: res, DataType: dvgrammar.TYPE_STRING}, nil
	}
	return nil, fmt.Errorf("Cannot apply [:] for %v", value)
}

func GetStringAtChar(s string, index int) *dvgrammar.ExpressionValue {
	n := len(s)
	if index < 0 || index >= n {
		return &dvgrammar.ExpressionValue{
			DataType: dvgrammar.TYPE_NULL,
			Value:    nil,
		}
	}
	return &dvgrammar.ExpressionValue{
		Value:    s[index : index+1],
		DataType: dvgrammar.TYPE_STRING,
	}
}

func GetExpressionValueChild(value *dvgrammar.ExpressionValue, index *dvgrammar.ExpressionValue, context *dvgrammar.ExpressionContext) (*dvgrammar.ExpressionValue, error) {
	if value == nil {
		return nil, errors.New("Cannot get child of undefined")
	}
	indexInt64, intOk := AnyToNumberInt(index)
	indexInt := int(indexInt64)
	if intOk {
		switch value.DataType {
		case dvgrammar.TYPE_STRING:
			return GetStringAtChar(AnyToString(value.Value), indexInt), nil
		case dvgrammar.TYPE_OBJECT:
			v := AnyToDvVariable(value.Value)
			if v == nil || v.Kind == FIELD_NULL {
				return nil, errors.New("Cannot get child of undefined")
			}
			if v.Kind == FIELD_STRING {
				return GetStringAtChar(string(v.Value), indexInt), nil
			}
			if v.Kind == FIELD_ARRAY {
				if indexInt < 0 || indexInt >= len(v.Fields) {
					return &dvgrammar.ExpressionValue{
						Value:    nil,
						DataType: dvgrammar.TYPE_NULL,
						Name: strconv.Itoa(indexInt),
					}, nil
				}
				fev:=v.Fields[indexInt].ToDvGrammarExpressionValue()
				if (context.VisitorOptions & dvgrammar.EVALUATE_OPTION_NAME)!=0 {
					fev.Name = strconv.Itoa(indexInt)
				}
				return fev, nil
			}
		}
	}
	child := AnyToString(index)
	v := AnyToDvVariable(value.Value)
	r := v.ReadSimpleChild(child)
	if r != nil {
		rev := r.ToDvGrammarExpressionValue()
		if r.Kind==FIELD_FUNCTION && r.Extra!=nil {
			dvf:=r.Extra.(*DvFunction)
			if dvf!=nil && dvf.Immediate {
				val, err := ExecuteAnyFunction(context, r, v, nil)
				if err!=nil {
					return nil, err
				}
				rev = AnyToDvGrammarExpressionValue(val)
			}
		}
		if rev != nil && (context.VisitorOptions&dvgrammar.EVALUATE_OPTION_NAME) != 0 {
			rev.Name = child
		}
		return rev, nil
	}
	fnMap := GetPrototypeForDvGrammarExpressionValue(value)
	vl, ok := fnMap.Get(child)
	if !ok {
		if (context.VisitorOptions & dvgrammar.EVALUATE_OPTION_NAME)!=0 {
			vlev:=&dvgrammar.ExpressionValue{
				DataType: dvgrammar.TYPE_NULL,
				Value: nil,
				Name: child,
			}
			return vlev, nil
		}
		return nil, nil
	}
	switch vl.(type) {
	case *DvFunction:
		return GetFunctionObjectVariable(vl.(*DvFunction), value, context)
	}
	gev := AnyToDvGrammarExpressionValue(vl)
	if gev != nil && (context.VisitorOptions&dvgrammar.EVALUATE_OPTION_NAME) != 0 {
		gev.Name = child
	}
	return gev, nil
}

func GetColumnSubNodes(node *dvgrammar.BuildNode) (*dvgrammar.BuildNode, *dvgrammar.BuildNode, bool) {
	if node.Operator == ":" && len(node.Children) == 2 {
		return node.Children[0], node.Children[1], true
	}
	return nil, nil, false
}
