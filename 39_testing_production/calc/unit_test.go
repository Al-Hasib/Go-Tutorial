package calc

import "testing"

// a test function's name must start with Test, take *testing.T, and live
// in a file ending in _test.go - `go test` finds and runs it automatically,
// no registration needed anywhere.
func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		//t.Errorf marks the test failed but keeps running the rest of this
		//function - use this by default, so one bad assertion doesn't hide
		//other useful failures further down
		t.Errorf("Add(2, 3) = %d, want %d", got, want)
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(10, 2)
	if err != nil {
		//t.Fatalf marks the test failed AND stops this function immediately -
		//use this when a later line would panic/be meaningless without this
		//check succeeding first (here: `got` only makes sense if err is nil)
		t.Fatalf("Divide(10, 2) returned unexpected error: %v", err)
	}
	if got != 5 {
		t.Errorf("Divide(10, 2) = %v, want 5", got)
	}
}

func TestDivide_ByZero(t *testing.T) {
	//testing that an error IS returned is just as important as testing the
	//happy path - a function that's supposed to fail should be tested failing
	_, err := Divide(10, 0)
	if err == nil {
		t.Error("Divide(10, 0) should have returned an error, got nil")
	}
}

func TestIsPalindrome(t *testing.T) {
	//t.Run creates a named SUBTEST - failures point at exactly which one
	//failed (e.g. "TestIsPalindrome/empty_string"), and `go test -run` can
	//target a single subtest by name
	t.Run("simple palindrome", func(t *testing.T) {
		if !IsPalindrome("racecar") {
			t.Error("expected \"racecar\" to be a palindrome")
		}
	})
	t.Run("not a palindrome", func(t *testing.T) {
		if IsPalindrome("hello") {
			t.Error("expected \"hello\" to NOT be a palindrome")
		}
	})
	t.Run("empty string", func(t *testing.T) {
		//an empty string trivially reads the same forwards and backwards
		if !IsPalindrome("") {
			t.Error("expected empty string to count as a palindrome")
		}
	})
}

// assertEqual is a small hand-rolled helper - testing has no built-in
// assertion library on purpose (see the philosophy note in table_driven_test.go)
func assertEqual(t *testing.T, got, want any) {
	t.Helper() // makes a failure here report the CALLER's line number, not this one
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestAdd_WithHelper(t *testing.T) {
	assertEqual(t, Add(1, 1), 2)
}
