package testdata

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

	type testCase struct {
		name string
	}

	tests := []testCase{
		{
			name: "a b c",
		},
		{
			name: "1 2 3",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
