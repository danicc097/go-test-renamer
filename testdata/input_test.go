package testpkg

import "testing"

func TestSomething(t *testing.T) {
	t.Run("test with spaces", func(t *testing.T) {
		t.Parallel()

		t.Run("deeply nested", func(t *testing.T) {
			t.Parallel()
		})
	})

	t.Run("another  t.Run()", func(t *testing.T) {
		t.Parallel()
	})
}
