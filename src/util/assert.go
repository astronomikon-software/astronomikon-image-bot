package util

import "testing"

func Assert[T comparable](t *testing.T, expected T, value T) {
	if expected != value {
		t.Errorf("Assertation failed, expected: %v, got: %v", expected, value)
	}
}
