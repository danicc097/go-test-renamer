package testpkg

import "testing"

func TestSomething(t *testing.T) {
	t.Run("test_with_spaces", func(t *testing.T) {
		t.Parallel()

		t.Run("deeply_nested", func(t *testing.T) {
			t.Parallel()
		})
	})

	t.Run("another__t.Run()", func(t *testing.T) {
		t.Parallel()
	})
}
