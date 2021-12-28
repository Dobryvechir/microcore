/***********************************************************************
MicroCore
Copyright 2017 - 2020 by Danyil Dobryvechir (dobrivecher@yahoo.com ddobryvechir@gmail.com)
************************************************************************/
package dvgrammar

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

func buildExpressionTree(tokens []Token, opt *GrammarBaseDefinition) (*BuildNode, error) {
	currentPreAttributes := make([]string, 0, 16)
	tree := newNode(nil, opt, nil)
	current := tree
	amount := len(tokens)
	for i := 0; i < amount; i++ {
		value := &tokens[i]
		operator := value.Value
		if value.DataType != TYPE_OPERATOR && value.DataType != TYPE_CONTROL {
			operator = dataOperator
		}
		properties, isOperator := opt.Operators[operator]
		if !isOperator && (current.Children != nil || current.Value != nil) && operator != ")" {
			if opt.DefaultOperator != "" {
				operator = opt.DefaultOperator
				properties = opt.Operators[operator]
				isOperator = true
				i--
			} else {
				fullTreeClean(tree)
				return nil, errorMessage("No operator between values", value)
			}
		}
		if isOperator {
			if len(currentPreAttributes) != 0 {
				fullTreeClean(tree)
				return nil, errorMessage("Unexpected unary operator before "+operator, value)
			}
			if current.Value != nil {
				valueNode := &BuildNode{
					Parent: current,
					PreAttributes: current.PreAttributes,
					PostAttributes: current.PostAttributes,
					Value: current.Value,
				}
				current.Children = make([]*BuildNode, 1, 2)
				current.Children[0] = valueNode
				current.Value = nil
			}
			if current.Children == nil {
				if _,isUniOper:=opt.UnaryOperators[operator];isUniOper {
					isOperator = false
				} else {
					fullTreeClean(tree)
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
						current.Children = append(current.Children, node)
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
								fullTreeClean(tree)
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
		if !isOperator {
			modifier, okModifier := opt.UnaryOperators[operator]
			if okModifier {
				if modifier.Post {

				}
				if modifier.Pre {
					currentPreAttributes = append(currentPreAttributes, operator)
				}
			} else {
				switch operator {
				case "(":
					{
						node := newNode(current, opt, currentPreAttributes)
						currentPreAttributes = currentPreAttributes[:0]
						current.Children = make([]*BuildNode, 1, 2)
						current.Children[0] = node
						current.Operator = ")"
						current = node
					}
				case ")":
					for current.Operator != operator {
						if current.Parent == nil {
							fullTreeClean(tree)
							return nil, errorMessage("Unexpected )", value)
						}
						current = current.Parent
					}
					if current.Parent != nil {
						index := indexOfNode(current.Parent.Children, current)
						if index < 0 {
							fullTreeClean(tree)
							return nil, errorMessage("Broken tree", value)
						}
						node := current.Children[0]
						node.Parent = current.Parent
						current.Parent.Children[index] = node
						current = node
					} else {
						tree = current.Children[0]
						tree.Parent = nil
						current = tree
					}
					current.closed = true
					if current.Value != nil && current.Parent != nil && (current.Parent.Operator == "" || opt.Operators[current.Parent.Operator] != nil) {
						current = current.Parent
					}
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
					fullTreeClean(tree)
					return nil, errorMessage("Not allowed operator "+operator, value)
				}
			}
		}
	}
	for current != nil && (current.Operator == "" || opt.Operators[current.Operator] != nil) {
		current = current.Parent
	}
	if current != nil {
		fullTreeClean(tree)
		return nil, errorMessage("Unexpected end of expression", &tokens[len(tokens)-1])
	}
	if len(currentPreAttributes) != 0 {
		fullTreeClean(tree)
		return nil, errorMessage("Unexpected unary operator at the end of expression", &tokens[len(tokens)-1])
	}
	halfTreeClean(tree)
	return tree, nil
}
