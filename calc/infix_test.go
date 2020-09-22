package calc

import (
	"testing"

	"github.com/pbkheiron/sweinterview/moretesting"
)

func Test_InfixEval_OK(t *testing.T) {
	cases := []struct {
		expr     string
		expected float64
	}{
		{
			expr: "( 1 + 2 )",
			expected: 3,
		},
		{
			expr: "( 1 + ( 2 * 3 ) )",
			expected: 7,
		},
		{
			expr: " ( ( 1 * 2 ) + 3 )",
			expected: 5,
		},
		{
			expr: " ( ( ( 1 + 1 ) / 10 ) - ( 1 * 2 ) )",
			expected: -1.8,
		},
		{
			expr: "3",
			expected: 3,
		},
		{
			expr: "2 + 2",
			expected: 4,
		},
		{
			expr: "5 + 5 - 2",
			expected: 8,
		},
		{
			expr: "5 - 2 + 12 / 2 * 3",
			expected: 21,
		},
		{
			expr: "( 5 * 2 ) + 2",
			expected: 12,
		},
		{
			expr: "5 * ( 2 + 2 )",
			expected: 20,
		},
		{
			expr: "( 5 * 2 ) + ( 5 * 3 )",
			expected: 25,
		},
		{
			expr: "( ( 5 * 2 ) + ( 5 * 3 ) )",
			expected: 25,
		},
		{
			expr: "( 5 * ( 2 + 5 ) * 3 )",
			expected: 105,
		},
		{
			expr: "0 - 100",
			expected: -100,
		},
		{
			expr: "8 + ( 2 - 12 )",
			expected: -2,
		},
		{
			expr: "8 / ( 0 - 2 )",
			expected: -4,
		},
		{
			expr: "( ( 4 + 2 ) + ( 3 * ( 5 - 1 ) ) )",
			expected: 18,
		},
	}

	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			result, err := InfixEval(c.expr)
			moretesting.AssertNoError(t, err)
			moretesting.AssertEqual(t, "invalid result", c.expected, result)
		})
	}
}


func Test_InfixParser_ParseOperatorOK(t *testing.T) {
	cases := []struct {
		expr     string
		expected *ExprNode
	}{
		{
			expr: "1 + 2",
			expected: &ExprNode{
				Type:     NodeTypeOpAdd,
				Operands: []*ExprNode{NewNumExpr(1), NewNumExpr(2)},
			},
		},
		{
			expr: "1 + 2 * 3",
			expected: &ExprNode{
				Type: NodeTypeOpAdd,
				Operands: []*ExprNode{
					NewNumExpr(1),
					{
						Type:     NodeTypeOpMul,
						Operands: []*ExprNode{NewNumExpr(2), NewNumExpr(3)},
					},
				},
			},
		},
		{
			expr: "1 + 2 * 3 - 4",
			expected: &ExprNode{
				Type: NodeTypeOpSub,
				Operands: []*ExprNode{
					{
						Type: NodeTypeOpAdd,
						Operands: []*ExprNode{
							NewNumExpr(1),
							{
								Type:     NodeTypeOpMul,
								Operands: []*ExprNode{NewNumExpr(2), NewNumExpr(3)},
							},
						},
					},
					NewNumExpr(4),
				},
			},
		},
		{
			expr: "1 - 2 - 3",
			expected: &ExprNode{
				Type: NodeTypeOpSub,
				Operands: []*ExprNode{
					{
						Type:     NodeTypeOpSub,
						Operands: []*ExprNode{NewNumExpr(1), NewNumExpr(2)},
					},
					NewNumExpr(3),
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			parser := InfixParser{}
			result, err := parser.Parse(c.expr)
			moretesting.AssertNoError(t, err)
			moretesting.AssertEqual(t, "invalid expression tree", c.expected, result)
		})
	}
}

func Test_InfixParser_ParseParenOK(t *testing.T) {
	cases := []struct {
		expr     string
		expected *ExprNode
	}{
		{
			expr: "( 1 + 2 )",
			expected: &ExprNode{
				Type:     NodeTypeOpAdd,
				Operands: []*ExprNode{NewNumExpr(1), NewNumExpr(2)},
			},
		},
		{
			expr: "( 1 + 2 ) * 3",
			expected: &ExprNode{
				Type: NodeTypeOpMul,
				Operands: []*ExprNode{
					{
						Type:     NodeTypeOpAdd,
						Operands: []*ExprNode{NewNumExpr(1), NewNumExpr(2)},
					},
					NewNumExpr(3),
				},
			},
		},
		{
			expr: "2 * ( 3 - ( 4 + 5 ) )",
			expected: &ExprNode{
				Type: NodeTypeOpMul,
				Operands: []*ExprNode{
					NewNumExpr(2),
					{
						Type: NodeTypeOpSub,
						Operands: []*ExprNode{
							NewNumExpr(3),
							{
								Type:     NodeTypeOpAdd,
								Operands: []*ExprNode{NewNumExpr(4), NewNumExpr(5)},
							},
						},
					},
				},
			},
		},
		{
			expr: "1 - ( 2 - 3 )",
			expected: &ExprNode{
				Type: NodeTypeOpSub,
				Operands: []*ExprNode{
					NewNumExpr(1),
					{
						Type:     NodeTypeOpSub,
						Operands: []*ExprNode{NewNumExpr(2), NewNumExpr(3)},
					},

				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			parser := InfixParser{}
			result, err := parser.Parse(c.expr)
			moretesting.AssertNoError(t, err)
			moretesting.AssertEqual(t, "invalid expression tree", c.expected, result)
		})
	}
}

