package testpkg

import "testing"

func TestSomething(t *testing.T) {
	someConstant := " something"
	someOtherConstant := " something else"
	t.Run("test with spaces", func(t *testing.T) {
		t.Parallel()

		t.Run("deeply nested"+someConstant, func(t *testing.T) {
			t.Parallel()
		})
	})

	t.Run("another  t.Run()", func(t *testing.T) {
		t.Parallel()
	})

	t.Run(someOtherConstant, func(t *testing.T) {
		t.Parallel()
	})
}
