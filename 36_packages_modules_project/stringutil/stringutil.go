// Package stringutil is a regular (non-internal) package - unlike
// mathutil, it does NOT live under an internal/ directory, so it could be
// imported by a totally different module too, not just this one. this
// file exists mainly to contrast with mathutil and round out the
// PROJECT STRUCTURE lesson: internal/ vs. a normal importable package.
package stringutil

import "strings"

// Reverse is exported - part of this package's public API.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Shout is exported and happens to use an unexported helper below.
func Shout(s string) string {
	return normalize(s) + "!!!"
}

// normalize is unexported - an internal detail of Shout, invisible and
// uncallable from outside this package.
func normalize(s string) string {
	return strings.ToUpper(strings.TrimSpace(s))
}
