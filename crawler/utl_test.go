package crawler

import (
	"testing"
)

var StencilTestCases = map[string][]string{
	"a  bc": []string{"a -- bc", "a xx bc"},
	"ab":    []string{"ab-d", "abc"},
	"abd":   []string{"ab-d", "abcdef"},
	"23.12 2017":   []string{"-----02.01 2006", "Sat. 23.12 2017 - Doors: 22:00"},
	"abc":   []string{"a-b-c", "axbxc"}}

func TestStencil(t *testing.T) {
	for expected, params := range StencilTestCases {
		if actual := Stencil(params[0], params[1]); actual != expected {
			fail(t, expected, actual)
		}
	}
}

var stripDashesTestCases = map[string]string{
	"a--b-c": "abc",
	"-abc":   "abc",
	"abc-":   "abc",
	"abc":    "abc"}

func TestStripDashes(t *testing.T) {
	for source, expected := range stripDashesTestCases {
		if actual := stripDashes(source); actual != expected {
			fail(t, expected, actual)
		}
	}
}

func fail(t *testing.T, expected, actual string) {
	t.Fatalf("Expected %q, not %q", expected, actual)
}
