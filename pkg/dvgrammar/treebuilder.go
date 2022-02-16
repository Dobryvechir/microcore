/***********************************************************************
MicroCore
Copyright 2017 - 2021 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

import "errors"

const dataOperator string = "DATA"

func invertStrings(data []string) []string {
	n := len(data)
	if n == 0 {
		return nil
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[n-i-1] = data[i]
	}
	return res
}

func newNode(parent *BuildNode, opt *GrammarBaseDefinition, preAttributes []string) *BuildNode {
	attr := invertStrings(preAttributes)
	return &BuildNode{Parent: parent, PreAttributes: attr}
}

func fullTreeClean(tree *BuildNode) {
	if tree != nil {
		tree.Parent = nil
		tree.Value = nil
		l := len(tree.Children)
		for i := 0; i < l; i++ {
			fullTreeClean(tree.Children[i])
			tree.Children[i] = nil
		}
	}
}

func fullTreeForestClean(forest []*BuildNode, tree *BuildNode) {
	n := len(forest)
	for i := 0; i < n; i++ {
		fullTreeClean(forest[i])
	}
	fullTreeClean(tree)
}

func halfTreeClean(tree *BuildNode) {
	if tree != nil {
		tree.Parent = nil
		l := len(tree.Children)
		for i := 0; i < l; i++ {
			halfTreeClean(tree.Children[i])
		}
	}
}

func indexOfNode(nodes []*BuildNode, node *BuildNode) int {
	for i, n := range nodes {
		if n == node {
			return i
		}
	}
	return -1
}

func findClosingTag(tokens []Token, pos int) (int, error) {
	n := len(tokens)
	var buf []byte
	tag := tokens[pos].Value[0]
	var closing byte
	switch tag {
	case '(':
		closing = ')'
	case '[':
		closing = ']'
	case '{':
		closing = '}'
	default:
		return -1, errors.New("Invalid bracket " + tokens[pos].Value)
	}
	var nextClosing byte
	var stackLen = 0
	for pos++; pos < n; pos++ {
		if len(tokens[pos].Value) == 0 {
			continue
		}
		tag = tokens[pos].Value[0]
		if tag == '(' {
			nextClosing = ')'
		} else if tag == '[' {
			nextClosing = ']'
		} else if tag == '{' {
			nextClosing = '}'
		} else if tag == ')' || tag == ']' || tag == '}' {
			if stackLen == 0 {
				if tag == closing {
					return pos, nil
				}
				return -1, errors.New("Expected " + string([]byte{closing}) + " but found " + string([]byte{tag}))
			}
			stackLen--
			if tag != buf[stackLen] {
				return -1, errors.New("Expected " + string([]byte{buf[stackLen]}) + " but found " + string([]byte{tag}))
			}
		} else {
			continue
		}
		if stackLen == 0 {
			buf = make([]byte, 1, n-pos)
			buf[0] = nextClosing
		} else if stackLen < len(buf) {
			buf[stackLen] = nextClosing
		} else {
			buf = append(buf, nextClosing)
		}
		stackLen++
	}
	return -1, errors.New("No closing tag " + string([]byte{closing}))
}

func buildExpressionTree(tokens []Token, opt *GrammarBaseDefinition) (forest []*BuildNode, err error) {
	forest = make([]*BuildNode, 0, 16)
	currentPreAttributes := make([]string, 0, 16)
	tree := newNode(nil, opt, nil)
	current := tree
	amount := len(tokens)
	group := 0
tokenRunner:
	for i := 0; i < amount; i++ {
		value := &tokens[i]
		operator := value.Value
		if value.DataType == TYPE_CONTROL {
			switch operator {
			case ";", ",":
				forest, err = placeTreeToForest(forest, current, tree, tokens, opt, currentPreAttributes, group)
				if err != nil {
					return
				}
				if operator == ";" {
					group++
				}
				tree = newNode(nil, opt, nil)
				current = tree
				continue tokenRunner
			case ".":
				holdDot := current
				for holdDot.Operator != "" && len(holdDot.Children) > 0 {
					holdDot = holdDot.Children[len(holdDot.Children)-1]
				}
				if holdDot.Value == nil {
					return nil, errors.New("Unexpected dot without previous variable")
				}
				i++
				if i == amount {
					return nil, errors.New("Unexpected dot without following name")
				}
				value = &tokens[i]
				operator = value.Value
				if operator == "" || value.DataType == TYPE_CONTROL || value.DataType == TYPE_OPERATOR || value.DataType == TYPE_OBJECT {
					return nil, errors.New("Dot must be followed by name")
				}
				subForest := make([]*BuildNode, 1)
				subForest[0] = &BuildNode{
					Value: &Token{DataType: TYPE_STRING, Value: operator},
				}
				node := &BuildNode{
					Children: subForest,
					Operator: "[",
				}
				holdDot.Children = append(holdDot.Children, node)
				continue tokenRunner
			case "(", "[", "{":
				pos, err := findClosingTag(tokens, i)
				if err != nil {
					fullTreeForestClean(forest, tree)
					return nil, err
				}
				var subForest []*BuildNode
				if i+1 < pos {
					subForest, err = buildExpressionTree(tokens[i+1:pos], opt)
					if err != nil {
						fullTreeForestClean(forest, tree)
						return nil, err
					}
				}
				i = pos
				node := &BuildNode{
					Children:      subForest,
					Operator:      operator,
					PreAttributes: invertStrings(currentPreAttributes),
				}
				currentPreAttributes = currentPreAttributes[:0]
				holderNode := current
				for holderNode.Operator != "" {
					m := len(holderNode.Children)
					if m == 0 {
						return nil, errors.New("Unexpected no place for " + operator)
					}
					holderNode = holderNode.Children[m-1]
				}
				holderNode.Children = append(holderNode.Children, node)
				if holderNode.Value == nil {
					holderNode.Value = &Token{DataType: TYPE_FUNCTION}
				}
				if current.Operator == "" && current.Parent != nil && (current.Parent.Operator == "" || opt.Operators[current.Parent.Operator] != nil) {
					current = current.Parent
				}
				continue tokenRunner
			}
		} else if value.DataType != TYPE_OPERATOR {
			operator = dataOperator
		}
		properties, isOperator := opt.Operators[operator]
		if !isOperator && (current.Children != nil || current.Value != nil) {
			modifier, okModifier := opt.UnaryOperators[operator]
			if okModifier && modifier.Post {
				node := current
				for node.Operator != "" && len(node.Children) > 0 {
					node = node.Children[len(node.Children)-1]
				}
				node.PostAttributes = append(node.PostAttributes, operator)
				continue
			} else if opt.DefaultOperator != "" {
				operator = opt.DefaultOperator
				properties = opt.Operators[operator]
				isOperator = true
				i--
			} else {
				fullTreeForestClean(forest, tree)
				return nil, errorMessage("No operator between values", value)
			}
		}
		if isOperator {
			if len(currentPreAttributes) != 0 {
				if _, isUniOper := opt.UnaryOperators[operator]; isUniOper {
					isOperator = false
				} else {
					fullTreeForestClean(forest, tree)
					return nil, errorMessage("Unexpected unary operator before "+operator, value)
				}
			} else {
				if current.Value != nil {
					valueNode := &BuildNode{
						Parent:         current,
						PreAttributes:  current.PreAttributes,
						PostAttributes: current.PostAttributes,
						Value:          current.Value,
						Children:       current.Children,
					}
					current.Children = make([]*BuildNode, 1, 2)
					current.Children[0] = valueNode
					current.Value = nil
					current.PreAttributes = nil
					current.PostAttributes = nil
				}
				if current.Children == nil {
					if _, isUniOper := opt.UnaryOperators[operator]; isUniOper {
						isOperator = false
					} else {
						fullTreeForestClean(forest, tree)
						return nil, errorMessage("Unexpected "+operator, value)
					}
				} else {
					node := newNode(current, opt, nil)
					if current.Operator == "" {
						current.Operator = operator
						current.Children = append(current.Children, node)
						current = node
					} else {
						precedence := properties.Precedence
						for current.Parent != nil && opt.Operators[current.Parent.Operator] != nil &&
							(current.closed || opt.Operators[current.Operator].Precedence > precedence) {
							current = current.Parent
						}
						if current.Operator == operator && properties.Multi {
							current.Children = append(current.Children, node)
							current = node
						} else if !current.closed && opt.Operators[current.Operator].Precedence < precedence {
							//reattach the node
							lastIndex := len(current.Children) - 1
							lastNode := current.Children[lastIndex]
							node.Children = make([]*BuildNode, 2, 3)
							node.Children[0] = lastNode
							node.Operator = operator
							lastNode.Parent = node
							current.Children[lastIndex] = node
							lastNode = newNode(node, opt, currentPreAttributes)
							currentPreAttributes = currentPreAttributes[:0]
							node.Children[1] = lastNode
							current = lastNode
						} else {
							//attach at up position
							node.Parent = current.Parent
							node.Children = make([]*BuildNode, 1, 2)
							node.Children[0] = current
							node.Operator = operator
							if node.Parent == nil {
								tree = node
							} else {
								index := indexOfNode(node.Parent.Children, current)
								if index < 0 {
									fullTreeForestClean(forest, tree)
									return nil, errorMessage("Tree broken", value)
								}
								node.Parent.Children[index] = node
							}
							current.Parent = node
							current = newNode(node, opt, nil)
							node.Children = append(node.Children, current)
						}
					}
				}
			}
		}
		if !isOperator {
			modifier, okModifier := opt.UnaryOperators[operator]
			if okModifier && modifier.Pre {
				currentPreAttributes = append(currentPreAttributes, operator)
			} else {
				switch operator {
				case dataOperator:
					current.Value = value
					if len(currentPreAttributes) != 0 {
						current.PreAttributes = invertStrings(currentPreAttributes)
						currentPreAttributes = currentPreAttributes[:0]
					}
					if current.Parent != nil && (current.Parent.Operator == "" || opt.Operators[current.Parent.Operator] != nil) {
						current = current.Parent
					}
				default:
					fullTreeForestClean(forest, tree)
					return nil, errorMessage("Not allowed operator "+operator, value)
				}
			}
		}
	}
	return placeTreeToForest(forest, current, tree, tokens, opt, currentPreAttributes, group)
}

func placeTreeToForest(forest []*BuildNode, current *BuildNode, tree *BuildNode,
	tokens []Token, opt *GrammarBaseDefinition, currentPreAttributes []string,
	group int) ([]*BuildNode, error) {
	for current != nil && (current.Operator == "" || opt.Operators[current.Operator] != nil) {
		current = current.Parent
	}
	if current != nil {
		fullTreeForestClean(forest, tree)
		return nil, errorMessage("Unexpected end of expression", &tokens[len(tokens)-1])
	}
	if len(currentPreAttributes) != 0 {
		fullTreeForestClean(forest, tree)
		return nil, errorMessage("Unexpected unary operator at the end of expression", &tokens[len(tokens)-1])
	}
	halfTreeClean(tree)
	tree.Group = group
	forest = append(forest, tree)
	return forest, nil
}

func (b *BuildNode) CloneFrom(other *BuildNode) {
	if b != nil {
		if other == nil {
			b.Value = nil
			b.Operator = ""
			b.Children = nil
			b.PreAttributes = nil
			b.PostAttributes = nil
			b.Parent = nil
			b.Group = 0
			b.closed = false
		} else {
			b.Value = other.Value
			b.Operator = other.Operator
			b.Children = other.Children
			b.PreAttributes = other.PreAttributes
			b.PostAttributes = other.PostAttributes
			b.Parent = other.Parent
			b.Group = other.Group
			b.closed = other.closed
		}
	}
}
