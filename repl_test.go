package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input 		string
		expected 	[]string
	}{
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Arremangalo arrempujaLo",
			expected: []string{"arremangalo", "arrempujalo"},
		},
	}

	for _, test_case := range cases {
		current_case := cleanInput(test_case.input);

		for i := range current_case {
			word := current_case[i];
			expected_word := test_case.expected[i];

			if word != expected_word {
				t.Errorf("Expected %s, but got %s", expected_word, word);
			}
		}
	}
}