package calc

import "fmt"

type NodeType string

const (
	NodeTypeNum NodeType = "num"
	NodeTypeOpAdd NodeType = "+"
	NodeTypeOpSub NodeType = "-"
	NodeTypeOpMul NodeType = "*"
	NodeTypeOpDiv NodeType = "/"
)

func (t NodeType) Operator() bool {
	return t == NodeTypeOpAdd || t == NodeTypeOpSub || t == NodeTypeOpMul || t == NodeTypeOpDiv
}


type ExprNode struct {
	Type NodeType
	// Number is the expression value if NodeType == NodeTypeNum
	Number float64
	// Operands if Type is one of OpXYZ
	Operands []*ExprNode
}

func NewNumExpr(value float64) *ExprNode{
	return &ExprNode{Type: NodeTypeNum, Number: value}
}

func (n *ExprNode) Eval() (float64, error) {
	if n.Type == NodeTypeNum {
		return n.Number, nil
	}

	op1, op2, err := n.evalOperands()
	if err != nil {
		return 0, err
	}
	switch n.Type {
	case NodeTypeOpAdd:
		return op1 + op2, nil
	case NodeTypeOpSub:
		return op1 - op2, nil
	case NodeTypeOpMul:
		return op1 * op2, nil
	case NodeTypeOpDiv:
		if op2 == 0.0 {
			return 0, fmt.Errorf("division by zero")
		}
		return op1 / op2, nil
	default:
		return 0, fmt.Errorf("unsupported type: %v", n.Type)
	}
}

func (n *ExprNode) evalOperands() (float64, float64, error) {
	if len(n.Operands) != 2 {
		return 0, 0, fmt.Errorf("operator %v: unexpected number of operands: %d", n.Type, len(n.Operands))
	}
	op1, err := n.Operands[0].Eval()
	if err != nil {
		return 0, 0, fmt.Errorf("operator %v: failed to eval operand 1: %s", n.Type, err.Error())
	}
	op2, err := n.Operands[1].Eval()
	if err != nil {
		return 0, 0, fmt.Errorf("operator %v: failed to eval operand 2: %s", n.Type, err.Error())
	}
	return op1, op2, nil
}