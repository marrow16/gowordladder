package test

import "testing"

func AssertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Panic expected")
		}
	}()
	f()
}

func AssertTrue(t *testing.T, val bool) {
	t.Helper()
	if !val {
		t.Error("Expected true")
	}
}

func AssertFalse(t *testing.T, val bool) {
	t.Helper()
	if val {
		t.Error("Expected false")
	}
}

func AssertEqualsString(t *testing.T, expected string, actual string) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %s but actual is %s", expected, actual)
	}
}

func AssertEqualsInt(t *testing.T, expected int, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %d but actual is %d", expected, actual)
	}
}
