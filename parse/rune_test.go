package parse

import (
	"fmt"
	"testing"

	"github.com/a-h/lexical/input"
)

func TestRune(t *testing.T) {
	tests := []struct {
		input    string
		r        rune
		expected bool
	}{
		{
			input:    "A",
			r:        'A',
			expected: true,
		},
		{
			input:    "b",
			r:        'A',
			expected: false,
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(fmt.Sprintf("test %d", i), test.input)
		result := Rune(pi, test.r)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if test.expected && result.Item != test.r {
			t.Errorf("test %v: for input '%v' expected to capture '%v' but got '%v'", i, test.input, test.r, result.Item)
		}
	}
}
