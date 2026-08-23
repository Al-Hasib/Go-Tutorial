// Package mathutil demonstrates PACKAGES and EXPORTED VS UNEXPORTED names.
//
// this doc comment right above "package mathutil" is a real convention -
// `go doc mathutil` (or hovering in an editor) shows exactly this text.
//
// mathutil lives under internal/, which is a special directory name the Go
// tool itself understands: any package under a path containing "internal/"
// can ONLY be imported by code inside the same module. a completely
// different module could import stringutil (a sibling package, not under
// internal/) but importing example.com/packagesdemo/internal/mathutil from
// outside this module would fail to even compile - the compiler enforces
// this, it's not just a naming convention.
package mathutil

// ---------- EXPORTED: starts with a Capital letter ----------

// Add is EXPORTED (capital A) - any other package that imports mathutil
// can call mathutil.Add(...). exported names are this package's public API.
func Add(a, b int) int {
	return a + roundToEven(b) // uses an unexported helper internally
}

// Multiply is also exported.
func Multiply(a, b int) int {
	return a * b
}

// Config is an exported struct - but notice its fields below are a mix of
// exported and unexported.
type Config struct {
	Precision int    // exported field - visible and settable from other packages
	label     string // unexported field - invisible outside this package entirely
}

// NewConfig is exported. it's the only way another package can set `label`,
// since the field itself is unexported - this is how Go does encapsulation:
// no "private" keyword, just capitalization, at the package level.
func NewConfig(precision int, label string) Config {
	return Config{Precision: precision, label: label}
}

// Label is an exported method that reads the unexported field for us -
// a "getter", the idiomatic way to expose an unexported field selectively.
func (c Config) Label() string {
	return c.label
}

// ---------- UNEXPORTED: starts with a lowercase letter ----------

// roundToEven is unexported (lowercase r) - it's an implementation detail
// of this package. code outside mathutil cannot call mathutil.roundToEven(...)
// at all; the compiler won't even let you reference the name.
func roundToEven(n int) int {
	if n%2 != 0 {
		n++
	}
	return n
}

// unexported names aren't a security boundary (anyone can read this source
// file), they're a DESIGN boundary: they mark what this package promises to
// keep working across versions (exported) vs. what it's free to change
// without warning (unexported).
