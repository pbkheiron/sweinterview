package calc

import (
	"fmt"
	"strconv"
	"strings"
)

type Token struct {
	Raw string

	Number bool
	Operator bool
	Paren bool

	NumValue float64
}

func Tokenize(expr string) ([]Token, error) {
	rawTokens := strings.Fields(expr)
	tokens := make([]Token, 0, len(rawTokens))
	for _, rawToken := range rawTokens {
		switch {
		case rawToken == "+" || rawToken == "-" || rawToken == "*" || rawToken == "/":
			tokens = append(tokens, Token{
				Raw: rawToken,
				Operator: true,
			})
		case rawToken == "(" || rawToken == ")":
			tokens = append(tokens, Token{
				Raw: rawToken,
				Paren: true,
			})
		default:
			// assume number
			numValue, err := strconv.ParseFloat(rawToken, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse token as number %s: %s", rawToken, err.Error())
			}
			tokens = append(tokens, Token{
				Raw: rawToken,
				Number: true,
				NumValue: numValue,
			})
		}
	}
	return tokens, nil
}
