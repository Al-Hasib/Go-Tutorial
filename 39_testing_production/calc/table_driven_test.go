package calc

import "testing"

// table-driven tests exist because most of the tests above are really the
// SAME test, run with different numbers - Go doesn't have a built-in
// assertion/parameterization library, so the idiomatic answer is: build a
// slice of cases, then loop over it with t.Run.
func TestDivideTable(t *testing.T) {
	cases := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr bool
	}{
		{name: "positive numbers", a: 10, b: 2, want: 5},
		{name: "negative dividend", a: -10, b: 2, want: -5},
		{name: "divides to a fraction", a: 1, b: 4, want: 0.25},
		{name: "division by zero errors", a: 10, b: 0, wantErr: true},
	}

	for _, tc := range cases {
		//capturing tc as a parameter isn't strictly needed for closures
		//since Go 1.22 (see the Goroutines lesson), but t.Run's own name
		//parameter still needs tc.name read at the right time either way
		t.Run(tc.name, func(t *testing.T) {
			got, err := Divide(tc.a, tc.b)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Divide(%v, %v) should have errored, got %v", tc.a, tc.b, got)
				}
				return // nothing else to check on the error path
			}

			if err != nil {
				t.Fatalf("Divide(%v, %v) returned unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.want {
				t.Errorf("Divide(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestIsPalindromeTable(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  bool
	}{
		{"palindrome", "racecar", true},
		{"not a palindrome", "hello", false},
		{"single character", "a", true},
		{"empty string", "", true},
		{"two identical characters", "aa", true},
		{"two different characters", "ab", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsPalindrome(tc.input); got != tc.want {
				t.Errorf("IsPalindrome(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

// adding one more case to either table above means adding one line, not
// writing (and remembering to write) a whole new Test function - that's
// the entire point of this pattern.
