package prometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveCharacters(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{input: "__under_score__", output: "under_score__"},
		{input: "c,omma", output: "comma"},
		{input: "test#it#", output: "testit"},
		{input: "white space", output: "whitespace"},
		{input: "normal_name", output: "normal_name"},
	}

	specialChars := "#, "

	for _, test := range tests {
		o := removeSpecialCharacters(test.input, specialChars)
		assert.Equal(t, test.output, o)
	}
}
