package calc

import (
	"testing"

	"github.com/pbkheiron/sweinterview/moretesting"
)

func Test_Expr_Eval(t *testing.T) {
	// (20 * 5) / ((7 + 3) - 2) - parenthesis to clarify the expression tree structure.
	root := &ExprNode{
		Type: NodeTypeOpDiv,
		Operands: []*ExprNode{
			{
				Type: NodeTypeOpMul,
				Operands: []*ExprNode{
					NewNumExpr(20),
					NewNumExpr(5.0),
				},
			},
			{
				Type: NodeTypeOpSub,
				Operands: []*ExprNode{
					{
						Type: NodeTypeOpAdd,
						Operands: []*ExprNode{
							NewNumExpr(7),
							NewNumExpr(5),
						},
					},
					NewNumExpr(2),
				},
			},
		},
	}

	result, err := root.Eval()
	moretesting.AssertNoError(t, err)
	moretesting.AssertEqual(t, "invalid result", result, 10.0)
}
