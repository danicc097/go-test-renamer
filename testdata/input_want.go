package testpkg

import "testing"

func TestSomething(t *testing.T) {
	someConstant := " something"
	someOtherConstant := " something else"
	t.Run("test_with_spaces", func(t *testing.T) {
		t.Parallel()

		t.Run("deeply nested"+someConstant, func(t *testing.T) {
			t.Parallel()
		})
	})

	t.Run("another__t.Run()", func(t *testing.T) {
		t.Parallel()
	})

	t.Run(someOtherConstant, func(t *testing.T) {
		t.Parallel()
	})
}
