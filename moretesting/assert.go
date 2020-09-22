package moretesting

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func AssertEqual(t *testing.T, msg string, expected, actual interface{}) {
	t.Helper()
	err := CheckEqual(msg, expected, actual)
	if err != nil {
		t.Error(err.Error())
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err.Error())
	}
}

func CheckEqual(msg string, expected, actual interface{}) error {
	deep.CompareUnexportedFields = true
	diff := deep.Equal(expected, actual)
	if diff != nil {
		diffLines := strings.Join(diff, "\n\t")
		return fmt.Errorf("%s\n%s", msg, diffLines)
	} else {
		return nil
	}
}