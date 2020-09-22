package calc

import (
	"testing"

	"github.com/pbkheiron/sweinterview/moretesting"
)

func Test_PrefixEval_OK(t *testing.T) {
	cases := []struct {
		expr     string
		expected float64
	}{
		{
			expr:     "3",
			expected: 3,
		},
		{
			expr:     "+ 1 2",
			expected: 3,
		},
		{
			expr:     "+ 1 * 2 3",
			expected: 7,
		},
		{
			expr:     "+ * 1 2 3",
			expected: 5,
		},
		{
			expr:     "- / 10 + 1 1 * 1 2",
			expected: 3,
		},
		{
			expr:     "- 0 3",
			expected: -3,
		},
		{
			expr:     "/ 3 2",
			expected: 1.5,
		},
	}

	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			result, err := PrefixEval(c.expr)
			moretesting.AssertNoError(t, err)
			moretesting.AssertEqual(t, "invalid result", c.expected, result)
		})
	}
}
