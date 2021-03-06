package scanner

import (
	"io"
	"reflect"
	"testing"
	"unicode"

	"github.com/a-h/lexical/input"
	"github.com/a-h/lexical/parse"
)

func TestScanning(t *testing.T) {
	stream := input.NewFromString(`<a>Example</a>`)

	scanner := New(stream, xmlTokens)
	var err error
	for {
		_, err := scanner.Next()
		if err != nil {
			break
		}
	}
	if err != nil && err != io.EOF {
		t.Error(err)
	}
}

var xmlTokens = parse.Any(xmlOpenElement, xmlText, xmlCloseElement)

var combineTagAndContents parse.MultipleResultCombiner = func(results []interface{}) (interface{}, bool) {
	name, ok := results[1].(string)
	return "name: " + name, ok
}

var letterOrDigit = parse.RuneInRanges(unicode.Letter, unicode.Number)

var xmlName = parse.Then(
	parse.WithStringConcatCombiner,
	parse.RuneWhere(unicode.IsLetter),
	parse.Many(parse.WithStringConcatCombiner,
		0,   // minimum match count
		500, // maxmum match count
		letterOrDigit),
)

var xmlOpenElement = parse.All(
	combineTagAndContents,
	parse.Rune('<'),
	xmlName,
	parse.Rune('>'),
)

var xmlText = parse.StringUntil(parse.Rune('<'))

var xmlCloseElement = parse.All(
	parse.WithStringConcatCombiner,
	parse.Rune('<'),
	parse.Rune('/'),
	parse.StringUntil(parse.Rune('>')),
)

func TestXMLName(t *testing.T) {
	tests := []struct {
		input        string
		expected     bool
		expectedItem string
	}{
		{
			input:        "AB",
			expected:     true,
			expectedItem: "AB",
		},
	}

	for i, test := range tests {
		pi := input.NewFromString(test.input)
		result := xmlName(pi)
		actual := result.Success
		if actual != test.expected {
			t.Errorf("test %v: for input '%v' expected %v but got %v", i, test.input, test.expected, actual)
		}
		if test.expected && result.Item != test.expectedItem {
			t.Errorf("test %v: for input '%v' expected item '%v' but got '%v' (%v)", i, test.input, test.expectedItem, result.Item, reflect.TypeOf(result.Item))
		}
	}
}
