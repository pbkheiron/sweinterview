package calc

import (
	"fmt"
)

func InfixEval(s string) (float64, error) {
	exprRoot, err := InfixParser{}.Parse(s)
	if err != nil {
		return 0, err
	}
	return exprRoot.Eval()
}


type InfixParser struct {}

type operator struct {
	op string
	priority int
}

var opPriorities = map[string]int{
	"+": 2,
	"-": 2,
	"*": 3,
	"/": 3,
	"(": 1,
}

func newOperator(op string) (operator, error) {
	priority, ok := opPriorities[op]
	if !ok {
		return operator{}, fmt.Errorf("no priority defined for operator: %s", op)
	}
	return operator{op: op, priority: priority}, nil
}

type infixContext struct {
	exprQueue []*ExprNode
	operatorStack []operator
}

func (p InfixParser) Parse(expr string) (*ExprNode, error) {
	tokens, err := Tokenize(expr)
	if err != nil {
		return nil, err
	}

	// new exprQueue with enough memory pre-allocated to store all tokens as ExprNodes.
	// A bit of overstretch, but good enough for first version.
	ctx := &infixContext{
		exprQueue: make([]*ExprNode, 0, len(tokens)),
		operatorStack: make([]operator, 0, len(tokens)/2),
	}

	for _, token := range tokens {
		switch {
		case token.Number:
			// TODO: here we can implement additional check, if there is enough operands in the queue
			// (to avoid parsing '+ * 1 2 3' as infix without error.
			ctx.appendToExprQueue(NewNumExpr(token.NumValue))
		case token.Operator:
			op, err := newOperator(token.Raw)
			if err != nil {
				return nil, err
			}
			for ctx.opStackLen() > 0 && ctx.opStackPeek().priority >= op.priority {
				err := ctx.popTopOperatorIntoExprQueue()
				if err != nil {
					return nil, err
				}
			}
			ctx.opPush(op)

		case token.Paren && token.Raw == "(":
			op, err := newOperator(token.Raw)
			if err != nil {
				return nil, err
			}
			ctx.opPush(op)
		case token.Paren && token.Raw == ")":
			for ctx.opStackLen() > 0 && ctx.opStackPeek().op != "(" {
				err := ctx.popTopOperatorIntoExprQueue()
				if err != nil {
					return nil, err
				}
			}
			if ctx.opStackLen() == 0 {
				return nil, fmt.Errorf("unmatched right parenthesis )")
			}
			ctx.opStackPop()
		}
	}

	// finalize - pop all remaining operator from stack and convert them into expressions
	for ctx.opStackLen() > 0 {
		err := ctx.popTopOperatorIntoExprQueue()
		if err != nil {
			return nil, err
		}
	}
	if len(ctx.exprQueue) != 1 {
		return nil, fmt.Errorf("invalid expression string: evaluated to %d final expressions", len(ctx.exprQueue))
	}
	return ctx.exprQueue[0], nil
}

func (c *infixContext) exprQueueLen() int {
	return len(c.exprQueue)
}

func (c *infixContext) appendToExprQueue(expr *ExprNode) {
	c.exprQueue = append(c.exprQueue, expr)
}

func (c *infixContext) popTopOperatorIntoExprQueue() error {
	if len(c.exprQueue) < 2 {
		return fmt.Errorf("not enough expression in queue to build binary expression")
	}

	op := c.opStackPop()
	exprType := NodeType(op.op)
	if !exprType.Operator() {
		return fmt.Errorf("invalid operator: %s", op.op)
	}

	exprN := len(c.exprQueue)
	expr1, expr2 := c.exprQueue[exprN-2], c.exprQueue[exprN-1]
	c.exprQueue = c.exprQueue[:exprN-2]

	opExpr := &ExprNode{
		Type: exprType,
		Operands: []*ExprNode{expr1, expr2},
	}
	c.appendToExprQueue(opExpr)
	return nil
}

func (c *infixContext) opStackLen() int {
	return len(c.operatorStack)
}

func (c *infixContext) opStackPeek() operator {
	return c.operatorStack[c.opStackLen()-1]
}

func (c *infixContext) opPush(op operator) {
	c.operatorStack = append(c.operatorStack, op)
}

func (c *infixContext) opStackPop() operator {
	n := c.opStackLen()
	op := c.operatorStack[c.opStackLen()-1]
	c.operatorStack = c.operatorStack[:n-1]
	return op
}
