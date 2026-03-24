package repl

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "           ",
			expected: []string{},
		},
		{
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   hello",
			expected: []string{"hello"},
		},
		{
			input:    "  he ll oo wooo",
			expected: []string{"he", "ll", "oo", "wooo"},
		},
		{
			input:    "HELLO WORLD   ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Fatalf("expected %d words, got %d: %#v", len(c.expected), len(actual), actual)
		}

		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Fatalf("at index %d: expected %q, got %q", i, c.expected[i], actual[i])
			}
		}
	}
}
