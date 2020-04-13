package match

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOperator(t *testing.T) {
	testCase := map[string]struct {
		op      string
		value   string
		express string
		result  bool
	}{
		"1": {
			"noEqu",
			"1",
			"2",
			true,
		},
		"2": {
			"noEqu",
			"2",
			"2",
			false,
		},
		"3": {
			"noLess",
			"3",
			"2",
			true,
		},
		"4": {
			"noLess",
			"3",
			"3",
			true,
		},
		"5": {
			"less",
			"3",
			"4",
			true,
		},
		"6": {
			"less",
			"3",
			"3",
			false,
		},
		"7": {
			"noGreater",
			"3",
			"3",
			true,
		},
		"8": {
			"noGreater",
			"2",
			"3",
			true,
		},
		"9": {
			"greater",
			"3",
			"3",
			false,
		},
		"10": {
			"greater",
			"3",
			"2",
			true,
		},
	}

	for _, tc := range testCase {
		f, err := operatorPlugin[tc.op]
		assert.NotNil(t, err)
		assert.Equal(t, tc.result, f(tc.value, tc.express),
			fmt.Sprintf("test value %s op %s exp %s faile", tc.value, tc.op, tc.express))
	}

}
