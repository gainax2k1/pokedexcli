package main

import (
	"fmt"
	"testing"
	"time"
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
func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}
