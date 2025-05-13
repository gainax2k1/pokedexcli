package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {

	// start by creating a slice of structs:

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "this is a test",
			expected: []string{"this", "is", "a", "test"},
		},
		{
			input:    " ThIs Is A cAmEl TeSt ",
			expected: []string{"this", "is", "a", "camel", "test"},
		},
		{
			input:    " tOO     MuCh    whiTESPace     ",
			expected: []string{"too", "much", "whitespace"},
		},

		// add more cases here

	}

	// then loop over cases and run tests

	for _, c := range cases {
		actual := cleanInput(c.input)
		// check the actual length vs expected
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q) return %d words, expected %d", c.input, len(actual), len(c.expected))
			continue
		}
		// if not a match, use t.Errorf to print msg and fail test

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// check each word in slice
			// if no match, uses t.Errorf.....
			if word != expectedWord {
				t.Errorf("cleanInput(%q)[%d] = %q, expected %q", c.input, i, word, expectedWord)
			}
		}
	}

}
