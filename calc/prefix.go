package calc

import "fmt"

// PrefixEval evaluates expression written in prefix notation.
//
// NOTE: the implementation could use single stack for evaluation.
// As I already have ExprNode implemented, it was easier to parse the expression as tree and evaluate it.
func PrefixEval(s string) (float64, error) {
	exprRoot, err := PrefixParser{}.Parse(s)
	if err != nil {
		return 0, err
	}
	return exprRoot.Eval()
}


type PrefixParser struct {}

func (PrefixParser) Parse(s string) (*ExprNode, error) {
	tokens, err := Tokenize(s)
	if err != nil {
		return nil, err
	}

	exprStack := make([]*ExprNode, 0, len(tokens))
	for i := len(tokens)-1; i >= 0; i-- {
		token := tokens[i]
		switch {
		case token.Number:
			exprStack = append(exprStack, NewNumExpr(token.NumValue))
		case token.Operator:
			if len(exprStack) < 2 {
				return nil, fmt.Errorf("not enough expression on stack to build binary expression (op = %s)", token.Raw)
			}
			n := len(exprStack)
			op1, op2 := exprStack[n-1], exprStack[n-2]
			exprStack = exprStack[:n-2]

			exprType := NodeType(token.Raw)
			if !exprType.Operator() {
				return nil, fmt.Errorf("invalid operator: %s", token.Raw)
			}

			exprStack = append(exprStack, &ExprNode{
				Type:     exprType,
				Operands: []*ExprNode{op1, op2},
			})
		default:
			return nil, fmt.Errorf("unsupported token: %s", token.Raw)
		}
	}

	if len(exprStack) != 1 {
		return nil, fmt.Errorf("invalid expression string: evaluated to %d final expressions", len(exprStack))
	}
	return exprStack[0], nil
}
